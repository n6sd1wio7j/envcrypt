// Package env provides utilities for managing .env files, including
// parsing, writing, diffing, merging, rotating, snapshotting, auditing,
// locking, validating, schema enforcement, template generation, history
// tracking, and tagging.
//
// # Tag
//
// Tags provide a lightweight mechanism to mark named points in the history
// of an encrypted environment file. Each tag references a snapshot by its
// path or identifier and carries an optional human-readable description.
//
// Tags are persisted as JSON in a file (default: .envcrypt/tags.json) and
// support the following operations:
//
//   - [LoadTagFile] — load all tags from disk (returns empty set if missing)
//   - [SaveTagFile] — persist tags to disk, creating directories as needed
//   - [AddTag]     — append a new uniquely-named tag
//   - [RemoveTag]  — delete a tag by name
//   - [FindTag]    — look up a tag by name
//
// Example usage:
//
//	tf, _ := env.LoadTagFile(".envcrypt/tags.json")
//	_ = env.AddTag(&tf, env.Tag{
//		Name:        "v2.0",
//		Description: "post-migration secrets",
//		SnapshotRef: ".envcrypt/snapshots/2024-01-15T10:00:00Z.env",
//	})
//	_ = env.SaveTagFile(".envcrypt/tags.json", tf)
package env
