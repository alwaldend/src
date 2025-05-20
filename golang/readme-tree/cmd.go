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
	config := &Config{}
	outputter := NewOutputter(config, stdout)
	parser := NewParse(config, outputter)

	cmd := &cobra.Command{
		Use:           "readme-tree",
		Short:         "Readme tree generator",
		Long:          "Readme tree generator",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if config.ChangeDir != "" {
				return os.Chdir(config.ChangeDir)
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(
		&config.ChangeDir,
		"change-dir",
		"C",
		"",
		"Change directory to the specified before running commands",
	)

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)

	cmdParse, err := NewParseCommand(parser, cmd)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(cmdParse)
	return cmd, nil
}

func NewParseCommand(parse *Parser, parent *cobra.Command) (*cobra.Command, error) {
	config := &ParseConfig{}
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse readme files",
		Long:  "Parse readme files",
		RunE: func(cmd *cobra.Command, args []string) error {
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}
			config.RootPath = workingDir
			return parse.ParsePathsAndOutput(args, config)
		},
	}
	flags := cmd.PersistentFlags()
	flags.VarP(
		&config.OutputType,
		"output-type",
		"o",
		"Output type",
	)
	flags.StringVarP(
		&config.ReadmeName,
		"readme-name",
		"n",
		"README.md",
		"Filename to process",
	)
	flags.BoolVarP(
		&config.UseGit,
		"git",
		"g",
		false,
		"If set, use git to get files",
	)
	return cmd, nil
}
