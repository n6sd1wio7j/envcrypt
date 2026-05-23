package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadLockFile_MissingFile(t *testing.T) {
	lf, err := LoadLockFile("/nonexistent/path/.envcrypt.lock")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if lf.Version != 1 {
		t.Errorf("expected version 1, got %d", lf.Version)
	}
	if len(lf.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(lf.Entries))
	}
}

func TestSaveAndLoadLockFile_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envcrypt.lock")

	lf := &LockFile{
		Version: 1,
		Entries: make(map[string]LockEntry),
	}
	entry := LockEntry{
		File:      ".env",
		Encrypted: ".env.age",
		Checksum:  "abc123",
		UpdatedAt: time.Now().UTC().Truncate(time.Second),
		UpdatedBy: "alice",
	}
	lf.Upsert(entry)

	if err := SaveLockFile(path, lf); err != nil {
		t.Fatalf("SaveLockFile: %v", err)
	}

	loaded, err := LoadLockFile(path)
	if err != nil {
		t.Fatalf("LoadLockFile: %v", err)
	}
	got, ok := loaded.Entries[".env"]
	if !ok {
		t.Fatal("expected entry for .env")
	}
	if got.Checksum != "abc123" {
		t.Errorf("checksum mismatch: got %q", got.Checksum)
	}
	if got.UpdatedBy != "alice" {
		t.Errorf("updatedBy mismatch: got %q", got.UpdatedBy)
	}
}

func TestLockFile_Upsert_And_Remove(t *testing.T) {
	lf := &LockFile{Version: 1, Entries: make(map[string]LockEntry)}
	lf.Upsert(LockEntry{File: ".env", Encrypted: ".env.age"})
	if _, ok := lf.Entries[".env"]; !ok {
		t.Error("expected entry after upsert")
	}
	lf.Remove(".env")
	if _, ok := lf.Entries[".env"]; ok {
		t.Error("expected entry removed")
	}
}

func TestSaveLockFile_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "nested", ".envcrypt.lock")
	lf := &LockFile{Version: 1, Entries: make(map[string]LockEntry)}
	if err := SaveLockFile(path, lf); err != nil {
		t.Fatalf("expected dir creation, got %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("lock file not created: %v", err)
	}
}
