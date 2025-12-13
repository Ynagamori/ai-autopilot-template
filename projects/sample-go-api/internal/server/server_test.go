package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHealth(t *testing.T) {
	srv := New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	assertStatus(t, rec.Code, http.StatusOK)
	var body map[string]string
	decodeBody(t, rec.Body, &body)
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body)
	}
}

func TestTaskLifecycle(t *testing.T) {
	srv := New()

	created := createTask(t, srv, "Write tests")
	if created.ID != 1 || created.Title != "Write tests" || created.Completed {
		t.Fatalf("unexpected created task: %+v", created)
	}

	fetched := getTask(t, srv, created.ID)
	if fetched != created {
		t.Fatalf("fetched task mismatch: %+v vs %+v", fetched, created)
	}

	updated := updateTask(t, srv, created.ID, map[string]any{"completed": true})
	if !updated.Completed || updated.Title != created.Title {
		t.Fatalf("unexpected updated task: %+v", updated)
	}

	updated = updateTask(t, srv, created.ID, map[string]any{"title": "Write more tests", "completed": false})
	if updated.Title != "Write more tests" || updated.Completed {
		t.Fatalf("unexpected updated task: %+v", updated)
	}

	list := listTasks(t, srv)
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("unexpected list result: %+v", list)
	}

	deleteTask(t, srv, created.ID)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusNotFound)
}

func TestValidation(t *testing.T) {
	srv := New()

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusBadRequest)

	req = httptest.NewRequest(http.MethodPut, "/tasks/abc", bytes.NewBufferString(`{"completed":true}`))
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusBadRequest)

	req = httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(`{}`))
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusBadRequest)
}

func createTask(t *testing.T, srv http.Handler, title string) Task {
	t.Helper()

	body := map[string]string{"title": title}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(data))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusCreated)

	var task Task
	decodeBody(t, rec.Body, &task)
	return task
}

func getTask(t *testing.T, srv http.Handler, id int) Task {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+itoa(id), nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusOK)

	var task Task
	decodeBody(t, rec.Body, &task)
	return task
}

func updateTask(t *testing.T, srv http.Handler, id int, payload map[string]any) Task {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+itoa(id), bytes.NewReader(data))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusOK)

	var task Task
	decodeBody(t, rec.Body, &task)
	return task
}

func listTasks(t *testing.T, srv http.Handler) []Task {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusOK)

	var response struct {
		Tasks []Task `json:"tasks"`
	}
	decodeBody(t, rec.Body, &response)
	return response.Tasks
}

func deleteTask(t *testing.T, srv http.Handler, id int) {
	t.Helper()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+itoa(id), nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	assertStatus(t, rec.Code, http.StatusNoContent)
}

func assertStatus(t *testing.T, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Fatalf("unexpected status: got %d want %d", actual, expected)
	}
}

func decodeBody(t *testing.T, body io.Reader, target any) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(target); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
