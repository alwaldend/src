package al_plugin

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/proto"
)

type pluginState struct {
	client *PluginClient
	start  *al_proto.PluginStartResponse
	config *al_proto.PluginConfig
	calls  map[string]*al_proto.PluginCall
}

type Manager struct {
	config    *al_proto.Config
	run       *runfiles.Runfiles
	plugins   map[string]*pluginState
	pluginsMx sync.Mutex
	stderr    io.Writer
	ctx       *al.CmdCtx
	lc        lifecycle.Manager
}

func NewManager(ctx *al.CmdCtx, config *al_proto.Config) (*Manager, error) {
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
		run:     run,
		ctx:     ctx,
		plugins: map[string]*pluginState{},
	}, nil
}

func (self *Manager) Lifecycle() *lifecycle.Manager {
	return &self.lc
}

func (self *Manager) Env() []string {
	res := []string{}
	for _, plugin := range self.plugins {
		if resp, ok := plugin.client.StartResponse(); ok {
			for k, v := range resp.Env {
				res = append(res, fmt.Sprintf("%s=%s", k, v))
			}
		}

	}
	return res
}

func (self *Manager) AddLabels(labelArgs []string) error {
	self.pluginsMx.Lock()
	defer self.pluginsMx.Unlock()
	labels := labelsMap(labelArgs)
	for _, plugin := range self.config.Plugins {
		if _, ok := self.plugins[plugin.Name]; ok {
			continue
		}
		if !labelsMatch(plugin.Labels, labels) {
			continue
		}
		state, err := self.newPlugin(plugin)
		if err != nil {
			return fmt.Errorf("could not create plugin %s: %w", plugin.Name, err)
		}
		self.plugins[plugin.Name] = state
		self.lc.Add(state.client)
	}
	for i, call := range self.config.PluginCalls {
		if !labelsMatch(call.Labels, labels) {
			continue
		}
		if call.Name == "" {
			return fmt.Errorf("call %d is missing a name", i)
		}
		if call.Plugin == "" {
			return fmt.Errorf("call %s is missing the plugin name", call.Name)
		}
		plugin, ok := self.plugins[call.Plugin]
		if !ok {
			pluginConfig, err := fp.Find(func(_ int, p *al_proto.PluginConfig) bool {
				return p.Name == call.Plugin
			})(self.config.Plugins).Get()
			if err != nil {
				return fmt.Errorf("could not find plugin %s: %w", call.Plugin, err)
			}
			plugin, err = self.newPlugin(pluginConfig)
			if err != nil {
				return fmt.Errorf("could not create plugin for call %s: %w", call.Name, err)
			}
			self.plugins[pluginConfig.Name] = plugin
			self.lc.Add(plugin.client)
		}
		if _, ok := plugin.calls[call.Name]; !ok {
			plugin.calls[call.Name] = call
			plugin.config.Calls = append(plugin.config.Calls, call)
		}
	}
	return nil
}

func (self *Manager) newPlugin(config *al_proto.PluginConfig) (*pluginState, error) {
	config = proto.CloneOf(config)
	client, err := NewPluginClient(self.ctx, self.run, self.config, config)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %w", err)
	}
	return &pluginState{
		config: config,
		calls:  map[string]*al_proto.PluginCall{},
		client: client,
	}, nil
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

func labelsMap(labels []string) map[string]string {
	res := map[string]string{}
	for _, labelFull := range labels {
		labelSplit := strings.SplitN(labelFull, "=", 2)
		labelKey, labelVal := labelSplit[0], ""
		if len(labelSplit) > 1 {
			labelVal = labelSplit[1]
		}
		res[labelKey] = labelVal
	}
	return res
}
