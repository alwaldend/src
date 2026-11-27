package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/spf13/cobra"
)

// main forwards the test runner arguments to the packaged linter.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	command := &cobra.Command{
		Use: "ansible-lint-runner", DisableFlagParsing: true, SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(cmd.Context(), args); err != nil {
				return fmt.Errorf("run Ansible lint: %w", err)
			}
			return nil
		},
	}
	err := command.ExecuteContext(ctx)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run prepares an isolated writable Ansible project and invokes the linter there.
func run(ctx context.Context, args []string) (result error) {
	workspaceRoot := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"))
	if os.Getenv("TEST_SRCDIR") == "" || os.Getenv("TEST_WORKSPACE") == "" {
		return fmt.Errorf("runfiles root is unavailable")
	}
	packageRoot := os.Getenv("ANSIBLE_LINT_PACKAGE_ROOT")
	if packageRoot == "" {
		return fmt.Errorf("ANSIBLE_LINT_PACKAGE_ROOT is unavailable")
	}

	binDir, err := os.MkdirTemp("", "ansible-lint-bin")
	if err != nil {
		return fmt.Errorf("create writable bin directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(binDir); err != nil {
			result = errors.Join(result, fmt.Errorf("remove %q: %w", binDir, err))
		}
	}()

	for name, target := range map[string]string{
		"ansible":          "ansible",
		"ansible-config":   "ansible_config",
		"ansible-galaxy":   "ansible_galaxy",
		"ansible-lint":     "ansible_lint",
		"ansible-playbook": "ansible_playbook",
	} {
		source, found := bazel.FindBinary("tools/ansible", target)
		if !found {
			return fmt.Errorf("find %s", target)
		}
		if err := os.Symlink(source, filepath.Join(binDir, name)); err != nil {
			return fmt.Errorf("link %s: %w", name, err)
		}
	}

	home, err := os.MkdirTemp("", "ansible-lint-home")
	if err != nil {
		return fmt.Errorf("create writable home directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(home); err != nil {
			result = errors.Join(result, fmt.Errorf("remove %q: %w", home, err))
		}
	}()

	projectRoot := filepath.Join(home, "project")
	if err := os.Mkdir(projectRoot, 0o700); err != nil {
		return fmt.Errorf("create project directory: %w", err)
	}
	if err := copySources(filepath.Join(workspaceRoot, packageRoot), projectRoot); err != nil {
		return fmt.Errorf("populate writable project: %w", err)
	}
	collectionRoot := filepath.Join(workspaceRoot, "projects/ansible_collection")
	if err := os.MkdirAll(filepath.Join(projectRoot, "collections/ansible_collections"), 0o700); err != nil {
		return fmt.Errorf("create collection path: %w", err)
	}
	if err := copySources(collectionRoot, filepath.Join(projectRoot, "collections/ansible_collections/alwaldend/main")); err != nil {
		return fmt.Errorf("populate collection: %w", err)
	}

	localTemp := filepath.Join(home, "ansible-local")
	cache := filepath.Join(home, "ansible-cache")
	for _, directory := range []string{localTemp, cache} {
		if err := os.Mkdir(directory, 0o700); err != nil {
			return fmt.Errorf("create writable directory: %w", err)
		}
	}
	config := fmt.Sprintf("[defaults]\ncollections_path = %s\nroles_path = %s\n", filepath.Join(projectRoot, "collections"), home)
	configPath := filepath.Join(home, "ansible.cfg")
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		return fmt.Errorf("write ansible.cfg: %w", err)
	}

	environment := os.Environ()
	environment = append(environment,
		"ANSIBLE_CACHE_PLUGIN_CONNECTION="+cache,
		"ANSIBLE_CONFIG="+configPath,
		"ANSIBLE_LOCAL_TEMP="+localTemp,
		"HOME="+home,
		"PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"),
	)

	arguments := make([]string, 0, len(args)+2)
	for _, argument := range args {
		arguments = append(arguments, strings.TrimPrefix(argument, packageRoot+"/"))
	}
	arguments = append(arguments, "--exclude", filepath.Join(projectRoot, "collections"))
	command := exec.CommandContext(ctx, filepath.Join(binDir, "ansible-lint"), arguments...)
	command.Dir = projectRoot
	command.Env = environment
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("ansible-lint: %w", err)
	}
	return nil
}

// copySources copies packaged files into the writable lint tree.
func copySources(sourceRoot, destinationRoot string) error {
	err := filepath.WalkDir(sourceRoot, func(source string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("visit %q: %w", source, walkErr)
		}
		if entry.IsDir() {
			return nil
		}
		relativePath, err := filepath.Rel(sourceRoot, source)
		if err != nil {
			return fmt.Errorf("resolve project path: %w", err)
		}
		destination := filepath.Join(destinationRoot, relativePath)
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return fmt.Errorf("create project path: %w", err)
		}
		if err := copyFile(source, destination); err != nil {
			return fmt.Errorf("copy %q: %w", relativePath, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("copy source tree %q: %w", sourceRoot, err)
	}
	return nil
}

// copyFile streams a packaged source into a private writable file.
func copyFile(source, destination string) (result error) {
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close source: %w", err))
		}
	}()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open destination %q: %w", destination, err)
	}
	_, copyErr := io.Copy(output, input)
	if err := errors.Join(copyErr, output.Close()); err != nil {
		return fmt.Errorf("write destination %q: %w", destination, err)
	}
	return nil
}
