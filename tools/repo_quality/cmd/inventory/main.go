package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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
	// Statuses contains the single project-owned lifecycle declaration.
	Statuses []string `yaml:"statuses"`
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

// statusField matches only the top-level YAML status field and its generated comment.
var statusField = regexp.MustCompile(`(?m)^(?:# Generated from project README by //tools/repo_quality:inventory_update\.\n)?statuses:[^\n]*(?:\n(?:[ \t]+|-[ \t])[^\n]*)*\n?`)

// frontmatter returns the YAML header boundary, requiring explicit delimiters.
func frontmatter(data []byte, path string) (int, error) {
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return 0, fmt.Errorf("YAML front matter required in %q", path)
	}
	end := bytes.Index(data[4:], []byte("\n---\n"))
	if end < 0 {
		return 0, fmt.Errorf("unterminated front matter in %q", path)
	}
	return 4 + end, nil
}

// statusCatalog reads the shared lifecycle identifiers and their definitions.
func statusCatalog(root string) (map[string]string, error) {
	path := filepath.Join(root, "projects/statuses.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read status catalog %q: %w", path, err)
	}
	var catalog map[string]string
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("decode status catalog %q: %w", path, err)
	}
	if len(catalog) == 0 {
		return nil, fmt.Errorf("empty status catalog %q", path)
	}
	for name, description := range catalog {
		if strings.TrimSpace(name) != name || name == "" || strings.TrimSpace(description) == "" {
			return nil, fmt.Errorf("invalid status definition %q in %q", name, path)
		}
	}
	return catalog, nil
}

// landingStatus derives only the lifecycle field, preserving the rest of the page.
func landingStatus(data []byte, path, status string) ([]byte, error) {
	end, err := frontmatter(data, path)
	if err != nil {
		return nil, fmt.Errorf("read landing metadata: %w", err)
	}
	header := data[4:end]
	var fields map[string]any
	if err := yaml.Unmarshal(header, &fields); err != nil {
		return nil, fmt.Errorf("decode landing metadata %q: %w", path, err)
	}
	header = append(bytes.Clone(header), '\n')
	field := []byte("# Generated from project README by //tools/repo_quality:inventory_update.\nstatuses:\n  - " + status + "\n")
	if _, exists := fields["statuses"]; exists {
		if !statusField.Match(header) {
			return nil, fmt.Errorf("top-level YAML statuses field required in %q", path)
		}
		header = statusField.ReplaceAllLiteral(header, field)
	} else {
		header = append(header, field...)
	}
	var projected metadata
	if err := yaml.Unmarshal(header, &projected); err != nil {
		return nil, fmt.Errorf("validate generated landing metadata %q: %w", path, err)
	}
	if len(projected.Statuses) != 1 || projected.Statuses[0] != status {
		return nil, fmt.Errorf("cannot safely replace landing status in %q", path)
	}
	result := append([]byte("---\n"), header...)
	return append(result, data[end+1:]...), nil
}

// cleanCell escapes Markdown table cells without changing their displayed text.
func cleanCell(value string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(value), " "), "|", "\\|")
}

// statusDefinitions renders the catalog as a generated documentation table.
func statusDefinitions(catalog map[string]string) []byte {
	names := make([]string, 0, len(catalog))
	for name := range catalog {
		names = append(names, name)
	}
	sort.Strings(names)
	rows := [][]string{{"Status", "Meaning"}}
	for _, name := range names {
		rows = append(rows, []string{"`" + cleanCell(name) + "`", cleanCell(catalog[name])})
	}
	return append([]byte("<!-- Generated by //tools/repo_quality:inventory_update from statuses.yaml. -->\n\n"), markdownTable(rows)...)
}

