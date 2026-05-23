package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	other := []Entry{{Key: "NEW", Value: "val"}}

	res, err := Merge(base, other, false)
	require.NoError(t, err)
	assert.Empty(t, res.Conflicts)
	assert.Equal(t, []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
		{Key: "NEW", Value: "val"},
	}, res.Entries)
}

func TestMerge_OverrideWins(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "old"}}
	other := []Entry{{Key: "FOO", Value: "new"}}

	res, err := Merge(base, other, false)
	require.NoError(t, err)
	require.Len(t, res.Conflicts, 1)
	assert.Equal(t, "FOO", res.Conflicts[0].Key)
	assert.Equal(t, "old", res.Conflicts[0].BaseVal)
	assert.Equal(t, "new", res.Conflicts[0].OtherVal)
	// override value wins
	assert.Equal(t, "new", res.Entries[0].Value)
}

func TestMerge_StrictConflictErrors(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "old"}}
	other := []Entry{{Key: "FOO", Value: "new"}}

	_, err := Merge(base, other, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "FOO")
}

func TestMerge_SameValueNoConflict(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "same"}}
	other := []Entry{{Key: "FOO", Value: "same"}}

	res, err := Merge(base, other, true)
	require.NoError(t, err)
	assert.Empty(t, res.Conflicts)
	assert.Equal(t, "same", res.Entries[0].Value)
}

func TestMerge_EmptyBase(t *testing.T) {
	other := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	res, err := Merge(nil, other, false)
	require.NoError(t, err)
	assert.Empty(t, res.Conflicts)
	assert.Equal(t, other, res.Entries)
}

func TestMerge_EmptyOther(t *testing.T) {
	base := []Entry{{Key: "X", Value: "y"}}
	res, err := Merge(base, nil, false)
	require.NoError(t, err)
	assert.Empty(t, res.Conflicts)
	assert.Equal(t, base, res.Entries)
}
