package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

type echoRequest struct {
	Message string `json:"message"`
}

type echoResponse struct {
	Echo string `json:"echo"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	respondJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	defer r.Body.Close()

	var payload echoRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
		return
	}

	if payload.Message == "" {
		respondError(w, http.StatusBadRequest, errors.New("message is required"))
		return
	}

	respondJSON(w, http.StatusOK, echoResponse{Echo: payload.Message})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, err error) {
	respondJSON(w, status, map[string]string{"error": err.Error()})
}

func newServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/echo", echoHandler)
	return mux
}

func main() {
	addr := ":8080"
	log.Printf("starting sample-go-api at %s", addr)

	if err := http.ListenAndServe(addr, newServer()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
