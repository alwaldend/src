package release

import (
	"context"
	"fmt"
	"io"
	"time"

	gitLib "git.alwaldend.com/alwaldend/src/projects/al/internal/gitinfo"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// newRootCommand assembles artifact operations and optional website commands.
func newRootCommand(
	ctx context.Context,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (*cobra.Command, error) {
	cmd := newCommandGroup("release", "Generate, publish, and inspect project releases")
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetIn(stdin)
	cmd.SetArgs(args)
	cmd.SetContext(ctx)
	cmdGen, err := newGenCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("create release subcommand: %w", err)
	}
	cmd.AddCommand(
		cmdGen,
		newOperationCommand(ctx, "publish", "upload", "Publish immutable artifact releases"),
		newOperationCommand(ctx, "delete", "delete", "Delete an inactive release and its extracted website"),
		newWebsiteCommand(ctx),
		newObservationCommand(ctx, false),
		newObservationCommand(ctx, true),
	)
	return cmd, nil
}

// newCommandGroup shows available operations while rejecting unknown input.
func newCommandGroup(name, summary string) *cobra.Command {
	return &cobra.Command{
		Use:           name,
		Short:         summary,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmd.Help(); err != nil {
				return fmt.Errorf("show %s commands: %w", cmd.Name(), err)
			}
			return nil
		},
	}
}

// newWebsiteCommand exposes optional website operations within a project release.
func newWebsiteCommand(ctx context.Context) *cobra.Command {
	cmd := newCommandGroup("website", "Extract and activate a published project website")
	cmd.AddCommand(
		newOperationCommand(ctx, "extract", "extract", "Extract and retain a published website archive"),
		newOperationCommand(ctx, "activate", "link", "Activate an extracted website version"),
	)
	return cmd
}

// newOperationCommand binds one public command to one fixed publication operation.
// Website previews reuse the offline planner and never create mutation receipts.
func newOperationCommand(ctx context.Context, name, step, summary string) *cobra.Command {
	deployer := NewDeployer()
	opts := &DeployerOpts{Ctx: ctx, SSHSteps: []string{step}}
	dryRun := false
	cmd := &cobra.Command{
		Use:   name,
		Short: summary,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx, opts.Stdout, opts.Stderr = cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr()
			if dryRun {
				if cmd.Flags().Changed("receipt") {
					return fmt.Errorf("--receipt cannot be combined with --dry-run")
				}
				if err := reportArtifacts(opts, false); err != nil {
					return fmt.Errorf("preview %s: %w", cmd.CommandPath(), err)
				}
				return nil
			}
			if err := deployer.Deploy(opts); err != nil {
				return fmt.Errorf("%s releases: %w", cmd.Name(), err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Receipt, "receipt", "", "Atomically checkpoint artifact operation evidence to this local file")
	if step == "extract" || step == "link" {
		cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview the website operation without contacting the destination")
	}
	if step == "extract" {
		cmd.Flags().StringVar(&opts.SiteArchive, "site_archive", "", "Published website archive filename")
	}
	publicationFlags(cmd, opts)
	return cmd
}

// newObservationCommand adapts offline preview or read-only destination inspection.
func newObservationCommand(ctx context.Context, inspect bool) *cobra.Command {
	opts := &DeployerOpts{Ctx: ctx, SSHSteps: []string{"upload"}}
	cmd := &cobra.Command{Use: "plan", Short: "Preview artifact publication metadata without contacting the destination", Args: cobra.NoArgs}
	if inspect {
		cmd.Use, cmd.Short = "inspect", "Verify local payloads and observe the SSH destination without modifying it"
	}
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		opts.Ctx, opts.Stdout, opts.Stderr = cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr()
		if err := reportArtifacts(opts, inspect); err != nil {
			return fmt.Errorf("%s artifacts: %w", cmd.Name(), err)
		}
		return nil
	}
	publicationFlags(cmd, opts)
	return cmd
}

// publicationFlags keeps destination selection identical for mutations and readers.
func publicationFlags(cmd *cobra.Command, opts *DeployerOpts) {
	flags := cmd.PersistentFlags()
	flags.StringArrayVar(&opts.Releases, "release_dir", nil, "Path to a release directory")
	flags.StringVar(&opts.Environment, "environment", "", "Explicit SSH environment")
	flags.StringVar(&opts.SSHDeployment, "ssh_deployment", "", "SSH deployment JSON selecting the release destination")
	flags.StringVar(&opts.SSHHost, "ssh_host", "", "SSH destination host override")
	flags.StringVar(&opts.SSHUser, "ssh_user", "", "Administrator SSH login override")
	flags.StringVar(&opts.SSHRoot, "ssh_root", "", "Remote content root override")
	flags.StringVar(&opts.PublishUser, "publish_user", "", "Remote content account override")
	flags.StringVar(&opts.Project, "project", "", "Public project name override")
	flags.StringVar(&opts.Namespace, "namespace", "", "Public artifact namespace: projects or third_party")
	flags.StringVar(&opts.SSHPath, "ssh_path", "ssh", "OpenSSH executable")
	flags.StringVar(&opts.RsyncPath, "rsync_path", "rsync", "Local rsync executable")
	flags.StringVar(&opts.SSHConfig, "ssh_config", "", "OpenSSH configuration file")
	flags.StringVar(&opts.SSHIdentity, "ssh_identity", "", "OpenSSH identity file")
	flags.BoolVar(&opts.SSHDisableAgent, "ssh_disable_agent", false, "Disable agent authentication and use only configured identities")
	flags.StringVar(&opts.SSHCertificate, "ssh_certificate", "", "OpenSSH identity certificate file")
	flags.StringVar(&opts.SSHKnownHosts, "ssh_known_hosts", "", "Known-hosts file (default: OpenSSH configuration)")
	flags.IntVar(&opts.SSHPort, "ssh_port", 0, "SSH port override")
	flags.DurationVar(&opts.Timeout, "timeout", 30*time.Minute, "Maximum time for each SSH release operation")
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

// NewCommand constructs release adapters for the shared AL command tree.
func NewCommand(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) (*cobra.Command, error) {
	command, err := newRootCommand(ctx, nil, stdin, stdout, stderr)
	if err != nil {
		return nil, fmt.Errorf("construct release command: %w", err)
	}
	return command, nil
}
