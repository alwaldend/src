package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// execute binds the caller’s arguments and streams to a fresh command.
func execute(
	ctx context.Context,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	root, err := newRootCommand(args, stdin, stdout, stderr)
	if err != nil {
		return fmt.Errorf("configure templating command: %w", err)
	}
	err = root.ExecuteContext(ctx)
	if err != nil {
		return fmt.Errorf("run templating command: %w", err)
	}
	return nil
}

// newRootCommand defines the working directory shared by template operations.
func newRootCommand(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	var changeDir string
	cmd := &cobra.Command{
		Use:           "template-files",
		Short:         "Tool to template file",
		Long:          "Tool to template file",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if changeDir != "" {
				if err := os.Chdir(changeDir); err != nil {
					return fmt.Errorf("change working directory to %q: %w", changeDir, err)
				}
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVarP(
		&changeDir,
		"change-dir",
		"C",
		"",
		"Change directory to the specified before running commands",
	)

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)

	cmdParse, err := newTemplateCommand()
	if err != nil {
		return nil, fmt.Errorf("configure template subcommand: %w", err)
	}
	cmd.AddCommand(cmdParse)
	return cmd, nil
}

// newTemplateCommand binds the rendering inputs and output path.
func newTemplateCommand() (*cobra.Command, error) {
	var dataPaths []string
	var templatePath string
	var outputPath string
	var extension string
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Load data files (-d), template the template (-t) and write it to the output (-o)",
		Long:  "Load data files (-d), template the template (-t) and write it to the output (-o)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := templateFiles(cmd.Context(), dataPaths, templatePath, outputPath, extension); err != nil {
				return fmt.Errorf("render %q: %w", templatePath, err)
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringSliceVarP(&dataPaths, "data", "d", []string{}, "Paths to data files")
	flags.StringVarP(&outputPath, "output", "o", "", "Output path")
	flags.StringVarP(&templatePath, "template", "t", "", "Template path")
	flags.StringVarP(&extension, "extension", "e", "", "Override data extension")
	for _, name := range []string{"template", "output"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return nil, fmt.Errorf("require flag %q: %w", name, err)
		}
	}
	return cmd, nil
}
