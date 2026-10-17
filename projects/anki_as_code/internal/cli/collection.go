package cli

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// newCollectionCommand groups native collection workflows and their shared namespace flag.
func newCollectionCommand(application application) *cobra.Command {
	command := &cobra.Command{
		Use:     "col",
		Aliases: []string{"collection"},
		Short:   "Inspect and manage Anki collections",
		Args:    cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return fmt.Errorf("col requires a subcommand")
		},
	}
	command.PersistentFlags().String(stateKeyFlag.name, state.DefaultKey, stateKeyFlag.usage)
	command.AddCommand(newExportCommand(application), newPlanCommand(application), newApplyCommand(application), newBuildCommand(application), newStateCommand(application))
	return command
}

// newExportCommand defines export flags and delegates execution to the application.
func newExportCommand(application application) *cobra.Command {
	var input, output string
	command := &cobra.Command{
		Use:   "export",
		Short: "Export a collection to editable TOML",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.ColExport(command.Context(), input, output, commandStateKey(command)); err != nil {
				return fmt.Errorf("export: %w", err)
			}
			if _, err := fmt.Fprintln(command.ErrOrStderr(), "Exported collection text to", output); err != nil {
				return fmt.Errorf("write export status: %w", err)
			}
			return nil
		},
	}
	collectionFlag.bindRequired(command, &input)
	outputFlag.bindRequired(command, &output)
	return command
}

// newPlanCommand defines plan flags and writes the resulting change plan as JSON.
func newPlanCommand(application application) *cobra.Command {
	var input, text string
	command := &cobra.Command{
		Use:   "plan",
		Short: "Inspect desired collection changes",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			plan, err := application.ColPlan(command.Context(), input, text, commandStateKey(command))
			if err != nil {
				return fmt.Errorf("plan: %w", err)
			}
			data, err := (protojson.MarshalOptions{Indent: "  "}).Marshal(plan)
			if err != nil {
				return fmt.Errorf("encode plan: %w", err)
			}
			if _, err := command.OutOrStdout().Write(append(data, '\n')); err != nil {
				return fmt.Errorf("write plan: %w", err)
			}
			return nil
		},
	}
	collectionFlag.bindRequired(command, &input)
	textFlag.bindRequired(command, &text)
	return command
}

// newApplyCommand defines apply flags and delegates in-place reconciliation.
func newApplyCommand(application application) *cobra.Command {
	var input, text string
	command := &cobra.Command{
		Use:   "apply",
		Short: "Reconcile an offline collection in place",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.ColApply(command.Context(), input, text, commandStateKey(command)); err != nil {
				return fmt.Errorf("apply: %w", err)
			}
			if _, err := fmt.Fprintln(command.ErrOrStderr(), "Reconciled collection", input); err != nil {
				return fmt.Errorf("write apply status: %w", err)
			}
			return nil
		},
	}
	collectionFlag.bindRequired(command, &input)
	textFlag.bindRequired(command, &text)
	return command
}

// newBuildCommand defines build flags and delegates collection archive creation.
func newBuildCommand(application application) *cobra.Command {
	var base, text, output string
	command := &cobra.Command{
		Use:   "build",
		Short: "Build a new archive from the original and desired text",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.ColBuild(command.Context(), base, text, output, commandStateKey(command)); err != nil {
				return fmt.Errorf("build: %w", err)
			}
			if _, err := fmt.Fprintln(command.ErrOrStderr(), "Built collection archive", output); err != nil {
				return fmt.Errorf("write build status: %w", err)
			}
			return nil
		},
	}
	baseFlag.bindRequired(command, &base)
	textFlag.bindRequired(command, &text)
	outputFlag.bindRequired(command, &output)
	return command
}
