package env_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/envcrypt/internal/env"
)

func ExampleAddTag() {
	dir, _ := os.MkdirTemp("", "envcrypt-tag-example")
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "tags.json")

	tf, _ := env.LoadTagFile(path)
	_ = env.AddTag(&tf, env.Tag{
		Name:        "v1.0",
		Description: "initial production release",
		SnapshotRef: ".envcrypt/snapshots/2024-01-01T00:00:00Z.env",
	})
	_ = env.SaveTagFile(path, tf)

	loaded, _ := env.LoadTagFile(path)
	fmt.Println(loaded.Tags[0].Name)
	// Output: v1.0
}

func ExampleFindTag() {
	tf := env.TagFile{
		Tags: []env.Tag{
			{Name: "staging", SnapshotRef: "snap_staging"},
			{Name: "prod", SnapshotRef: "snap_prod"},
		},
	}
	tag, err := env.FindTag(tf, "prod")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(tag.SnapshotRef)
	// Output: snap_prod
}

func ExampleRemoveTag() {
	tf := env.TagFile{
		Tags: []env.Tag{
			{Name: "old", SnapshotRef: "snap_old"},
			{Name: "current", SnapshotRef: "snap_current"},
		},
	}
	_ = env.RemoveTag(&tf, "old")
	fmt.Println(len(tf.Tags))
	// Output: 1
}
