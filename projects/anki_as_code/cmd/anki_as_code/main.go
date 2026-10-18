package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/cli"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiserver"
)

// run wires App dependencies and executes the complete CLI.
func run(ctx context.Context, args []string) error {
	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve application executable: %w", err)
	}
	var application ankiascodeapi.App = ankiascode.New(ankiascode.Options{
		BackendBinary: binary, BackendArguments: []string{"server", "serve"}, ServerRunner: runNativeServer,
	})
	command := cli.New(application)
	command.SetArgs(args)
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

// runNativeServer adapts interrupt and broken-pipe signals for the linked native CLI.
func runNativeServer(arguments []string) error {
	signal.Reset(os.Interrupt)
	signal.Ignore(syscall.SIGPIPE)
	if err := ankiserver.Run(arguments); err != nil {
		return fmt.Errorf("host linked native CLI: %w", err)
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
