package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultServerTimeout         = 30 * time.Second
	defaultServerShutdownTimeout = 10 * time.Second
)

func Listen(ctx context.Context, address string) error {
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         address,
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  defaultServerTimeout,
		WriteTimeout: defaultServerTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), defaultServerShutdownTimeout)
		defer shutdownCancel()

		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
