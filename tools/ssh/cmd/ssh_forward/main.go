package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

// resourceType is the AL capability this provider serves.
const resourceType = "ssh_forward"

// provider owns tunnels created for its caller's invocation.
type provider struct {
	// UnimplementedResourceServiceServer provides defaults for future RPCs.
	al_proto.UnimplementedResourceServiceServer
	// tunnels closes the listeners and SSH sessions during provider shutdown.
	tunnels lifecycle.Manager
}

// Start initializes tunnel lifetime management.
func (p *provider) Start(ctx context.Context) error {
	if err := p.tunnels.Start(ctx); err != nil {
		return fmt.Errorf("start tunnel manager: %w", err)
	}
	return nil
}

// Stop releases every tunnel registered by a resource request.
func (p *provider) Stop(ctx context.Context) error {
	if err := p.tunnels.Stop(ctx); err != nil {
		return fmt.Errorf("stop tunnels: %w", err)
	}
	return nil
}

// PluginStart advertises the forwarding resource without opening connections.
func (p *provider) PluginStart(context.Context, *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	return &al_proto.PluginStartResponse{ResourceTypes: []string{resourceType}}, nil
}

// ResourceFetch opens a tunnel using the certificate supplied by its dependency.
func (p *provider) ResourceFetch(ctx context.Context, req *al_proto.ResourceFetchRequest) (*al_proto.ResourceResult, error) {
	if req.GetType() != resourceType {
		return nil, fmt.Errorf("unsupported SSH provider resource type %q", req.GetType())
	}
	if len(req.GetDependencies()) != 1 {
		return nil, errors.New("SSH forward requires one client credential dependency")
	}
	args := req.GetArgs().GetValMap().GetVal()
	credentials := req.GetDependencies()[0].GetData().GetValMap().GetVal()
	cfg := forwardConfig{
		host:     args["host"].GetValString(),
		user:     args["user"].GetValString(),
		target:   args["target"].GetValString(),
		ca:       args["host_ca"].GetValString(),
		keyPath:  credentials["private_key"].GetValString(),
		certPath: credentials["public_key_signed"].GetValString(),
	}
	tunnel, err := openForward(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open SSH resource %q: %w", req.GetName(), err)
	}
	address, err := al.ToPbJson(map[string]any{"address": tunnel.listener.Addr().String()}).Get()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("encode forward address: %w", err), tunnel.stop(ctx))
	}
	if err := p.tunnels.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(tunnel.stop)); err != nil {
		return nil, errors.Join(fmt.Errorf("register forward cleanup: %w", err), tunnel.stop(ctx))
	}
	return &al_proto.ResourceResult{Name: req.GetName(), Data: address}, nil
}

// main runs the provider and ties its lifetime to process cancellation.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	command := al.NewCmdCtx(ctx, "ssh_forward ")
	plugin := &provider{}
	var lc lifecycle.Manager
	err := lc.Add(plugin, al_plugin.NewPluginServer(command, plugin))
	if err == nil {
		err = lc.Run(ctx, 10*time.Second)
	}
	if err != nil {
		command.Logger.Printf("SSH forward provider failed: %s", err)
		os.Exit(1)
	}
}
