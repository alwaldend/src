package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// main runs the hook command with process streams and interrupt cancellation.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := execute(ctx, os.Args[1:], os.Getenv("BUILD_WORKSPACE_DIRECTORY"), os.Stdout, os.Stderr)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
