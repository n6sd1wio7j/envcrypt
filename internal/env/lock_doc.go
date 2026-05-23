// Package env provides utilities for parsing, writing, diffing, merging,
// rotating, snapshotting, auditing, and locking .env files used by envcrypt.
//
// # Lock File
//
// The lock file (.envcrypt.lock) tracks the relationship between plaintext
// .env files and their encrypted counterparts. It records:
//
//   - The source .env file path
//   - The encrypted output file path
//   - A checksum of the plaintext content (for change detection)
//   - The timestamp and author of the last encryption
//
// This allows envcrypt to detect when a .env file has changed since it was
// last encrypted, warn about stale ciphertext, and provide an audit trail
// of who last updated each secret file.
//
// # Usage
//
//	lf, err := env.LoadLockFile(".envcrypt.lock")
//	if err != nil { ... }
//
//	lf.Upsert(env.LockEntry{
//		File:      ".env",
//		Encrypted: ".env.age",
//		Checksum:  checksum,
//		UpdatedAt: time.Now().UTC(),
//		UpdatedBy: alias,
//	})
//
//	if err := env.SaveLockFile(".envcrypt.lock", lf); err != nil { ... }
package env
