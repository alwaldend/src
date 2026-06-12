package main

import (
	"context"
	"fmt"
	"maps"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
)

type Plugin struct {
	ctx *al.CmdCtx
	lc  lifecycle.Manager
}

func NewPlugin(ctx *al.CmdCtx) *Plugin {
	return &Plugin{
		ctx: ctx,
	}
}

func (self *Plugin) Start(ctx context.Context) error {
	return self.lc.Start(ctx)
}

func (self *Plugin) Stop(ctx context.Context) error {
	return self.lc.Stop(ctx)
}

func (self *Plugin) handleCall(req *al_proto.PluginStartRequest, call *al_proto.PluginCall) (map[string]string, error) {
	config := &tf_backend_proto.Config{}
	if _, err := al.FromPbJsonToPb(call.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin call data: %w", err)
	}
	backend, err := NewTfBackend(self.ctx, req.Config, config)
	if err != nil {
		return nil, fmt.Errorf("could not create a terraform backend: %w", err)
	}
	if err := backend.Start(self.ctx.Ctx); err != nil {
		return nil, fmt.Errorf("could not start terraform backend: %w", err)
	}
	self.lc.Add(lifecycle.StoppableFunc(backend.Stop))
	env, err := backend.Env()
	if err != nil {
		return nil, fmt.Errorf("could not create env variables: %w", err)
	}
	return env, nil
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	res := &al_proto.PluginStartResponse{Env: map[string]string{}}
	for _, call := range req.Plugin.Calls {
		self.ctx.Logger.Printf("handling call %s", call.Name)
		env, err := self.handleCall(req, call)
		if err != nil {
			return nil, fmt.Errorf("could not process call %s: %w", call.Name, err)
		}
		maps.Copy(res.Env, env)
	}
	return res, nil
}
