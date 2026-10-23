package cli

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"github.com/spf13/cobra"
)

// newMCPCommand delegates protocol hosting to the injected MCP interface.
func newMCPCommand(protocol ankiascodeapi.MCP) *cobra.Command {
	command := &cobra.Command{
		Use: "mcp", Short: "Expose workflows and native Anki API tools over MCP", Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return fmt.Errorf("mcp requires the serve subcommand") },
	}
	command.AddCommand(&cobra.Command{
		Use: "serve", Short: "Serve MCP over stdin/stdout", Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if protocol == nil {
				return fmt.Errorf("MCP server is not configured")
			}
			if err := protocol.Serve(command.Context()); err != nil {
				return fmt.Errorf("serve MCP: %w", err)
			}
			return nil
		},
	})
	return command
}
