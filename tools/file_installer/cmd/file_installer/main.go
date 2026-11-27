package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	cmd "git.alwaldend.com/alwaldend/src/tools/file_installer/internal/command"
)

// main runs the installer and reports failures to the caller.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := cmd.Execute(ctx)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
