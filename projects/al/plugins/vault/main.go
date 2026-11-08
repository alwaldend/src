package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/injector"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

// run serves provider RPCs until cancellation and drains before cleanup.
func run(ctx *al.CmdCtx) error {
	cleaner := injector.NewCleaner()
	plugin := NewPlugin(ctx, cleaner)
	server := al_plugin.NewPluginServer(ctx, plugin)
	var lc lifecycle.Manager
	if err := lc.Add(cleaner, plugin, server); err != nil {
		return fmt.Errorf("register Vault provider lifecycle: %w", err)
	}
	if err := lc.Run(ctx.Ctx, time.Second*10); err != nil {
		return fmt.Errorf("run Vault provider: %w", err)
	}
	return nil
}

// main connects process signals to the provider lifecycle.
func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.projects.al.plugins.vault ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("error: %s", err)
		os.Exit(1)
	}
}
