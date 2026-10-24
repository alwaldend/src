package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// mcpServer supplies the protocol operations used by MCP commands.
type mcpServer interface {
	// Serve hosts protocol requests until disconnect or cancellation.
	Serve(ctx context.Context) error
}

// newMCPCommand delegates protocol hosting to the supplied server.
func newMCPCommand(protocol mcpServer) *cobra.Command {
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
