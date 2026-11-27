package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"git.alwaldend.com/alwaldend/src/tools/go/pkg/flagcompat"
)

// main checks or updates tracked module files using the configured Git executable.
func main() {
	var version, git, marker string
	var checkOnly bool
	command := &cobra.Command{
		Use: "go_mod", Args: cobra.NoArgs, SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version == "" {
				return fmt.Errorf("--version requires a nonempty value")
			}
			root, err := rootDir(marker)
			if err != nil {
				return fmt.Errorf("locate module workspace: %w", err)
			}
			gitPath, err := runfiles.Rlocation(git)
			if err != nil {
				return fmt.Errorf("locate Git runfile %q: %w", git, err)
			}
			if checkOnly {
				if err := check(cmd.Context(), root, gitPath, version); err != nil {
					return fmt.Errorf("check module versions: %w", err)
				}
			} else if err := update(cmd.Context(), root, gitPath, version); err != nil {
				return fmt.Errorf("update module versions: %w", err)
			}
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&version, "version", "", "Go version to set in go.mod files")
	flags.StringVar(&git, "git", "", "Git executable runfile")
	flags.StringVar(&marker, "workspace-marker", "", "Workspace module runfile")
	flags.BoolVar(&checkOnly, "check", false, "Check without writing")
	for _, name := range []string{"version", "git", "workspace-marker"} {
		if err := command.MarkFlagRequired(name); err != nil {
			fmt.Fprintf(os.Stderr, "configure required flag %q: %v\n", name, err)
			os.Exit(1)
		}
	}
	flagcompat.SetArgs(command, os.Args[1:])
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := command.ExecuteContext(ctx)
	stop()
	if err != nil {
		os.Exit(1)
	}
}

// rootDir resolves the workspace marker to its source directory.
func rootDir(marker string) (string, error) {
	path, err := runfiles.Rlocation(marker)
	if err != nil {
		return "", fmt.Errorf("locate marker %q: %w", marker, err)
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve marker %q: %w", marker, err)
	}
	return filepath.Dir(path), nil
}

// pathList holds the NUL-delimited module names emitted by Git.
type pathList []byte

// Write bounds the storage used to collect the workspace module list.
func (list *pathList) Write(chunk []byte) (int, error) {
	const maxPathBytes = 4 << 20
	if len(chunk) > maxPathBytes-len(*list) {
		return 0, fmt.Errorf("tracked module paths exceed %d bytes", maxPathBytes)
	}
	*list = append(*list, chunk...)
	return len(chunk), nil
}

// discover lists tracked module files using NUL-delimited Git output.
func discover(ctx context.Context, root, git string) ([]string, error) {
	cmd := exec.CommandContext(ctx, git, "--git-dir", filepath.Join(root, ".git"), "-C", root, "ls-files", "-z", "go.mod", "**/go.mod")
	var output pathList
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("discover go.mod files: %w", err)
	}
	files := make([]string, 0, bytes.Count(output, []byte{0}))
	for name := range strings.SplitSeq(string(output), "\x00") {
		if name != "" {
			files = append(files, name)
		}
	}
	return files, nil
}

// readModule reads one module file within the parser's input budget.
func readModule(path string) (data []byte, result error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close %q: %w", path, err))
		}
	}()
	const maxModuleBytes = 16 << 20
	data, err = io.ReadAll(io.LimitReader(file, maxModuleBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	if len(data) > maxModuleBytes {
		return nil, fmt.Errorf("module %q exceeds %d bytes", path, maxModuleBytes)
	}
	return data, nil
}

// parseVersion reads the Go directive and distinguishes its absence.
func parseVersion(path string, data []byte) (string, bool, error) {
	file, err := modfile.Parse(path, data, nil)
	if err != nil {
		return "", false, fmt.Errorf("parse %q: %w", path, err)
	}
	if file.Go == nil {
		return "", false, nil
	}
	return file.Go.Version, true, nil
}

// check reports every tracked module that lacks the requested version.
func check(ctx context.Context, root, git, version string) error {
	files, err := discover(ctx, root, git)
	if err != nil {
		return fmt.Errorf("list modules in %q: %w", root, err)
	}
	var mismatches []string
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("process module %q: %w", file, err)
		}
		data, err := readModule(filepath.Join(root, filepath.FromSlash(file)))
		if err != nil {
			return fmt.Errorf("read module %q: %w", file, err)
		}
		current, hasDirective, err := parseVersion(file, data)
		if err != nil {
			return fmt.Errorf("inspect module version: %w", err)
		}
		if !hasDirective || current != version {
			mismatches = append(mismatches, fmt.Sprintf("%s: go %s, want go %s", file, current, version))
		}
	}
	if len(mismatches) > 0 {
		return fmt.Errorf("go.mod files do not match go %s:\n%s", version, strings.Join(mismatches, "\n"))
	}
	return nil
}

// update rewrites modules whose Go directive differs.
func update(ctx context.Context, root, git, version string) error {
	files, err := discover(ctx, root, git)
	if err != nil {
		return fmt.Errorf("list modules in %q: %w", root, err)
	}
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("process module %q: %w", file, err)
		}
		path := filepath.Join(root, filepath.FromSlash(file))
		data, err := readModule(path)
		if err != nil {
			return fmt.Errorf("read module %q: %w", file, err)
		}
		current, hasDirective, err := parseVersion(file, data)
		if err != nil {
			return fmt.Errorf("inspect module version: %w", err)
		}
		if hasDirective && current == version {
			continue
		}
		if err := updateFile(path, data, version); err != nil {
			return fmt.Errorf("update module %q: %w", file, err)
		}
	}
	return nil
}

// updateFile sets a Go directive while preserving other module declarations.
func updateFile(path string, data []byte, version string) (result error) {
	module, err := modfile.Parse(path, data, nil)
	if err != nil {
		return fmt.Errorf("parse module %q: %w", path, err)
	}
	if err := module.AddGoStmt(version); err != nil {
		return fmt.Errorf("set Go directive in %q: %w", path, err)
	}
	output, err := module.Format()
	if err != nil {
		return fmt.Errorf("format module %q: %w", path, err)
	}
	// A tracked symlink names the module to update. Preserve that link when replacing the file.
	destination, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("resolve module %q: %w", path, err)
	}
	info, err := os.Stat(destination)
	if err != nil {
		return fmt.Errorf("inspect module %q: %w", path, err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(destination), ".go-mod-*")
	if err != nil {
		return fmt.Errorf("create temporary module for %q: %w", path, err)
	}
	defer func() {
		if result != nil {
			if err := os.Remove(temporary.Name()); err != nil {
				result = errors.Join(result, fmt.Errorf("remove temporary module: %w", err))
			}
		}
	}()
	_, writeErr := temporary.Write(output)
	modeErr := temporary.Chmod(info.Mode().Perm())
	if err := errors.Join(writeErr, modeErr, temporary.Close()); err != nil {
		return fmt.Errorf("stage module %q: %w", path, err)
	}
	if err := os.Rename(temporary.Name(), destination); err != nil {
		return fmt.Errorf("replace module %q: %w", path, err)
	}
	return nil
}
