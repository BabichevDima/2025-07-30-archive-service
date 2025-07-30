package models

import (
	"time"
)

type Task struct {
	ID        string
	Name      string
	Status    TaskStatus
	URLs      []string
	Errors    []string
	ZipPath   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskStatus string

const (
	StatusCreated    TaskStatus = "created"
	StatusProcessing TaskStatus = "in process"
	StatusCompleted  TaskStatus = "completed"
	// check if "failed" status neaded
	// StatusFailed     TaskStatus = "failed"
)
