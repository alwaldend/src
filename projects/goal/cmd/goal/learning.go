package main

import (
	"fmt"
	"os"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	"github.com/spf13/cobra"
)

func newLearningCommand() *cobra.Command {
	stdout := os.Stdout
	input := ""
	output := ""
	command := &cobra.Command{
		Use:   "learning-proposal",
		Short: "Validate a bounded learning proposal from friction records",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if input == "" {
				return fmt.Errorf("--input is required")
			}
			content, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("read proposal: %w", err)
			}
			proposal, err := v1alpha1.DecodeLearningProposal(content)
			if err != nil {
				return err
			}
			if output == "" {
				return writeJSON(stdout, proposal)
			}
			encoded, err := v1alpha1.CanonicalLearningProposalJSON(proposal)
			if err != nil {
				return err
			}
			return os.WriteFile(output, encoded, 0o644)
		},
	}
	command.Flags().StringVar(&input, "input", "", "learning proposal JSON file")
	command.Flags().StringVar(&output, "output", "", "canonical output JSON file")
	return command
}
