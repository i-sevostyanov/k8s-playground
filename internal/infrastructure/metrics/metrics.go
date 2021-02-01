package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(ctx context.Context, address string) error {
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- server.ListenAndServe()
	}()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	select {
	case <-ctx.Done():
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
