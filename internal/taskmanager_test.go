package internal

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateTask(t *testing.T) {
	client := NewClient()
	taskManager := NewInMemoryTaskManager(client)

	url := "http://google.com"
	method := "GET"
	headers := map[string]string{}
	body := ""

	id, err := taskManager.CreateTask(context.Background(), CreateTaskRequest{
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
	})

	if err != nil {
		t.Error(err)
	}

	if task, ok := taskManager.tasks[id]; !ok {
		t.Errorf("the task with %s was not appear in sheduled tasks", id)
	} else {
		err = compareTasks(task, &Task{URL: url, Method: method, Body: body})

		if err != nil {
			t.Error(err)
		}
	}
}

func TestGetTaskNotFound(t *testing.T) {
	client := NewClient()
	taskManager := NewInMemoryTaskManager(client)

	url := "http://google.com"
	method := "GET"
	headers := map[string]string{}
	body := ""

	id, err := taskManager.CreateTask(context.Background(), CreateTaskRequest{
		URL:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
	})

	if err != nil {
		t.Error(err)
	}

	if task, err := taskManager.GetTask(context.Background(), id); err != nil {
		t.Error(err)
	} else {
		err = compareTasks(task, &Task{URL: url, Method: method, Body: body})

		if err != nil {
			t.Error(err)
		}
	}
}

func compareTasks(t1 *Task, t2 *Task) error {
	if t2.URL != t1.URL {
		return fmt.Errorf("expected %s url, but got %s", t2.URL, t1.URL)
	}

	if t2.Method != t1.Method {
		return fmt.Errorf("expected %s metho, but got %s", t2.Method, t1.Method)

	}

	if t2.Body != t1.Body {
		return fmt.Errorf("expected %s body, but got %s", t2.Body, t1.Body)
	}

	return nil
}
