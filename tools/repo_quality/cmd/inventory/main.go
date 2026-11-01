package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/bazelbuild/buildtools/build"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// metadata contains the project-owned human-facing index fields.
type metadata struct {
	// Title is the component's published name.
	Title string `yaml:"title"`
	// Description summarizes the component's responsibility.
	Description string `yaml:"description"`
}

// parse reads a Starlark declaration file without evaluating it.
func parse(path string) (*build.File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read declarations %q: %w", path, err)
	}
	file, err := build.ParseBuild(path, data)
	if err != nil {
		return nil, fmt.Errorf("parse declarations %q: %w", path, err)
	}
	return file, nil
}

// namedCall returns whether a declaration calls function with the given name attribute.
func namedCall(file *build.File, function, name string) bool {
	for _, statement := range file.Stmt {
		call, ok := statement.(*build.CallExpr)
		if !ok {
			continue
		}
		identifier, ok := call.X.(*build.Ident)
		if !ok || identifier.Name != function {
			continue
		}
		for _, argument := range call.List {
			assignment, ok := argument.(*build.AssignExpr)
			if !ok {
				continue
			}
			key, ok := assignment.LHS.(*build.Ident)
			value, stringValue := assignment.RHS.(*build.StringExpr)
			if ok && stringValue && key.Name == "name" && value.Value == name {
				return true
			}
		}
	}
	return false
}

// moduleName reads the module identity used by an enclosing standalone workspace.
func moduleName(path string) (string, error) {
	file, err := parse(path)
	if err != nil {
		return "", fmt.Errorf("resolve module identity %q: %w", path, err)
	}
	for _, statement := range file.Stmt {
		call, ok := statement.(*build.CallExpr)
		if !ok {
			continue
		}
		identifier, ok := call.X.(*build.Ident)
		if !ok || identifier.Name != "module" {
			continue
		}
		for _, argument := range call.List {
			assignment, ok := argument.(*build.AssignExpr)
			if !ok {
				continue
			}
			key, ok := assignment.LHS.(*build.Ident)
			value, stringValue := assignment.RHS.(*build.StringExpr)
			if ok && stringValue && key.Name == "name" {
				return value.Value, nil
			}
		}
	}
	return "", fmt.Errorf("module name missing in %q", path)
}

// sourceLabel resolves an owner export using its nearest module boundary.
func sourceLabel(root, owner string) (string, error) {
	for directory := owner; directory != "."; directory = filepath.Dir(directory) {
		path := filepath.Join(root, directory, "MODULE.bazel")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("inspect module %q: %w", path, err)
		}
		name, err := moduleName(path)
		if err != nil {
			return "", fmt.Errorf("resolve owner module %q: %w", owner, err)
		}
		relative, err := filepath.Rel(directory, filepath.Join(owner, "openspec"))
		if err != nil {
			return "", fmt.Errorf("resolve source owner %q: %w", owner, err)
		}
		return "@" + name + "//" + filepath.ToSlash(relative) + ":source", nil
	}
	return "//" + filepath.ToSlash(owner) + "/openspec:source", nil
}

// workspaceInventory discovers declared specification sources, excluding draft-only workspaces.
func workspaceInventory(root string) ([]byte, error) {
	entries := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("visit %q: %w", path, walkErr)
		}
		if entry.IsDir() {
			if path != root && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "out" || entry.Name() == "node_modules" || strings.HasPrefix(entry.Name(), "bazel-")) {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() != "BUILD.bazel" || filepath.Base(filepath.Dir(path)) != "openspec" {
			return nil
		}
		file, err := parse(path)
		if err != nil {
			return fmt.Errorf("read source export %q: %w", path, err)
		}
		if !namedCall(file, "openspec_workspace", "source") && !namedCall(file, "filegroup", "source") {
			return nil
		}
		owner, err := filepath.Rel(root, filepath.Dir(filepath.Dir(path)))
		if err != nil {
			return fmt.Errorf("resolve owner %q: %w", path, err)
		}
		label, err := sourceLabel(root, owner)
		if err != nil {
			return fmt.Errorf("resolve source label %q: %w", owner, err)
		}
		key := strings.NewReplacer("/", "_", ".", "_", "-", "_").Replace(filepath.ToSlash(owner))
		if owner == filepath.Join("infra", "src") {
			key = "repository"
		}
		if previous, exists := entries[key]; exists {
			return fmt.Errorf("inventory key %q collides: %s and %s", key, previous, label)
		}
		entries[key] = label
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover specification exports: %w", err)
	}
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var output strings.Builder
	output.WriteString("# Generated by //tools/repo_quality:inventory_update. Do not edit.\nWORKSPACE_SOURCES = {\n")
	for _, key := range keys {
		fmt.Fprintf(&output, "    %q: %q,\n", key, entries[key])
	}
	output.WriteString("}\n")
	file, err := build.ParseBzl("workspaces.bzl", []byte(output.String()))
	if err != nil {
		return nil, fmt.Errorf("format inventory: %w", err)
	}
	return build.Format(file), nil
}

