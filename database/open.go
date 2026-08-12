package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/somethingsoftware/cadencereader/config"
)

func Open(ctx context.Context, cfg config.Config) (*pgx.Conn, *Queries, error) {
	conn, err := pgx.Connect(ctx, cfg.DatabaseURL.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect via pgx: %w", err)
	}
	queries := New(conn)

	db, err := sql.Open("pgx", cfg.DatabaseURL.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sql/db for migration: %w", err)
	}
	defer db.Close()

	m, err := openMigration(db, cfg)
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

func openMigration(db *sql.DB, cfg config.Config) (*migrate.Migrate, error) {
	instance, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return nil, err
	}

	fSrc, err := (&file.File{}).Open(cfg.MigrationFolder)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		return nil, err
	}
	return m, nil
}
