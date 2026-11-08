package main

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"git.alwaldend.com/alwaldend/src/projects/al/plugins/vault/vault_proto"
)

// secretFetcher retrieves a KV v2 resource with the selected Vault identity.
type secretFetcher struct {
	// vault owns lazy clients for the selected Vault identities.
	vault *al.VaultStore
}

// newSecretFetcher binds a lazy Vault store to KV fetches.
func newSecretFetcher(vault *al.VaultStore) *secretFetcher {
	return &secretFetcher{vault}
}

var _ resourceFetcher = (*secretFetcher)(nil)

// String identifies the fetcher without exposing configuration values.
func (self *secretFetcher) String() string {
	return "com.alwaldend.src.projects.al.plugins.vault.secretFetcher"
}

// Get evaluates one resource using the supplied dependency results.
func (self *secretFetcher) Get(ctx context.Context, r *vault_proto.Resource, d []*injector.ResourceResult) (*injector.ResourceResult, error) {
	secret := r.GetKv()
	if secret == nil {
		return nil, fmt.Errorf("missing kv config")
	}
	if secret.Mount == "" {
		return nil, fmt.Errorf("missing kv mount")
	}
	if secret.Path == "" {
		return nil, fmt.Errorf("missing kv path")
	}
	client, err := self.vault.Client(ctx, r.VaultConn, r.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not get a client for the secret: %w", al.SanitizeVaultError(err))
	}
	data, err := client.Client.KVv2(secret.Mount).Get(ctx, secret.Path)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret: %w", al.SanitizeVaultError(err))
	}
	res := &injector.ResourceResult{
		Name: r.Name,
		Data: data.Data,
	}
	return res, nil
}
