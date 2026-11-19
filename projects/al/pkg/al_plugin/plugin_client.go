package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net"
	"os"
	"os/exec"
	"slices"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PluginClient struct {
	stopOnce sync.Once
	stopErr  error
	wait     chan error
	cmd      *exec.Cmd
	bin      string
	client   al_proto.PluginServiceClient
	conn     *grpc.ClientConn
	ioConn   *IOConn
	req      *al_proto.PluginStartRequest
	resp     *al_proto.PluginStartResponse
	ctx      *al.CmdCtx
}

func NewPluginClient(ctx *al.CmdCtx, run *runfiles.Runfiles, config *al_proto.Config, plugin *al_proto.PluginConfig) (*PluginClient, error) {
	bin, err := run.Rlocation(plugin.Bin)
	if err != nil {
		return nil, fmt.Errorf("could not find location of %s: %w", plugin.Bin, err)
	}
	cmd := exec.Command(bin, plugin.Args...)
	constructed := false
	defer func() {
		if !constructed {
			_ = closeUnstartedChildPipes(cmd)
		}
	}()
	cmd.WaitDelay = time.Second * 10
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, run.Env()...)
	for key, value := range plugin.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not create stdout pipe: %w", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		_ = stdout.Close()
		return nil, fmt.Errorf("could not create stdin pipe: %w", err)
	}
	ioConn, err := NewIOConn(stdout, stdin, NewIOAddr(stdout, stdin))
	if err != nil {
		_ = stdout.Close()
		_ = stdin.Close()
		return nil, fmt.Errorf("could not create an io connection: %w", err)
	}
	conn, err := grpc.NewClient(
		fmt.Sprintf("passthrough://%s", plugin.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return ioConn, nil
		}),
	)
	if err != nil {
		_ = ioConn.Close()
		return nil, fmt.Errorf("could not create grpc client: %w", err)
	}
	client := al_proto.NewPluginServiceClient(conn)
	request := &al_proto.PluginStartRequest{
		Config: config,
		Plugin: plugin,
	}
	constructed = true
	return &PluginClient{
		req:    request,
		ioConn: ioConn,
		cmd:    cmd,
		client: client,
		conn:   conn,
		ctx:    ctx,
		wait:   make(chan error, 1),
	}, nil
}

func (self *PluginClient) Stop(ctx context.Context) error {
	self.stopOnce.Do(func() { self.stopErr = self.stop(ctx) })
	return self.stopErr
}

func (self *PluginClient) stop(ctx context.Context) error {
	var errs []error
	errs = append(errs, lifecycle.StopProcess(ctx, self.cmd, self.wait))
	errs = append(errs, closeUnstartedChildPipes(self.cmd))
	if self.conn != nil {
		errs = append(errs, self.conn.Close())
	}
	if self.ioConn != nil {
		errs = append(errs, self.ioConn.Close())
	}
	return errors.Join(errs...)
}

func (self *PluginClient) StartResponse() (*al_proto.PluginStartResponse, bool) {
	if self.resp == nil {
		return nil, false
	}
	return self.resp, true
}

func (self *PluginClient) Start(ctx context.Context) error {
	self.ctx.Logger.Printf("starting plugin %s", self.req.Plugin.Name)
	if err := self.cmd.Start(); err != nil {
		return fmt.Errorf("could not start the plugin binary: %w", err)
	}
	go func() { self.wait <- self.cmd.Wait() }()
	resp, err := self.client.PluginStart(ctx, self.req)
	if err != nil {
		return fmt.Errorf("could not execute plugin start request: %w", err)
	}
	self.ctx.Logger.Printf("finished starting plugin %s, env: %s", self.req.Plugin.Name, slices.Collect(maps.Keys(resp.Env)))
	self.resp = resp
	return nil
}

// exec.Start normally closes the child-facing pipe ends. If initialization
// stops before Start, they remain owned by cmd and must be closed explicitly.
// Stderr is inherited, not an owned pipe, and must remain open.
func closeUnstartedChildPipes(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process != nil {
		return nil
	}
	var errs []error
	for _, endpoint := range []any{cmd.Stdin, cmd.Stdout} {
		if closer, ok := endpoint.(io.Closer); ok {
			if err := closer.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
