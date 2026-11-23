package worker

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
)

const (
	// stateReady is the initial state of a worker invocation.
	stateReady machine.State = "worker-ready"
	// stateServing invokes the worker service entry point.
	stateServing machine.State = "worker-serve"
	// serve selects the worker service entry point.
	serve machine.Event = "serve"
)

// Worker owns the machine for rental-local worker commands.
type Worker struct {
	// definition dispatches worker commands.
	definition *machine.Machine[struct{}]
}

// New configures the worker machine.
func New() *Worker {
	definition := machine.New[struct{}](stateReady)
	definition.Configure(stateReady).Permit(serve, stateServing)
	definition.Configure(stateServing).OnEntry(serveWorker)
	return &Worker{definition: definition}
}

// Run invokes the worker service through its machine.
func (worker *Worker) Run(ctx context.Context) error {
	state := stateReady
	if err := worker.definition.Fire(ctx, &state, serve, struct{}{}); err != nil {
		return fmt.Errorf("worker serve: %w", err)
	}
	return nil
}

// Definition exposes the configured worker states for inspection.
func (worker *Worker) Definition() machine.Definition {
	return worker.definition
}

// serveWorker reports the worker service's current implementation status.
func serveWorker(_ context.Context, _ ...any) error {
	return errors.New("worker service is not implemented")
}
