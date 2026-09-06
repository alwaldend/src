package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type ProcessFetcher struct {
	ctx    *al.CmdCtx
	config *al_proto.Config
	lc     *lifecycle.Manager
}

func NewProcessFetcher(ctx *al.CmdCtx, config *al_proto.Config, lc *lifecycle.Manager) *ProcessFetcher {
	return &ProcessFetcher{ctx: ctx, config: config, lc: lc}
}

var _ ResourceFetcher = (*ProcessFetcher)(nil)

func (self *ProcessFetcher) String() string {
	return "com.alwaldend.src.tools.vault.injector.ProcessFetcher"
}

func (self *ProcessFetcher) Get(ctx context.Context, r *injector_proto.Resource, d []*ResourceResult) (*ResourceResult, error) {
	process := r.GetProcess()
	if process == nil {
		return nil, fmt.Errorf("missing process config")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	cmd := exec.Command(process.Name, process.Args...)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	for _, res := range d {
		for key, value := range res.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start process: %w", err)
	}
	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()
	stop := func(ctx context.Context) error { return lifecycle.StopProcess(ctx, cmd, wait) }
	if err := self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(stop)); err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return nil, errors.Join(err, stop(cleanupCtx))
	}

	res := &ResourceResult{
		Name: r.Name,
		Data: map[string]any{
			"pid":  cmd.Process.Pid,
			"name": process.Name,
			"args": process.Args,
		},
	}
	return res, nil
}
