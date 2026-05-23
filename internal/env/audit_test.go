package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendAndLoadAuditEvent(t *testing.T) {
	dir := t.TempDir()
	path := DefaultAuditPath(dir)

	event := AuditEvent{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Action:    "encrypt",
		File:      ".env",
		User:      "alice",
		Details:   "initial encryption",
	}

	if err := AppendAuditEvent(path, event); err != nil {
		t.Fatalf("AppendAuditEvent: %v", err)
	}

	log, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(log.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(log.Events))
	}
	if log.Events[0].Action != "encrypt" {
		t.Errorf("expected action 'encrypt', got %q", log.Events[0].Action)
	}
	if log.Events[0].User != "alice" {
		t.Errorf("expected user 'alice', got %q", log.Events[0].User)
	}
}

func TestAppendAuditEvent_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := DefaultAuditPath(dir)

	actions := []string{"encrypt", "decrypt", "rotate"}
	for _, a := range actions {
		e := AuditEvent{Timestamp: time.Now().UTC(), Action: a, File: ".env"}
		if err := AppendAuditEvent(path, e); err != nil {
			t.Fatalf("AppendAuditEvent(%s): %v", a, err)
		}
	}

	log, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(log.Events) != len(actions) {
		t.Fatalf("expected %d events, got %d", len(actions), len(log.Events))
	}
	for i, a := range actions {
		if log.Events[i].Action != a {
			t.Errorf("event[%d]: expected %q, got %q", i, a, log.Events[i].Action)
		}
	}
}

func TestLoadAuditLog_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	log, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(log.Events) != 0 {
		t.Errorf("expected empty log, got %d events", len(log.Events))
	}
}

func TestSaveAuditLog_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	path := filepath.Join(dir, "audit.json")
	log := &AuditLog{Events: []AuditEvent{{Action: "test", File: ".env", Timestamp: time.Now().UTC()}}}
	if err := SaveAuditLog(path, log); err != nil {
		t.Fatalf("SaveAuditLog: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
