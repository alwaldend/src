package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

type Plugin struct {
	ctx        context.Context
	v          *al.VaultStore
	wg         *sync.WaitGroup
	returnCode int
	logger     *log.Logger
	cleaner    *Cleaner
	manager    *ResourceManager
}

func (self Plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	vault := al.NewVault(req.Config)
	templater := &Templater{}
	self.manager = &ResourceManager{
		env:      &EnvFetcher{templater: templater},
		file:     &FileFetcher{vault: vault, templater: templater, cleaner: self.cleaner},
		op:       &OpFetcher{vault: vault},
		secret:   &SecretFetcher{vault: vault},
		vaultEnv: &VaultEnvFetcher{templater: templater, config: req.Config, vault: vault},
	}
	config := &injector_proto.Config{}
	if _, err := al.FromPbJsonToPb(req.Plugin.Data, config).Get(); err != nil {
		return nil, fmt.Errorf("could not parse plugin data: %w", err)
	}
	for _, call := range req.Plugin.Calls {
		curConfig := &injector_proto.Config{}
		if _, err := al.FromPbJsonToPb(call.Data, curConfig).Get(); err != nil {
			return nil, fmt.Errorf("could not parse plugin call data of %s: %w", call.Name, err)
		}
		config.Res = append(config.Res, curConfig.Res...)
	}
	resource, err := self.manager.Get(ctx, config).Get()
	if err != nil {
		return nil, fmt.Errorf("could not process request: %w", err)
	}
	res := &al_proto.PluginStartResponse{
		Env: resource.Env,
	}
	return res, nil
}

func (self *Plugin) Run() int {
	self.wg.Go(func() {
		self.cleaner.Run(self.ctx)
	})
	self.wg.Go(func() {
		if _, err := al_plugin.Serve(self.ctx, os.Stdin, os.Stdout, self).Get(); err != nil {
			self.logger.Printf("could not serve the plugin: %s", err)
			self.returnCode = 1
		}
	})
	self.wg.Wait()
	self.logger.Printf("Stopped running with code: %d\n", self.returnCode)
	return self.returnCode
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	p := &Plugin{
		wg:      &sync.WaitGroup{},
		logger:  log.New(os.Stderr, "com.alwaldend.src.tools.vault.injector.plugin ", log.Flags()),
		cleaner: NewCleaner(os.Stderr),
		ctx:     ctx,
	}
	go func() {
		<-ctx.Done()
		p.logger.Printf("Got a cancel signal")
	}()
	p.logger.Printf("Starting plugin")
	code := p.Run()
	p.logger.Printf("Stopped plugin, code: %d", code)
	os.Exit(code)
}
