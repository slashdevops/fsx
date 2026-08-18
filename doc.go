// Package fsx provides small filesystem helpers built only on the Go standard
// library and targets the latest Go toolchain version declared by the module.
//
// The package focuses on behavior that is useful across command-line tools and
// services:
//
//   - expanding user-facing paths that start with ~ or contain $HOME
//   - checking whether paths exist and whether they are regular files or directories
//   - checking whether a regular file is marked readable or writable by its owner
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
// # Permission Predicates
//
// [IsReadable] and [IsWritable] inspect the file mode, not the effective access of the calling process.
// That is the portable answer the standard library can give: os.Stat exposes mode bits everywhere, while
// asking "can I open this?" needs access(2) or an attempted open. A caller that must be certain should
// open the file and handle the error — the answer can change between any check and the open regardless.
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
// fsx has zero third-party dependencies and no external module requirements. It
// uses os, path/filepath, strings, and other standard library packages, and the
// module currently targets Go 1.26.3.
package fsx
