package main

import (
	"context"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type ResourceType string

type ResourceResult struct {
	Name  string
	Env   map[string]string
	Files []string
	Data  any
}

type ResourceFetcher interface {
	Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) fp.Either[*ResourceResult]
}

type ResourceManager struct {
	env    *env
	file   *file
	op     *vaultOp
	secret *secret
}

func (self *ResourceManager) Get(ctx context.Context, name string) fp.Either[*ResourceResult] {
}
