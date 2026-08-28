package main

import (
	"fmt"
	"os"
)

func main() {
	if err := execute(os.Args[1:], os.Getenv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
