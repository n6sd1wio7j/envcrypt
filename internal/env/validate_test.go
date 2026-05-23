package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateEntries_Valid(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "PORT", Value: "8080"},
	}
	result := ValidateEntries(entries)
	if !result.OK() {
		t.Fatalf("expected no errors, got: %s", result.Error())
	}
}

func TestValidateEntries_EmptyKey(t *testing.T) {
	entries := []Entry{
		{Key: "", Value: "value"},
	}
	result := ValidateEntries(entries)
	if result.OK() {
		t.Fatal("expected error for empty key")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
}

func TestValidateEntries_InvalidKeyFormat(t *testing.T) {
	cases := []string{"lower_case", "123START", "has-hyphen", "has space"}
	for _, k := range cases {
		entries := []Entry{{Key: k, Value: "v"}}
		result := ValidateEntries(entries)
		if result.OK() {
			t.Errorf("expected error for key %q", k)
		}
	}
}

func TestValidateEntries_DuplicateKey(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "FOO", Value: "baz"},
	}
	result := ValidateEntries(entries)
	if result.OK() {
		t.Fatal("expected duplicate key error")
	}
	found := false
	for _, e := range result.Errors {
		if e.Key == "FOO" && e.Message == "duplicate key" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate key error for FOO, got: %s", result.Error())
	}
}

func TestValidateFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "APP_NAME=envcrypt\nDEBUG=false\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	result, err := ValidateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.OK() {
		t.Fatalf("expected no validation errors, got: %s", result.Error())
	}
}

func TestValidateFile_MissingFile(t *testing.T) {
	_, err := ValidateFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidationResult_ErrorString(t *testing.T) {
	r := &ValidationResult{
		Errors: []ValidationError{
			{Key: "bad-key", Message: "key must match pattern [A-Z_][A-Z0-9_]*"},
			{Key: "FOO", Message: "duplicate key"},
		},
	}
	s := r.Error()
	if s == "" {
		t.Fatal("expected non-empty error string")
	}
}
