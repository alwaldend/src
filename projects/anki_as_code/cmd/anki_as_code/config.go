package main

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/app"
	"github.com/spf13/cobra"
)

// newConfigCommand groups operations on text configuration independently of Anki collections.
func newConfigCommand(application *app.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Inspect and manipulate collection configuration files",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return fmt.Errorf("config requires a subcommand")
		},
	}
	command.AddCommand(newConfigDumpCommand(application))
	return command
}

// newConfigDumpCommand defines a read-only configuration dump with selectable output format.
func newConfigDumpCommand(application *app.App) *cobra.Command {
	var text, format string
	command := &cobra.Command{
		Use:   "dump",
		Short: "Print parsed collection configurations",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			data, err := application.Dump(command.Context(), text, format)
			if err != nil {
				return fmt.Errorf("dump configuration: %w", err)
			}
			if _, err := command.OutOrStdout().Write(data); err != nil {
				return fmt.Errorf("write configuration dump: %w", err)
			}
			return nil
		},
	}
	textFlag.bindRequired(command, &text)
	command.Flags().StringVar(&format, "format", "json", "Output format: json or toml")
	return command
}
