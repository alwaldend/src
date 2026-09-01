// Command capability_check compiles the bounded deterministic
// CapabilityCatalog over the registry-declared skills, runtime tools, direct
// binaries, operation providers, and their discovery links.
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
	"os"
	"path/filepath"
	"sort"
	"strings"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

const registrySchema = "agents.alwaldend.com/phase1-registry/v1alpha1"

const discoveryDir = ".agents/skills"
const skillsDir = "projects/agents/skills"
const runtimeToolSource = "projects/mcp_cordis/internal/mcp.mjs"
const fullCheckSkillPath = "projects/agents/skills/full-repo-check"

// discoveryOnlyLimitation is emitted once whenever a discovered skill has no
// registry entry. The discovery set is intentionally a superset of the
// registry set, so the emitted catalog stays complete by design while the
// limitation keeps the identity difference explicit.
const discoveryOnlyLimitation = "discovered skill has no registry entry; " +
	"emitted as discovery-only with explicit identity"

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
	flags := flag.NewFlagSet("capability_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.registryPath, "registry", "tools/agents/declarations/registry.json",
		"registry JSON path relative to workspace root")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/capability.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/capability.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on completeness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: content-addressed inputs)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.capability-compiler",
		"producer reference for the catalog")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	return opts, nil
}

type registrySkill struct {
	ID                   string   `json:"id"`
	Owner                string   `json:"owner"`
	Layer                string   `json:"layer"`
	Activation           string   `json:"activation"`
	Exclusions           []string `json:"exclusions"`
	CapabilityRefs       []string `json:"capabilityRefs"`
	Dependencies         []string `json:"dependencies"`
	Conflicts            []string `json:"conflicts"`
	ProviderRequirements []string `json:"providerRequirements"`
	ContextCost          string   `json:"contextCost"`
	EvaluationMaturity   string   `json:"evaluationMaturity"`
}

type runtimeTool struct {
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Classification string `json:"classification"`
}

type directBinary struct {
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Path           string `json:"path"`
	Classification string `json:"classification"`
}

type operationProviderFile struct {
	Schema     string `json:"schema"`
	Owner      string `json:"owner"`
	Provider   string `json:"provider"`
	Definition string `json:"definition"`
}

type registry struct {
	Schema             string                  `json:"schema"`
	Skills             []registrySkill         `json:"skills"`
	RuntimeTools       []runtimeTool           `json:"runtimeTools"`
	DirectBinaries     []directBinary          `json:"directBinaries"`
	OperationFiles     []string                `json:"operationFiles"`
	Authorities        []registryAuthority     `json:"authorities"`
	OperationProviders []operationProviderFile `json:"operationProviders"`
}

type registryAuthority struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Source string `json:"source"`
}

type compiler struct {
	root          string
	opts          options
	inputs        []catalogv1alpha1.CatalogInput
	skills        []catalogv1alpha1.CapabilitySkill
	providers     []catalogv1alpha1.CapabilityProvider
	conflicts     []catalogv1alpha1.CatalogConflict
	problems      []string
	discoveryOnly bool
	eligible      int
	emitted       int
	unavailable   int
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

func (c *compiler) loadRegistry(registryPath string) (*registry, []byte, error) {
	full := filepath.Join(c.root, filepath.FromSlash(registryPath))
	content, err := os.ReadFile(full)
	if err != nil {
		return nil, nil, fmt.Errorf("read registry: %w", err)
	}
	var registry registry
	decoder := json.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&registry); err != nil {
		return nil, nil, fmt.Errorf("decode registry: %w", err)
	}
	if registry.Schema != registrySchema {
		return nil, nil, fmt.Errorf("registry schema mismatch: %s", registry.Schema)
	}
	return &registry, content, nil
}

