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
	ctx    *al.CmdCtx
	server *al_plugin.PluginServer
	plugin *Plugin
	logger *log.Logger
}

func NewApp(ctx *al.CmdCtx) (*App, error) {
	plugin := &Plugin{ctx: ctx}
	server, err := al_plugin.NewPluginServer(ctx.Stdin, ctx.Stdout, ctx.Stderr, plugin)
	if err != nil {
		return nil, fmt.Errorf("could not create the plugin server: %w", err)
	}
	return &App{
		ctx:    ctx,
		server: server,
		plugin: plugin,
		logger: log.New(ctx.Stderr, "com.alwaldend.src.tools.vault.tf_backend.App ", ctx.LogFlags),
	}, nil
}

func (self *App) Shutdown(ctx context.Context) error {
	var wg fp.WaitGroupE
	self.logger.Printf("Shutting down the app")
	wg.Go(func() error {
		if err := self.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("plugin server error: %w", err)
		}
		return nil
	})
	wg.Go(func() error {
		if err := self.plugin.Shutdown(ctx); err != nil {
			return fmt.Errorf("plugin error: %w", err)
		}
		return nil
	})
	if err := wg.WaitCtx(ctx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}

func (self *App) Start(ctx context.Context) error {
	self.logger.Printf("Starting the app")
	if err := self.server.Start(); err != nil {
		return fmt.Errorf("could not start the plugin server: %w", err)
	}
	return nil
}
