package env

import (
	"fmt"
	"os"
	"strings"
)

// ExportFormat defines the output format for exported env entries.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatShell  ExportFormat = "shell"
	FormatJSON   ExportFormat = "json"
)

// ExportOptions controls how entries are exported.
type ExportOptions struct {
	Format ExportFormat
	Keys   []string // if non-empty, only export these keys
}

// Export writes env entries to w in the specified format.
// If opts.Keys is non-empty, only those keys are included.
func Export(entries []Entry, opts ExportOptions) (string, error) {
	filtered := filterEntries(entries, opts.Keys)

	switch opts.Format {
	case FormatShell:
		return exportShell(filtered), nil
	case FormatJSON:
		return exportJSON(filtered), nil
	case FormatDotenv, "":
		return exportDotenv(filtered), nil
	default:
		return "", fmt.Errorf("unknown export format: %q", opts.Format)
	}
}

// ExportToFile writes exported content to a file.
func ExportToFile(entries []Entry, opts ExportOptions, dest string) error {
	content, err := Export(entries, opts)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, []byte(content), 0600)
}

func filterEntries(entries []Entry, keys []string) []Entry {
	if len(keys) == 0 {
		return entries
	}
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := make([]Entry, 0, len(keys))
	for _, e := range entries {
		if _, ok := set[e.Key]; ok {
			out = append(out, e)
		}
	}
	return out
}

func exportDotenv(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

func exportShell(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.Value)
	}
	return sb.String()
}

func exportJSON(entries []Entry) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", e.Key, e.Value, comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}
