package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strings"
	"unicode"
)

func extractWebsite(sourcePath, destination string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open website archive %q: %w", sourcePath, err)
	}
	defer source.Close()
	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("inspect website archive %q: %w", sourcePath, err)
	}
	root, err := os.OpenRoot(destination)
	if err != nil {
		return fmt.Errorf("open extraction root %q: %w", destination, err)
	}
	defer root.Close()
	if err := extractArchive(root, source, sourcePath, info.Size()); err != nil {
		return fmt.Errorf("unpack website archive %q: %w", sourcePath, err)
	}
	index, err := root.Lstat("index.html")
	if err != nil {
		return fmt.Errorf("find root website index.html: %w", err)
	}
	if !index.Mode().IsRegular() {
		return fmt.Errorf("website archive needs index.html at its root")
	}
	return nil
}

func extractArchive(root *os.Root, source *os.File, sourceName string, size int64) error {
	seen := map[string]bool{}
	entry := func(name string, directory bool, size int64, content io.Reader) error {
		if strings.HasPrefix(name, "/") || strings.Contains(name, "\\") || strings.IndexFunc(name, unicode.IsControl) >= 0 {
			return fmt.Errorf("unsafe archive path %q", name)
		}
		for _, part := range strings.Split(name, "/") {
			if part == ".." {
				return fmt.Errorf("archive traversal %q", name)
			}
		}
		name = path.Clean(name)
		if name == "." && directory {
			return nil
		}
		if name == "." || seen[name] {
			return fmt.Errorf("invalid or duplicate archive path %q", name)
		}
		seen[name] = true
		if len(seen) > 100000 {
			return fmt.Errorf("site exceeds 100000 archive entries")
		}
		if directory {
			if err := root.MkdirAll(name, 0o755); err != nil {
				return fmt.Errorf("create archive directory %q: %w", name, err)
			}
			return nil
		}
		if parent := path.Dir(name); parent != "." {
			if err := root.MkdirAll(parent, 0o755); err != nil {
				return fmt.Errorf("create archive parent %q: %w", parent, err)
			}
		}
		file, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return fmt.Errorf("create archive file %q: %w", name, err)
		}
		_, copyErr := io.CopyN(file, content, size)
		if copyErr == nil {
			// Read the end too: zip checks its CRC at EOF.
			var extra [1]byte
			n, err := content.Read(extra[:])
			if err != nil && err != io.EOF {
				copyErr = fmt.Errorf("verify archive entry checksum: %w", err)
			} else if n != 0 {
				copyErr = fmt.Errorf("archive entry length differs")
			}
		}
		if copyErr == nil {
			copyErr = file.Sync()
		}
		closeErr := file.Close()
		if copyErr != nil {
			return fmt.Errorf("write archive file %q: %w", name, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close archive output %q: %w", name, closeErr)
		}
		return nil
	}
	if strings.HasSuffix(sourceName, ".zip") {
		archive, err := zip.NewReader(source, size)
		if err != nil {
			return fmt.Errorf("read zip archive: %w", err)
		}
		for _, file := range archive.File {
			if !file.Mode().IsRegular() && !file.Mode().IsDir() {
				return fmt.Errorf("unsupported zip entry %q", file.Name)
			}
			if file.UncompressedSize64 > math.MaxInt64 {
				return fmt.Errorf("archive entry too large")
			}
			content, err := file.Open()
			if err != nil {
				return fmt.Errorf("open zip entry %q: %w", file.Name, err)
			}
			err = entry(file.Name, file.Mode().IsDir(), int64(file.UncompressedSize64), content)
			closeErr := content.Close()
			if err != nil {
				return fmt.Errorf("extract zip entry %q: %w", file.Name, err)
			}
			if closeErr != nil {
				return fmt.Errorf("close zip entry %q: %w", file.Name, closeErr)
			}
		}
		return nil
	}
	gz, err := gzip.NewReader(source)
	if err != nil {
		return fmt.Errorf("open gzip stream: %w", err)
	}
	defer gz.Close()
	archive := tar.NewReader(gz)
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar header: %w", err)
		}
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeDir {
			return fmt.Errorf("unsupported tar entry %q", header.Name)
		}
		if err := entry(header.Name, header.Typeflag == tar.TypeDir, header.Size, archive); err != nil {
			return fmt.Errorf("extract tar entry %q: %w", header.Name, err)
		}
	}
	_, err = io.Copy(io.Discard, gz)
	if err != nil {
		return fmt.Errorf("verify gzip trailer: %w", err)
	}
	return nil
}
