package release

import (
	"context"
	"fmt"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/gitinfo"
	"github.com/spf13/cobra"
)

// newGitInfoCommand exposes the merged checkout metadata generator.
func newGitInfoCommand() *cobra.Command {
	opts := &gitinfo.GitInfoGenerateOptions{}
	var timeout time.Duration
	cmd := &cobra.Command{
		Use: "git-info", Short: "Generate repository history metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			if err := gitinfo.NewGitInfo().Generate(ctx, opts); err != nil {
				return fmt.Errorf("generate Git information: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Output, "output", "", "Output JSON path")
	cmd.Flags().StringVar(&opts.GitRoot, "git_root", "", "Checkout root")
	cmd.Flags().DurationVar(&timeout, "timeout", time.Minute, "Metadata generation deadline")
	cmd.MarkFlagsRequiredTogether("output", "git_root")
	cmd.MarkFlagsOneRequired("output", "git_root")
	return cmd
}
