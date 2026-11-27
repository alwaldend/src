package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/tools/go/pkg/sensitiveerr"
	"git.alwaldend.com/alwaldend/src/tools/vault/internal/privatefile"
	"git.alwaldend.com/alwaldend/src/tools/vault/pkg/vaultcmd"
	"github.com/spf13/cobra"
)

// main accepts either a new root-token operation or an existing nonce.
func main() {
	var vault, key, nonce string
	command := &cobra.Command{
		Use: "gen_root_token", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := generate(cmd.Context(), vault, key, nonce); err != nil {
				return fmt.Errorf("generate root token: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&vault, "vault", "", "Vault executable")
	command.Flags().StringVar(&key, "pgp_key", "", "Public PGP key in base64")
	command.Flags().StringVar(&nonce, "nonce", "", "Existing operation nonce")
	if err := vaultcmd.Execute(command, "vault"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// generate submits an unseal share and persists the completed token privately.
func generate(ctx context.Context, vault, key, nonce string) error {
	if nonce == "" && key == "" {
		return errors.New("--pgp_key is required to begin a new operation")
	}
	share, err := vaultcmd.ReadKey(ctx, "")
	if err != nil {
		return fmt.Errorf("read unseal share: %w", err)
	}
	if nonce == "" {
		output, err := vaultcmd.Output(ctx, vault, nil, "operator", "generate-root", "-init", "-pgp-key="+key, "-format", "json")
		if err != nil {
			return fmt.Errorf("initialize generation: %w", err)
		}
		var result struct {
			// Nonce identifies the generation operation created by Vault.
			Nonce string `json:"nonce"`
		}
		if err := json.Unmarshal(output, &result); err != nil {
			return sensitiveerr.WithContext("decode operation nonce", err)
		}
		nonce = result.Nonce
		if nonce == "" {
			return errors.New("Vault returned an empty operation nonce")
		}
	}
	output, err := vaultcmd.Output(ctx, vault, bytes.NewReader(share), "operator", "generate-root", "-nonce="+nonce, "-format", "json", "-")
	if err != nil {
		return fmt.Errorf("submit generation share: %w", err)
	}
	var response struct {
		// EncodedToken is the PGP-encrypted token returned when generation completes.
		EncodedToken string `json:"encoded_token"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return sensitiveerr.WithContext("decode generation response", err)
	}
	if response.EncodedToken == "" {
		return errors.New("Vault has not returned a completed root token")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(response.EncodedToken)
	if err != nil {
		return fmt.Errorf("decode encrypted token: %w", err)
	}
	token, err := vaultcmd.Output(ctx, "gpg", bytes.NewReader(ciphertext), "--decrypt")
	if err != nil {
		return fmt.Errorf("decrypt root token: %w", err)
	}
	if len(bytes.TrimSpace(token)) == 0 {
		return errors.New("empty decrypted token")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("locate token directory: %w", err)
	}
	if err := privatefile.Write(filepath.Join(home, ".vault-token"), token); err != nil {
		return fmt.Errorf("save root token: %w", err)
	}
	return nil
}
