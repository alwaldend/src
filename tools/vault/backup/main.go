package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/pkg/vaultcmd"
	"github.com/spf13/cobra"
)

// main configures the snapshot command and reports backup failures.
func main() {
	var rclone, vault string
	command := &cobra.Command{
		Use: "backup", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := backup(cmd.Context(), vault, rclone, os.Getenv("RCLONE_S3_BUCKET")); err != nil {
				return fmt.Errorf("back up Vault: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&rclone, "rclone", "", "Rclone executable")
	command.Flags().StringVar(&vault, "vault", "", "Vault executable")
	if err := vaultcmd.Execute(command, "vault", "rclone"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// backup stages a snapshot in a private directory and transfers it to the configured bucket.
func backup(ctx context.Context, vault, rclone, bucket string) (result error) {
	if bucket == "" {
		return errors.New("RCLONE_S3_BUCKET is required")
	}
	tmp, err := os.MkdirTemp("", "infra_dc1_vault_backup_*")
	if err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmp); err != nil {
			result = errors.Join(result, fmt.Errorf("remove snapshot directory: %w", err))
		}
	}()
	name := fmt.Sprintf("vault-backup-%s.snap", time.Now().Format(time.RFC3339))
	path := filepath.Join(tmp, name)
	if err := al.RunCommand(al.CommandArgs{Ctx: ctx, Name: vault, Args: []string{"operator", "raft", "snapshot", "save", path}}); err != nil {
		return fmt.Errorf("save Raft snapshot: %w", err)
	}
	if err := al.RunCommand(al.CommandArgs{Ctx: ctx, Name: rclone, Args: []string{"copyto", "-v", "--check-first", "--s3-no-check-bucket", path, fmt.Sprintf("remote:%s/vault/%s", bucket, name)}}); err != nil {
		return fmt.Errorf("upload snapshot: %w", err)
	}
	return nil
}
