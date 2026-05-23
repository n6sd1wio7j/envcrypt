package env

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) OK() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

var validKeyRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// ValidateEntries checks a slice of Entry values for common issues:
// - empty keys
// - keys that don't match the conventional KEY_NAME format
// - duplicate keys
func ValidateEntries(entries []Entry) *ValidationResult {
	result := &ValidationResult{}
	seen := make(map[string]int)

	for _, e := range entries {
		if e.Key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     "(empty)",
				Message: "key must not be empty",
			})
			continue
		}

		if !validKeyRe.MatchString(e.Key) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     e.Key,
				Message: "key must match pattern [A-Z_][A-Z0-9_]*",
			})
		}

		seen[e.Key]++
		if seen[e.Key] == 2 {
			result.Errors = append(result.Errors, ValidationError{
				Key:     e.Key,
				Message: "duplicate key",
			})
		}
	}

	return result
}

// ValidateFile parses the given file and validates its entries.
func ValidateFile(path string) (*ValidationResult, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return ValidateEntries(entries), nil
}
