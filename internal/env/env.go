package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
	Raw   string
}

// File represents a parsed .env file.
type File struct {
	Entries []Entry
}

// Parse reads and parses a .env file from the given path.
func Parse(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		entry, ok := parseLine(line)
		if ok {
			entries = append(entries, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}
	return &File{Entries: entries}, nil
}

// Write serialises the File back to disk at the given path.
func Write(path string, ef *File) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating env file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range ef.Entries {
		if _, err := fmt.Fprintln(w, e.Raw); err != nil {
			return fmt.Errorf("writing env entry: %w", err)
		}
	}
	return w.Flush()
}

// parseLine parses a single line into an Entry.
// Returns (entry, true) for valid KEY=VALUE lines, (zero, false) otherwise.
func parseLine(line string) (Entry, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return Entry{}, false
	}
	parts := strings.SplitN(trimmed, "=", 2)
	if len(parts) != 2 {
		return Entry{}, false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" {
		return Entry{}, false
	}
	return Entry{Key: key, Value: value, Raw: trimmed}, true
}
