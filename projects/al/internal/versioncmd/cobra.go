package versioncmd

import (
	"fmt"
	"os"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/internal/versioning"
	"github.com/spf13/cobra"
)

// versionOptions holds shared checkout and version-selection flags.
type versionOptions struct {
	// repository identifies the local Git checkout.
	repository string
	// trunk names the repository's trunk branch.
	trunk string
	// channel disambiguates exact release and nightly tags.
	channel string
	// release identifies a detached release-line checkout.
	release string
}

// NewCommand exposes native Cobra commands for version operations.
func NewCommand() *cobra.Command {
	return newCommand(time.Now)
}

// NewReleaseCreateCommand prepares local weekly release refs without changing HEAD.
// The CLI mounts this command under release; versioning owns its Git transaction.
func NewReleaseCreateCommand() *cobra.Command {
	options := &versionOptions{channel: "auto"}
	cmd := newMutationCommand(options, time.Now, "create")
	cmd.Short = "Create a local release branch and initial tag without changing HEAD"
	options.addRepositoryFlags(cmd)
	return cmd
}

// newCommand constructs version commands with an explicit clock for ref preparation.
func newCommand(now func() time.Time) *cobra.Command {
	options := &versionOptions{}
	root := &cobra.Command{
		Use: "version", Short: "Inspect, prepare, publish, and stamp repository versions",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if options.channel != "auto" && options.channel != "release" && options.channel != "nightly" {
				return fmt.Errorf("invalid channel %q; want auto, release, or nightly", options.channel)
			}
			if options.release != "" {
				if _, err := versioning.ParseCalendar(options.release); err != nil {
					return fmt.Errorf("invalid release context: %w", err)
				}
			}
			return nil
		},
	}
	options.addRepositoryFlags(root)
	flags := root.PersistentFlags()
	flags.StringVar(&options.channel, "channel", "auto", "Version channel: auto, release, or nightly")
	flags.StringVar(&options.release, "release", "", "Detached release line (YYYY.W)")
	root.AddCommand(newShowCommand(options), newMutationCommand(options, now, "nightly-tag"), newMutationCommand(options, now, "release-tag"), newPlanCommand(options), newPublishCommand(options))
	root.AddCommand(&cobra.Command{
		Use: "bazel-status", Short: "Emit deterministic Bazel workspace-status keys", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			status, err := options.repositoryFor(cmd).BazelStatus()
			if err != nil {
				return fmt.Errorf("calculate workspace status: %w", err)
			}
			if _, err := fmt.Fprint(cmd.OutOrStdout(), status); err != nil {
				return fmt.Errorf("write workspace status: %w", err)
			}
			return nil
		},
	})
	return root
}

// addRepositoryFlags shares checkout selection between version operations and release creation.
func (options *versionOptions) addRepositoryFlags(cmd *cobra.Command) {
	repository := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if repository == "" {
		repository = "."
	}
	flags := cmd.PersistentFlags()
	flags.StringVar(&options.repository, "repository", repository, "Git repository path")
	flags.StringVar(&options.trunk, "trunk", "master", "Trunk branch name")
}

// repositoryFor binds all Git subprocesses to the selected command's lifetime.
func (options *versionOptions) repositoryFor(cmd *cobra.Command) versioning.Repository {
	return versioning.Repository{
		Git:         versioning.ExecRunner{Directory: options.repository, Context: cmd.Context()},
		TrunkBranch: options.trunk, RequestedChannel: options.channel, RequestedRelease: options.release,
	}
}
