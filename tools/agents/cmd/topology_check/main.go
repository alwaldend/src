// Command topology_check compiles the bounded deterministic TopologyCatalog
// over the registered projects universe and tracked Bzlmod workspaces.
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

const registrySchema = "agents.alwaldend.com/phase1-registry/v1alpha1"

var boundaryNames = map[string]string{
	"projects":    "product",
	"infra":       "repository_internal",
	"tools":       "tool",
	"data":        "data",
	"third_party": "third_party",
	"users":       "user",
}

type options struct {
	workspaceRoot  string
	registryPath   string
	outputPath     string
	markdownPath   string
	check          bool
	sourceRevision string
	producerRef    string
}

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("topology_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.registryPath, "registry", "tools/agents/declarations/registry.json",
		"registry JSON path relative to workspace root")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/topology.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/topology.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on completeness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: git rev-parse HEAD)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.topology-compiler",
		"producer reference for the catalog")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	return opts, nil
}

type registry struct {
	Schema           string   `json:"schema"`
	CriteriaRevision int      `json:"criteriaRevision"`
	OperationFiles   []string `json:"operationFiles"`
}

type moduleManifest struct {
	moduleName string
	modulePath string
}

var moduleNamePattern = regexp.MustCompile(`(?m)^\s*(?:name|module\(\s*name)\s*=\s*"([^"]+)"`)

func readModuleName(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	match := moduleNamePattern.FindSubmatch(content)
	if match == nil {
		return "", fmt.Errorf("no module name in %s", path)
	}
	return string(match[1]), nil
}

func hashContent(content []byte) string {
	value := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(value[:])
}

func (compiler *compiler) input(path, role string, content []byte) {
	relative, err := filepath.Rel(compiler.root, path)
	if err != nil {
		relative = path
	}
	compiler.inputs = append(compiler.inputs, catalogv1alpha1.CatalogInput{
		Path:   filepath.ToSlash(relative),
		Role:   role,
		Digest: hashContent(content),
	})
}

func (compiler *compiler) inside(root, path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute path is forbidden: %s", path)
	}
	joined := filepath.Clean(filepath.Join(root, path))
	relative, err := filepath.Rel(root, joined)
	if err != nil || relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes workspace: %s", path)
	}
	return joined, nil
}

type compiler struct {
	root        string
	opts        options
	registry    registry
	inputs      []catalogv1alpha1.CatalogInput
	trees       []catalogv1alpha1.TopologyTree
	components  []catalogv1alpha1.TopologyComponent
	workspaces  []catalogv1alpha1.TopologyWorkspace
	problems    []string
	conflicts   []catalogv1alpha1.CatalogConflict
	eligible    int
	emitted     int
	unavailable int
}

func (compiler *compiler) problem(format string, args ...any) {
	compiler.problems = append(compiler.problems, fmt.Sprintf(format, args...))
}

func (compiler *compiler) loadRegistry(registryPath string) error {
	path, err := compiler.inside(compiler.root, registryPath)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read registry: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&compiler.registry); err != nil {
		return fmt.Errorf("decode registry: %w", err)
	}
	if compiler.registry.Schema != registrySchema {
		return fmt.Errorf("registry schema mismatch: %s", compiler.registry.Schema)
	}
	compiler.input(registryPath, "registry", content)
	return nil
}

func (compiler *compiler) compile() error {
	if err := compiler.loadRegistry(compiler.opts.registryPath); err != nil {
		return err
	}
	// Top-level boundary trees with their authoritative READMEs.
	boundaryReadmes := map[string]string{
		"projects":    "projects/README.md",
		"infra":       "infra/README.md",
		"tools":       "tools/README.md",
		"data":        "data/README.md",
		"third_party": "third_party/README.md",
		"users":       "users/README.md",
	}
	boundaryOrder := []string{"projects", "infra", "tools", "data", "third_party", "users"}
	for _, name := range boundaryOrder {
		readme := boundaryReadmes[name]
		path, err := compiler.inside(compiler.root, readme)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			compiler.problem("tree README missing: %s", readme)
			compiler.eligible++
			compiler.unavailable++
			continue
		}
		compiler.input(readme, "tree-boundary-readme", content)
		compiler.trees = append(compiler.trees, catalogv1alpha1.TopologyTree{
			ID:         name,
			Path:       name,
			ReadmePath: readme,
			Boundary:   catalogv1alpha1.TreeBoundaryClass(boundaryNames[name]),
		})
		compiler.eligible++
		compiler.emitted++
	}
	// Registered projects universe: immediate projects/* directories.
	projectsDir := filepath.Join(compiler.root, "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return fmt.Errorf("read projects dir: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == ".git" || strings.HasPrefix(name, ".") {
			continue
		}
		compiler.eligible++
		component, ok := compiler.readComponent(name)
		if !ok {
			compiler.unavailable++
			continue
		}
		compiler.components = append(compiler.components, component)
		compiler.emitted++
	}
	// Tracked workspaces: MODULE.bazel roots.
	modules, err := walkModuleRoots(compiler.root)
	if err != nil {
		return err
	}
	for _, module := range modules {
		compiler.eligible++
		workspace, ok := compiler.readWorkspace(module)
		if !ok {
			compiler.unavailable++
			continue
		}
		compiler.workspaces = append(compiler.workspaces, workspace)
		compiler.emitted++
	}
	return nil
}

