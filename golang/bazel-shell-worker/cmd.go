package main

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func Execute(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	root, err := newRootCommand(args, stdin, stdout, stderr)
	if err != nil {
		return fmt.Errorf("could not create commands: %w", err)
	}
	err = root.Execute()
	if err != nil {
		return err
	}
	return nil
}

func newRootCommand(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "bazel-genrule-worker",
		Short:         "Bazel worker that runs shell commands",
		Long:          "Bazel worker that runs shell commands",
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmdParse, err := newRunCommand(stdout, stdin)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(cmdParse)
	return cmd, nil
}

func newRunCommand(stdout io.Writer, stdin io.Reader) (*cobra.Command, error) {
	protocol := NewWorkerProtocol(stdout, stdin)
	worker := NewWorker(protocol)
	persist := false
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run",
		Long:  "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			worker.Run(persist)
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.BoolVarP(&persist, "persistent_worker", "p", false, "If set, the worker will persist")
	return cmd, nil
}
