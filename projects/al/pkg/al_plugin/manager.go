package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"net"
	"os"
	"os/exec"
	"runtime/debug"
	"slices"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
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
		logger:  log.New(stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.Manager ", log.Flags()),
	}, nil
}

func (self *Manager) StartPlugins(ctx context.Context, args []string) ([]*al_proto.PluginStartResponse, error) {
	plugins := map[string]*al_proto.PluginConfig{}
	for _, arg := range args {
		if arg == "" {
			continue
		}
		parts := strings.SplitN(arg, ".", 2)
		plugin, ok := plugins[parts[0]]
		if !ok {
			plugin_config, ok := self.plugins[parts[0]]
			if !ok {
				return nil, fmt.Errorf("missing config for plugin %s", parts[0])
			}
			plugin = &al_proto.PluginConfig{}
			proto.Merge(plugin, plugin_config)
			plugins[parts[0]] = plugin
		}
		if len(parts) > 1 {
			valueParts := strings.SplitN(parts[1], "=", 2)
			if valueParts[0] == "bin" && len(valueParts) > 1 {
				plugin.Bin = valueParts[1]
			} else if valueParts[0] == "arg" && len(valueParts) > 1 {
				plugin.Args = append(plugin.Args, valueParts[1])
			}
		}
	}
	for _, plugin := range plugins {
		for _, extendKey := range plugin.Extends {
			extend, ok := self.plugins[extendKey]
			if !ok {
				return nil, fmt.Errorf("extends contains missing plugin %s", extendKey)
			}
			proto.Merge(plugin, extend)
		}
	}
	res, err := self.startPlugins(ctx, slices.Collect(maps.Values(plugins)))
	if err != nil {
		return nil, fmt.Errorf("could not start plugins: %w", err)
	}
	return res, nil
}

func (self *Manager) startPlugins(ctx context.Context, plugins []*al_proto.PluginConfig) ([]*al_proto.PluginStartResponse, error) {
	var wg sync.WaitGroup
	errsWorker := make(chan error, len(plugins))
	responsesWorker := make(chan *al_proto.PluginStartResponse, len(plugins))
	errsConsumer := make(chan error, len(plugins))
	responsesConsumer := make(chan *al_proto.PluginStartResponse, len(plugins))
	for _, plugin := range plugins {
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					errsWorker <- fmt.Errorf("panicked: %s\n%s", r, debug.Stack())
				}
			}()
			resp, err := self.startPlugin(ctx, plugin)
			if err == nil {
				responsesWorker <- resp
			} else {
				errsWorker <- fmt.Errorf("could not start plugin %s: %w", plugin.Name, err)
			}
		})
	}
	go func() {
		for err := range errsWorker {
			errsConsumer <- err
		}
		close(errsConsumer)
	}()
	go func() {
		for resp := range responsesWorker {
			responsesConsumer <- resp
		}
		close(responsesConsumer)
	}()
	wg.Wait()
	close(errsWorker)
	close(responsesWorker)
	var err error
	res := []*al_proto.PluginStartResponse{}
	for curErr := range errsConsumer {
		err = errors.Join(err, curErr)
	}
	for resp := range responsesConsumer {
		res = append(res, resp)
	}
	if err != nil {
		return nil, fmt.Errorf("could not start plugins: %w", err)
	}
	return res, nil

}

func (self *Manager) startPlugin(ctx context.Context, plugin *al_proto.PluginConfig) (*al_proto.PluginStartResponse, error) {
	pluginJson, err := protojson.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("could not marshal plugin: %w", err)
	}
	self.logger.Printf("Starting plugin: %s", string(pluginJson))
	bin, err := self.run.Rlocation(plugin.Bin)
	if err != nil {
		return nil, fmt.Errorf("could not get rlocation of %s: %w", plugin.Bin, err)
	}
	cmd := exec.CommandContext(ctx, bin, plugin.Args...)
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
			self.logger.Printf("Plugin %s finished with an error: %s", plugin.Name, err.Error())
		} else {
			self.logger.Printf("Plugin %s finished", plugin.Name)
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
	self.logger.Printf("Finished starting plugin %s", plugin.Name)
	return response, nil
}
