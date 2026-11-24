package jobrunner

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/jobrepository"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/workload"
	"google.golang.org/protobuf/proto"
)

const (
	// absent denotes a job with no database record.
	absent machine.State = "job-absent"
	// running reserves the job ID while the runner prepares and runs the command.
	running machine.State = "running"
	// finished records a successful command exit.
	finished machine.State = "finished"
	// failed records a failed preparation, launch, or command exit.
	failed machine.State = "failed"
	// killed records confirmed cancellation.
	killed machine.State = "killed"
	// unknown records an interrupted listener with an unconfirmed outcome.
	unknown machine.State = "unknown"
)

const (
	// submit admits a workload or checks whether it replays one already stored.
	submit machine.Event = "submit"
	// recordResult records an observed workload outcome or an unconfirmed result.
	recordResult machine.Event = "record-result"
)

// ErrConflict indicates a submission reused a job ID with different inputs.
var ErrConflict = errors.New("job ID already belongs to another workload")

// jobEvent carries the inputs to one transition.
type jobEvent struct {
	// record is the job loaded by the active transaction.
	record *jobrepository.Record
	// kind selects the event to fire.
	kind machine.Event
	// workload holds a submitted request.
	workload *jobpb.Workload
	// outcome holds the workload completion report.
	outcome workload.Outcome
	// cancelled reports that the workload context was cancelled.
	cancelled bool
	// unconfirmed reports that startup found the job without a recorded result.
	unconfirmed bool
}

// jobMachine applies durable job transitions to a transaction's record.
type jobMachine struct {
	// Machine supplies the scoped input to the named callbacks.
	*machine.Machine[jobEvent]
}

// newJobMachine configures the job transitions independently of storage.
func newJobMachine() *jobMachine {
	jobs := &jobMachine{Machine: machine.New[jobEvent](absent)}
	jobs.Configure(absent).Permit(submit, running)
	jobs.Configure(running).
		OnEntryFrom(submit, jobs.recordSubmission).
		InternalTransition(submit, jobs.recordSubmission).
		Permit(recordResult, finished, jobs.completedSuccessfully).
		Permit(recordResult, failed, jobs.completedWithFailure).
		Permit(recordResult, killed, jobs.cancellationConfirmed).
		Permit(recordResult, unknown, jobs.exitUnconfirmed)
	for _, state := range []machine.State{finished, failed, killed, unknown} {
		jobs.Configure(state).
			InternalTransition(submit, jobs.recordSubmission).
			OnEntryFrom(recordResult, jobs.recordExit).
			Ignore(recordResult)
	}
	return jobs
}

// apply advances the phase of the record before its transaction commits.
func (jobs *jobMachine) apply(ctx context.Context, record *jobrepository.Record, event jobEvent) error {
	phase := machine.State(record.Phase)
	event.record = record
	ctx = jobs.WithInvocation(ctx, &phase, event)
	if err := jobs.FireCtx(ctx, event.kind); err != nil {
		return fmt.Errorf("apply %s: %w", event.kind, err)
	}
	record.Phase = string(phase)
	return nil
}

// recordSubmission stores a new workload or validates an identical replay.
func (jobs *jobMachine) recordSubmission(ctx context.Context, _ ...any) error {
	event, err := jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("read submission: %w", err)
	}
	if machine.State(event.record.Phase) == absent {
		event.record.Workload = event.workload
		return nil
	}
	if !proto.Equal(event.record.Workload, event.workload) {
		return ErrConflict
	}
	return nil
}

// completionState classifies the observed outcome.
func (jobs *jobMachine) completionState(ctx context.Context) machine.State {
	update, err := jobs.Input(ctx)
	if err != nil {
		return ""
	}
	outcome := update.outcome
	switch {
	case update.unconfirmed || (outcome.Started && !outcome.Observed):
		return unknown
	case update.cancelled && (!outcome.Started || outcome.Stopped):
		return killed
	case outcome.ExitCode != nil && *outcome.ExitCode == 0:
		return finished
	default:
		return failed
	}
}

// completedSuccessfully matches a confirmed zero exit code.
func (jobs *jobMachine) completedSuccessfully(ctx context.Context, _ ...any) bool {
	return jobs.completionState(ctx) == finished
}

// completedWithFailure matches an unsuccessful exit or launch.
func (jobs *jobMachine) completedWithFailure(ctx context.Context, _ ...any) bool {
	return jobs.completionState(ctx) == failed
}

// cancellationConfirmed matches a workload stopped by cancellation or shutdown.
func (jobs *jobMachine) cancellationConfirmed(ctx context.Context, _ ...any) bool {
	return jobs.completionState(ctx) == killed
}

// exitUnconfirmed matches a workload whose exit could not be observed.
func (jobs *jobMachine) exitUnconfirmed(ctx context.Context, _ ...any) bool {
	return jobs.completionState(ctx) == unknown
}

// recordExit preserves the exit code in the same transaction as the terminal state.
func (jobs *jobMachine) recordExit(ctx context.Context, _ ...any) error {
	update, err := jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("record exit result: %w", err)
	}
	update.record.ExitCode = update.outcome.ExitCode
	return nil
}

// Graph describes the job lifecycle and its workload without opening SQLite.
func Graph(execution machine.GraphSection) machine.GraphSection {
	return machine.GraphSection{
		ID: "jobs", Title: "Worker jobs", DOT: newJobMachine().ToGraph(),
		Calls: []machine.GraphCall{
			{From: running, Label: "prepares and runs workload", Target: execution},
		},
	}
}
