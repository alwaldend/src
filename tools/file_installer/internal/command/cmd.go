package cmd

import (
	"archive/tar"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Execute parses installer arguments and attempts each requested installation in order.
func Execute(ctx context.Context) error {
	var target string
	root := &cobra.Command{
		Use: "file-installer", Short: "File installer",
		SilenceErrors: true, SilenceUsage: true,
		Run: func(*cobra.Command, []string) {},
	}
	install := &cobra.Command{
		Use: "install [files-to-install...]", Short: "Install files", Args: cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, files []string) error {
			var failures []error
			for _, file := range files {
				if err := command.Context().Err(); err != nil {
					failures = append(failures, fmt.Errorf("install remaining files: %w", err))
					break
				}
				if err := installArg(command.Context(), file, target); err != nil {
					failures = append(failures, fmt.Errorf("install %q: %w", file, err))
				}
			}
			return errors.Join(failures...)
		},
	}
	install.Flags().StringVarP(&target, "target", "t", "", "default target for all files")
	root.AddCommand(install)
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("file installer: %w", err)
	}
	return nil
}

// installArg selects directory linking or archive extraction from the input's type.
func installArg(ctx context.Context, source, target string) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("inspect source: %w", err)
	}
	if info.IsDir() {
		if err := installDirectory(ctx, source, target); err != nil {
			return fmt.Errorf("link directory contents: %w", err)
		}
		return nil
	}
	if !info.Mode().IsRegular() || filepath.Ext(source) != ".tar" {
		return fmt.Errorf("expected a directory or regular tar archive, got %q", source)
	}
	if err := installTar(ctx, source, filepath.Clean(target)); err != nil {
		return fmt.Errorf("extract archive: %w", err)
	}
	return nil
}

// installTar streams regular entries beneath the destination and owns both open handles.
func installTar(ctx context.Context, source, target string) (result error) {
	log.Printf("installing tar %s to %s", source, target)
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close archive: %w", err))
		}
	}()
	root, err := os.OpenRoot(target)
	if err != nil {
		return fmt.Errorf("open destination %q: %w", target, err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close destination %q: %w", target, err))
		}
	}()
	reader := tar.NewReader(file)
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("read archive: %w", err)
		}
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read archive header: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if err := extractFile(ctx, root, header, reader); err != nil {
			return fmt.Errorf("extract %q: %w", header.Name, err)
		}
	}
}

// extractFile stages one archive entry before replacing its destination.
func extractFile(ctx context.Context, root *os.Root, header *tar.Header, reader io.Reader) (result error) {
	// Create the temporary file through the rooted handle to confine writes.
	directory := filepath.Dir(header.Name)
	if err := root.MkdirAll(directory, 0o777); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}
	// The private directory protects partial data even when the requested mode
	// makes the completed entry readable by other users.
	staging := filepath.Join(directory, ".file-installer-"+rand.Text())
	if err := root.Mkdir(staging, 0o700); err != nil {
		return fmt.Errorf("create private staging directory: %w", err)
	}
	defer func() {
		if err := root.RemoveAll(staging); err != nil {
			result = errors.Join(result, fmt.Errorf("remove staging directory: %w", err))
		}
	}()
	// A caller may mask owner permissions. The private staging directory must
	// remain traversable even when the final entry is created with mode 000.
	if err := root.Chmod(staging, 0o700); err != nil {
		return fmt.Errorf("set staging directory permissions: %w", err)
	}
	temporary := filepath.Join(staging, "entry")
	tmp, err := root.OpenFile(temporary, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fs.FileMode(header.Mode).Perm())
	if err != nil {
		return fmt.Errorf("open staged entry: %w", err)
	}
	_, copyErr := io.CopyN(tmp, archiveReader{ctx: ctx, source: reader}, header.Size)
	// Creation already applied umask. Replacing an existing regular file
	// retains that file's permissions, as an in-place write would.
	existing, statErr := root.Lstat(header.Name)
	var modeErr error
	switch {
	case statErr == nil && existing.Mode().IsRegular():
		modeErr = tmp.Chmod(existing.Mode().Perm())
	case statErr != nil && !errors.Is(statErr, fs.ErrNotExist):
		modeErr = fmt.Errorf("inspect destination permissions: %w", statErr)
	}
	if err := errors.Join(copyErr, modeErr, tmp.Close()); err != nil {
		return fmt.Errorf("stage archive contents: %w", err)
	}
	if err := root.Rename(temporary, header.Name); err != nil {
		return fmt.Errorf("replace extracted file: %w", err)
	}
	return nil
}

