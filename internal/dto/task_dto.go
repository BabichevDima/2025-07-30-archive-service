package dto

import (
	"time"
)

type RequestTask struct {
	Name string `json:"name" binding:"required"`
}

type ResponseTask struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	URLs      []string  `json:"urls"`
	Errors    []string  `json:"errors"`
	ZipPath   string    `json:"zip_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type TaskStatusResponse struct {
    Status    string    `json:"status"`
}