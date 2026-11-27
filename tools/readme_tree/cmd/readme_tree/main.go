package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// main connects the command to process streams and interrupts.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := execute(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
