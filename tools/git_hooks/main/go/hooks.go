package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed precommit.sh
var precommit []byte

type hook struct {
	name     string
	contents []byte
}

var repositoryHooks = []hook{
	{
		name:     "pre-commit",
		contents: precommit,
	},
}

func execute(
	args []string,
	getenv func(string) string,
	stdout io.Writer,
) error {
	command := "install"
	if len(args) > 1 {
		return fmt.Errorf("usage: git_hooks [install|test]")
	}
	if len(args) == 1 {
		command = args[0]
	}

	if command == "help" || command == "-h" || command == "--help" {
		_, err := fmt.Fprintln(stdout, "usage: git_hooks [install|test]")
		return err
	}
	if command != "install" && command != "test" {
		return fmt.Errorf("unknown command %q; usage: git_hooks [install|test]", command)
	}

	hooksDirectory, err := resolveHooksDirectory(
		getenv("BUILD_WORKSPACE_DIRECTORY"),
	)
	if err != nil {
		return err
	}

	for _, repositoryHook := range repositoryHooks {
		hookPath := filepath.Join(hooksDirectory, repositoryHook.name)
		if command == "install" {
			if err := installHook(hookPath, repositoryHook.contents); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(stdout, "installed %s\n", hookPath); err != nil {
				return err
			}
			continue
		}
		if err := checkHook(hookPath, repositoryHook.contents); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(stdout, "verified %s\n", hookPath); err != nil {
			return err
		}
	}
	return nil
}

func resolveHooksDirectory(workspaceDirectory string) (string, error) {
	command := exec.Command("git", "rev-parse", "--git-path", "hooks")
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

func installHook(hookPath string, contents []byte) error {
	if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
		return fmt.Errorf("could not create hooks directory: %w", err)
	}

	temporary, err := os.CreateTemp(filepath.Dir(hookPath), ".git-hook-*")
	if err != nil {
		return fmt.Errorf("could not create temporary hook: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return fmt.Errorf("could not write temporary hook: %w", err)
	}
	if err := temporary.Chmod(0o755); err != nil {
		temporary.Close()
		return fmt.Errorf("could not make temporary hook executable: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("could not sync temporary hook: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("could not close temporary hook: %w", err)
	}
	if err := os.Rename(temporaryPath, hookPath); err != nil {
		return fmt.Errorf("could not install hook %q: %w", hookPath, err)
	}
	return nil
}

func checkHook(hookPath string, contents []byte) error {
	installed, err := os.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("hook %q is not installed: %w", hookPath, err)
	}
	if !bytes.Equal(installed, contents) {
		return fmt.Errorf("hook %q is out of date", hookPath)
	}

	info, err := os.Stat(hookPath)
	if err != nil {
		return fmt.Errorf("could not inspect hook %q: %w", hookPath, err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		return fmt.Errorf("hook %q is not executable", hookPath)
	}
	return nil
}
