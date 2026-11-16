package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// command adapts the apply-time environment to the credential bootstrap.
func command() *cobra.Command {
	root := &cobra.Command{
		Use:           "runner_registration",
		Short:         "Store one Forgejo runner registration in Vault",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(&cobra.Command{
		Use:   "ensure",
		Short: "Reuse cached credentials or register a runner once",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
			defer cancel()
			cfg := configuration{
				instance: os.Getenv("FORGEJO_HOST"), repository: os.Getenv("FORGEJO_RUNNER_REPOSITORY"),
				name: os.Getenv("FORGEJO_RUNNER_NAME"), mount: os.Getenv("FORGEJO_RUNNER_VAULT_MOUNT"),
				path: os.Getenv("FORGEJO_RUNNER_VAULT_PATH"),
			}
			cached, err := cfg.ensure(ctx)
			if err != nil {
				return fmt.Errorf("runner credential bootstrap: %w", err)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Runner credentials stored in Vault (cached=%t).\n", cached)
			if err != nil {
				return fmt.Errorf("write completion report: %w", err)
			}
			return nil
		},
	})
	return root
}

// main propagates cancellation and reports only sanitized operational errors.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
