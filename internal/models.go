package internal

import "sync"

type TaskStatus string

const (
	TaskStatusPending TaskStatus = "new"
	TaskStatusRunning TaskStatus = "in_progress"
	TaskStatusDone    TaskStatus = "done"
	TaskStatusError   TaskStatus = "error"
)

type Task struct {
	ID      string            `json:"id"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Status  TaskStatus        `json:"status"`
	Error   error             `json:"error,omitempty"`

	mutex sync.Mutex
}

type CreateTaskRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
