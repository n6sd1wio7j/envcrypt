package env

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two env files.
type CompareResult struct {
	OnlyInA  []Entry
	OnlyInB  []Entry
	Changed  []ChangedEntry
	Identical []Entry
}

// ChangedEntry represents a key whose value differs between two env files.
type ChangedEntry struct {
	Key    string
	ValueA string
	ValueB string
}

// Compare compares two slices of Entry and returns a CompareResult describing
// keys that are unique to each side, changed between sides, or identical.
func Compare(a, b []Entry) CompareResult {
	aMap := toMap(a)
	bMap := toMap(b)

	result := CompareResult{}

	for _, e := range a {
		if bVal, ok := bMap[e.Key]; ok {
			if e.Value == bVal {
				result.Identical = append(result.Identical, e)
			} else {
				result.Changed = append(result.Changed, ChangedEntry{
					Key:    e.Key,
					ValueA: e.Value,
					ValueB: bVal,
				})
			}
		} else {
			result.OnlyInA = append(result.OnlyInA, e)
		}
	}

	for _, e := range b {
		if _, ok := aMap[e.Key]; !ok {
			result.OnlyInB = append(result.OnlyInB, e)
		}
	}

	return result
}

// FormatCompareResult returns a human-readable summary of a CompareResult.
func FormatCompareResult(r CompareResult, labelA, labelB string) string {
	var sb strings.Builder

	keys := func(entries []Entry) []string {
		out := make([]string, len(entries))
		for i, e := range entries {
			out[i] = e.Key
		}
		sort.Strings(out)
		return out
	}

	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(&sb, "Only in %s: %s\n", labelA, strings.Join(keys(r.OnlyInA), ", "))
	}
	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(&sb, "Only in %s: %s\n", labelB, strings.Join(keys(r.OnlyInB), ", "))
	}
	for _, c := range r.Changed {
		fmt.Fprintf(&sb, "Changed: %s (%s=%q | %s=%q)\n", c.Key, labelA, c.ValueA, labelB, c.ValueB)
	}
	if len(r.Identical) > 0 {
		fmt.Fprintf(&sb, "Identical: %d key(s)\n", len(r.Identical))
	}
	if sb.Len() == 0 {
		return "No differences found.\n"
	}
	return sb.String()
}
