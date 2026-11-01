package main

import (
	"context"
	"io"
	"os"

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
		return err
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
	var changeDir string
	cmd := &cobra.Command{
		Use:           "template-files",
		Short:         "Tool to template file",
		Long:          "Tool to template file",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if changeDir != "" {
				return os.Chdir(changeDir)
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
		return nil, err
	}
	cmd.AddCommand(cmdParse)
	return cmd, nil
}

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
			return NewTemplater().TemplateFiles(dataPaths, templatePath, outputPath, extension)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringSliceVarP(&dataPaths, "data", "d", []string{}, "Paths to data files")
	flags.StringVarP(&outputPath, "output", "o", "", "Output path")
	flags.StringVarP(&templatePath, "template", "t", "", "Template path")
	flags.StringVarP(&extension, "extension", "e", "", "Override data extension")
	return cmd, nil
}
