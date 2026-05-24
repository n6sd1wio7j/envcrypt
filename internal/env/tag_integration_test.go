package env

import (
	"path/filepath"
	"testing"
)

// TestTagLifecycle exercises the full add → save → load → find → remove → save cycle.
func TestTagLifecycle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envcrypt", "tags.json")

	// Start with an empty tag file.
	tf, err := LoadTagFile(path)
	if err != nil {
		t.Fatalf("LoadTagFile: %v", err)
	}

	// Add two tags.
	for _, name := range []string{"v1", "v2"} {
		if err := AddTag(&tf, Tag{Name: name, SnapshotRef: "snap_" + name}); err != nil {
			t.Fatalf("AddTag %q: %v", name, err)
		}
	}

	if err := SaveTagFile(path, tf); err != nil {
		t.Fatalf("SaveTagFile: %v", err)
	}

	// Reload and verify.
	tf2, err := LoadTagFile(path)
	if err != nil {
		t.Fatalf("LoadTagFile after save: %v", err)
	}
	if len(tf2.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tf2.Tags))
	}

	// Find a specific tag.
	tag, err := FindTag(tf2, "v1")
	if err != nil {
		t.Fatalf("FindTag: %v", err)
	}
	if tag.SnapshotRef != "snap_v1" {
		t.Errorf("unexpected ref: %q", tag.SnapshotRef)
	}

	// Remove a tag and persist.
	if err := RemoveTag(&tf2, "v1"); err != nil {
		t.Fatalf("RemoveTag: %v", err)
	}
	if err := SaveTagFile(path, tf2); err != nil {
		t.Fatalf("SaveTagFile after remove: %v", err)
	}

	// Reload and confirm only v2 remains.
	tf3, err := LoadTagFile(path)
	if err != nil {
		t.Fatalf("LoadTagFile final: %v", err)
	}
	if len(tf3.Tags) != 1 || tf3.Tags[0].Name != "v2" {
		t.Errorf("expected only 'v2', got %+v", tf3.Tags)
	}

	// Ensure v1 is gone.
	if _, err := FindTag(tf3, "v1"); err == nil {
		t.Error("expected error finding removed tag 'v1'")
	}
}

// TestAddTag_SetsTimestampAutomatically ensures CreatedAt is populated when zero.
func TestAddTag_SetsTimestampAutomatically(t *testing.T) {
	tf := TagFile{}
	_ = AddTag(&tf, Tag{Name: "auto-ts", SnapshotRef: "snap_x"})
	if tf.Tags[0].CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be auto-set")
	}
}
