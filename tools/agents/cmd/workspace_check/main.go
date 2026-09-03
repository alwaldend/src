// Command workspace_check compiles the bounded deterministic
// WorkspaceCheckCatalog over tracked MODULE.bazel roots and their configured
// check projections.
//
// It performs no network or stateful operation. Outputs are checked
// repository artifacts: portable JSON plus a human Markdown render that
// states the JSON digest.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

const (
	bazelIgnorePath     = ".bazelignore"
	fullCheckScriptPath = "projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go"
)

var moduleNamePattern = regexp.MustCompile(`(?m)^\s*(?:name|module\(\s*name)\s*=\s*"([^"]+)"`)

type options struct {
	workspaceRoot  string
	outputPath     string
	markdownPath   string
	check          bool
	sourceRevision string
	producerRef    string
}

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("workspace_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/workspace-check.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/workspace-check.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on completeness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: content-addressed inputs)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.workspace-check-compiler",
		"producer reference for the catalog")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	return opts, nil
}

type compiler struct {
	root        string
	opts        options
	inputs      []catalogv1alpha1.CatalogInput
	workspaces  []catalogv1alpha1.WorkspaceRecord
	conflicts   []catalogv1alpha1.CatalogConflict
	problems    []string
	ignored     map[string]bool
	eligible    int
	emitted     int
	unavailable int
}

func (c *compiler) problem(format string, args ...any) {
	c.problems = append(c.problems, fmt.Sprintf(format, args...))
}

func (c *compiler) input(relativePath, role string, content []byte) {
	value := sha256.Sum256(content)
	c.inputs = append(c.inputs, catalogv1alpha1.CatalogInput{
		Path:   filepath.ToSlash(relativePath),
		Role:   role,
		Digest: "sha256:" + hex.EncodeToString(value[:]),
	})
}

func readModuleName(content []byte) (string, error) {
	match := moduleNamePattern.FindSubmatch(content)
	if match == nil {
		return "", fmt.Errorf("no module name")
	}
	return string(match[1]), nil
}

func workspaceID(modulePath string) string {
	dir := filepath.ToSlash(filepath.Dir(modulePath))
	if dir == "." {
		return "root"
	}
	return strings.ReplaceAll(dir, "/", ".")
}

// walkModuleRoots returns the sorted relative paths of every tracked
// MODULE.bazel root below root, skipping build and scratch directories.
func walkModuleRoots(root string) ([]string, error) {
	var modules []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			name := entry.Name()
			if path != root && (name == ".git" || name == "out" ||
				strings.HasPrefix(name, "bazel-")) {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() == "MODULE.bazel" {
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			modules = append(modules, filepath.ToSlash(relative))
		}
		return nil
	})
	sort.Strings(modules)
	return modules, err
}

func (c *compiler) readBazelIgnore() map[string]bool {
	ignored := map[string]bool{}
	full := filepath.Join(c.root, filepath.FromSlash(bazelIgnorePath))
	content, err := os.ReadFile(full)
	if err != nil {
		c.problem("bazelignore missing: %s", bazelIgnorePath)
		c.eligible++
		c.unavailable++
		return ignored
	}
	c.input(bazelIgnorePath, "bazel-ignore", content)
	c.eligible++
	c.emitted++
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignored[line] = true
	}
	return ignored
}

func (c *compiler) readFullCheckScript() ([]byte, bool) {
	full := filepath.Join(c.root, filepath.FromSlash(fullCheckScriptPath))
	content, err := os.ReadFile(full)
	if err != nil {
		c.problem("full check script missing: %s", fullCheckScriptPath)
		c.eligible++
		c.unavailable++
		return nil, false
	}
	c.input(fullCheckScriptPath, "full-check-script", content)
	c.eligible++
	c.emitted++
	return content, true
}

