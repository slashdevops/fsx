// Package fsx provides small filesystem helpers built only on the Go standard
// library.
//
// The package focuses on behavior that is useful across command-line tools and
// services:
//
//   - expanding user-facing paths that start with ~ or contain $HOME
//   - checking whether paths exist and whether they are regular files or directories
//   - confirming that a path remains inside a base directory after normalization
//   - matching file extensions case-insensitively
//   - writing files atomically by replacing the destination with a temporary file
//
// # Path Expansion
//
// [ExpandPath] expands a leading ~ with the current user's home directory and
// replaces $HOME references with the HOME environment variable. If the home
// directory cannot be determined, the original path is returned unchanged.
//
// # Atomic Writes
//
// [WriteFileAtomic] writes data to a temporary file in the destination directory,
// sets the requested permissions, and renames the temporary file over the target.
// This prevents readers from observing partially written file contents on the
// same filesystem.
//
// # Dependencies
//
// fsx has zero third-party dependencies. It uses os, path/filepath, strings, and
// other standard library packages.
package fsx
