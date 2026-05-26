package env

import (
	"fmt"
	"strings"
)

// LintSeverity indicates the severity level of a lint finding.
type LintSeverity string

const (
	LintWarn  LintSeverity = "warn"
	LintError LintSeverity = "error"
)

// LintFinding represents a single lint issue found in an env file.
type LintFinding struct {
	Line     int
	Key      string
	Message  string
	Severity LintSeverity
}

func (f LintFinding) String() string {
	if f.Key != "" {
		return fmt.Sprintf("%s (line %d) [%s]: %s", f.Key, f.Line, f.Severity, f.Message)
	}
	return fmt.Sprintf("line %d [%s]: %s", f.Line, f.Severity, f.Message)
}

// LintOptions controls which lint checks are enabled.
type LintOptions struct {
	WarnOnNoValue     bool // warn when a key has no value and no trailing =
	WarnOnLowercase   bool // warn when a key contains lowercase letters
	ErrorOnDuplicate  bool // error when a key appears more than once
}

// DefaultLintOptions returns a sensible default lint configuration.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		WarnOnNoValue:    true,
		WarnOnLowercase:  true,
		ErrorOnDuplicate: true,
	}
}

// Lint reads an env file and returns a list of findings based on opts.
func Lint(path string, opts LintOptions) ([]LintFinding, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("lint: %w", err)
	}

	var findings []LintFinding
	seen := make(map[string]int)

	for i, e := range entries {
		lineNum := i + 1

		if opts.ErrorOnDuplicate {
			if prev, ok := seen[e.Key]; ok {
				findings = append(findings, LintFinding{
					Line:     lineNum,
					Key:      e.Key,
					Message:  fmt.Sprintf("duplicate key, first seen at entry %d", prev+1),
					Severity: LintError,
				})
			}
			seen[e.Key] = i
		}

		if opts.WarnOnNoValue && e.Value == "" {
			findings = append(findings, LintFinding{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "key has an empty value",
				Severity: LintWarn,
			})
		}

		if opts.WarnOnLowercase && strings.ToUpper(e.Key) != e.Key {
			findings = append(findings, LintFinding{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "key contains lowercase letters; convention is UPPER_SNAKE_CASE",
				Severity: LintWarn,
			})
		}
	}

	return findings, nil
}
