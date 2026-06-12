package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

func run(ctx *al.CmdCtx) error {
	cleaner := NewCleaner(ctx)
	plugin := NewPlugin(ctx, cleaner)
	server := al_plugin.NewPluginServer(ctx, plugin)
	var lc lifecycle.Manager
	lc.Add(cleaner, server)
	return lc.Run(ctx.Ctx, time.Second*10)
}

func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx, "com.alwaldend.src.tools.vault.injector ")
	if err := run(ctx); err != nil {
		ctx.Logger.Printf("error: %s", err)
		os.Exit(1)
	}
}
