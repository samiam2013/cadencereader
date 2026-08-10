package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samiam2013/cadencereader/config"
	"github.com/samiam2013/cadencereader/database"
)

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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthCheck(db))
	mux.HandleFunc("GET /",
		static(fmt.Sprintf("%s/%s", cfg.ViewFolder, "index.html")))
	mux.HandleFunc("GET /add-blog",
		static(fmt.Sprintf("%s/%s", cfg.ViewFolder, "add-blog.html")))
	mux.HandleFunc("POST /add-blog", addBlog(db))

	intrC := make(chan os.Signal, 1)
	signal.Notify(intrC, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPport),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-intrC
		slog.Warn("received interrupt, shutting down gracefully")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
		}
	}()

	slog.Info("server starting")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server exit", "error", err)
		os.Exit(1)
	}

	// srv.Shutdown() has returned (or timed out) by this point — safe to close DB now
	if err := db.Close(); err != nil {
		slog.Error("failed to close database", "error", err)
		os.Exit(1)
	}
	slog.Info("shutdown complete")
}

func static(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fh, err := os.Open(path)
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
