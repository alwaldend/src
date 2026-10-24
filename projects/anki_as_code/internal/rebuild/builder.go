package rebuild

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
)

// Builder publishes a copied package only after its injected apply operation succeeds.
type Builder struct {
	// filesystem owns source inspection, staging, copying, publication and cleanup.
	filesystem filesystem.FS
	// apply updates the private archive and returns after closing its native collection.
	apply func(context.Context, string) error
}

// New uses the supplied non-nil apply operation without owning its native services.
// storage supplies staging and publication; its zero value selects the OS filesystem.
func New(apply func(context.Context, string) error, storage filesystem.FS) *Builder {
	return &Builder{apply: apply, filesystem: storage}
}

// Build applies changes to a private copy and atomically replaces the destination.
// Paths use ordinary filesystem notation; the output parent must exist.
// Destinations aliasing the source are rejected. Failures before publication
// leave both paths intact. Cancellation retains the context error and its cause.
func (builder *Builder) Build(ctx context.Context, base, output string) (err error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start build: %w", errors.Join(err, context.Cause(ctx)))
	}
	source, err := builder.filesystem.ResolvePath(base, false)
	if err != nil {
		return fmt.Errorf("resolve build source aliases: %w", err)
	}
	if !strings.EqualFold(filepath.Ext(source), ".colpkg") {
		return fmt.Errorf("build requires a collection package")
	}
	output, err = resolveDestination(builder.filesystem, output)
	if err != nil {
		return fmt.Errorf("resolve build destination: %w", err)
	}
	original, err := builder.filesystem.Stat(source)
	if err != nil {
		return fmt.Errorf("inspect build source: %w", err)
	}
	if !original.Mode().IsRegular() {
		return fmt.Errorf("build source must be a regular collection package")
	}
	target, err := builder.filesystem.Stat(output)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect build destination: %w", err)
	}
	if source == output || target != nil && os.SameFile(original, target) {
		return fmt.Errorf("output must differ from the base archive")
	}
	staging, err := builder.filesystem.TempDir(filepath.Dir(output), ".anki-build-")
	if err != nil {
		return fmt.Errorf("create build staging directory: %w", err)
	}
	defer func() {
		if cleanupErr := builder.filesystem.RemoveAll(staging); cleanupErr != nil {
			err = errors.Join(err, fmt.Errorf("remove build staging directory: %w", cleanupErr))
		}
	}()
	archive := filepath.Join(staging, "collection.colpkg")
	if err := copyArchive(builder.filesystem, source, archive); err != nil {
		return fmt.Errorf("copy build source: %w", err)
	}
	if err := builder.apply(ctx, archive); err != nil {
		return fmt.Errorf("apply staged build: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish build: %w", errors.Join(err, context.Cause(ctx)))
	}
	if err := builder.filesystem.Rename(archive, output); err != nil {
		return fmt.Errorf("publish build: %w", err)
	}
	return nil
}

// resolveDestination follows an existing output or resolves its required parent directory.
func resolveDestination(storage filesystem.FS, output string) (string, error) {
	resolved, err := storage.ResolvePath(output, false)
	if err == nil {
		return resolved, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("resolve output aliases: %w", err)
	}
	// Split retains unresolved parent components and the final dangling-link entry.
	directory, name := filepath.Split(output)
	parent, err := storage.ResolvePath(directory, false)
	if err != nil {
		return "", fmt.Errorf("resolve output parent: %w", err)
	}
	return filepath.Join(parent, name), nil
}

// copyArchive streams a source package into a new staging file and reports close errors.
func copyArchive(storage filesystem.FS, source, destination string) (err error) {
	input, err := storage.Open(source)
	if err != nil {
		return fmt.Errorf("open source archive: %w", err)
	}
	defer func() {
		if closeErr := input.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close source archive: %w", closeErr))
		}
	}()
	output, err := storage.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create staging archive: %w", err)
	}
	defer func() {
		if closeErr := output.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close staging archive: %w", closeErr))
		}
	}()
	if _, err := io.Copy(output, input); err != nil {
		return fmt.Errorf("copy archive bytes: %w", err)
	}
	return nil
}
