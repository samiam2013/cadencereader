package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world\n"))
}

var hcIdx int = 0

func healthCheck(w http.ResponseWriter, r *http.Request) {
	hcIdx++
	fmt.Println("health check count", hcIdx)
	if hcIdx%7 == 0 {
		fmt.Println("inducing failure")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func readyCheck(w http.ResponseWriter, r *http.Request) {
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
