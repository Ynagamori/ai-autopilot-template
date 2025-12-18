package server

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

var (
	// ErrTaskNotFound is returned when a task does not exist in the store.
	ErrTaskNotFound = errors.New("task not found")
)

// Task represents a simple todo item stored by the in-memory API.
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// TaskStore manages Task records in memory with concurrency safety.
type TaskStore struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

// NewTaskStore creates a TaskStore with initialized fields.
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

// List returns a stable, ID-sorted slice of stored tasks.
func (s *TaskStore) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		items = append(items, task)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items
}

// Create inserts a new task with the provided title.
func (s *TaskStore) Create(title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, errors.New("title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	task := Task{
		ID:    id,
		Title: title,
		Done:  false,
	}
	s.tasks[id] = task

	return task, nil
}

// Complete marks the target task as done and returns the updated record.
func (s *TaskStore) Complete(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	if task.Done {
		return task, nil
	}

	task.Done = true
	s.tasks[id] = task

	return task, nil
}
