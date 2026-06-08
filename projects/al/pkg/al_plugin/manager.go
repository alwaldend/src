package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

type Manager struct {
	logger  *log.Logger
	config  *al_proto.Config
	run     *runfiles.Runfiles
	plugins []*PluginClient
	stderr  io.Writer
	ctx     context.Context
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
		logger: log.New(stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.Manager ", log.Flags()),
		stderr: stderr,
	}, nil
}

func (self *Manager) Shutdown(ctx context.Context) error {
	err := make(chan error, len(self.plugins))
	done := make(chan struct{})
	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %w", err)
	}
	return nil
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
	ctx, cancel := context.WithCancel(ctx)
	mx := sync.Mutex{}
	res := []fp.Either[*al_proto.PluginStartResponse]{}
	for _, plugin := range plugins {
		func() {
			var err error
			var client *PluginClient
			var resp *al_proto.PluginStartResponse
			client, err = NewPluginClient(self.stderr, self.run, self.config, plugin)
			if err != nil {
				err = fmt.Errorf("could not create plugin client for %s: %w", plugin.Name, err)
			}
			if err == nil {
				resp, err = client.Start(ctx)
				if err != nil {
					err = fmt.Errorf("could not start plugin %s: %w", err)
				}
			}
			mx.Lock()
			defer mx.Unlock()
			if err == nil {
				self.plugins = append(self.plugins, client)
			} else {
				cancel()
			}
			res = append(res, fp.EitherOf(resp, err))
		}()
	}
	var err error
	resp := []*al_proto.PluginStartResponse{}
	for _, curRes := range res {
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
