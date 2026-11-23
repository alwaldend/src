package main

import (
	"context"
	"fmt"
	"io"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/app"
	"github.com/spf13/cobra"
)

// executeCommand registers the Cobra routes and invokes the selected handler.
func executeCommand(
	ctx context.Context,
	arguments []string,
	runService func(context.Context, app.Command) error,
	exportGraph func(context.Context, io.Writer) error,
) error {
	root := &cobra.Command{
		Use:           "mlflow-vastai",
		Short:         "Run MLflow Vast.ai services",
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	serve := func(kind app.CommandKind) func(*cobra.Command, []string) error {
		return func(command *cobra.Command, _ []string) error {
			if err := runService(command.Context(), app.Command{Kind: kind}); err != nil {
				return fmt.Errorf("dispatch %s: %w", command.CommandPath(), err)
			}
			return nil
		}
	}

	server := &cobra.Command{Use: "server", Short: "Manage the job server"}
	server.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Run the job server",
		Args:  cobra.NoArgs,
		RunE:  serve(app.ServeJobs),
	})

	workerCommand := &cobra.Command{Use: "worker", Short: "Manage the rental-local worker"}
	workerCommand.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Run the rental-local worker",
		Args:  cobra.NoArgs,
		RunE:  serve(app.RunWorker),
	})

	graph := &cobra.Command{Use: "state-machine", Short: "Inspect the application state machines"}
	graph.AddCommand(&cobra.Command{
		Use:   "graph",
		Short: "Print the composed application graph as DOT",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := exportGraph(command.Context(), command.OutOrStdout()); err != nil {
				return fmt.Errorf("export state machine graph: %w", err)
			}
			return nil
		},
	})

	root.AddCommand(server, workerCommand, graph)
	root.SetArgs(arguments)
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("process command: %w", err)
	}
	return nil
}
