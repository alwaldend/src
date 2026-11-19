// Command extract obtains the web application from the pinned Drawio AppImage.
package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type entry struct {
	Files    map[string]entry `json:"files"`
	Size     int64            `json:"size"`
	Offset   string           `json:"offset"`
	Link     string           `json:"link"`
	Unpacked bool             `json:"unpacked"`
}

func main() {
	drawio := flag.String("drawio", "", "pinned Drawio AppImage executable")
	out := flag.String("out", "", "new output directory")
	flag.Parse()
	if *drawio == "" || *out == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "required: --drawio EXECUTABLE --out DIRECTORY")
		os.Exit(2)
	}
	if err := run(*drawio, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(drawio, out string) error {
	drawio, err := filepath.Abs(drawio)
	if err != nil {
		return err
	}
	out, err = filepath.Abs(out)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	// Keep extraction inside the action's output tree, independent of host TMPDIR.
	scratch, err := os.MkdirTemp(filepath.Dir(out), ".drawio-extract-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratch)
	cmd := exec.Command(drawio, "--appimage-extract", "resources/app.asar")
	cmd.Dir = scratch
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract AppImage: %w", err)
	}
	archive := filepath.Join(scratch, "squashfs-root", "resources", "app.asar")
	for path := archive; path != scratch; path = filepath.Dir(path) {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in extracted archive path: %s", path)
		}
	}
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("archive is not a regular file")
	}
	if err := os.Mkdir(out, 0o755); err != nil {
		if !os.IsExist(err) {
			return err
		}
		info, statErr := os.Lstat(out)
		if statErr != nil {
			return statErr
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output is not a real directory")
		}
		entries, readErr := os.ReadDir(out)
		if readErr != nil {
			return readErr
		}
		if len(entries) != 0 {
			return fmt.Errorf("output directory must be empty")
		}
	}
	if err := extract(f, info.Size(), out); err != nil {
		os.RemoveAll(out)
		return err
	}
	return nil
}

func extract(r io.ReaderAt, size int64, out string) error {
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
		return err
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
	count := 0
	var walk func(entry, string, int) error
	walk = func(node entry, dest string, depth int) error {
		count++
		if count > 100000 || depth > 128 {
			return fmt.Errorf("ASAR tree exceeds limits")
		}
		if node.Link != "" || node.Unpacked {
			return fmt.Errorf("unsupported ASAR link or unpacked entry: %s", dest)
		}
		if node.Files != nil {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return err
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
					return err
				}
			}
			return nil
		}
		offset, err := strconv.ParseInt(node.Offset, 10, 64)
		if err != nil || offset < 0 || node.Size < 0 || offset > size-dataStart || node.Size > size-dataStart-offset {
			return fmt.Errorf("invalid ASAR data bounds: %s", dest)
		}
		f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return err
		}
		_, copyErr := io.CopyN(f, io.NewSectionReader(r, dataStart+offset, node.Size), node.Size)
		closeErr := f.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	}
	return walk(root, out, 0)
}
