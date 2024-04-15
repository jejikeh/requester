package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jejikeh/requester/internal"
)

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Router struct {
	taskManager internal.TaskManager
}

func NewRouter(taskManager internal.TaskManager) *Router {
	return &Router{taskManager: taskManager}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router := mux.NewRouter()

	router.HandleFunc("/task", r.createTaskHandler).Methods(http.MethodPost)

	router.HandleFunc("/task/{id}", r.getTaskHandler).Methods(http.MethodGet)

	router.ServeHTTP(w, req)
}

func (r *Router) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	var task internal.CreateTaskRequest
	err := json.NewDecoder(req.Body).Decode(&task)

	if err != nil {
		writeErrorMessage(w, err)

		return
	}

	id, err := r.taskManager.CreateTask(req.Context(), task)

	if err != nil {
		writeErrorMessage(w, err)

		return
	}

	response := CreateTaskResponse{
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (r *Router) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	task, err := r.taskManager.GetTask(req.Context(), id)

	if err != nil {
		writeErrorMessage(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(task)
}

func writeErrorMessage(w http.ResponseWriter, err error) {
	var code int

	switch err {
	case internal.ErrTaskNotFound:
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}

	response := ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(response)
}
