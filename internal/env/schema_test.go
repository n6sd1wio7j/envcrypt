package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSchemaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	return p
}

func TestLoadSchema_Valid(t *testing.T) {
	path := writeSchemaFile(t, "DB_URL:\n  required: true\n  description: database connection URL\n")
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s["DB_URL"]; !ok {
		t.Error("expected DB_URL in schema")
	}
	if !s["DB_URL"].Required {
		t.Error("expected DB_URL to be required")
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSchema_InvalidYAML(t *testing.T) {
	path := writeSchemaFile(t, ": invalid: [yaml")
	_, err := LoadSchema(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestValidateAgainstSchema_NoViolations(t *testing.T) {
	schema := Schema{
		"DB_URL": {Required: true},
		"PORT":   {Required: false},
	}
	entries := []Entry{
		{Key: "DB_URL", Value: "postgres://localhost/db"},
		{Key: "PORT", Value: "5432"},
	}
	violations := ValidateAgainstSchema(entries, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidateAgainstSchema_MissingRequired(t *testing.T) {
	schema := Schema{
		"DB_URL": {Required: true},
		"SECRET": {Required: true},
	}
	entries := []Entry{
		{Key: "DB_URL", Value: "postgres://localhost/db"},
	}
	violations := ValidateAgainstSchema(entries, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET" {
		t.Errorf("expected violation for SECRET, got %s", violations[0].Key)
	}
}

func TestValidateAgainstSchema_DefaultSatisfiesRequired(t *testing.T) {
	schema := Schema{
		"LOG_LEVEL": {Required: true, Default: "info"},
	}
	entries := []Entry{}
	violations := ValidateAgainstSchema(entries, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations when default is set, got %v", violations)
	}
}
