package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type EnvFetcher struct {
	templater *Templater
}

var _ ResourceFetcher = (*EnvFetcher)(nil)

func (self *EnvFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.EnvFetcher"
}

func (self *EnvFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) fp.Either[*ResourceResult] {
	if r.GetEnv() == nil {
		return fp.Left[*ResourceResult](fmt.Errorf("missing env config"))
	}
	content, err := self.templater.Template(ctx, r.GetEnv().Value, d).Get()
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not template: %w", err))
	}
	res := &ResourceResult{
		Name: r.Name,
		Env:  map[string]string{r.Name: content},
	}
	return fp.Right(res)
}
