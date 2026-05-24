package env

import (
	"fmt"
	"sort"
)

// ChangeKind describes the type of change between two env files.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeRemoved ChangeKind = "removed"
	ChangeUpdated ChangeKind = "updated"
)

// Change represents a single difference between two env files.
type Change struct {
	Kind     ChangeKind
	Key      string
	OldValue string
	NewValue string
}

// String returns a human-readable representation of the change.
func (c Change) String() string {
	switch c.Kind {
	case ChangeAdded:
		return fmt.Sprintf("+ %s=%s", c.Key, c.NewValue)
	case ChangeRemoved:
		return fmt.Sprintf("- %s=%s", c.Key, c.OldValue)
	case ChangeUpdated:
		return fmt.Sprintf("~ %s: %s -> %s", c.Key, c.OldValue, c.NewValue)
	default:
		return fmt.Sprintf("? %s", c.Key)
	}
}

// Diff computes the differences between an old and new env File.
// The returned changes are sorted by key for deterministic output.
func Diff(old, new *File) []Change {
	oldMap := toMap(old)
	newMap := toMap(new)

	var changes []Change

	for key, oldVal := range oldMap {
		newVal, exists := newMap[key]
		if !exists {
			changes = append(changes, Change{Kind: ChangeRemoved, Key: key, OldValue: oldVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Kind: ChangeUpdated, Key: key, OldValue: oldVal, NewValue: newVal})
		}
	}

	for key, newVal := range newMap {
		if _, exists := oldMap[key]; !exists {
			changes = append(changes, Change{Kind: ChangeAdded, Key: key, NewValue: newVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return changes
}

func toMap(ef *File) map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		m[e.Key] = e.Value
	}
	return m
}
