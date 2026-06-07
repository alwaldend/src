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
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Manager struct {
	logger *log.Logger
	wg     *fp.WaitGroupE
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
		run:    run,
		wg:     &fp.WaitGroupE{},
		logger: log.New(stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.Manager ", log.Flags()),
	}, nil
}

func (self *Manager) Shutdown() error {
	return self.wg.Wait()
}

func (self *Manager) StartPlugins(ctx context.Context, labels []string) ([]*al_proto.PluginStartResponse, error) {
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
			return nil, fmt.Errorf("duplicate plugin: %s", plugin.Name)
		}
		plugins[plugin.Name] = plugin
	}
	for _, call := range self.config.PluginCalls {
		if !labelsMatch(call.Labels, labelsMap) {
			continue
		}
		if call.Plugin == "" {
			return nil, fmt.Errorf("call %s is missing plugin name", call.Name)
		}
		plugin, ok := plugins[call.Plugin]
		if !ok {
			var err error
			plugin, err = fp.Find(func(_ int, p *al_proto.PluginConfig) bool {
				return p.Name == call.Plugin
			})(self.config.Plugins).Get()
			if err != nil {
				return nil, fmt.Errorf("could not find plugin %s: %w", call.Plugin, err)
			}
			plugins[plugin.Name] = plugin
		}
		plugin.Calls = append(plugin.Calls, call)
	}
	res, err := self.startPlugins(ctx, slices.Collect(maps.Values(plugins)))
	if err != nil {
		return nil, fmt.Errorf("could not start plugins: %w", err)
	}
	return res, nil
}

func (self *Manager) startPlugins(ctx context.Context, plugins []*al_proto.PluginConfig) ([]*al_proto.PluginStartResponse, error) {
	resCh := make(chan fp.Either[*al_proto.PluginStartResponse], len(plugins))
	for _, plugin := range plugins {
		func() {
			resp, err := self.startPlugin(ctx, plugin)
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
		return nil, fmt.Errorf("some plugins failed: %w", err)
	}
	return res, nil

}

func (self *Manager) startPlugin(ctx context.Context, plugin *al_proto.PluginConfig) (*al_proto.PluginStartResponse, error) {
	pluginJson, err := protojson.MarshalOptions{}.Marshal(plugin)
	if err != nil {
		return nil, fmt.Errorf("could not marshal plugin: %w", err)
	}
	self.logger.Printf("Starting plugin %s: %s", plugin.Name, string(pluginJson))
	bin, err := self.run.Rlocation(plugin.Bin)
	if err != nil {
		return nil, fmt.Errorf("could not get rlocation of %s: %w", plugin.Bin, err)
	}
	cmd := exec.CommandContext(ctx, bin, plugin.Args...)
	cmd.WaitDelay = time.Second * 10
	cmd.Cancel = func() error {
		return syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}
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
		Config: self.config,
		Plugin: plugin,
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start plugin: %w", err)
	}
	self.wg.Go(func() error {
		if err := cmd.Wait(); err != nil {
			err = fmt.Errorf("plugin %s finished with an error: %w", plugin.Name, err)
			self.logger.Println(err)
			return err
		} else {
			self.logger.Printf("plugin %s finished", plugin.Name)
		}
		return nil
	})
	response, err := client.PluginStart(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not execute plugin start request: %w", err)
	}
	if err := conn.Close(); err != nil {
		return nil, fmt.Errorf("could not close the connection")
	}
	self.logger.Printf("Finished starting plugin %s", plugin.Name)
	return response, nil
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
