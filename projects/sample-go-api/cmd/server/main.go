package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sample-go-api/internal/server"
)

func main() {
	srv := &http.Server{
		Addr:         addr(),
		Handler:      server.New(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func addr() string {
	if value := os.Getenv("PORT"); value != "" {
		return ":" + value
	}
	return ":8080"
}
