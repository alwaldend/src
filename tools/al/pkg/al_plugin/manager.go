package al_plugin

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
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Manager struct {
	cmds    []*exec.Cmd
	logger  *log.Logger
	config  *al_proto.Config
	plugins map[string]*al_proto.PluginConfig
	clients map[string]al_proto.PluginServiceClient
	run     *runfiles.Runfiles
}

func NewManager(config *al_proto.Config, stderr io.Writer) (*Manager, error) {
	if stderr == nil {
		stderr = os.Stderr
	}
	run, err := runfiles.New()
	if err != nil {
		return nil, fmt.Errorf("could not create runfiles: %w", err)
	}
	plugins := make(map[string]*al_proto.PluginConfig, len(config.Plugins))
	for _, plugin := range config.Plugins {
		_, ok := plugins[plugin.Name]
		if ok {
			return nil, fmt.Errorf("duplicate plugin name: %s", plugin.Name)
		}
		plugins[plugin.Name] = plugin
	}
	return &Manager{
		config:  config,
		cmds:    []*exec.Cmd{},
		plugins: plugins,
		run:     run,
		clients: make(map[string]al_proto.PluginServiceClient, len(plugins)),
		logger:  log.New(stderr, "com.alwaldend.src.tools.al.pkg.al_plugin.Manager ", log.Flags()),
	}, nil
}

func (self *Manager) StartPlugins(ctx context.Context, args []string, env []string) ([]*al_proto.PluginStartResponse, error) {
	plugins := []*al_proto.PluginConfig{}
	res := []*al_proto.PluginStartResponse{}
	for i, arg := range args {
		if arg == "" {
			return nil, fmt.Errorf("empty argument %d", i)
		}
		parts := strings.SplitN(arg, ".", 2)
		plugin_config, ok := self.plugins[parts[0]]
		if !ok {
			return nil, fmt.Errorf("missing config for plugin %s", parts[0])
		}
		plugin := &al_proto.PluginConfig{}
		proto.Merge(plugin, plugin_config)
		plugins = append(plugins, plugin_config)
		if len(parts) < 2 {
			continue
		}
		valueParts := strings.SplitN(parts[1], "=", 2)
		if valueParts[0] == "bin" && len(valueParts) > 1 {
			plugin_config.Bin = valueParts[1]
		} else if valueParts[0] == "arg" && len(valueParts) > 1 {
			plugin_config.Args = append(plugin_config.Args, valueParts[1])
		}
	}
	for i, plugin := range plugins {
		resp, err := self.StartPlugin(ctx, plugin)
		if err != nil {
			return nil, fmt.Errorf("could not start plugin %d (%s): %w", i, plugin.Name, err)
		}
		res = append(res, resp)
	}
	return res, nil
}

func (self *Manager) StartPlugin(ctx context.Context, plugin *al_proto.PluginConfig) (*al_proto.PluginStartResponse, error) {
	pluginJson, err := protojson.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("could not marshal plugin: %w", err)
	}
	self.logger.Printf("Starting plugin: %s", string(pluginJson))
	bin, err := self.run.Rlocation(plugin.Bin)
	if err != nil {
		return nil, fmt.Errorf("could not get rlocation of %s: %w", plugin.Bin, err)
	}
	cmd := exec.Command(bin, plugin.Args...)
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, self.run.Env()...)
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
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start plugin: %w", err)
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			self.logger.Printf("plugin %s finished with an error: %s", plugin.Name, err.Error())
		} else {
			self.logger.Printf("plugin %s finished", plugin.Name)
		}
	}()
	conn, err := grpc.NewClient(
		fmt.Sprintf("passthrough://%s", plugin.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			conn, err := NewIOConn(stdout, stdin, NewIOAddr(stdout, stdin))
			if err != nil {
				return nil, fmt.Errorf("could not create the io connection: %w", err)
			}
			return conn, nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc client: %w", err)
	}
	client := al_proto.NewPluginServiceClient(conn)
	request := &al_proto.PluginStartRequest{
		Config: self.config,
		Plugin: plugin,
	}
	response, err := client.PluginStart(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not execute plugin start request: %w", err)
	}
	return response, nil
}
