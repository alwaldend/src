package main

import (
	"context"
	"fmt"
	"io"
	"maps"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
)

type Plugin struct {
	wg  *fp.WaitGroupE
	ctx context.Context
}

func NewPlugin(ctx context.Context) *Plugin {
	return &Plugin{ctx: ctx, wg: &fp.WaitGroupE{}}
}

func (self *Plugin) Run(stdin io.Reader, stdout io.Writer) error {
	self.wg.Go(func() error {
		if _, err := al_plugin.Serve(self.ctx, stdin, stdout, self).Get(); err != nil {
			err = fmt.Errorf("could not serve the plugin: %w", err)
		}
		logger.Printf("Finished serving the plugin")
		return nil
	})
	if err := self.wg.Wait(); err != nil {
		return fmt.Errorf("could not run the plugin: %w", err)
	}
	return nil
}

func (self *Plugin) handleCall(req *al_proto.PluginStartRequest, call *al_proto.PluginCall) (map[string]string, error) {
	config := &tf_backend_proto.Config{}
	if _, err := al.FromPbJsonToPb(call.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin call data: %w", err)
	}
	backend, err := NewTfBackend(self.ctx, req.Config, config, self.wg)
	if err != nil {
		return nil, fmt.Errorf("could not create a terraform backend: %w", err)
	}
	if err := backend.Start(self.ctx); err != nil {
		return nil, fmt.Errorf("could not start terraform backend: %w", err)
	}
	env, err := backend.Env()
	if err != nil {
		return nil, fmt.Errorf("could not create env variables: %w", err)
	}
	return env, nil
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	res := &al_proto.PluginStartResponse{Env: map[string]string{}}
	for _, call := range req.Plugin.Calls {
		logger.Printf("Handling call %s", call.Name)
		env, err := self.handleCall(req, call)
		if err != nil {
			return nil, fmt.Errorf("could not process call %s: %w", call.Name, err)
		}
		maps.Copy(res.Env, env)
	}
	logger.Printf("Finished request")
	return res, nil
}
