package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// main reports archive extraction failures at the test-tool boundary.
func main() {
	command := &cobra.Command{
		Use:           "extract <archive> <directory>",
		Args:          cobra.ExactArgs(2),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := extract(args[0], args[1]); err != nil {
				return fmt.Errorf("extract browser distribution: %w", err)
			}
			return nil
		},
	}
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// extract expands regular files into a caller-owned directory without accepting
// archive paths that escape it or entries that create filesystem aliases.
func extract(archive, directory string) (result error) {
	input, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("open archive %q: %w", archive, err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close archive %q: %w", archive, err))
		}
	}()
	reader := tar.NewReader(input)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read archive header: %w", err)
		}
		if !filepath.IsLocal(header.Name) {
			return fmt.Errorf("archive entry %q is not a local path", header.Name)
		}
		filename := filepath.Join(directory, filepath.FromSlash(header.Name))
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filename, 0o755); err != nil {
				return fmt.Errorf("create archive directory %q: %w", header.Name, err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := writeEntry(filename, reader); err != nil {
				return fmt.Errorf("extract archive entry %q: %w", header.Name, err)
			}
		default:
			return fmt.Errorf("archive entry %q has unsupported type %d", header.Name, header.Typeflag)
		}
	}
}

// writeEntry writes one file and retains both copy and close failures.
func writeEntry(filename string, contents io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}
	output, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	_, copyErr := io.Copy(output, contents)
	closeErr := output.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}
