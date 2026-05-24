// Package env provides utilities for managing .env file entries.
//
// # Promote
//
// Promote copies entries from a source environment into a destination
// environment, with fine-grained control over which keys are transferred
// and how conflicts are resolved.
//
// Basic usage — promote all keys, skip existing:
//
//	result, out, err := env.Promote(stagingEntries, prodEntries, env.PromoteOptions{})
//
// Promote a subset of keys and overwrite existing values:
//
//	opts := env.PromoteOptions{
//		Keys:      []string{"API_URL", "TIMEOUT"},
//		Overwrite: true,
//	}
//	out, result, err := env.Promote(src, dst, opts)
//
// Dry-run mode returns what would change without modifying the destination:
//
//	opts := env.PromoteOptions{DryRun: true}
//	_, result, err := env.Promote(src, dst, opts)
//	fmt.Println("Would promote:", result.Promoted)
package env
