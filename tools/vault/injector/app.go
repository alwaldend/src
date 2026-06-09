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
	errsCh := make(chan error, 2)
	errs := []error{}
	go func() {
		err := self.server.Shutdown(ctx)
		if err != nil {
			err = fmt.Errorf("the plugin server shut down with an error: %w", err)
		}
		errsCh <- err
	}()
	go func() {
		err := self.cleaner.Shutdown(ctx)
		if err != nil {
			err = fmt.Errorf("cleaner shut down with an error: %w", err)
		}
		errsCh <- err
	}()
	for range len(errsCh) {
		select {
		case err := <-errsCh:
			errs = append(errs, err)
		case <-ctx.Done():
			return fmt.Errorf("shutdown timed out")
		}
	}
	return errors.Join(errs...)
}
