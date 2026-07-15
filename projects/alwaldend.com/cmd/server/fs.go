// https://github.com/go-again/tarfs/blob/main/tarfs.go
package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"time"
)

// Compile-time interface assertions.
var (
	_ fs.FS         = (*FS)(nil)
	_ fs.ReadDirFS  = (*FS)(nil)
	_ fs.ReadFileFS = (*FS)(nil)
)

// FS is a read-only in-memory filesystem loaded from a tar archive.
// Construct it with [New], [NewGzip], [NewZstd], or [NewFromReader].
// All methods are safe for concurrent use after construction.
type FS struct {
	entries map[string]*entry   // clean path → entry ("." for root)
	roots   map[string][]string // dir path → sorted child full paths
}

type entry struct {
	name    string
	path    string // full clean path, e.g. "dist/index.html"; "." for root
	isDir   bool
	size    int64
	mode    fs.FileMode
	modTime time.Time
	data    []byte // nil for directories
}

// New builds an FS from raw (uncompressed) tar data.
func New(data []byte) (*FS, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("tarfs: empty archive")
	}
	return NewFromReader(bytes.NewReader(data))
}


// NewFromReader builds an FS by reading a raw (uncompressed) tar stream from r.
// For compressed archives, wrap r with the appropriate decompressor before calling.
func NewFromReader(r io.Reader) (*FS, error) {
	tfs := &FS{
		entries: make(map[string]*entry),
		roots:   make(map[string][]string),
	}
	// Root always exists.
	tfs.entries["."] = &entry{
		name:  ".",
		path:  ".",
		isDir: true,
		mode:  fs.ModeDir | 0o555,
	}

	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tarfs: tar read: %w", err)
		}

		// Normalize: strip leading "./" and clean.
		clean := path.Clean(hdr.Name)
		if clean == "." || clean == "" {
			continue
		}

		mode := hdr.FileInfo().Mode()
		isDir := hdr.Typeflag == tar.TypeDir || mode.IsDir()

		e := &entry{
			name:    path.Base(clean),
			path:    clean,
			isDir:   isDir,
			size:    hdr.Size,
			mode:    mode,
			modTime: hdr.ModTime,
		}

		if !isDir {
			e.data, err = io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("tarfs: read %s: %w", clean, err)
			}
			e.size = int64(len(e.data))
		}

		tfs.entries[clean] = e
		tfs.ensureParents(clean, hdr.ModTime)
	}

	// Build sorted child lists for every directory.
	for p := range tfs.entries {
		if p == "." {
			continue
		}
		parent := path.Dir(p)
		if parent == "" {
			parent = "."
		}
		tfs.roots[parent] = append(tfs.roots[parent], p)
	}
	for dir := range tfs.roots {
		sort.Strings(tfs.roots[dir])
	}

	return tfs, nil
}

// ensureParents synthesizes directory entries for every ancestor of p that
// does not yet exist in the entries map.
func (f *FS) ensureParents(p string, modTime time.Time) {
	for {
		parent := path.Dir(p)
		if parent == p {
			break
		}
		p = parent
		if p == "." {
			break
		}
		if _, ok := f.entries[p]; !ok {
			f.entries[p] = &entry{
				name:    path.Base(p),
				path:    p,
				isDir:   true,
				mode:    fs.ModeDir | 0o555,
				modTime: modTime,
			}
		}
	}
}

// Open opens the named file, satisfying [fs.FS].
func (f *FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	e, ok := f.entries[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	h := &fileHandle{entry: e, fsRef: f}
	if !e.isDir {
		h.reader = bytes.NewReader(e.data)
	}
	return h, nil
}

// ReadDir returns the sorted directory entries for name, satisfying [fs.ReadDirFS].
func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}
	e, ok := f.entries[name]
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	if !e.isDir {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: errors.New("not a directory")}
	}
	children := f.roots[e.path]
	result := make([]fs.DirEntry, len(children))
	for i, child := range children {
		result[i] = &dirEntry{f.entries[child]}
	}
	return result, nil
}

// ReadFile returns a copy of the contents of the named file, satisfying [fs.ReadFileFS].
func (f *FS) ReadFile(name string) ([]byte, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: fs.ErrInvalid}
	}
	e, ok := f.entries[name]
	if !ok {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: fs.ErrNotExist}
	}
	if e.isDir {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: errors.New("is a directory")}
	}
	out := make([]byte, len(e.data))
	copy(out, e.data)
	return out, nil
}

// ---------------------------------------------------------------------------
// fileHandle — returned by Open; implements fs.File, fs.ReadDirFile, io.Seeker
// ---------------------------------------------------------------------------

type fileHandle struct {
	entry  *entry
	reader *bytes.Reader // nil for directories; bytes.Reader supports io.ReadSeeker
	dirPos int
	fsRef  *FS
}

func (h *fileHandle) Stat() (fs.FileInfo, error) { return fileInfo{h.entry}, nil }

func (h *fileHandle) Read(b []byte) (int, error) {
	if h.entry.isDir {
		return 0, &fs.PathError{Op: "read", Path: h.entry.path, Err: errors.New("is a directory")}
	}
	return h.reader.Read(b)
}

// Seek implements io.Seeker. http.FileServer uses this for range requests.
func (h *fileHandle) Seek(offset int64, whence int) (int64, error) {
	if h.entry.isDir {
		return 0, &fs.PathError{Op: "seek", Path: h.entry.path, Err: errors.New("is a directory")}
	}
	return h.reader.Seek(offset, whence)
}

func (h *fileHandle) Close() error { return nil }

// ReadDir satisfies [fs.ReadDirFile], required by [net/http.FileServer].
func (h *fileHandle) ReadDir(n int) ([]fs.DirEntry, error) {
	if !h.entry.isDir {
		return nil, &fs.PathError{Op: "readdir", Path: h.entry.path, Err: errors.New("not a directory")}
	}
	children := h.fsRef.roots[h.entry.path]
	total := len(children)
	start := h.dirPos

	if n <= 0 {
		h.dirPos = total
		result := make([]fs.DirEntry, total-start)
		for i, child := range children[start:] {
			result[i] = &dirEntry{h.fsRef.entries[child]}
		}
		return result, nil
	}

	end := start + n
	if end > total {
		end = total
	}
	h.dirPos = end
	result := make([]fs.DirEntry, end-start)
	for i, child := range children[start:end] {
		result[i] = &dirEntry{h.fsRef.entries[child]}
	}
	if end == total && n > 0 {
		return result, io.EOF
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// fileInfo — implements fs.FileInfo
// ---------------------------------------------------------------------------

type fileInfo struct{ e *entry }

func (fi fileInfo) Name() string       { return fi.e.name }
func (fi fileInfo) Size() int64        { return fi.e.size }
func (fi fileInfo) Mode() fs.FileMode  { return fi.e.mode }
func (fi fileInfo) ModTime() time.Time { return fi.e.modTime }
func (fi fileInfo) IsDir() bool        { return fi.e.isDir }
func (fi fileInfo) Sys() any           { return nil }

// ---------------------------------------------------------------------------
// dirEntry — implements fs.DirEntry
// ---------------------------------------------------------------------------

type dirEntry struct{ e *entry }

func (d *dirEntry) Name() string               { return d.e.name }
func (d *dirEntry) IsDir() bool                { return d.e.isDir }
func (d *dirEntry) Type() fs.FileMode          { return d.e.mode.Type() }
func (d *dirEntry) Info() (fs.FileInfo, error) { return fileInfo{d.e}, nil }
