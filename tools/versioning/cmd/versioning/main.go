package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := execute(os.Args[1:], time.Now, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "versioning:", err)
		os.Exit(processExitCode(err))
	}
}

func processExitCode(err error) int {
	var bazelError *bazelExitError
	if errors.As(err, &bazelError) && bazelError.ExitCode() > 0 {
		return bazelError.ExitCode()
	}
	return 1
}
