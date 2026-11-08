package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/projects/al/plugins/vault/vault_proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// resourceFetcher is internal to the Vault provider, never part of core dispatch.
type resourceFetcher interface {
	Get(context.Context, *vault_proto.Resource, []*injector.ResourceResult) (*injector.ResourceResult, error)
}

// Plugin owns lazy Vault clients, capability handlers and provider cleanup.
type Plugin struct {
	// UnimplementedResourceServiceServer supplies forward-compatible RPC defaults.
	al_proto.UnimplementedResourceServiceServer
	// ctx borrows command cancellation and diagnostic streams.
	ctx *al.CmdCtx
	// cleaner owns temporary paths until provider or engine shutdown.
	cleaner *injector.Cleaner
	// lc owns dependent process or credential lifetimes.
	lc lifecycle.Manager
	// fetchers maps resource types to their registered implementation.
	fetchers map[string]resourceFetcher
	// mx serializes startup, fetches or cleanup state.
	mx sync.Mutex
}

// NewPlugin creates an unstarted provider without authenticating.
func NewPlugin(ctx *al.CmdCtx, cleaner *injector.Cleaner) *Plugin {
	return &Plugin{ctx: ctx, cleaner: cleaner}
}

// Start initializes the provider lifecycle without authenticating.
func (p *Plugin) Start(ctx context.Context) error { return p.lc.Start(ctx) }

// Stop revokes invocation-owned credentials after requests drain.
func (p *Plugin) Stop(ctx context.Context) error { return p.lc.Stop(ctx) }

// PluginStart initializes handlers without fetching data or producing output.
func (p *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.fetchers != nil {
		return nil, errors.New("Vault provider already started")
	}
	vault := al.NewVault(req.Config)
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(vault.Stop)); err != nil {
		return nil, errors.Join(err, vault.Stop(ctx))
	}
	op := newOpFetcher(vault)
	p.fetchers = map[string]resourceFetcher{
		"kv": newSecretFetcher(vault), "op": op, "oidc": newOidcFetcher(vault),
		"vault_env": newVaultEnvFetcher(req.Config, vault),
		"vault_ssh": newVaultSshFetcher(p.ctx, req.Config, vault, op, p.cleaner),
	}
	return &al_proto.PluginStartResponse{ResourceTypes: slices.Sorted(maps.Keys(p.fetchers))}, nil
}

// ResourceFetch evaluates only an advertised Vault type and returns generic results.
func (p *Plugin) ResourceFetch(ctx context.Context, req *al_proto.ResourceFetchRequest) (*al_proto.ResourceResult, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	fetcher, ok := p.fetchers[req.Type]
	if !ok {
		return nil, fmt.Errorf("Vault provider does not support resource type %q", req.Type)
	}
	value, err := al.FromPbJson(req.Args).Get()
	if err != nil {
		return nil, errors.New("invalid Vault resource arguments")
	}
	args, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("Vault resource arguments must be an object")
	}
	settings := map[string]any{}
	fields := map[string]any{"name": req.Name, req.Type: settings}
	for k, v := range args {
		if k == "vault_auth" || k == "vault_conn" {
			fields[k] = v
		} else {
			settings[k] = v
		}
	}
	if req.Type == "op" {
		if _, exists := settings["data"]; exists {
			return nil, errors.New("Vault operation payload must use resource data")
		}
		data, err := al.FromPbJson(req.Data).Get()
		if err != nil {
			return nil, errors.New("invalid Vault operation data")
		}
		settings["data"] = data
	}
	raw, err := json.Marshal(fields)
	if err != nil {
		return nil, errors.New("invalid Vault resource settings")
	}
	resource := &vault_proto.Resource{}
	if err := protojson.Unmarshal(raw, resource); err != nil {
		return nil, errors.New("invalid Vault resource settings")
	}
	deps := []*injector.ResourceResult{}
	for _, dep := range req.Dependencies {
		result, err := injector.FromProto(dep)
		if err != nil {
			return nil, err
		}
		deps = append(deps, result)
	}
	result, err := fetcher.Get(ctx, resource, deps)
	if err != nil {
		return nil, err
	}
	if result.Name == "" {
		result.Name = req.Name
	}
	wire, err := injector.ToProto(result)
	if err != nil {
		return nil, errors.New("invalid Vault resource result")
	}
	return wire, nil
}
