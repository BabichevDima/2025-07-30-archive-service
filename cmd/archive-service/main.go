package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	router "github.com/BabichevDima/2025-07-30-archive-service/internal/http"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/handlers"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/middleware"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/repository"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/service"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/usecase"
	"github.com/BabichevDima/2025-07-30-archive-service/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.L.Sync()

	storPath := "./storage"

	taskRepo := repository.NewTaskRepository()
	archiveService := service.NewArchiveServiceImpl(taskRepo, storPath)
	taskUsecase := usecase.NewTaskUsecase(taskRepo, archiveService, 3)
	taskHandler := handlers.NewTaskHandler(taskUsecase)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, taskHandler)
	handler := middleware.RequestLogger(logger.L, mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("Graceful shutdown timed out.. forcing exit")
			}
		}()

		logger.Info("Starting graceful shutdown...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server shutdown failed", zap.Error(err))
		}
		serverStopCtx()
	}()

	logger.Info("HTTP server is listening",
		zap.String("address", "http://localhost"+server.Addr),
	)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	<-serverCtx.Done()
	logger.Info("Server stopped gracefully")
}
