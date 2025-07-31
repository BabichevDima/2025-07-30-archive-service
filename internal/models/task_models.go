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
	StatusCreated   TaskStatus = "Created"
	StatusInProcess TaskStatus = "In process"
	StatusCompleted TaskStatus = "Completed"
	StatusFailed    TaskStatus = "Failed"
)
