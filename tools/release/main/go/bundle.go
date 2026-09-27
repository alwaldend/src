package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/release/main/proto/contracts"
)

func copyBundleFile(source, destination string) (result error) {
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close source: %w", err))
		}
	}()
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	if _, err := io.Copy(output, input); err != nil {
		result = fmt.Errorf("copy payload: %w", err)
	}
	if err := output.Close(); err != nil {
		result = errors.Join(result, fmt.Errorf("close destination: %w", err))
	}
	return result
}

func readReleaseVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read workspace status %q: %w", path, err)
	}
	version := ""
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "STABLE_VERSION" {
			continue
		}
		if len(fields) != 2 || version != "" {
			return "", fmt.Errorf("workspace status must contain exactly one STABLE_VERSION value")
		}
		version = fields[1]
	}
	if err := pathComponent(version); err != nil {
		return "", fmt.Errorf("invalid STABLE_VERSION: %w", err)
	}
	return version, nil
}

func (self *Generator) writeBundle(opts *GenerateOpts, release *contracts.Release) error {
	if err := pathComponent(release.Name); err != nil {
		return fmt.Errorf("release version: %w", err)
	}
	if err := pathComponent(strings.TrimPrefix(release.GetProject().GetSubdir(), "projects/")); err != nil {
		return fmt.Errorf("release project: %w", err)
	}
	if len(release.Items) == 0 {
		return fmt.Errorf("release bundle requires at least one file")
	}
	sources := make(map[string]string)
	for _, source := range opts.AddFiles {
		sources[filepath.Base(source)] = source
	}
	for _, item := range release.Items {
		if item.File == nil || sources[item.File.Name] == "" {
			return fmt.Errorf("every bundled file must have an --add_file source")
		}
		if err := pathComponent(item.File.Name); err != nil {
			return fmt.Errorf("bundle filename: %w", err)
		}
	}
	entries, err := os.ReadDir(opts.OutputDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("inspect output directory: %w", err)
	}
	if len(entries) != 0 {
		return fmt.Errorf("output directory must be empty")
	}
	filesDir := filepath.Join(opts.OutputDir, "files")
	if err := os.MkdirAll(filesDir, 0o755); err != nil {
		return fmt.Errorf("create bundle directory: %w", err)
	}
	for _, item := range release.Items {
		source := sources[item.File.Name]
		if err := copyBundleFile(source, filepath.Join(filesDir, item.File.Name)); err != nil {
			return fmt.Errorf("copy bundle file %q: %w", item.File.Name, err)
		}
	}
	if err := self.writeMessage(opts, release, filepath.Join(opts.OutputDir, "release.json")); err != nil {
		return fmt.Errorf("write bundle manifest: %w", err)
	}
	return nil
}
