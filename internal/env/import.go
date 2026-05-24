package env

import (
	"fmt"
	"os"
	"strings"
)

// ImportFormat represents the source format for importing env entries.
type ImportFormat string

const (
	ImportDotenv ImportFormat = "dotenv"
	ImportShell  ImportFormat = "shell"
)

// ImportOptions controls how entries are imported.
type ImportOptions struct {
	// Format is the source format (dotenv or shell).
	Format ImportFormat
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// Keys restricts import to specific keys; empty means all.
	Keys []string
}

// Import reads entries from src according to opts and merges them into dst,
// returning the resulting entries. If dst does not exist it is created.
func Import(src, dst string, opts ImportOptions) ([]Entry, error) {
	raw, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("import: read source %q: %w", src, err)
	}

	incoming, err := parseImport(string(raw), opts.Format)
	if err != nil {
		return nil, fmt.Errorf("import: parse source: %w", err)
	}

	if len(opts.Keys) > 0 {
		incoming = filterEntries(incoming, opts.Keys)
	}

	var base []Entry
	if _, statErr := os.Stat(dst); statErr == nil {
		base, err = Parse(dst)
		if err != nil {
			return nil, fmt.Errorf("import: parse destination %q: %w", dst, err)
		}
	}

	mode := MergeSkip
	if opts.Overwrite {
		mode = MergeOverride
	}

	merged, err := Merge(base, incoming, mode)
	if err != nil {
		return nil, fmt.Errorf("import: merge: %w", err)
	}

	if err := Write(dst, merged); err != nil {
		return nil, fmt.Errorf("import: write destination %q: %w", dst, err)
	}

	return merged, nil
}

// parseImport converts raw text in the given format to entries.
func parseImport(raw string, format ImportFormat) ([]Entry, error) {
	switch format {
	case ImportShell:
		return parseShellExport(raw)
	default:
		// dotenv and unknown formats fall through to standard parser
		return parseLines(raw)
	}
}

// parseShellExport handles lines like: export KEY=VALUE
func parseShellExport(raw string) ([]Entry, error) {
	var entries []Entry
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		e, ok := parseLine(line)
		if ok {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

// parseLines is a helper that parses dotenv-style text without requiring a file.
func parseLines(raw string) ([]Entry, error) {
	var entries []Entry
	for _, line := range strings.Split(raw, "\n") {
		if e, ok := parseLine(strings.TrimSpace(line)); ok {
			entries = append(entries, e)
		}
	}
	return entries, nil
}
