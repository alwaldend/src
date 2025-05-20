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
	state := &State{
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  stdin,
		Config: &Config{
			Parse: &ParseConfig{},
		},
	}
	outputter := NewOutputter(state)
	parser := NewParse(state, outputter)

	cmd := &cobra.Command{
		Use:           "readme-tree",
		Short:         "Readme tree generator",
		Long:          "Readme tree generator",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if state.Config.ChangeDir != "" {
				return os.Chdir(state.Config.ChangeDir)
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(
		&state.Config.ChangeDir,
		"change-dir",
		"C",
		"",
		"Change directory to the specified before running commands",
	)
	flags.StringVarP(
		&state.Config.ReadmeName,
		"readme-name",
		"n",
		"README.md",
		"Filename to process",
	)
	flags.BoolVarP(
		&state.Config.UseGit,
		"git",
		"g",
		false,
		"If set, use git to get files",
	)

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)

	err := AddParseCommand(parser, state, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func AddParseCommand(parse *Parser, state *State, parent *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse readme files",
		Long:  "Parse readme files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}
			return parse.ParsePathsAndOutput(args, workingDir)
		},
	}
	flags := cmd.PersistentFlags()
	flags.VarP(
		&state.Config.Parse.OutputType,
		"output-type",
		"o",
		"Output type",
	)
	parent.AddCommand(cmd)
	return nil
}
