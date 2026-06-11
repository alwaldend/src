package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"slices"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

type pluginState struct {
	client *PluginClient
	start  *al_proto.PluginStartResponse
	config *al_proto.PluginConfig
	calls  map[string]*al_proto.PluginCall
}

type Manager struct {
	logger    *log.Logger
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
		logger:  log.New(ctx.Stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.Manager ", ctx.LogFlags),
		ctx:     ctx,
		plugins: map[string]*pluginState{},
	}, nil
}

func (self *Manager) Stop(ctx context.Context) error {
	if err := self.lc.Stop(ctx); err != nil {
		return fmt.Errorf("stop error: %w", err)
	}
	return nil
}

func (self *Manager) Start(ctx context.Context) error {
	if err := self.lc.Start(ctx); err != nil {
		return fmt.Errorf("stop error: %w", err)
	}
	return nil
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
		self.lc.AddRunnable(state.client)
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
			plugin = &pluginState{
				config: pluginConfig,
				calls:  map[string]*al_proto.PluginCall{},
			}
			self.plugins[pluginConfig.Name] = plugin
		}
		if _, ok := plugin.calls[call.Name]; !ok {
			plugin.calls[call.Name] = call
		}
	}
	return nil
}

func (self *Manager) Start(ctx context.Context) error {
	res, err := self.startPlugins(ctx, slices.Collect(maps.Values(plugins)))
	if err != nil {
		return fmt.Errorf("could not start plugins: %w", err)
	}
	return res
}

func (self *Manager) newPlugin(config *al_proto.PluginConfig) (*pluginState, error) {
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
