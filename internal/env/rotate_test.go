package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotate_CreatesBackupAndWritesNew(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, ".env.age")

	original := []byte("original encrypted content")
	if err := os.WriteFile(encPath, original, 0o600); err != nil {
		t.Fatal(err)
	}

	newContent := []byte("new encrypted content")
	opts := RotateOptions{Suffix: "20240101T000000Z"}

	if err := Rotate(encPath, newContent, opts); err != nil {
		t.Fatalf("Rotate returned error: %v", err)
	}

	// New content should be in place.
	got, err := os.ReadFile(encPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(newContent) {
		t.Errorf("expected new content %q, got %q", newContent, got)
	}

	// Backup should exist.
	backupPath := filepath.Join(dir, "backups", ".env.age.20240101T000000Z")
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("backup file missing: %v", err)
	}
	if string(backupData) != string(original) {
		t.Errorf("backup content mismatch: got %q", backupData)
	}
}

func TestRotate_NoExistingFile(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, ".env.age")

	newContent := []byte("brand new")
	opts := RotateOptions{Suffix: "ts"}

	if err := Rotate(encPath, newContent, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := os.ReadFile(encPath)
	if string(got) != string(newContent) {
		t.Errorf("expected %q, got %q", newContent, got)
	}

	backups, _ := ListBackups(encPath, "")
	if len(backups) != 0 {
		t.Errorf("expected no backups, got %v", backups)
	}
}

func TestRotate_CustomBackupDir(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, ".env.age")
	customBackup := filepath.Join(dir, "archive")

	_ = os.WriteFile(encPath, []byte("old"), 0o600)

	opts := RotateOptions{BackupDir: customBackup, Suffix: "v1"}
	if err := Rotate(encPath, []byte("new"), opts); err != nil {
		t.Fatal(err)
	}

	backups, err := ListBackups(encPath, customBackup)
	if err != nil {
		t.Fatal(err)
	}
	if len(backups) != 1 || !strings.HasSuffix(backups[0], ".v1") {
		t.Errorf("unexpected backups: %v", backups)
	}
}

func TestListBackups_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, ".env.age")
	backupDir := filepath.Join(dir, "backups")
	_ = os.MkdirAll(backupDir, 0o700)

	for _, suffix := range []string{"ts1", "ts2", "ts3"} {
		_ = os.WriteFile(filepath.Join(backupDir, ".env.age."+suffix), []byte(suffix), 0o600)
	}

	backups, err := ListBackups(encPath, backupDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(backups) != 3 {
		t.Errorf("expected 3 backups, got %d", len(backups))
	}
}
