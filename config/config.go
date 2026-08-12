package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type AppEnvironment uint8

const (
	_ AppEnvironment = iota
	local
	prod
)

type Config struct {
	AppEnv          AppEnvironment
	HTTPport        uint16
	DatabaseURL     *url.URL
	MigrationFolder string
	Extra           map[string]string
}

func Load(extraEnv ...string) (Config, error) {
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

	hp := getEnv("HTTP_PORT")
	httpPort, err := strconv.ParseUint(hp, 10, 16)
	if err != nil {
		slog.Error("Failed to parse port number", "error", err)
		os.Exit(1)
	}

	dbu := getEnv("DATABASE_URL")
	dbURL, err := url.Parse(dbu)
	if err != nil {
		return Config{}, fmt.Errorf("failed parsing database url: %w", err)
	}

	migrationFolder := getEnv("MIGRATION_FOLDER")

	extraEnvMap := make(map[string]string)
	for _, key := range extraEnv {
		value := getEnv(key)
		extraEnvMap[key] = value
	}

	return Config{
		AppEnv:          appEnv,
		HTTPport:        uint16(httpPort),
		DatabaseURL:     dbURL,
		MigrationFolder: migrationFolder,
		Extra:           extraEnvMap,
	}, nil
}
