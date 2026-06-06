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
	"slices"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Manager struct {
	cmds   []*exec.Cmd
	logger *log.Logger
	config *al_proto.Config
	run    *runfiles.Runfiles
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
		config: config,
		cmds:   []*exec.Cmd{},
		run:    run,
		logger: log.New(stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.Manager ", log.Flags()),
	}, nil
}

func (self *Manager) StartPlugins(ctx context.Context, labels []string) fp.Either[[]*al_proto.PluginStartResponse] {
	left := func(err error) fp.Either[[]*al_proto.PluginStartResponse] {
		return fp.Left[[]*al_proto.PluginStartResponse](err)
	}
	plugins := map[string]*al_proto.PluginConfig{}
	labelsMap := map[string]string{}
	for _, labelFull := range labels {
		labelSplit := strings.SplitN(labelFull, "=", 2)
		labelKey, labelVal := labelSplit[0], ""
		if len(labelSplit) > 1 {
			labelVal = labelSplit[1]
		}
		labelsMap[labelKey] = labelVal
	}
	for _, plugin := range self.config.Plugins {
		if !labelsMatch(plugin.Labels, labelsMap) {
			continue
		}
		_, ok := plugins[plugin.Name]
		if ok {
			return left(fmt.Errorf("duplicate plugin: %s", plugin.Name))
		}
		plugins[plugin.Name] = plugin
	}
	for _, call := range self.config.PluginCalls {
		if !labelsMatch(call.Labels, labelsMap) {
			continue
		}
		if call.Plugin == "" {
			return left(fmt.Errorf("call %s is missing plugin name", call.Name))
		}
		plugin, ok := plugins[call.Plugin]
		if !ok {
			var err error
			plugin, err = fp.Find(func(_ int, p *al_proto.PluginConfig) bool {
				return p.Name == call.Plugin
			})(self.config.Plugins).Get()
			if err != nil {
				return left(fmt.Errorf("could not find plugin %s: %w", call.Plugin, err))
			}
			plugins[plugin.Name] = plugin
		}
		plugin.Calls = append(plugin.Calls, call)
	}
	res, err := self.startPlugins(ctx, slices.Collect(maps.Values(plugins))).Get()
	if err != nil {
		return left(fmt.Errorf("could not start plugins: %w", err))
	}
	return fp.Right(res)
}

func (self *Manager) startPlugins(ctx context.Context, plugins []*al_proto.PluginConfig) fp.Either[[]*al_proto.PluginStartResponse] {
	resCh := make(chan fp.Either[*al_proto.PluginStartResponse], len(plugins))
	for _, plugin := range plugins {
		func() {
			resp, err := self.startPlugin(ctx, plugin).Get()
			if err != nil {
				err = fmt.Errorf("could not start plugin %s: %w", plugin.Name, err)
			}
			resCh <- fp.EitherOf(resp, err)
		}()
	}
	var err error
	res := []*al_proto.PluginStartResponse{}
	for range len(plugins) {
		val, curErr := (<-resCh).Get()
		if curErr != nil {
			err = errors.Join(err, curErr)
		} else {
			res = append(res, val)
		}
	}
	close(resCh)
	if err != nil {
		return fp.Left[[]*al_proto.PluginStartResponse](fmt.Errorf("some plugins failed: %w", err))
	}
	return fp.Right(res)

}

func (self *Manager) startPlugin(ctx context.Context, plugin *al_proto.PluginConfig) fp.Either[*al_proto.PluginStartResponse] {
	left := func(err error) fp.Either[*al_proto.PluginStartResponse] {
		return fp.Left[*al_proto.PluginStartResponse](err)
	}
	pluginJson, err := protojson.Marshal(plugin)
	if err != nil {
		return left(fmt.Errorf("could not marshal plugin: %w", err))
	}
	self.logger.Printf("Starting plugin %s: %s\n", plugin.Name, string(pluginJson))
	bin, err := self.run.Rlocation(plugin.Bin)
	if err != nil {
		return left(fmt.Errorf("could not get rlocation of %s: %w", plugin.Bin, err))
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
		return left(fmt.Errorf("could not create stdout pipe: %w", err))
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return left(fmt.Errorf("could not create stdin pipe: %w", err))
	}
	err = cmd.Start()
	if err != nil {
		return left(fmt.Errorf("could not start plugin: %w", err))
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
		return left(fmt.Errorf("could not create grpc client: %w", err))
	}
	client := al_proto.NewPluginServiceClient(conn)
	request := &al_proto.PluginStartRequest{
		Config: self.config,
		Plugin: plugin,
	}
	response, err := client.PluginStart(ctx, request)
	if err != nil {
		return left(fmt.Errorf("could not execute plugin start request: %w", err))
	}
	self.logger.Printf("Finished starting plugin %s", plugin.Name)
	return fp.Right(response)
}

// Return true if any labels in `want` match labels in `have`
func labelsMatch(have map[string]string, want map[string]string) bool {
	for k, v := range want {
		vHave, ok := have[k]
		if ok && (v == "" || v == vHave) {
			return true
		}
	}
	return false
}
