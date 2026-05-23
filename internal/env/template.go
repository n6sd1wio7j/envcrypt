package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TemplateEntry represents a key with an optional description and required flag.
type TemplateEntry struct {
	Key         string
	Description string
	Required    bool
}

// Template holds the expected keys for an environment file.
type Template struct {
	Entries []TemplateEntry
}

// GenerateTemplate creates a Template from an existing .env file,
// treating all keys as required and preserving inline comments as descriptions.
func GenerateTemplate(path string) (*Template, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("template: open %q: %w", path, err)
	}
	defer f.Close()

	var tmpl Template
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, _, desc := parseLineWithComment(line)
		if key == "" {
			continue
		}
		tmpl.Entries = append(tmpl.Entries, TemplateEntry{
			Key:         key,
			Description: desc,
			Required:    true,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("template: scan %q: %w", path, err)
	}
	return &tmpl, nil
}

// CheckTemplate verifies that all required template keys are present in entries.
// Returns a list of missing required keys.
func CheckTemplate(tmpl *Template, entries []Entry) []string {
	present := make(map[string]bool, len(entries))
	for _, e := range entries {
		present[e.Key] = true
	}
	var missing []string
	for _, te := range tmpl.Entries {
		if te.Required && !present[te.Key] {
			missing = append(missing, te.Key)
		}
	}
	return missing
}

// parseLineWithComment parses a key=value line and extracts an inline comment.
func parseLineWithComment(line string) (key, value, comment string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", ""
	}
	key = strings.TrimSpace(parts[0])
	rest := parts[1]
	if idx := strings.Index(rest, " #"); idx >= 0 {
		value = strings.TrimSpace(rest[:idx])
		comment = strings.TrimSpace(rest[idx+2:])
	} else {
		value = strings.TrimSpace(rest)
	}
	return key, value, comment
}
