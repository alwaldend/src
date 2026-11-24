package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	jobserver "git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/server"
	"github.com/spf13/cobra"
)

// executeCommand registers the Cobra routes and invokes the selected handler.
func executeCommand(ctx context.Context, arguments []string, cli *cli) error {
	root := &cobra.Command{
		Use:           "mlflow-vastai",
		Short:         "Run MLflow Vast.ai services",
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	route := func(command *cobra.Command, event machine.Event, operation func(context.Context) error) error {
		if err := cli.selectRoute(command.Context(), event, operation); err != nil {
			return fmt.Errorf("%s: %w", command.CommandPath(), err)
		}
		return nil
	}

	server := &cobra.Command{Use: "server", Short: "Manage the job server"}
	var config jobserver.Config
	serve := &cobra.Command{
		Use:   "serve",
		Short: "Run the job server",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return route(command, dispatch, func(ctx context.Context) error {
				if err := cli.application.ServeJobs(ctx, config); err != nil {
					return fmt.Errorf("serve coordinator: %w", err)
				}
				return nil
			})
		},
	}
	serve.Flags().StringVar(&config.Listen, "listen", "127.0.0.1:50051", "Local gRPC address")
	serve.Flags().StringVar(&config.Directory, "state-dir", "", "Private directory for coordinator state")
	serve.Flags().StringSliceVar(&config.Workers, "workers", nil, "Allowed Vast.ai rental IDs or exact labels")
	serve.Flags().StringVar(&config.VastEndpoint, "vast-endpoint", "https://console.vast.ai/api/v0", "Vast.ai API endpoint")
	serve.Flags().StringVar(&config.TrackingURI, "tracking-uri", os.Getenv("MLFLOW_TRACKING_URI"), "MLflow tracking endpoint")
	serve.Flags().StringVar(&config.SSHUser, "ssh-user", "root", "SSH account on the rentals")
	server.AddCommand(serve)

	workerCommand := &cobra.Command{Use: "worker", Short: "Run or contact the rental-local worker"}
	var directory string
	workerCommand.PersistentFlags().StringVar(&directory, "state-dir", "", "Directory containing the worker socket and state")
	workerRoute := func(name, description string, operation func(context.Context, string, io.Reader, io.Writer) error) *cobra.Command {
		return &cobra.Command{
			Use:   name,
			Short: description,
			Args:  cobra.NoArgs,
			RunE: func(command *cobra.Command, _ []string) error {
				return route(command, dispatch, func(ctx context.Context) error {
					if err := operation(ctx, directory, command.InOrStdin(), command.OutOrStdout()); err != nil {
						return fmt.Errorf("run worker %s: %w", name, err)
					}
					return nil
				})
			},
		}
	}
	workerCommand.AddCommand(
		&cobra.Command{
			Use:   "serve",
			Short: "Run the rental-local worker",
			Args:  cobra.NoArgs,
			RunE: func(command *cobra.Command, _ []string) error {
				return route(command, dispatch, func(ctx context.Context) error {
					if err := cli.application.ServeWorker(ctx, directory); err != nil {
						return fmt.Errorf("serve worker: %w", err)
					}
					return nil
				})
			},
		},
		workerRoute("submit", "Submit a protobuf job request from stdin", cli.application.Submit),
		workerRoute("status", "Request the durable status of a job", cli.application.Status),
		workerRoute("cancel", "Request cancellation of a job", cli.application.Cancel),
	)

	graph := &cobra.Command{Use: "state-machine", Short: "Inspect the application state machines"}
	graph.AddCommand(&cobra.Command{
		Use:   "graph",
		Short: "Print the composed application graph as DOT",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return route(command, showGraph, func(_ context.Context) error {
				if err := cli.writeGraph(command.OutOrStdout()); err != nil {
					return fmt.Errorf("export state machine graph: %w", err)
				}
				return nil
			})
		},
	})

	root.AddCommand(server, workerCommand, graph)
	root.SetArgs(arguments)
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("process command: %w", err)
	}
	return nil
}
