package main

import (
	"context"
	"fmt"
	"log"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

type App struct {
	ctx     *al.CmdCtx
	logger  *log.Logger
	cleaner *Cleaner
	server  *al_plugin.PluginServer
}

func NewApp(ctx *al.CmdCtx) (*App, error) {
	plugin := &Plugin{}
	server, err := al_plugin.NewPluginServer(ctx.Stdin, ctx.Stdout, ctx.Stderr, plugin)
	if err != nil {
		return nil, fmt.Errorf("could not create the plugin server: %w", err)
	}
	return &App{
		ctx:     ctx,
		logger:  log.New(ctx.Stderr, "com.alwaldend.src.tools.vault.injector.App ", ctx.LogFlags),
		cleaner: NewCleaner(ctx),
		server:  server,
	}, nil
}

func (self *App) Start(ctx context.Context) error {
	if err := self.server.Start(); err != nil {
		return fmt.Errorf("could not start the plugin server: %w", err)
	}
	return nil
}

func (self *App) Shutdown(ctx context.Context) error {
	var wg fp.WaitGroupE
	wg.Go(func() error {
		if err := self.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("the plugin server shut down with an error: %w", err)
		}
		return nil
	})
	wg.Go(func() error {
		if err := self.cleaner.Shutdown(ctx); err != nil {
			err = fmt.Errorf("cleaner shut down with an error: %w", err)
		}
		return nil
	})
	if err := wg.WaitCtx(ctx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}
