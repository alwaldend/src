package archive

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

const (
	// collectionDatabaseEntry names the supported compressed database entry in collection archives.
	collectionDatabaseEntry = "collection.anki21b"
	// supportedPackageMetadata is the supported archive metadata payload.
	supportedPackageMetadata = "\x08\x03"
)

// Source owns a private database copy and optional source archive resources.
type Source struct {
	// zip owns the open base archive, or is nil for an offline database.
	zip *zip.ReadCloser
	// dir, database identify the private scratch directory and its editable database file.
	dir, database string
}

// FileHash computes the SHA-256 digest of a collection file.
func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open archive %q: %w", path, err)
	}
	h := sha256.New()
	_, copyErr := io.Copy(h, f)
	closeErr := f.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return "", fmt.Errorf("hash archive %q: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// VerifyUnchanged rejects publication when the input no longer matches its recorded digest.
func VerifyUnchanged(path, expectedHash string) error {
	current, err := FileHash(path)
	if err != nil {
		return fmt.Errorf("verify apply target: %w", err)
	}
	if current != expectedHash {
		return fmt.Errorf("collection changed during apply; refusing replacement")
	}
	return nil
}

// OpenPackage extracts the supported archive database into a private scratch directory.
func OpenPackage(ctx context.Context, path, scratch string) (p *Source, err error) {
	p = &Source{}
	defer func() {
		if err != nil {
			err = errors.Join(err, p.Close())
		}
	}()
	p.zip, err = zip.OpenReader(path)
	if err != nil {
		return p, fmt.Errorf("open package %q: %w", path, err)
	}
	var database *zip.File
	seen := map[string]bool{}
	for _, f := range p.zip.File {
		if seen[f.Name] {
			return p, fmt.Errorf("duplicate archive entry %q", f.Name)
		}
		seen[f.Name] = true
		if f.Name == collectionDatabaseEntry {
			database = f
		}
	}
	if database == nil {
		return p, fmt.Errorf("unsupported package: expected modern %s", collectionDatabaseEntry)
	}
	meta, err := p.zip.Open("meta")
	if err != nil {
		return p, fmt.Errorf("read package metadata: %w", err)
	}
	mb, readErr := io.ReadAll(io.LimitReader(meta, int64(len(supportedPackageMetadata)+1)))
	closeErr := meta.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return p, fmt.Errorf("read package metadata: %w", err)
	}
	if string(mb) != supportedPackageMetadata {
		return p, fmt.Errorf("unsupported package metadata %x", mb)
	}
	p.dir, err = os.MkdirTemp(scratch, "anki-database-")
	if err != nil {
		return p, fmt.Errorf("create database scratch: %w", err)
	}
	p.database = filepath.Join(p.dir, "collection.sqlite")
	r, err := database.Open()
	if err != nil {
		return p, fmt.Errorf("read collection database: %w", err)
	}
	z, err := zstd.NewReader(r, zstd.WithDecoderConcurrency(1))
	if err != nil {
		_ = r.Close()
		return p, fmt.Errorf("decode database: %w", err)
	}
	f, err := os.OpenFile(p.database, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		z.Close()
		_ = r.Close()
		return p, fmt.Errorf("create database: %w", err)
	}
	_, copyErr := io.Copy(f, z)
	z.Close()
	if err := errors.Join(copyErr, f.Close(), r.Close()); err != nil {
		return p, fmt.Errorf("extract database: %w", err)
	}
	return p, nil
}

// Close closes the source archive and removes its private scratch directory.
func (p *Source) Close() error {
	var errs []error
	if p.zip != nil {
		if err := p.zip.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close archive: %w", err))
		}
		p.zip = nil
	}
	if p.dir != "" {
		if err := os.RemoveAll(p.dir); err != nil {
			errs = append(errs, fmt.Errorf("remove database scratch: %w", err))
		}
		p.dir = ""
	}
	return errors.Join(errs...)
}

// write rebuilds an archive and replaces the destination after verifying the input digest.
func (p *Source) write(path string, mode os.FileMode, expectedHash string) (err error) {
	f, err := os.CreateTemp(filepath.Dir(path), ".anki-output-")
	if err != nil {
		return fmt.Errorf("create archive output: %w", err)
	}
	tmp := f.Name()
	defer func() {
		_ = f.Close()
		_ = os.Remove(tmp)
	}()
	w := zip.NewWriter(f)
	for _, entry := range p.zip.File {
		if entry.Name != collectionDatabaseEntry {
			if err := w.Copy(entry); err != nil {
				_ = w.Close()
				return fmt.Errorf("copy archive entry %q: %w", entry.Name, err)
			}
			continue
		}
		s, err := w.CreateHeader(&zip.FileHeader{Name: entry.Name, Method: zip.Store})
		if err != nil {
			_ = w.Close()
			return fmt.Errorf("create compressed database entry: %w", err)
		}
		z, err := zstd.NewWriter(s, zstd.WithEncoderConcurrency(1))
		if err != nil {
			_ = w.Close()
			return fmt.Errorf("encode database: %w", err)
		}
		r, err := os.Open(p.database)
		if err != nil {
			_ = z.Close()
			_ = w.Close()
			return fmt.Errorf("open patched database: %w", err)
		}
		_, copyErr := io.Copy(z, r)
		if err := errors.Join(copyErr, r.Close(), z.Close()); err != nil {
			_ = w.Close()
			return fmt.Errorf("compress patched database: %w", err)
		}
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("finish archive: %w", err)
	}
	if err := f.Chmod(mode); err != nil {
		return fmt.Errorf("set archive permissions: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync archive: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close archive output: %w", err)
	}
	if err := p.zip.Close(); err != nil {
		return fmt.Errorf("close input archive: %w", err)
	}
	p.zip = nil
	if expectedHash != "" {
		if err := VerifyUnchanged(path, expectedHash); err != nil {
			return fmt.Errorf("verify archive before replacement: %w", err)
		}
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("publish archive %q: %w", path, err)
	}
	return nil
}

// DatabasePath identifies the private extracted or copied database.
func (p *Source) DatabasePath() string {
	return p.database
}

// Publish replaces the target only after the caller closes its SQLite connection.
func (p *Source) Publish(path string, mode os.FileMode, expectedHash string) error {
	if p.zip != nil {
		if err := p.write(path, mode, expectedHash); err != nil {
			return fmt.Errorf("publish collection archive: %w", err)
		}
		return nil
	}
	if err := os.Chmod(p.database, mode); err != nil {
		return fmt.Errorf("preserve database permissions: %w", err)
	}
	if expectedHash != "" {
		if err := VerifyUnchanged(path, expectedHash); err != nil {
			return fmt.Errorf("verify database before replacement: %w", err)
		}
	}
	if err := os.Rename(p.database, path); err != nil {
		return fmt.Errorf("publish reconciled database: %w", err)
	}
	return nil
}
