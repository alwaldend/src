package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

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
	self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc0(func() error {
		if err := syscall.Kill(cmd.Process.Pid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("could not kill process %d: %w", cmd.Process.Pid, err)
		}
		return nil
	}))
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
