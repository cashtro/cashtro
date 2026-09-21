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

	"github.com/cashtro/cashtro/internal/ultron"
)

func main() {
	addr := flag.String("addr", envOr("ULTRON_ADDR", ":9090"), "Ultron IDE listen address")
	data := flag.String("data", envOr("ULTRON_DATA", "data/ultron.json"), "Ultron disk image")
	cashtro := flag.String("cashtro", envOr("CASHTRO_URL", "http://127.0.0.1:8080"), "Cashtro OS base URL")
	bootPW := flag.String("boot-password", envOr("ULTRON_OWNER_PASSWORD", "ultron-change-me"), "owner password on first boot")
	flag.Parse()

	if err := run(*addr, *data, *cashtro, *bootPW); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func run(addr, data, cashtroURL, bootPW string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	p := ultron.New(
		ultron.WithPersistPath(data),
		ultron.WithCashtroURL(cashtroURL),
		ultron.WithBootOwnerPassword(bootPW),
	)
	if err := ultron.LoadFile(data, p); err != nil {
		return err
	}
	if err := p.Boot(); err != nil {
		return err
	}
	_ = ultron.SaveFile(data, p)

	_, ok := p.PingCashtro()
	about := p.About(ok)
	log.Printf("%s %s · %d companies · %d agents · cashtro %s (%v)",
		about.Name, about.Version, about.Companies, about.Agents, about.CashtroURL, about.CashtroOK)

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           ultron.Handler(p),
		ReadHeaderTimeout: 5 * time.Second,
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("ultron IDE listening on http://%s · no Cursor token required", ln.Addr())

	errCh := make(chan error, 1)
	go func() { errCh <- httpSrv.Serve(ln) }()

	select {
	case <-ctx.Done():
		_ = ultron.SaveFile(data, p)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		_ = ultron.SaveFile(data, p)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
