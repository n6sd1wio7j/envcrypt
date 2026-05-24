package env

import (
	"regexp"
	"strings"
)

// RedactOptions controls how sensitive values are masked.
type RedactOptions struct {
	// Keys is an explicit list of keys whose values should be redacted.
	Keys []string
	// Patterns is a list of regex patterns matched against key names.
	Patterns []string
	// Mask is the string used to replace sensitive values. Defaults to "***".
	Mask string
}

// defaultSensitivePatterns are common patterns for sensitive env var names.
var defaultSensitivePatterns = []string{
	`(?i)password`,
	`(?i)secret`,
	`(?i)token`,
	`(?i)api_?key`,
	`(?i)private_?key`,
	`(?i)auth`,
	`(?i)credential`,
}

// Redact returns a copy of entries with sensitive values replaced by a mask.
// It applies both explicit key matches and regex pattern matches against key names.
// If opts.Mask is empty, "***" is used.
func Redact(entries []Entry, opts RedactOptions) []Entry {
	mask := opts.Mask
	if mask == "" {
		mask = "***"
	}

	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = struct{}{}
	}

	patterns := make([]*regexp.Regexp, 0, len(opts.Patterns)+len(defaultSensitivePatterns))
	for _, p := range append(defaultSensitivePatterns, opts.Patterns...) {
		if re, err := regexp.Compile(p); err == nil {
			patterns = append(patterns, re)
		}
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		result[i] = e
		if isSensitive(e.Key, explicit, patterns) {
			result[i].Value = mask
		}
	}
	return result
}

// RedactMap returns a copy of the map with sensitive values replaced.
func RedactMap(m map[string]string, opts RedactOptions) map[string]string {
	mask := opts.Mask
	if mask == "" {
		mask = "***"
	}

	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = struct{}{}
	}

	patterns := make([]*regexp.Regexp, 0, len(opts.Patterns)+len(defaultSensitivePatterns))
	for _, p := range append(defaultSensitivePatterns, opts.Patterns...) {
		if re, err := regexp.Compile(p); err == nil {
			patterns = append(patterns, re)
		}
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if isSensitive(k, explicit, patterns) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

func isSensitive(key string, explicit map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := explicit[key]; ok {
		return true
	}
	upper := strings.ToUpper(key)
	_ = upper
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
