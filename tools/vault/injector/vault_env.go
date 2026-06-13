package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type VaultEnvFetcher struct {
	templater *Templater
	config    *al_proto.Config
	vault     *al.VaultStore
}

func NewVaultEnvFetcher(templater *Templater, config *al_proto.Config, vault *al.VaultStore) *VaultEnvFetcher {
	return &VaultEnvFetcher{templater, config, vault}
}

var _ ResourceFetcher = (*VaultEnvFetcher)(nil)

func (self *VaultEnvFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.VaultEnvFetcher"
}

// Create vault environment variables
// https://developer.hashicorp.com/vault/docs/commands#configure-environment-variables
func (self *VaultEnvFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error) {
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
	res := &ResourceResult{
		Env: map[string]string{
			"VAULT_ADDR": vault.Config.Address,
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
