package archive

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
)

// Source owns a private collection and borrows the caller's running connection.
// The caller keeps that connection alive and serializes collection operations
// until Close or Publish. Source never starts or stops a worker.
type Source struct {
	// closed prevents reuse after cleanup or a publication attempt.
	closed bool
	// collections controls the native collection lifetime.
	collections ankicol.BackendCollectionServiceClient
	// packages imports and exports Anki's supported collection packages.
	packages import_export.BackendImportExportServiceClient
	// opened records a collection that still requires native closing.
	opened bool
	// packaged selects native package export rather than an offline database copy.
	packaged bool
	// dir is the private directory removed after native work completes.
	dir string
	// database is the private collection path understood only by Anki.
	database string
}

// FileHash computes the SHA-256 digest used by collection replacement guards.
func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open collection file %q: %w", path, err)
	}
	h := sha256.New()
	_, copyErr := io.Copy(h, f)
	if err := errors.Join(copyErr, f.Close()); err != nil {
		return "", fmt.Errorf("hash collection file %q: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// verifyUnchanged rejects replacement after another writer changes the target.
func verifyUnchanged(path, expectedHash string) error {
	current, err := FileHash(path)
	if err != nil {
		return fmt.Errorf("verify publication target: %w", err)
	}
	if current != expectedHash {
		return fmt.Errorf("collection changed before publication; refusing replacement")
	}
	return nil
}

// OpenPackage restores a package through Anki and opens its private collection.
// The supplied connection must be running with no collection open.
// Collection Build uses this constructor to exclude raw database inputs.
func OpenPackage(ctx context.Context, connection grpc.ClientConnInterface, path, scratch string) (*Source, error) {
	p, err := open(ctx, connection, path, scratch, true)
	if err != nil {
		return nil, fmt.Errorf("open collection package %q: %w", path, err)
	}
	return p, nil
}

// Close closes the private collection and removes its scratch directory.
// It leaves the caller's worker running and permits repeated cleanup. Cleanup
// awaits native closing without cancellation before removing collection files.
func (p *Source) Close() error {
	p.closed = true
	if err := p.closeCollection(context.Background()); err != nil {
		return fmt.Errorf("close collection before removing scratch: %w", err)
	}
	p.database = ""
	if p.dir != "" {
		if err := os.RemoveAll(p.dir); err != nil {
			return fmt.Errorf("remove collection scratch: %w", err)
		}
		p.dir = ""
	}
	return nil
}

// closeCollection retires a known open collection before raw publication or cleanup.
func (p *Source) closeCollection(ctx context.Context) error {
	if !p.opened {
		return nil
	}
	if _, err := p.collections.CloseCollection(ctx, &ankicol.CloseCollectionRequest{}); err != nil {
		return fmt.Errorf("close private collection: %w", err)
	}
	p.opened = false
	return nil
}

// Publish consumes this source and replaces a destination with a completed native
// package or offline database. The caller supplies mode and an optional expected
// FileHash. Digest checking and rename are separate; callers serialize writers.
// Cancellation is checked between stages. Started native calls complete before
// temporary outputs are removed; caller-owned workers remain running.
func (p *Source) Publish(ctx context.Context, path string, mode os.FileMode, expectedHash string) (err error) {
	if p.closed || p.database == "" || !p.opened {
		return fmt.Errorf("collection source is closed")
	}
	var temporary string
	defer func() {
		err = errors.Join(err, p.Close())
		if temporary != "" {
			if removeErr := os.Remove(temporary); removeErr != nil && !os.IsNotExist(removeErr) {
				err = errors.Join(err, fmt.Errorf("remove collection staging: %w", removeErr))
			}
		}
	}()
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("publish collection: %w", err)
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".anki-output-")
	if err != nil {
		return fmt.Errorf("create collection staging: %w", err)
	}
	temporary, err = filepath.Abs(f.Name())
	closeErr := f.Close()
	if err := errors.Join(err, closeErr); err != nil {
		return fmt.Errorf("prepare collection staging: %w", err)
	}
	if p.packaged {
		// ExportCollectionPackage takes the open collection even when export fails.
		p.opened = false
		if _, err := p.packages.ExportCollectionPackage(context.WithoutCancel(ctx), &import_export.ExportCollectionPackageRequest{
			OutPath: temporary, IncludeMedia: true,
		}); err != nil {
			return fmt.Errorf("export collection package: %w", err)
		}
	} else {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("publish offline collection: %w", err)
		}
		if err := p.closeCollection(context.WithoutCancel(ctx)); err != nil {
			return fmt.Errorf("prepare offline collection: %w", err)
		}
		if err := copyFile(p.database, temporary); err != nil {
			return fmt.Errorf("stage offline collection: %w", err)
		}
	}
	f, err = os.OpenFile(temporary, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("open completed staging: %w", err)
	}
	modeErr := f.Chmod(mode)
	syncErr := f.Sync()
	if err := errors.Join(modeErr, syncErr, f.Close()); err != nil {
		return fmt.Errorf("finish collection staging: %w", err)
	}
	if err := p.Close(); err != nil {
		return fmt.Errorf("close source before publication: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("publish collection: %w", err)
	}
	if expectedHash != "" {
		if err := verifyUnchanged(path, expectedHash); err != nil {
			return fmt.Errorf("verify collection before replacement: %w", err)
		}
	}
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("publish collection %q: %w", path, err)
	}
	return nil
}
