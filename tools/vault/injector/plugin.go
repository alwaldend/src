package main

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

type Plugin struct {
	ctx     *al.CmdCtx
	cleaner *Cleaner
	lc      lifecycle.Manager
}

func NewPlugin(ctx *al.CmdCtx, cleaner *Cleaner) *Plugin {
	return &Plugin{
		ctx:     ctx,
		cleaner: cleaner,
	}
}

func (self *Plugin) Start(ctx context.Context) error {
	return self.lc.Start(ctx)
}

func (self *Plugin) Stop(ctx context.Context) error {
	return self.lc.Stop(ctx)
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	vault := al.NewVault(req.Config)
	templater := &Templater{}
	opFetcher := NewOpFetcher(vault)
	envFetcher := NewEnvFetcher(templater)
	fileFetcher := NewFileFetcher(vault, templater, self.cleaner)
	secretFetcher := NewSecretFetcher(vault)
	sshFetcher := NewVaultSshFetcher(self.ctx, req.Config, vault, opFetcher, self.cleaner)
	vaultEnvFetcher := NewVaultEnvFetcher(templater, req.Config, vault)
	oidcFetcher := NewOidcFetcher(vault)
	processFetcher := NewProcessFetcher(self.ctx, req.Config, &self.lc)
	manager := NewResourceManager(self.ctx, envFetcher, fileFetcher, opFetcher, secretFetcher, vaultEnvFetcher, sshFetcher, processFetcher, oidcFetcher)
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
	resource, err := manager.Get(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not process request: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: resource.Env,
	}
	return res, nil
}
