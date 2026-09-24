package database

import (
	"strings"
	"testing"
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
