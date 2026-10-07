package apps

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRepositoryFindAllHandlesNullCreatedAt(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE apps (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			original_name TEXT NOT NULL,
			package_name TEXT,
			version_name TEXT,
			version_code TEXT,
			min_sdk INTEGER,
			target_sdk INTEGER,
			size INTEGER,
			sha256 TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO apps (
			id, filename, original_name, package_name, version_name, version_code,
			min_sdk, target_sdk, size, sha256, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, "app-1", "app.apk", "app.apk", "com.example", "1.0", "1", 21, 21, 10, "abc", nil)
	if err != nil {
		t.Fatalf("insert app: %v", err)
	}

	repo := NewRepository(db)
	apps, err := repo.FindAll()
	if err != nil {
		t.Fatalf("find all: %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}

	if !apps[0].CreatedAt.IsZero() {
		t.Fatalf("expected zero time for null created_at, got %v", apps[0].CreatedAt)
	}
}
