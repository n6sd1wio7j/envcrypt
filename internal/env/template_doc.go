// Package env provides utilities for managing .env files.
//
// # Template
//
// The template feature allows generating a canonical list of expected
// environment variable keys from an existing .env file, and validating
// that a given set of entries satisfies all required keys.
//
// Usage:
//
//	// Generate a template from a reference .env file
//	tmpl, err := env.GenerateTemplate(".env.example")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Load the actual .env entries
//	entries, err := env.Parse(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check for missing required keys
//	missing := env.CheckTemplate(tmpl, entries)
//	if len(missing) > 0 {
//		fmt.Println("Missing keys:", missing)
//	}
package env
