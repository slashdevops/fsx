package fsx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands a user-facing filesystem path.
//
// It replaces a leading ~ with the current user's home directory and replaces
// all $HOME references with the HOME environment variable. Empty paths are
// returned unchanged. If the home directory cannot be determined, the original
// path is returned unchanged.
func ExpandPath(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		return homeDir + path[1:]
	}

	if strings.Contains(path, "$HOME") {
		homeDir := os.Getenv("HOME")
		if homeDir != "" {
			return strings.ReplaceAll(path, "$HOME", homeDir)
		}
	}

	return path
}

// Exists reports whether a file, directory, or other filesystem entry exists.
//
// The path is expanded with [ExpandPath] before it is checked. Empty paths and
// paths that cannot be statted return false.
func Exists(path string) bool {
	if path == "" {
		return false
	}

	_, err := os.Stat(ExpandPath(path))
	return err == nil
}

// IsDir reports whether path exists and is a directory.
//
// The path is expanded with [ExpandPath] before it is checked.
func IsDir(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(ExpandPath(path))
	if err != nil {
		return false
	}

	return info.IsDir()
}

// IsFile reports whether path exists and is a regular file.
//
// The path is expanded with [ExpandPath] before it is checked.
func IsFile(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(ExpandPath(path))
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}

// IsWithin reports whether target resolves to a path contained inside base.
//
// Both paths are expanded, cleaned, and made absolute before comparison. The
// function returns false when either path is empty, either path cannot be
// resolved, or the relative path from base to target escapes base with "..".
// A target equal to base is considered within base.
func IsWithin(base, target string) bool {
	if base == "" || target == "" {
		return false
	}

	absBase, err := filepath.Abs(ExpandPath(base))
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(ExpandPath(target))
	if err != nil {
		return false
	}

	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return false
	}

	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// HasExtension reports whether path has one of the provided extensions.
//
// Extension matching is case-insensitive. Extensions may be passed with or
// without a leading dot. The path does not need to exist.
func HasExtension(path string, extensions ...string) bool {
	if path == "" || len(extensions) == 0 {
		return false
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}

	for _, candidate := range extensions {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if !strings.HasPrefix(candidate, ".") {
			candidate = "." + candidate
		}
		if ext == strings.ToLower(candidate) {
			return true
		}
	}

	return false
}

// Dir returns the directory component of path after expanding it with
// [ExpandPath].
func Dir(path string) string {
	return filepath.Dir(ExpandPath(path))
}

// WriteFileAtomic writes data to path by replacing it with a completed temporary
// file in the same directory.
//
// The destination directory must already exist. The temporary file is created in
// that directory, written, closed, chmodded to perm, and renamed to path. If any
// step fails before the rename, the temporary file is removed.
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	path = ExpandPath(path)
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, ".fsx-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	tempPath := tempFile.Name()
	removeTemp := true
	defer func() {
		if removeTemp {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("write temporary file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}
	if err := os.Chmod(tempPath, perm); err != nil {
		return fmt.Errorf("set temporary file permissions: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("rename temporary file: %w", err)
	}

	removeTemp = false
	return nil
}
