package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/internal/cli"
	"git.alwaldend.com/alwaldend/src/tools/rules_al/pkg/al"
)

// main owns process signals and delegates command execution to the CLI adapter.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	os.Exit(cli.Execute(al.NewCmdCtx(ctx, "al: ")))
}
