package router

import (
	"net/http"

	_ "github.com/BabichevDima/2025-07-30-archive-service/internal/docs"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux, taskHandler *handlers.TaskHandler) {
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	mux.Handle("POST /api/tasks", http.HandlerFunc(taskHandler.Create))
	mux.Handle("GET /api/tasks", http.HandlerFunc(taskHandler.GetAllTasks))
	mux.Handle("POST /api/tasks/{id}/urls", http.HandlerFunc(taskHandler.AddURL))
	mux.Handle("GET /api/tasks/{id}/status", http.HandlerFunc(taskHandler.GetTaskStatus))
}
