package env

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// DefaultLockPath is the default location for the lock file.
const DefaultLockPath = ".envcrypt.lock"

// LockEntry records metadata about an encrypted env file.
type LockEntry struct {
	File      string    `json:"file"`
	Encrypted string    `json:"encrypted"`
	Checksum  string    `json:"checksum"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

// LockFile holds all lock entries indexed by env file path.
type LockFile struct {
	Version int                  `json:"version"`
	Entries map[string]LockEntry `json:"entries"`
}

// LoadLockFile reads and parses the lock file at path.
// Returns an empty LockFile if the file does not exist.
func LoadLockFile(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &LockFile{Version: 1, Entries: make(map[string]LockEntry)}, nil
		}
		return nil, err
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, err
	}
	if lf.Entries == nil {
		lf.Entries = make(map[string]LockEntry)
	}
	return &lf, nil
}

// SaveLockFile writes the lock file to path, creating directories as needed.
func SaveLockFile(path string, lf *LockFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Upsert adds or updates the lock entry for the given env file.
func (lf *LockFile) Upsert(entry LockEntry) {
	lf.Entries[entry.File] = entry
}

// Remove deletes the lock entry for the given env file key.
func (lf *LockFile) Remove(file string) {
	delete(lf.Entries, file)
}
