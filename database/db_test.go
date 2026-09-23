package database

import (
	"path/filepath"
	"testing"
)

func TestInitDBCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "envs", "app.db")

	InitDB(path)
	if DB == nil {
		t.Fatal("expected database connection to be initialized")
	}

	if _, err := DB.DB(); err != nil {
		t.Fatalf("expected database connection to be usable: %v", err)
	}

	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.Close()
	}
}
