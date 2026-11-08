package main

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"git.alwaldend.com/alwaldend/src/projects/al/plugins/vault/vault_proto"
	"github.com/hashicorp/vault/api"
)

// opFetcher performs the configured Vault operation with cancellation.
type opFetcher struct {
	// vault owns lazy clients for the selected Vault identities.
	vault *al.VaultStore
}

// newOpFetcher binds a lazy Vault store to operation fetches.
func newOpFetcher(vault *al.VaultStore) *opFetcher {
	return &opFetcher{vault}
}

var _ resourceFetcher = (*opFetcher)(nil)

// String identifies the fetcher without exposing configuration values.
func (self *opFetcher) String() string {
	return "com.alwaldend.src.projects.al.plugins.vault.opFetcher"
}

// Get evaluates one resource using the supplied dependency results.
func (self *opFetcher) Get(ctx context.Context, r *vault_proto.Resource, d []*injector.ResourceResult) (*injector.ResourceResult, error) {
	op := r.GetOp()
	if op == nil {
		return nil, fmt.Errorf("missing op config")
	}
	client, err := self.vault.Client(ctx, r.VaultConn, r.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create vault client for op %s: %w", r.Name, al.SanitizeVaultError(err))
	}
	data := map[string]any{}
	for key, value := range op.Data {
		data[key] = value
	}
	var resp *api.Secret
	logical := client.Client.Logical()
	switch op.Method {
	case "", "read":
		resp, err = logical.ReadWithContext(ctx, op.Path)
		if err != nil {
			return nil, fmt.Errorf("read error: %w", al.SanitizeVaultError(err))
		}
	case "write":
		resp, err = logical.WriteWithContext(ctx, op.Path, data)
		if err != nil {
			return nil, fmt.Errorf("write error: %w", al.SanitizeVaultError(err))
		}
	default:
		return nil, fmt.Errorf("invalid method %s", op.Method)
	}
	res := &injector.ResourceResult{Name: r.Name}
	if resp != nil {
		res.Data = resp.Data
	}
	return res, nil
}
