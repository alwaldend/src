package main

import (
	"context"
	"fmt"
	"io"

	"git.alwaldend.com/src/tools/release/contracts"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
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
		return fmt.Errorf("could not execute commands: %w", err)
	}
	return nil
}

func newRootCommand(
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "release",
		Short:         "Release tool",
		Long:          "Release tool",
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmdGen, err := newGenCommand()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(cmdGen)
	return cmd, nil
}

func newGenCommand() (*cobra.Command, error) {
	generator := NewGenerator()
	release := &contracts.Release{}
	var output string
	var items []string
	marshalOptions := &protojson.MarshalOptions{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  "Generate manifest files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generator.Generate(release, items, output, marshalOptions)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&release.Tags, "tag", nil, "Release tags")
	flags.StringArrayVar(&items, "item", nil, "Add an item to the manifest (<type>:<path>, for example file:file.txt)")
	flags.StringVar(&output, "output", "", "Output path")
	flags.StringVar(&release.Title, "title", "", "Release title")
	flags.StringVar(&marshalOptions.Indent, "indent", "    ", "Json indent")
	cmd.MarkFlagsOneRequired("output")
	return cmd, nil
}
