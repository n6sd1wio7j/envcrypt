// Package env provides utilities for parsing, writing, diffing, merging,
// rotating, snapshotting, auditing, locking, and validating .env files.
//
// # Validation
//
// The validate module checks env entry collections for common mistakes before
// they are encrypted or committed:
//
//   - Empty keys are rejected.
//   - Keys must follow the conventional ALL_CAPS_WITH_UNDERSCORES format
//     (regex: [A-Z_][A-Z0-9_]*).
//   - Duplicate keys are flagged so that the last-write-wins behaviour of
//     many env loaders does not silently hide mistakes.
//
// Usage:
//
//	result, err := env.ValidateFile(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !result.OK() {
//		log.Fatalf("validation failed: %s", result.Error())
//	}
//
// You can also validate an already-parsed slice of entries:
//
//	entries, _ := env.Parse(".env")
//	result := env.ValidateEntries(entries)
package env
