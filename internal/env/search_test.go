package env

import (
	"path/filepath"
	"testing"
)

func TestSearch_AllEntries(t *testing.T) {
	f := writeTemp(t, "KEY1=val1\nKEY2=val2\n")
	res, err := Search(f, SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestSearch_KeyPattern(t *testing.T) {
	f := writeTemp(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")
	res, err := Search(f, SearchOptions{KeyPattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	for _, r := range res {
		if r.Entry.Key == "APP_NAME" {
			t.Error("APP_NAME should not be in results")
		}
	}
}

func TestSearch_ValuePattern(t *testing.T) {
	f := writeTemp(t, "KEY1=production\nKEY2=staging\nKEY3=production_db\n")
	res, err := Search(f, SearchOptions{ValuePattern: "^production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	f := writeTemp(t, "SECRET_KEY=abc\nother=def\n")
	res, err := Search(f, SearchOptions{KeyPattern: "secret", CaseSensitive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Entry.Key != "SECRET_KEY" {
		t.Fatalf("expected SECRET_KEY, got %+v", res)
	}
}

func TestSearch_CaseSensitiveNoMatch(t *testing.T) {
	f := writeTemp(t, "SECRET_KEY=abc\n")
	res, err := Search(f, SearchOptions{KeyPattern: "secret", CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 results, got %d", len(res))
	}
}

func TestSearch_MissingFile(t *testing.T) {
	_, err := Search(filepath.Join(t.TempDir(), "missing.env"), SearchOptions{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSearch_InvalidPattern(t *testing.T) {
	f := writeTemp(t, "KEY=val\n")
	_, err := Search(f, SearchOptions{KeyPattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSearchMultiple_AggregatesResults(t *testing.T) {
	f1 := writeTemp(t, "DB_HOST=localhost\n")
	f2 := writeTemp(t, "DB_PORT=5432\nAPP=web\n")
	res, err := SearchMultiple([]string{f1, f2}, SearchOptions{KeyPattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestFormatSearchResults(t *testing.T) {
	results := []SearchResult{
		{Entry: Entry{Key: "KEY", Value: "val"}, File: ".env"},
	}
	out := FormatSearchResults(results)
	if out != ".env: KEY=val\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}
