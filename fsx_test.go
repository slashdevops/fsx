package fsx

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}

	t.Setenv("HOME", homeDir)

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "empty path",
			path: "",
			want: "",
		},
		{
			name: "tilde only",
			path: "~",
			want: homeDir,
		},
		{
			name: "tilde prefix",
			path: "~/config.yaml",
			want: filepath.Join(homeDir, "config.yaml"),
		},
		{
			name: "HOME only",
			path: "$HOME",
			want: homeDir,
		},
		{
			name: "HOME prefix",
			path: "$HOME/config.yaml",
			want: filepath.Join(homeDir, "config.yaml"),
		},
		{
			name: "multiple HOME references",
			path: "$HOME/$HOME/test",
			want: homeDir + string(filepath.Separator) + homeDir + string(filepath.Separator) + "test",
		},
		{
			name: "plain absolute path",
			path: filepath.FromSlash("/etc/config.yaml"),
			want: filepath.FromSlash("/etc/config.yaml"),
		},
		{
			name: "relative path",
			path: filepath.FromSlash("config/file.yaml"),
			want: filepath.FromSlash("config/file.yaml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandPath(tt.path); got != tt.want {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestExpandPathKeepsHOMEWhenUnset(t *testing.T) {
	t.Setenv("HOME", "")

	const path = "$HOME/config.yaml"
	if got := ExpandPath(path); got != path {
		t.Errorf("ExpandPath(%q) = %q, want original path", path, got)
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(file, []byte("content"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "empty path",
			path: "",
			want: false,
		},
		{
			name: "existing file",
			path: file,
			want: true,
		},
		{
			name: "existing directory",
			path: tmpDir,
			want: true,
		},
		{
			name: "missing path",
			path: filepath.Join(tmpDir, "missing.txt"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.path); got != tt.want {
				t.Errorf("Exists(%q) = %t, want %t", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(file, []byte("content"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "empty path",
			path: "",
			want: false,
		},
		{
			name: "directory",
			path: tmpDir,
			want: true,
		},
		{
			name: "file",
			path: file,
			want: false,
		},
		{
			name: "missing path",
			path: filepath.Join(tmpDir, "missing"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDir(tt.path); got != tt.want {
				t.Errorf("IsDir(%q) = %t, want %t", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(file, []byte("content"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "empty path",
			path: "",
			want: false,
		},
		{
			name: "regular file",
			path: file,
			want: true,
		},
		{
			name: "directory",
			path: tmpDir,
			want: false,
		},
		{
			name: "missing path",
			path: filepath.Join(tmpDir, "missing.txt"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.path); got != tt.want {
				t.Errorf("IsFile(%q) = %t, want %t", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsWithin(t *testing.T) {
	tmpDir := t.TempDir()
	base := filepath.Join(tmpDir, "base")
	if err := os.MkdirAll(filepath.Join(base, "nested"), 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	tests := []struct {
		name   string
		base   string
		target string
		want   bool
	}{
		{
			name:   "same path",
			base:   base,
			target: base,
			want:   true,
		},
		{
			name:   "file directly inside base",
			base:   base,
			target: filepath.Join(base, "file.txt"),
			want:   true,
		},
		{
			name:   "nested file inside base",
			base:   base,
			target: filepath.Join(base, "nested", "file.txt"),
			want:   true,
		},
		{
			name:   "sibling path outside base",
			base:   base,
			target: filepath.Join(tmpDir, "sibling", "file.txt"),
			want:   false,
		},
		{
			name:   "parent traversal outside base",
			base:   base,
			target: filepath.Join(base, "..", "sibling", "file.txt"),
			want:   false,
		},
		{
			name:   "empty base",
			base:   "",
			target: filepath.Join(base, "file.txt"),
			want:   false,
		},
		{
			name:   "empty target",
			base:   base,
			target: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWithin(tt.base, tt.target); got != tt.want {
				t.Errorf("IsWithin(%q, %q) = %t, want %t", tt.base, tt.target, got, tt.want)
			}
		})
	}
}

func TestHasExtension(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		extensions []string
		want       bool
	}{
		{
			name:       "empty path",
			path:       "",
			extensions: []string{"txt"},
			want:       false,
		},
		{
			name:       "no extensions",
			path:       "file.txt",
			extensions: nil,
			want:       false,
		},
		{
			name:       "extension without dot",
			path:       "file.txt",
			extensions: []string{"txt"},
			want:       true,
		},
		{
			name:       "extension with dot",
			path:       "file.txt",
			extensions: []string{".txt"},
			want:       true,
		},
		{
			name:       "case insensitive",
			path:       "file.YAML",
			extensions: []string{"yaml"},
			want:       true,
		},
		{
			name:       "one of many extensions",
			path:       "file.md",
			extensions: []string{"txt", "md"},
			want:       true,
		},
		{
			name:       "mismatch",
			path:       "file.md",
			extensions: []string{"txt", "yaml"},
			want:       false,
		},
		{
			name:       "path without extension",
			path:       "Makefile",
			extensions: []string{"txt"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasExtension(tt.path, tt.extensions...); got != tt.want {
				t.Errorf("HasExtension(%q, %q) = %t, want %t", tt.path, tt.extensions, got, tt.want)
			}
		})
	}
}

func TestDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}

	t.Setenv("HOME", homeDir)

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "absolute file",
			path: filepath.FromSlash("/etc/config/file.yaml"),
			want: filepath.FromSlash("/etc/config"),
		},
		{
			name: "relative file",
			path: filepath.FromSlash("config/file.yaml"),
			want: "config",
		},
		{
			name: "tilde path",
			path: "~/config/file.yaml",
			want: filepath.Join(homeDir, "config"),
		},
		{
			name: "HOME path",
			path: "$HOME/config/file.yaml",
			want: filepath.Join(homeDir, "config"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dir(tt.path); got != tt.want {
				t.Errorf("Dir(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestWriteFileAtomic(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "data.txt")
	data := []byte("atomic content")

	if err := WriteFileAtomic(path, data, 0o600); err != nil {
		t.Fatalf("WriteFileAtomic() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("os.ReadFile() = %q, want %q", got, data)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Errorf("file permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestWriteFileAtomicReplacesExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "data.txt")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	if err := WriteFileAtomic(path, []byte("new"), 0o644); err != nil {
		t.Fatalf("WriteFileAtomic() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != "new" {
		t.Errorf("os.ReadFile() = %q, want %q", got, "new")
	}
}

func TestWriteFileAtomicEmptyPath(t *testing.T) {
	err := WriteFileAtomic("", []byte("content"), 0o644)
	if err == nil {
		t.Fatal("WriteFileAtomic() error = nil, want non-nil error")
	}
	if !strings.Contains(err.Error(), "path is empty") {
		t.Errorf("WriteFileAtomic() error = %v, want empty path error", err)
	}
}

func TestWriteFileAtomicMissingDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "data.txt")
	err := WriteFileAtomic(path, []byte("content"), 0o644)
	if err == nil {
		t.Fatal("WriteFileAtomic() error = nil, want non-nil error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("WriteFileAtomic() error = %v, want os.ErrNotExist", err)
	}
}

func BenchmarkExpandPath(b *testing.B) {
	for b.Loop() {
		ExpandPath("~/config.yaml")
	}
}

func BenchmarkExists(b *testing.B) {
	path := filepath.Join(b.TempDir(), "data.txt")
	if err := os.WriteFile(path, []byte("content"), 0o644); err != nil {
		b.Fatalf("os.WriteFile() error = %v", err)
	}

	for b.Loop() {
		Exists(path)
	}
}

func BenchmarkWriteFileAtomic(b *testing.B) {
	path := filepath.Join(b.TempDir(), "data.txt")
	data := []byte("content")

	for b.Loop() {
		if err := WriteFileAtomic(path, data, 0o644); err != nil {
			b.Fatalf("WriteFileAtomic() error = %v", err)
		}
	}
}
