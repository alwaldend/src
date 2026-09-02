package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if err := Execute(
		context.Background(),
		os.Args[1:],
		os.Getenv,
		os.Stdin,
		os.Stdout,
		os.Stderr,
		&execRunner{},
	); err != nil {
		fmt.Fprintf(os.Stderr, "repo_delivery: %v\n", err)
		os.Exit(1)
	}
}
