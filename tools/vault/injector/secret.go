package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type SecretFetcher struct {
	vault *al.VaultStore
}

func NewSecretFetcher(vault *al.VaultStore) *SecretFetcher {
	return &SecretFetcher{vault}
}

var _ ResourceFetcher = (*SecretFetcher)(nil)

func (self *SecretFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.SecretFetcher"
}

func (self *SecretFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error) {
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
		return nil, fmt.Errorf("could not get a client for the secret: %w", err)
	}
	data, err := client.Client.KVv2(secret.Mount).Get(ctx, secret.Path)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret: %w", err)
	}
	res := &ResourceResult{
		Name: r.Name,
		Data: data.Data,
	}
	return res, nil
}
