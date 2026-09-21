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
	"strings"
	"syscall"
	"time"

	"github.com/cashtro/cashtro/internal/agents"
	"github.com/cashtro/cashtro/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	if err := run(listenAddr(*addr)); err != nil {
		log.Fatal(err)
	}
}

// listenAddr prefers Azure App Service PORT / WEBSITES_PORT over the flag default.
func listenAddr(flagAddr string) string {
	for _, key := range []string{"PORT", "WEBSITES_PORT"} {
		if p := strings.TrimSpace(os.Getenv(key)); p != "" {
			if strings.HasPrefix(p, ":") {
				return p
			}
			return ":" + p
		}
	}
	if flagAddr == "" {
		return ":8080"
	}
	return flagAddr
}

func run(addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	k, err := agents.Boot()
	if err != nil {
		return err
	}
	about := k.About()
	log.Printf("%s %s · %d agentics online", about.Name, about.Version, about.Running)

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
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
