package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry represents a single recorded change to an env file.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"` // e.g. "encrypt", "decrypt", "rotate", "merge"
	File      string    `json:"file"`
	User      string    `json:"user,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// DefaultHistoryPath is the default location for the history log.
const DefaultHistoryPath = ".envcrypt/history.json"

// AppendHistory appends a new entry to the history log at the given path.
func AppendHistory(path string, entry HistoryEntry) error {
	entries, err := LoadHistory(path)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	entries = append(entries, entry)
	return SaveHistory(path, entries)
}

// LoadHistory reads the history log from path. Returns an empty slice if the
// file does not exist.
func LoadHistory(path string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []HistoryEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history file: %w", err)
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse history file: %w", err)
	}
	return entries, nil
}

// SaveHistory writes entries to path, creating parent directories as needed.
func SaveHistory(path string, entries []HistoryEntry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create history dir: %w", err)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write history file: %w", err)
	}
	return nil
}
