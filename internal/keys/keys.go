package keys

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
)

const (
	DefaultKeysDir  = ".envcrypt"
	DefaultKeysFile = "keys.txt"
)

// Recipient wraps an age recipient with an optional alias.
type Recipient struct {
	Alias     string
	PublicKey string
	recipient age.Recipient
}

// LoadRecipients reads public keys from a keys file.
// Each line may be in the format: "alias=age1..." or just "age1...".
func LoadRecipients(path string) ([]Recipient, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var recipients []Recipient
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		alias := ""
		pubKey := line
		if idx := strings.Index(line, "="); idx != -1 {
			alias = line[:idx]
			pubKey = line[idx+1:]
		}

		r, err := age.ParseX25519Recipient(pubKey)
		if err != nil {
			return nil, err
		}
		recipients = append(recipients, Recipient{
			Alias:     alias,
			PublicKey: pubKey,
			recipient: r,
		})
	}

	if len(recipients) == 0 {
		return nil, errors.New("no valid recipients found in keys file")
	}
	return recipients, nil
}

// AgeRecipients returns the underlying age.Recipient slice.
func AgeRecipients(rs []Recipient) []age.Recipient {
	out := make([]age.Recipient, len(rs))
	for i, r := range rs {
		out[i] = r.recipient
	}
	return out
}

// DefaultKeysPath returns the default path to the keys file.
func DefaultKeysPath() string {
	return filepath.Join(DefaultKeysDir, DefaultKeysFile)
}

// GenerateIdentity creates a new age X25519 identity and returns it.
func GenerateIdentity() (*age.X25519Identity, error) {
	return age.GenerateX25519Identity()
}
