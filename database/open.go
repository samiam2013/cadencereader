package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

func Open(ctx context.Context, dbURL string) (*pgx.Conn, *Queries, error) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect via pgx: %w", err)
	}
	queries := New(conn)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sql/db for migration: %w", err)
	}
	defer db.Close()

	m, err := openMigration(db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open migrations: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, nil, fmt.Errorf("failed to run migration: %w", err)
		}
	}
	return conn, queries, nil
}

func openMigration(db *sql.DB) (*migrate.Migrate, error) {
	instance, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return nil, err
	}

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
