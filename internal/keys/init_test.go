package keys_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envcrypt/internal/keys"
)

func TestInit_CreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	// Override the default keys dir for this test.
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	result, err := keys.Init("testuser")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if result.PublicKey == "" {
		t.Error("expected non-empty public key")
	}
	if !strings.HasPrefix(result.PublicKey, "age1") {
		t.Errorf("public key should start with 'age1', got %q", result.PublicKey)
	}

	// Check private key file exists and is mode 0600.
	info, err := os.Stat(result.PrivateKeyPath)
	if err != nil {
		t.Fatalf("private key file missing: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("private key perm: got %o, want 0600", info.Mode().Perm())
	}

	// Check keys file contains alias and public key.
	data, err := os.ReadFile(result.KeysFilePath)
	if err != nil {
		t.Fatalf("keys file missing: %v", err)
	}
	if !strings.Contains(string(data), "testuser="+result.PublicKey) {
		t.Errorf("keys file missing expected entry; content:\n%s", data)
	}
}

func TestInit_NoAlias(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(origDir) })

	result, err := keys.Init("")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	data, _ := os.ReadFile(result.KeysFilePath)
	if !strings.Contains(string(data), result.PublicKey) {
		t.Errorf("keys file should contain bare public key")
	}
	if strings.Contains(string(data), "="+result.PublicKey) {
		t.Errorf("keys file should not contain alias prefix when alias is empty")
	}
}
