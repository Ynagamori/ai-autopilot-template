package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestTaskLifecycle(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	client := srv.Client()

	// Initially empty.
	tasks := fetchTasks(t, client, srv.URL)
	if len(tasks) != 0 {
		t.Fatalf("expected no tasks initially, got %d", len(tasks))
	}

	// Create two tasks.
	first := createTask(t, client, srv.URL, "Write tests")
	second := createTask(t, client, srv.URL, "Ship features")

	if first.ID == 0 || second.ID == 0 || first.ID == second.ID {
		t.Fatalf("task IDs should be unique and non-zero, got %v and %v", first.ID, second.ID)
	}

	// Mark the first task complete.
	completed := completeTask(t, client, srv.URL, first.ID)
	if !completed.Done {
		t.Fatalf("task should be marked done after completion")
	}

	// Update second task title.
	updated := updateTaskTitle(t, client, srv.URL, second.ID, "Ship polished features")
	if updated.Title != "Ship polished features" {
		t.Fatalf("expected updated title, got %+v", updated)
	}

	// Delete first task.
	deleted := deleteTask(t, client, srv.URL, first.ID)
	if deleted.ID != first.ID {
		t.Fatalf("expected to delete task %d, got %+v", first.ID, deleted)
	}

	// Confirm list includes updates.
	tasks = fetchTasks(t, client, srv.URL)
	if len(tasks) != 1 {
		t.Fatalf("expected one task, got %d", len(tasks))
	}

	if tasks[0].ID != second.ID || tasks[0].Title != "Ship polished features" {
		t.Fatalf("expected updated second task, got %+v", tasks[0])
	}
}

func TestTaskListPagination(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.Routes())
	defer srv.Close()

	client := srv.Client()

	createTask(t, client, srv.URL, "First")
	createTask(t, client, srv.URL, "Second")
	createTask(t, client, srv.URL, "Third")

	tasks := fetchTasksWithQuery(t, client, srv.URL, "?offset=1&limit=1")
	if len(tasks) != 1 {
		t.Fatalf("expected one task, got %d", len(tasks))
	}
	if tasks[0].Title != "Second" {
		t.Fatalf("expected second task, got %+v", tasks[0])
	}
}

func fetchTasks(t *testing.T, client *http.Client, baseURL string) []Task {
	t.Helper()

	resp, err := client.Get(baseURL + "/tasks")
	if err != nil {
		t.Fatalf("failed to fetch tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for list, got %d", resp.StatusCode)
	}

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode tasks: %v", err)
	}

	return tasks
}

func fetchTasksWithQuery(t *testing.T, client *http.Client, baseURL, query string) []Task {
	t.Helper()

	resp, err := client.Get(baseURL + "/tasks" + query)
	if err != nil {
		t.Fatalf("failed to fetch tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for list, got %d", resp.StatusCode)
	}

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode tasks: %v", err)
	}

	return tasks
}

func createTask(t *testing.T, client *http.Client, baseURL, title string) Task {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"title": title})
	resp, err := client.Post(baseURL+"/tasks", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 for create, got %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode task: %v", err)
	}

	return task
}

func completeTask(t *testing.T, client *http.Client, baseURL string, id int) Task {
	t.Helper()

	req, err := http.NewRequest(http.MethodPatch, baseURL+"/tasks/"+itoa(id)+"/complete", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("complete task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for complete, got %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode completed task: %v", err)
	}

	return task
}

func updateTaskTitle(t *testing.T, client *http.Client, baseURL string, id int, title string) Task {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"title": title})
	req, err := http.NewRequest(http.MethodPatch, baseURL+"/tasks/"+itoa(id), bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("update task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for update, got %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode updated task: %v", err)
	}

	return task
}

func deleteTask(t *testing.T, client *http.Client, baseURL string, id int) Task {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, baseURL+"/tasks/"+itoa(id), nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("delete task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for delete, got %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode deleted task: %v", err)
	}

	return task
}

func itoa(n int) string {
	return strconv.FormatInt(int64(n), 10)
}
