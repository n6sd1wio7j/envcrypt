package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadTagFile_MissingFile(t *testing.T) {
	tf, err := LoadTagFile("/nonexistent/path/tags.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(tf.Tags) != 0 {
		t.Errorf("expected empty tags, got %d", len(tf.Tags))
	}
}

func TestSaveAndLoadTagFile_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")

	tf := TagFile{
		Tags: []Tag{
			{Name: "v1.0", Description: "initial release", SnapshotRef: "snap_001", CreatedAt: time.Now().UTC().Truncate(time.Second)},
		},
	}
	if err := SaveTagFile(path, tf); err != nil {
		t.Fatalf("SaveTagFile: %v", err)
	}

	loaded, err := LoadTagFile(path)
	if err != nil {
		t.Fatalf("LoadTagFile: %v", err)
	}
	if len(loaded.Tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(loaded.Tags))
	}
	if loaded.Tags[0].Name != "v1.0" {
		t.Errorf("expected name 'v1.0', got %q", loaded.Tags[0].Name)
	}
}

func TestAddTag_Success(t *testing.T) {
	tf := TagFile{}
	err := AddTag(&tf, Tag{Name: "release", SnapshotRef: "snap_002"})
	if err != nil {
		t.Fatalf("AddTag: %v", err)
	}
	if len(tf.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tf.Tags))
	}
	if tf.Tags[0].CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestAddTag_Duplicate(t *testing.T) {
	tf := TagFile{Tags: []Tag{{Name: "v1", SnapshotRef: "snap_001"}}}
	err := AddTag(&tf, Tag{Name: "v1", SnapshotRef: "snap_002"})
	if err == nil {
		t.Fatal("expected error for duplicate tag")
	}
}

func TestRemoveTag_Success(t *testing.T) {
	tf := TagFile{Tags: []Tag{{Name: "v1"}, {Name: "v2"}}}
	if err := RemoveTag(&tf, "v1"); err != nil {
		t.Fatalf("RemoveTag: %v", err)
	}
	if len(tf.Tags) != 1 || tf.Tags[0].Name != "v2" {
		t.Errorf("unexpected tags after remove: %+v", tf.Tags)
	}
}

func TestRemoveTag_NotFound(t *testing.T) {
	tf := TagFile{}
	if err := RemoveTag(&tf, "missing"); err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestFindTag(t *testing.T) {
	tf := TagFile{Tags: []Tag{{Name: "prod", SnapshotRef: "snap_prod"}}}
	tag, err := FindTag(tf, "prod")
	if err != nil {
		t.Fatalf("FindTag: %v", err)
	}
	if tag.SnapshotRef != "snap_prod" {
		t.Errorf("unexpected ref: %q", tag.SnapshotRef)
	}
	_, err = FindTag(tf, "missing")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestSaveTagFile_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "tags.json")
	if err := SaveTagFile(path, TagFile{}); err != nil {
		t.Fatalf("SaveTagFile: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
