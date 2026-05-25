// Package env provides utilities for parsing, writing, and manipulating .env files.
//
// # Search
//
// The Search function allows filtering entries in a .env file by key and/or
// value using regular expressions.
//
// Basic usage:
//
//	results, err := env.Search(".env", env.SearchOptions{
//		KeyPattern: "^DB_",
//	})
//
// Both key and value patterns may be provided simultaneously:
//
//	results, err := env.Search(".env", env.SearchOptions{
//		KeyPattern:   "SECRET",
//		ValuePattern: "^prod",
//		CaseSensitive: false,
//	})
//
// To search across multiple files use SearchMultiple:
//
//	results, err := env.SearchMultiple([]string{".env", ".env.prod"}, opts)
//
// Results include both the matched Entry and the source file path, making it
// straightforward to report findings across environments.
package env
