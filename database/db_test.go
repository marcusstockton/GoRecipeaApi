package database

import (
	"strings"
	"testing"

	"gorm.io/gorm/logger"
)

func TestDefaultDSNUsesPostgresSettings(t *testing.T) {
	dsn := defaultDSN()

	for _, want := range []string{
		"host=localhost",
		"user=recipea",
		"password=recipea",
		"dbname=recipea_dev",
		"port=5432",
		"sslmode=disable",
	} {
		if !strings.Contains(dsn, want) {
			t.Fatalf("default DSN missing %q in %q", want, dsn)
		}
	}
}

func TestGormLogModeFromAppLogLevel(t *testing.T) {
	cases := map[string]logger.LogLevel{
		"debug": logger.Info,
		"info":  logger.Warn,
		"warn":  logger.Warn,
		"error": logger.Error,
		"":      logger.Warn,
	}

	for input, want := range cases {
		if got := gormLogModeFromAppLogLevel(input); got != want {
			t.Fatalf("gormLogModeFromAppLogLevel(%q) = %v, want %v", input, got, want)
		}
	}
}
