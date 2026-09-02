// Command action_check compiles the bounded deterministic ActionCatalog over
// the registry-declared operation files and their owner-local definitions.
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

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

const (
	registrySchema   = "agents.alwaldend.com/phase1-registry/v1alpha1"
	operationsSchema = "agents.alwaldend.com/operations/v1alpha1"
)

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
	flags := flag.NewFlagSet("action_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.registryPath, "registry", "tools/agents/declarations/registry.json",
		"registry JSON path relative to workspace root")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/action.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/action.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on completeness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: content-addressed inputs)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.action-compiler",
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
	Schema         string   `json:"schema"`
	OperationFiles []string `json:"operationFiles"`
}

type operationFile struct {
	Schema     string           `json:"schema"`
	Owner      string           `json:"owner"`
	Provider   string           `json:"provider"`
	Definition string           `json:"definition"`
	Operations []operationEntry `json:"operations"`
}

type operationEntry struct {
	ID                  string   `json:"id"`
	Selector            string   `json:"selector"`
	Classification      string   `json:"classification"`
	Effects             []string `json:"effects"`
	Inputs              []string `json:"inputs"`
	Outputs             []string `json:"outputs"`
	Information         []string `json:"information"`
	CredentialUse       string   `json:"credentialUse"`
	NetworkUse          string   `json:"networkUse"`
	EnvironmentSelector string   `json:"environmentSelector"`
	AuthorityGate       string   `json:"authorityGate"`
	Preflight           string   `json:"preflight"`
	Verification        string   `json:"verification"`
	Cost                string   `json:"cost"`
	Cacheability        string   `json:"cacheability"`
	Cancellation        string   `json:"cancellation"`
}

type compiler struct {
	root        string
	opts        options
	inputs      []catalogv1alpha1.CatalogInput
	providers   []catalogv1alpha1.ActionProvider
	actions     []catalogv1alpha1.ActionRecord
	aliases     []catalogv1alpha1.ActionAlias
	conflicts   []catalogv1alpha1.CatalogConflict
	problems    []string
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

func (c *compiler) compile() error {
	registryPath := filepath.Join(c.root, filepath.FromSlash(c.opts.registryPath))
	registryContent, err := os.ReadFile(registryPath)
	if err != nil {
		return fmt.Errorf("read registry: %w", err)
	}
	var registry registry
	decoder := json.NewDecoder(bytes.NewReader(registryContent))
	if err := decoder.Decode(&registry); err != nil {
		return fmt.Errorf("decode registry: %w", err)
	}
	if registry.Schema != registrySchema {
		return fmt.Errorf("registry schema mismatch: %s", registry.Schema)
	}
	c.input(c.opts.registryPath, "registry", registryContent)
	c.eligible++
	c.emitted++
	for _, operationPath := range registry.OperationFiles {
		c.eligible++
		full := filepath.Join(c.root, filepath.FromSlash(operationPath))
		content, err := os.ReadFile(full)
		if err != nil {
			c.problem("operation file missing: %s", operationPath)
			c.unavailable++
			continue
		}
		c.input(operationPath, "operation-file", content)
		var operations operationFile
		decoder := json.NewDecoder(bytes.NewReader(content))
		if err := decoder.Decode(&operations); err != nil {
			c.problem("operation file %s is malformed: %v", operationPath, err)
			c.unavailable++
			continue
		}
		if operations.Schema != operationsSchema {
			c.problem("operation file %s has schema mismatch: %s",
				operationPath, operations.Schema)
			c.unavailable++
			continue
		}
		c.emitted++
		providerID := operations.Provider
		c.providers = append(c.providers, catalogv1alpha1.ActionProvider{
			ID:             providerID,
			Owner:          operations.Owner,
			DefinitionPath: operations.Definition,
		})
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
		for _, entry := range operations.Operations {
			owner := operations.Owner
			c.actions = append(c.actions, catalogv1alpha1.ActionRecord{
				ID:                  entry.ID,
				ProviderRef:         providerID,
				Owner:               owner,
				SourcePath:          operationPath,
				Selector:            entry.Selector,
				Classification:      entry.Classification,
				Effects:             entry.Effects,
				Inputs:              entry.Inputs,
				Outputs:             entry.Outputs,
				Information:         entry.Information,
				CredentialUse:       entry.CredentialUse,
				NetworkUse:          entry.NetworkUse,
				EnvironmentSelector: entry.EnvironmentSelector,
				AuthorityGate:       entry.AuthorityGate,
				Preflight:           entry.Preflight,
				Verification:        entry.Verification,
				Cost:                entry.Cost,
				Cacheability:        entry.Cacheability,
				Cancellation:        entry.Cancellation,
			})
		}
	}
	return nil
}

func (c *compiler) catalog() (catalogv1alpha1.ActionCatalog, error) {
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
	providers := c.providers
	if providers == nil {
		providers = []catalogv1alpha1.ActionProvider{}
	}
	actions := c.actions
	if actions == nil {
		actions = []catalogv1alpha1.ActionRecord{}
	}
	aliases := c.aliases
	if aliases == nil {
		aliases = []catalogv1alpha1.ActionAlias{}
	}
	sort.Slice(providers, func(i, j int) bool { return providers[i].ID < providers[j].ID })
	sort.Slice(actions, func(i, j int) bool { return actions[i].ID < actions[j].ID })
	return catalogv1alpha1.ActionCatalog{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            catalogv1alpha1.APIVersion + "/" + catalogv1alpha1.KindActionCatalog,
			Kind:              catalogv1alpha1.KindActionCatalog,
			ID:                "agent-system.action",
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
		Providers: providers,
		Actions:   actions,
		Aliases:   aliases,
	}, nil
}

func (c *compiler) buildOutputs(catalog catalogv1alpha1.ActionCatalog) ([]byte, string, error) {
	jsonContent, err := catalogv1alpha1.CanonicalJSONAction(catalog)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	var digestValue catalogv1alpha1.ActionCatalog
	if err := catalogv1alpha1.DecodeStrict(jsonContent, &digestValue); err != nil {
		return nil, "", err
	}
	catalog.Digest = digestValue.Digest
	return jsonContent, catalogv1alpha1.RenderActionMarkdown(catalog), nil
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
			return fmt.Errorf("action catalog is incomplete: %d problem(s)",
				len(compiler.problems))
		}
		jsonPath := filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked action JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked action JSON is stale; rerun the generator")
		}
		markdownPath := filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked action Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked action Markdown is stale; rerun the generator")
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
		fmt.Fprintln(os.Stderr, "action_check:", err)
		os.Exit(1)
	}
}
