package agent_system

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

type coverageOptions struct {
	workspaceRoot string
	catalog       string
	input         string
	output        string
}

func parseCoverageFlags(args []string) (coverageOptions, error) {
	var opts coverageOptions
	flags := flag.NewFlagSet("coverage", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"base for relative paths (default: BUILD_WORKSPACE_DIRECTORY or working directory)")
	flags.StringVar(&opts.catalog, "catalog", "",
		"capability catalog JSON file")
	flags.StringVar(&opts.input, "input", "", "coverage entries JSON file")
	flags.StringVar(&opts.output, "output", "", "coverage matrix JSON output")
	if err := flags.Parse(args); err != nil {
		return coverageOptions{}, err
	}
	if opts.catalog == "" || opts.input == "" || opts.output == "" {
		return coverageOptions{}, fmt.Errorf(
			"--catalog, --input, and --output are required",
		)
	}
	if flags.NArg() != 0 {
		return coverageOptions{}, fmt.Errorf("unexpected positional arguments")
	}
	err := workspaceFilePaths(opts.workspaceRoot, &opts.catalog, &opts.input, &opts.output)
	return opts, err
}

func runCoverage(args []string, stdout io.Writer) error {
	opts, err := parseCoverageFlags(args)
	if err != nil {
		return err
	}
	catalogContent, err := os.ReadFile(opts.catalog)
	if err != nil {
		return fmt.Errorf("read capability catalog: %w", err)
	}
	catalog, err := catalogv1alpha1.DecodeCapabilityStrict(catalogContent)
	if err != nil {
		return fmt.Errorf("decode capability catalog: %w", err)
	}
	content, err := os.ReadFile(opts.input)
	if err != nil {
		return fmt.Errorf("read coverage input: %w", err)
	}
	var entries []v1alpha1.CoverageEntry
	if err := json.Unmarshal(content, &entries); err != nil {
		return fmt.Errorf("decode coverage entries: %w", err)
	}
	matrix, err := configuredCoverageMatrix(catalog, entries,
		"coverage/agent-system")
	if err != nil {
		return err
	}
	encoded, err := v1alpha1.CanonicalCoverageMatrixJSON(matrix)
	if err != nil {
		return err
	}
	if err := os.WriteFile(opts.output, encoded, 0o644); err != nil {
		return fmt.Errorf("write coverage matrix: %w", err)
	}
	return writeJSONLine(stdout, map[string]any{
		"entries": len(entries),
		"output":  opts.output,
	})
}

// Neither inventory command loads or verifies a model execution result.
// Reject observed-evidence claims until a result-verifying importer exists.
func configuredCoverageMatrix(
	catalog catalogv1alpha1.CapabilityCatalog,
	entries []v1alpha1.CoverageEntry,
	id string,
) (v1alpha1.CoverageMatrix, error) {
	skillIDs := make(map[string]bool, len(catalog.Skills))
	for _, skill := range catalog.Skills {
		skillIDs[skill.ID] = true
	}
	coveredSkills := make(map[string]bool)
	for _, entry := range entries {
		if entry.EvidenceTier != v1alpha1.TierConfigured ||
			entry.State != v1alpha1.CoverageConfigured {
			return v1alpha1.CoverageMatrix{}, fmt.Errorf(
				"case %q: coverage inventories configured cases only; execution evidence requires a result-verifying importer",
				entry.CaseID,
			)
		}
		if !skillIDs[entry.SkillID] {
			return v1alpha1.CoverageMatrix{}, fmt.Errorf(
				"case %q references skill %q absent from the bound catalog",
				entry.CaseID, entry.SkillID,
			)
		}
		coveredSkills[entry.SkillID] = true
	}
	return v1alpha1.CoverageMatrix{
		APIVersion: v1alpha1.APIVersion,
		Kind:       "CoverageMatrix",
		ID:         id,
		CatalogRef: v1alpha1.Reference{
			Kind:   v1alpha1.ReferenceArtifact,
			ID:     "artifact/agent-system-capability",
			Digest: catalog.Digest,
		},
		Entries:       entries,
		Total:         len(entries),
		Truncated:     false,
		TotalSkills:   len(catalog.Skills),
		CoveredSkills: len(coveredSkills),
	}, nil
}
