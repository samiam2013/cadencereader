package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func Open(dbURL *url.URL) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database after opening: %w", err)
	}

	m, err := openMigration(db)
	if err != nil {
		return nil, fmt.Errorf("failed to open migrations: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("failed to run migration: %w", err)
		}
	}
	return db, nil
}

func openMigration(db *sql.DB) (*migrate.Migrate, error) {
	instance, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return nil, err
	}

	// TODO: this path is for the docker continer?
	fSrc, err := (&file.File{}).Open(os.Getenv("MIGRATION_FOLDER"))
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		return nil, err
	}
	return m, nil
}
