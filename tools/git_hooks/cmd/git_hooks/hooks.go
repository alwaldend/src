package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// precommit contains the hook installed into the checkout's effective hooks directory.
//
//go:embed precommit.sh
var precommit []byte

// execute dispatches hook installation (also the default) or verification.
func execute(ctx context.Context, args []string, workspace string, stdout, stderr io.Writer) error {
	install := func(cmd *cobra.Command, _ []string) error {
		if err := processHook(cmd.Context(), workspace, cmd.OutOrStdout(), true); err != nil {
			return fmt.Errorf("install repository hook: %w", err)
		}
		return nil
	}
	root := &cobra.Command{
		Use: "git_hooks", Short: "Install or verify repository hooks", Args: cobra.NoArgs,
		RunE: install, SilenceErrors: true, SilenceUsage: true,
	}
	root.AddCommand(&cobra.Command{Use: "install", Short: "Install the repository hook", Args: cobra.NoArgs, RunE: install})
	root.AddCommand(&cobra.Command{
		Use: "test", Short: "Verify the installed hook", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := processHook(cmd.Context(), workspace, cmd.OutOrStdout(), false); err != nil {
				return fmt.Errorf("verify repository hook: %w", err)
			}
			return nil
		},
	})
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("run git_hooks: %w", err)
	}
	return nil
}

// processHook applies the requested operation to the hook selected by Git.
func processHook(ctx context.Context, workspace string, stdout io.Writer, install bool) error {
	dir, err := resolveHooksDirectory(ctx, workspace)
	if err != nil {
		return fmt.Errorf("locate hook: %w", err)
	}
	path := filepath.Join(dir, "pre-commit")
	action := "verified"
	if install {
		if err := installHook(path, precommit); err != nil {
			return fmt.Errorf("replace %q: %w", path, err)
		}
		action = "installed"
	} else if err := checkHook(path, precommit); err != nil {
		return fmt.Errorf("check installed contents: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "%s %s\n", action, path); err != nil {
		return fmt.Errorf("report hook %q: %w", path, err)
	}
	return nil
}

// resolveHooksDirectory honors worktree configuration and core.hooksPath.
func resolveHooksDirectory(ctx context.Context, workspaceDirectory string) (string, error) {
	command := exec.CommandContext(ctx, "git", "rev-parse", "--git-path", "hooks")
	if workspaceDirectory != "" {
		command.Dir = workspaceDirectory
	}
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"could not resolve Git hooks directory: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	hooksDirectory := strings.TrimSpace(string(output))
	if hooksDirectory == "" {
		return "", fmt.Errorf("Git returned an empty hooks directory")
	}
	if filepath.IsAbs(hooksDirectory) {
		return filepath.Clean(hooksDirectory), nil
	}

	baseDirectory := workspaceDirectory
	if baseDirectory == "" {
		baseDirectory, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not get working directory: %w", err)
		}
	}
	return filepath.Clean(filepath.Join(baseDirectory, hooksDirectory)), nil
}

// installHook stages a complete executable in the destination directory before replacement.
func installHook(hookPath string, contents []byte) (resultErr error) {
	if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
		return fmt.Errorf("could not create hooks directory: %w", err)
	}

	temporary, err := os.CreateTemp(filepath.Dir(hookPath), ".git-hook-*")
	if err != nil {
		return fmt.Errorf("could not create temporary hook: %w", err)
	}
	temporaryPath := temporary.Name()
	defer func() {
		if err := os.Remove(temporaryPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove temporary hook: %w", err))
		}
	}()
	writeErr := writeHook(temporary, contents)
	closeErr := temporary.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("prepare temporary hook: %w", err)
	}
	if err := os.Rename(temporaryPath, hookPath); err != nil {
		return fmt.Errorf("could not install hook %q: %w", hookPath, err)
	}
	return nil
}

// writeHook finishes all writes and permissions before the caller closes the file.
func writeHook(file *os.File, contents []byte) error {
	if _, err := file.Write(contents); err != nil {
		return fmt.Errorf("write contents: %w", err)
	}
	if err := file.Chmod(0o755); err != nil {
		return fmt.Errorf("set executable permissions: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync contents: %w", err)
	}
	return nil
}

// checkHook checks file contents and executable permissions without changing them.
func checkHook(hookPath string, contents []byte) (resultErr error) {
	file, err := os.Open(hookPath)
	if err != nil {
		return fmt.Errorf("hook %q is not installed: %w", hookPath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close hook %q: %w", hookPath, err))
		}
	}()
	installed, err := io.ReadAll(io.LimitReader(file, int64(len(contents))+1))
	if err != nil {
		return fmt.Errorf("read hook %q: %w", hookPath, err)
	}
	if !bytes.Equal(installed, contents) {
		return fmt.Errorf("hook %q is out of date", hookPath)
	}

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not inspect hook %q: %w", hookPath, err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		return fmt.Errorf("hook %q is not executable", hookPath)
	}
	return nil
}
