package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker"
)

// main runs the root CLI machine and reports command failures.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	rental := worker.New()
	application := app.New(rental.Run)
	cli := newMachine(application.Run, func(definition machine.Definition) (string, error) {
		return composeStateMachineGraph(definition, application.Definition(), rental.Definition())
	})
	if err := cli.run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
