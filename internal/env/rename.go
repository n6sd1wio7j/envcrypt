package env

import "fmt"

// RenameOptions controls the behaviour of Rename.
type RenameOptions struct {
	// DryRun reports what would change without writing anything.
	DryRun bool
}

// RenameResult holds the outcome of a single rename operation.
type RenameResult struct {
	OldKey string
	NewKey string
	Found  bool
}

// Rename renames one or more keys inside an env file.
// Each element of pairs must be a [2]string{oldKey, newKey}.
// Returns an error if newKey already exists in the file (to prevent
// silent overwrites) or if any oldKey is not found.
func Rename(path string, pairs [][2]string, opts RenameOptions) ([]RenameResult, error) {
	entries, err := Parse(path)
	if err != nil {
		return nil, fmt.Errorf("rename: parse %q: %w", path, err)
	}

	// Build an index for O(1) lookups.
	idx := make(map[string]int, len(entries))
	for i, e := range entries {
		idx[e.Key] = i
	}

	results := make([]RenameResult, 0, len(pairs))

	for _, p := range pairs {
		old, nw := p[0], p[1]
		res := RenameResult{OldKey: old, NewKey: nw}

		if _, exists := idx[old]; !exists {
			return nil, fmt.Errorf("rename: key %q not found in %q", old, path)
		}
		res.Found = true

		if old != nw {
			if _, conflict := idx[nw]; conflict {
				return nil, fmt.Errorf("rename: key %q already exists in %q", nw, path)
			}
		}

		results = append(results, res)
	}

	if opts.DryRun {
		return results, nil
	}

	// Apply renames.
	for _, p := range pairs {
		old, nw := p[0], p[1]
		i := idx[old]
		entries[i].Key = nw
		// Keep idx consistent for subsequent iterations.
		delete(idx, old)
		idx[nw] = i
	}

	if err := Write(path, entries); err != nil {
		return nil, fmt.Errorf("rename: write %q: %w", path, err)
	}

	return results, nil
}
