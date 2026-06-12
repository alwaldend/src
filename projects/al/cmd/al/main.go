package main

import (
	"context"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	os.Exit(Execute(al.NewCmdCtx(ctx, "com.alwaldend.src.projects.al.cmd.al ")))
}
