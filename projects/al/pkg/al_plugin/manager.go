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
	"github.com/bazelbuild/rules_go/go/runfiles"
)

type Manager struct {
	logger  *log.Logger
	config  *al_proto.Config
	run     *runfiles.Runfiles
	plugins []*PluginClient
	stderr  io.Writer
	ctx     *al.CmdCtx
	mx      sync.Mutex
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
		config: config,
		run:    run,
		logger: log.New(ctx.Stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.Manager ", ctx.LogFlags),
		ctx:    ctx,
	}, nil
}

func (self *Manager) Shutdown(ctx context.Context) error {
	var wg fp.WaitGroupE
	for _, plugin := range self.plugins {
		wg.Go(func() error {
			if err := plugin.Shutdown(ctx); err != nil {
				return fmt.Errorf("plugin %s shut down with an error: %w", plugin.req.Plugin.Name, err)
			}
			return nil
		})
	}
	if err := wg.WaitCtx(ctx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}

func (self *Manager) Start(ctx context.Context, labelArgs []string) ([]*al_proto.PluginStartResponse, error) {
	plugins := map[string]*al_proto.PluginConfig{}
	labels := labelsMap(labelArgs)
	for _, plugin := range self.config.Plugins {
		if !labelsMatch(plugin.Labels, labels) {
			continue
		}
		if _, ok := plugins[plugin.Name]; ok {
			return nil, fmt.Errorf("duplicate plugin: %s", plugin.Name)
		}
		plugins[plugin.Name] = plugin
	}
	for _, call := range self.config.PluginCalls {
		if !labelsMatch(call.Labels, labels) {
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

func (self *Manager) startPlugin(ctx context.Context, plugin *al_proto.PluginConfig) (*al_proto.PluginStartResponse, error) {
	client, err := NewPluginClient(self.ctx, self.run, self.config, plugin)
	if err != nil {
		return nil, fmt.Errorf("could not create client: %w", err)
	}
	resp, err := client.Start(ctx)
	if err != nil {
		err = fmt.Errorf("could not start: %w", err)
	}
	self.mx.Lock()
	defer self.mx.Unlock()
	self.plugins = append(self.plugins, client)
	return resp, nil
}

func (self *Manager) startPlugins(ctx context.Context, plugins []*al_proto.PluginConfig) ([]*al_proto.PluginStartResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	res := make(chan fp.Either[*al_proto.PluginStartResponse], len(plugins))
	for _, plugin := range plugins {
		func() {
			resp, err := self.startPlugin(ctx, plugin)
			if err != nil {
				err = fmt.Errorf("could not start plugin %s: %w", plugin.Name, err)
				cancel()
			}
			res <- fp.EitherOf(resp, err)
		}()
	}
	var err error
	resp := []*al_proto.PluginStartResponse{}
	for curRes := range res {
		val, curErr := curRes.Get()
		if curErr != nil {
			err = errors.Join(err, curErr)
		} else {
			resp = append(resp, val)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("some plugins failed: %w", err)
	}
	return resp, nil

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
