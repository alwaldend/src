package filesystem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"
)

// FS owns application path handling and I/O over caller-supplied storage.
// Its zero value selects OS storage. The backing filesystem is never exposed.
type FS struct {
	// storage supplies every underlying file operation and optional link capability.
	storage afero.Fs
}

// New selects the supplied filesystem or OS storage when it is nil.
func New(storage afero.Fs) FS {
	return FS{storage: storage}
}

// backend returns the supplied storage or the zero-value OS default.
func (storage FS) backend() afero.Fs {
	if storage.storage == nil {
		return afero.NewOsFs()
	}
	return storage.storage
}

// expandPath translates host notation and rooted names without cleaning unresolved components.
func (storage FS) expandPath(input string) (string, error) {
	if _, rooted := storage.backend().(*afero.BasePathFs); rooted {
		if filepath.VolumeName(input) != "" {
			return "", fmt.Errorf("rooted filesystem path %q contains a volume: %w", input, os.ErrPermission)
		}
		input = strings.TrimLeft(input, string(filepath.Separator))
		if input == "" {
			input = "."
		}
		return input, nil
	}
	if _, host := storage.backend().(*afero.OsFs); !host {
		if input == "" {
			input = "."
		}
		return input, nil
	}
	if input == "~" || strings.HasPrefix(input, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		input = home + string(filepath.Separator) + strings.TrimPrefix(strings.TrimPrefix(input, "~"), "/")
	}
	if filepath.IsAbs(input) {
		return input, nil
	}
	volume := filepath.VolumeName(input)
	remainder := strings.TrimPrefix(input, volume)
	anchor := volume + "."
	if strings.HasPrefix(remainder, string(filepath.Separator)) {
		anchor = volume + string(filepath.Separator)
		remainder = strings.TrimLeft(remainder, string(filepath.Separator))
	}
	base, err := filepath.Abs(anchor)
	if err != nil {
		return "", fmt.Errorf("resolve filesystem path: %w", err)
	}
	return base + string(filepath.Separator) + remainder, nil
}

// lstat inspects an already checked path without resolving its final link.
func (storage FS) lstat(path string) (os.FileInfo, error) {
	var info os.FileInfo
	var err error
	if capable, ok := storage.backend().(afero.Lstater); ok {
		info, _, err = capable.LstatIfPossible(path)
	} else {
		info, err = storage.backend().Stat(path)
	}
	if err != nil {
		return nil, fmt.Errorf("inspect filesystem path %q: %w", path, err)
	}
	return info, nil
}

// ResolvePath follows aliases in supplied storage, optionally retaining nonexistent components.
// Link targets are processed component by component: a parent component applies after preceding links.
func (storage FS) ResolvePath(input string, allowMissing bool) (string, error) {
	path, err := storage.resolvePath(input, allowMissing, false)
	if err != nil {
		return "", fmt.Errorf("resolve filesystem path %q: %w", input, err)
	}
	return path, nil
}

// ResolvePathWithoutLinks resolves components in order while rejecting every encountered symbolic link.
// Missing components may be retained for virtual parent normalization before callers inspect the result.
func (storage FS) ResolvePathWithoutLinks(input string, allowMissing bool) (string, error) {
	path, err := storage.resolvePath(input, allowMissing, true)
	if err != nil {
		return "", fmt.Errorf("resolve unlinked filesystem path %q: %w", input, err)
	}
	return path, nil
}

