package env

import (
	"fmt"
	"strings"
)

// PromoteOptions controls how values are promoted between environments.
type PromoteOptions struct {
	// Keys restricts promotion to a specific set of keys. If empty, all keys are promoted.
	Keys []string
	// Overwrite controls whether existing keys in the destination are overwritten.
	Overwrite bool
	// DryRun returns what would change without modifying the destination.
	DryRun bool
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Overwritten []string
}

// Promote copies entries from src to dst according to opts.
// It returns a PromoteResult summarising what changed.
func Promote(src, dst []Entry, opts PromoteOptions) ([]Entry, PromoteResult, error) {
	if src == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: source entries must not be nil")
	}

	allowedKeys := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowedKeys[strings.TrimSpace(k)] = true
	}

	dstMap := make(map[string]int, len(dst))
	for i, e := range dst {
		dstMap[e.Key] = i
	}

	result := PromoteResult{}
	out := make([]Entry, len(dst))
	copy(out, dst)

	for _, e := range src {
		if len(allowedKeys) > 0 && !allowedKeys[e.Key] {
			continue
		}

		if idx, exists := dstMap[e.Key]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, e.Key)
				continue
			}
			if !opts.DryRun {
				out[idx] = e
			}
			result.Overwritten = append(result.Overwritten, e.Key)
		} else {
			if !opts.DryRun {
				out = append(out, e)
				dstMap[e.Key] = len(out) - 1
			}
			result.Promoted = append(result.Promoted, e.Key)
		}
	}

	return out, result, nil
}
