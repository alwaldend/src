package cli

import (
	"fmt"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"github.com/spf13/cobra"
)

// newServerCommand defines presentation and forwards opaque server arguments to App operations.
func newServerCommand(application ankiascodeapi.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "server",
		Short: "Manage the Anki server",
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
		Short:              "Run the linked Anki server CLI",
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
		Short:              "Run the Anki gRPC server",
		DisableFlagParsing: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			if err := application.ServerRun(command.Context(), run.CommandPath(), append([]string{"serve"}, arguments...)); err != nil {
				return fmt.Errorf("execute server serve: %w", err)
			}
			return nil
		},
	})
	return command
}
