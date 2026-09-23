package config

import (
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppEnv       string
	Host         string
	Port         string
	JWTSecret    string
	DatabasePath string
	GinMode      string
	LogLevel     string
}

func Load() Config {
	appEnv := strings.TrimSpace(strings.ToLower(envOrDefault("APP_ENV", "dev")))
	if appEnv == "" {
		appEnv = "dev"
	}

	return Config{
		AppEnv:       appEnv,
		Host:         envOrDefault("HOST", "0.0.0.0"),
		Port:         envOrDefault("PORT", "8080"),
		JWTSecret:    envOrDefault("JWT_SECRET", "change-me-in-production"),
		DatabasePath: resolveDatabasePath(envOrDefault("DATABASE_PATH", "database/database.db")),
		GinMode:      envOrDefault("GIN_MODE", defaultGinMode(appEnv)),
		LogLevel:     envOrDefault("LOG_LEVEL", defaultLogLevel(appEnv)),
	}
}

func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "prod"
}

func defaultGinMode(appEnv string) string {
	switch appEnv {
	case "prod", "uat":
		return "release"
	default:
		return "debug"
	}
}

func defaultLogLevel(appEnv string) string {
	switch appEnv {
	case "prod":
		return "warn"
	case "uat":
		return "info"
	default:
		return "debug"
	}
}

func resolveDatabasePath(path string) string {
	if path == "" {
		path = "database/database.db"
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	root, err := findProjectRoot()
	if err != nil {
		return filepath.Clean(path)
	}

	return filepath.Join(root, filepath.Clean(path))
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", os.ErrNotExist
		}
		wd = parent
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
