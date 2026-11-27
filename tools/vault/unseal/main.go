package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/pkg/vaultcmd"
	"github.com/spf13/cobra"
)

// main submits one decrypted share to each requested Vault server and reports failures together.
func main() {
	var vault string
	command := &cobra.Command{
		Use: "unseal [address...]", Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, addresses []string) error {
			key, err := vaultcmd.ReadKey(cmd.Context(), os.Getenv("UNSEAL_KEY"))
			if err != nil {
				return fmt.Errorf("read unseal key: %w", err)
			}
			var failures []error
			for i, address := range addresses {
				if err := cmd.Context().Err(); err != nil {
					failures = append(failures, fmt.Errorf("unseal remaining servers: %w", err))
					break
				}
				if err := al.RunCommand(al.CommandArgs{
					Ctx: cmd.Context(), Name: vault,
					Args: []string{"write", "-address", address, "sys/unseal", "key=-"}, Stdin: bytes.NewReader(key),
				}); err != nil {
					failures = append(failures, fmt.Errorf("unseal server %d: %w", i+1, err))
				}
			}
			return errors.Join(failures...)
		},
	}
	command.Flags().StringVar(&vault, "vault", "", "Vault executable")
	if err := vaultcmd.Execute(command, "vault"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
