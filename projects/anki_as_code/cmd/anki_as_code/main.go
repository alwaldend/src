package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/cli"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiserver"
	"github.com/spf13/cobra"
)

// run executes the command tree; each workflow owns its native lifetime.
func run(ctx context.Context, args []string) error {
	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve application executable: %w", err)
	}
	application := ankiascode.New(ankiascode.Options{BackendBinary: binary, BackendArguments: []string{"server", "serve"}})
	command := cli.New(application)
	serverCommand := &cobra.Command{
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
	runCommand := &cobra.Command{
		Use:                "run",
		Short:              "Run the linked native Anki CLI",
		DisableFlagParsing: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			// The blocking native CLI owns this process and must retain default interruption.
			signal.Reset(os.Interrupt)
			if err := ankiserver.Run(append([]string{command.CommandPath()}, arguments...)); err != nil {
				return fmt.Errorf("run linked Anki CLI: %w", err)
			}
			return nil
		},
	}
	serverCommand.AddCommand(runCommand, &cobra.Command{
		Use:                "serve",
		Short:              "Run the native Anki gRPC server",
		DisableFlagParsing: true,
		RunE: func(_ *cobra.Command, arguments []string) error {
			if err := runCommand.RunE(runCommand, append([]string{"serve"}, arguments...)); err != nil {
				return fmt.Errorf("serve Anki: %w", err)
			}
			return nil
		},
	})
	command.AddCommand(serverCommand)
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
