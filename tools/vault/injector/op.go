package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type vaultOp struct {
	vault *al.VaultStore
}

var _ ResourceFetcher = (*vaultOp)(nil)

func (self *vaultOp) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) fp.Either[*ResourceResult] {
	op := r.GetOp()
	if op == nil {
		return fp.Left[*ResourceResult](fmt.Errorf("missing op config"))
	}
	client, err := self.vault.Client(ctx, r.VaultConn, r.VaultAuth).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not create vault client for op %s: %w", r.Name, err))
	}
	if op.Method == "" || op.Method == "write" {
		data := map[string]any{}
		for key, value := range op.Data {
			data[key] = value
		}
		resp, err := client.Client.Logical().Write(op.Path, data)
		if err != nil {
			return fp.Left[*ResourceResult](fmt.Errorf("write error: %w", err))
		}
		res := &ResourceResult{Name: r.Name, Data: resp.Data}
		return fp.Right(res)
	} else {
		return fp.Left[*ResourceResult](fmt.Errorf("invalid method %s", op.Method))
	}
}
