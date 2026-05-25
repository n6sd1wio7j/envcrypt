package env

import (
	"regexp"
	"strings"
)

// SearchOptions controls how Search filters entries.
type SearchOptions struct {
	// KeyPattern is an optional regex applied to entry keys.
	KeyPattern string
	// ValuePattern is an optional regex applied to entry values.
	ValuePattern string
	// CaseSensitive controls whether patterns are case-sensitive.
	CaseSensitive bool
}

// SearchResult holds a matched entry and the file it came from.
type SearchResult struct {
	Entry Entry
	File  string
}

// Search filters entries from a .env file according to the given options.
// It returns all entries whose key and/or value match the provided patterns.
// If both patterns are empty, all entries are returned.
func Search(path string, opts SearchOptions) ([]SearchResult, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, err
	}

	flags := ""
	if !opts.CaseSensitive {
		flags = "(?i)"
	}

	var keyRe, valRe *regexp.Regexp
	if opts.KeyPattern != "" {
		keyRe, err = regexp.Compile(flags + opts.KeyPattern)
		if err != nil {
			return nil, err
		}
	}
	if opts.ValuePattern != "" {
		valRe, err = regexp.Compile(flags + opts.ValuePattern)
		if err != nil {
			return nil, err
		}
	}

	var results []SearchResult
	for _, e := range entries {
		if keyRe != nil && !keyRe.MatchString(e.Key) {
			continue
		}
		if valRe != nil && !valRe.MatchString(e.Value) {
			continue
		}
		results = append(results, SearchResult{Entry: e, File: path})
	}
	return results, nil
}

// SearchMultiple runs Search across multiple files and aggregates results.
func SearchMultiple(paths []string, opts SearchOptions) ([]SearchResult, error) {
	var all []SearchResult
	for _, p := range paths {
		res, err := Search(p, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, res...)
	}
	return all, nil
}

// FormatSearchResults returns a human-readable string for a slice of results.
func FormatSearchResults(results []SearchResult) string {
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(r.File)
		sb.WriteString(": ")
		sb.WriteString(r.Entry.Key)
		sb.WriteString("=")
		sb.WriteString(r.Entry.Value)
		sb.WriteString("\n")
	}
	return sb.String()
}
