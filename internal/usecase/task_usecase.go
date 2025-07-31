package usecase

import (
	"fmt"
	"sync"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/dto"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/models"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/repository"
	// "errors"
)

type TaskUsecase struct {
	repo     *repository.TaskRepository
	maxTasks int
	active   int
	mu       sync.Mutex
}

func NewTaskUsecase(repo *repository.TaskRepository, maxTasks int) *TaskUsecase {
	return &TaskUsecase{
		repo:     repo,
		maxTasks: maxTasks,
	}
}

func (u *TaskUsecase) Create(request dto.RequestTask) (dto.ResponseTask, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.active >= u.maxTasks {
		return dto.ResponseTask{}, fmt.Errorf("server is busy (max %d tasks allowed)", u.maxTasks)
	}

	resp := &models.Task{
		Name: request.Name,
	}

	taskResp, err := u.repo.Create(resp)
	if err != nil {
		return dto.ResponseTask{}, err
	}

	u.active++
	return dto.ResponseTask{
		ID:        taskResp.ID,
		Name:      taskResp.Name,
		Status:    string(taskResp.Status),
		URLs:      taskResp.URLs,
		Errors:    taskResp.Errors,
		ZipPath:   taskResp.ZipPath,
		CreatedAt: taskResp.CreatedAt,
		UpdatedAt: taskResp.UpdatedAt,
	}, nil
}

func (u *TaskUsecase) GetAllTasks() ([]*models.Task, error) {
	return u.repo.GetAllTasks()
}

func (uc *TaskUsecase) AddURL(taskID string, url string) error {
	if err := uc.repo.AddURL(taskID, url); err != nil {
		return fmt.Errorf("failed to add URL: %w", err)
	}
	return nil
}

func (uc *TaskUsecase) GetTaskStatus(taskID string) (dto.TaskStatusResponse, error) {
	task, err := uc.repo.GetTaskByID(taskID)
	if err != nil {
		return dto.TaskStatusResponse{}, fmt.Errorf("failed to get task status: %w", err)
	}
	return dto.TaskStatusResponse{
		Status: string(task.Status),
	}, nil
}
