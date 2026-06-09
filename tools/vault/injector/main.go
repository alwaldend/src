package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

func main() {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	ctx := al.NewCmdCtx(shutdownCtx)
	logger := log.New(ctx.Stderr, "com.alwaldend.src.tools.vault.injector.global ", ctx.LogFlags)
	app, err := NewApp(ctx)
	if err != nil {
		logger.Printf("could not create the app: %s", err)
		os.Exit(1)
	}
	if err := fp.Run(shutdownCtx, app.Start, app.Shutdown, time.Second*10); err != nil {
		logger.Printf("could not run the app: %s", err)
		os.Exit(1)
	}
}
