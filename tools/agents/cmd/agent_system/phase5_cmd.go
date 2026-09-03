package agent_system

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

type normalizeOptions struct {
	input  string
	output string
}

func parseNormalizeFlags(args []string) (normalizeOptions, error) {
	var opts normalizeOptions
	flags := flag.NewFlagSet("normalize", flag.ContinueOnError)
	flags.StringVar(&opts.input, "input", "", "input skill-case JSON file")
	flags.StringVar(&opts.output, "output", "", "normalized output JSON file")
	if err := flags.Parse(args); err != nil {
		return normalizeOptions{}, err
	}
	if opts.input == "" || opts.output == "" {
		return normalizeOptions{}, fmt.Errorf("--input and --output are required")
	}
	return opts, nil
}

func runNormalize(args []string, stdout io.Writer) error {
	opts, err := parseNormalizeFlags(args)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(opts.input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	var inputs []v1alpha1.SkillCase
	if err := json.Unmarshal(content, &inputs); err != nil {
		return fmt.Errorf("decode skill cases: %w", err)
	}
	if len(inputs) == 0 {
		return fmt.Errorf("input must contain at least one skill case")
	}
	seen := make(map[string]bool, len(inputs))
	normalized := make([]v1alpha1.SkillCase, 0, len(inputs))
	for _, input := range inputs {
		if seen[input.ID] {
			return fmt.Errorf("duplicate normalized case %q", input.ID)
		}
		seen[input.ID] = true
		content, err := v1alpha1.CanonicalSkillCaseJSON(input)
		if err != nil {
			return err
		}
		var encoded v1alpha1.SkillCase
		if err := json.Unmarshal(content, &encoded); err != nil {
			return err
		}
		normalized = append(normalized, encoded)
	}
	output, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return err
	}
	output = append(output, '\n')
	if err := os.WriteFile(opts.output, output, 0o644); err != nil {
		return fmt.Errorf("write normalized cases: %w", err)
	}
	return writeJSONLine(stdout, map[string]any{
		"accepted": len(normalized),
		"output":   opts.output,
	})
}

func writeJSONLine(writer io.Writer, value any) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = writer.Write(append(content, '\n'))
	return err
}
