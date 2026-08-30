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
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintf(os.Stderr, "goal: %v\n", err)
		os.Exit(1)
	}
}
