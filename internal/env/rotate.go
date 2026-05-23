package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions configures the rotation behaviour.
type RotateOptions struct {
	// BackupDir is the directory where the old encrypted file is archived.
	// If empty, a "backups" sub-directory next to the source file is used.
	BackupDir string

	// Suffix is appended to the backup filename (e.g. a timestamp).
	// Defaults to the current UTC time formatted as 20060102T150405Z.
	Suffix string
}

// Rotate archives the current encrypted env file and replaces it with the
// newly encrypted content supplied in newEncrypted.
//
// The archived file is named <original>.<suffix> inside BackupDir.
func Rotate(encryptedPath string, newEncrypted []byte, opts RotateOptions) error {
	if opts.Suffix == "" {
		opts.Suffix = time.Now().UTC().Format("20060102T150405Z")
	}

	backupDir := opts.BackupDir
	if backupDir == "" {
		backupDir = filepath.Join(filepath.Dir(encryptedPath), "backups")
	}

	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		return fmt.Errorf("rotate: create backup dir: %w", err)
	}

	// Archive the existing file only when it already exists.
	if _, err := os.Stat(encryptedPath); err == nil {
		base := filepath.Base(encryptedPath)
		backupPath := filepath.Join(backupDir, base+"."+opts.Suffix)

		existing, err := os.ReadFile(encryptedPath)
		if err != nil {
			return fmt.Errorf("rotate: read existing file: %w", err)
		}

		if err := os.WriteFile(backupPath, existing, 0o600); err != nil {
			return fmt.Errorf("rotate: write backup: %w", err)
		}
	}

	// Write the new encrypted content.
	if err := os.WriteFile(encryptedPath, newEncrypted, 0o600); err != nil {
		return fmt.Errorf("rotate: write new file: %w", err)
	}

	return nil
}

// ListBackups returns the paths of all backup files for the given encrypted
// file, sorted lexicographically (oldest first when using timestamp suffixes).
func ListBackups(encryptedPath string, backupDir string) ([]string, error) {
	if backupDir == "" {
		backupDir = filepath.Join(filepath.Dir(encryptedPath), "backups")
	}

	base := filepath.Base(encryptedPath)
	pattern := filepath.Join(backupDir, base+".*")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	return matches, nil
}
