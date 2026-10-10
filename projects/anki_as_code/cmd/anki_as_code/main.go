package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/app"
)

func newRootCommand(application *app.App) *cobra.Command {
	command := &cobra.Command{
		Use:           "anki_as_code",
		Short:         "Anki as code CLI scaffold (operations not implemented)",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	command.CompletionOptions.DisableDefaultCmd = true
	command.AddCommand(newExportCommand(application), newPlanCommand(application), newApplyCommand(application),
		newBuildCommand(application), newGenerateIDCommand(application))
	return command
}

func newExportCommand(application *app.App) *cobra.Command {
	var input, output string
	command := &cobra.Command{
		Use:   "export",
		Short: "Export a collection to editable TOML",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.Export(command.Context(), input, output); err != nil {
				return fmt.Errorf("export: %w", err)
			}
			fmt.Fprintln(command.ErrOrStderr(), "Exported collection text to", output)
			return nil
		},
	}
	inputFlag.bindRequired(command, &input)
	outputFlag.bindRequired(command, &output)
	return command
}

func newPlanCommand(application *app.App) *cobra.Command {
	var input, text string
	command := &cobra.Command{
		Use:   "plan",
		Short: "Inspect desired collection changes",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			plan, err := application.Plan(command.Context(), input, text)
			if err != nil {
				return fmt.Errorf("plan: %w", err)
			}
			encoder := json.NewEncoder(command.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(plan); err != nil {
				return fmt.Errorf("write plan: %w", err)
			}
			return nil
		},
	}
	inputFlag.bindRequired(command, &input)
	textFlag.bindRequired(command, &text)
	return command
}

func newApplyCommand(application *app.App) *cobra.Command {
	var input, text string
	command := &cobra.Command{
		Use:   "apply",
		Short: "Reconcile an offline collection in place",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.Apply(command.Context(), input, text); err != nil {
				return fmt.Errorf("apply: %w", err)
			}
			fmt.Fprintln(command.ErrOrStderr(), "Reconciled collection", input)
			return nil
		},
	}
	inputFlag.bindRequired(command, &input)
	textFlag.bindRequired(command, &text)
	return command
}

func newBuildCommand(application *app.App) *cobra.Command {
	var base, text, output string
	command := &cobra.Command{
		Use:   "build",
		Short: "Build a new archive from the original and desired text",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := application.Build(command.Context(), base, text, output); err != nil {
				return fmt.Errorf("build: %w", err)
			}
			fmt.Fprintln(command.ErrOrStderr(), "Built collection archive", output)
			return nil
		},
	}
	baseFlag.bindRequired(command, &base)
	textFlag.bindRequired(command, &text)
	outputFlag.bindRequired(command, &output)
	return command
}

func newGenerateIDCommand(application *app.App) *cobra.Command {
	var config string
	command := &cobra.Command{
		Use:   "generate-id <note-file-or-directory>...",
		Short: "Give note files new note, card, and synchronization identities",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, paths []string) error {
			count, err := application.GenerateIDsWithConfig(command.Context(), paths, config)
			if err != nil {
				return fmt.Errorf("generate note identities: %w", err)
			}
			fmt.Fprintf(command.ErrOrStderr(), "Generated identities for %d notes\n", count)
			return nil
		},
	}
	collectionFlag.bind(command, &config)
	return command
}

func run(ctx context.Context, args []string) error {
	application := &app.App{}
	command := newRootCommand(application)
	command.SetArgs(args)
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
