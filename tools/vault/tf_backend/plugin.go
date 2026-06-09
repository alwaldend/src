package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/tf_backend/tf_backend_proto"
)

type Plugin struct {
	ctx      *al.CmdCtx
	logger   *log.Logger
	backends []*TfBackend
}

func NewPlugin(ctx *al.CmdCtx) *Plugin {
	return &Plugin{
		ctx:    ctx,
		logger: log.New(ctx.Stderr, "com.alwaldend.src.tools.vault.tf_backend.Plugin ", ctx.LogFlags),
	}
}

func (self *Plugin) Shutdown(ctx context.Context) error {
	self.logger.Printf("Shutting down the plugin")
	errsCh := make(chan error, len(self.backends))
	errs := []error{}
	for _, backend := range self.backends {
		go func() {
			err := backend.Shutdown(ctx)
			if err != nil {
				err = fmt.Errorf("could not shut down a backend: %w", err)
			}
			errsCh <- err
		}()
	}
	for range len(errsCh) {
		select {
		case err := <-errsCh:
			errs = append(errs, err)
		case <-ctx.Done():
			return fmt.Errorf("shutdown timed out")
		}
	}
	return errors.Join(errs...)
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
	self.backends = append(self.backends, backend)
	env, err := backend.Env()
	if err != nil {
		return nil, fmt.Errorf("could not create env variables: %w", err)
	}
	return env, nil
}

func (self *Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	res := &al_proto.PluginStartResponse{Env: map[string]string{}}
	for _, call := range req.Plugin.Calls {
		self.logger.Printf("Handling call %s", call.Name)
		env, err := self.handleCall(req, call)
		if err != nil {
			return nil, fmt.Errorf("could not process call %s: %w", call.Name, err)
		}
		maps.Copy(res.Env, env)
	}
	self.logger.Printf("Finished request")
	return res, nil
}
