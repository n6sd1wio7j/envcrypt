package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envcrypt/internal/env"
)

// TestExportRoundtrip parses a .env file and exports it back to dotenv
// format, verifying the round-trip produces equivalent content.
func TestExportRoundtrip(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	_ = os.WriteFile(src, []byte("FOO=bar\nBAZ=qux\n"), 0600)

	entries, err := env.Parse(src)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	out, err := env.Export(entries, env.ExportOptions{Format: env.FormatDotenv})
	if err != nil {
		t.Fatalf("Export: %v", err)
	}

	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("round-trip output missing expected keys: %q", out)
	}
}

// TestExportToFile_PermissionsAreRestricted ensures the exported file
// is written with mode 0600.
func TestExportToFile_PermissionsAreRestricted(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "exported.env")
	entries := []env.Entry{{Key: "SECRET", Value: "shh"}}

	if err := env.ExportToFile(entries, env.ExportOptions{}, dest); err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}
