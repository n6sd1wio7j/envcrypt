package env

import (
	"testing"
)

func makeFile(entries ...Entry) *File {
	return &File{Entries: entries}
}

func entry(k, v string) Entry {
	return Entry{Key: k, Value: v, Raw: k + "=" + v}
}

func TestDiff_NoChanges(t *testing.T) {
	old := makeFile(entry("FOO", "bar"))
	new := makeFile(entry("FOO", "bar"))
	changes := Diff(old, new)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestDiff_Added(t *testing.T) {
	old := makeFile()
	new := makeFile(entry("NEW_KEY", "value"))
	changes := Diff(old, new)
	if len(changes) != 1 || changes[0].Kind != ChangeAdded {
		t.Errorf("expected 1 added change, got %+v", changes)
	}
	if changes[0].Key != "NEW_KEY" || changes[0].NewValue != "value" {
		t.Errorf("unexpected change content: %+v", changes[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	old := makeFile(entry("OLD_KEY", "val"))
	new := makeFile()
	changes := Diff(old, new)
	if len(changes) != 1 || changes[0].Kind != ChangeRemoved {
		t.Errorf("expected 1 removed change, got %+v", changes)
	}
}

func TestDiff_Updated(t *testing.T) {
	old := makeFile(entry("KEY", "old"))
	new := makeFile(entry("KEY", "new"))
	changes := Diff(old, new)
	if len(changes) != 1 || changes[0].Kind != ChangeUpdated {
		t.Errorf("expected 1 updated change, got %+v", changes)
	}
	if changes[0].OldValue != "old" || changes[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", changes[0])
	}
}

func TestChangeString(t *testing.T) {
	tests := []struct {
		c    Change
		want string
	}{
		{Change{Kind: ChangeAdded, Key: "K", NewValue: "v"}, "+ K=v"},
		{Change{Kind: ChangeRemoved, Key: "K", OldValue: "v"}, "- K=v"},
		{Change{Kind: ChangeUpdated, Key: "K", OldValue: "a", NewValue: "b"}, "~ K: a -> b"},
	}
	for _, tt := range tests {
		if got := tt.c.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}
