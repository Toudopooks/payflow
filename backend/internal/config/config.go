package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port            string
	DatabaseDSN     string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
	DemoMode        bool
}

func Load() (Config, error) {
	timeout, err := time.ParseDuration(env("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse SHUTDOWN_TIMEOUT: %w", err)
	}
	demo, err := strconv.ParseBool(env("DEMO_MODE", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("parse DEMO_MODE: %w", err)
	}
	return Config{
		Port:            env("PORT", "8080"),
		DatabaseDSN:     env("DATABASE_DSN", "payflow:payflow@tcp(localhost:3306)/payflow?parseTime=true&multiStatements=true"),
		AllowedOrigins:  strings.Split(env("ALLOWED_ORIGINS", "http://localhost:5173"), ","),
		ShutdownTimeout: timeout,
		DemoMode:        demo,
	}, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