// projectIndex renders descriptions from each direct project's README metadata.
func projectIndex(root string) ([]byte, error) {
	entries, err := os.ReadDir(filepath.Join(root, "projects"))
	if err != nil {
		return nil, fmt.Errorf("list project owners: %w", err)
	}
	landings := map[string]bool{}
	registryPath := filepath.Join(root, "projects/projects.bzl")
	if _, err := os.Stat(registryPath); err == nil {
		registry, err := parse(registryPath)
		if err != nil {
			return nil, fmt.Errorf("read landing registry: %w", err)
		}
		for _, statement := range registry.Stmt {
			assignment, ok := statement.(*build.AssignExpr)
			if !ok {
				continue
			}
			key, ok := assignment.LHS.(*build.Ident)
			if !ok || key.Name != "PROJECTS" {
				continue
			}
			list, ok := assignment.RHS.(*build.ListExpr)
			if !ok {
				return nil, fmt.Errorf("PROJECTS in %q must be a literal list", registryPath)
			}
			for _, item := range list.List {
				value, ok := item.(*build.StringExpr)
				if !ok {
					return nil, fmt.Errorf("PROJECTS in %q must contain strings", registryPath)
				}
				landings[value.Value] = true
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("inspect landing registry: %w", err)
	}
	rows := [][3]string{{"Project", "Description", "Landing page"}}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(root, "projects", entry.Name(), "README.md")
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read project metadata %q: %w", path, err)
		}
		value := metadata{Title: entry.Name()}
		if bytes.HasPrefix(data, []byte("---\n")) {
			end := bytes.Index(data[4:], []byte("\n---"))
			if end < 0 {
				return nil, fmt.Errorf("unterminated frontmatter in %q", path)
			}
			if err := yaml.Unmarshal(data[4:4+end], &value); err != nil {
				return nil, fmt.Errorf("decode project metadata %q: %w", path, err)
			}
		} else {
			paragraphs := strings.Split(strings.TrimSpace(string(data)), "\n\n")
			if len(paragraphs) > 0 && strings.HasPrefix(paragraphs[0], "# ") {
				value.Title = strings.TrimPrefix(paragraphs[0], "# ")
			}
			if len(paragraphs) > 1 {
				value.Description = paragraphs[1]
			}
		}
		clean := func(value string) string {
			return strings.ReplaceAll(strings.Join(strings.Fields(value), " "), "|", "\\|")
		}
		landing := "None"
		if landings[entry.Name()] {
			landing = "[Site](https://alwaldend.com/projects/" + entry.Name() + "/)"
		}
		if entry.Name() == "alwaldend.com" {
			landing = "[Site](https://alwaldend.com/)"
		}
		rows = append(rows, [3]string{fmt.Sprintf("[%s](%s/README.md)", clean(value.Title), entry.Name()), clean(value.Description), landing})
	}
	widths := [3]int{}
	for _, row := range rows {
		for column, value := range row {
			if width := utf8.RuneCountInString(value); width > widths[column] {
				widths[column] = width
			}
		}
	}
	var output strings.Builder
	output.WriteString("<!-- Generated by //tools/repo_quality:inventory_update from project READMEs. -->\n\n")
	for index, row := range rows {
		output.WriteString("|")
		for column, value := range row {
			fmt.Fprintf(&output, " %s%s |", value, strings.Repeat(" ", widths[column]-utf8.RuneCountInString(value)))
		}
		output.WriteByte('\n')
		if index == 0 {
			output.WriteString("|")
			for _, width := range widths {
				fmt.Fprintf(&output, " %s |", strings.Repeat("-", width))
			}
			output.WriteByte('\n')
		}
	}
	return []byte(output.String()), nil
}

// synchronize checks or updates both generated inventory projections.
func synchronize(root string, check bool) error {
	workspaces, err := workspaceInventory(root)
	if err != nil {
		return fmt.Errorf("generate workspace inventory: %w", err)
	}
	index, err := projectIndex(root)
	if err != nil {
		return fmt.Errorf("generate project index: %w", err)
	}
	index = append(index, '\n')
	readmePath := filepath.Join(root, "projects/README.md")
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("read project index %q: %w", readmePath, err)
	}
	const begin = "<!-- BEGIN GENERATED PROJECT INDEX -->"
	const end = "<!-- END GENERATED PROJECT INDEX -->"
	start, finish := bytes.Index(readme, []byte(begin)), bytes.Index(readme, []byte(end))
	if start < 0 || finish < start {
		return fmt.Errorf("generated project index markers missing in %q", readmePath)
	}
	readme = append(append(append([]byte{}, readme[:start+len(begin)]...), '\n'), append(index, readme[finish:]...)...)
	for path, content := range map[string][]byte{"infra/src/openspec/validation/workspaces.bzl": workspaces, "projects/README.md": readme} {
		path = filepath.Join(root, path)
		previous, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("read generated inventory %q: %w", path, err)
		}
		if bytes.Equal(previous, content) {
			continue
		}
		if check {
			return fmt.Errorf("stale inventory %q; run //tools/repo_quality:inventory_update", path)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create inventory directory %q: %w", path, err)
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			return fmt.Errorf("write inventory %q: %w", path, err)
		}
	}
	return nil
}

// main exposes the inventory updater and its non-mutating freshness mode.
func main() {
	root := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	check := false
	command := &cobra.Command{Use: "inventory", SilenceUsage: true, RunE: func(_ *cobra.Command, args []string) error {
		if root == "" {
			return fmt.Errorf("--root or BUILD_WORKSPACE_DIRECTORY is required")
		}
		return synchronize(root, check)
	}}
	command.Args = cobra.NoArgs
	command.Flags().StringVar(&root, "root", root, "Repository root to inspect")
	command.Flags().BoolVar(&check, "check", false, "Reject stale inventories without writing")
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
