package env

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaEntry defines the expected shape of a single env variable.
type SchemaEntry struct {
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
}

// Schema maps variable names to their schema definitions.
type Schema map[string]SchemaEntry

// SchemaViolation describes a single schema violation.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// LoadSchema reads a YAML schema file from the given path.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema file not found: %s", path)
		}
		return nil, fmt.Errorf("read schema: %w", err)
	}
	var schema Schema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parse schema: %w", err)
	}
	return schema, nil
}

// ValidateAgainstSchema checks that entries conform to the schema.
// It returns a slice of violations; an empty slice means valid.
func ValidateAgainstSchema(entries []Entry, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	present := make(map[string]string, len(entries))
	for _, e := range entries {
		present[e.Key] = e.Value
	}

	for key, def := range schema {
		val, ok := present[key]
		if !ok || strings.TrimSpace(val) == "" {
			if def.Required && strings.TrimSpace(def.Default) == "" {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: "required variable is missing or empty",
				})
			}
		}
	}

	return violations
}
