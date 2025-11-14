package main

import (
	"context"
	"fmt"
	"io"

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
	var output string
	var items []string
	var manifests []string
	marshalOptions := &protojson.MarshalOptions{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  "Merge several manifests into one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generator.Generate(items, manifests, output, marshalOptions)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&items, "item", nil, "Path to a file with ReleaseItem json")
	flags.StringArrayVar(&manifests, "manifest", nil, "Path to a file with Release json")
	flags.StringVar(&output, "output", "", "Output path")
	flags.StringVar(&marshalOptions.Indent, "indent", "    ", "Json indent")
	cmd.MarkFlagsOneRequired("output")
	return cmd, nil
}
