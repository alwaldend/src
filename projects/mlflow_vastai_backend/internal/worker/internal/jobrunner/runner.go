package jobrunner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/jobrepository"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/workload"
)

// ErrNotFound indicates that the worker has no record for a job ID.
var ErrNotFound = jobrepository.ErrNotFound

// Executor prepares and runs one workload in its job directory.
type Executor interface {
	Execute(ctx context.Context, directory, id string, specification *jobpb.Workload) (workload.Outcome, error)
}

// Runner owns workload execution and commits job state transitions.
type Runner struct {
	// repository stores the durable job records.
	repository *jobrepository.Repository
	// jobs advances each record before it is written.
	jobs *jobMachine
	// executor prepares and runs the workload command.
	executor Executor
	// directory holds the per-job files.
	directory string
	// ctx sets the lifetime of accepted workloads.
	ctx context.Context
	// stop cancels the runner and its workloads.
	stop context.CancelFunc
	// mu serializes admission, cancellation, and shutdown.
	mu sync.Mutex
	// active holds the cancellation function for each registered workload.
	active map[string]context.CancelFunc
	// running joins executions and result writes on shutdown.
	running sync.WaitGroup
}

// Open reconciles unfinished jobs before accepting new submissions.
// The caller owns the repository's connection and closes it after Close returns.
func Open(ctx context.Context, repository *jobrepository.Repository, directory string, executor Executor) (*Runner, error) {
	runner := &Runner{
		repository: repository, jobs: newJobMachine(), executor: executor, directory: directory,
		active: make(map[string]context.CancelFunc),
	}
	if err := runner.recover(ctx); err != nil {
		return nil, fmt.Errorf("reconcile interrupted jobs: %w", err)
	}
	runner.ctx, runner.stop = context.WithCancel(ctx)
	return runner, nil
}

// Close stops admission, cancels active workloads, and waits for their results.
func (runner *Runner) Close() {
	runner.mu.Lock()
	runner.stop()
	runner.mu.Unlock()
	runner.running.Wait()
}

// update commits a job event and returns the record and its prior phase.
func (runner *Runner) update(ctx context.Context, id string, event jobEvent) (*jobrepository.Record, machine.State, error) {
	var current *jobrepository.Record
	var prior machine.State
	err := runner.repository.Update(ctx, id, func(record *jobrepository.Record, exists bool) (bool, error) {
		if !exists {
			if event.kind != submit {
				return false, ErrNotFound
			}
			record.Phase = string(absent)
		}
		current = record
		prior = machine.State(record.Phase)
		if err := runner.jobs.apply(ctx, record, event); err != nil {
			return false, fmt.Errorf("apply job transition: %w", err)
		}
		return machine.State(record.Phase) != prior, nil
	})
	if err != nil {
		return nil, "", fmt.Errorf("update job %s: %w", id, err)
	}
	return current, prior, nil
}

// recover marks jobs left running by a prior listener as unknown.
func (runner *Runner) recover(ctx context.Context) error {
	ids, err := runner.repository.IDsInPhase(ctx, string(running))
	if err != nil {
		return fmt.Errorf("list interrupted jobs: %w", err)
	}
	for _, id := range ids {
		if _, _, err := runner.update(ctx, id, jobEvent{kind: recordResult, unconfirmed: true}); err != nil {
			return fmt.Errorf("record interrupted job %s: %w", id, err)
		}
	}
	return nil
}

// Submit reserves and starts a job or returns the existing observation for an identical request.
// The caller validates the run ID as a filesystem component.
func (runner *Runner) Submit(ctx context.Context, id string, specification *jobpb.Workload) (*jobpb.Job, error) {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	if err := runner.ctx.Err(); err != nil {
		return nil, fmt.Errorf("submit job: %w", err)
	}
	record, prior, err := runner.update(ctx, id, jobEvent{kind: submit, workload: specification})
	if err == nil && prior == absent {
		runner.start(id, specification)
	}
	return observation(record, err)
}

// Get returns the latest committed job observation.
func (runner *Runner) Get(ctx context.Context, id string) (*jobpb.Job, error) {
	return observation(runner.repository.Get(ctx, id))
}

// Cancel cancels an active workload and returns its latest committed observation.
func (runner *Runner) Cancel(ctx context.Context, id string) (*jobpb.Job, error) {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	if stop := runner.active[id]; stop != nil {
		stop()
	}
	return observation(runner.repository.Get(ctx, id))
}

// start registers a workload while the caller holds mu.
func (runner *Runner) start(id string, specification *jobpb.Workload) {
	ctx, stop := context.WithCancel(runner.ctx)
	runner.active[id] = stop
	runner.running.Go(func() {
		defer func() {
			stop()
			runner.mu.Lock()
			delete(runner.active, id)
			runner.mu.Unlock()
		}()
		outcome, err := runner.executor.Execute(ctx, filepath.Join(runner.directory, id), id, specification)
		if err != nil {
			log.Printf("job %s: %v", id, err)
		}
		if err := runner.commitResult(id, jobEvent{kind: recordResult, outcome: outcome, cancelled: ctx.Err() != nil}); err != nil {
			log.Printf("commit job %s result: %v", id, err)
		}
	})
}

// commitResult retains the observed outcome across retries and bounds the shutdown flush.
func (runner *Runner) commitResult(id string, event jobEvent) error {
	ctx := runner.ctx
	flushing := false
	for {
		if runner.ctx.Err() != nil && !flushing {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.WithoutCancel(runner.ctx), 10*time.Second)
			defer cancel()
			flushing = true
		}
		_, _, err := runner.update(ctx, id, event)
		if err == nil {
			return nil
		}
		if flushing && ctx.Err() != nil {
			return fmt.Errorf("flush job result: %w", err)
		}
		log.Printf("commit job %s result, will retry: %v", id, err)
		select {
		case <-ctx.Done():
		case <-time.After(time.Second):
		}
	}
}

// observation projects the persisted phase to the worker protocol.
func observation(record *jobrepository.Record, err error) (*jobpb.Job, error) {
	if err != nil {
		return nil, fmt.Errorf("read job observation: %w", err)
	}
	response := &jobpb.Job{JobId: record.ID, ExitCode: record.ExitCode}
	switch machine.State(record.Phase) {
	case running:
		response.State = jobpb.JobState_JOB_STATE_RUNNING
	case finished:
		response.State = jobpb.JobState_JOB_STATE_FINISHED
	case failed:
		response.State = jobpb.JobState_JOB_STATE_FAILED
	case killed:
		response.State = jobpb.JobState_JOB_STATE_KILLED
	case unknown:
		response.State = jobpb.JobState_JOB_STATE_UNKNOWN
		response.ObservationError = "worker could not confirm the workload outcome"
	default:
		return nil, errors.New("worker state is invalid")
	}
	return response, nil
}
