package agent_system

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

type aggregateOptions struct {
	catalog  string
	input    string
	output   string
	markdown string
}

func parseAggregateFlags(args []string) (aggregateOptions, error) {
	var opts aggregateOptions
	flags := flag.NewFlagSet("aggregate", flag.ContinueOnError)
	flags.StringVar(&opts.catalog, "catalog", "",
		"capability catalog JSON file")
	flags.StringVar(&opts.input, "input", "skill-cases.json",
		"normalized skill-case JSON file")
	flags.StringVar(&opts.output, "output", "coverage-matrix.json",
		"coverage matrix JSON output")
	flags.StringVar(&opts.markdown, "markdown", "",
		"coverage matrix Markdown output")
	if err := flags.Parse(args); err != nil {
		return aggregateOptions{}, err
	}
	if opts.catalog == "" || opts.output == "" {
		return aggregateOptions{}, fmt.Errorf(
			"--catalog and --output are required")
	}
	if opts.markdown != "" && opts.input == "" {
		return aggregateOptions{}, fmt.Errorf(
			"--markdown requires --input")
	}
	return opts, nil
}

func runAggregate(args []string, stdout io.Writer) error {
	opts, err := parseAggregateFlags(args)
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
	catalogSkillIDs := make(map[string]bool, len(catalog.Skills))
	for _, skill := range catalog.Skills {
		catalogSkillIDs[skill.ID] = true
	}
	totalSkills := len(catalog.Skills)
	if totalSkills == 0 {
		return fmt.Errorf("capability catalog contains no skills")
	}
	caseContent, err := os.ReadFile(opts.input)
	if err != nil {
		return fmt.Errorf("read skill cases: %w", err)
	}
	var cases []v1alpha1.SkillCase
	if err := json.Unmarshal(caseContent, &cases); err != nil {
		return fmt.Errorf("decode skill cases: %w", err)
	}
	if len(cases) == 0 {
		return fmt.Errorf("skill-case input must contain at least one case")
	}
	seen := make(map[string]bool, len(cases))
	for _, value := range cases {
		if err := value.Validate(); err != nil {
			return fmt.Errorf("skill case %q: %w", value.ID, err)
		}
		if !catalogSkillIDs[value.SkillID] {
			return fmt.Errorf(
				"skill case %q references skill %q absent from the bound catalog",
				value.ID, value.SkillID,
			)
		}
		identity := value.SkillID + "\x00" + value.ID
		if seen[identity] {
			return fmt.Errorf("duplicate skill case %q for skill %q",
				value.ID, value.SkillID)
		}
		seen[identity] = true
	}
	sort.Slice(cases, func(left int, right int) bool {
		if cases[left].SkillID != cases[right].SkillID {
			return cases[left].SkillID < cases[right].SkillID
		}
		return cases[left].ID < cases[right].ID
	})
	entries := make([]v1alpha1.CoverageEntry, 0, len(cases))
	coveredSkills := make(map[string]bool)
	for _, value := range cases {
		coveredSkills[value.SkillID] = true
		entries = append(entries, v1alpha1.CoverageEntry{
			SkillID:      value.SkillID,
			CaseID:       value.ID,
			State:        v1alpha1.CoverageState(value.EvidenceTier),
			Metric:       value.Metric,
			EvidenceTier: value.EvidenceTier,
			EvidenceRef:  value.SourceRef,
		})
	}
	total := totalSkills
	truncated := len(coveredSkills) < totalSkills
	if len(entries) > totalSkills {
		total = len(entries)
	}
	matrix := v1alpha1.CoverageMatrix{
		APIVersion: v1alpha1.APIVersion,
		Kind:       "CoverageMatrix",
		ID:         "coverage/repository-skills",
		CatalogRef: v1alpha1.Reference{
			Kind:   v1alpha1.ReferenceArtifact,
			ID:     "artifact/agent-system-capability",
			Digest: catalog.Digest,
		},
		Entries:   entries,
		Total:     total,
		Truncated: truncated,
	}
	encoded, err := v1alpha1.CanonicalCoverageMatrixJSON(matrix)
	if err != nil {
		return err
	}
	var finalMatrix v1alpha1.CoverageMatrix
	if err := json.Unmarshal(encoded, &finalMatrix); err != nil {
		return fmt.Errorf("decode coverage matrix: %w", err)
	}
	if err := os.WriteFile(opts.output, encoded, 0o644); err != nil {
		return fmt.Errorf("write coverage matrix: %w", err)
	}
	if opts.markdown != "" {
		if err := os.WriteFile(
			opts.markdown,
			[]byte(renderCoverageMarkdown(finalMatrix, totalSkills)),
			0o644,
		); err != nil {
			return fmt.Errorf("write coverage Markdown: %w", err)
		}
	}
	return writeJSONLine(stdout, map[string]any{
		"cases":     len(cases),
		"total":     totalSkills,
		"truncated": finalMatrix.Truncated,
		"output":    opts.output,
		"markdown":  opts.markdown,
		"digest":    finalMatrix.Digest,
	})
}

func renderCoverageMarkdown(
	matrix v1alpha1.CoverageMatrix,
	totalSkills int,
) string {
	content := fmt.Sprintf(
		"# Skill coverage matrix\n\n"+
			"- Case entries: %d\n"+
			"- Capability skills: %d\n"+
			"- Truncated: %t\n"+
			"- Catalog digest: `%s`\n"+
			"- Matrix digest: `%s`\n\n",
		len(matrix.Entries),
		totalSkills,
		matrix.Truncated,
		matrix.CatalogRef.Digest,
		matrix.Digest,
	)
	content += "| Skill | Case | Metric | Evidence |\n"
	content += "| --- | --- | --- | --- |\n"
	for _, entry := range matrix.Entries {
		evidence := "none"
		if entry.EvidenceRef.ID != "" {
			evidence = entry.EvidenceRef.ID
		}
		content += fmt.Sprintf(
			"| `%s` | `%s` | `%s` | `%s` |\n",
			entry.SkillID,
			entry.CaseID,
			entry.Metric,
			evidence,
		)
	}
	content += "\nThis projection reports only normalized fixture-tested cases.\n" +
		"It does not claim live, scheduled, or complete behavioral coverage.\n"
	return content
}
