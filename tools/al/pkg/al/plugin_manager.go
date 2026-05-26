package al

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PluginManager struct {
	cmds   []*exec.Cmd
	logger *log.Logger
}

func NewPluginManager(stderr io.Writer) *PluginManager {
	if stderr == nil {
		stderr = os.Stderr
	}
	return &PluginManager{
		cmds:   []*exec.Cmd{},
		logger: log.New(stderr, "al.pkg.al.PluginManager", log.Flags()),
	}
}

type PluginManagerStart struct {
	name string
	bin  string
	args []string
}

func (self *PluginManager) StartPlugins(ctx context.Context, args []string) ([]*al_proto.PluginStartResponse, error) {
	plugins := make(map[string]*PluginManagerStart)
	res := []*al_proto.PluginStartResponse{}
	for i, arg := range args {
		if arg == "" {
			return nil, fmt.Errorf("empty argument %d", i)
		}
		parts := strings.SplitN(arg, ".", 2)
		plugin, ok := plugins[parts[0]]
		if !ok {
			plugin = &PluginManagerStart{}
			plugins[parts[0]] = plugin
		}
		if len(parts) < 2 {
			continue
		}
		valueParts := strings.SplitN(parts[1], "=", 2)
		if valueParts[0] == "bin" && len(valueParts) > 1 {
			plugin.bin = valueParts[1]
		} else if valueParts[0] == "arg" && len(valueParts) > 1 {
			plugin.args = append(plugin.args, valueParts[0])
		}
	}
	for name, plugin := range plugins {
		plugin.name = name
		resp, err := self.StartPlugin(ctx, plugin)
		if err != nil {
			return nil, fmt.Errorf("could not start plugin %s: %w", name, err)
		}
		res = append(res, resp)
	}
	return res, nil
}

func (self *PluginManager) StartPlugin(ctx context.Context, plugin *PluginManagerStart) (*al_proto.PluginStartResponse, error) {
	self.logger.Output(0, fmt.Sprintf("Starting plugin: %s", plugin))
	cmd := exec.Command(plugin.bin, plugin.args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not create stdout pipe: %w", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("could not create stdin pipe: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start plugin: %w", err)
	}
	conn, err := grpc.NewClient(
		"io",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			conn, err := NewIOConn(stdout, stdin, NewIOAddr(stdout, stdin))
			return conn, err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc client: %w", err)
	}
	client := al_proto.NewPluginServiceClient(conn)
	request := &al_proto.PluginStartRequest{}
	response, err := client.PluginStart(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not execute plugin start request: %w", err)
	}
	return response, nil
}
