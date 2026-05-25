package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempRename(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempRename: %v", err)
	}
	return p
}

func TestRename_SingleKey(t *testing.T) {
	p := writeTempRename(t, "FOO=bar\nBAZ=qux\n")
	_, err := Rename(p, [][2]string{{"FOO", "FOO_NEW"}}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := Parse(p)
	if entries[0].Key != "FOO_NEW" {
		t.Errorf("expected FOO_NEW, got %q", entries[0].Key)
	}
	if entries[0].Value != "bar" {
		t.Errorf("value should be preserved, got %q", entries[0].Value)
	}
}

func TestRename_KeyNotFound(t *testing.T) {
	p := writeTempRename(t, "FOO=bar\n")
	_, err := Rename(p, [][2]string{{"MISSING", "NEW"}}, RenameOptions{})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRename_ConflictWithExistingKey(t *testing.T) {
	p := writeTempRename(t, "FOO=1\nBAR=2\n")
	_, err := Rename(p, [][2]string{{"FOO", "BAR"}}, RenameOptions{})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestRename_DryRun_DoesNotModify(t *testing.T) {
	original := "FOO=bar\nBAZ=qux\n"
	p := writeTempRename(t, original)
	results, err := Rename(p, [][2]string{{"FOO", "FOO_NEW"}}, RenameOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Found {
		t.Errorf("expected 1 found result, got %+v", results)
	}
	got, _ := os.ReadFile(p)
	if string(got) != original {
		t.Errorf("dry run should not modify file; got %q", string(got))
	}
}

func TestRename_MultipleKeys(t *testing.T) {
	p := writeTempRename(t, "A=1\nB=2\nC=3\n")
	_, err := Rename(p, [][2]string{{"A", "AA"}, {"C", "CC"}}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := Parse(p)
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	expected := []string{"AA", "B", "CC"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("index %d: expected %q got %q", i, k, keys[i])
		}
	}
}

func TestRename_SameKeyIsNoOp(t *testing.T) {
	p := writeTempRename(t, "FOO=bar\n")
	_, err := Rename(p, [][2]string{{"FOO", "FOO"}}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := Parse(p)
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("same-key rename should be a no-op, got %+v", entries[0])
	}
}
