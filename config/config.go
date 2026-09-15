package config

import (
	"net"
	"os"
)

type Config struct {
	Host         string
	Port         string
	JWTSecret    string
	DatabasePath string
	GinMode      string
}

func Load() Config {
	return Config{
		Host:         envOrDefault("HOST", "0.0.0.0"),
		Port:         envOrDefault("PORT", "8080"),
		JWTSecret:    envOrDefault("JWT_SECRET", "change-me-in-production"),
		DatabasePath: envOrDefault("DATABASE_PATH", "database/database.db"),
		GinMode:      envOrDefault("GIN_MODE", "debug"),
	}
}

func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
