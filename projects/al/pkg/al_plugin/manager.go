package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/protobuf/proto"
)

// pluginState retains one selected plugin and its configured calls.
type pluginState struct {
	// client owns the plugin process and RPC connection.
	client *PluginClient
	// config is a private copy with selected calls keyed by name.
	config *al_proto.PluginConfig
	// name is the selected executable declaration key.
	name string
}

// Manager selects plugins and keeps their resources alive until consumers stop.
type Manager struct {
	// engine owns generic injection resources independently of provider processes.
	engine *injector.Engine
	// injectionCalls retains selected groups by their configuration keys.
	injectionCalls map[string]*al_proto.InjectorCall
	// injectionEnv holds exports after graph evaluation.
	injectionEnv map[string]string
	// orchestration exposes phased startup and ordered shutdown to the caller.
	orchestration *lifecycle.Manager
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
	if err := al.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("validate plugin configuration: %w", err)
	}
	manager := &Manager{
		config:         config,
		run:            run,
		ctx:            ctx,
		plugins:        map[string]*pluginState{},
		lc:             lifecycle.NewParallelStopManager(),
		orchestration:  lifecycle.NewParallelStopManager(),
		engine:         injector.New(),
		injectionCalls: map[string]*al_proto.InjectorCall{},
	}
	if err := manager.orchestration.Add(manager); err != nil {
		return nil, fmt.Errorf("register injection lifecycle: %w", err)
	}
	return manager, nil
}

// Lifecycle exposes the selected plugin resources to their enclosing command.
func (self *Manager) Lifecycle() *lifecycle.Manager {
	return self.orchestration
}

