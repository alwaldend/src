package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
)

// newRootCommand assembles the Cobra command tree around the supplied application.
func newRootCommand(application *app.App) *cobra.Command {
	command := &cobra.Command{
		Use:           "anki_as_code",
		Short:         "Manage Anki collection state and inspect configs",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	command.CompletionOptions.DisableDefaultCmd = true
	stateOptions := &app.StateOptions{}
	command.PersistentFlags().StringVar(&stateOptions.Key, "state-key", state.DefaultKey, "Collection state namespace key")
	command.AddCommand(newConfigCommand(application), newExportCommand(application), newPlanCommand(application), newApplyCommand(application),
		newBuildCommand(application), newStateCommand(application, stateOptions))
	return command
}

// newExportCommand defines export flags and delegates execution to the application.
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

// newPlanCommand defines plan flags and writes the resulting change plan as JSON.
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

// newApplyCommand defines apply flags and delegates in-place reconciliation.
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

// newBuildCommand defines build flags and delegates collection archive creation.
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

// run executes the command tree and always closes any initialized application.
func run(ctx context.Context, args []string) (err error) {
	application := &app.App{}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if stopErr := application.Stop(cleanup); stopErr != nil {
			err = errors.Join(err, fmt.Errorf("close application: %w", stopErr))
		}
	}()
	command := newRootCommand(application)
	command.SetArgs(args)
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

// main installs cancellation handling and reports command failures to stderr.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
