package main

import (
	"context"
	"fmt"
	"os"
)

// main runs the template command with process streams and reports failures.
func main() {
	ctx := context.Background()
	if err := execute(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
