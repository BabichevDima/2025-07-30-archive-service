package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BabichevDima/2025-07-30-archive-service/pkg/logger"
	"go.uber.org/zap"
)

type ShutdownFunc func(ctx context.Context) error

func GracefulShutdown(
	cancel context.CancelFunc,
	server *http.Server,
	log *zap.Logger,
	timeout time.Duration,
	shutdownFuncs ...ShutdownFunc,
) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancelTimeout := context.WithTimeout(context.Background(), timeout)
	defer cancelTimeout()

	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown failed", zap.Error(err))
		}
	}

	for _, shutdownFn := range shutdownFuncs {
		if shutdownFn != nil {
			if err := shutdownFn(ctx); err != nil {
				logger.Error("Error during shutdown", zap.Error(err))
			}
		}
	}

	if cancel != nil {
		cancel()
	}

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			logger.Warn("Graceful shutdown timed out")
		}
	default:
	}

	logger.Info("Server exited properly")
}
