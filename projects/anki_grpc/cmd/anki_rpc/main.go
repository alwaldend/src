package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/cli"
)

// run assembles the generated commands and executes them with caller cancellation.
func run(ctx context.Context) error {
	command, err := cli.NewCommand()
	if err != nil {
		return fmt.Errorf("construct commands: %w", err)
	}
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

// main translates interrupts into cancellation and reports failures on stderr.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
