package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	// maxHeaderBytes bounds YAML input before decoding one document.
	maxHeaderBytes = 1 << 20
	// maxReadmes bounds the number of records retained for one output tree.
	maxReadmes = 100_000
	// maxTreeMetadataBytes caps cumulative YAML retained for rendering.
	maxTreeMetadataBytes = 64 << 20
	// maxWalkDepth bounds recursive directory descent from the requested path.
	maxWalkDepth = 256
)

// parseResult holds the ordered tree used by the public JSON and template formats.
type parseResult struct {
	// Dir names the root used for relative paths.
	Dir string `json:"dir"`
	// Readmes lists accepted metadata records in traversal order.
	Readmes []*readme `json:"readmes"`
	// metadataBytes counts accepted YAML bytes across all requested paths.
	metadataBytes int
}

// collect appends an explicit file or discovers matching files in a directory.
func (result *parseResult) collect(ctx context.Context, path string, config parseConfig) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect %q: %w", path, err)
	}
	if !info.IsDir() {
		if err := result.addFile(ctx, path, config); err != nil {
			return fmt.Errorf("parse explicit file: %w", err)
		}
		return nil
	}
	if config.useGit {
		if err := result.collectTracked(ctx, path, config); err != nil {
			return fmt.Errorf("inspect Git index: %w", err)
		}
		return nil
	}
	baseDepth := strings.Count(filepath.Clean(path), string(os.PathSeparator))
	err = filepath.WalkDir(path, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("visit %q: %w", path, walkErr)
		}
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("traverse %q: %w", path, err)
		}
		if excluded(path, config.exclude) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if strings.Count(path, string(os.PathSeparator))-baseDepth > maxWalkDepth {
				return fmt.Errorf("directory depth exceeds %d at %q", maxWalkDepth, path)
			}
			return nil
		}
		if !matchesReadme(path, config) {
			return nil
		}
		if err := result.addFile(ctx, path, config); err != nil {
			return fmt.Errorf("read discovered file: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk %q: %w", path, err)
	}
	return nil
}

// collectTracked streams NUL-delimited paths so quoting and newlines are preserved.
func (result *parseResult) collectTracked(ctx context.Context, dir string, config parseConfig) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "ls-files", "-z", "--", ".")
	cmd.Dir = dir
	output, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open Git output: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start git ls-files in %q: %w", dir, err)
	}
	paths := bufio.NewScanner(output)
	paths.Split(splitPath)
	for paths.Scan() {
		path := filepath.Join(dir, paths.Text())
		if !matchesReadme(path, config) || excluded(path, config.exclude) {
			continue
		}
		if err = result.addFile(ctx, path, config); err != nil {
			break
		}
	}
	if err == nil {
		err = paths.Err()
	}
	if err != nil {
		cancel()
		// Reap Git even if parsing or reading its output failed.
		_ = cmd.Wait()
		return fmt.Errorf("read tracked files in %q: %w", dir, err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("list tracked paths in %q: %w", dir, err)
	}
	return nil
}

// splitPath extracts one terminated pathname from git ls-files -z.
func splitPath(data []byte, atEOF bool) (int, []byte, error) {
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF && len(data) != 0 {
		return 0, nil, fmt.Errorf("unterminated Git pathname")
	}
	return 0, nil, nil
}

// matchesReadme ignores the root README that consumes the generated tree.
func matchesReadme(path string, config parseConfig) bool {
	return filepath.Base(path) == config.name && path != filepath.Join(config.root, config.name)
}

// excluded tests whole path components, including the named subtree root.
func excluded(path string, excludes []string) bool {
	for _, exclude := range excludes {
		rel, err := filepath.Rel(exclude, path)
		if err == nil && filepath.IsLocal(rel) {
			return true
		}
	}
	return false
}

// addFile decodes a README's front matter and appends its serialized path metadata.
func (result *parseResult) addFile(ctx context.Context, path string, config parseConfig) (resultErr error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("read %q: %w", path, err)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open README %q: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close README %q: %w", path, err))
		}
	}()
	lines := bufio.NewScanner(file)
	lines.Buffer(nil, maxHeaderBytes)
	if !lines.Scan() || lines.Text() != "---" {
		if err := lines.Err(); err != nil {
			return fmt.Errorf("read opening delimiter in %q: %w", path, err)
		}
		return nil
	}
	var header bytes.Buffer
	for lines.Scan() {
		if lines.Text() == "---" {
			if len(result.Readmes) >= maxReadmes {
				return fmt.Errorf("README tree exceeds %d records", maxReadmes)
			}
			if header.Len() > maxTreeMetadataBytes-result.metadataBytes {
				return fmt.Errorf("tree metadata exceeds %d bytes", maxTreeMetadataBytes)
			}
			var data map[string]any
			if err := yaml.Unmarshal(header.Bytes(), &data); err != nil {
				return fmt.Errorf("decode metadata in %q: %w", path, err)
			}
			rel, err := filepath.Rel(config.root, path)
			if err != nil {
				return fmt.Errorf("make %q relative to %q: %w", path, config.root, err)
			}
			dir := filepath.Dir(rel)
			result.metadataBytes += header.Len()
			result.Readmes = append(result.Readmes, &readme{
				Path: path, PathRelative: rel, Dir: filepath.Dir(path), DirRelative: dir,
				DirRelativeSplit: strings.Split(dir, string(os.PathSeparator)), DirName: filepath.Base(dir), Data: data,
			})
			return nil
		}
		if header.Len()+len(lines.Bytes())+1 > maxHeaderBytes {
			return fmt.Errorf("metadata in %q exceeds %d bytes", path, maxHeaderBytes)
		}
		header.Write(lines.Bytes())
		header.WriteByte('\n')
	}
	if err := lines.Err(); err != nil {
		return fmt.Errorf("read metadata in %q: %w", path, err)
	}
	return nil
}
