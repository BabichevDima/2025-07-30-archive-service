package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/dto"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/response"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/usecase"
)

type TaskHandler struct {
	usecase *usecase.TaskUsecase
}

func NewTaskHandler(u *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{usecase: u}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	request := dto.RequestTask{}
	err := decoder.Decode(&request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown field"):
			response.RespondWithError(w, http.StatusTooManyRequests, "Unknown field", err)
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

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasksResponse, err := h.usecase.GetAllTasks()
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unavailable"):
			response.RespondWithError(w, http.StatusServiceUnavailable, "Service unavailable", err)
		default:
			response.RespondWithError(w, http.StatusInternalServerError, "Failed to get tasks", err)
		}
		return
	}
	response.RespondWithJSON(w, http.StatusOK, tasksResponse)
}

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
			response.RespondWithError(w, http.StatusTooManyRequests, "Unknown field", err)
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
