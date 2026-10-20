package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/cli"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
)

// run composes the CLI with one server adapter and a same-executable worker.
func run(ctx context.Context, arguments []string) error {
	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve application executable: %w", err)
	}
	hostedRunner := func(arguments []string) error {
		// The server loop does not observe main's Go cancellation context, so
		// restore SIGINT termination instead of only cancelling that context.
		signal.Reset(os.Interrupt)
		// Native writes to a disconnected client can raise SIGPIPE; let the
		// server handle the write error instead of terminating the worker.
		signal.Ignore(syscall.SIGPIPE)
		if err := runServer(arguments); err != nil {
			return fmt.Errorf("host configured Anki server: %w", err)
		}
		return nil
	}
	workerFactory := ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](ankigrpc.NewSubprocessClient)
	factory := ankigrpc.NewFactory(ankigrpc.SubprocessAddress(binary, "server", "serve"))
	for _, scheme := range subprocessSchemes {
		factory.RegisterClientFactory(scheme, workerFactory)
	}
	application := ankiascode.New(ankiascode.Options{
		NewClient:    factory,
		ServerRunner: hostedRunner,
	})
	command := cli.New(application)
	command.SetArgs(arguments)
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

// main owns executable cancellation, command invocation, and failure reporting.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
