package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of an env file.
type Snapshot struct {
	Timestamp time.Time
	Path      string
	Entries   []Entry
}

// TakeSnapshot reads the given env file and returns a Snapshot.
func TakeSnapshot(envPath string) (*Snapshot, error) {
	entries, err := Parse(envPath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: parse %q: %w", envPath, err)
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Path:      envPath,
		Entries:   entries,
	}, nil
}

// SaveSnapshot writes the snapshot entries to a timestamped file inside dir.
// It returns the path of the written snapshot file.
func SaveSnapshot(snap *Snapshot, dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("snapshot: mkdir %q: %w", dir, err)
	}
	name := fmt.Sprintf("snapshot_%s.env", snap.Timestamp.Format("20060102T150405Z"))
	dest := filepath.Join(dir, name)
	if err := Write(dest, snap.Entries); err != nil {
		return "", fmt.Errorf("snapshot: write %q: %w", dest, err)
	}
	return dest, nil
}

// LoadSnapshot reads a previously saved snapshot file and returns a Snapshot.
func LoadSnapshot(path string) (*Snapshot, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: load %q: %w", path, err)
	}
	ts, err := timestampFromPath(filepath.Base(path))
	if err != nil {
		ts = time.Time{}
	}
	return &Snapshot{
		Timestamp: ts,
		Path:      path,
		Entries:   entries,
	}, nil
}

// timestampFromPath parses a timestamp from a snapshot filename of the form
// snapshot_20060102T150405Z.env.
func timestampFromPath(name string) (time.Time, error) {
	const prefix = "snapshot_"
	const suffix = ".env"
	if len(name) <= len(prefix)+len(suffix) {
		return time.Time{}, fmt.Errorf("unexpected filename %q", name)
	}
	raw := name[len(prefix) : len(name)-len(suffix)]
	return time.Parse("20060102T150405Z", raw)
}
