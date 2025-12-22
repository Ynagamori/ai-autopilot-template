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

	// Confirm list includes updates.
	tasks = fetchTasks(t, client, srv.URL)
	if len(tasks) != 2 {
		t.Fatalf("expected two tasks, got %d", len(tasks))
	}

	if tasks[0].ID != first.ID || tasks[0].Done != true {
		t.Fatalf("expected first task done, got %+v", tasks[0])
	}

	if tasks[1].ID != second.ID || tasks[1].Done {
		t.Fatalf("expected second task not done, got %+v", tasks[1])
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

func itoa(n int) string {
	return strconv.FormatInt(int64(n), 10)
}