func (c *compiler) fullCheckWorkspace(script []byte, modulePath string) bool {
	fullScript := string(script)
	pathQuote := filepath.ToSlash(filepath.Dir(modulePath))
	pathQuote = strings.TrimPrefix(pathQuote, "./")
	if pathQuote == "." {
		return true
	}
	return strings.Contains(fullScript, `path: "`+pathQuote+`"`)
}

// projectWorkspace reports whether a module path sits directly under the
// projects universe.
func projectWorkspace(modulePath string) bool {
	dir := filepath.ToSlash(filepath.Dir(modulePath))
	return strings.HasPrefix(dir, "projects/")
}

// moduleDirProjection returns the deterministic projection for one workspace
// root. Uncertain signals resolve to false and add no problems, keeping the
// projection conservative.
func (c *compiler) moduleDirProjection(modulePath string, script []byte) catalogv1alpha1.WorkspaceProjection {
	dir := filepath.ToSlash(filepath.Dir(modulePath))
	projection := catalogv1alpha1.WorkspaceProjection{
		BazelIgnore: c.ignored[dir],
	}
	if dir == "." {
		projection.FullCheck = true
	} else {
		projection.RootOverride = true
		if fullCheck := c.fullCheckWorkspace(script, modulePath); fullCheck {
			projection.FullCheck = true
		}
		docsBuild := filepath.Join(c.root, filepath.FromSlash(dir), "docs", "BUILD.bazel")
		if _, err := os.Stat(docsBuild); err == nil {
			projection.DocsAggregation = true
		}
	}
	return projection
}

func (c *compiler) compile() error {
	c.ignored = c.readBazelIgnore()
	script, ok := c.readFullCheckScript()
	if !ok {
		return nil
	}
	modules, err := walkModuleRoots(c.root)
	if err != nil {
		return err
	}
	for _, modulePath := range modules {
		c.eligible++
		full := filepath.Join(c.root, filepath.FromSlash(modulePath))
		content, err := os.ReadFile(full)
		if err != nil {
			c.problem("workspace module unreadable: %s", modulePath)
			c.unavailable++
			continue
		}
		moduleName, err := readModuleName(content)
		if err != nil {
			c.problem("workspace module name missing: %s", modulePath)
			c.unavailable++
			continue
		}
		c.input(modulePath, "workspace-module", content)
		c.emitted++
		projection := c.moduleDirProjection(modulePath, script)
		if projectWorkspace(modulePath) {
			projection.RootOverride = true
		}
		id := workspaceID(modulePath)
		workspace := catalogv1alpha1.WorkspaceRecord{
			ID:          id,
			Path:        filepath.ToSlash(filepath.Dir(modulePath)),
			ModulePath:  modulePath,
			ModuleName:  moduleName,
			Projections: projection,
			Phases: []catalogv1alpha1.WorkspacePhase{{
				ID:              id + ".check",
				ProviderRef:     "repository.bazel-operations",
				CommandTemplate: "bazel_agent bazel test //...",
			}},
		}
		c.workspaces = append(c.workspaces, workspace)
	}
	return nil
}

