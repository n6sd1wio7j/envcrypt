package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemplateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGenerateTemplate_BasicKeys(t *testing.T) {
	p := writeTemplateEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nSECRET=abc\n")
	tmpl, err := GenerateTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(tmpl.Entries))
	}
	if tmpl.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", tmpl.Entries[0].Key)
	}
	for _, e := range tmpl.Entries {
		if !e.Required {
			t.Errorf("expected entry %q to be required", e.Key)
		}
	}
}

func TestGenerateTemplate_InlineComment(t *testing.T) {
	p := writeTemplateEnv(t, "API_KEY=changeme # your API key\n")
	tmpl, err := GenerateTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tmpl.Entries))
	}
	if tmpl.Entries[0].Description != "your API key" {
		t.Errorf("expected description 'your API key', got %q", tmpl.Entries[0].Description)
	}
}

func TestGenerateTemplate_SkipsComments(t *testing.T) {
	p := writeTemplateEnv(t, "# comment\nKEY=val\n")
	tmpl, err := GenerateTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tmpl.Entries))
	}
}

func TestGenerateTemplate_MissingFile(t *testing.T) {
	_, err := GenerateTemplate("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckTemplate_NoMissing(t *testing.T) {
	tmpl := &Template{
		Entries: []TemplateEntry{
			{Key: "A", Required: true},
			{Key: "B", Required: true},
		},
	}
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	missing := CheckTemplate(tmpl, entries)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestCheckTemplate_MissingRequired(t *testing.T) {
	tmpl := &Template{
		Entries: []TemplateEntry{
			{Key: "A", Required: true},
			{Key: "B", Required: true},
		},
	}
	entries := []Entry{{Key: "A", Value: "1"}}
	missing := CheckTemplate(tmpl, entries)
	if len(missing) != 1 || missing[0] != "B" {
		t.Errorf("expected [B] missing, got %v", missing)
	}
}

func TestCheckTemplate_OptionalNotRequired(t *testing.T) {
	tmpl := &Template{
		Entries: []TemplateEntry{
			{Key: "A", Required: true},
			{Key: "OPT", Required: false},
		},
	}
	entries := []Entry{{Key: "A", Value: "1"}}
	missing := CheckTemplate(tmpl, entries)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}
