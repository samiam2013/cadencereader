package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/somethingsoftware/cadencereader/config"
	"github.com/somethingsoftware/cadencereader/database"
)

func main() {
	slog.Info("CadenceReader Starting")

	cfg, err := config.Load("VIEW_FOLDER")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	conn, queries, err := database.Open(ctx, cfg)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthCheck(conn))
	mux.HandleFunc("GET /",
		static(fmt.Sprintf("%s/%s", cfg.Extra["VIEW_FOLDER"], "index.html")))
	mux.HandleFunc("GET /add-blog",
		static(fmt.Sprintf("%s/%s", cfg.Extra["VIEW_FOLDER"], "add-blog.html")))
	mux.HandleFunc("POST /add-blog", addBlog(queries))

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
