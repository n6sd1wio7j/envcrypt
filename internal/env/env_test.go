package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return p
}

func TestParse_ValidFile(t *testing.T) {
	p := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "FOO" || ef.Entries[0].Value != "bar" {
		t.Errorf("unexpected first entry: %+v", ef.Entries[0])
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	p := writeTemp(t, "# comment\n\nKEY=value\n")
	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "KEY" {
		t.Errorf("unexpected key: %s", ef.Entries[0].Key)
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParse_EmptyFile(t *testing.T) {
	p := writeTemp(t, "")
	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(ef.Entries))
	}
}

func TestWriteRoundtrip(t *testing.T) {
	orig := "HOST=localhost\nPORT=5432\n"
	p := writeTemp(t, orig)

	ef, err := Parse(p)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	out := filepath.Join(t.TempDir(), ".env.out")
	if err := Write(out, ef); err != nil {
		t.Fatalf("write error: %v", err)
	}

	ef2, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parse error: %v", err)
	}
	if len(ef2.Entries) != len(ef.Entries) {
		t.Errorf("entry count mismatch: %d vs %d", len(ef2.Entries), len(ef.Entries))
	}
}
