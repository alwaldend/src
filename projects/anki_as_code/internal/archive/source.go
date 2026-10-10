package archive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// isDatabase recognizes an offline SQLite collection by its file header.
func isDatabase(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open collection %q: %w", path, err)
	}
	b := make([]byte, 16)
	_, readErr := io.ReadFull(f, b)
	closeErr := f.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return false, fmt.Errorf("read collection header: %w", err)
	}
	return string(b) == "SQLite format 3\x00", nil
}

// OpenSource opens an archive or offline database as a private editable source.
func OpenSource(ctx context.Context, path, scratch string) (*Source, error) {
	isDB, err := isDatabase(path)
	if err != nil {
		return nil, fmt.Errorf("identify collection: %w", err)
	}
	if !isDB {
		p, err := OpenPackage(ctx, path, scratch)
		if err != nil {
			return nil, fmt.Errorf("load archive: %w", err)
		}
		return p, nil
	}
	for _, suffix := range []string{"-wal", "-journal"} {
		if _, err := os.Stat(path + suffix); err == nil {
			return nil, fmt.Errorf("collection has %s; close Anki and checkpoint first", suffix)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("check database journal: %w", err)
		}
	}
	dir, err := os.MkdirTemp(scratch, "anki-database-")
	if err != nil {
		return nil, fmt.Errorf("create collection scratch: %w", err)
	}
	p := &Source{dir: dir, database: filepath.Join(dir, "collection.sqlite")}
	if err := CopyFile(path, p.database); err != nil {
		_ = p.Close()
		return nil, fmt.Errorf("copy offline collection: %w", err)
	}
	return p, nil
}

// CopyFile copies file content and reports read, write, and close failures.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open copy input: %w", err)
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		_ = in.Close()
		return fmt.Errorf("create copy output: %w", err)
	}
	_, copyErr := io.Copy(out, in)
	if err := errors.Join(copyErr, in.Close(), out.Close()); err != nil {
		return fmt.Errorf("copy collection: %w", err)
	}
	return nil
}
