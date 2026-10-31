package versioncmd

import (
	"fmt"
	"os"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/versioning"
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
	repository := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	if repository == "" {
		repository = "."
	}
	flags := root.PersistentFlags()
	flags.StringVar(&options.repository, "repository", repository, "Git repository path")
	flags.StringVar(&options.repository, "repo", repository, "Alias for --repository")
	flags.StringVar(&options.trunk, "trunk", "master", "Trunk branch name")
	flags.StringVar(&options.channel, "channel", "auto", "Version channel: auto, release, or nightly")
	flags.StringVar(&options.release, "release", "", "Detached release line (YYYY.W)")
	root.AddCommand(newShowCommand(options), newMutationCommand(options, now, "nightly-tag"), newMutationCommand(options, now, "release-start"), newMutationCommand(options, now, "release-tag"), newPlanCommand(options), newPublishCommand(options))
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
	root.AddCommand(&cobra.Command{
		Use: "bazel -- ARGS", Short: "Run Bazel with repository version stamping",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || cmd.ArgsLenAtDash() != 0 {
				return fmt.Errorf("bazel requires -- followed by a Bazel command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := launchBazel(cmd.Context(), options.repository, options.trunk, options.channel, options.release, args); err != nil {
				return fmt.Errorf("run versioned build: %w", err)
			}
			return nil
		},
	})
	return root
}

// repositoryFor binds all Git subprocesses to the selected command's lifetime.
func (options *versionOptions) repositoryFor(cmd *cobra.Command) versioning.Repository {
	return versioning.Repository{
		Git:         versioning.ExecRunner{Directory: options.repository, Context: cmd.Context()},
		TrunkBranch: options.trunk, RequestedChannel: options.channel, RequestedRelease: options.release,
	}
}
