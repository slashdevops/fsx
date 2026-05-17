package fsx_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/slashdevops/fsx"
)

// ExampleExpandPath demonstrates expanding a user-facing path.
func ExampleExpandPath() {
	expanded := fsx.ExpandPath("~/config.yaml")
	fmt.Println(filepath.IsAbs(expanded))

	// Output:
	// true
}

// ExampleIsWithin demonstrates a containment check before deleting a file.
func ExampleIsWithin() {
	base := filepath.Join(os.TempDir(), "workspace")
	target := filepath.Join(base, "page.md")

	fmt.Println(fsx.IsWithin(base, target))
	fmt.Println(fsx.IsWithin(base, filepath.Join(base, "..", "escape.md")))

	// Output:
	// true
	// false
}

// ExampleHasExtension demonstrates extension matching.
func ExampleHasExtension() {
	fmt.Println(fsx.HasExtension("config.YAML", "yaml", "json"))
	fmt.Println(fsx.HasExtension("README", "md"))

	// Output:
	// true
	// false
}

// ExampleWriteFileAtomic demonstrates an atomic file replacement.
func ExampleWriteFileAtomic() {
	dir, err := os.MkdirTemp("", "fsx-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			log.Fatal(err)
		}
	}()

	path := filepath.Join(dir, "config.txt")
	if err := fsx.WriteFileAtomic(path, []byte("ready"), 0o600); err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	// Output:
	// ready
}
