package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
)

// configuration describes the rule-generated executable and argument prefix.
type configuration struct {
	// Schema identifies the supported generated configuration format.
	Schema string `json:"schema"`
	// Binary is the canonical runfiles key of the wrapped executable.
	Binary string `json:"binary"`
	// Arguments precede caller-supplied application arguments.
	Arguments []string `json:"arguments"`
	// Properties are the consumer's systemd NAME=VALUE assignments.
	Properties []string `json:"properties"`
}

// main parses wrapper options and reports failures before process replacement.
func main() {
	configured, err := launchConfigured()
	if err != nil {
		fmt.Fprintf(os.Stderr, "systemd_run: %s\n", err)
		os.Exit(1)
	}
	if configured {
		return
	}
	var properties []string
	command := &cobra.Command{
		Use:           "systemd_run --property NAME=VALUE -- executable [arguments...]",
		Short:         "Run an executable in a resource-limited user systemd scope",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			if err := launch(properties, arguments); err != nil {
				return fmt.Errorf("launch scoped command: %w", err)
			}
			return nil
		},
	}
	command.Flags().SetInterspersed(false)
	command.Flags().StringArrayVar(&properties, "property", nil, "Systemd scope property as NAME=VALUE (repeatable)")
	if err := command.MarkFlagRequired("property"); err != nil {
		fmt.Fprintf(os.Stderr, "systemd_run: configure required property: %s\n", err)
		os.Exit(1)
	}
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "systemd_run: %s\n", err)
		os.Exit(1)
	}
}

// launchConfigured loads the adjacent rule configuration or selects CLI mode.
func launchConfigured() (bool, error) {
	executable, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("resolve launcher executable: %w", err)
	}
	content, err := os.ReadFile(executable + ".json")
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return true, fmt.Errorf("read launcher configuration: %w", err)
	}
	var configured configuration
	if err := json.Unmarshal(content, &configured); err != nil {
		return true, fmt.Errorf("decode launcher configuration: %w", err)
	}
	if configured.Schema != "systemd_run/v1" || configured.Binary == "" || len(configured.Properties) == 0 {
		return true, fmt.Errorf("invalid launcher configuration")
	}
	files, err := runfiles.New()
	if err != nil {
		return true, fmt.Errorf("resolve application runfiles: %w", err)
	}
	binary, err := files.Rlocation(configured.Binary)
	if err != nil {
		return true, fmt.Errorf("resolve configured executable: %w", err)
	}
	environment := files.Env()
	for _, assignment := range environment {
		name, value, _ := strings.Cut(assignment, "=")
		if err := os.Setenv(name, value); err != nil {
			return true, fmt.Errorf("set child runfiles environment %q: %w", name, err)
		}
	}
	arguments := append([]string{binary}, configured.Arguments...)
	arguments = append(arguments, os.Args[1:]...)
	if err := launch(configured.Properties, arguments); err != nil {
		return true, fmt.Errorf("launch configured command: %w", err)
	}
	return true, nil
}

// launch validates assignments and replaces the wrapper with systemd-run.
func launch(properties, arguments []string) error {
	for _, property := range properties {
		name, value, found := strings.Cut(property, "=")
		if !found || name == "" || value == "" {
			return fmt.Errorf("property must be a nonempty NAME=VALUE assignment")
		}
	}
	executable, err := exec.LookPath("systemd-run")
	if err != nil {
		return fmt.Errorf("find host systemd-run: %w", err)
	}
	command := []string{executable, "--user", "--scope", "--expand-environment=no"}
	for _, property := range properties {
		command = append(command, "--property", property)
	}
	command = append(command, "--")
	command = append(command, arguments...)
	if err := syscall.Exec(executable, command, os.Environ()); err != nil {
		return fmt.Errorf("execute host systemd-run: %w", err)
	}
	return nil
}
