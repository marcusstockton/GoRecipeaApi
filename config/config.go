package config

import (
	"net"
	"os"
	"strings"
)

type Config struct {
	AppEnv      string
	Host        string
	Port        string
	JWTSecret   string
	DatabaseDSN string
	GinMode     string
	LogLevel    string
}

func Load() Config {
	appEnv := strings.TrimSpace(strings.ToLower(envOrDefault("APP_ENV", "dev")))
	if appEnv == "" {
		appEnv = "dev"
	}

	return Config{
		AppEnv:      appEnv,
		Host:        envOrDefault("HOST", "0.0.0.0"),
		Port:        envOrDefault("PORT", "8080"),
		JWTSecret:   envOrDefault("JWT_SECRET", "change-me-in-production"),
		DatabaseDSN: envOrDefault("DATABASE_DSN", "host=localhost user=recipea password=recipea dbname=recipea_dev port=5432 sslmode=disable"),
		GinMode:     envOrDefault("GIN_MODE", defaultGinMode(appEnv)),
		LogLevel:    envOrDefault("LOG_LEVEL", defaultLogLevel(appEnv)),
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

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
