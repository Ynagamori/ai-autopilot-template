package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Task represents a unit of work managed by the API.
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Server exposes HTTP handlers for task management.
type Server struct {
	mux    *http.ServeMux
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

// New creates a ready-to-use Server instance.
func New() *Server {
	s := &Server{
		mux:    http.NewServeMux(),
		tasks:  make(map[int]Task),
		nextID: 1,
	}
	s.registerRoutes()
	return s
}

// ServeHTTP satisfies the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /tasks", s.handleListTasks)
	s.mux.HandleFunc("POST /tasks", s.handleCreateTask)
	s.mux.HandleFunc("GET /tasks/{id}", s.handleGetTask)
	s.mux.HandleFunc("PUT /tasks/{id}", s.handleUpdateTask)
	s.mux.HandleFunc("DELETE /tasks/{id}", s.handleDeleteTask)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListTasks(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]int, 0, len(s.tasks))
	for id := range s.tasks {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	tasks := make([]Task, 0, len(s.tasks))
	for _, id := range ids {
		tasks = append(tasks, s.tasks[id])
	}

	s.writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	payload.Title = strings.TrimSpace(payload.Title)
	if payload.Title == "" {
		s.writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	task := Task{ID: id, Title: payload.Title, Completed: false}
	s.tasks[id] = task

	s.writeJSON(w, http.StatusCreated, task)
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	task, ok := s.tasks[id]
	s.mu.Unlock()

	if !ok {
		s.writeError(w, http.StatusNotFound, "task not found")
		return
	}

	s.writeJSON(w, http.StatusOK, task)
}

func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if payload.Title == nil && payload.Completed == nil {
		s.writeError(w, http.StatusBadRequest, "nothing to update")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		s.writeError(w, http.StatusNotFound, "task not found")
		return
	}

	if payload.Title != nil {
		title := strings.TrimSpace(*payload.Title)
		if title == "" {
			s.writeError(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		task.Title = title
	}

	if payload.Completed != nil {
		task.Completed = *payload.Completed
	}

	s.tasks[id] = task
	s.writeJSON(w, http.StatusOK, task)
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		s.writeError(w, http.StatusNotFound, "task not found")
		return
	}

	delete(s.tasks, id)
	w.WriteHeader(http.StatusNoContent)
}

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid task id")
	}
	return id, nil
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}
