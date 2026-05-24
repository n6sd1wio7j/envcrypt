// Package env provides utilities for working with .env files, including
// parsing, writing, diffing, merging, snapshotting, auditing, and comparing.
//
// # Compare
//
// The Compare function performs a detailed key-level comparison between two
// slices of Entry values, categorising each key as:
//
//   - OnlyInA: present in the first set but not the second
//   - OnlyInB: present in the second set but not the first
//   - Changed: present in both but with differing values
//   - Identical: present in both with the same value
//
// Example usage:
//
//	a, _ := env.Parse(".env.staging")
//	b, _ := env.Parse(".env.production")
//	result := env.Compare(a, b)
//	fmt.Print(env.FormatCompareResult(result, "staging", "production"))
//
// FormatCompareResult produces a human-readable summary suitable for CLI
// output, showing which keys differ and how.
package env
