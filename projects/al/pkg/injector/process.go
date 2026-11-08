package injector

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

// processFetcher keeps resource processes alive until engine shutdown.
type processFetcher struct {
	// lc owns dependent process or credential lifetimes.
	lc *lifecycle.Manager
}

// newProcessFetcher binds process lifetime ownership.
func newProcessFetcher(lc *lifecycle.Manager) *processFetcher {
	return &processFetcher{lc: lc}
}

var _ ResourceFetcher = (*processFetcher)(nil)

// Get evaluates one resource using the supplied dependency results.
func (self *processFetcher) Get(ctx context.Context, r *Resource, d []*ResourceResult) (*ResourceResult, error) {
	process := &al_proto.InjectorProcess{}
	if err := decodeArgs(r.Args, process); err != nil {
		return nil, fmt.Errorf("decode process settings: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("start resource process: %w", err)
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
	stop := func(ctx context.Context) error { return lifecycle.StopResourceProcess(ctx, cmd, wait) }
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
