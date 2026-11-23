package workload

import (
	"context"
	"fmt"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"github.com/spf13/afero"
)

const (
	// stateReady awaits a request to stage the project.
	stateReady machine.State = "workload-ready"
	// statePreparing stages the project and configures its command.
	statePreparing machine.State = "workload-preparing"
	// statePrepared awaits authorization from the caller.
	statePrepared machine.State = "workload-prepared"
	// stateRunning starts and observes the prepared command.
	stateRunning machine.State = "workload-running"
)

const (
	// prepareProjectCommand begins project staging.
	prepareProjectCommand machine.Event = "prepare"
	// preparedProjectCommand records successful project staging.
	preparedProjectCommand machine.Event = "prepared"
	// runProjectCommand runs the prepared command after authorization.
	runProjectCommand machine.Event = "run"
)

// invocation holds one workload's staging inputs and execution result.
type invocation struct {
	// mu serializes attempts to run the same prepared command.
	mu sync.Mutex
	// state is the current workload state.
	state machine.State
	// directory contains the project and its outputs.
	directory string
	// id identifies the MLflow run.
	id string
	// specification describes the project and its entry point.
	specification *jobpb.Workload
	// prepared holds the command after successful staging.
	prepared *prepared
	// outcome holds the facts observed while running the command.
	outcome Outcome
}

// Runner owns the workload machine shared by the listener's jobs.
type Runner struct {
	// definition applies transitions to individual workload invocations.
	definition *machine.Machine[*invocation]
	// filesystem stages projects and opens their logs.
	filesystem afero.Fs
}

// New configures the workload state machine.
func New(filesystem afero.Fs) *Runner {
	runner := &Runner{definition: machine.New[*invocation](stateReady), filesystem: filesystem}
	runner.definition.Configure(stateReady).Permit(prepareProjectCommand, statePreparing)
	runner.definition.Configure(statePreparing).OnEntry(runner.stageProject).
		Permit(preparedProjectCommand, statePrepared)
	runner.definition.Configure(statePrepared).Permit(runProjectCommand, stateRunning)
	runner.definition.Configure(stateRunning).OnEntry(runner.runProject)
	return runner
}

// Definition exposes the workload transitions for inspection.
func (runner *Runner) Definition() machine.Definition {
	return runner.definition
}

// Prepare stages the workload and returns its run action after successful preparation.
func (runner *Runner) Prepare(ctx context.Context, directory, id string, specification *jobpb.Workload) (func(context.Context) (Outcome, error), error) {
	current := &invocation{state: stateReady, directory: directory, id: id, specification: specification}
	if err := runner.definition.Fire(ctx, &current.state, prepareProjectCommand, current); err != nil {
		return nil, fmt.Errorf("prepare workload %s: %w", id, err)
	}
	return func(ctx context.Context) (Outcome, error) {
		current.mu.Lock()
		defer current.mu.Unlock()
		if err := runner.definition.Fire(ctx, &current.state, runProjectCommand, current); err != nil {
			return current.outcome, fmt.Errorf("run workload %s: %w", id, err)
		}
		return current.outcome, nil
	}, nil
}

// stageProject prepares the command before the run trigger is permitted.
func (runner *Runner) stageProject(ctx context.Context, _ ...any) error {
	current, err := runner.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read staging input: %w", err)
	}
	current.prepared, err = prepare(ctx, runner.filesystem, current.directory, current.id, current.specification)
	if err != nil {
		return fmt.Errorf("stage project: %w", err)
	}
	if err := runner.definition.Advance(ctx, preparedProjectCommand); err != nil {
		return fmt.Errorf("record prepared project: %w", err)
	}
	return nil
}

// runProject starts the prepared process and keeps the observed outcome.
func (runner *Runner) runProject(ctx context.Context, _ ...any) error {
	current, err := runner.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read process input: %w", err)
	}
	current.outcome, err = current.prepared.run(ctx)
	if err != nil {
		return fmt.Errorf("observe workload process: %w", err)
	}
	return nil
}
