package main

import (
	"log/slog"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world\n"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// TODO ping database
	w.WriteHeader(http.StatusOK)
}

func readyCheck(w http.ResponseWriter, r *http.Request) {
	// TODO make some meaningful readiness check
	w.WriteHeader(http.StatusOK)
}

func main() {
	slog.Info("binary started")

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/ready", readyCheck)

	http.HandleFunc("/", index)
	slog.Info("server starting")
	err := http.ListenAndServe(":80", nil)
	slog.Error("HTTP server exit", "error", err)
}
