package env

import "fmt"

// MergeResult holds the outcome of merging two env files.
type MergeResult struct {
	Entries  []Entry
	Conflicts []Conflict
}

// Conflict represents a key present in both base and override with different values.
type Conflict struct {
	Key      string
	BaseVal  string
	OtherVal string
}

// Merge combines base and other env entries.
// Keys in other override keys in base. If a key exists in both with different
// values, it is recorded as a conflict but the other (override) value wins.
// Pass strict=true to return an error instead of silently overriding conflicts.
func Merge(base, other []Entry, strict bool) (*MergeResult, error) {
	result := &MergeResult{}

	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		baseMap[e.Key] = e.Value
	}

	otherMap := make(map[string]string, len(other))
	for _, e := range other {
		otherMap[e.Key] = e.Value
	}

	// Start with base entries, applying overrides.
	seen := make(map[string]bool)
	for _, e := range base {
		if oval, ok := otherMap[e.Key]; ok {
			if oval != e.Value {
				result.Conflicts = append(result.Conflicts, Conflict{
					Key:      e.Key,
					BaseVal:  e.Value,
					OtherVal: oval,
				})
				if strict {
					return nil, fmt.Errorf("merge conflict on key %q: %q vs %q", e.Key, e.Value, oval)
				}
			}
			result.Entries = append(result.Entries, Entry{Key: e.Key, Value: oval})
		} else {
			result.Entries = append(result.Entries, e)
		}
		seen[e.Key] = true
	}

	// Append keys only present in other.
	for _, e := range other {
		if !seen[e.Key] {
			result.Entries = append(result.Entries, e)
		}
	}

	return result, nil
}
