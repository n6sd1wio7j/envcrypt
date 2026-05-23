package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTakeSnapshot_ValidFile(t *testing.T) {
	f := writeTemp(t, "KEY1=val1\nKEY2=val2\n")
	snap, err := TakeSnapshot(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap.Entries))
	}
	if snap.Path != f {
		t.Errorf("expected path %q, got %q", f, snap.Path)
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTakeSnapshot_MissingFile(t *testing.T) {
	_, err := TakeSnapshot("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoadSnapshot_Roundtrip(t *testing.T) {
	f := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	snap, err := TakeSnapshot(f)
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	dir := t.TempDir()
	savedPath, err := SaveSnapshot(snap, dir)
	if err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	if _, err := os.Stat(savedPath); err != nil {
		t.Fatalf("saved file not found: %v", err)
	}

	loaded, err := LoadSnapshot(savedPath)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if len(loaded.Entries) != len(snap.Entries) {
		t.Fatalf("entry count mismatch: want %d, got %d", len(snap.Entries), len(loaded.Entries))
	}
	for i, e := range loaded.Entries {
		if e.Key != snap.Entries[i].Key || e.Value != snap.Entries[i].Value {
			t.Errorf("entry %d mismatch: want %v, got %v", i, snap.Entries[i], e)
		}
	}
}

func TestTimestampFromPath_Valid(t *testing.T) {
	name := "snapshot_20240115T120000Z.env"
	ts, err := timestampFromPath(name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts)
	}
}

func TestSaveSnapshot_CreatesDir(t *testing.T) {
	f := writeTemp(t, "A=1\n")
	snap, _ := TakeSnapshot(f)
	dir := filepath.Join(t.TempDir(), "nested", "snapshots")
	_, err := SaveSnapshot(snap, dir)
	if err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}