// resolvePath shares component traversal for alias-following and link-rejecting callers.
func (storage FS) resolvePath(input string, allowMissing, rejectLinks bool) (string, error) {
	path, err := storage.expandPath(input)
	if err != nil {
		return "", fmt.Errorf("expand filesystem path: %w", err)
	}
	prefix, pending := splitPath(path)
	links := 0
	for len(pending) != 0 {
		part := pending[0]
		pending = pending[1:]
		if part == "" {
			continue
		}
		path = filepath.Join(prefix, part)
		if _, rooted := storage.backend().(*afero.BasePathFs); rooted && !filepath.IsLocal(path) {
			return "", fmt.Errorf("path %q leaves the filesystem base: %w", input, os.ErrPermission)
		}
		info, err := storage.lstat(path)
		if err != nil {
			if allowMissing && errors.Is(err, os.ErrNotExist) {
				prefix = path
				continue
			}
			return "", fmt.Errorf("resolve filesystem path %q: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			if rejectLinks {
				return "", fmt.Errorf("symlink in filesystem path %q", path)
			}
			links++
			if links > 255 {
				return "", fmt.Errorf("too many symlinks resolving filesystem path %q", input)
			}
			reader, ok := storage.backend().(afero.LinkReader)
			if !ok {
				return "", fmt.Errorf("read filesystem link %q: %w", path, afero.ErrNoReadlink)
			}
			target, err := reader.ReadlinkIfPossible(path)
			if err != nil {
				return "", fmt.Errorf("read filesystem link %q: %w", path, err)
			}
			target, err = storage.linkTarget(path, target)
			if err != nil {
				return "", fmt.Errorf("resolve filesystem link %q: %w", path, err)
			}
			targetPrefix, parts := splitPath(target)
			if filepath.IsAbs(target) {
				prefix = targetPrefix
				if _, rooted := storage.backend().(*afero.BasePathFs); rooted {
					prefix = "."
				}
			}
			pending = append(parts, pending...)
			continue
		}
		if len(pending) != 0 && !info.IsDir() {
			return "", fmt.Errorf("filesystem ancestor %q must be a directory", path)
		}
		prefix = path
	}
	return prefix, nil
}

// splitPath separates a root from uncleaned components so link targets preserve parent semantics.
func splitPath(path string) (string, []string) {
	volume := filepath.VolumeName(path)
	prefix := volume + "."
	remainder := strings.TrimPrefix(path, volume)
	if filepath.IsAbs(path) {
		prefix = volume + string(filepath.Separator)
		remainder = strings.TrimLeft(remainder, string(filepath.Separator))
	}
	return prefix, strings.Split(remainder, string(filepath.Separator))
}

// linkTarget maps absolute backing-store targets into the rooted virtual namespace.
// Relative targets stay relative to the link's directory, without prematurely cleaning parent components.
func (storage FS) linkTarget(path, target string) (string, error) {
	if rooted, ok := storage.backend().(*afero.BasePathFs); ok && filepath.IsAbs(target) {
		root := filepath.Clean(afero.FullBaseFsPath(rooted, "."))
		if target == root {
			return string(filepath.Separator), nil
		}
		prefix := root
		if !strings.HasSuffix(prefix, string(filepath.Separator)) {
			prefix += string(filepath.Separator)
		}
		if !strings.HasPrefix(target, prefix) {
			return "", fmt.Errorf("symlink %q leaves the filesystem base: %w", path, os.ErrPermission)
		}
		return string(filepath.Separator) + strings.TrimPrefix(target, prefix), nil
	}
	return target, nil
}

// entryPath resolves an entry's ancestors while preserving the final entry for lstat, rename or removal.
func (storage FS) entryPath(input string) (string, error) {
	path, err := storage.expandPath(input)
	if err != nil {
		return "", fmt.Errorf("expand filesystem entry: %w", err)
	}
	directory, name := filepath.Split(path)
	parent, err := storage.ResolvePath(directory, true)
	if err != nil {
		return "", fmt.Errorf("resolve filesystem entry parent: %w", err)
	}
	path = filepath.Join(parent, name)
	if _, rooted := storage.backend().(*afero.BasePathFs); rooted && !filepath.IsLocal(path) {
		return "", fmt.Errorf("entry %q leaves the filesystem base: %w", input, os.ErrPermission)
	}
	return path, nil
}

// Lstat inspects an entry without following its final symlink when the backend exposes that capability.
func (storage FS) Lstat(path string) (os.FileInfo, error) {
	path, err := storage.entryPath(path)
	if err != nil {
		return nil, fmt.Errorf("resolve inspected entry: %w", err)
	}
	info, err := storage.lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect entry: %w", err)
	}
	return info, nil
}

// Stat inspects the resolved target within the supplied namespace.
func (storage FS) Stat(path string) (os.FileInfo, error) {
	path, err := storage.ResolvePath(path, false)
	if err != nil {
		return nil, fmt.Errorf("resolve stat target: %w", err)
	}
	info, err := storage.backend().Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %q: %w", path, err)
	}
	return info, nil
}

// Open opens an existing resolved file for reading.
func (storage FS) Open(path string) (afero.File, error) {
	path, err := storage.ResolvePath(path, false)
	if err != nil {
		return nil, fmt.Errorf("resolve input file: %w", err)
	}
	file, err := storage.backend().Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input %q: %w", path, err)
	}
	return file, nil
}

