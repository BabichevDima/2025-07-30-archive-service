package handlers

import (
	"encoding/json"
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

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	decoder := json.NewDecoder(r.Body)
	request := dto.RequestTask{}
	err := decoder.Decode(&request)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if request.Name == "" {
		response.RespondWithError(w, http.StatusBadRequest, "Name is required", err)
		return
	}

	taskResponse, err := h.usecase.Create(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, taskResponse)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	
	// TODO: !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	tasks := h.usecase.GetAllTasks()

	response.RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) AddURL(w http.ResponseWriter, r *http.Request) {
	
	// TODO: !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	taskId := r.PathValue("id")
    if taskId == "" {
        response.RespondWithError(w, http.StatusBadRequest, "task ID is required", nil)
        return
    }

    var req dto.AddURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

    if err := h.usecase.AddURL(taskId, req.URL); err != nil {
        status := http.StatusInternalServerError
        if strings.Contains(err.Error(), "not found") {
            status = http.StatusNotFound
        }
		response.RespondWithError(w, status, "Invalid request payload", err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	
	// TODO: !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	taskId := r.PathValue("id")
    if taskId == "" {
        response.RespondWithError(w, http.StatusBadRequest, "task ID is required", nil)
        return
    }

    statusResponse, err := h.usecase.GetTaskStatus(taskId)
    if err != nil {
        statusCode := http.StatusInternalServerError
        if strings.Contains(err.Error(), "not found") {
            statusCode = http.StatusNotFound
        }
		response.RespondWithError(w, statusCode, "Invalid request payload", err)
        return
    }


	response.RespondWithJSON(w, http.StatusOK, statusResponse)
}
