package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// main packages all consolidated rulesets from the repository revision at HEAD.
func main() {
	var destination string
	command := &cobra.Command{
		Use:           "release_archives",
		Short:         "Package the root-owned Bazel rulesets",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, _ []string) error {
			root := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
			if root == "" || !filepath.IsAbs(destination) {
				return errors.New("run with Bazel and supply an absolute output directory")
			}
			if err := packageRepositories(command.Context(), root, destination); err != nil {
				return fmt.Errorf("package ruleset archives: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&destination, "output-directory", "", "Directory for source archives")
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// packageRepositories builds archives from project-owned release manifests.
func packageRepositories(ctx context.Context, root, destination string) error {
	revision, err := gitOutput(ctx, root, "rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("resolve archive revision: %w", err)
	}
	commit := strings.TrimSpace(string(revision))
	layout, err := readArchiveLayout(ctx, root, commit)
	if err != nil {
		return fmt.Errorf("read release inventory: %w", err)
	}
	if len(layout.sources) == 0 {
		return errors.New("no release manifests use the archive builder")
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return fmt.Errorf("create archive output directory: %w", err)
	}
	slices.Sort(layout.sources)
	for _, path := range layout.sources {
		name := layout.modules[path]
		rewrite := archiveLabels(path, layout.modules)
		if err := packageRepository(ctx, root, commit, path, filepath.Join(destination, name+".tar.gz"), rewrite); err != nil {
			return fmt.Errorf("archive %s: %w", name, err)
		}
	}
	return nil
}

// packageRepository writes an archive to a temporary file before publishing it locally.
func packageRepository(ctx context.Context, root, commit, path, destination string, rewrite *strings.Replacer) (resultErr error) {
	file, err := os.CreateTemp(filepath.Dir(destination), ".ruleset-*.tar.gz")
	if err != nil {
		return fmt.Errorf("create temporary archive: %w", err)
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil && !os.IsNotExist(err) {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove temporary archive: %w", err))
		}
	}()
	compressed := gzip.NewWriter(file)
	archive := tar.NewWriter(compressed)
	writeErr := writeSources(ctx, root, commit, path, archive, rewrite)
	if err := errors.Join(writeErr, archive.Close(), compressed.Close(), file.Close()); err != nil {
		return fmt.Errorf("finish source archive: %w", err)
	}
	if err := os.Rename(file.Name(), destination); err != nil {
		return fmt.Errorf("move completed archive: %w", err)
	}
	return nil
}

// writeSources packages Git-tracked sources and installs the standalone module manifest.
func writeSources(ctx context.Context, root, commit, path string, output *tar.Writer, rewrite *strings.Replacer) (resultErr error) {
	ctx, stop := context.WithCancel(ctx)
	defer stop()
	command := exec.CommandContext(ctx, "git", "archive", "--format=tar", commit+":"+filepath.ToSlash(path))
	command.Dir = root
	command.Stderr = os.Stderr
	source, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open Git archive stream: %w", err)
	}
	if err := command.Start(); err != nil {
		return fmt.Errorf("start Git archive: %w", err)
	}
	defer func() {
		if resultErr != nil {
			stop()
		}
		if err := command.Wait(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("read Git archive: %w", err))
		}
	}()
	input := tar.NewReader(source)
	foundManifest := false
	for {
		header, err := input.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read Git archive entry: %w", err)
		}
		var contents io.Reader = input
		switch {
		case header.Name == ".bazelignore":
			continue
		case header.Name == "release.MODULE.bazel":
			header.Name = "MODULE.bazel"
			foundManifest = true
		case header.Typeflag == tar.TypeReg && containsBazelLabels(header.Name):
			source, err := io.ReadAll(input)
			if err != nil {
				return fmt.Errorf("read source %q: %w", header.Name, err)
			}
			updated := rewrite.Replace(string(source))
			if header.Name == "openspec/BUILD.bazel" || header.Name == "openspec/BUILD" {
				updated = relocateOpenSpecPackage(updated)
			}
			header.Size = int64(len(updated))
			contents = strings.NewReader(updated)
		}
		header.ModTime = time.Unix(0, 0)
		if err := output.WriteHeader(header); err != nil {
			return fmt.Errorf("write archive header %q: %w", header.Name, err)
		}
		if _, err := io.Copy(output, contents); err != nil {
			return fmt.Errorf("write archive entry %q: %w", header.Name, err)
		}
	}
	if _, err := io.Copy(io.Discard, source); err != nil {
		return fmt.Errorf("drain Git archive stream: %w", err)
	}
	if !foundManifest {
		return fmt.Errorf("missing release.MODULE.bazel in %q", path)
	}
	for _, entry := range [][2]string{
		{".bazeliskrc", ".bazeliskrc"},
		{".bazelignore", ".bazelignore"},
		{"internal/license/LICENSE.txt", "LICENSE.txt"},
	} {
		data, err := gitOutput(ctx, root, "show", commit+":"+entry[0])
		if err != nil {
			return fmt.Errorf("read release metadata %q: %w", entry[0], err)
		}
		if err := output.WriteHeader(&tar.Header{Name: entry[1], Mode: 0o644, Size: int64(len(data))}); err != nil {
			return fmt.Errorf("write metadata header %q: %w", entry[1], err)
		}
		if _, err := output.Write(data); err != nil {
			return fmt.Errorf("write release metadata %q: %w", entry[1], err)
		}
	}
	return nil
}

// gitOutput reads committed release inputs from the source worktree.
func gitOutput(ctx context.Context, root string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = root
	command.Stderr = os.Stderr
	return command.Output()
}
