package usecase

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/dto"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/models"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/repository"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/http/service"
	"github.com/gabriel-vasile/mimetype"
)

type TaskUsecase struct {
	repo       *repository.TaskRepository
	archiveSvc service.ArchiveService
	maxTasks   int
	active     int
	mu         sync.Mutex
}

func NewTaskUsecase(repo *repository.TaskRepository, archiveSvc service.ArchiveService, maxTasks int) *TaskUsecase {
	return &TaskUsecase{
		repo:       repo,
		archiveSvc: archiveSvc,
		maxTasks:   maxTasks,
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

func (u *TaskUsecase) AddURL(taskID string, url string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	respFileData, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file by path: %v", url)
	}
	defer respFileData.Body.Close()

	if respFileData.StatusCode != http.StatusOK {
		return fmt.Errorf("file unavailable (status %d): %s", respFileData.StatusCode, url)
	}

	limitedReader := io.LimitReader(respFileData.Body, 512)
	mime, err := mimetype.DetectReader(limitedReader)
	if err != nil {
		return fmt.Errorf("failed to detect file type: %v", err)
	}

	if mime.String() == "application/pdf" || mime.String() == "image/jpeg" {
		if err := u.repo.AddURL(taskID, url); err != nil {
			return fmt.Errorf("failed to add URL: %w", err)
		}

		task, err := u.repo.GetTask(taskID)
		if err != nil {
			return err
		}

		if len(task.URLs) == 3 {
			go func() {
				if err := u.archiveSvc.CreateArchive(taskID, task.URLs); err != nil {
					log.Printf("Archive failed for task %s: %v", taskID, err)
					u.repo.UpdateTaskStatus(taskID, models.StatusFailed)
				} else {
					u.active--
				}
			}()
		}

		return nil
	}

	if err := u.repo.AddURL(taskID, url); err != nil {
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
