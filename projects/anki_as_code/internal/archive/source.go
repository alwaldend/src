package archive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
)

// isDatabase recognizes an offline collection by its opaque SQLite file header.
func isDatabase(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open collection %q: %w", path, err)
	}
	b := make([]byte, 16)
	_, readErr := io.ReadFull(f, b)
	if err := errors.Join(readErr, f.Close()); err != nil {
		return false, fmt.Errorf("read collection header: %w", err)
	}
	return string(b) == "SQLite format 3\x00", nil
}

// OpenSource opens an archive or offline database as a private Anki collection.
// It borrows an existing connection with no open collection. Cancellation is
// checked between stages; started native calls finish before scratch cleanup.
func OpenSource(ctx context.Context, connection grpc.ClientConnInterface, path, scratch string) (*Source, error) {
	p, err := open(ctx, connection, path, scratch, false)
	if err != nil {
		return nil, fmt.Errorf("open collection source %q: %w", path, err)
	}
	return p, nil
}

// open establishes private files and a collection on the borrowed connection.
func open(ctx context.Context, connection grpc.ClientConnInterface, path, scratch string, packageOnly bool) (p *Source, err error) {
	p = &Source{}
	defer func() {
		if err != nil {
			err = errors.Join(err, p.Close())
		}
	}()
	if connection == nil {
		return p, fmt.Errorf("running Anki connection is required")
	}
	if err := ctx.Err(); err != nil {
		return p, fmt.Errorf("open private collection: %w", err)
	}
	input, err := filepath.Abs(path)
	if err != nil {
		return p, fmt.Errorf("resolve collection path: %w", err)
	}
	raw, err := isDatabase(input)
	if err != nil {
		return p, fmt.Errorf("identify collection: %w", err)
	}
	if packageOnly && raw {
		return p, fmt.Errorf("expected collection package, got offline database")
	}
	p.packaged = !raw
	if raw {
		for _, suffix := range []string{"-wal", "-journal"} {
			if _, err := os.Stat(input + suffix); err == nil {
				return p, fmt.Errorf("collection has %s; close Anki and checkpoint first", suffix)
			} else if !os.IsNotExist(err) {
				return p, fmt.Errorf("check database journal: %w", err)
			}
		}
	}
	p.dir, err = os.MkdirTemp(scratch, "anki-collection-")
	if err != nil {
		return p, fmt.Errorf("create collection scratch: %w", err)
	}
	p.database, err = filepath.Abs(filepath.Join(p.dir, "collection.anki2"))
	if err != nil {
		return p, fmt.Errorf("resolve private collection: %w", err)
	}
	if raw {
		if err := copyFile(input, p.database); err != nil {
			return p, fmt.Errorf("copy offline collection: %w", err)
		}
	}
	p.collections = ankicol.NewBackendCollectionServiceClient(connection)
	p.packages = import_export.NewBackendImportExportServiceClient(connection)
	nativeCtx := context.WithoutCancel(ctx)
	if p.packaged {
		if _, err := p.packages.ImportCollectionPackage(nativeCtx, &import_export.ImportCollectionPackageRequest{
			ColPath: p.database, BackupPath: input,
			MediaFolder: p.database + ".media", MediaDb: p.database + ".media.db",
		}); err != nil {
			return p, fmt.Errorf("import collection package: %w", err)
		}
	}
	if err := ctx.Err(); err != nil {
		return p, fmt.Errorf("open private collection: %w", err)
	}
	if _, err := p.collections.OpenCollection(nativeCtx, &ankicol.OpenCollectionRequest{
		CollectionPath:  p.database,
		MediaFolderPath: p.database + ".media", MediaDbPath: p.database + ".media.db",
	}); err != nil {
		return p, fmt.Errorf("open private collection: %w", err)
	}
	p.opened = true
	if err := ctx.Err(); err != nil {
		return p, fmt.Errorf("open private collection: %w", err)
	}
	return p, nil
}

// copyFile copies opaque bytes and reports read, write, and close failures.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open copy input: %w", err)
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return errors.Join(fmt.Errorf("open copy output: %w", err), in.Close())
	}
	_, copyErr := io.Copy(out, in)
	if err := errors.Join(copyErr, in.Close(), out.Close()); err != nil {
		return fmt.Errorf("copy collection: %w", err)
	}
	return nil
}
