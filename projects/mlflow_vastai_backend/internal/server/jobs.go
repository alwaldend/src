package server

import (
	"context"
	"errors"
	"fmt"
	"hash/maphash"
	"io/fs"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	vastai "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	// stateUnassigned has no reservation in the coordinator database.
	stateUnassigned machine.State = "job-unassigned"
	// stateReserved has an assignment and awaits a confirmed worker reply.
	stateReserved machine.State = "job-reserved"
	// stateActive has a confirmed running worker.
	stateActive machine.State = "job-active"
	// stateResult has a terminal or unconfirmed outcome.
	stateResult machine.State = "coordinator-result"

	// submit reserves a rental or checks an identical replay.
	submit machine.Event = "submit"
	// query refreshes the job status from its assigned worker.
	query machine.Event = "status"
	// cancel requests cancellation on the assigned worker.
	cancel machine.Event = "cancel"
	// reconcile refreshes execution or publishes a result to tracking.
	reconcile machine.Event = "reconcile"
	// observedRunning records the first confirmed running result.
	observedRunning machine.Event = "worker-running"
	// observedResult records a terminal or unconfirmed result.
	observedResult machine.Event = "worker-result"
)

var (
	// errInvalidJob indicates a request that cannot be admitted.
	errInvalidJob = errors.New("invalid job ID or workload")
	// errJobMissing indicates a run without a coordinator reservation.
	errJobMissing = errors.New("job not found")
	// errJobConflict indicates incompatible inputs for an existing run ID.
	errJobConflict = errors.New("run ID already belongs to a different submission")
	// errWorkerUnavailable indicates that the worker did not confirm the operation.
	errWorkerUnavailable = errors.New("assigned worker did not confirm the operation")
	// runID limits the ID to a valid worker path component.
	runID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)
)

// operation carries the state and request used by one coordinator transition.
type operation struct {
	// state is the coordinator phase of the loaded job.
	state machine.State
	// row holds the durable assignment when one exists.
	row *assignment
	// request is present on submissions.
	request *jobpb.SubmitJobRequest
	// response reports the confirmed observation and any communication error.
	response *jobpb.Job
}

// coordinator serves job RPCs and reconciles the same assignments after restart.
type coordinator struct {
	// UnimplementedJobServiceServer provides the generated RPC defaults.
	jobpb.UnimplementedJobServiceServer
	// jobs handles coordinator transitions.
	jobs *machine.Machine[*operation]
	// repository preserves reservations and results.
	repository *repository
	// rentals discovers and contacts the assigned instance.
	rentals rentals
	// tracking reports terminal results to MLflow.
	tracking tracking
	// seed distributes job IDs across locks.
	seed maphash.Seed
	// locks serialize changes to the same assignment while permitting parallel RPCs.
	locks [64]sync.Mutex
}

// newCoordinator creates the service and its state machine.
func newCoordinator(repository *repository, rentals rentals, tracking tracking) *coordinator {
	service := &coordinator{repository: repository, rentals: rentals, tracking: tracking, seed: maphash.MakeSeed()}
	service.jobs = machine.New[*operation](stateUnassigned)
	service.jobs.Configure(stateUnassigned).Permit(submit, stateReserved)
	service.jobs.Configure(stateReserved).OnEntryFrom(submit, service.submit).
		Permit(observedRunning, stateActive).Permit(observedResult, stateResult)
	service.jobs.Configure(stateActive).OnEntryFrom(observedRunning, service.saveObservation).Permit(observedResult, stateResult)
	service.jobs.Configure(stateResult).OnEntryFrom(observedResult, service.saveObservation)
	for _, state := range []machine.State{stateReserved, stateActive, stateResult} {
		service.jobs.Configure(state).InternalTransition(submit, service.submit).
			InternalTransition(query, service.query).InternalTransition(cancel, service.cancel).
			InternalTransition(reconcile, service.reconcileJob)
	}
	return service
}

// jobGraph describes coordinator job transitions without accessing live services.
func jobGraph() machine.GraphSection {
	return machine.GraphSection{ID: "coordinator_jobs", Title: "Coordinator jobs", DOT: newCoordinator(nil, rentals{}, tracking{}).jobs.ToGraph()}
}

// SubmitJob validates and reserves an execution or returns its existing state.
func (service *coordinator) SubmitJob(ctx context.Context, request *jobpb.SubmitJobRequest) (*jobpb.Job, error) {
	return response(service.run(ctx, request.GetJobId(), submit, request))
}

// GetJob observes a previously submitted job.
func (service *coordinator) GetJob(ctx context.Context, reference *jobpb.JobReference) (*jobpb.Job, error) {
	return response(service.run(ctx, reference.GetJobId(), query, nil))
}

