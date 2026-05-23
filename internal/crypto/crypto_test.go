package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/yourorg/envcrypt/internal/crypto"
)

func generateTestIdentity(t *testing.T) (*age.X25519Identity, *age.X25519Recipient) {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating identity: %v", err)
	}
	return id, id.Recipient()
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, ".env")
	encFile := filepath.Join(dir, ".env.age")
	dstFile := filepath.Join(dir, ".env.decrypted")

	original := []byte("SECRET=hello\nANOTHER=world\n")
	if err := os.WriteFile(srcFile, original, 0600); err != nil {
		t.Fatalf("writing source: %v", err)
	}

	id, recipient := generateTestIdentity(t)

	if err := crypto.EncryptFile(srcFile, encFile, []age.Recipient{recipient}); err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}

	if _, err := os.Stat(encFile); err != nil {
		t.Fatalf("encrypted file not created: %v", err)
	}

	if err := crypto.DecryptFile(encFile, dstFile, []age.Identity{id}); err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}

	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("reading decrypted file: %v", err)
	}

	if string(got) != string(original) {
		t.Errorf("roundtrip mismatch: got %q, want %q", got, original)
	}
}

func TestEncryptFile_MissingSource(t *testing.T) {
	dir := t.TempDir()
	_, recipient := generateTestIdentity(t)

	err := crypto.EncryptFile(
		filepath.Join(dir, "nonexistent.env"),
		filepath.Join(dir, "out.age"),
		[]age.Recipient{recipient},
	)
	if err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}

func TestDecryptFile_WrongIdentity(t *testing.T) {
	dir := t.TempDir()
	srcFile := filepath.Join(dir, ".env")
	encFile := filepath.Join(dir, ".env.age")
	dstFile := filepath.Join(dir, ".env.out")

	if err := os.WriteFile(srcFile, []byte("KEY=value\n"), 0600); err != nil {
		t.Fatalf("writing source: %v", err)
	}

	_, recipient := generateTestIdentity(t)
	wrongID, _ := generateTestIdentity(t)

	if err := crypto.EncryptFile(srcFile, encFile, []age.Recipient{recipient}); err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}

	err := crypto.DecryptFile(encFile, dstFile, []age.Identity{wrongID})
	if err == nil {
		t.Error("expected error when decrypting with wrong identity, got nil")
	}
}
