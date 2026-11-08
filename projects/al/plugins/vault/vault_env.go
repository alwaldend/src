package main

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"git.alwaldend.com/alwaldend/src/projects/al/plugins/vault/vault_proto"
)

// vaultEnvFetcher exports a selected Vault connection and identity.
type vaultEnvFetcher struct {
	// config borrows the configured Vault connections and authentication methods.
	config *al_proto.Config
	// vault owns lazy clients for the selected Vault identities.
	vault *al.VaultStore
}

// newVaultEnvFetcher borrows Vault configuration and client ownership.
func newVaultEnvFetcher(config *al_proto.Config, vault *al.VaultStore) *vaultEnvFetcher {
	return &vaultEnvFetcher{config, vault}
}

var _ resourceFetcher = (*vaultEnvFetcher)(nil)

// String identifies the fetcher without exposing configuration values.
func (self *vaultEnvFetcher) String() string {
	return "com.alwaldend.src.projects.al.plugins.vault.vaultEnvFetcher"
}

// Get exports the selected Vault connection, TLS paths and token.
// https://developer.hashicorp.com/vault/docs/commands#configure-environment-variables
func (self *vaultEnvFetcher) Get(ctx context.Context, r *vault_proto.Resource, d []*injector.ResourceResult) (*injector.ResourceResult, error) {
	vaultConn, vaultAuth := r.GetVaultEnv().Conn, r.GetVaultEnv().Auth
	if vaultConn == "" {
		vaultConn = al.VAULT_DEFAULT_NAME
	}
	if vaultAuth == "" {
		vaultAuth = al.VAULT_DEFAULT_NAME
	}
	vault, err := al.VaultByName(self.config, vaultConn)
	if err != nil {
		return nil, fmt.Errorf("missing vault connection: %w", err)
	}
	auth, err := al.VaultAuthByName(self.config, vaultAuth)
	if err != nil {
		return nil, fmt.Errorf("missing vault auth: %w", err)
	}
	tlsConfig, err := fp.Get(al.VaultTlsConfig(vault))
	if err != nil {
		return nil, fmt.Errorf("could not create tls config: %w", err)
	}
	res := &injector.ResourceResult{
		Env: map[string]string{
			"VAULT_ADDR":  vault.Config.Address,
			"VAULT_TOKEN": "",
		},
	}
	if tlsConfig.CACert != "" {
		res.Env["VAULT_CACERT"] = tlsConfig.CACert
	}
	if tlsConfig.ClientKey != "" {
		res.Env["VAULT_CLIENT_KEY"] = tlsConfig.ClientKey
	}
	if tlsConfig.ClientCert != "" {
		res.Env["VAULT_CLIENT_CERT"] = tlsConfig.ClientCert
	}
	if !auth.NoAuth {
		client, err := fp.Get(self.vault.Client(ctx, vaultConn, vaultAuth))
		if err != nil {
			return nil, fmt.Errorf("could not create client for %s/%s: %w", vaultConn, vaultAuth, err)
		}
		res.Env["VAULT_TOKEN"] = client.Client.Token()
	}
	return res, nil
}
