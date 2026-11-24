package worker

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/sqlite"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/jobrunner"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/workload"
	"github.com/spf13/afero"
)

const (
	// listenerReady awaits the command that starts the listener.
	listenerReady machine.State = "listener-ready"
	// workerStarting acquires the listener resources.
	workerStarting machine.State = "worker-starting"
	// stateListenerServing runs the listener and its workloads.
	stateListenerServing machine.State = "worker-serving"
	// workerStopping releases listener resources after serving ends.
	workerStopping machine.State = "worker-stopping"
	// listenerStopped records a clean exit.
	listenerStopped machine.State = "listener-stopped"
	// listenerFailed records a startup, serving, or shutdown error.
	listenerFailed machine.State = "listener-failed"
)

const (
	// startListener starts a listener invocation.
	startListener machine.Event = "start-listener"
	// opened indicates that the listener's resources are ready.
	opened machine.Event = "listener-opened"
	// stop requests listener shutdown.
	stop machine.Event = "stop-listener"
	// done records a clean exit.
	done machine.Event = "done"
	// abort records an unsuccessful exit.
	abort machine.Event = "failed"
)

// listenerMachine owns service startup and shutdown.
type listenerMachine struct {
	// definition holds the listener lifecycle.
	definition *machine.Machine[struct{}]
	// fs supplies the filesystem operations.
	fs afero.Fs
	// workloads stages and runs jobs owned by the listener.
	workloads *workload.Runner
}

// newListenerMachine configures the listener independently of command dispatch.
func newListenerMachine(fs afero.Fs) *listenerMachine {
	listener := &listenerMachine{definition: machine.New[struct{}](listenerReady), fs: fs, workloads: workload.New(fs)}
	definition := listener.definition
	definition.Configure(listenerReady).Permit(startListener, workerStarting)
	definition.Configure(workerStarting).Permit(opened, stateListenerServing).Permit(abort, listenerFailed)
	definition.Configure(stateListenerServing).Permit(stop, workerStopping).Permit(abort, listenerFailed)
	definition.Configure(workerStopping).Permit(done, listenerStopped).Permit(abort, listenerFailed)
	definition.Configure(listenerStopped)
	definition.Configure(listenerFailed)
	return listener
}

// graph describes listener lifecycle and the jobs and workloads it runs.
func (listener *listenerMachine) graph() machine.GraphSection {
	return machine.GraphSection{
		ID: "listener", Title: "Worker listener", DOT: listener.definition.ToGraph(),
		Calls: []machine.GraphCall{
			{From: workerStarting, Label: "opens database", Target: sqlite.Graph()},
			{From: stateListenerServing, Label: "delegates job operations", Target: jobrunner.Graph(listener.workloads.Graph())},
		},
	}
}

// run drives the listener phases and closes resources owned by this invocation.
func (listener *listenerMachine) run(ctx context.Context, directory string) (err error) {
	state := listenerReady
	ctx = listener.definition.WithInvocation(ctx, &state, struct{}{})
	if err := listener.definition.FireCtx(ctx, startListener); err != nil {
		return fmt.Errorf("start listener: %w", err)
	}
	active, err := openListener(ctx, listener.fs, directory, listener.workloads)
	defer func() {
		if active != nil {
			err = errors.Join(err, active.close())
		}
		event := done
		if err != nil {
			event = abort
		}
		if finishErr := listener.definition.FireCtx(ctx, event); finishErr != nil {
			err = errors.Join(err, fmt.Errorf("finish listener: %w", finishErr))
		}
	}()
	if err != nil {
		return fmt.Errorf("open rental worker: %w", err)
	}
	if err := listener.definition.FireCtx(ctx, opened); err != nil {
		return fmt.Errorf("enter worker serving state: %w", err)
	}
	serveErr := active.serve(ctx)
	if err := listener.definition.FireCtx(ctx, stop); err != nil {
		return errors.Join(serveErr, fmt.Errorf("enter worker shutdown: %w", err))
	}
	return serveErr
}
