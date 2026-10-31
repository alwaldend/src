package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/proto"
)

// pluginState retains one selected plugin and its configured calls.
type pluginState struct {
	// client owns the plugin process and RPC connection.
	client *PluginClient
	// config is a private copy with selected calls appended.
	config *al_proto.PluginConfig
	// calls prevents registering a selected call twice.
	calls map[string]*al_proto.PluginCall
}

// Manager selects plugins and keeps their resources alive until consumers stop.
type Manager struct {
	// config supplies the explicit plugin and call declarations.
	config *al_proto.Config
	// run resolves relative plugin executables when Bazel runtime files are available.
	run *runfiles.Runfiles
	// plugins records selected plugins by their unique configuration names.
	plugins map[string]*pluginState
	// pluginsMx serializes plugin selection and export collection.
	pluginsMx sync.Mutex
	// ctx supplies command cancellation and borrowed streams.
	ctx *al.CmdCtx
	// lc owns the independent plugin process lifetimes.
	lc *lifecycle.Manager
}

// NewManager validates plugin identities without requiring an optional runfiles context.
func NewManager(ctx *al.CmdCtx, config *al_proto.Config) (*Manager, error) {
	run, err := al.FindRunfiles()
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
		lc:      lifecycle.NewParallelStopManager(),
	}, nil
}

// Lifecycle exposes the selected plugin resources to their enclosing command.
func (self *Manager) Lifecycle() *lifecycle.Manager {
	return self.lc
}

// Env combines exports after selected plugin startup has completed. Equal values
// coalesce; conflicting values fail deterministically without disclosing values.
// Inherited command environment is deliberately outside this conflict check.
func (self *Manager) Env() ([]string, error) {
	self.pluginsMx.Lock()
	defer self.pluginsMx.Unlock()
	names := make([]string, 0, len(self.plugins))
	for name := range self.plugins {
		names = append(names, name)
	}
	slices.Sort(names)
	values := map[string]string{}
	owners := map[string]string{}
	for _, name := range names {
		plugin := self.plugins[name]
		if resp, ok := plugin.client.StartResponse(); ok {
			keys := make([]string, 0, len(resp.Env))
			for key := range resp.Env {
				keys = append(keys, key)
			}
			slices.Sort(keys)
			for _, key := range keys {
				value := resp.Env[key]
				if previous, exists := values[key]; exists && previous != value {
					return nil, fmt.Errorf("conflicting environment variable %q from plugins %q and %q", key, owners[key], name)
				}
				if _, exists := values[key]; !exists {
					values[key], owners[key] = value, name
				}
			}
		}
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	res := make([]string, 0, len(keys))
	for _, key := range keys {
		res = append(res, key+"="+values[key])
	}
	return res, nil
}

// AddLabels registers plugins and calls matching any supplied label selector.
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
		if err := self.lc.Add(state.client); err != nil {
			return fmt.Errorf("register plugin: %w", errors.Join(err, state.client.Stop(context.Background())))
		}
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
			if err := self.lc.Add(plugin.client); err != nil {
				return fmt.Errorf("register plugin: %w", errors.Join(err, plugin.client.Stop(context.Background())))
			}
		}
		if _, ok := plugin.calls[call.Name]; !ok {
			plugin.calls[call.Name] = call
			plugin.config.Calls = append(plugin.config.Calls, call)
		}
	}
	return nil
}

// newPlugin creates an owned client from an independent copy of its declaration.
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

// labelsMatch reports whether any requested key and optional value match.
func labelsMatch(have map[string]string, want map[string]string) bool {
	for k, v := range want {
		vHave, ok := have[k]
		if ok && (v == "" || v == vHave) {
			return true
		}
	}
	return false
}

// labelsMap keeps the final selector value supplied for each label key.
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
