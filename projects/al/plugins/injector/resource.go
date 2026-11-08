package main

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/projects/al/plugins/injector/injector_proto"
)

type ResourceType string

type ResourceResult struct {
	Name  string
	Env   map[string]string
	Files []string
	Data  any
}

type ResourceFetcher interface {
	Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error)
	String() string
}

type ResourceManager struct {
	env      *EnvFetcher
	file     *FileFetcher
	op       *OpFetcher
	secret   *SecretFetcher
	vaultEnv *VaultEnvFetcher
	process  *ProcessFetcher
	vaultSsh *VaultSshFetcher
	oidc     *OidcFetcher
	ctx      *al.CmdCtx
}

func NewResourceManager(
	ctx *al.CmdCtx,
	env *EnvFetcher,
	file *FileFetcher,
	op *OpFetcher,
	secret *SecretFetcher,
	vaultEnv *VaultEnvFetcher,
	vaultSsh *VaultSshFetcher,
	process *ProcessFetcher,
	oidc *OidcFetcher,
) *ResourceManager {
	return &ResourceManager{
		ctx:      ctx,
		env:      env,
		file:     file,
		op:       op,
		secret:   secret,
		vaultEnv: vaultEnv,
		vaultSsh: vaultSsh,
		process:  process,
		oidc:     oidc,
	}
}

// get evaluates a resource in dependency order and publishes its exports once.
func (self *ResourceManager) get(ctx context.Context, config *injector_proto.Config, resConfig *injector_proto.Resource, cache map[string]*ResourceResult, active map[string]bool, environment map[string]string) fp.Either[*ResourceResult] {
	if res, ok := cache[resConfig.Name]; ok {
		return fp.Right(res)
	}
	if active[resConfig.Name] {
		return fp.Left[*ResourceResult](fmt.Errorf("dependency cycle at resource %q", resConfig.Name))
	}
	active[resConfig.Name] = true
	defer delete(active, resConfig.Name)
	deps := []*ResourceResult{}
	for _, name := range resConfig.Deps {
		curR, ok := config.Res[name]
		if !ok {
			return fp.Left[*ResourceResult](fmt.Errorf("resource %q references missing dependency %q", resConfig.Name, name))
		}
		dep, err := self.get(ctx, config, curR, cache, active, environment).Get()
		if err != nil {
			return fp.Left[*ResourceResult](fmt.Errorf("fetch dependency %q: %w", name, err))
		}
		deps = append(deps, dep)
	}
	var fetcher ResourceFetcher
	switch resConfig.Res.(type) {
	case *injector_proto.Resource_Env:
		fetcher = self.env
	case *injector_proto.Resource_File:
		fetcher = self.file
	case *injector_proto.Resource_Op:
		fetcher = self.op
	case *injector_proto.Resource_Kv:
		fetcher = self.secret
	case *injector_proto.Resource_VaultEnv:
		fetcher = self.vaultEnv
	case *injector_proto.Resource_VaultSsh:
		fetcher = self.vaultSsh
	case *injector_proto.Resource_Process:
		fetcher = self.process
	case *injector_proto.Resource_Oidc:
		fetcher = self.oidc
	default:
		return fp.Left[*ResourceResult](fmt.Errorf("resource %s is missing config", resConfig.Name))
	}
	self.ctx.Logger.Printf("fetching resource %s using %s", resConfig.Name, fetcher)
	res, err := fetcher.Get(ctx, resConfig, deps)
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not execute %s: %w", fetcher, err))
	}
	cache[resConfig.Name] = res
	maps.Copy(environment, res.Env)
	return fp.Right(res)
}

// Get merges dependency exports before dependent exports in a named graph.
func (self *ResourceManager) Get(ctx context.Context, config *injector_proto.Config) (*ResourceResult, error) {
	resRes := &ResourceResult{Env: map[string]string{}}
	cache := map[string]*ResourceResult{}
	active := map[string]bool{}
	for _, name := range slices.Sorted(maps.Keys(config.Res)) {
		resConfig := config.Res[name]
		_, err := self.get(ctx, config, resConfig, cache, active, resRes.Env).Get()
		if err != nil {
			return nil, fmt.Errorf("could not process resource %s: %w", resConfig.Name, err)
		}
	}
	return resRes, nil
}
