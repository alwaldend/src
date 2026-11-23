package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker"
)

const (
	// serveJobsEvent starts the coordinator.
	serveJobsEvent machine.Event = "server serve"
	// runWorkerEvent delegates a rental-local operation.
	runWorkerEvent machine.Event = "worker"
)

const (
	// stateReady is the initial state of an invocation.
	stateReady machine.State = "app-ready"
	// stateServer selects coordinator service.
	stateServer machine.State = "server"
	// stateDispatchWorker runs a command on the worker.
	stateDispatchWorker machine.State = "app-worker"
)

// App owns application dispatch and invokes its worker child.
type App struct {
	// definition dispatches application requests.
	definition *machine.Machine[func(context.Context) error]
	// worker handles rental-local operations.
	worker *worker.Worker
}

// New configures the application around its worker.
func New(worker *worker.Worker) *App {
	app := &App{definition: machine.New[func(context.Context) error](stateReady), worker: worker}
	definition := app.definition
	definition.Configure(stateReady).Permit(serveJobsEvent, stateServer).Permit(runWorkerEvent, stateDispatchWorker)
	definition.Configure(stateServer).OnEntry(serveJobs)
	definition.Configure(stateDispatchWorker).OnEntry(app.runWorker)
	return app
}

// ServeJobs starts the job server.
func (app *App) ServeJobs(ctx context.Context) error {
	if err := app.run(ctx, serveJobsEvent, nil); err != nil {
		return fmt.Errorf("start job server: %w", err)
	}
	return nil
}

// ServeWorker starts the rental-local listener.
func (app *App) ServeWorker(ctx context.Context, directory string) error {
	err := app.run(ctx, runWorkerEvent, func(ctx context.Context) error {
		if err := app.worker.Serve(ctx, directory); err != nil {
			return fmt.Errorf("serve rental-local worker: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("start worker service: %w", err)
	}
	return nil
}

// Submit sends a job request to the rental-local listener.
func (app *App) Submit(ctx context.Context, directory string, input io.Reader, output io.Writer) error {
	if err := app.callWorker(ctx, directory, input, output, app.worker.Submit); err != nil {
		return fmt.Errorf("submit job: %w", err)
	}
	return nil
}

// Status queries a job's latest committed state.
func (app *App) Status(ctx context.Context, directory string, input io.Reader, output io.Writer) error {
	if err := app.callWorker(ctx, directory, input, output, app.worker.Status); err != nil {
		return fmt.Errorf("query job status: %w", err)
	}
	return nil
}

// Cancel sends a cancellation request to the rental-local listener.
func (app *App) Cancel(ctx context.Context, directory string, input io.Reader, output io.Writer) error {
	if err := app.callWorker(ctx, directory, input, output, app.worker.Cancel); err != nil {
		return fmt.Errorf("cancel job: %w", err)
	}
	return nil
}

// Graph describes the application and the call to its worker.
func (app *App) Graph() machine.GraphSection {
	return machine.GraphSection{
		ID: "application", Title: "Application", DOT: app.definition.ToGraph(),
		Calls: []machine.GraphCall{
			{From: stateDispatchWorker, Label: "runs worker command", Target: app.worker.Graph()},
		},
	}
}

// callWorker binds the request streams to an application invocation.
func (app *App) callWorker(ctx context.Context, directory string, input io.Reader, output io.Writer, operation func(context.Context, string, io.Reader, io.Writer) error) error {
	err := app.run(ctx, runWorkerEvent, func(ctx context.Context) error {
		if err := operation(ctx, directory, input, output); err != nil {
			return fmt.Errorf("call rental-local worker: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("dispatch worker request: %w", err)
	}
	return nil
}

// run fires the selected application event with its invocation input.
func (app *App) run(ctx context.Context, event machine.Event, operation func(context.Context) error) error {
	state := stateReady
	ctx = app.definition.WithInvocation(ctx, &state, operation)
	if err := app.definition.FireCtx(ctx, event); err != nil {
		return fmt.Errorf("run application event %s: %w", event, err)
	}
	return nil
}

// runWorker invokes the named worker operation selected by the application method.
func (app *App) runWorker(ctx context.Context, _ ...any) error {
	operation, err := app.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read worker invocation: %w", err)
	}
	if err := operation(ctx); err != nil {
		return fmt.Errorf("run worker command: %w", err)
	}
	return nil
}

// serveJobs reports the coordinator's current implementation status.
func serveJobs(_ context.Context, _ ...any) error {
	return errors.New("job server is not implemented")
}
