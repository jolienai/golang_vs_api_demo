package app

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	LogLevel      string
	AppEnv        string
	MigrationsDir string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		HTTPAddr:      envOrDefault("HTTP_ADDR", ":8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		LogLevel:      envOrDefault("LOG_LEVEL", "info"),
		AppEnv:        envOrDefault("APP_ENV", "local"),
		MigrationsDir: envOrDefault("MIGRATIONS_DIR", "migrations"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
