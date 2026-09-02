package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

type coverageOptions struct {
	catalog string
	input   string
	output  string
}

func parseCoverageFlags(args []string) (coverageOptions, error) {
	var opts coverageOptions
	flags := flag.NewFlagSet("coverage", flag.ContinueOnError)
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
	return opts, nil
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
	catalogEnvelope := struct {
		Digest string     `json:"digest"`
		Skills []struct{} `json:"skills"`
	}{}
	if err := json.Unmarshal(catalogContent, &catalogEnvelope); err != nil {
		return fmt.Errorf("decode capability catalog: %w", err)
	}
	totalSkills := len(catalogEnvelope.Skills)
	content, err := os.ReadFile(opts.input)
	if err != nil {
		return fmt.Errorf("read coverage input: %w", err)
	}
	var entries []v1alpha1.CoverageEntry
	if err := json.Unmarshal(content, &entries); err != nil {
		return fmt.Errorf("decode coverage entries: %w", err)
	}
	matrix := v1alpha1.CoverageMatrix{
		APIVersion: v1alpha1.APIVersion,
		Kind:       "CoverageMatrix",
		ID:         "coverage/agent-system",
		CatalogRef: v1alpha1.Reference{
			Kind:   v1alpha1.ReferenceArtifact,
			ID:     "artifact/agent-system-capability",
			Digest: catalogEnvelope.Digest,
		},
		Entries:   entries,
		Total:     totalSkills,
		Truncated: len(entries) < totalSkills,
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
