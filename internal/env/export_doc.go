// Package env provides utilities for parsing, writing, diffing,
// merging, validating, templating, snapshotting, auditing, and exporting
// .env files used by envcrypt.
//
// # Export
//
// The Export function serialises a slice of [Entry] values into one of
// three formats:
//
//   - FormatDotenv  – standard KEY=VALUE lines (default)
//   - FormatShell   – export KEY="VALUE" lines suitable for sourcing
//   - FormatJSON    – a simple JSON object mapping keys to values
//
// An optional key filter lets callers export only a subset of entries:
//
//	content, err := env.Export(entries, env.ExportOptions{
//		Format: env.FormatShell,
//		Keys:   []string{"DB_HOST", "DB_PORT"},
//	})
//
// ExportToFile is a convenience wrapper that writes the result directly
// to a file with mode 0600.
package env
