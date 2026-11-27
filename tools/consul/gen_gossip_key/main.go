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

// main generates a Consul key and writes it through Vault's stdin input.
func main() {
	var vault, consul, secret, mount string
	command := &cobra.Command{
		Use: "gen_gossip_key", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			key, err := vaultcmd.Output(cmd.Context(), consul, nil, "keygen")
			if err != nil {
				return fmt.Errorf("generate gossip key: %w", err)
			}
			key = bytes.TrimSpace(key)
			if len(key) == 0 {
				return errors.New("Consul returned an empty gossip key")
			}
			if err := al.RunCommand(al.CommandArgs{Ctx: cmd.Context(), Name: vault, Stdin: bytes.NewReader(key), Args: []string{"kv", "put", "-mount", mount, secret, "key=-"}}); err != nil {
				return fmt.Errorf("store gossip key: %w", err)
			}
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&vault, "vault", "", "Vault executable")
	flags.StringVar(&consul, "consul", "", "Consul executable")
	flags.StringVar(&secret, "secret", "", "Secret path")
	flags.StringVar(&mount, "secret_mount", "", "Secret mount")
	if err := vaultcmd.Execute(command, "vault", "consul", "secret", "secret_mount"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