// CancelJob records cancellation intent and contacts the assigned worker.
func (service *coordinator) CancelJob(ctx context.Context, reference *jobpb.JobReference) (*jobpb.Job, error) {
	return response(service.run(ctx, reference.GetJobId(), cancel, nil))
}

// response translates domain errors without exposing SSH or API diagnostics.
func response(job *jobpb.Job, err error) (*jobpb.Job, error) {
	switch {
	case err == nil:
		return job, nil
	case errors.Is(err, errInvalidJob):
		return nil, status.Error(codes.InvalidArgument, errInvalidJob.Error())
	case errors.Is(err, errJobMissing):
		return nil, status.Error(codes.NotFound, errJobMissing.Error())
	case errors.Is(err, errJobConflict):
		return nil, status.Error(codes.AlreadyExists, errJobConflict.Error())
	case errors.Is(err, errNoRental):
		return nil, status.Error(codes.FailedPrecondition, errNoRental.Error())
	case errors.Is(err, errWorkerUnavailable):
		return nil, status.Error(codes.Unavailable, errWorkerUnavailable.Error())
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, status.FromContextError(err).Err()
	default:
		log.Printf("job coordinator: %v", err)
		return nil, status.Error(codes.Internal, "coordinator operation failed")
	}
}

// run binds the persisted assignment and dispatches one machine event.
func (service *coordinator) run(ctx context.Context, id string, event machine.Event, request *jobpb.SubmitJobRequest) (*jobpb.Job, error) {
	if !runID.MatchString(id) {
		return nil, errInvalidJob
	}
	ctx, stop := context.WithTimeout(ctx, 15*time.Second)
	defer stop()
	lock := &service.locks[maphash.String(service.seed, id)%uint64(len(service.locks))]
	lock.Lock()
	defer lock.Unlock()
	row, err := service.repository.find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load job %s: %w", id, err)
	}
	current := &operation{row: row, state: stateUnassigned, request: request}
	if row != nil {
		current.state = phase(row.job.GetState())
		current.response = row.job
	} else if event != submit {
		return nil, errJobMissing
	}
	ctx = service.jobs.WithInvocation(ctx, &current.state, current)
	if err := service.jobs.FireCtx(ctx, event); err != nil {
		return nil, fmt.Errorf("apply %s to job %s: %w", event, id, err)
	}
	return current.response, nil
}

// phase maps the recorded worker state to the coordinator lifecycle.
func phase(state jobpb.JobState) machine.State {
	switch state {
	case jobpb.JobState_JOB_STATE_SCHEDULED:
		return stateReserved
	case jobpb.JobState_JOB_STATE_RUNNING:
		return stateActive
	default:
		return stateResult
	}
}

// submit reserves the assignment before its first dispatch and validates replay.
func (service *coordinator) submit(ctx context.Context, _ ...any) error {
	current, err := service.jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("read submission action: %w", err)
	}
	request := current.request
	if current.row == nil {
		workload := request.GetWorkload()
		path := workload.GetProjectArtifactPath()
		if path == "." || !fs.ValidPath(path) || strings.ContainsRune(path, 0) ||
			strings.TrimSpace(workload.GetCommand()) == "" || strings.TrimSpace(workload.GetExperimentId()) == "" ||
			strings.TrimRight(workload.GetTrackingUri(), "/") != service.tracking.endpoint {
			return errInvalidJob
		}
		instance, err := service.rentals.selectRental(ctx, request.GetWorker())
		if err != nil {
			return fmt.Errorf("select rental: %w", err)
		}
		current.row = &assignment{
			submission: request, instance: instance.ID,
			job: &jobpb.Job{JobId: request.GetJobId(), State: jobpb.JobState_JOB_STATE_SCHEDULED},
		}
		if err := service.repository.reserve(ctx, current.row); err != nil {
			return fmt.Errorf("persist assignment: %w", err)
		}
		if err := service.observe(ctx, current, instance); err != nil {
			return fmt.Errorf("dispatch reserved job: %w", err)
		}
	} else {
		if !proto.Equal(current.row.submission, request) {
			return errJobConflict
		}
		if err := service.refresh(ctx, current); err != nil {
			return fmt.Errorf("refresh submitted job: %w", err)
		}
	}
	if current.response.GetObservationError() != "" {
		return errWorkerUnavailable
	}
	return nil
}

// query reads the latest result, refreshing the worker while the job is active.
func (service *coordinator) query(ctx context.Context, _ ...any) error {
	current, err := service.jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("read query action: %w", err)
	}
	if err := service.refresh(ctx, current); err != nil {
		return fmt.Errorf("refresh queried job: %w", err)
	}
	return nil
}

