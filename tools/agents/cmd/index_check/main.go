// Command index_check compiles the bounded deterministic AgentSystemIndex
// over the checked agent-system catalog descriptors. The index is an
// inventory only: it references catalog digests and never embeds catalog
// payloads.
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

type options struct {
	workspaceRoot  string
	outputPath     string
	markdownPath   string
	check          bool
	sourceRevision string
	producerRef    string
}

// catalogFile is one checked catalog descriptor input.
type catalogFile struct {
	kind string
	id   string
	path string
}

var indexCatalogFiles = []catalogFile{
	{
		kind: catalogv1alpha1.KindTopologyCatalog,
		id:   "agent-system.topology",
		path: "tools/agents/catalogs/topology.json",
	},
	{
		kind: catalogv1alpha1.KindPolicyCatalog,
		id:   "agent-system.policy",
		path: "tools/agents/catalogs/policy.json",
	},
	{
		kind: catalogv1alpha1.KindActionCatalog,
		id:   "agent-system.action",
		path: "tools/agents/catalogs/action.json",
	},
	{
		kind: catalogv1alpha1.KindCapabilityCatalog,
		id:   "agent-system.capability",
		path: "tools/agents/catalogs/capability.json",
	},
	{
		kind: catalogv1alpha1.KindWorkspaceCheckCatalog,
		id:   "agent-system.workspace-check",
		path: "tools/agents/catalogs/workspace-check.json",
	},
	{
		kind: catalogv1alpha1.KindGoalCatalog,
		id:   "agent-system.goal",
		path: "tools/agents/catalogs/goal.json",
	},
}

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("index_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/index.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/index.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on incompleteness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: content-addressed inputs)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.index-compiler",
		"producer reference for the index")
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
	descriptors []catalogv1alpha1.IndexCatalogDescriptor
	conflicts   []catalogv1alpha1.CatalogConflict
	limitations []string
	eligible    int
	emitted     int
	unavailable int
}

func (c *compiler) input(relativePath, role string, content []byte) {
	value := sha256.Sum256(content)
	c.inputs = append(c.inputs, catalogv1alpha1.CatalogInput{
		Path:   filepath.ToSlash(relativePath),
		Role:   role,
		Digest: "sha256:" + hex.EncodeToString(value[:]),
	})
}

// compile reads the six checked catalog descriptor files. Missing or
// undecodable catalogs become truthful unavailable descriptors with a stable
// limitation; a present catalog never embeds its payload.
func (c *compiler) compile() error {
	for _, file := range indexCatalogFiles {
		c.eligible++
		content, err := os.ReadFile(filepath.Join(c.root, filepath.FromSlash(file.path)))
		if err != nil {
			c.descriptors = append(c.descriptors, unavailableDescriptor(
				file.kind, file.id, "file missing",
			))
			c.unavailable++
			continue
		}
		c.input(file.path, "catalog-descriptor", content)
		descriptor, childConflicts, ok := c.deriveDescriptor(file, content)
		c.descriptors = append(c.descriptors, descriptor)
		c.conflicts = append(c.conflicts, childConflicts...)
		if ok {
			c.emitted++
		} else {
			c.unavailable++
		}
	}
	return nil
}

// unavailableDescriptor is a stable descriptor for a catalog that is not
// present or not decodable.
func unavailableDescriptor(kind, id, reason string) catalogv1alpha1.IndexCatalogDescriptor {
	return catalogv1alpha1.IndexCatalogDescriptor{
		ID:                id,
		Kind:              kind,
		Schema:            catalogv1alpha1.APIVersion + "/" + kind,
		DerivationVersion: "0.0.0",
		Digest:            "sha256:" + strings.Repeat("0", 64),
		InputDigests:      []string{"sha256:" + strings.Repeat("0", 64)},
		Completeness:      "unavailable",
		QueryRoutes: []string{
			"compile:" + kind,
			"check:" + kind,
		},
		Limitations: []string{reason},
	}
}

// envelope carries the shared catalog fields the index projects.
type envelope struct {
	Schema            string                            `json:"schema"`
	Kind              string                            `json:"kind"`
	ID                string                            `json:"id"`
	DerivationVersion string                            `json:"derivationVersion"`
	Digest            string                            `json:"digest"`
	Inputs            []catalogv1alpha1.CatalogInput    `json:"inputs"`
	Completeness      catalogv1alpha1.Completeness      `json:"completeness"`
	Limitations       []string                          `json:"limitations"`
	Conflicts         []catalogv1alpha1.CatalogConflict `json:"conflicts"`
	Bounds            catalogv1alpha1.CatalogBounds     `json:"bounds"`
	ProducerRef       string                            `json:"producerRef"`
	SourceRevision    string                            `json:"sourceRevision"`
}

