package main

import (
	"context"
	"fmt"
	"log"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

type Plugin struct {
	ctx     *al.CmdCtx
	logger  *log.Logger
	cleaner *Cleaner
}

func NewPlugin(ctx *al.CmdCtx, cleaner *Cleaner) *Plugin {
	return &Plugin{
		ctx:     ctx,
		logger:  log.New(ctx.Stderr, "com.alwaldend.src.tools.vault.injector.Plugin ", ctx.LogFlags),
		cleaner: cleaner,
	}
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	vault := al.NewVault(req.Config)
	templater := &Templater{}
	manager := &ResourceManager{
		env:      &EnvFetcher{templater: templater},
		file:     &FileFetcher{vault: vault, templater: templater, cleaner: self.cleaner},
		op:       &OpFetcher{vault: vault},
		secret:   &SecretFetcher{vault: vault},
		vaultEnv: &VaultEnvFetcher{templater: templater, config: req.Config, vault: vault},
	}
	config := &injector_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.Plugin.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin data: %w", err)
	}
	for _, call := range req.Plugin.Calls {
		curConfig := &injector_proto.Config{}
		if _, err := al.FromPbJsonToPb(call.Data, curConfig).Get(); err != nil {
			return nil, fmt.Errorf("could not parse plugin call data of %s: %w", call.Name, err)
		}
		config.Res = append(config.Res, curConfig.Res...)
	}
	resource, err := manager.Get(ctx, config).Get()
	if err != nil {
		return nil, fmt.Errorf("could not process request: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: resource.Env,
	}
	return res, nil
}
