# fsx

[![main branch](https://github.com/slashdevops/fsx/actions/workflows/main.yml/badge.svg)](https://github.com/slashdevops/fsx/actions/workflows/main.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/slashdevops/fsx?style=plastic)
[![Go Reference](https://pkg.go.dev/badge/github.com/slashdevops/fsx.svg)](https://pkg.go.dev/github.com/slashdevops/fsx)
[![Go Report Card](https://goreportcard.com/badge/github.com/slashdevops/fsx)](https://goreportcard.com/report/github.com/slashdevops/fsx)
[![license](https://img.shields.io/github/license/slashdevops/fsx.svg)](https://github.com/slashdevops/fsx/blob/main/LICENSE)
[![Release](https://github.com/slashdevops/fsx/actions/workflows/release.yml/badge.svg)](https://github.com/slashdevops/fsx/actions/workflows/release.yml)

`fsx` is a small Go package for filesystem helpers that sit just above `os` and `path/filepath`. It is designed for applications that need compact, easy-to-review helpers for expanded user paths, containment checks, file predicates, extension matching, and atomic writes.

The package is standard-library-only and intentionally keeps one short package name instead of exposing broad `utils` packages.

## Features

- User path expansion for leading `~` and `$HOME`
- Path existence checks for files, directories, and any filesystem entry
- Directory containment checks that normalize relative traversal
- Case-insensitive file extension matching
- Atomic file replacement with permissions
- Zero third-party dependencies
- Apache-2.0 licensed

## Installation

```sh
go get github.com/slashdevops/fsx
```

Update to the latest available version:

```sh
go get -u github.com/slashdevops/fsx
```

## Requirements

- Go 1.26.3 or newer
- No external Go modules

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/slashdevops/fsx"
)

func main() {
    configPath := fsx.ExpandPath("~/app/config.yaml")
    if !fsx.HasExtension(configPath, "yaml") {
        log.Fatal("config file must be YAML")
    }

    if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
        log.Fatal(err)
    }

    if err := fsx.WriteFileAtomic(configPath, []byte("enabled: true\n"), 0o600); err != nil {
        log.Fatal(err)
    }

    fmt.Println(filepath.Base(configPath))
}
```

## API

### ExpandPath

```go
func ExpandPath(path string) string
```

`ExpandPath` replaces a leading `~` with the current user's home directory and replaces `$HOME` references with the `HOME` environment variable. Empty paths are returned unchanged.

```go
configPath := fsx.ExpandPath("~/app/config.yaml")
```

### Exists

```go
func Exists(path string) bool
```

`Exists` reports whether a file, directory, or other filesystem entry exists. The path is expanded before it is checked.

```go
if fsx.Exists("~/app/config.yaml") {
    // read config
}
```

### IsDir

```go
func IsDir(path string) bool
```

`IsDir` reports whether the expanded path exists and is a directory.

```go
if !fsx.IsDir("~/app") {
    return errors.New("app directory does not exist")
}
```

### IsFile

```go
func IsFile(path string) bool
```

`IsFile` reports whether the expanded path exists and is a regular file.

```go
if !fsx.IsFile("~/app/config.yaml") {
    return errors.New("config file does not exist")
}
```

### IsWithin

```go
func IsWithin(base, target string) bool
```

`IsWithin` reports whether `target` resolves to a path contained inside `base`. Both paths are expanded, cleaned, and made absolute before comparison. A target equal to base is considered within base.

This is useful before deleting or overwriting paths read from user-editable state files.

```go
if !fsx.IsWithin(outputDir, stateFilePath) {
    return fmt.Errorf("refusing to remove file outside output directory")
}
```

### HasExtension

```go
func HasExtension(path string, extensions ...string) bool
```

`HasExtension` reports whether `path` has one of the provided extensions. Matching is case-insensitive, extensions may be passed with or without a leading dot, and the path does not need to exist.

```go
if !fsx.HasExtension("config.YAML", "yaml", "yml") {
    return errors.New("config must be YAML")
}
```

### Dir

```go
func Dir(path string) string
```

`Dir` returns the directory component of `path` after expanding it.

```go
dir := fsx.Dir("~/app/config.yaml")
```

### WriteFileAtomic

```go
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error
```

`WriteFileAtomic` writes data to a temporary file in the destination directory, sets the requested permissions, and renames the temporary file over the target path.

The destination directory must already exist. If any step fails before the rename, the temporary file is removed.

```go
if err := fsx.WriteFileAtomic("~/app/config.yaml", data, 0o600); err != nil {
    return fmt.Errorf("save config: %w", err)
}
```

## Quality Automation

The repository follows the same GitHub quality practices used by `slashdevops/e5t`:

- `main.yml` validates pushes to `main` with format checks, `go vet`, race-enabled tests, coverage, and `go build`.
- `pr.yml` runs the same quality gate for pull requests targeting `main`.
- `codeql.yml` runs CodeQL for Go and GitHub Actions on pushes, pull requests, and a weekly schedule.
- `release.yml` validates tagged releases and publishes GitHub releases with generated notes.
- `dependabot.yml` checks Go module and GitHub Actions updates weekly.

## Project Structure

```text
.
|-- .github/         GitHub Actions, CodeQL, Dependabot, and release metadata
|-- .golangci.yaml   Optional local golangci-lint configuration
|-- doc.go           Package documentation rendered by pkg.go.dev
|-- example_test.go  Executable Go examples for documentation
|-- fsx.go           Public filesystem API
|-- fsx_test.go      Unit tests and benchmarks
|-- go.mod           Module definition with no external requirements
|-- LICENSE          Apache License 2.0
|-- README.md        Project overview and usage guide
`-- SECURITY.md      Vulnerability reporting policy
```

## Testing

Run the test suite:

```sh
go test ./...
```

Run benchmarks:

```sh
go test -bench=. ./...
```

Check test coverage:

```sh
go test -cover ./...
```

Run the same local quality checks used by CI:

```sh
go fix ./...
go fmt ./...
go vet ./...
go test -race -coverprofile=/tmp/fsx-coverage.txt -covermode=atomic ./...
go build ./...
```

## License

`fsx` is licensed under the [Apache License 2.0](LICENSE).

## Contributing

Issues and pull requests are welcome at [github.com/slashdevops/fsx](https://github.com/slashdevops/fsx). Please keep changes small, idiomatic, tested, documented, and dependency-free unless there is a clear reason to expand the project scope.
