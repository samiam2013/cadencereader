package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

type AppEnvironment uint8

const (
	_                    = iota
	local AppEnvironment = iota
	prod  AppEnvironment = iota
)

type Config struct {
	AppEnv          AppEnvironment
	DatabaseURL     *url.URL
	MigrationFolder string
	ViewFolder      string
}

func Load() (Config, error) {
	getEnv := func(key string) string {
		val := os.Getenv(key)
		if len(strings.TrimSpace(val)) == 0 {
			slog.Error("Env var empty or all whitespace", "key", key)
			os.Exit(1)
		}
		return val
	}

	ae := getEnv("APP_ENV")
	var appEnv AppEnvironment
	switch ae {
	case "local":
		appEnv = local
	case "prod":
		appEnv = prod
	default:
		slog.Error("Failed to parse app environment", "APP_ENV", ae)
		os.Exit(1)
	}

	dbu := getEnv("DATABASE_URL")
	dbURL, err := url.Parse(dbu)
	if err != nil {
		return Config{}, fmt.Errorf("failed parsing database url: %w", err)
	}

	migFold := getEnv("MIGRATION_FOLDER")
	viewFold := getEnv("VIEW_FOLDER")

	return Config{
		AppEnv:          appEnv,
		DatabaseURL:     dbURL,
		MigrationFolder: migFold,
		ViewFolder:      viewFold,
	}, nil

}
