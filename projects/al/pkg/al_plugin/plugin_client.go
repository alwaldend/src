package al_plugin

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
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
	req        *al_proto.PluginStartRequest
	ctx        context.Context
	cancel     context.CancelFunc
	res        chan error
}

func NewPluginClient(stderr io.Writer, run *runfiles.Runfiles, config *al_proto.Config, plugin *al_proto.PluginConfig) (*PluginClient, error) {
	ctx, cancel := context.WithCancel(context.Background())
	configJson, err := protojson.MarshalOptions{}.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("could not marshal plugin config: %w", err)
	}
	bin, err := run.Rlocation(plugin.Bin)
	if err != nil {
		return nil, fmt.Errorf("could not find location of %s: %w", plugin.Bin, err)
	}
	cmd := exec.CommandContext(ctx, bin, plugin.Args...)
	cmd.WaitDelay = time.Second * 10
	cmd.Cancel = func() error {
		return syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}
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
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	client := al_proto.NewPluginServiceClient(conn)
	request := &al_proto.PluginStartRequest{
		Config: config,
		Plugin: plugin,
	}
	return &PluginClient{
		req:        request,
		configJson: configJson,
		logger:     log.New(stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.PluginClient ", log.Flags()),
		cmd:        cmd,
		client:     client,
		ctx:        ctx,
		cancel:     cancel,
		res:        make(chan error, 1),
	}, nil
}

func (self *PluginClient) Shutdown(ctx context.Context) error {
	var err error
	self.cancel()
	select {
	case err = <-self.res:
		if err != nil {
			err = fmt.Errorf("plugin %s shut down with an error: %w", self.req.Plugin.Name, err)
		}
	case <-ctx.Done():
		err = fmt.Errorf("shutdown timed out")
	}
	return err
}

func (self *PluginClient) Start(ctx context.Context) (*al_proto.PluginStartResponse, error) {
	self.logger.Printf("Starting plugin %s: %s", self.req.Plugin.Name, string(self.configJson))
	if err := self.cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start the plugin binary: %w", err)
	}
	func() {
		var err error
		defer func() {
			self.cancel()
			self.res <- err
		}()
		err = self.cmd.Wait()
		if err != nil {
			self.logger.Printf("Plugin %s shut down with an error: %s", self.req.Plugin.Name, err)
		} else {
			self.logger.Printf("Plugin %s shut down", self.req.Plugin.Name)
		}
	}()
	response, err := self.client.PluginStart(ctx, self.req)
	if err != nil {
		self.cancel()
		return nil, fmt.Errorf("could not execute plugin start request: %w", err)
	}
	self.logger.Printf("Finished starting plugin %s", self.req.Plugin.Name)
	return response, nil
}
