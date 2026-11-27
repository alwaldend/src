package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/go/pkg/flagcompat"
	"github.com/spf13/cobra"
)

// entry describes one ASAR header node using its JSON field names.
type entry struct {
	// Files maps directory children by archive basename.
	Files map[string]entry `json:"files"`
	// Size is the length in bytes of a packed file.
	Size int64 `json:"size"`
	// Offset is the decimal offset relative to the packed data.
	Offset string `json:"offset"`
	// Link names a link target, which the extractor rejects.
	Link string `json:"link"`
	// Unpacked indicates that the payload resides outside the ASAR file.
	Unpacked bool `json:"unpacked"`
}

// main binds the extraction action's paths and cancels the child process on interrupt.
func main() {
	var drawio, out string
	status := 2
	command := &cobra.Command{
		Use: "extract", Args: cobra.NoArgs, SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if drawio == "" || out == "" {
				return fmt.Errorf("--drawio and --out must be nonempty")
			}
			status = 1
			if err := run(cmd.Context(), drawio, out); err != nil {
				return fmt.Errorf("extract webapp to %q: %w", out, err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&drawio, "drawio", "", "pinned Drawio AppImage executable")
	command.Flags().StringVar(&out, "out", "", "new output directory")
	for _, name := range []string{"drawio", "out"} {
		if err := command.MarkFlagRequired(name); err != nil {
			fmt.Fprintf(os.Stderr, "configure --%s: %v\n", name, err)
			os.Exit(1)
		}
	}
	flagcompat.SetArgs(command, os.Args[1:])
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := command.ExecuteContext(ctx)
	stop()
	if err != nil {
		os.Exit(status)
	}
}

// run unpacks the AppImage in a private workspace and extracts only its web application.
func run(ctx context.Context, drawio, out string) (result error) {
	drawio, err := filepath.Abs(drawio)
	if err != nil {
		return fmt.Errorf("resolve AppImage path: %w", err)
	}
	out, err = filepath.Abs(out)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("create output parent %q: %w", filepath.Dir(out), err)
	}
	// Keep extraction inside the action's output tree, independent of host TMPDIR.
	scratch, err := os.MkdirTemp(filepath.Dir(out), ".drawio-extract-")
	if err != nil {
		return fmt.Errorf("create extraction workspace: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(scratch); err != nil {
			result = errors.Join(result, fmt.Errorf("remove extraction workspace: %w", err))
		}
	}()
	cmd := exec.CommandContext(ctx, drawio, "--appimage-extract", "resources/app.asar")
	cmd.Dir = scratch
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract AppImage: %w", err)
	}
	archive := filepath.Join(scratch, "squashfs-root", "resources", "app.asar")
	for path := archive; path != scratch; path = filepath.Dir(path) {
		info, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("inspect extracted path %q: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in extracted archive path: %s", path)
		}
	}
	f, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("open archive %q: %w", archive, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close archive: %w", err))
		}
	}()
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("inspect archive %q: %w", archive, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("archive is not a regular file")
	}
	if err := os.Mkdir(out, 0o755); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("create output %q: %w", out, err)
		}
		info, statErr := os.Lstat(out)
		if statErr != nil {
			return fmt.Errorf("inspect output %q: %w", out, statErr)
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output is not a real directory")
		}
		directory, openErr := os.Open(out)
		if openErr != nil {
			return fmt.Errorf("open output directory %q: %w", out, openErr)
		}
		entries, readErr := directory.ReadDir(1)
		if readErr == io.EOF {
			readErr = nil
		}
		readErr = errors.Join(readErr, directory.Close())
		if readErr != nil {
			return fmt.Errorf("list output %q: %w", out, readErr)
		}
		if len(entries) != 0 {
			return fmt.Errorf("output directory must be empty")
		}
	}
	if err := extract(ctx, f, info.Size(), out); err != nil {
		failure := fmt.Errorf("extract ASAR %q: %w", archive, err)
		if cleanupErr := os.RemoveAll(out); cleanupErr != nil {
			return errors.Join(failure, fmt.Errorf("remove partial extraction %q: %w", out, cleanupErr))
		}
		return failure
	}
	return nil
}

// extract validates the ASAR tree and streams packed files to their destination.
func extract(ctx context.Context, r io.ReaderAt, size int64, out string) error {
	var prefix [16]byte
	if _, err := r.ReadAt(prefix[:], 0); err != nil {
		return fmt.Errorf("ASAR prefix: %w", err)
	}
	headerSize := int64(binary.LittleEndian.Uint32(prefix[4:8]))
	jsonSize := int64(binary.LittleEndian.Uint32(prefix[12:16]))
	dataStart := 8 + headerSize
	if binary.LittleEndian.Uint32(prefix[:4]) != 4 || headerSize < 8 || jsonSize > 16<<20 || jsonSize > headerSize-8 || dataStart > size {
		return fmt.Errorf("invalid ASAR header bounds")
	}
	data := make([]byte, jsonSize)
	if _, err := r.ReadAt(data, 16); err != nil {
		return fmt.Errorf("read ASAR header: %w", err)
	}
	var root entry
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("ASAR JSON: %w", err)
	}
	for _, name := range []string{"drawio", "src", "main", "webapp"} {
		next, ok := root.Files[name]
		if !ok || next.Files == nil || next.Link != "" || next.Unpacked {
			return fmt.Errorf("missing or invalid ASAR webapp directory: %s", name)
		}
		root = next
	}
	// Entries may reference overlapping byte ranges. Bound total writes independently
	// of the ASAR file size as well as the number and depth of extracted paths.
	const maxWebappBytes int64 = 8 << 30
	var bytesWritten int64
	count := 0
	var walk func(entry, string, int) error
	walk = func(node entry, dest string, depth int) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("extract %q: %w", dest, err)
		}
		count++
		if count > 100000 || depth > 128 {
			return fmt.Errorf("ASAR tree exceeds limits")
		}
		if node.Link != "" || node.Unpacked {
			return fmt.Errorf("unsupported ASAR link or unpacked entry: %s", dest)
		}
		if node.Files != nil {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return fmt.Errorf("create directory %q: %w", dest, err)
			}
			names := make([]string, 0, len(node.Files))
			for name := range node.Files {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00") {
					return fmt.Errorf("unsafe ASAR filename %q", name)
				}
				if err := walk(node.Files[name], filepath.Join(dest, name), depth+1); err != nil {
					return fmt.Errorf("extract child %q: %w", name, err)
				}
			}
			return nil
		}
		offset, err := strconv.ParseInt(node.Offset, 10, 64)
		if err != nil {
			return fmt.Errorf("parse ASAR offset for %q: %w", dest, err)
		}
		if offset < 0 || node.Size < 0 || offset > size-dataStart || node.Size > size-dataStart-offset {
			return fmt.Errorf("invalid ASAR data bounds: %s", dest)
		}
		if node.Size > maxWebappBytes-bytesWritten {
			return fmt.Errorf("ASAR web application exceeds %d bytes", maxWebappBytes)
		}
		f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return fmt.Errorf("create file %q: %w", dest, err)
		}
		_, copyErr := io.CopyN(f, io.NewSectionReader(r, dataStart+offset, node.Size), node.Size)
		closeErr := f.Close()
		if err := errors.Join(copyErr, closeErr); err != nil {
			return fmt.Errorf("write extracted file %q: %w", dest, err)
		}
		bytesWritten += node.Size
		return nil
	}
	if err := walk(root, out, 0); err != nil {
		return fmt.Errorf("unpack webapp: %w", err)
	}
	return nil
}