func (c *compiler) deriveDescriptor(file catalogFile, content []byte) (
	catalogv1alpha1.IndexCatalogDescriptor,
	[]catalogv1alpha1.CatalogConflict,
	bool,
) {
	var envelope envelope
	decoder := json.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&envelope); err != nil {
		return unavailableDescriptor(file.kind, file.id, "undecodable envelope"),
			nil, false
	}
	if envelope.Kind != file.kind || envelope.ID != file.id {
		return unavailableDescriptor(file.kind, file.id, "identity mismatch"),
			nil, false
	}
	if !validDigest(envelope.Digest) {
		return unavailableDescriptor(file.kind, file.id, "malformed digest"),
			nil, false
	}
	descriptor := catalogv1alpha1.IndexCatalogDescriptor{
		ID:                file.id,
		Kind:              file.kind,
		Schema:            envelope.Schema,
		DerivationVersion: envelope.DerivationVersion,
		Digest:            envelope.Digest,
		InputDigests:      inputDigests(envelope.Inputs),
		Completeness:      string(envelope.Completeness),
		QueryRoutes: []string{
			"compile:" + file.kind,
			"check:" + file.kind,
		},
		Limitations: envelope.Limitations,
	}
	conflicts := make([]catalogv1alpha1.CatalogConflict, 0, len(envelope.Conflicts))
	for _, conflict := range envelope.Conflicts {
		conflicts = append(conflicts, catalogv1alpha1.CatalogConflict{
			ID:          file.id + "." + conflict.ID,
			Code:        conflict.Code,
			SourcePaths: conflict.SourcePaths,
		})
	}
	return descriptor, conflicts, true
}

// inputDigests projects a catalog's input digests as a sorted, deduplicated
// set. The index references inputs only by digest and path role.
func inputDigests(inputs []catalogv1alpha1.CatalogInput) []string {
	digests := make([]string, 0, len(inputs))
	for _, input := range inputs {
		digests = append(digests, input.Digest)
	}
	return catalogv1alpha1.SortedUnique(digests)
}

func validDigest(value string) bool {
	return len(value) == len("sha256:")+64 &&
		strings.HasPrefix(value, "sha256:")
}

func (c *compiler) catalog() (catalogv1alpha1.AgentSystemIndex, error) {
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
	descriptors := c.descriptors
	if descriptors == nil {
		descriptors = []catalogv1alpha1.IndexCatalogDescriptor{}
	}
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].ID < descriptors[j].ID
	})
	conflicts := c.conflicts
	if conflicts == nil {
		conflicts = []catalogv1alpha1.CatalogConflict{}
	}
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].ID < conflicts[j].ID
	})
	return catalogv1alpha1.AgentSystemIndex{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            catalogv1alpha1.APIVersion + "/" + catalogv1alpha1.KindAgentSystemIndex,
			Kind:              catalogv1alpha1.KindAgentSystemIndex,
			ID:                "agent-system.index",
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
			Completeness: catalogv1alpha1.CompletenessComplete,
			Limitations:  []string{},
			Conflicts:    conflicts,
		},
		Catalogs: descriptors,
	}, nil
}

func (c *compiler) buildOutputs(index catalogv1alpha1.AgentSystemIndex) ([]byte, string, error) {
	jsonContent, err := catalogv1alpha1.CanonicalJSONIndex(index)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	var digestValue catalogv1alpha1.AgentSystemIndex
	if err := catalogv1alpha1.DecodeStrict(jsonContent, &digestValue); err != nil {
		return nil, "", err
	}
	index.Digest = digestValue.Digest
	return jsonContent, catalogv1alpha1.RenderIndexMarkdown(index), nil
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
	index, err := compiler.catalog()
	if err != nil {
		return err
	}
	jsonContent, markdown, err := compiler.buildOutputs(index)
	if err != nil {
		return err
	}
	if opts.check {
		if len(compiler.limitations) > 0 || compiler.unavailable > 0 {
			return fmt.Errorf("agent system index is incomplete: %d problem(s)",
				len(compiler.limitations)+compiler.unavailable)
		}
		jsonPath := filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked index JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked index JSON is stale; rerun the generator")
		}
		markdownPath := filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked index Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked index Markdown is stale; rerun the generator")
		}
	}
	if err := os.MkdirAll(
		filepath.Dir(filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))),
		0o755,
	); err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath)),
		jsonContent,
		0o644,
	); err != nil {
		return err
	}
	if err := os.MkdirAll(
		filepath.Dir(filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))),
		0o755,
	); err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath)),
		[]byte(markdown),
		0o644,
	); err != nil {
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
		CatalogID:    index.ID,
		Output:       opts.outputPath,
		Markdown:     opts.markdownPath,
		Completeness: string(index.Completeness),
		Status:       "ok",
		Problems:     compiler.limitations,
	}
	if len(compiler.limitations) > 0 {
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
		fmt.Fprintln(os.Stderr, "index_check:", err)
		os.Exit(1)
	}
}
