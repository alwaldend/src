package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"golang.org/x/mod/modfile"
)

var (
	goVersion       = flag.String("version", "", "Go version to set in go.mod files")
	gitRunfile      = flag.String("git", "", "Git executable runfile")
	checkOnly       = flag.Bool("check", false, "check without writing")
	workspaceMarker = flag.String("workspace-marker", "", "workspace module runfile")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if *goVersion == "" {
		return fmt.Errorf("--version is required")
	}
	root, err := rootDir(*workspaceMarker)
	if err != nil {
		return err
	}
	gitPath, err := runfiles.Rlocation(*gitRunfile)
	if err != nil {
		return err
	}
	if *checkOnly {
		return check(root, gitPath, *goVersion)
	}
	return update(root, gitPath, *goVersion)
}

func rootDir(marker string) (string, error) {
	path, err := runfiles.Rlocation(marker)
	if err != nil {
		return "", err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

func discover(root, git string) ([]string, error) {
	cmd := exec.Command(git, "--git-dir", filepath.Join(root, ".git"), "-C", root, "ls-files", "-z", "go.mod", "**/go.mod")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("discover go.mod files: %w", err)
	}
	var files []string
	for _, name := range strings.Split(string(output), "\x00") {
		if name != "" {
			files = append(files, name)
		}
	}
	return files, nil
}

func parseVersion(path string, data []byte) (string, bool, error) {
	file, err := modfile.Parse(path, data, nil)
	if err != nil {
		return "", false, err
	}
	if file.Go == nil {
		return "", false, nil
	}
	return file.Go.Version, true, nil
}

func check(root, git, version string) error {
	files, err := discover(root, git)
	if err != nil {
		return err
	}
	var mismatches []string
	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file)))
		if err != nil {
			return err
		}
		current, hasDirective, err := parseVersion(file, data)
		if err != nil {
			return fmt.Errorf("%s: %v", file, err)
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

func update(root, git, version string) error {
	files, err := discover(root, git)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(root, filepath.FromSlash(file))
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		current, hasDirective, err := parseVersion(file, data)
		if err != nil {
			return fmt.Errorf("%s: %v", file, err)
		}
		if hasDirective && current == version {
			continue
		}
		if err := updateFile(path, data, version); err != nil {
			return err
		}
	}
	return nil
}

func updateFile(path string, data []byte, version string) error {
	module, err := modfile.Parse(path, data, nil)
	if err != nil {
		return err
	}
	if err := module.AddGoStmt(version); err != nil {
		return err
	}
	output, err := module.Format()
	if err != nil {
		return err
	}
	return os.WriteFile(path, output, 0o644)
}
