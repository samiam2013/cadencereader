package main

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func index(w http.ResponseWriter, r *http.Request) {
	fh, err := os.Open("/views/index.html")
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

	dbURL := os.Getenv("DATABASE_URL")
	if len(strings.TrimSpace(dbURL)) == 0 {
		slog.Error("DATABASE_URL empty and required")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Could not connect to database", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		slog.Error("DB ping failed after opening", "error", err)
		os.Exit(1)
	}

	m, err := OpenMigration(db)
	if err != nil {
		slog.Error("Failed to open migrations", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("No migration changes")
		} else {
			slog.Error("Failed during up migration", "error", err)
			os.Exit(1)
		}
	}

	http.HandleFunc("/health", healthCheck(db))

	http.HandleFunc("/", index)

	slog.Info("server starting")
	if err := http.ListenAndServe(":80", nil); err != nil {
		slog.Error("HTTP server exit", "error", err)
		os.Exit(1)
	}
}

func OpenMigration(db *sql.DB) (*migrate.Migrate, error) {

	instance, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return nil, err
	}

	// TODO: this path is for the docker continer?
	fSrc, err := (&file.File{}).Open("/migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		return nil, err
	}
	return m, nil
}
