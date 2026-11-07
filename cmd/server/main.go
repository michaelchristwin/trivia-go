package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/michaelchristwin/trivia-go/internal/config"
)

func run(ctx context.Context, w io.Writer) error {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	router := SetupRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		fmt.Fprintf(w, "server started | addr=%s | pid=%d\n", srv.Addr, os.Getpid())
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		// shutdown signal received
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		return nil

	case err := <-errCh:

		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
