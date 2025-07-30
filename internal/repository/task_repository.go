package repository

import (
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
		ID:        generateTaskID(),
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

func generateTaskID() string {
	return uuid.New().String()
}

func (r *TaskRepository) GetAllTasks() []*models.Task {
    r.mu.Lock()
    defer r.mu.Unlock()

    tasks := make([]*models.Task, 0, len(r.tasks))
    for _, task := range r.tasks {
        tasks = append(tasks, task)
    }

    return tasks
}

func (r *TaskRepository) AddURL(taskID string, url string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    task, exists := r.tasks[taskID]
    if !exists {
        return errors.New("task not found")
    }

    task.URLs = append(task.URLs, url)
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