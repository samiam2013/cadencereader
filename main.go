package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world\n"))
}

func healthCheck(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(context.Background()); err != nil {
			slog.Error("Database ping failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	slog.Info("CadenceReader Starting")

	dbURL := os.Getenv("DATABASE_URL")
	if len(strings.TrimSpace(dbURL)) == 0 {
		slog.Error("DATABASE_URL empty and required")
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		slog.Error("Could not connect to database", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close(ctx) }()

	if err := db.Ping(ctx); err != nil {
		slog.Error("DB ping failed after opening", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/health", healthCheck(db))

	http.HandleFunc("/", index)

	slog.Info("server starting")
	if err := http.ListenAndServe(":80", nil); err != nil {
		slog.Error("HTTP server exit", "error", err)
		os.Exit(1)
	}
}
