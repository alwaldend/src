package main

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	root, err := NewRootCommand(args, stdin, stdout, stderr)
	if err != nil {
		return err
	}
	err = root.Execute()
	if err != nil {
		return err
	}
	return nil
}

func NewRootCommand(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	config := &State{Stdout: stdout, Stderr: stderr, Stdin: stdin}
	changeDir := ""
	readmeName := ""
	useGit := false

	cmdRoot := &cobra.Command{
		Use:           "readme-tree",
		Short:         "Readme tree generator",
		Long:          "Readme tree generator",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if changeDir != "" {
				return os.Chdir(changeDir)
			}
			return nil
		},
	}
	cmdParse := &cobra.Command{
		Use:   "parse",
		Short: "Parse readme files",
		Long:  "Parse readme files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}
			return NewParse(config, readmeName, useGit).ParsePathsAndOutput(args, workingDir)
		},
	}

	cmdRoot.PersistentFlags().StringVarP(
		&changeDir,
		"change-dir",
		"C",
		"",
		"Change directory to the specified before running commands",
	)
	cmdRoot.PersistentFlags().StringVarP(
		&readmeName,
		"readme-name",
		"n",
		"README.md",
		"Filename to process",
	)
	cmdRoot.PersistentFlags().BoolVarP(
		&useGit,
		"git",
		"g",
		false,
		"If set, use git to get files",
	)

	cmdRoot.SetOut(stdout)
	cmdRoot.SetErr(stderr)
	cmdRoot.SetIn(stdin)
	cmdRoot.SetArgs(args)
	cmdRoot.AddCommand(cmdParse)
	return cmdRoot, nil
}
