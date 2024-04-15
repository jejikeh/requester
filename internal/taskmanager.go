package internal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskManager interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) (string, error)
	GetTask(ctx context.Context, id string) (*Task, error)
}

type InMemoryTaskManager struct {
	tasks map[string]*Task
	mutex sync.Mutex
}

func NewInMemoryTaskManager() *InMemoryTaskManager {
	return &InMemoryTaskManager{
		tasks: make(map[string]*Task),
	}
}

func (t *InMemoryTaskManager) CreateTask(ctx context.Context, req CreateTaskRequest) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	id := uuid.New().String()

	task := Task{
		ID:      id,
		Method:  req.Method,
		URL:     req.URL,
		Headers: req.Headers,
		Body:    req.Body,
		Status:  TaskStatusPending,
	}

	t.tasks[id] = &task

	return id, nil
}

func (t *InMemoryTaskManager) GetTask(ctx context.Context, id string) (*Task, error) {
	task, ok := t.tasks[id]

	if !ok {
		return &Task{}, ErrTaskNotFound
	}

	return task, nil
}