// archiveReader interrupts a streamed archive entry between reads.
type archiveReader struct {
	// ctx carries command cancellation.
	ctx context.Context
	// source supplies data for the current entry.
	source io.Reader
}

// Read preserves EOF for io.CopyN and reports cancellation before each read.
func (r archiveReader) Read(buffer []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, fmt.Errorf("read archive entry: %w", err)
	}
	n, err := r.source.Read(buffer)
	if err != nil && err != io.EOF {
		return n, fmt.Errorf("read archive entry: %w", err)
	}
	return n, err
}

// installDirectory mirrors a source tree with symlinks to its files.
func installDirectory(ctx context.Context, source, target string) (result error) {
	absolute, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("resolve source %q: %w", source, err)
	}
	absolute, err = filepath.EvalSymlinks(absolute)
	if err != nil {
		return fmt.Errorf("follow source %q: %w", source, err)
	}
	target = filepath.Clean(target)
	if err := os.MkdirAll(target, 0o777); err != nil {
		return fmt.Errorf("create destination %q: %w", target, err)
	}
	destination, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve destination %q: %w", target, err)
	}
	destination, err = filepath.EvalSymlinks(destination)
	if err != nil {
		return fmt.Errorf("follow destination %q: %w", target, err)
	}
	if absolute == destination {
		return fmt.Errorf("source and destination are the same directory %q", absolute)
	}
	root, err := os.OpenRoot(destination)
	if err != nil {
		return fmt.Errorf("open destination %q: %w", target, err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close destination %q: %w", target, err))
		}
	}()
	if err := linkDirectory(ctx, root, absolute, ".", destination, 0); err != nil {
		return fmt.Errorf("install directory %q: %w", source, err)
	}
	return nil
}

// linkDirectory processes directory entries in bounded batches and skips a nested destination.
func linkDirectory(ctx context.Context, root *os.Root, source, relative, destination string, depth int) (result error) {
	const maximumDepth = 128
	if depth >= maximumDepth {
		return fmt.Errorf("source directory %q exceeds depth limit %d", source, maximumDepth)
	}
	directory, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open directory %q: %w", source, err)
	}
	defer func() {
		if err := directory.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close directory %q: %w", source, err))
		}
	}()
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("scan %q: %w", source, err)
		}
		entries, err := directory.ReadDir(64)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read directory %q: %w", source, err)
		}
		for _, entry := range entries {
			path := filepath.Join(source, entry.Name())
			name := filepath.Join(relative, entry.Name())
			if entry.IsDir() {
				if path == destination {
					continue
				}
				if err := linkDirectory(ctx, root, path, name, destination, depth+1); err != nil {
					return fmt.Errorf("descend into %q: %w", name, err)
				}
				continue
			}
			if err := root.MkdirAll(relative, 0o777); err != nil {
				return fmt.Errorf("create parent for %q: %w", name, err)
			}
			if err := root.Symlink(path, name); err != nil {
				if errors.Is(err, fs.ErrExist) {
					if existing, readErr := root.Readlink(name); readErr == nil && existing == path {
						continue
					}
				}
				return fmt.Errorf("link %q to %q: %w", name, path, err)
			}
		}
	}
}