func (c *compiler) compile() error {
	registry, registryContent, err := c.loadRegistry(c.opts.registryPath)
	if err != nil {
		return err
	}
	c.input(c.opts.registryPath, "registry", registryContent)
	c.eligible++
	c.emitted++

	discoveryLinkIDs, err := c.readDiscoveryLinks()
	if err != nil {
		return err
	}

	registryByID := map[string]registrySkill{}
	for _, skill := range registry.Skills {
		registryByID[skill.ID] = skill
	}

	discovered := map[string]bool{}
	for _, id := range discoveryLinkIDs {
		discovered[id] = true
	}

	emittedSkills := map[string]bool{}
	for _, id := range discoveryLinkIDs {
		entry, ok := registryByID[id]
		if !ok {
			if !c.discoveryOnly {
				c.discoveryOnly = true
				c.problems = append(c.problems, discoveryOnlyLimitation)
			}
			c.eligible++
			c.emitted++
			skill, ok := c.readDiscoveredOnlySkill(id)
			if !ok {
				c.unavailable++
				continue
			}
			c.skills = append(c.skills, skill)
			emittedSkills[id] = true
			continue
		}
		c.eligible++
		skill, ok := c.readRegisteredSkill(entry)
		if !ok {
			c.unavailable++
			continue
		}
		c.skills = append(c.skills, skill)
		emittedSkills[id] = true
	}

	for _, entry := range registry.Skills {
		if discovered[entry.ID] || emittedSkills[entry.ID] {
			continue
		}
		c.eligible++
		skill, ok := c.readRegisteredSkill(entry)
		if !ok {
			c.unavailable++
			continue
		}
		c.skills = append(c.skills, skill)
		emittedSkills[entry.ID] = true
	}

	// Runtime tools: the shared MCP server source is a single input.
	runtimePath := runtimeToolSource
	fullRuntime, err := filepath.Abs(filepath.Join(c.root, filepath.FromSlash(runtimePath)))
	if err != nil {
		return err
	}
	relativeRuntime, err := filepath.Rel(c.root, fullRuntime)
	if err != nil || strings.HasPrefix(relativeRuntime, "..") {
		return fmt.Errorf("runtime tool source escapes workspace: %s", runtimePath)
	}
	runtimeContent, err := os.ReadFile(fullRuntime)
	runtimeRead := false
	if err != nil {
		c.problem("runtime tool source missing: %s", runtimePath)
		c.eligible++
		c.unavailable++
	} else {
		c.input(filepath.ToSlash(relativeRuntime), "runtime-tool-source", runtimeContent)
		c.eligible++
		c.emitted++
		runtimeRead = true
	}
	seenProviders := map[string]bool{}
	for _, tool := range registry.RuntimeTools {
		provider := catalogv1alpha1.CapabilityProvider{
			ID:             tool.ID,
			Owner:          tool.Owner,
			Kind:           "runtime_tool",
			SourcePath:     runtimePath,
			Classification: tool.Classification,
		}
		c.providers = append(c.providers, provider)
		seenProviders[provider.ID] = true
		if runtimeRead {
			c.emitted++
		}
	}

	// Direct binaries: each source path is read individually.
	for _, binary := range registry.DirectBinaries {
		provider := catalogv1alpha1.CapabilityProvider{
			ID:             binary.ID,
			Owner:          binary.Owner,
			Kind:           "direct_binary",
			SourcePath:     binary.Path,
			Classification: binary.Classification,
		}
		providerFull := filepath.Join(c.root, filepath.FromSlash(binary.Path))
		c.eligible++
		content, err := os.ReadFile(providerFull)
		if err != nil {
			c.problem("direct binary source missing: %s", binary.Path)
			c.unavailable++
		} else {
			c.input(binary.Path, "direct-binary-source", content)
			c.emitted++
		}
		seenProviders[provider.ID] = true
		c.providers = append(c.providers, provider)
	}

	// Operation providers: one entry per operation file with definition.
	operationOrder := append([]string(nil), registry.OperationFiles...)
	sort.Strings(operationOrder)
	readOperationFiles := map[string]bool{}
	for _, operationPath := range operationOrder {
		c.eligible++
		fullOperation := filepath.Join(c.root, filepath.FromSlash(operationPath))
		content, err := os.ReadFile(fullOperation)
		if err != nil {
			c.problem("operation file missing: %s", operationPath)
			c.unavailable++
			continue
		}
		c.input(operationPath, "operation-file", content)
		var operations operationProviderFile
		decoder := json.NewDecoder(bytes.NewReader(content))
		if err := decoder.Decode(&operations); err != nil {
			c.problem("operation file %s is malformed: %v", operationPath, err)
			c.unavailable++
			continue
		}
		if operations.Schema != "agents.alwaldend.com/operations/v1alpha1" {
			c.problem("operation file %s has schema mismatch: %s",
				operationPath, operations.Schema)
			c.unavailable++
			continue
		}
		readOperationFiles[operationPath] = true
		c.emitted++
		provider := catalogv1alpha1.CapabilityProvider{
			ID:             operations.Provider,
			Owner:          operations.Owner,
			Kind:           "operation_provider",
			SourcePath:     operations.Definition,
			Classification: "classified",
		}
		if operations.Definition != "" {
			fullDef := filepath.Join(c.root, filepath.FromSlash(operations.Definition))
			if definitionContent, err := os.ReadFile(fullDef); err == nil {
				c.input(operations.Definition, "operation-definition", definitionContent)
				c.eligible++
				c.emitted++
			} else {
				c.problem("operation definition missing: %s", operations.Definition)
				c.unavailable++
			}
		}
		if !seenProviders[provider.ID] {
			seenProviders[provider.ID] = true
			c.providers = append(c.providers, provider)
		}
	}
	_ = readOperationFiles

	// Sort and dedup providers by ID, keeping the first occurrence; the
	// registry and operation sets intentionally overlap (for example the
	// bazel-agent binary and the bazel-agent skill are distinct identities,
	// but a direct binary and an operation provider can share an ID).
	seen := map[string]bool{}
	deduped := c.providers[:0]
	for _, provider := range c.providers {
		if seen[provider.ID] {
			continue
		}
		seen[provider.ID] = true
		deduped = append(deduped, provider)
	}
	c.providers = deduped
	return nil
}

