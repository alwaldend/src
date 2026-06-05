package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type secret struct {
	vault *al.VaultStore
}

var _ ResourceFetcher = (*secret)(nil)

func (self *secret) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) fp.Either[*ResourceResult] {
	secret := r.GetKv()
	if secret == nil {
		return fp.Left[*ResourceResult](fmt.Errorf("missing kv config"))
	}
	if secret.Mount == "" {
		return fp.Left[*ResourceResult](fmt.Errorf("missing kv mount"))
	}
	if secret.Path == "" {
		return fp.Left[*ResourceResult](fmt.Errorf("missing kv path"))
	}
	client, err := self.vault.Client(ctx, r.VaultConn, r.VaultAuth).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not get a client for the secret: %w", err))
	}
	data, err := client.Client.KVv2(secret.Mount).Get(ctx, secret.Path)
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not fetch secret: %w", err))
	}
	res := &ResourceResult{Name: r.Name, Data: data.Data}
	return fp.Right(res)
}