// OpenFile resolves a path and opens it with the supplied flags and creation permissions.
func (storage FS) OpenFile(path string, flags int, mode os.FileMode) (afero.File, error) {
	path, err := storage.ResolvePath(path, flags&os.O_CREATE != 0)
	if err != nil {
		return nil, fmt.Errorf("resolve opened file: %w", err)
	}
	file, err := storage.backend().OpenFile(path, flags, mode)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}
	return file, nil
}

// ReadFile reads a complete file and retains read and close failures.
func (storage FS) ReadFile(path string) ([]byte, error) {
	file, err := storage.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file contents %q: %w", path, err)
	}
	data, readErr := io.ReadAll(file)
	if err := errors.Join(readErr, file.Close()); err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}
	return data, nil
}

// WriteFile replaces file contents and retains write and close failures.
func (storage FS) WriteFile(path string, data []byte, mode os.FileMode) error {
	file, err := storage.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("open output %q: %w", path, err)
	}
	written, writeErr := file.Write(data)
	if written != len(data) && writeErr == nil {
		writeErr = io.ErrShortWrite
	}
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return fmt.Errorf("write file %q: %w", path, err)
	}
	return nil
}

// ReadDir returns sorted directory entries and retains read and close failures.
func (storage FS) ReadDir(path string) ([]os.FileInfo, error) {
	file, err := storage.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open directory %q: %w", path, err)
	}
	entries, readErr := file.Readdir(-1)
	if err := errors.Join(readErr, file.Close()); err != nil {
		return nil, fmt.Errorf("read directory %q: %w", path, err)
	}
	sort.Slice(entries, func(left, right int) bool { return entries[left].Name() < entries[right].Name() })
	return entries, nil
}

// MkdirAll creates a resolved directory and any missing ancestors.
func (storage FS) MkdirAll(path string, mode os.FileMode) error {
	path, err := storage.ResolvePath(path, true)
	if err != nil {
		return fmt.Errorf("resolve directory creation: %w", err)
	}
	if err := storage.backend().MkdirAll(path, mode); err != nil {
		return fmt.Errorf("create directory %q: %w", path, err)
	}
	return nil
}

// Rename moves entries within supplied storage without following their final symlinks.
func (storage FS) Rename(source, destination string) error {
	source, err := storage.entryPath(source)
	if err != nil {
		return fmt.Errorf("resolve rename source: %w", err)
	}
	destination, err = storage.entryPath(destination)
	if err != nil {
		return fmt.Errorf("resolve rename destination: %w", err)
	}
	if err := storage.backend().Rename(source, destination); err != nil {
		return fmt.Errorf("rename %q to %q: %w", source, destination, err)
	}
	return nil
}

// Remove removes an entry without following its final symlink.
func (storage FS) Remove(path string) error {
	path, err := storage.entryPath(path)
	if err != nil {
		return fmt.Errorf("resolve removed entry: %w", err)
	}
	if err := storage.backend().Remove(path); err != nil {
		return fmt.Errorf("remove %q: %w", path, err)
	}
	return nil
}

// RemoveAll removes an entry tree without following its final symlink.
func (storage FS) RemoveAll(path string) error {
	path, err := storage.entryPath(path)
	if err != nil {
		return fmt.Errorf("resolve removed tree: %w", err)
	}
	if err := storage.backend().RemoveAll(path); err != nil {
		return fmt.Errorf("remove tree %q: %w", path, err)
	}
	return nil
}

// TempFile creates a file and returns its stable virtual pathname separately from the layered handle name.
// The directory is explicit; an empty directory selects ".", never host temp storage.
func (storage FS) TempFile(directory, pattern string) (afero.File, string, error) {
	directory, err := storage.ResolvePath(directory, false)
	if err != nil {
		return nil, "", fmt.Errorf("resolve temporary file directory: %w", err)
	}
	file, err := afero.TempFile(storage.backend(), directory, pattern)
	if err != nil {
		return nil, "", fmt.Errorf("create temporary file in %q: %w", directory, err)
	}
	return file, filepath.Join(directory, filepath.Base(file.Name())), nil
}

// TempDir creates a private directory within an explicitly resolved filesystem directory.
func (storage FS) TempDir(directory, prefix string) (string, error) {
	directory, err := storage.ResolvePath(directory, false)
	if err != nil {
		return "", fmt.Errorf("resolve temporary directory parent: %w", err)
	}
	path, err := afero.TempDir(storage.backend(), directory, prefix)
	if err != nil {
		return "", fmt.Errorf("create temporary directory in %q: %w", directory, err)
	}
	return path, nil
}
