package repository

import (
	// "fmt"
	"sync"
	"time"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/models"
	"github.com/google/uuid"
	"github.com/pingcap/errors"
)

type TaskRepository struct {
	tasks map[string]*models.Task
	mu    sync.Mutex
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*models.Task),
	}
}

func (r *TaskRepository) Create(task *models.Task) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newTask := &models.Task{
		ID:        generateID(),
		Name:      task.Name,
		Status:    models.StatusCreated,
		URLs:      []string{},
		Errors:    []string{},
		ZipPath:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.tasks[newTask.ID] = newTask
	return newTask, nil
}

func generateID() string {
	return uuid.New().String()
}

func (r *TaskRepository) GetAllTasks() ([]*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.tasks == nil {
		return nil, errors.New("storage temporarily unavailable")
	}

	tasks := make([]*models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) AddURL(taskID string, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	if len(task.URLs) >= 3 {
		return errors.New("Validation Error. max 3 files per task")
	}

	task.URLs = append(task.URLs, url)
	if task.Status == models.StatusCreated {
		task.Status = models.StatusInProcess
	}

	if len(task.URLs) == 3 {
		task.Status = models.StatusCompleted
	}

	task.UpdatedAt = time.Now()
	return nil
}

func (r *TaskRepository) GetTaskByID(id string) (*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}
