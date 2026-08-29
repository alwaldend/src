package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type lookPathFunc func(string) (string, error)
type replaceProcessFunc func(string, []string, []string) error

func bazelArguments(args []string) []string {
	result := make([]string, 0, len(args)+2)
	result = append(result, "--batch")
	if len(args) == 0 {
		return result
	}
	result = append(result, args[0], "--config=agent")
	return append(result, args[1:]...)
}

func run(
	args []string,
	environment []string,
	lookPath lookPathFunc,
	replaceProcess replaceProcessFunc,
) error {
	bazelPath, err := lookPath("bazel")
	if err != nil {
		return fmt.Errorf("find bazel in PATH: %w", err)
	}
	processArgs := append([]string{bazelPath}, bazelArguments(args)...)
	if err := replaceProcess(bazelPath, processArgs, environment); err != nil {
		return fmt.Errorf("execute bazel: %w", err)
	}
	return nil
}

func main() {
	err := run(os.Args[1:], os.Environ(), exec.LookPath, syscall.Exec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bazel_agent: %v\n", err)
		os.Exit(1)
	}
}
