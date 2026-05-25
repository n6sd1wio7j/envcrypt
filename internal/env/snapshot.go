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

// Diff returns the keys that were added, removed, or changed between s and
// other. It returns three slices: added keys (present in other but not s),
// removed keys (present in s but not other), and changed keys (present in
// both but with different values).
func (s *Snapshot) Diff(other *Snapshot) (added, removed, changed []string) {
	base := make(map[string]string, len(s.Entries))
	for _, e := range s.Entries {
		base[e.Key] = e.Value
	}
	next := make(map[string]string, len(other.Entries))
	for _, e := range other.Entries {
		next[e.Key] = e.Value
	}
	for k, v := range next {
		if old, ok := base[k]; !ok {
			added = append(added, k)
		} else if old != v {
			changed = append(changed, k)
		}
	}
	for k := range base {
		if _, ok := next[k]; !ok {
			removed = append(removed, k)
		}
	}
	return added, removed, changed
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
