package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/dto"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/response"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/usecase"
)

type TaskHandler struct {
	usecase *usecase.TaskUsecase
}

func NewTaskHandler(u *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{usecase: u}
}

// @title Archive Service API
// @version 1.0
// @description Сервис для создания ZIP архивов из файлов

// @contact.name API Support
// @contact.email support@archive-service.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

// Create godoc
// @Summary Создать новую задачу архивации
// @Description Создает новую задачу для последующего добавления URL файлов
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dto.RequestTask true "Данные для создания задачи"
// @Success 201 {object} dto.ResponseTask
// @Failure 400 {object} response.BadRequestError
// @Failure 429 {object} response.ServerBusyRequestError
// @Failure 500 {object} response.InternalServerError
// @Router /api/tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	request := dto.RequestTask{}
	err := decoder.Decode(&request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown field"):
			response.RespondWithError(w, http.StatusBadRequest, "Unknown field", err)
		default:
			response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		}
		return
	}

	if request.Name == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Name is required", err)
		return
	}

	taskResponse, err := h.usecase.Create(request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "server is busy"):
			response.RespondWithError(w, http.StatusTooManyRequests, "Server is busy", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		}
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, taskResponse)
}

// GetAllTasks godoc
// @Summary Получить список всех задач
// @Description Возвращает список всех задач архивации
// @Tags tasks
// @Produce json
// @Success 200 {array} dto.ResponseTask
// @Failure 500 {object} response.InternalServerError
// @Router /api/tasks [get]
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasksResponse, err := h.usecase.GetAllTasks()
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unavailable"):
			response.RespondWithError(w, http.StatusServiceUnavailable, "Service unavailable", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		}
		return
	}
	response.RespondWithJSON(w, http.StatusOK, tasksResponse)
}

// AddURL godoc
// @Summary Добавить URL в задачу
// @Description Добавляет URL файла для загрузки в указанную задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID задачи"
// @Param request body dto.URLRequest true "URL файла"
// @Success 204
// @Failure 400 {object} response.BadRequestError
// @Failure 404 {object} response.NotFoundRequestError
// @Failure 422 {object} response.ConstrainsErrorResponse
// @Failure 500 {object} response.InternalServerError
// @Router /api/tasks/{id}/urls [post]
func (h *TaskHandler) AddURL(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "taskID is required", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req dto.URLRequest
	err := decoder.Decode(&req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown field"):
			response.RespondWithError(w, http.StatusBadRequest, "Unknown field", err)
		default:
			response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		}
		return
	}

	if err := h.usecase.AddURL(taskID, req.URL); err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			response.RespondWithError(w, http.StatusNotFound, "Task not found or was deleted", err)
		case strings.Contains(err.Error(), "Validation Error"):
			response.RespondWithError(w, http.StatusUnprocessableEntity, "You can only upload up to 3 files per task", err)
		case strings.Contains(err.Error(), "unavailable"):
			response.RespondWithError(w, http.StatusUnprocessableEntity, fmt.Sprintf("URL is not available: %v", req.URL), err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		}
		return
	}

	response.RespondWithJSON(w, http.StatusNoContent, nil)
}

// GetTaskStatus godoc
// @Summary Получить статус задачи
// @Description Возвращает текущий статус задачи и ссылку на архив (если готов)
// @Tags tasks
// @Produce json
// @Param id path string true "ID задачи"
// @Success 200 {object} dto.TaskStatusResponse
// @Failure 400 {object} response.BadRequestError
// @Failure 404 {object} response.NotFoundRequestError
// @Failure 500 {object} response.InternalServerError
// @Router /api/tasks/{id}/status [get]
func (h *TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("id")
	if taskID == "" {
		response.RespondWithError(w, http.StatusBadRequest, "task ID is required", nil)
		return
	}

	statusResponse, err := h.usecase.GetTaskStatus(taskID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			response.RespondWithError(w, http.StatusNotFound, "not found", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		}
		return
	}

	response.RespondWithJSON(w, http.StatusOK, statusResponse)
}
