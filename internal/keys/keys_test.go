package keys_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envcrypt/internal/keys"
)

func TestLoadRecipients_ValidFile(t *testing.T) {
	// Generate a real identity to get a valid public key.
	id, err := keys.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity: %v", err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "keys.txt")
	content := "alice=" + id.Recipient().String() + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	recipients, err := keys.LoadRecipients(path)
	if err != nil {
		t.Fatalf("LoadRecipients: %v", err)
	}
	if len(recipients) != 1 {
		t.Fatalf("expected 1 recipient, got %d", len(recipients))
	}
	if recipients[0].Alias != "alice" {
		t.Errorf("expected alias 'alice', got %q", recipients[0].Alias)
	}
	if recipients[0].PublicKey != id.Recipient().String() {
		t.Errorf("public key mismatch")
	}
}

func TestLoadRecipients_SkipsComments(t *testing.T) {
	id, _ := keys.GenerateIdentity()
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.txt")
	content := "# this is a comment\n" + id.Recipient().String() + "\n"
	os.WriteFile(path, []byte(content), 0600)

	rs, err := keys.LoadRecipients(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs) != 1 {
		t.Errorf("expected 1 recipient, got %d", len(rs))
	}
}

func TestLoadRecipients_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.txt")
	os.WriteFile(path, []byte("# only comments\n"), 0600)

	_, err := keys.LoadRecipients(path)
	if err == nil {
		t.Error("expected error for empty keys file, got nil")
	}
}

func TestLoadRecipients_MissingFile(t *testing.T) {
	_, err := keys.LoadRecipients("/nonexistent/path/keys.txt")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestAgeRecipients(t *testing.T) {
	id, _ := keys.GenerateIdentity()
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.txt")
	os.WriteFile(path, []byte(id.Recipient().String()+"\n"), 0600)

	rs, _ := keys.LoadRecipients(path)
	ageRs := keys.AgeRecipients(rs)
	if len(ageRs) != len(rs) {
		t.Errorf("length mismatch: %d vs %d", len(ageRs), len(rs))
	}
}
