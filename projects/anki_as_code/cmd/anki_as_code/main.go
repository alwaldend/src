package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/cli"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
)

// run executes the command tree; each workflow owns its native lifetime.
func run(ctx context.Context, args []string) error {
	application := &ankiascode.App{}
	command := cli.New(application)
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