func (compiler *compiler) readComponent(name string) (catalogv1alpha1.TopologyComponent, bool) {
	componentDir := "projects/" + name
	readme, err := compiler.inside(compiler.root, componentDir+"/README.md")
	if err != nil {
		compiler.problem("project path escapes: %s", componentDir)
		return catalogv1alpha1.TopologyComponent{}, false
	}
	readmeContent, err := os.ReadFile(readme)
	if err != nil {
		compiler.problem("project README missing: %s", componentDir)
		return catalogv1alpha1.TopologyComponent{}, false
	}
	build, err := compiler.inside(compiler.root, componentDir+"/BUILD.bazel")
	if err != nil {
		compiler.problem("project build path escapes: %s", componentDir)
		return catalogv1alpha1.TopologyComponent{}, false
	}
	buildContent, err := os.ReadFile(build)
	if err != nil {
		compiler.problem("project BUILD missing: %s", componentDir)
		return catalogv1alpha1.TopologyComponent{}, false
	}
	title, description, lifecycle, err := parseProjectREADME(name, readmeContent)
	if err != nil {
		compiler.problem("project README malformed: %s: %v", componentDir, err)
		return catalogv1alpha1.TopologyComponent{}, false
	}
	compiler.input(readme, "component-readme", readmeContent)
	compiler.input(build, "component-build", buildContent)
	return catalogv1alpha1.TopologyComponent{
		ID:          name,
		Path:        componentDir,
		OwnerReadme: componentDir + "/README.md",
		BuildPath:   componentDir + "/BUILD.bazel",
		Title:       title,
		Description: description,
		Lifecycle:   lifecycle,
		DocsState:   "owned",
	}, true
}

func parseProjectREADME(name string, content []byte) (title, description, lifecycle string, err error) {
	inFrontmatter := false
	closedFrontmatter := false
	inStatuses := false
	for _, rawLine := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(rawLine)
		if strings.HasPrefix(line, "---") {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			inFrontmatter = false
			closedFrontmatter = true
			continue
		}
		if inFrontmatter {
			switch {
			case strings.HasPrefix(line, "title:"):
				title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
				inStatuses = false
			case strings.HasPrefix(line, "description:"):
				description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
				inStatuses = false
			case line == "statuses:":
				inStatuses = true
			case inStatuses && strings.HasPrefix(line, "-") && lifecycle == "":
				lifecycle = strings.TrimSpace(strings.TrimPrefix(line, "-"))
			default:
				inStatuses = false
			}
			continue
		}
		if closedFrontmatter && strings.HasPrefix(line, "# ") && title == "" {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	if title == "" {
		title = name
	}
	if description == "" {
		description = name
	}
	if lifecycle == "" {
		return "", "", "", fmt.Errorf("README %s lacks a statuses frontmatter entry", name)
	}
	return title, description, lifecycle, nil
}

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

func (compiler *compiler) readWorkspace(modulePath string) (catalogv1alpha1.TopologyWorkspace, bool) {
	full, err := compiler.inside(compiler.root, modulePath)
	if err != nil {
		compiler.problem("workspace path escapes: %s", modulePath)
		return catalogv1alpha1.TopologyWorkspace{}, false
	}
	content, err := os.ReadFile(full)
	if err != nil {
		compiler.problem("workspace module unreadable: %s", modulePath)
		return catalogv1alpha1.TopologyWorkspace{}, false
	}
	moduleName, err := readModuleName(full)
	if err != nil {
		compiler.problem("workspace module name missing: %s", modulePath)
		return catalogv1alpha1.TopologyWorkspace{}, false
	}
	compiler.input(modulePath, "workspace-module", content)
	return catalogv1alpha1.TopologyWorkspace{
		ID:         workspaceID(modulePath),
		Path:       filepath.Dir(modulePath),
		ModulePath: modulePath,
		ModuleName: moduleName,
	}, true
}

func workspaceID(modulePath string) string {
	dir := filepath.ToSlash(filepath.Dir(modulePath))
	if dir == "." {
		return "root"
	}
	return strings.ReplaceAll(dir, "/", ".")
}

func (compiler *compiler) catalog(opts options) (catalogv1alpha1.TopologyCatalog, error) {
	sourceRevision := opts.sourceRevision
	if sourceRevision == "" {
		var aggregate bytes.Buffer
		for _, input := range compiler.inputs {
			aggregate.WriteString(input.Path)
			aggregate.WriteString("\x00")
			aggregate.WriteString(input.Role)
			aggregate.WriteString("\x00")
			aggregate.WriteString(input.Digest)
			aggregate.WriteString("\n")
		}
		// The source revision is a short content-addressed identity over the
		// input universe. It changes only when the owner-local facts change,
		// not on unrelated commits.
		sourceRevision = shortenDigest(hashContent(aggregate.Bytes()))
	}
	completeness := catalogv1alpha1.CompletenessComplete
	limitations := []string{}
	if len(compiler.problems) > 0 {
		completeness = catalogv1alpha1.CompletenessPartial
		limitations = compiler.problems
	}
	inputs := compiler.inputs
	if inputs == nil {
		inputs = []catalogv1alpha1.CatalogInput{}
	}
	sort.Slice(inputs, func(i, j int) bool {
		if inputs[i].Path != inputs[j].Path {
			return inputs[i].Path < inputs[j].Path
		}
		return inputs[i].Role < inputs[j].Role
	})
	conflicts := compiler.conflicts
	if conflicts == nil {
		conflicts = []catalogv1alpha1.CatalogConflict{}
	}
	trees := compiler.trees
	if trees == nil {
		trees = []catalogv1alpha1.TopologyTree{}
	}
	components := compiler.components
	if components == nil {
		components = []catalogv1alpha1.TopologyComponent{}
	}
	workspaces := compiler.workspaces
	if workspaces == nil {
		workspaces = []catalogv1alpha1.TopologyWorkspace{}
	}
	return catalogv1alpha1.TopologyCatalog{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            catalogv1alpha1.APIVersion + "/" + catalogv1alpha1.KindTopologyCatalog,
			Kind:              catalogv1alpha1.KindTopologyCatalog,
			ID:                "agent-system.topology",
			DerivationVersion: "1.0.0",
			ProducerRef:       opts.producerRef,
			SourceRevision:    sourceRevision,
			Inputs:            inputs,
			Bounds: catalogv1alpha1.CatalogBounds{
				Eligible:       compiler.eligible,
				Emitted:        compiler.emitted,
				Unavailable:    compiler.unavailable,
				MaxItems:       1000,
				MaxInputBytes:  32 << 20,
				MaxOutputBytes: 1 << 20,
			},
			Completeness: completeness,
			Limitations:  limitations,
			Conflicts:    conflicts,
		},
		Trees:      trees,
		Components: components,
		Workspaces: workspaces,
	}, nil
}

