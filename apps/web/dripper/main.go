package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/somethingsoftware/cadencereader/config"
	"github.com/somethingsoftware/cadencereader/database"
)

func main() {
	slog.Info("Dripper Starting")

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
	mux.HandleFunc("GET /post/{post_id}/{segment_type}/{index}", drip(ctx, cfg, queries))
	intrC := make(chan os.Signal, 1)
	signal.Notify(intrC, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8081), // TODO add this port to config
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

var segTypes = []string{"word", "sentence", "paragraph"}

func drip(ctx context.Context, cfg config.Config, queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		postIDstr := r.PathValue("post_id")
		postID, err := strconv.Atoi(postIDstr)
		if err != nil {
			http.Error(w, "failed to parse post id in path", http.StatusInternalServerError)
			slog.Error("Failed to parse post id", "error", err)
			return
		}

		segmentType := r.PathValue("segment_type")
		if !slices.Contains(segTypes, segmentType) {
			http.Error(w, "failed to recognize segment type", http.StatusInternalServerError)
			slog.Error("Segment type not recognized", "segment_type", segmentType)
			return
		}

		idxStr := r.PathValue("index")
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			http.Error(w, "failed to parse post id in path", http.StatusInternalServerError)
			slog.Error("Failed to parse post id", "error", err)
			return
		}

		// TODO: cache the post lookup
		post, err := queries.GetBlogPost(ctx, int32(postID))
		if err != nil {
			http.Error(w, "failed to find post", http.StatusInternalServerError)
			slog.Error("Failed to get blog post", "error", err)
			return
		}

		var parts []string
		switch segmentType {
		case "word":
			parts = strings.Split(post.Content, " ")
		case "sentence":
			parts = strings.Split(post.Content, ".?!;")
		case "paragraph":
			parts = strings.Split(post.Content, "\n")
		}

		if idx >= len(parts) {
			http.Error(w, "index beyond parts length by segment type", http.StatusBadRequest)
			slog.Error("Parts length didn't allow for index requested", "idx", idx)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(parts[idx]))
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
