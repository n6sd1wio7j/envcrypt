package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendAndLoadHistory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	entry := HistoryEntry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Operation: "encrypt",
		File:      ".env",
		User:      "alice",
	}

	if err := AppendHistory(path, entry); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Operation != "encrypt" {
		t.Errorf("expected operation 'encrypt', got %q", entries[0].Operation)
	}
	if entries[0].User != "alice" {
		t.Errorf("expected user 'alice', got %q", entries[0].User)
	}
}

func TestAppendHistory_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	ops := []string{"encrypt", "rotate", "decrypt"}
	for _, op := range ops {
		if err := AppendHistory(path, HistoryEntry{Operation: op, File: ".env"}); err != nil {
			t.Fatalf("AppendHistory(%s): %v", op, err)
		}
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, op := range ops {
		if entries[i].Operation != op {
			t.Errorf("entry %d: expected %q, got %q", i, op, entries[i].Operation)
		}
	}
}

func TestLoadHistory_MissingFile(t *testing.T) {
	entries, err := LoadHistory("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestSaveHistory_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "history.json")

	if err := SaveHistory(path, []HistoryEntry{{Operation: "merge", File: ".env"}}); err != nil {
		t.Fatalf("SaveHistory: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestAppendHistory_SetsTimestamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	before := time.Now().UTC()
	if err := AppendHistory(path, HistoryEntry{Operation: "encrypt", File: ".env"}); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}
	after := time.Now().UTC()

	entries, _ := LoadHistory(path)
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}