// Env combines exports after selected plugin startup has completed. Equal values
// coalesce; explicit injection overrides startup defaults. Other conflicting
// startup exports fail deterministically without disclosing values.
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
		if resp := plugin.client.StartResponse().GetResult(); resp != nil {
			keys := make([]string, 0, len(resp.Env))
			for key := range resp.Env {
				keys = append(keys, key)
			}
			slices.Sort(keys)
			for _, key := range keys {
				if _, overridden := self.injectionEnv[key]; overridden {
					continue
				}
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
	maps.Copy(values, self.injectionEnv)
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

// AddLabels selects plugins and dedicated injector calls by label.
func (self *Manager) AddLabels(labelArgs []string) error {
	self.pluginsMx.Lock()
	defer self.pluginsMx.Unlock()
	labels := labelsMap(labelArgs)
	for name, config := range self.config.PluginConfigs {
		if _, exists := self.plugins[name]; exists || !labelsMatch(config.Labels, labels) {
			continue
		}
		state, err := self.newPlugin(name, config)
		if err != nil {
			return fmt.Errorf("create plugin %q: %w", name, err)
		}
		self.plugins[name] = state
		if err := self.lc.Add(state.client); err != nil {
			return fmt.Errorf("register plugin %q: %w", name, errors.Join(err, state.client.Stop(context.Background())))
		}
	}
	for name, call := range self.config.Plugins {
		if !labelsMatch(call.Labels, labels) {
			continue
		}
		if call.Plugin == "" {
			return fmt.Errorf("call %q is missing the plugin name", name)
		}
		plugin, exists := self.plugins[call.Plugin]
		if !exists {
			config := self.config.PluginConfigs[call.Plugin]
			if config == nil {
				return fmt.Errorf("could not find plugin %q", call.Plugin)
			}
			var err error
			plugin, err = self.newPlugin(call.Plugin, config)
			if err != nil {
				return fmt.Errorf("create plugin for call %q: %w", name, err)
			}
			self.plugins[call.Plugin] = plugin
			if err := self.lc.Add(plugin.client); err != nil {
				return fmt.Errorf("register plugin for call %q: %w", name, errors.Join(err, plugin.client.Stop(context.Background())))
			}
		}
		plugin.config.Calls[name] = call
	}
	for name, call := range self.config.Inject {
		if labelsMatch(call.Labels, labels) {
			self.injectionCalls[name] = call
		}
	}

	return nil
}

// newPlugin creates an owned client from an independent copy of its declaration.
func (self *Manager) newPlugin(name string, config *al_proto.PluginConfig) (*pluginState, error) {
	config = proto.CloneOf(config)
	if config.Calls == nil {
		config.Calls = map[string]*al_proto.PluginCall{}
	}
	client, err := NewPluginClient(self.ctx, self.run, self.config, name, config)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %w", err)
	}
	return &pluginState{
		config: config,
		name:   name,
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

// Start starts selected environment plugins and only discovers providers when needed.
func (self *Manager) Start(ctx context.Context) error {
	if err := self.engine.Load(self.injectionCalls, self.config.Inject); err != nil {
		return fmt.Errorf("load injector calls: %w", err)
	}
	if self.engine.NeedsProviders() {
		for name, config := range self.config.PluginConfigs {
			if !config.ResourceProvider {
				continue
			}
			if _, exists := self.plugins[name]; exists {
				continue
			}
			state, err := self.newPlugin(name, config)
			if err != nil {
				return fmt.Errorf("create resource provider %q: %w", name, err)
			}
			self.plugins[name] = state
			if err := self.lc.Add(state.client); err != nil {
				return fmt.Errorf("register resource provider %q: %w", name, errors.Join(err, state.client.Stop(ctx)))
			}
		}
	}
	if err := self.lc.Start(ctx); err != nil {
		return fmt.Errorf("start selected plugins: %w", err)
	}
	for _, name := range slices.Sorted(maps.Keys(self.plugins)) {
		state := self.plugins[name]
		types := state.client.StartResponse().GetResourceTypes()
		if len(types) > 0 && !state.config.ResourceProvider {
			return fmt.Errorf("plugin %q advertised resources without resource_provider", name)
		}
		if state.config.ResourceProvider && len(types) == 0 {
			return fmt.Errorf("resource provider %q advertised no resource types", name)
		}
		for _, kind := range types {
			if err := self.engine.Register(kind, injector.FetchFunc(state.fetchResource)); err != nil {
				return fmt.Errorf("register plugin %q: %w", name, err)
			}
		}
	}
	env, err := self.engine.Evaluate(ctx)
	if err != nil {
		return fmt.Errorf("evaluate selected resources: %w", err)
	}
	self.injectionEnv = env
	return nil
}

// fetchResource adapts a completed graph node and its ordered dependencies to RPC.
func (self *pluginState) fetchResource(ctx context.Context, resource *injector.Resource, deps []*injector.ResourceResult) (*injector.ResourceResult, error) {
	args, err := al.ToPbJson(resource.Args).Get()
	if err != nil {
		return nil, fmt.Errorf("encode provider arguments: %w", err)
	}
	data, err := al.ToPbJson(resource.Data).Get()
	if err != nil {
		return nil, fmt.Errorf("encode provider data: %w", err)
	}
	request := &al_proto.ResourceFetchRequest{Name: resource.Name, Type: resource.Type, Args: args, Data: data}
	for _, result := range deps {
		value, err := injector.ToProto(result)
		if err != nil {
			return nil, fmt.Errorf("encode provider dependency %q: %w", result.Name, err)
		}
		request.Dependencies = append(request.Dependencies, value)
	}
	response, err := self.client.FetchResource(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("fetch from provider %q: %w", self.name, err)
	}
	result, err := injector.FromProto(response)
	if err != nil {
		return nil, fmt.Errorf("decode provider %q response: %w", self.name, err)
	}
	return result, nil
}

// Stop releases built-in consumers before terminating resource provider processes.
func (self *Manager) Stop(ctx context.Context) error {
	return errors.Join(self.engine.Stop(ctx), self.lc.Stop(ctx))
}
