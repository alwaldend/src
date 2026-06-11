package al_plugin

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type PluginClient struct {
	configJson []byte
	cmd        *exec.Cmd
	bin        string
	logger     *log.Logger
	client     al_proto.PluginServiceClient
	conn       *grpc.ClientConn
	ioConn     *IOConn
	req        *al_proto.PluginStartRequest
	resp       *al_proto.PluginStartResponse
	ctx        *al.CmdCtx
	res        chan error
}

func NewPluginClient(ctx *al.CmdCtx, run *runfiles.Runfiles, config *al_proto.Config, plugin *al_proto.PluginConfig) (*PluginClient, error) {
	configJson, err := protojson.MarshalOptions{}.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("could not marshal plugin config: %w", err)
	}
	bin, err := run.Rlocation(plugin.Bin)
	if err != nil {
		return nil, fmt.Errorf("could not find location of %s: %w", plugin.Bin, err)
	}
	cmd := exec.Command(bin, plugin.Args...)
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
		return nil, fmt.Errorf("could not create stdin pipe: %w", err)
	}
	ioConn, err := NewIOConn(stdout, stdin, NewIOAddr(stdout, stdin))
	if err != nil {
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
		return nil, fmt.Errorf("could not create grpc client: %w", err)
	}
	client := al_proto.NewPluginServiceClient(conn)
	request := &al_proto.PluginStartRequest{
		Config: config,
		Plugin: plugin,
	}
	return &PluginClient{
		req:        request,
		configJson: configJson,
		logger:     log.New(ctx.Stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.PluginClient ", ctx.LogFlags),
		cmd:        cmd,
		client:     client,
		conn:       conn,
		ctx:        ctx,
		res:        make(chan error, 1),
	}, nil
}

func (self *PluginClient) Stop(ctx context.Context) error {
	var wg fp.WaitGroupE
	wg.Go(func() error {
		if err := syscall.Kill(self.cmd.Process.Pid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("could not kill process %d: %w", self.cmd.Process.Pid, err)
		}
		return nil
	})
	wg.Go(func() error {
		if err := self.conn.Close(); err != nil {
			return fmt.Errorf("could not close the grpc client connection: %w", err)
		}
		return nil
	})
	if err := wg.WaitCtx(ctx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}

func (self *PluginClient) StartResponse() (*al_proto.PluginStartResponse, bool) {
	if self.resp == nil {
		return nil, false
	}
	return self.resp, true
}

func (self *PluginClient) Start(ctx context.Context) error {
	self.logger.Printf("Starting plugin %s: %s", self.req.Plugin.Name, string(self.configJson))
	if err := self.cmd.Start(); err != nil {
		return fmt.Errorf("could not start the plugin binary: %w", err)
	}
	resp, err := self.client.PluginStart(ctx, self.req)
	if err != nil {
		return fmt.Errorf("could not execute plugin start request: %w", err)
	}
	self.logger.Printf("Finished starting plugin %s", self.req.Plugin.Name)
	self.resp = resp
	return nil
}
