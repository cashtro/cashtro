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
	pulseEvery := flag.Duration("pulse", 8*time.Second, "never-stop pulse interval; 0 disables")
	flag.Parse()

	if err := run(*addr, *pulseEvery); err != nil {
		log.Fatal(err)
	}
}

func run(addr string, pulseEvery time.Duration) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	k, err := agents.Boot()
	if err != nil {
		return err
	}
	about := k.About()
	log.Printf("%s %s · %d agentics online · kimi+glm dual · pulse %s", about.Name, about.Version, about.Running, pulseEvery)

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
	if pulseEvery > 0 {
		go runPulse(ctx, k, pulseEvery)
	}

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

func runPulse(ctx context.Context, k *kernel.Kernel, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	tick := func() {
		res, err := k.Invoke(ctx, "pulse", kernel.Call{Capability: "pulse.tick"})
		if err != nil {
			log.Printf("pulse: %v", err)
			return
		}
		log.Printf("%s", res.Message)
	}
	tick()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			tick()
		}
	}
}
