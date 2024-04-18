package internal

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskManager interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) (string, error)
	GetTask(ctx context.Context, id string) (*Task, error)
}

type InMemoryTaskManager struct {
	tasks  map[string]*Task
	client *Client

	mutex sync.Mutex
	wg    *sync.WaitGroup
}

func NewInMemoryTaskManager(client *Client) *InMemoryTaskManager {
	return &InMemoryTaskManager{
		tasks:  make(map[string]*Task),
		wg:     &sync.WaitGroup{},
		client: client,
	}
}

func (t *InMemoryTaskManager) CreateTask(ctx context.Context, req CreateTaskRequest) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Printf("Creating task with %s %s", req.Method, req.URL)

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

	resp, err := t.client.MakeRequest(task)

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

	log.Printf("Task %s done", task.ID)

	return nil
}
