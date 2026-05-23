// Package env provides utilities for parsing, writing, diffing, merging, and
// rotating .env files used by envcrypt.
//
// # Rotation
//
// Rotation is the process of replacing an encrypted .env file with a
// freshly-encrypted version (e.g. after adding a new team member's public key
// or revoking an old one) while preserving the previous version as a backup.
//
// Typical usage:
//
//	// 1. Decrypt the current file.
//	// 2. Re-encrypt with the updated recipient list.
//	// 3. Call Rotate to archive the old ciphertext and write the new one.
//	//
//	//   newCiphertext, err := crypto.EncryptFile(plaintext, recipients)
//	//   if err != nil { ... }
//	//   err = env.Rotate(".env.age", newCiphertext, env.RotateOptions{})
//
// Backups are stored in a "backups" directory next to the encrypted file by
// default and are named <original>.<timestamp> (UTC, format 20060102T150405Z).
//
// Use ListBackups to enumerate previous versions for auditing or rollback.
package env
