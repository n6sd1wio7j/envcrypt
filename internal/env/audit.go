package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEvent represents a single recorded action on an env file.
type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	File      string    `json:"file"`
	User      string    `json:"user,omitempty"`
	Details   string    `json:"details,omitempty"`
}

// AuditLog holds an ordered list of audit events.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

// DefaultAuditPath returns the default path for the audit log file.
func DefaultAuditPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_audit.json")
}

// AppendAuditEvent loads an existing audit log (or creates a new one),
// appends the given event, and writes it back to path.
func AppendAuditEvent(path string, event AuditEvent) error {
	log, err := LoadAuditLog(path)
	if err != nil {
		return fmt.Errorf("audit: load log: %w", err)
	}
	log.Events = append(log.Events, event)
	return SaveAuditLog(path, log)
}

// LoadAuditLog reads and parses an audit log from path.
// If the file does not exist an empty log is returned.
func LoadAuditLog(path string) (*AuditLog, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, err
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("audit: parse log: %w", err)
	}
	return &log, nil
}

// SaveAuditLog serialises the log to path, creating parent directories as needed.
func SaveAuditLog(path string, log *AuditLog) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
