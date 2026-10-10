package collection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/repository"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// Export exports collection resources into editable TOML declarations.
func Export(ctx context.Context, input, output string) (err error) {
	input, err = textstore.ExpandPath(input)
	if err != nil {
		return fmt.Errorf("resolve input path: %w", err)
	}
	output, err = textstore.ExpandPath(output)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	p, store, err := openCollection(ctx, input, "", false)
	if err != nil {
		return fmt.Errorf("export collection: %w", err)
	}
	defer func() {
		err = errors.Join(err, store.Close(), p.Close())
	}()
	s, err := store.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("read export state: %w", err)
	}
	if err := textstore.Export(s, output); err != nil {
		return fmt.Errorf("export collection text: %w", err)
	}
	return nil
}

// Build rebuilds a collection archive from its base and desired declarations.
func Build(ctx context.Context, base, dir, output string) (err error) {
	base, err = textstore.ExpandPath(base)
	if err != nil {
		return fmt.Errorf("resolve base path: %w", err)
	}
	dir, err = textstore.ExpandPath(dir)
	if err != nil {
		return fmt.Errorf("resolve dir path: %w", err)
	}
	output, err = textstore.ExpandPath(output)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	a, err := filepath.Abs(base)
	if err != nil {
		return fmt.Errorf("resolve base path: %w", err)
	}
	b, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	if a == b {
		return fmt.Errorf("output must differ from base archive")
	}
	if old, e := os.Stat(base); e == nil {
		if target, e := os.Stat(output); e == nil && os.SameFile(old, target) {
			return fmt.Errorf("output aliases base archive")
		}
	}
	d, err := textstore.Load(dir)
	if err != nil {
		return fmt.Errorf("read build state: %w", err)
	}
	p, store, err := openCollection(ctx, base, "", true)
	if err != nil {
		return fmt.Errorf("read base archive: %w", err)
	}
	defer func() {
		err = errors.Join(err, store.Close(), p.Close())
	}()
	if err := store.Reconcile(ctx, d); err != nil {
		return fmt.Errorf("build desired collection: %w", err)
	}
	if err := store.Close(); err != nil {
		return fmt.Errorf("close rebuilt database: %w", err)
	}
	if err := p.Publish(output, 0o600, ""); err != nil {
		return fmt.Errorf("write rebuilt archive: %w", err)
	}
	if err := textstore.PersistDeckIDs(d.Decks); err != nil {
		return fmt.Errorf("archive built but deck identity recording failed: %w", err)
	}
	return nil
}

// Apply reconciles a collection in place after validating its desired declarations.
func Apply(ctx context.Context, input, dir string) (err error) {
	input, err = textstore.ExpandPath(input)
	if err != nil {
		return fmt.Errorf("resolve input path: %w", err)
	}
	dir, err = textstore.ExpandPath(dir)
	if err != nil {
		return fmt.Errorf("resolve dir path: %w", err)
	}
	info, err := os.Lstat(input)
	if err != nil {
		return fmt.Errorf("inspect apply target: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("apply target must be a regular file")
	}
	d, err := textstore.Load(dir)
	if err != nil {
		return fmt.Errorf("load apply state: %w", err)
	}
	expectedHash, err := archive.FileHash(input)
	if err != nil {
		return fmt.Errorf("identify apply target: %w", err)
	}
	p, store, err := openCollection(ctx, input, filepath.Dir(input), false)
	if err != nil {
		return fmt.Errorf("read apply target: %w", err)
	}
	defer func() {
		err = errors.Join(err, store.Close(), p.Close())
	}()
	if err := store.Reconcile(ctx, d); err != nil {
		return fmt.Errorf("apply desired state: %w", err)
	}
	if err := store.Close(); err != nil {
		return fmt.Errorf("close reconciled database: %w", err)
	}
	if err := p.Publish(input, info.Mode().Perm(), expectedHash); err != nil {
		return fmt.Errorf("publish reconciled collection: %w", err)
	}
	if err := textstore.PersistDeckIDs(d.Decks); err != nil {
		return fmt.Errorf("collection applied but deck identity recording failed: %w", err)
	}
	return nil
}

// openCollection joins independent file-resource and database lifetimes.
func openCollection(ctx context.Context, input, scratch string, packageOnly bool) (*archive.Source, *repository.Repository, error) {
	var source *archive.Source
	var err error
	if packageOnly {
		source, err = archive.OpenPackage(ctx, input, scratch)
	} else {
		source, err = archive.OpenSource(ctx, input, scratch)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("open collection resource: %w", err)
	}
	store, err := repository.Open(source.DatabasePath())
	if err != nil {
		return nil, nil, errors.Join(fmt.Errorf("open collection database: %w", err), source.Close())
	}
	return source, store, nil
}
