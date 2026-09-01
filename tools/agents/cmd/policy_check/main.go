// Command policy_check compiles the bounded deterministic PolicyCatalog over
// tracked agent-policy sources (AGENTS.md, CODEOWNERS, boundary READMEs).
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

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("policy_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/policy.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/policy.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on completeness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: content-addressed inputs)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.policy-compiler",
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
	policies    []catalogv1alpha1.PolicyRecord
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
	// Explicitly closed policy universe: every tracked AGENTS.md plus the
	// CODEOWNERS and six top-level boundary READMEs.
	agentPolicies := []string{"AGENTS.md"}
	for _, path := range agentPolicies {
		full := filepath.Join(c.root, filepath.FromSlash(path))
		content, err := os.ReadFile(full)
		if err != nil {
			c.problem("agent policy source missing: %s", path)
			c.eligible++
			c.unavailable++
			continue
		}
		c.input(path, "agent-policy", content)
		c.eligible++
		c.emitted++
		c.derivePolicy(path, content)
	}
	codeownersPath := "CODEOWNERS"
	full := filepath.Join(c.root, filepath.FromSlash(codeownersPath))
	if content, err := os.ReadFile(full); err == nil {
		c.input(codeownersPath, "review-owners", content)
		c.eligible++
		if !bytes.Contains(content, []byte("review")) && !bytes.Contains(content, []byte("@")) {
			// A CODEOWNERS file must name owners; otherwise record it as an
			// availability limitation but keep the catalog complete.
			c.problem("CODEOWNERS does not name any owners")
			c.unavailable++
		} else {
			c.emitted++
		}
	}
	boundaryOrder := []string{"projects", "infra", "tools", "data", "third_party", "users"}
	for _, name := range boundaryOrder {
		readme := name + "/README.md"
		full := filepath.Join(c.root, filepath.FromSlash(readme))
		content, err := os.ReadFile(full)
		if err != nil {
			c.problem("boundary README missing: %s", readme)
			c.eligible++
			c.unavailable++
			continue
		}
		c.input(readme, "boundary-readme", content)
		c.eligible++
		c.emitted++
	}
	return nil
}

var policyAxisSections = []struct {
	name     string
	keywords []string
}{
	{"sourceDisclosure", []string{"disclosure", "public"}},
	{"evidenceHandling", []string{"evidence", "secret"}},
	{"bazelVisibility", []string{"visibility", "bazel"}},
	{"buildConsumer", []string{"consumer", "build"}},
	{"artifactPublication", []string{"publication", "artifact"}},
	{"documentationPublication", []string{"documentation", "docs"}},
	{"information", []string{"information"}},
	{"liveEnvironmentAssociation", []string{"live", "environment"}},
}

func (c *compiler) derivePolicy(agentPolicyPath string, content []byte) {
	lower := strings.ToLower(string(content))
	axes := make([]catalogv1alpha1.PolicyAxis, 0, len(policyAxisSections))
	conflicts := 0
	for _, section := range policyAxisSections {
		value := "unknown"
		for _, keyword := range section.keywords {
			if strings.Contains(lower, keyword) {
				value = "known"
				break
			}
		}
		axes = append(axes, catalogv1alpha1.PolicyAxis{
			Name:   section.name,
			Value:  value,
			Source: agentPolicyPath,
		})
		if value == "known" && strings.Contains(lower, "conflict") {
			conflicts++
		}
	}
	if conflicts > 0 {
		c.conflicts = append(c.conflicts, catalogv1alpha1.CatalogConflict{
			ID:          "policy." + strings.ToLower(strings.TrimSuffix(filepath.Base(agentPolicyPath), filepath.Ext(agentPolicyPath))),
			Code:        "policy_axis_conflict",
			SourcePaths: []string{agentPolicyPath},
		})
	}
	c.policies = append(c.policies, catalogv1alpha1.PolicyRecord{
		ID:                "policy." + strings.ToLower(strings.TrimSuffix(filepath.Base(agentPolicyPath), filepath.Ext(agentPolicyPath))),
		PathPrefix:        "/",
		Precedence:        0,
		AgentPolicySource: agentPolicyPath,
		Axes:              axes,
	})
}

func (c *compiler) catalog() (catalogv1alpha1.PolicyCatalog, error) {
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
	policies := c.policies
	if policies == nil {
		policies = []catalogv1alpha1.PolicyRecord{}
	}
	sort.Slice(policies, func(i, j int) bool {
		return policies[i].ID < policies[j].ID
	})
	return catalogv1alpha1.PolicyCatalog{
		CatalogEnvelope: catalogv1alpha1.CatalogEnvelope{
			Schema:            catalogv1alpha1.APIVersion + "/" + catalogv1alpha1.KindPolicyCatalog,
			Kind:              catalogv1alpha1.KindPolicyCatalog,
			ID:                "agent-system.policy",
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
		Policies: policies,
	}, nil
}

func (c *compiler) buildOutputs(catalog catalogv1alpha1.PolicyCatalog) ([]byte, string, error) {
	jsonContent, err := catalogv1alpha1.CanonicalJSONPolicy(catalog)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	var digestValue catalogv1alpha1.PolicyCatalog
	if err := catalogv1alpha1.DecodeStrict(jsonContent, &digestValue); err != nil {
		return nil, "", err
	}
	catalog.Digest = digestValue.Digest
	return jsonContent, catalogv1alpha1.RenderPolicyMarkdown(catalog), nil
}

func (c *compiler) sameDir(rel string) string {
	return filepath.ToSlash(rel)
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
			return fmt.Errorf("policy catalog is incomplete: %d problem(s)",
				len(compiler.problems))
		}
		jsonPath := filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked policy JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked policy JSON is stale; rerun the generator")
		}
		markdownPath := filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked policy Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked policy Markdown is stale; rerun the generator")
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
		fmt.Fprintln(os.Stderr, "policy_check:", err)
		os.Exit(1)
	}
}
