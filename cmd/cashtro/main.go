package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cashtro/cashtro/internal/agents"
	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	data := flag.String("data", envOr("CASHTRO_DATA", "data/cashtro.json"), "OS disk image")
	flag.Parse()

	if err := run(*addr, *data); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func run(addr, data string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	k, err := agents.Boot(kernel.WithPersistPath(data))
	if err != nil {
		return err
	}
	about := k.About()
	log.Printf("%s %s · %d agentics online · image %s", about.Name, about.Version, about.Running, data)

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           server.New(k),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("cashtro listening on http://%s", ln.Addr())

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpSrv.Serve(ln)
	}()

	select {
	case <-ctx.Done():
		k.Close()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		k.Close()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
