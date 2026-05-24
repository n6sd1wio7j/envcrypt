package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImport_DotenvIntoDst(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.env")
	dst := filepath.Join(dir, "dest.env")

	_ = os.WriteFile(src, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	entries, err := Import(src, dst, ImportOptions{Format: ImportDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestImport_ShellFormat(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.sh")
	dst := filepath.Join(dir, "dest.env")

	_ = os.WriteFile(src, []byte("export KEY=value\nexport OTHER=123\n"), 0600)

	entries, err := Import(src, dst, ImportOptions{Format: ImportShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestImport_OverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.env")
	dst := filepath.Join(dir, "dest.env")

	_ = os.WriteFile(dst, []byte("FOO=old\n"), 0600)
	_ = os.WriteFile(src, []byte("FOO=new\n"), 0600)

	entries, err := Import(src, dst, ImportOptions{Format: ImportDotenv, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if e.Key == "FOO" && e.Value != "new" {
			t.Errorf("expected FOO=new, got FOO=%s", e.Value)
		}
	}
}

func TestImport_SkipExistingWithoutOverwrite(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.env")
	dst := filepath.Join(dir, "dest.env")

	_ = os.WriteFile(dst, []byte("FOO=old\n"), 0600)
	_ = os.WriteFile(src, []byte("FOO=new\n"), 0600)

	entries, err := Import(src, dst, ImportOptions{Format: ImportDotenv, Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if e.Key == "FOO" && e.Value != "old" {
			t.Errorf("expected FOO=old (skip), got FOO=%s", e.Value)
		}
	}
}

func TestImport_KeyFilter(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.env")
	dst := filepath.Join(dir, "dest.env")

	_ = os.WriteFile(src, []byte("A=1\nB=2\nC=3\n"), 0600)

	entries, err := Import(src, dst, ImportOptions{Format: ImportDotenv, Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 filtered entries, got %d", len(entries))
	}
}

func TestImport_MissingSource(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "dest.env")

	_, err := Import(filepath.Join(dir, "nonexistent.env"), dst, ImportOptions{})
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}
