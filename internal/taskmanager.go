package internal

import (
	"context"
	"errors"
	"io"
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
	wg    *sync.WaitGroup
}

func NewInMemoryTaskManager() *InMemoryTaskManager {
	return &InMemoryTaskManager{
		tasks: make(map[string]*Task),
		wg:    &sync.WaitGroup{},
	}
}

func (t *InMemoryTaskManager) CreateTask(ctx context.Context, req CreateTaskRequest) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	id := uuid.New().String()

	task := Task{
		ID:              id,
		Method:          req.Method,
		URL:             req.URL,
		RequestHeaders:  req.Headers,
		ResponseHeaders: map[string]string{},
		Body:            req.Body,
		Status:          TaskStatusPending,
	}

	t.tasks[id] = &task

	if t.wg != nil {
		t.wg.Add(1)

		go func() {
			defer t.wg.Done()

			err := t.processTask(t.tasks[id])

			if err != nil {
				task.Error = err.Error()
				task.Status = TaskStatusError
			}
		}()
	}

	return id, nil
}

func (t *InMemoryTaskManager) GetTask(ctx context.Context, id string) (*Task, error) {
	task, ok := t.tasks[id]

	if !ok {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (t *InMemoryTaskManager) processTask(task *Task) error {
	task.Status = TaskStatusRunning

	resp, err := MakeRequest(task)

	if err != nil {
		task.Status = TaskStatusError
		task.Error = err.Error()

		return err
	}

	defer resp.Body.Close()

	task.Status = TaskStatusDone

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	task.Body = string(bodyBytes)

	for k, v := range resp.Header {
		task.ResponseHeaders[k] = v[0]
	}

	task.Status = TaskStatusDone

	return nil
}
