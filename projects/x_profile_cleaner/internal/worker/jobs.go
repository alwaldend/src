package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/script"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// maxRetainedJobs bounds all running and terminal job records.
	maxRetainedJobs = 16
	// maxResultBytes bounds retained successful result payloads across jobs.
	maxResultBytes = 32 * 1024 * 1024
)

// evaluationJob owns one transient execution and its immutable published status.
type evaluationJob struct {
	// status is replaced atomically under the service lock, never mutated afterward.
	status *workerpb.EvaluationStatus
	// cancel stops only this evaluation, independently of submission RPC lifetime.
	cancel context.CancelFunc
	// expires fixes terminal retention; polling never changes it.
	expires time.Time
	// timer removes this exact terminal record when retention ends.
	timer *time.Timer
	// bytes is the successful result's contribution to retained payload accounting.
	bytes int
}

// service owns bounded transient jobs; durable account/job state belongs upstream.
type service struct {
	// UnsafeArchiveWorkerServer requires every current RPC to be implemented.
	workerpb.UnsafeArchiveWorkerServer
	// evaluate is the trusted composition root for generic scripts and application helpers.
	evaluate Evaluator
	// ctx is the worker-owned lifetime, independent of any submission connection.
	ctx context.Context
	// cancel stops worker-owned computation during shutdown.
	cancel context.CancelFunc
	// active admits one archive fetch or evaluation, with no waiting queue.
	active chan struct{}
	// timeout bounds archive fetching, parsing and computation for each admitted job.
	timeout time.Duration
	// retention bounds availability after a terminal state is published.
	retention time.Duration
	// mu protects job publication, deletion, timers, and result accounting.
	mu sync.Mutex
	// storage owns the configured archive source providers.
	storage *sources
	// jobs indexes backend-provided identifiers without retaining expired history.
	jobs map[string]*evaluationJob
	// resultBytes counts only complete retained successful result payloads.
	resultBytes int
	// closed prevents job admission after shutdown starts.
	closed bool
	// workers tracks detached evaluation goroutines for bounded shutdown cleanup.
	workers sync.WaitGroup
}

// newService creates the worker's transient admission and result owner.
func newService(ctx context.Context, timeout, retention time.Duration, storage *sources, evaluate Evaluator) *service {
	ctx, cancel := context.WithCancel(ctx)
	return &service{ctx: ctx, cancel: cancel, active: make(chan struct{}, 1), timeout: timeout, retention: retention, jobs: map[string]*evaluationJob{}, storage: storage, evaluate: evaluate}
}

// SubmitEvaluation validates a reference and detaches all archive I/O from the RPC.
func (s *service) SubmitEvaluation(ctx context.Context, value *workerpb.Evaluation) (*workerpb.SubmitEvaluationResponse, error) {
	request, err := requestFromProto(value)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid evaluation metadata")
	}
	if err := s.storage.validate(request.SourceName, request.ArchiveKey); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid archive source reference")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.ctx.Err() != nil {
		return nil, status.Error(codes.Unavailable, "worker is stopping")
	}
	if err := ctx.Err(); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	for id, old := range s.jobs {
		if old.status.State != workerpb.EvaluationState_RUNNING && !time.Now().Before(old.expires) {
			s.removeLocked(id, old)
		}
	}
	if _, exists := s.jobs[request.OperationID]; exists {
		return nil, status.Error(codes.AlreadyExists, "job ID is already retained")
	}
	if len(s.jobs) >= maxRetainedJobs {
		return nil, status.Error(codes.ResourceExhausted, "worker job retention is full")
	}
	select {
	case s.active <- struct{}{}:
	default:
		return nil, status.Error(codes.ResourceExhausted, "worker is busy")
	}
	jobContext, stop := context.WithTimeout(s.ctx, s.timeout)
	job := &evaluationJob{status: &workerpb.EvaluationStatus{JobId: request.OperationID, AccountId: request.AccountID, State: workerpb.EvaluationState_RUNNING, SourceName: request.SourceName, ArchiveKey: request.ArchiveKey, EvaluatedAt: timestamppb.New(request.EvaluatedAt)}, cancel: stop}
	s.jobs[request.OperationID] = job
	s.workers.Add(1)
	go s.run(jobContext, job, request)
	return &workerpb.SubmitEvaluationResponse{JobId: request.OperationID}, nil
}

// run retains the source through evaluation and publishes only complete typed plans.
func (s *service) run(ctx context.Context, job *evaluationJob, request Request) {
	defer s.workers.Done()
	result := &workerpb.EvaluationStatus{JobId: request.OperationID, AccountId: request.AccountID, State: workerpb.EvaluationState_FAILED, SourceName: request.SourceName, ArchiveKey: request.ArchiveKey, EvaluatedAt: timestamppb.New(request.EvaluatedAt)}
	data, release, err := s.storage.load(ctx, request.SourceName, request.ArchiveKey, request.AccountID)
	if err != nil {
		result.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_INVALID_ARCHIVE, Message: "archive loading or validation failed"}
		s.finish(ctx, job, result)
		return
	}
	value, err := s.evaluate.Evaluate(ctx, data, request.Script, request.EvaluatedAt)
	closeErr := release()
	if errors.Is(err, context.DeadlineExceeded) {
		result.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_DEADLINE_EXCEEDED, Message: "script evaluation deadline exceeded"}
	} else if err != nil {
		result.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_INVALID_SCRIPT, Message: script.Diagnostic(err)}
	} else if closeErr != nil {
		result.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_INTERNAL, Message: "archive cleanup failed"}
	} else if err := plan.Validate(value); err != nil {
		result.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_INTERNAL, Message: "worker produced invalid operation plan"}
	} else {
		result.State = workerpb.EvaluationState_SUCCEEDED
		result.Plan = value
	}
	s.finish(ctx, job, result)
}

