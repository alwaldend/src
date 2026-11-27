package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
)

// plugin owns the backends created by AL plugin calls.
type plugin struct {
	// ctx provides the plugin logger and invocation context.
	ctx *al.CmdCtx
	// lc tracks each running backend for ordered shutdown.
	lc lifecycle.Manager
}

// newPlugin binds an empty lifecycle manager to this AL invocation.
func newPlugin(ctx *al.CmdCtx) *plugin {
	return &plugin{
		ctx: ctx,
	}
}

// Start starts the lifecycle manager before the RPC server accepts calls.
func (p *plugin) Start(ctx context.Context) error {
	if err := p.lc.Start(ctx); err != nil {
		return fmt.Errorf("start managed backends: %w", err)
	}
	return nil
}

// Stop stops every backend registered during the invocation.
func (p *plugin) Stop(ctx context.Context) error {
	if err := p.lc.Stop(ctx); err != nil {
		return fmt.Errorf("stop managed backends: %w", err)
	}
	return nil
}

// handleCall starts one backend and transfers cleanup ownership to the manager.
func (p *plugin) handleCall(ctx context.Context, req *al_proto.PluginStartRequest, call *al_proto.PluginCall) (map[string]string, error) {
	config := &tf_backend_proto.Config{}
	if _, err := al.FromPbJsonToPb(call.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin call data: %w", err)
	}
	backend, err := newBackend(ctx, p.ctx, req.Config, config)
	if err != nil {
		return nil, fmt.Errorf("could not create a terraform backend: %w", err)
	}
	if err := backend.start(ctx); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return nil, errors.Join(fmt.Errorf("could not start terraform backend: %w", err), backend.stop(cleanupCtx))
	}
	if err := p.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(backend.stop)); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return nil, errors.Join(fmt.Errorf("register backend for shutdown: %w", err), backend.stop(cleanupCtx))
	}
	env, err := backend.environment()
	if err != nil {
		return nil, fmt.Errorf("could not create env variables: %w", err)
	}
	return env, nil
}

// PluginStart returns the combined Terraform environment for the selected calls.
func (p *plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	res := &al_proto.ResourceResult{Env: map[string]string{}}
	for _, name := range slices.Sorted(maps.Keys(req.Plugin.Calls)) {
		call := req.Plugin.Calls[name]
		p.ctx.Logger.Printf("handling call %s", name)
		env, err := p.handleCall(ctx, req, call)
		if err != nil {
			return nil, fmt.Errorf("could not process call %s: %w", name, err)
		}
		for key, value := range env {
			if previous, exists := res.Env[key]; exists && previous != value {
				return nil, fmt.Errorf("conflicting backend environment export %q from selected calls", key)
			}
			res.Env[key] = value
		}
	}
	return &al_proto.PluginStartResponse{Result: res}, nil
}