func (c *compiler) readDiscoveryLinks() ([]string, error) {
	discoveryRoot := filepath.Join(c.root, filepath.FromSlash(discoveryDir))
	entries, err := os.ReadDir(discoveryRoot)
	if err != nil {
		if os.IsNotExist(err) {
			c.problem("discovery directory missing: %s", discoveryDir)
			return nil, nil
		}
		return nil, fmt.Errorf("read discovery dir: %w", err)
	}
	var ids []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		id := entry.Name()
		if id == "" || strings.ContainsAny(id, "/\\") {
			continue
		}
		linkPath := filepath.ToSlash(discoveryDir + "/" + id)
		linkFull := filepath.Join(c.root, filepath.FromSlash(linkPath))
		if _, err := os.Lstat(linkFull); err != nil {
			continue
		}
		if !strings.HasPrefix(id, ".") {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids, nil
}

func (c *compiler) readDiscoveredOnlySkill(id string) (catalogv1alpha1.CapabilitySkill, bool) {
	discoveryPath := discoveryDir + "/" + id
	canonicalPath := canonicalSkillPath(discoveryPath, c.root)
	if canonicalPath == "" {
		c.problem("skill %s discovery link is not a relative symlink", id)
		return catalogv1alpha1.CapabilitySkill{}, false
	}
	if err := c.readSkillDocuments(id, discoveryPath, canonicalPath); err != nil {
		c.problem("skill %s documents unavailable: %v", id, err)
		return catalogv1alpha1.CapabilitySkill{}, false
	}
	return catalogv1alpha1.CapabilitySkill{
		ID:                   id,
		Owner:                "projects/agents",
		CanonicalPath:        canonicalPath,
		DiscoveryPath:        discoveryPath,
		Layer:                "procedure",
		Activation:           "discovered skill for repository agents",
		Exclusions:           []string{},
		CapabilityRefs:       []string{},
		Dependencies:         []string{},
		Conflicts:            []string{},
		ProviderRequirements: []string{},
		ContextCost:          "medium",
		EvaluationMaturity:   "discovered",
	}, true
}

func canonicalSkillPath(discoveryPath, root string) string {
	linkFull := filepath.Join(root, filepath.FromSlash(discoveryPath))
	target, err := os.Readlink(linkFull)
	if err != nil {
		return ""
	}
	if filepath.IsAbs(target) {
		return ""
	}
	canonicalFull := filepath.Clean(filepath.Join(filepath.Dir(linkFull), target))
	relative, err := filepath.Rel(root, canonicalFull)
	if err != nil || relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return ""
	}
	if _, err := os.Stat(canonicalFull); err != nil {
		return ""
	}
	return filepath.ToSlash(relative)
}

func (c *compiler) readRegisteredSkill(entry registrySkill) (catalogv1alpha1.CapabilitySkill, bool) {
	id := entry.ID
	discoveryPath := discoveryDir + "/" + id
	linkFull := filepath.Join(c.root, filepath.FromSlash(discoveryPath))
	if _, err := os.Lstat(linkFull); err != nil {
		c.problem("skill discovery link missing: %s", id)
		return catalogv1alpha1.CapabilitySkill{}, false
	}
	canonicalPath := canonicalSkillPath(discoveryPath, c.root)
	if canonicalPath == "" {
		c.problem("skill %s discovery link is not a relative symlink", id)
		return catalogv1alpha1.CapabilitySkill{}, false
	}
	if err := c.readSkillDocuments(id, discoveryPath, canonicalPath); err != nil {
		c.problem("skill %s documents unavailable: %v", id, err)
		return catalogv1alpha1.CapabilitySkill{}, false
	}
	return catalogv1alpha1.CapabilitySkill{
		ID:                   id,
		Owner:                entry.Owner,
		CanonicalPath:        canonicalPath,
		DiscoveryPath:        discoveryPath,
		Layer:                entry.Layer,
		Activation:           entry.Activation,
		Exclusions:           catalogv1alpha1.SortedUnique(entry.Exclusions),
		CapabilityRefs:       catalogv1alpha1.SortedUnique(entry.CapabilityRefs),
		Dependencies:         catalogv1alpha1.SortedUnique(entry.Dependencies),
		Conflicts:            catalogv1alpha1.SortedUnique(entry.Conflicts),
		ProviderRequirements: catalogv1alpha1.SortedUnique(entry.ProviderRequirements),
		ContextCost:          entry.ContextCost,
		EvaluationMaturity:   entry.EvaluationMaturity,
	}, true
}

