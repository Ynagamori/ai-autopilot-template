package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()

	newServer().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", res.Code)
	}

	var payload healthResponse
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Status != "ok" {
		t.Fatalf("unexpected status: %s", payload.Status)
	}
}

func TestEchoHandler(t *testing.T) {
	body := bytes.NewBufferString(`{"message":"hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/echo", body)
	res := httptest.NewRecorder()

	newServer().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", res.Code)
	}

	var payload echoResponse
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Echo != "hello" {
		t.Fatalf("unexpected echo message: %s", payload.Echo)
	}
}

func TestEchoHandlerValidation(t *testing.T) {
	t.Run("rejects non-POST", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/echo", nil)
		res := httptest.NewRecorder()

		newServer().ServeHTTP(res, req)

		if res.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405, got %d", res.Code)
		}
	})

	t.Run("rejects missing message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(`{"message":""}`))
		res := httptest.NewRecorder()

		newServer().ServeHTTP(res, req)

		if res.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", res.Code)
		}
	})
}
