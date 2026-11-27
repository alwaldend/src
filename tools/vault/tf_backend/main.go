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
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

// run registers the plugin before its dependent RPC server for ordered shutdown.
func run(ctx *al.CmdCtx) error {
	var lc lifecycle.Manager
	plugin := newPlugin(ctx)
	server := al_plugin.NewPluginServer(ctx, plugin)
	if err := lc.Add(plugin, server); err != nil {
		return fmt.Errorf("register plugin and RPC server: %w", err)
	}
	if err := lc.Run(ctx.Ctx, time.Second*10); err != nil {
		return fmt.Errorf("run Terraform backend plugin: %w", err)
	}
	return nil
}

// main serves the AL plugin protocol until process shutdown.
func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.tf_backend ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("failed: %s", err)
		os.Exit(1)
	}
}
