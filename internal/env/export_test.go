package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExport_DotenvFormat(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	out, err := Export(entries, ExportOptions{Format: FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("unexpected dotenv output: %q", out)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	entries := []Entry{{Key: "API_KEY", Value: "secret"}}
	out, err := Export(entries, ExportOptions{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export API_KEY=") {
		t.Errorf("expected shell export, got: %q", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	entries := []Entry{{Key: "HOST", Value: "localhost"}}
	out, err := Export(entries, ExportOptions{Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"HOST"`) || !strings.Contains(out, `"localhost"`) {
		t.Errorf("unexpected JSON output: %q", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	_, err := Export([]Entry{}, ExportOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExport_KeyFilter(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}
	out, err := Export(entries, ExportOptions{Format: FormatDotenv, Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "B=") {
		t.Errorf("key B should have been filtered out: %q", out)
	}
	if !strings.Contains(out, "A=1") || !strings.Contains(out, "C=3") {
		t.Errorf("expected A and C in output: %q", out)
	}
}

func TestExportToFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.env")
	entries := []Entry{{Key: "X", Value: "42"}}
	if err := ExportToFile(entries, ExportOptions{}, dest); err != nil {
		t.Fatalf("ExportToFile error: %v", err)
	}
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !strings.Contains(string(data), "X=42") {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestExport_DefaultFormatIsDotenv(t *testing.T) {
	entries := []Entry{{Key: "Z", Value: "99"}}
	out, err := Export(entries, ExportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Z=99") {
		t.Errorf("expected dotenv output by default: %q", out)
	}
}