// finish publishes one immutable terminal result before releasing admission.
func (s *service) finish(ctx context.Context, job *evaluationJob, result *workerpb.EvaluationStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		job.cancel()
		job.cancel = nil
		<-s.active
		return
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		result.State = workerpb.EvaluationState_CANCELLED
		result.Plan = nil
		result.Failure = nil
	} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		result.State = workerpb.EvaluationState_FAILED
		result.Plan = nil
		result.Failure = &workerpb.EvaluationFailure{Code: workerpb.EvaluationFailure_DEADLINE_EXCEEDED, Message: "worker evaluation deadline exceeded"}
	}
	job.cancel()
	job.cancel = nil
	job.expires = time.Now().Add(s.retention)
	result.ExpiresAt = timestamppb.New(job.expires)
	size := proto.Size(result)
	if result.State == workerpb.EvaluationState_SUCCEEDED && (size > maxPlanBytes || size > maxResultBytes-s.resultBytes) {
		code := workerpb.EvaluationFailure_OUTPUT_LIMIT
		if size <= maxPlanBytes {
			code = workerpb.EvaluationFailure_RESULT_CAPACITY
		}
		result.State = workerpb.EvaluationState_FAILED
		result.Plan = nil
		result.Failure = &workerpb.EvaluationFailure{Code: code, Message: "worker result capacity exceeded"}
	}
	if result.State == workerpb.EvaluationState_SUCCEEDED {
		job.bytes = size
		s.resultBytes += size
	}
	job.status = result
	job.timer = time.AfterFunc(s.retention, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.jobs[result.JobId] == job {
			s.removeLocked(result.JobId, job)
		}
	})
	<-s.active
}

// lookupLocked verifies account binding and never returns expired data.
func (s *service) lookupLocked(ref *workerpb.EvaluationRef) (*evaluationJob, error) {
	if ref == nil || !validOperationID(ref.JobId) || !model.ValidID(ref.AccountId) {
		return nil, status.Error(codes.InvalidArgument, "invalid job reference")
	}
	job := s.jobs[ref.JobId]
	if job == nil || job.status.AccountId != ref.AccountId {
		return nil, status.Error(codes.NotFound, "evaluation not found")
	}
	if job.status.State != workerpb.EvaluationState_RUNNING && !time.Now().Before(job.expires) {
		s.removeLocked(ref.JobId, job)
		return nil, status.Error(codes.NotFound, "evaluation not found")
	}
	return job, nil
}

// GetEvaluation returns a cheap immutable snapshot without extending retention.
func (s *service) GetEvaluation(_ context.Context, ref *workerpb.EvaluationRef) (*workerpb.EvaluationStatus, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.lookupLocked(ref)
	if err != nil {
		return nil, fmt.Errorf("get evaluation: %w", err)
	}
	return job.status, nil
}

// CancelEvaluation requests cancellation once; terminal results remain unchanged.
func (s *service) CancelEvaluation(_ context.Context, ref *workerpb.EvaluationRef) (*workerpb.CancelEvaluationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.lookupLocked(ref)
	if err != nil {
		return nil, fmt.Errorf("cancel evaluation: %w", err)
	}
	if job.status.State == workerpb.EvaluationState_RUNNING {
		job.cancel()
	}
	return &workerpb.CancelEvaluationResponse{}, nil
}

// DeleteEvaluation removes only a terminal result after the backend has copied it.
func (s *service) DeleteEvaluation(_ context.Context, ref *workerpb.EvaluationRef) (*workerpb.DeleteEvaluationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, err := s.lookupLocked(ref)
	if err != nil {
		return nil, fmt.Errorf("delete evaluation: %w", err)
	}
	if job.status.State == workerpb.EvaluationState_RUNNING {
		return nil, status.Error(codes.FailedPrecondition, "cancel running evaluation before deleting it")
	}
	s.removeLocked(ref.JobId, job)
	return &workerpb.DeleteEvaluationResponse{}, nil
}

// removeLocked releases one exact record's timer and retained payload allowance.
func (s *service) removeLocked(id string, job *evaluationJob) {
	if job.timer != nil {
		job.timer.Stop()
	}
	delete(s.jobs, id)
	s.resultBytes -= job.bytes
}

// close rejects further admission, cancels computation, and drops transient data.
func (s *service) close() error {
	s.mu.Lock()
	s.closed = true
	s.cancel()
	s.mu.Unlock()
	done := make(chan struct{})
	go func() { s.workers.Wait(); close(done) }()
	var err error
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("evaluation did not stop within shutdown grace period")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, job := range s.jobs {
		if job.timer != nil {
			job.timer.Stop()
		}
	}
	s.jobs = map[string]*evaluationJob{}
	s.resultBytes = 0
	if err != nil {
		return fmt.Errorf("stop transient evaluations: %w", err)
	}
	return nil
}
