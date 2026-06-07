package main

import (
	"context"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	os.Exit(Execute(&al.CmdCtx{
		Ctx:    ctx,
		Args:   os.Args,
		Getenv: os.Getenv,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}))
}
