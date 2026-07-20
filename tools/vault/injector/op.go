package main

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
	"github.com/hashicorp/vault/api"
)

type OpFetcher struct {
	vault *al.VaultStore
}

func NewOpFetcher(vault *al.VaultStore) *OpFetcher {
	return &OpFetcher{vault}
}

var _ ResourceFetcher = (*OpFetcher)(nil)

func (self *OpFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.OpFetcher"
}

func (self *OpFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error) {
	op := r.GetOp()
	if op == nil {
		return nil, fmt.Errorf("missing op config")
	}
	client, err := self.vault.Client(ctx, r.VaultConn, r.VaultAuth).Get()
	if err != nil {
		return nil, fmt.Errorf("could not create vault client for op %s: %w", r.Name, err)
	}
	data := map[string]any{}
	for key, value := range op.Data {
		data[key] = value
	}
	var resp *api.Secret
	logical := client.Client.Logical()
	switch op.Method {
	case "", "read":
		resp, err = logical.Read(op.Path)
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}
	case "write":
		resp, err = logical.Write(op.Path, data)
		if err != nil {
			return nil, fmt.Errorf("write error: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid method %s", op.Method)
	}
	res := &ResourceResult{Name: r.Name, Data: resp.Data}
	return res, nil
}
