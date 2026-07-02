package apps

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveAPKPathPrefersStorageDirectory(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(wd)

	if err := os.MkdirAll(filepath.Join("storage", "apks"), 0o755); err != nil {
		t.Fatalf("mkdir storage/apks: %v", err)
	}
	if err := os.MkdirAll("uploads", 0o755); err != nil {
		t.Fatalf("mkdir uploads: %v", err)
	}

	storagePath := filepath.Join("storage", "apks", "app.apk")
	uploadPath := filepath.Join("uploads", "app.apk")
	if err := os.WriteFile(storagePath, []byte("storage"), 0o644); err != nil {
		t.Fatalf("write storage apk: %v", err)
	}
	if err := os.WriteFile(uploadPath, []byte("upload"), 0o644); err != nil {
		t.Fatalf("write upload apk: %v", err)
	}

	got := resolveAPKPath("app.apk")
	if got != storagePath {
		t.Fatalf("expected %s, got %s", storagePath, got)
	}
}
