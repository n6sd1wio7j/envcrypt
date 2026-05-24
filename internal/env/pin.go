package env

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// PinEntry records a pinned snapshot reference for a named environment.
type PinEntry struct {
	Name      string    `json:"name"`
	SnapshotPath string `json:"snapshot_path"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by,omitempty"`
}

// PinFile holds all pinned entries.
type PinFile struct {
	Pins []PinEntry `json:"pins"`
}

// DefaultPinPath is the default location for the pin file.
const DefaultPinPath = ".envcrypt/pins.json"

// LoadPinFile loads the pin file from path. Returns an empty PinFile if the
// file does not exist.
func LoadPinFile(path string) (*PinFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &PinFile{}, nil
		}
		return nil, err
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	return &pf, nil
}

// SavePinFile writes the PinFile to path, creating parent directories as needed.
func SavePinFile(path string, pf *PinFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// AddPin adds or replaces a pin entry by name.
func AddPin(pf *PinFile, name, snapshotPath, pinnedBy string) {
	entry := PinEntry{
		Name:         name,
		SnapshotPath: snapshotPath,
		PinnedAt:     time.Now().UTC(),
		PinnedBy:     pinnedBy,
	}
	for i, p := range pf.Pins {
		if p.Name == name {
			pf.Pins[i] = entry
			return
		}
	}
	pf.Pins = append(pf.Pins, entry)
}

// RemovePin removes a pin entry by name. Returns false if not found.
func RemovePin(pf *PinFile, name string) bool {
	for i, p := range pf.Pins {
		if p.Name == name {
			pf.Pins = append(pf.Pins[:i], pf.Pins[i+1:]...)
			return true
		}
	}
	return false
}

// FindPin looks up a pin entry by name. Returns nil if not found.
func FindPin(pf *PinFile, name string) *PinEntry {
	for i, p := range pf.Pins {
		if p.Name == name {
			return &pf.Pins[i]
		}
	}
	return nil
}
