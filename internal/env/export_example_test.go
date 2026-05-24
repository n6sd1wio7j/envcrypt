package env_test

import (
	"fmt"

	"github.com/yourorg/envcrypt/internal/env"
)

func ExampleExport_dotenv() {
	entries := []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
	out, _ := env.Export(entries, env.ExportOptions{Format: env.FormatDotenv})
	fmt.Print(out)
	// Output:
	// DB_HOST=localhost
	// DB_PORT=5432
}

func ExampleExport_shell() {
	entries := []env.Entry{
		{Key: "TOKEN", Value: "abc123"},
	}
	out, _ := env.Export(entries, env.ExportOptions{Format: env.FormatShell})
	fmt.Print(out)
	// Output:
	// export TOKEN="abc123"
}

func ExampleExport_keyFilter() {
	entries := []env.Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}
	out, _ := env.Export(entries, env.ExportOptions{
		Format: env.FormatDotenv,
		Keys:   []string{"A", "C"},
	})
	fmt.Print(out)
	// Output:
	// A=1
	// C=3
}
