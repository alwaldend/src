package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newServerCommand defines presentation and forwards opaque native arguments to App operations.
func newServerCommand(application application) *cobra.Command {
	command := &cobra.Command{
		Use:   "server",
		Short: "Manage the native Anki server",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := command.Help(); err != nil {
				return fmt.Errorf("show server help: %w", err)
			}
			return nil
		},
	}
	run := &cobra.Command{
		Use:                "run",
		Short:              "Run the linked native Anki CLI",
		DisableFlagParsing: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			if err := application.ServerRun(command.Context(), command.CommandPath(), arguments); err != nil {
				return fmt.Errorf("run server: %w", err)
			}
			return nil
		},
	}
	command.AddCommand(run, &cobra.Command{
		Use:                "serve",
		Short:              "Run the native Anki gRPC server",
		DisableFlagParsing: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			if err := application.ServerServe(command.Context(), run.CommandPath(), arguments); err != nil {
				return fmt.Errorf("execute server serve: %w", err)
			}
			return nil
		},
	})
	return command
}
