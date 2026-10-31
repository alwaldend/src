package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
)

// errFindings distinguishes unmatched results from an incomplete scan.
var errFindings = errors.New("unmatched secret findings")

// options contains the caller's tools, worktree, and immutable finding policy.
type options struct {
	// scanner is the TruffleHog executable.
	scanner string
	// git is an optional explicit Git executable.
	git string
	// repository is the Git worktree root for direct command use.
	repository string
	// marker locates a root file through Bazel runfiles.
	marker string
	// allowlist is the strict JSON policy path.
	allowlist string
	// allRefs includes every local ref instead of only HEAD history.
	allRefs bool
	// timeout bounds the scanner process.
	timeout time.Duration
}

// main translates cancellation and command errors into stable process statuses.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "trufflehog_scan:", err)
		if errors.Is(err, errFindings) {
			os.Exit(183)
		}
		os.Exit(1)
	}
}

// command adapts flags and rule-provided runfiles to the scan operation.
func command() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use: "trufflehog_scan", Short: "Scan Git history without credential verification",
		Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := opts.resolve(); err != nil {
				return fmt.Errorf("resolve scan inputs: %w", err)
			}
			if err := scan(cmd.Context(), opts, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
				return fmt.Errorf("scan repository: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.scanner, "scanner", "", "TruffleHog executable")
	cmd.Flags().StringVar(&opts.git, "git", "", "Explicit Git executable")
	cmd.Flags().StringVar(&opts.repository, "repository", "", "Git worktree root")
	cmd.Flags().StringVar(&opts.marker, "repository-marker", "", "File at the Git worktree root")
	cmd.Flags().StringVar(&opts.allowlist, "allowlist", "", "JSON finding allowlist")
	cmd.Flags().BoolVar(&opts.allRefs, "all-refs", os.Getenv("TRUFFLEHOG_SCAN_ALL_REFS") == "true", "Scan every local Git ref")
	cmd.Flags().DurationVar(&opts.timeout, "timeout", 15*time.Minute, "Scanner deadline")
	return cmd
}

// resolve supplies runfiles defaults and validates unambiguous command inputs.
func (o *options) resolve() error {
	for _, item := range []struct {
		// destination receives an explicit or runfiles-resolved path.
		destination *string
		// variable names the rule's runfiles key.
		variable string
	}{
		{&o.scanner, "TRUFFLEHOG_SCAN_SCANNER"},
		{&o.git, "TRUFFLEHOG_SCAN_GIT"},
		{&o.allowlist, "TRUFFLEHOG_SCAN_ALLOWLIST"},
		{&o.marker, "TRUFFLEHOG_SCAN_MARKER"},
	} {
		if *item.destination != "" || os.Getenv(item.variable) == "" {
			continue
		}
		resolved, err := runfiles.Rlocation(os.Getenv(item.variable))
		if err != nil {
			return fmt.Errorf("resolve %s runfile: %w", item.variable, err)
		}
		*item.destination = resolved
	}
	if o.scanner == "" || o.allowlist == "" || o.timeout <= 0 {
		return errors.New("scanner, allowlist, and a positive timeout are required")
	}
	if o.repository != "" && o.marker != "" {
		return errors.New("repository and repository-marker are mutually exclusive")
	}
	if o.marker != "" {
		resolved, err := filepath.EvalSymlinks(o.marker)
		if err != nil {
			return fmt.Errorf("resolve repository marker: %w", err)
		}
		o.repository = filepath.Dir(resolved)
	}
	if o.repository == "" {
		return errors.New("repository or repository-marker is required")
	}
	absolute, err := filepath.Abs(o.repository)
	if err != nil {
		return fmt.Errorf("resolve repository path: %w", err)
	}
	o.repository = absolute
	return nil
}