// cancel records intent before issuing a potentially ambiguous SSH command.
func (service *coordinator) cancel(ctx context.Context, _ ...any) error {
	current, err := service.jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("read cancellation action: %w", err)
	}
	if current.state == stateResult {
		return nil
	}
	if !current.row.cancelling {
		current.row.cancelling = true
		if err := service.repository.save(ctx, current.row); err != nil {
			return fmt.Errorf("persist cancellation intent: %w", err)
		}
	}
	if err := service.refresh(ctx, current); err != nil {
		return fmt.Errorf("refresh cancelled job: %w", err)
	}
	return nil
}

// refresh contacts only the assigned rental and preserves prior state on temporary failure.
func (service *coordinator) refresh(ctx context.Context, current *operation) error {
	if current.state == stateResult {
		return nil
	}
	instance, err := service.rentals.resolve(ctx, current.row.instance)
	if errors.Is(err, errRentalLost) {
		current.row.job = &jobpb.Job{
			JobId: current.row.job.GetJobId(), State: jobpb.JobState_JOB_STATE_UNKNOWN,
			ObservationError: "assigned rental no longer exists",
		}
		return service.fireObservation(ctx, observedResult)
	}
	if err != nil {
		service.unavailable(current)
		return nil
	}
	if err := service.observe(ctx, current, instance); err != nil {
		return fmt.Errorf("observe assigned worker: %w", err)
	}
	return nil
}

// observe performs the appropriate operation on the assigned worker.
func (service *coordinator) observe(ctx context.Context, current *operation, instance *vastai.Instance) error {
	action := "status"
	var request proto.Message = &jobpb.JobReference{JobId: current.row.job.GetJobId()}
	if current.row.cancelling {
		action = "cancel"
	} else if current.state == stateReserved {
		action, request = "submit", current.row.submission
	}
	job, err := service.rentals.call(ctx, instance, action, request)
	if err != nil || job.GetJobId() != current.row.job.GetJobId() {
		service.unavailable(current)
		return nil
	}
	switch job.GetState() {
	case jobpb.JobState_JOB_STATE_RUNNING:
		if job.GetObservationError() != "" {
			service.unavailable(current)
			return nil
		}
		current.response = job
		if current.state == stateReserved {
			current.row.job = job
			return service.fireObservation(ctx, observedRunning)
		}
		return nil
	case jobpb.JobState_JOB_STATE_FINISHED, jobpb.JobState_JOB_STATE_FAILED,
		jobpb.JobState_JOB_STATE_KILLED, jobpb.JobState_JOB_STATE_UNKNOWN:
		current.row.job = job
		return service.fireObservation(ctx, observedResult)
	default:
		service.unavailable(current)
		return nil
	}
}

// fireObservation advances and commits an authoritative worker result.
func (service *coordinator) fireObservation(ctx context.Context, event machine.Event) error {
	if err := service.jobs.FireCtx(ctx, event); err != nil {
		return fmt.Errorf("record confirmed observation: %w", err)
	}
	return nil
}

// saveObservation writes a result before exposing it to callers.
func (service *coordinator) saveObservation(ctx context.Context, _ ...any) error {
	current, err := service.jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("read observation action: %w", err)
	}
	current.row.settled = current.row.job.GetState() == jobpb.JobState_JOB_STATE_UNKNOWN
	if err := service.repository.save(ctx, current.row); err != nil {
		return fmt.Errorf("persist observed job: %w", err)
	}
	current.response = current.row.job
	return nil
}

// unavailable reports an observation failure without replacing the confirmed result.
func (service *coordinator) unavailable(current *operation) {
	current.response = proto.CloneOf(current.row.job)
	current.response.ObservationError = errWorkerUnavailable.Error()
}

// reconcileJob refreshes jobs without a client and retries unreported terminal results.
func (service *coordinator) reconcileJob(ctx context.Context, _ ...any) error {
	current, err := service.jobs.Input(ctx)
	if err != nil {
		return fmt.Errorf("read reconciliation action: %w", err)
	}
	if err := service.refresh(ctx, current); err != nil {
		return fmt.Errorf("refresh pending job: %w", err)
	}
	if current.state != stateResult || current.row.settled {
		return nil
	}
	if err := service.tracking.publish(ctx, current.row.job); err != nil {
		return fmt.Errorf("publish terminal status: %w", err)
	}
	current.row.settled = true
	if err := service.repository.save(ctx, current.row); err != nil {
		return fmt.Errorf("commit tracking publication: %w", err)
	}
	return nil
}

// reconcile observes pending jobs until the server stops.
func (service *coordinator) reconcile(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for ctx.Err() == nil {
		ids, err := service.repository.pending(ctx)
		if err != nil && ctx.Err() == nil {
			log.Printf("list pending jobs: %v", err)
		}
		for _, id := range ids {
			if ctx.Err() != nil {
				return
			}
			if _, err := service.run(ctx, id, reconcile, nil); err != nil && ctx.Err() == nil {
				log.Printf("reconcile pending job: %v", err)
			}
		}
		select {
		case <-ctx.Done():
		case <-ticker.C:
		}
	}
}
