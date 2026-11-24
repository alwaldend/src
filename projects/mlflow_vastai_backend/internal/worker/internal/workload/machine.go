package workload

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"github.com/spf13/afero"
)

const (
	// stateReady awaits a request to stage the project.
	stateReady machine.State = "workload-ready"
	// statePreparing stages the project and configures its command.
	statePreparing machine.State = "workload-preparing"
	// statePrepared holds the staged command.
	statePrepared machine.State = "workload-prepared"
	// stateRunning starts and observes the prepared command.
	stateRunning machine.State = "workload-running"
)

const (
	// prepareProjectCommand begins project staging.
	prepareProjectCommand machine.Event = "prepare"
	// preparedProjectCommand records successful project staging.
	preparedProjectCommand machine.Event = "prepared"
	// runProjectCommand runs the prepared command.
	runProjectCommand machine.Event = "run"
)

// invocation holds one workload's staging inputs and execution result.
type invocation struct {
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

// Graph describes the workload staging and execution transitions.
func (runner *Runner) Graph() machine.GraphSection {
	return machine.GraphSection{ID: "workload", Title: "Workload", DOT: runner.definition.ToGraph()}
}

// Execute stages the workload and runs the prepared command.
func (runner *Runner) Execute(ctx context.Context, directory, id string, specification *jobpb.Workload) (outcome Outcome, err error) {
	current := &invocation{state: stateReady, directory: directory, id: id, specification: specification}
	defer func() {
		if current.prepared != nil && (err != nil || outcome.ExitCode == nil || *outcome.ExitCode != 0) {
			err = errors.Join(err, discardProject(runner.filesystem, current.prepared.command.Dir))
		}
	}()
	ctx = runner.definition.WithInvocation(ctx, &current.state, current)
	if err := runner.definition.FireCtx(ctx, prepareProjectCommand); err != nil {
		return Outcome{}, fmt.Errorf("prepare workload: %w", err)
	}
	if err := runner.definition.FireCtx(ctx, runProjectCommand); err != nil {
		return current.outcome, fmt.Errorf("run workload: %w", err)
	}
	return current.outcome, nil
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
	if err := runner.definition.FireCtx(ctx, preparedProjectCommand); err != nil {
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
