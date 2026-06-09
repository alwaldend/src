package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
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
	self.logger.Printf("Shutting down the app")
	errs := []error{}
	if err := self.server.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("could not shut down the plugin server: %w", err))
	}
	if err := self.plugin.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("could not shut down the plugin: %w", err))
	}
	return errors.Join(errs...)
}

func (self *App) Start(ctx context.Context) error {
	self.logger.Printf("Starting the app")
	if err := self.server.Start(); err != nil {
		return fmt.Errorf("could not start the plugin server: %w", err)
	}
	return nil
}