func (c *compiler) catalog() (catalogv1alpha1.WorkspaceCheckCatalog, error) {
	sourceRevision := c.opts.sourceRevision
	if sourceRevision == "" {
		var aggregate bytes.Buffer
		for _, input := range c.inputs {
			aggregate.WriteString(input.Path)
			aggregate.WriteString("\x00")
			aggregate.WriteString(input.Digest)
			aggregate.WriteString("\n")
		}
		value := sha256.Sum256(aggregate.Bytes())
		sourceRevision = hex.EncodeToString(value[:16])
	}
	completeness := catalogv1alpha1.CompletenessComplete
	limitations := []string{}
	if len(c.problems) > 0 {
		completeness = catalogv1alpha1.CompletenessPartial
		limitations = c.problems
	}
	inputs := c.inputs
	if inputs == nil {
		inputs = []catalogv1alpha1.CatalogInput{}
	}
	sort.Slice(inputs, func(i, j int) bool {
		if inputs[i].Path != inputs[j].Path {
			return inputs[i].Path < inputs[j].Path
		}
		return inputs[i].Role < inputs[j].Role
	})
	conflicts := c.conflicts
	if conflicts == nil {
		conflicts = []catalogv1alpha1.CatalogConflict{}
	}
	workspaces := c.workspaces
	if workspaces == nil {
		workspaces = []catalogv1alpha1.WorkspaceRecord{}
	}
	sort.Slice(workspaces, func(i, j int) bool {
		return workspaces[i].ID < workspaces[j].ID
	})
	return catalogv1alpha1.WorkspaceCheckCatalog{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            catalogv1alpha1.APIVersion + "/" + catalogv1alpha1.KindWorkspaceCheckCatalog,
			Kind:              catalogv1alpha1.KindWorkspaceCheckCatalog,
			ID:                "agent-system.workspace-check",
			DerivationVersion: "1.0.0",
			ProducerRef:       c.opts.producerRef,
			SourceRevision:    sourceRevision,
			Inputs:            inputs,
			Bounds: catalogv1alpha1.CatalogBounds{
				Eligible:       c.eligible,
				Emitted:        c.emitted,
				Unavailable:    c.unavailable,
				MaxItems:       1000,
				MaxInputBytes:  32 << 20,
				MaxOutputBytes: 1 << 20,
			},
			Completeness: completeness,
			Limitations:  limitations,
			Conflicts:    conflicts,
		},
		Workspaces: workspaces,
	}, nil
}

func (c *compiler) buildOutputs(catalog catalogv1alpha1.WorkspaceCheckCatalog) ([]byte, string, error) {
	jsonContent, err := catalogv1alpha1.CanonicalJSONWorkspaceCheck(catalog)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	var digestValue catalogv1alpha1.WorkspaceCheckCatalog
	if err := catalogv1alpha1.DecodeStrict(jsonContent, &digestValue); err != nil {
		return nil, "", err
	}
	catalog.Digest = digestValue.Digest
	return jsonContent, catalogv1alpha1.RenderWorkspaceCheckMarkdown(catalog), nil
}

func run(args []string, stdout io.Writer) error {
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}
	compiler := &compiler{root: opts.workspaceRoot, opts: opts}
	if err := compiler.compile(); err != nil {
		return err
	}
	catalog, err := compiler.catalog()
	if err != nil {
		return err
	}
	jsonContent, markdown, err := compiler.buildOutputs(catalog)
	if err != nil {
		return err
	}
	if opts.check {
		if len(compiler.problems) > 0 {
			return fmt.Errorf("workspace-check catalog is incomplete: %d problem(s)",
				len(compiler.problems))
		}
		jsonPath := filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked workspace-check JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked workspace-check JSON is stale; run //tools/agents/cmd/workspace_check:workspace_check_update")
		}
		markdownPath := filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked workspace-check Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked workspace-check Markdown is stale; run //tools/agents/cmd/workspace_check:workspace_check_update")
		}
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath)), jsonContent, 0o644); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath)), []byte(markdown), 0o644); err != nil {
		return err
	}
	report := struct {
		CatalogID    string   `json:"catalogID"`
		Output       string   `json:"output"`
		Markdown     string   `json:"markdown"`
		Completeness string   `json:"completeness"`
		Status       string   `json:"status"`
		Problems     []string `json:"problems,omitempty"`
	}{
		CatalogID:    catalog.ID,
		Output:       opts.outputPath,
		Markdown:     opts.markdownPath,
		Completeness: string(catalog.Completeness),
		Status:       "ok",
		Problems:     compiler.problems,
	}
	if len(compiler.problems) > 0 {
		report.Status = "incomplete"
	}
	content, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	if _, err := stdout.Write(content); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "workspace_check:", err)
		os.Exit(1)
	}
}
