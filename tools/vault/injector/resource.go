package main

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
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
	}
}

func (self *ResourceManager) get(ctx context.Context, config *injector_proto.Config, resConfig *injector_proto.Resource, cache map[string]*ResourceResult) fp.Either[*ResourceResult] {
	if res, ok := cache[resConfig.Name]; ok {
		return fp.Right(res)
	}
	deps := []*ResourceResult{}
	for _, curR := range config.Res {
		if curR.Name == resConfig.Name || !slices.Contains(resConfig.Deps, curR.Name) {
			continue
		}
		dep, err := self.get(ctx, config, curR, cache).Get()
		if err != nil {
			return fp.Left[*ResourceResult](fmt.Errorf("could not fetch dependency %s: %w", curR.Name, err))
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
	default:
		return fp.Left[*ResourceResult](fmt.Errorf("resource %s is missing config", resConfig.Name))
	}
	self.ctx.Logger.Printf("fetching resource %s using %s", resConfig.Name, fetcher)
	res, err := fetcher.Get(ctx, resConfig, deps)
	if err != nil {
		return fp.Left[*ResourceResult](fmt.Errorf("could not execute %s: %w", fetcher, err))
	}
	cache[resConfig.Name] = res
	return fp.Right(res)
}

func (self *ResourceManager) Get(ctx context.Context, config *injector_proto.Config) (*ResourceResult, error) {
	resRes := &ResourceResult{Env: map[string]string{}}
	cache := map[string]*ResourceResult{}
	for _, resConfig := range config.Res {
		res, err := self.get(ctx, config, resConfig, cache).Get()
		if err != nil {
			return nil, fmt.Errorf("could not process resource %s: %w", resConfig.Name, err)
		}
		maps.Copy(resRes.Env, res.Env)
	}
	return resRes, nil
}
