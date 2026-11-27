package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// run validates one skill, returning 2 for usage errors and 1 for failed validation or output.
func run(ctx context.Context, arguments []string, stdout, stderr io.Writer) int {
	var skillPath, openAIPath, expectedName, outputPath string
	var declaredFiles []string
	exitCode := 2
	command := &cobra.Command{
		Use:           "skill_validator",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			for _, name := range []string{"skill", "expected-name", "output"} {
				if cmd.Flags().Lookup(name).Value.String() == "" {
					return fmt.Errorf("--%s is required", name)
				}
			}
			if len(declaredFiles) == 0 {
				return fmt.Errorf("at least one --file is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			exitCode = 1
			issues := validateSkill(skillPath, openAIPath, expectedName, declaredFiles)
			if len(issues) != 0 {
				return fmt.Errorf("skill %s: %s", expectedName, strings.Join(issues, "\nERROR: skill "+expectedName+": "))
			}
			report := fmt.Sprintf("Skill is valid: %s\n", expectedName)
			if err := os.WriteFile(outputPath, []byte(report), 0o644); err != nil {
				return fmt.Errorf("write validation report %q: %w", outputPath, err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), report); err != nil {
				return fmt.Errorf("print validation report: %w", err)
			}
			return nil
		},
	}
	command.SetOut(stdout)
	command.SetErr(stderr)
	flags := command.Flags()
	flags.StringVar(&skillPath, "skill", "", "Path to SKILL.md")
	flags.StringVar(&openAIPath, "openai-yaml", "", "Optional path to agents/openai.yaml")
	flags.StringVar(&expectedName, "expected-name", "", "Skill name expected from the Bazel package")
	flags.StringVar(&outputPath, "output", "", "Validation report output")
	flags.StringArrayVar(&declaredFiles, "file", nil, "Declared file (repeatable)")
	for _, name := range []string{"skill", "expected-name", "output", "file"} {
		if err := command.MarkFlagRequired(name); err != nil {
			fmt.Fprintf(stderr, "ERROR: configure flag %q: %v\n", name, err)
			return 2
		}
	}
	// The original command also accepts a single dash before long option names.
	args := append([]string{}, arguments...)
	for i := 0; i < len(args); i++ {
		if args[i] == "--" {
			break
		}
		option, _, inline := strings.Cut(args[i], "=")
		if flags.Lookup(strings.TrimLeft(option, "-")) == nil || !strings.HasPrefix(option, "-") {
			continue
		}
		if !strings.HasPrefix(option, "--") {
			args[i] = "-" + args[i]
		}
		if !inline {
			i++
		}
	}
	command.SetArgs(args)
	if err := command.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return exitCode
	}
	return 0
}

// main validates the Bazel action inputs and exits with the corresponding status.
func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}
