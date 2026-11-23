package app

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker"
)

// CommandKind identifies the event chosen by the CLI.
type CommandKind = machine.Event

const (
	// ServeJobs starts the coordinator.
	ServeJobs CommandKind = "server serve"
	// RunWorker delegates a rental-local command to the worker machine.
	RunWorker CommandKind = "worker"
)

// Command carries an application request from its caller.
type Command struct {
	// Kind selects the application event to fire.
	Kind CommandKind
	// Worker holds the typed operation for the worker child.
	Worker worker.Command
}

const (
	// stateReady is the initial state of an invocation.
	stateReady machine.State = "app-ready"
	// stateServer selects coordinator service.
	stateServer machine.State = "server"
	// StateDispatchWorker runs a command on the worker child machine.
	StateDispatchWorker machine.State = "app-worker"
)

// App owns application dispatch and invokes its worker child.
type App struct {
	// definition dispatches typed application requests.
	definition *machine.Machine[worker.Command]
	// triggerWorker starts the worker child.
	triggerWorker func(context.Context, worker.Command) error
}

// New configures the application machine around its worker child.
func New(triggerWorker func(context.Context, worker.Command) error) *App {
	app := &App{definition: machine.New[worker.Command](stateReady), triggerWorker: triggerWorker}
	definition := app.definition
	definition.Configure(stateReady).Permit(ServeJobs, stateServer).Permit(RunWorker, StateDispatchWorker)
	definition.Configure(stateServer).OnEntry(serveJobs)
	definition.Configure(StateDispatchWorker).OnEntry(app.runWorker)
	return app
}

// Run handles a typed request on the application machine.
func (app *App) Run(ctx context.Context, command Command) error {
	state := stateReady
	if err := app.definition.Fire(ctx, &state, command.Kind, command.Worker); err != nil {
		return fmt.Errorf("%s: %w", command.Kind, err)
	}
	return nil
}

// Definition exposes the configured application states for inspection.
func (app *App) Definition() machine.Definition {
	return app.definition
}

// runWorker invokes the worker machine with the parsed operation.
func (app *App) runWorker(ctx context.Context, _ ...any) error {
	command, err := app.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read worker command: %w", err)
	}
	if err := app.triggerWorker(ctx, command); err != nil {
		return fmt.Errorf("run worker command: %w", err)
	}
	return nil
}

// serveJobs reports the coordinator's current implementation status.
func serveJobs(_ context.Context, _ ...any) error {
	return errors.New("job server is not implemented")
}
