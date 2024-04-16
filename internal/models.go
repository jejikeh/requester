package internal

type TaskStatus string

const (
	TaskStatusPending TaskStatus = "new"
	TaskStatusRunning TaskStatus = "in_progress"
	TaskStatusDone    TaskStatus = "done"
	TaskStatusError   TaskStatus = "error"
)

type Task struct {
	ID              string            `json:"id"`
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	RequestHeaders  map[string]string `json:"req_headers,omitempty"`
	ResponseHeaders map[string]string `json:"res_headers"`
	Body            string            `json:"body"`
	Status          TaskStatus        `json:"status"`
	Error           string            `json:"error,omitempty"`
}

type CreateTaskRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}
