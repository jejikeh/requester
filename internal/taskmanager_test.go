package internal

import (
	"context"
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
		if task.URL != url {
			t.Errorf("the url %s in task manager doesnt equal in CreateTaskRequest %s", task.URL, url)
		}

		if task.Method != method {
			t.Errorf("expected %s method, but got %s", method, task.Method)
		}

		if task.Body != body {
			t.Errorf("expected %s body, but got %s", body, task.Body)
		}
	}
}
