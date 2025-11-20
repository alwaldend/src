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
	opts := &GenerateOpts{MarshalOptions: &protojson.MarshalOptions{}}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  "Merge several manifests into one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generator.Generate(opts)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&opts.MergeItems, "merge_item", nil, "Path to a file with ReleaseItem json")
	flags.StringArrayVar(&opts.MergeManifests, "merge_manifest", nil, "Path to a file with Release json")
	flags.StringVar(&opts.OutputManifest, "output_manifest", "", "Write the combined manifest to this path")
	flags.StringVar(&opts.OutputReleasePage, "output_release_page", "", "Write the release page to this path")
	flags.StringVar(&opts.OutputFileMode, "output_file_mode", "0444", "Create output files with this file mode")
	flags.StringVar(&opts.GitRoot, "git_root", "", "Directory with .git")
	flags.StringVar(&opts.MarshalOptions.Indent, "indent", "    ", "Json indent")
	return cmd, nil
}
