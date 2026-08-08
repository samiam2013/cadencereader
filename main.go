package main

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/samiam2013/cadencereader/config"
	"github.com/samiam2013/cadencereader/database"
)

func index(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fh, err := os.Open(cfg.ViewFolder + "/index.html")
		if err != nil {
			slog.Error("Failed to open index html view", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := io.ReadAll(fh)
		if err != nil {
			slog.Error("Failed to read index view", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			slog.Error("Failed to write index view", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func healthCheck(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			slog.Error("Database ping failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	slog.Info("CadenceReader Starting")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	http.HandleFunc("/health", healthCheck(db))

	http.HandleFunc("/", index(cfg))

	slog.Info("server starting")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("HTTP server exit", "error", err)
		os.Exit(1)
	}
}
