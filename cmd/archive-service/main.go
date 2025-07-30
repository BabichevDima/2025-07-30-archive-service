package main

import (
	"context"
	"net/http"
	"time"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/handlers"
	router "github.com/BabichevDima/2025-07-30-archive-service/internal/http"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/middleware"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/repository"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/usecase"
	"github.com/BabichevDima/2025-07-30-archive-service/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.L.Sync()

	taskRepo := repository.NewTaskRepository()
	taskUsecase := usecase.NewTaskUsecase(taskRepo, 3)
	taskHandler := handlers.NewTaskHandler(taskUsecase)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, taskHandler)
	handler := middleware.RequestLogger(logger.L, mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info("HTTP server is listening",
			zap.String("address", "http://localhost"+server.Addr),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown failed", zap.Error(err))
	}
	logger.Info("Server stopped")
}
