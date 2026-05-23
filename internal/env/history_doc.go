// Package env provides utilities for managing .env files, including
// parsing, writing, diffing, merging, validating, and tracking history.
//
// # History
//
// The history module records operations performed on env files over time.
// Each entry captures the operation type, the target file, an optional
// user identifier, and a timestamp.
//
// Typical usage:
//
//	err := env.AppendHistory(".envcrypt/history.json", env.HistoryEntry{
//		Operation: "encrypt",
//		File:      ".env",
//		User:      "alice",
//	})
//
// Supported operations include: "encrypt", "decrypt", "rotate", "merge".
//
// History files are stored as JSON arrays and are created automatically
// if they do not yet exist.
package env
