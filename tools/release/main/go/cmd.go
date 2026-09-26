package main

import (
	"context"
	"fmt"
	"io"
	"time"

	gitLib "git.alwaldend.com/alwaldend/src/tools/git/main/go"
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
		return nil, fmt.Errorf("create release subcommand: %w", err)
	}
	cmdDeploy, err := newDeployCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("create release subcommand: %w", err)
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
			opts.Stdout, opts.Stderr = cmd.OutOrStdout(), cmd.ErrOrStderr()
			if err := deployer.Deploy(opts); err != nil {
				return fmt.Errorf("deploy releases: %w", err)
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringVar(&opts.SorasPath, "soras_path", "", "Path to the soras binary")
	flags.StringArrayVar(&opts.Releases, "release_dir", nil, "Path to a release directory")
	flags.StringVar(&opts.Environment, "environment", "", "Explicit SSH environment (OCI when omitted)")
	flags.StringVar(&opts.SSHDeployment, "ssh_deployment", "", "SSH deployment JSON; publish all release files using it")
	flags.StringVar(&opts.SSHHost, "ssh_host", "", "SSH host; publish all release files using this host")
	flags.StringVar(&opts.SSHUser, "ssh_user", "", "Administrator SSH login override")
	flags.StringVar(&opts.SSHRoot, "ssh_root", "", "Remote content root override")
	flags.StringVar(&opts.PublishUser, "publish_user", "", "Remote content account override")
	flags.StringVar(&opts.Project, "project", "", "Public project name override")
	flags.StringVar(&opts.SiteArchive, "site_archive", "", "Published archive filename for the extract step")
	flags.StringSliceVar(&opts.SSHSteps, "steps", []string{"upload"}, "SSH steps: ordered subset of upload,extract,link")
	flags.StringVar(&opts.SSHPath, "ssh_path", "ssh", "OpenSSH executable")
	flags.StringVar(&opts.RsyncPath, "rsync_path", "rsync", "Local rsync executable")
	flags.StringVar(&opts.SSHConfig, "ssh_config", "", "OpenSSH configuration file")
	flags.StringVar(&opts.SSHIdentity, "ssh_identity", "", "OpenSSH identity file")
	flags.StringVar(&opts.SSHKnownHosts, "ssh_known_hosts", "", "Known-hosts file (default: OpenSSH configuration)")
	flags.IntVar(&opts.SSHPort, "ssh_port", 0, "SSH port override")
	flags.DurationVar(&opts.Timeout, "timeout", 30*time.Minute, "Maximum time for each SSH release deployment")
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
			if err := generator.Generate(opts); err != nil {
				return fmt.Errorf("generate release: %w", err)
			}
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&opts.AddFiles, "add_file", nil, "Path to a file that will be added to the manifest")
	flags.StringArrayVar(&opts.DeploymentInfo, "deployment", nil, "Path to a deployment info file")
	flags.StringArrayVar(&opts.MergeManifests, "merge_manifest", nil, "Path to a file with Release json")
	flags.StringArrayVar(&opts.OutputManifests, "output_manifest", nil, "Write the combined manifest to this path")
	flags.StringArrayVar(&opts.OutputReleasePages, "output_release_page", nil, "Write the release page to this path")
	flags.StringVar(&opts.OutputFileMode, "output_file_mode", "0444", "Create output files with this file mode")
	flags.StringVar(&opts.GitBundle, "git_bundle", "", "Git bundle with release history")
	flags.StringVar(&opts.Project, "project", "", "Release project subdirectory")
	flags.StringVar(&opts.VersionFile, "version_file", "", "Workspace status file containing STABLE_VERSION")
	flags.StringVar(&opts.OutputDir, "output_dir", "", "Write a portable release.json and files/ bundle to an empty directory")
	flags.StringVar(&opts.MarshalOptions.Indent, "indent", "    ", "Json indent")
	return cmd, nil
}
