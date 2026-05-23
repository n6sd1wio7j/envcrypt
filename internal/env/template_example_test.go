package env_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourorg/envcrypt/internal/env"
)

func ExampleGenerateTemplate() {
	dir, _ := os.MkdirTemp("", "envcrypt-example")
	defer os.RemoveAll(dir)

	example := filepath.Join(dir, ".env.example")
	_ = os.WriteFile(example, []byte("DB_HOST=localhost # database host\nDB_PORT=5432\n"), 0600)

	tmpl, err := env.GenerateTemplate(example)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, e := range tmpl.Entries {
		fmt.Printf("key=%s required=%v desc=%q\n", e.Key, e.Required, e.Description)
	}

	// Output:
	// key=DB_HOST required=true desc="database host"
	// key=DB_PORT required=true desc=""
}

func ExampleCheckTemplate() {
	tmpl := &env.Template{
		Entries: []env.TemplateEntry{
			{Key: "DB_HOST", Required: true},
			{Key: "DB_PORT", Required: true},
			{Key: "OPTIONAL_KEY", Required: false},
		},
	}

	entries := []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
	}

	missing := env.CheckTemplate(tmpl, entries)
	fmt.Println("missing:", missing)

	// Output:
	// missing: [DB_PORT]
}
