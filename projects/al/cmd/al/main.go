package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	os.Exit(Execute(al.NewCmdCtx(ctx, "com.alwaldend.src.projects.al.cmd.al ")))
}
