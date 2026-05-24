package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPinFile_MissingFile(t *testing.T) {
	pf, err := LoadPinFile("/nonexistent/path/pins.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(pf.Pins) != 0 {
		t.Errorf("expected empty pins, got %d", len(pf.Pins))
	}
}

func TestSaveAndLoadPinFile_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envcrypt", "pins.json")

	pf := &PinFile{}
	AddPin(pf, "production", "/snapshots/prod-001.env", "alice")

	if err := SavePinFile(path, pf); err != nil {
		t.Fatalf("SavePinFile: %v", err)
	}

	loaded, err := LoadPinFile(path)
	if err != nil {
		t.Fatalf("LoadPinFile: %v", err)
	}
	if len(loaded.Pins) != 1 {
		t.Fatalf("expected 1 pin, got %d", len(loaded.Pins))
	}
	if loaded.Pins[0].Name != "production" {
		t.Errorf("expected name 'production', got %q", loaded.Pins[0].Name)
	}
	if loaded.Pins[0].PinnedBy != "alice" {
		t.Errorf("expected pinnedBy 'alice', got %q", loaded.Pins[0].PinnedBy)
	}
}

func TestAddPin_ReplacesExisting(t *testing.T) {
	pf := &PinFile{}
	AddPin(pf, "staging", "/snap/v1", "bob")
	AddPin(pf, "staging", "/snap/v2", "carol")

	if len(pf.Pins) != 1 {
		t.Fatalf("expected 1 pin after replace, got %d", len(pf.Pins))
	}
	if pf.Pins[0].SnapshotPath != "/snap/v2" {
		t.Errorf("expected updated snapshot path, got %q", pf.Pins[0].SnapshotPath)
	}
}

func TestRemovePin_Success(t *testing.T) {
	pf := &PinFile{}
	AddPin(pf, "dev", "/snap/dev", "")

	removed := RemovePin(pf, "dev")
	if !removed {
		t.Error("expected RemovePin to return true")
	}
	if len(pf.Pins) != 0 {
		t.Errorf("expected 0 pins after removal, got %d", len(pf.Pins))
	}
}

func TestRemovePin_NotFound(t *testing.T) {
	pf := &PinFile{}
	if RemovePin(pf, "ghost") {
		t.Error("expected RemovePin to return false for missing entry")
	}
}

func TestFindPin(t *testing.T) {
	pf := &PinFile{}
	AddPin(pf, "prod", "/snap/prod", "dave")

	e := FindPin(pf, "prod")
	if e == nil {
		t.Fatal("expected to find pin 'prod'")
	}
	if e.SnapshotPath != "/snap/prod" {
		t.Errorf("unexpected snapshot path: %q", e.SnapshotPath)
	}
	if FindPin(pf, "missing") != nil {
		t.Error("expected nil for missing pin")
	}
}

func TestSavePinFile_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "pins.json")

	if err := SavePinFile(path, &PinFile{}); err != nil {
		t.Fatalf("SavePinFile should create directories: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
