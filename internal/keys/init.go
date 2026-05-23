package keys

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitResult holds the output of initializing a new key pair.
type InitResult struct {
	PublicKey      string
	PrivateKeyPath string
	KeysFilePath   string
}

// Init generates a new age identity, saves the private key to disk,
// and creates the keys file with the public key pre-populated.
// alias is an optional label for the key (e.g. the user's name).
func Init(alias string) (*InitResult, error) {
	id, err := GenerateIdentity()
	if err != nil {
		return nil, fmt.Errorf("generate identity: %w", err)
	}

	if err := os.MkdirAll(DefaultKeysDir, 0700); err != nil {
		return nil, fmt.Errorf("create keys dir: %w", err)
	}

	// Write private key.
	privPath := filepath.Join(DefaultKeysDir, "identity.txt")
	privContent := fmt.Sprintf("# age secret key — keep this private, do NOT commit\n%s\n", id.String())
	if err := os.WriteFile(privPath, []byte(privContent), 0600); err != nil {
		return nil, fmt.Errorf("write identity: %w", err)
	}

	// Write public key to keys file.
	keysPath := DefaultKeysPath()
	var line string
	if alias != "" {
		line = fmt.Sprintf("%s=%s\n", alias, id.Recipient().String())
	} else {
		line = id.Recipient().String() + "\n"
	}

	keysContent := fmt.Sprintf("# envcrypt team keys — commit this file\n%s", line)
	if err := os.WriteFile(keysPath, []byte(keysContent), 0644); err != nil {
		return nil, fmt.Errorf("write keys file: %w", err)
	}

	return &InitResult{
		PublicKey:      id.Recipient().String(),
		PrivateKeyPath: privPath,
		KeysFilePath:   keysPath,
	}, nil
}
