package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// API exposes HTTP handlers for the task service.
type API struct {
	store *TaskStore
}

// NewAPI constructs an API with a fresh in-memory store.
func NewAPI() *API {
	return &API{
		store: NewTaskStore(),
	}
}

// Routes returns the configured HTTP handler.
func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.healthHandler)
	mux.HandleFunc("/tasks", a.tasksHandler)
	mux.HandleFunc("/tasks/", a.taskActionHandler)
	return mux
}

func (a *API) healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := a.store.List()
		writeJSON(w, http.StatusOK, tasks)
	case http.MethodPost:
		var payload struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}

		task, err := a.store.Create(payload.Title)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, task)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *API) taskActionHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/tasks/") {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	if len(parts) == 0 {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	idStr := parts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if len(parts) == 2 && parts[1] == "complete" && r.Method == http.MethodPatch {
		task, err := a.store.Complete(id)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				writeError(w, http.StatusNotFound, "task not found")
				return
			}
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, task)
		return
	}

	writeError(w, http.StatusNotFound, fmt.Sprintf("no action for %s", r.URL.Path))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