func (c *compiler) readSkillDocuments(id, discoveryPath, canonicalPath string) error {
	skillDoc := canonicalPath + "/SKILL.md"
	fullDoc := filepath.Join(c.root, filepath.FromSlash(skillDoc))
	content, err := os.ReadFile(fullDoc)
	if err != nil {
		return fmt.Errorf("SKILL.md: %w", err)
	}
	c.input(skillDoc, "skill-doc", content)

	buildPath := canonicalPath + "/BUILD.bazel"
	fullBuild := filepath.Join(c.root, filepath.FromSlash(buildPath))
	buildContent, err := os.ReadFile(fullBuild)
	if err != nil {
		return fmt.Errorf("BUILD.bazel: %w", err)
	}
	c.input(buildPath, "skill-build", buildContent)

	target, err := os.Readlink(filepath.Join(c.root, filepath.FromSlash(discoveryPath)))
	if err != nil {
		return fmt.Errorf("discovery link: %w", err)
	}
	c.input(discoveryPath, "skill-discovery", []byte(target))
	return nil
}

// orderSkills returns a deterministic dependency-respecting order: each
// skill appears after all skills it depends on, and otherwise skills are
// ordered by ID. The shared schema validates dependencies as
// already-seen, so the emitted order must be topological.
func orderSkills(skills []catalogv1alpha1.CapabilitySkill) []catalogv1alpha1.CapabilitySkill {
	byID := map[string]catalogv1alpha1.CapabilitySkill{}
	dependents := map[string][]string{}
	remainingDeps := map[string]int{}
	for _, skill := range skills {
		byID[skill.ID] = skill
		remainingDeps[skill.ID] = 0
	}
	for _, skill := range skills {
		for _, dep := range skill.Dependencies {
			if dep == skill.ID {
				continue
			}
			if _, ok := byID[dep]; ok {
				remainingDeps[skill.ID]++
				dependents[dep] = append(dependents[dep], skill.ID)
			}
		}
	}
	var ordered []catalogv1alpha1.CapabilitySkill
	emitted := map[string]bool{}
	for len(ordered) < len(skills) {
		next := ""
		for _, skill := range skills {
			if emitted[skill.ID] || remainingDeps[skill.ID] > 0 {
				continue
			}
			if next == "" || skill.ID < next {
				next = skill.ID
			}
		}
		if next == "" {
			break
		}
		ordered = append(ordered, byID[next])
		emitted[next] = true
		for _, dependent := range dependents[next] {
			remainingDeps[dependent]--
		}
	}
	return ordered
}

func (c *compiler) catalog() (catalogv1alpha1.CapabilityCatalog, error) {
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
		limitations = append([]string(nil), c.problems...)
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
	skills := c.skills
	if skills == nil {
		skills = []catalogv1alpha1.CapabilitySkill{}
	}
	providers := c.providers
	if providers == nil {
		providers = []catalogv1alpha1.CapabilityProvider{}
	}
	skills = orderSkills(skills)
	if skills == nil {
		skills = []catalogv1alpha1.CapabilitySkill{}
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i].ID < providers[j].ID })
	if providers == nil {
		providers = []catalogv1alpha1.CapabilityProvider{}
	}
	return catalogv1alpha1.CapabilityCatalog{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            catalogv1alpha1.APIVersion + "/" + catalogv1alpha1.KindCapabilityCatalog,
			Kind:              catalogv1alpha1.KindCapabilityCatalog,
			ID:                "agent-system.capability",
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
		Skills:    skills,
		Providers: providers,
	}, nil
}

func (c *compiler) buildOutputs(catalog catalogv1alpha1.CapabilityCatalog) ([]byte, string, error) {
	jsonContent, err := catalogv1alpha1.CanonicalJSONCapability(catalog)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	var digestValue catalogv1alpha1.CapabilityCatalog
	if err := catalogv1alpha1.DecodeStrict(jsonContent, &digestValue); err != nil {
		return nil, "", err
	}
	catalog.Digest = digestValue.Digest
	return jsonContent, catalogv1alpha1.RenderCapabilityMarkdown(catalog), nil
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
		jsonPath := filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked capability JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked capability JSON is stale; rerun the generator")
		}
		markdownPath := filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked capability Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked capability Markdown is stale; rerun the generator")
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
		fmt.Fprintln(os.Stderr, "capability_check:", err)
		os.Exit(1)
	}
}