func (compiler *compiler) buildOutputs(catalog catalogv1alpha1.TopologyCatalog) ([]byte, string, error) {
	jsonContent, err := catalogv1alpha1.CanonicalJSONTopology(catalog)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	var digestValue catalogv1alpha1.TopologyCatalog
	if err := catalogv1alpha1.DecodeStrict(jsonContent, &digestValue); err != nil {
		return nil, "", err
	}
	catalog.Digest = digestValue.Digest
	return jsonContent, catalogv1alpha1.RenderTopologyMarkdown(catalog), nil
}

func (compiler *compiler) writeArtifacts(opts options, jsonContent []byte, markdown string) error {
	jsonPath, err := compiler.inside(compiler.root, opts.outputPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, jsonContent, 0o644); err != nil {
		return err
	}
	markdownPath, err := compiler.inside(compiler.root, opts.markdownPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(markdownPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(markdownPath, []byte(markdown), 0o644)
}

func run(args []string, stdout io.Writer, stderr io.Writer) error {
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}
	compiler := &compiler{root: opts.workspaceRoot, opts: opts}
	if err := compiler.compile(); err != nil {
		return err
	}
	catalog, err := compiler.catalog(opts)
	if err != nil {
		return err
	}
	jsonContent, markdown, err := compiler.buildOutputs(catalog)
	if err != nil {
		return err
	}
	status := "ok"
	if len(compiler.problems) > 0 {
		status = "incomplete"
	}
	if opts.check {
		if len(compiler.problems) > 0 {
			return fmt.Errorf("topology catalog is incomplete: %d problem(s)",
				len(compiler.problems))
		}
		jsonPath, err := compiler.inside(compiler.root, opts.outputPath)
		if err != nil {
			return err
		}
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked topology JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked topology JSON is stale; rerun the generator")
		}
		markdownPath, err := compiler.inside(compiler.root, opts.markdownPath)
		if err != nil {
			return err
		}
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked topology Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked topology Markdown is stale; rerun the generator")
		}
	}
	if err := compiler.writeArtifacts(opts, jsonContent, markdown); err != nil {
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
		Status:       status,
		Problems:     compiler.problems,
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
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "topology_check:", err)
		os.Exit(1)
	}
}

// shortenDigest returns the first 16 hex characters of a sha256 digest,
// suitable for a short deterministic source-revision identity.
func shortenDigest(value string) string {
	const prefix = "sha256:"
	if len(value) > len(prefix)+16 {
		return value[len(prefix) : len(prefix)+16]
	}
	return value
}
