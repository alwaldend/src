package worker

import (
	"context"
	"errors"
	"fmt"
	"io"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
)

const (
	// serveEvent starts the listener.
	serveEvent machine.Event = "serve"
	// submitEvent forwards a SubmitJobRequest.
	submitEvent machine.Event = "submit"
	// statusEvent forwards a JobReference to query a job.
	statusEvent machine.Event = "status"
	// cancelEvent forwards a JobReference to cancel a job.
	cancelEvent machine.Event = "cancel"
)

const (
	// stateReady is the initial state of a worker invocation.
	stateReady machine.State = "worker-ready"
	// stateServing invokes the worker service entry point.
	stateServing machine.State = "worker-serve"
	// stateForwardRPC selects the worker RPC entry point.
	stateForwardRPC machine.State = "forward-rpc"
)

// workerRun holds the inputs and state of a worker invocation.
type workerRun struct {
	// state is the worker command's current state.
	state machine.State
	// operation selects the worker transition.
	operation machine.Event
	// directory contains the worker's socket and job records.
	directory string
	// input supplies worker RPC request bytes.
	input io.Reader
	// output receives the worker RPC response.
	output io.Writer
}

// Worker owns rental-local worker operations.
type Worker struct {
	// definition dispatches worker commands.
	definition *machine.Machine[*workerRun]
}

// New configures the worker state machine.
func New() *Worker {
	definition := machine.New[*workerRun](stateReady)
	worker := &Worker{definition: definition}
	definition.Configure(stateReady).
		Permit(serveEvent, stateServing).
		Permit(submitEvent, stateForwardRPC).
		Permit(statusEvent, stateForwardRPC).
		Permit(cancelEvent, stateForwardRPC)
	definition.Configure(stateServing).OnEntry(serveWorker)
	definition.Configure(stateForwardRPC).OnEntry(forwardRPC)
	return worker
}

// Serve starts the rental-local listener.
func (worker *Worker) Serve(ctx context.Context, directory string) error {
	if err := worker.run(ctx, workerRun{operation: serveEvent, directory: directory}); err != nil {
		return fmt.Errorf("serve worker: %w", err)
	}
	return nil
}

// Submit forwards a protobuf job submission to the listener.
func (worker *Worker) Submit(ctx context.Context, directory string, input io.Reader, output io.Writer) error {
	if err := worker.run(ctx, workerRun{operation: submitEvent, directory: directory, input: input, output: output}); err != nil {
		return fmt.Errorf("submit to worker: %w", err)
	}
	return nil
}

// Status queries the listener for the persisted state of a job.
func (worker *Worker) Status(ctx context.Context, directory string, input io.Reader, output io.Writer) error {
	if err := worker.run(ctx, workerRun{operation: statusEvent, directory: directory, input: input, output: output}); err != nil {
		return fmt.Errorf("query worker status: %w", err)
	}
	return nil
}

// Cancel forwards a job cancellation to the listener.
func (worker *Worker) Cancel(ctx context.Context, directory string, input io.Reader, output io.Writer) error {
	if err := worker.run(ctx, workerRun{operation: cancelEvent, directory: directory, input: input, output: output}); err != nil {
		return fmt.Errorf("cancel worker job: %w", err)
	}
	return nil
}

// Graph describes the worker command states.
func (worker *Worker) Graph() machine.GraphSection {
	return machine.GraphSection{
		ID: "worker", Title: "Worker", DOT: worker.definition.ToGraph(),
	}
}

// run fires the selected worker transition.
func (worker *Worker) run(ctx context.Context, current workerRun) error {
	current.state = stateReady
	ctx = worker.definition.WithInvocation(ctx, &current.state, &current)
	if err := worker.definition.FireCtx(ctx, current.operation); err != nil {
		return fmt.Errorf("run worker event %s: %w", current.operation, err)
	}
	return nil
}

// serveWorker reports the worker service's current implementation status.
func serveWorker(_ context.Context, _ ...any) error {
	return errors.New("worker service is not implemented")
}

// forwardRPC reports the worker RPC client's current implementation status.
func forwardRPC(_ context.Context, _ ...any) error {
	return errors.New("worker RPC client is not implemented")
}
