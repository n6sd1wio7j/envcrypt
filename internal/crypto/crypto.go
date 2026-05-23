package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// EncryptFile encrypts the plaintext file at src to dst using the provided recipients.
func EncryptFile(src, dst string, recipients []age.Recipient) error {
	plaintext, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)

	w, err := age.Encrypt(armorWriter, recipients...)
	if err != nil {
		return fmt.Errorf("creating age encryptor: %w", err)
	}

	if _, err := w.Write(plaintext); err != nil {
		return fmt.Errorf("encrypting data: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("finalizing encryption: %w", err)
	}

	if err := armorWriter.Close(); err != nil {
		return fmt.Errorf("closing armor writer: %w", err)
	}

	if err := os.WriteFile(dst, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("writing encrypted file: %w", err)
	}

	return nil
}

// DecryptFile decrypts the armored age file at src to dst using the provided identities.
func DecryptFile(src, dst string, identities []age.Identity) error {
	ciphertext, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading encrypted file: %w", err)
	}

	armorReader := armor.NewReader(bytes.NewReader(ciphertext))

	r, err := age.Decrypt(armorReader, identities...)
	if err != nil {
		return fmt.Errorf("decrypting data: %w", err)
	}

	plaintext, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading decrypted data: %w", err)
	}

	if err := os.WriteFile(dst, plaintext, 0600); err != nil {
		return fmt.Errorf("writing decrypted file: %w", err)
	}

	return nil
}
