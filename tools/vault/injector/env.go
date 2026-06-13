package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type EnvFetcher struct {
	templater *Templater
}

func NewEnvFetcher(templater *Templater) *EnvFetcher {
	return &EnvFetcher{templater}
}

var _ ResourceFetcher = (*EnvFetcher)(nil)

func (self *EnvFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.EnvFetcher"
}

func (self *EnvFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error) {
	if r.GetEnv() == nil {
		return nil, fmt.Errorf("missing env config")
	}
	content, err := self.templater.Template(ctx, r.GetEnv().Value, d).Get()
	if err != nil {
		return nil, fmt.Errorf("could not template: %w", err)
	}
	res := &ResourceResult{
		Name: r.Name,
		Env:  map[string]string{r.Name: content},
	}
	return res, nil
}
