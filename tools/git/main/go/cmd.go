package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

func Execute(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	root, err := newRootCommand(ctx, args, stdin, stdout, stderr)
	if err != nil {
		return fmt.Errorf("could not create commands: %w", err)
	}
	err = root.Execute()
	if err != nil {
		return fmt.Errorf("could not execute commands: %w", err)
	}
	return nil
}

func newRootCommand(
	ctx context.Context,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "git",
		Short:         "Git tool",
		Long:          "Git tool",
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmdGen, err := newGitInfoCmd(ctx)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(cmdGen)
	return cmd, nil
}

func newGitInfoCmd(ctx context.Context) (*cobra.Command, error) {
	gitInfo := NewGitInfo()
	var output string
	var gitRoot string
	var timeout int
	cmd := &cobra.Command{
		Use:   "git_info",
		Short: "Generate GitInfo file",
		Long:  "Generate GitInfo",
		RunE: func(cmd *cobra.Command, args []string) error {
			curCtx := ctx
			if timeout != 0 {
				var cancel context.CancelFunc
				curCtx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
				defer cancel()
			}
			return gitInfo.Generate(curCtx, gitRoot, output)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVar(&output, "output", "", "Output path")
	flags.StringVar(&gitRoot, "git_root", "", "Directory with .git")
	flags.IntVar(&timeout, "timeout", 0, "Timeout in seconds (default: not timeout)")
	cmd.MarkFlagsOneRequired("output", "git_root")
	return cmd, nil
}
