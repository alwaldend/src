package worker

import (
	"context"
	"errors"
	"fmt"
	"io"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
)

// Operation identifies a worker command.
type Operation string

const (
	// Serve starts the listener.
	Serve Operation = "serve"
	// Submit forwards a SubmitJobRequest.
	Submit Operation = "submit"
	// Status forwards a JobReference to query a job.
	Status Operation = "status"
	// Cancel forwards a JobReference to cancel a job.
	Cancel Operation = "cancel"
)

// Command supplies the inputs to a worker invocation.
type Command struct {
	// Operation selects the command to run.
	Operation Operation
	// Directory overrides the worker's state directory.
	Directory string
	// Input supplies worker RPC request bytes.
	Input io.Reader
	// Output receives the worker RPC response.
	Output io.Writer
}

const (
	// stateReady is the initial state of a worker invocation.
	stateReady machine.State = "worker-ready"
	// stateServing invokes the worker service entry point.
	stateServing machine.State = "worker-serve"
	// stateForwardRPC selects the worker RPC entry point.
	stateForwardRPC machine.State = "forward-rpc"
)

// Worker owns the machine for rental-local worker commands.
type Worker struct {
	// definition dispatches worker commands.
	definition *machine.Machine[struct{}]
}

// New configures the worker machine.
func New() *Worker {
	definition := machine.New[struct{}](stateReady)
	definition.Configure(stateReady).
		Permit(machine.Event(Serve), stateServing).
		Permit(machine.Event(Submit), stateForwardRPC).
		Permit(machine.Event(Status), stateForwardRPC).
		Permit(machine.Event(Cancel), stateForwardRPC)
	definition.Configure(stateServing).OnEntry(serveWorker)
	definition.Configure(stateForwardRPC).OnEntry(forwardRPC)
	return &Worker{definition: definition}
}

// Run invokes the worker service through its machine.
func (worker *Worker) Run(ctx context.Context, command Command) error {
	state := stateReady
	if err := worker.definition.Fire(ctx, &state, machine.Event(command.Operation), struct{}{}); err != nil {
		return fmt.Errorf("worker %s: %w", command.Operation, err)
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

// forwardRPC reports the worker RPC client's current implementation status.
func forwardRPC(_ context.Context, _ ...any) error {
	return errors.New("worker RPC client is not implemented")
}
