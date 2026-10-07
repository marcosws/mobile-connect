package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractApkMetadataHandlesNonApkFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(path, []byte("not an apk"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	meta, err := extractApkMetadata(path)
	if err != nil {
		t.Fatalf("expected metadata extraction to tolerate invalid input, got error: %v", err)
	}

	if meta == nil {
		t.Fatal("expected metadata object, got nil")
	}
}