// projectIndex renders descriptions from each direct project's README metadata.
func projectIndex(root string, catalog map[string]string, projections map[string][]byte) ([]byte, error) {
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
	rows := [][]string{{"Project", "Status", "Description", "Landing page"}}
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
		end, err := frontmatter(data, path)
		if err != nil {
			return nil, fmt.Errorf("read project metadata: %w", err)
		}
		if err := yaml.Unmarshal(data[4:end], &value); err != nil {
			return nil, fmt.Errorf("decode project metadata %q: %w", path, err)
		}
		if len(value.Statuses) != 1 {
			return nil, fmt.Errorf("exactly one lifecycle status required in %q", path)
		}
		status := value.Statuses[0]
		if _, exists := catalog[status]; !exists {
			return nil, fmt.Errorf("unknown lifecycle status %q in %q; see projects/statuses.yaml", status, path)
		}
		landingPath := filepath.Join("projects", entry.Name(), "site/content/_index.md")
		landingData, err := os.ReadFile(filepath.Join(root, landingPath))
		if err == nil {
			projected, err := landingStatus(landingData, landingPath, status)
			if err != nil {
				return nil, fmt.Errorf("project landing status: %w", err)
			}
			projections[landingPath] = projected
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read landing %q: %w", landingPath, err)
		}
		landing := "None"
		if landings[entry.Name()] {
			landing = "[Site](https://alwaldend.com/projects/" + entry.Name() + "/)"
		}
		if entry.Name() == "alwaldend.com" {
			landing = "[Site](https://alwaldend.com/)"
		}
		rows = append(rows, []string{fmt.Sprintf("[%s](%s/README.md)", cleanCell(value.Title), entry.Name()), status, cleanCell(value.Description), landing})
	}
	return append([]byte("<!-- Generated by //tools/repo_quality:inventory_update from project READMEs. -->\n\n"), markdownTable(rows)...), nil
}

// markdownTable pads generated cells to match the repository Markdown formatter.
func markdownTable(rows [][]string) []byte {
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for column, value := range row {
			if width := utf8.RuneCountInString(value); width > widths[column] {
				widths[column] = width
			}
		}
	}
	var output strings.Builder
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
	return []byte(output.String())
}

// replaceSection updates one generated README block without changing surrounding prose.
func replaceSection(data []byte, name string, content []byte) ([]byte, error) {
	begin := []byte("<!-- BEGIN GENERATED PROJECT " + name + " -->")
	end := []byte("<!-- END GENERATED PROJECT " + name + " -->")
	start, finish := bytes.Index(data, begin), bytes.Index(data, end)
	if start < 0 || finish < start {
		return nil, fmt.Errorf("generated project %s markers missing", name)
	}
	result := append(bytes.Clone(data[:start+len(begin)]), '\n')
	result = append(result, content...)
	result = append(result, '\n')
	return append(result, data[finish:]...), nil
}

// synchronize checks or updates generated inventory and lifecycle projections.
func synchronize(root string, check bool) error {
	workspaces, err := workspaceInventory(root)
	if err != nil {
		return fmt.Errorf("generate workspace inventory: %w", err)
	}
	catalog, err := statusCatalog(root)
	if err != nil {
		return fmt.Errorf("load project lifecycle: %w", err)
	}
	projections := map[string][]byte{"infra/src/openspec/validation/workspaces.bzl": workspaces}
	index, err := projectIndex(root, catalog, projections)
	if err != nil {
		return fmt.Errorf("generate project index: %w", err)
	}
	readmePath := filepath.Join(root, "projects/README.md")
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("read project index %q: %w", readmePath, err)
	}
	readme, err = replaceSection(readme, "INDEX", index)
	if err != nil {
		return fmt.Errorf("update project index: %w", err)
	}
	readme, err = replaceSection(readme, "STATUSES", statusDefinitions(catalog))
	if err != nil {
		return fmt.Errorf("update status definitions: %w", err)
	}
	projections["projects/README.md"] = readme
	paths := make([]string, 0, len(projections))
	for path := range projections {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		content := projections[path]
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
