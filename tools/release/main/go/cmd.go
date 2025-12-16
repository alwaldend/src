package main

import (
	"context"
	"fmt"
	"io"

	gitLib "git.alwaldend.com/src/tools/git/main/go"
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
		Use:           "release",
		Short:         "Release tool",
		Long:          "Release tool",
		SilenceErrors: true,
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmdGen, err := newGenCommand(ctx)
	if err != nil {
		return nil, err
	}
	cmdDeploy, err := newDeployCommand(ctx)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(cmdGen, cmdDeploy)
	return cmd, nil
}

func newDeployCommand(ctx context.Context) (*cobra.Command, error) {
	deployer := NewDeployer()
	opts := &DeployerOpts{Ctx: ctx}
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy",
		Long:  "Deploy releases",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deployer.Deploy(opts)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVar(&opts.SorasPath, "soras_path", "", "Path to the soras binary")
	flags.StringArrayVar(&opts.Releases, "release_dir", nil, "Path to a release directory")
	cmd.MarkFlagRequired("soras_path")
	return cmd, nil
}

func newGenCommand(ctx context.Context) (*cobra.Command, error) {
	gitInfo := &gitLib.GitInfo{}
	generator := NewGenerator(gitInfo)
	opts := &GenerateOpts{MarshalOptions: &protojson.MarshalOptions{}, Ctx: ctx}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  "Merge several manifests into one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generator.Generate(opts)
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&opts.AddFiles, "add_file", nil, "Path to a file that will be added to the manifest")
	flags.StringArrayVar(&opts.DeploymentInfo, "deployment", nil, "Path to a deployment info file")
	flags.StringArrayVar(&opts.MergeManifests, "merge_manifest", nil, "Path to a file with Release json")
	flags.StringArrayVar(&opts.OutputManifests, "output_manifest", nil, "Write the combined manifest to this path")
	flags.StringArrayVar(&opts.OutputReleasePages, "output_release_page", nil, "Write the release page to this path")
	flags.StringVar(&opts.OutputFileMode, "output_file_mode", "0444", "Create output files with this file mode")
	flags.StringVar(&opts.GitRoot, "git_root", "", "Directory with .git")
	flags.StringVar(&opts.MarshalOptions.Indent, "indent", "    ", "Json indent")
	return cmd, nil
}
