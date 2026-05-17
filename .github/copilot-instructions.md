# Development Guidelines

This document contains the critical information about working with the project codebase.
Follow these guidelines precisely to ensure consistency and maintainability of the code.

## Stack

- Language: Go (Go 1.26.3)
- Framework: Go standard library
- Testing: Go's built-in testing package
- Dependency Management: Go modules
- Version Control: Git
- Documentation: go doc
- Code Review: Pull requests on GitHub
- CI/CD: GitHub Actions

## Project Structure

Since this is a small Go library, files are organized in the root directory following standard Go package layout.

- Library files are located in the root directory.
- `.github/` contains GitHub-specific files such as workflows for CI/CD.
- `.gitignore` specifies files and directories to be ignored by Git.
- `LICENSE` is the license file for the project.
- `README.md` provides an overview of the project, installation instructions, usage examples, and other relevant information.
- `go.mod` declares the module and Go toolchain version.
- `go.sum` should not exist unless external dependencies are intentionally introduced.
- `*.go` files contain the main source code of the library.
- `*_test.go` files contain tests, benchmarks, and executable examples.

## Code Style

- Follow Go's idiomatic style defined in
  - <https://google.github.io/styleguide/go/guide>
  - <https://google.github.io/styleguide/go/decisions>
  - <https://google.github.io/styleguide/go/best-practices>
  - <https://golang.org/doc/effective_go.html>
- Keep package names short, lowercase, and non-stuttering.
- Use meaningful names for variables, functions, and packages.
- Keep functions small and focused on a single task.
- Add comments for exported identifiers and complex behavior.
- Do not use `interface{}`; use `any` when an unconstrained type is required.
- Do not add external module dependencies without explicit approval.

## Go 1.26 Practices

- Prefer `errors.AsType[T](err)` over `errors.As(err, &target)` in new code.
- Use `new(value)` when allocating and initializing pointer values.
- Use `for b.Loop()` in benchmarks.
- Run `go fix ./...` after larger changes to apply available modernizations.

## Post-Change Checklist

Use standard Go commands after making changes:

```bash
go fix ./...
go fmt ./...
go vet ./...
betteralign -apply ./...
go test -race -coverprofile=/tmp/fsx-coverage.txt -covermode=atomic ./...
go build ./...
```
