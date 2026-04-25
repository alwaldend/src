package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	err := Execute(ctx, os.Args[1:], os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
