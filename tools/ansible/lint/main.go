package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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
	defer os.RemoveAll(binDir)

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
	defer os.RemoveAll(home)

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

	arguments := make([]string, 0, len(os.Args)-1)
	for _, argument := range os.Args[1:] {
		arguments = append(arguments, strings.TrimPrefix(argument, packageRoot+"/"))
	}
	arguments = append(arguments, "--exclude", filepath.Join(projectRoot, "collections"))
	command := exec.Command(filepath.Join(binDir, "ansible-lint"), arguments...)
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

func copySources(sourceRoot, destinationRoot string) error {
	return filepath.WalkDir(sourceRoot, func(source string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
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
		content, err := os.ReadFile(source)
		if err != nil {
			return fmt.Errorf("read %s: %w", relativePath, err)
		}
		if err := os.WriteFile(destination, content, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", relativePath, err)
		}
		return nil
	})
}
