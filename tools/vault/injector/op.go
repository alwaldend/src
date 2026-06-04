package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
	"log"
)

type vaultOp struct {
	vault  *al.VaultStore
	plugin *injector_proto.Config
	logger log.Logger
}

var _ Resource = (*vaultOp)(nil)

func (self *vaultOp) Deps(ctx context.Context, name string) fp.Either[[]string] {
	return fp.Right[[]string](nil)
}

func (self *vaultOp) Get(ctx context.Context, name string, _ []*ResourceResult) fp.Either[*ResourceResult] {
	self.logger.Printf("executing a vault operation: %s", name)
	op, err := self.find(name).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not find by name: %w", err))
	}
	client, err := self.vault.Client(ctx, op.Vault, op.Auth).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not create vault client for op %s: %w", name, err))
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
		res := &ResourceResult{Name: name, Data: resp.Data}
		return fp.Right(res)
	} else {
		return fp.Left[*ResourceResult](fmt.Errorf("invalid method %s", op.Method))
	}
}

func (self *vaultOp) find(name string) fp.Either[*injector_proto.Op] {
	for i := range self.plugin.Ops {
		curOp := self.plugin.Ops[len(self.plugin.Ops)-1-i]
		if curOp.Name == name {
			return fp.Right(curOp)
		}
	}
	return fp.Left[*injector_proto.Op](fmt.Errorf("missing vault op with name %s", name))
}
