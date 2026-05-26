package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLintTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeLintTemp: %v", err)
	}
	return p
}

func TestLint_NoFindings(t *testing.T) {
	p := writeLintTemp(t, "FOO=bar\nBAZ=qux\n")
	findings, err := Lint(p, DefaultLintOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_WarnOnEmptyValue(t *testing.T) {
	p := writeLintTemp(t, "FOO=\nBAR=hello\n")
	findings, err := Lint(p, DefaultLintOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", findings[0].Key)
	}
	if findings[0].Severity != LintWarn {
		t.Errorf("expected warn severity, got %s", findings[0].Severity)
	}
}

func TestLint_WarnOnLowercaseKey(t *testing.T) {
	p := writeLintTemp(t, "foo=bar\nBAZ=qux\n")
	findings, err := Lint(p, DefaultLintOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "foo" {
		t.Errorf("expected key foo, got %s", findings[0].Key)
	}
}

func TestLint_ErrorOnDuplicateKey(t *testing.T) {
	p := writeLintTemp(t, "FOO=first\nBAR=baz\nFOO=second\n")
	findings, err := Lint(p, DefaultLintOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var dupFound bool
	for _, f := range findings {
		if f.Key == "FOO" && f.Severity == LintError {
			dupFound = true
		}
	}
	if !dupFound {
		t.Errorf("expected duplicate-key error finding for FOO, got: %v", findings)
	}
}

func TestLint_MissingFile(t *testing.T) {
	_, err := Lint("/nonexistent/.env", DefaultLintOptions())
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLint_DisabledChecks(t *testing.T) {
	p := writeLintTemp(t, "foo=\nfoo=dup\n")
	opts := LintOptions{
		WarnOnNoValue:    false,
		WarnOnLowercase:  false,
		ErrorOnDuplicate: false,
	}
	findings, err := Lint(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings with all checks disabled, got %d", len(findings))
	}
}
