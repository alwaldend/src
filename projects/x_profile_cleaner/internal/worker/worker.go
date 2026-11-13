package worker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// maxScriptBytes bounds UTF-8 Lua source before parsing.
	maxScriptBytes = 64 * 1024
	// maxPlanBytes bounds complete RPC output, including trusted metadata.
	maxPlanBytes = 16 * 1024 * 1024
	// maxRequestBytes bounds unary metadata and Lua source, without archive data.
	maxRequestBytes = 128 * 1024
	// maxLifetime bounds every evaluation even without a client deadline.
	maxLifetime = 90 * time.Second
)

// Request identifies a configured archive object and script, without X authorization.
type Request struct {
	// OperationID binds the response to the backend's operation.
	OperationID string
	// AccountID is the expected archive owner selected by the backend.
	AccountID string
	// EvaluatedAt fixes age predicates independently of worker wall time.
	EvaluatedAt time.Time
	// Script is UTF-8 Lua source collecting typed operations without returning values.
	Script string
	// SourceName selects a provider configured by the worker operator.
	SourceName string
	// ArchiveKey is an opaque object key interpreted by the selected provider.
	ArchiveKey string
}

// Evaluator computes a typed plan from one borrowed export under the job deadline.
// The composition root may install trusted application libraries; X access is forbidden.
type Evaluator interface {
	// Evaluate consumes one borrowed export and returns a complete typed plan.
	Evaluate(context.Context, *archive.Export, string, time.Time) (*workerpb.Plan, error)
}

// Serve listens on a new Unix socket in an existing private directory. Existing
// socket paths are never removed or replaced. Cancellation stops active RPCs.
func Serve(ctx context.Context, socket string, timeout, retention time.Duration, configured map[string]string, evaluate Evaluator) (resultErr error) {
	if evaluate == nil || !filepath.IsAbs(socket) || timeout <= 0 || timeout > maxLifetime || retention <= 0 || retention > time.Hour {
		return fmt.Errorf("worker requires an absolute socket, timeout in (0, 90s], and result TTL in (0, 1h]")
	}
	info, err := os.Stat(filepath.Dir(socket))
	if err != nil {
		return fmt.Errorf("inspect worker socket directory: %w", err)
	}
	if !info.IsDir() || info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("worker socket directory must be private (0700)")
	}
	storage, err := openSources(configured)
	if err != nil {
		return fmt.Errorf("configure worker sources: %w", err)
	}
	defer func() {
		if err := storage.close(); err != nil {
			resultErr = errors.Join(resultErr, err)
		}
	}()
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: socket, Net: "unix"})
	if err != nil {
		return fmt.Errorf("listen on worker socket: %w", err)
	}
	defer listener.Close()
	if err := os.Chmod(socket, 0o600); err != nil {
		return fmt.Errorf("restrict worker socket: %w", err)
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(maxRequestBytes), grpc.MaxSendMsgSize(maxPlanBytes), grpc.MaxConcurrentStreams(2))
	jobs := newService(ctx, timeout, retention, storage, evaluate)
	defer func() {
		if err := jobs.close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("stop worker jobs: %w", err))
		}
	}()
	defer server.Stop()
	workerpb.RegisterArchiveWorkerServer(server, jobs)
	finished := make(chan struct{})
	defer close(finished)
	go func() {
		select {
		case <-ctx.Done():
			server.Stop()
		case <-finished:
		}
	}()
	if err := server.Serve(listener); err != nil && ctx.Err() == nil {
		return fmt.Errorf("serve worker RPCs: %w", err)
	}
	return nil
}

// validate checks metadata without reading data or accepting caller-selected roots.
func (r Request) validate() error {
	if !validOperationID(r.OperationID) {
		return fmt.Errorf("invalid operation ID")
	}
	if !model.ValidID(r.AccountID) || r.EvaluatedAt.IsZero() || r.EvaluatedAt.Year() < 1 || r.EvaluatedAt.Year() > 9999 {
		return fmt.Errorf("invalid account or evaluation time")
	}
	if len(r.Script) > maxScriptBytes || !utf8.ValidString(r.Script) || !validSourceName(r.SourceName) || !validArchiveKey(r.ArchiveKey) {
		return fmt.Errorf("invalid script or archive reference")
	}
	return nil
}

// requestFromProto validates a unary job request without accessing its archive.
func requestFromProto(value *workerpb.Evaluation) (Request, error) {
	if value == nil {
		return Request{}, fmt.Errorf("missing evaluation")
	}
	stamp, err := wireTime(value.EvaluatedAt)
	if err != nil {
		return Request{}, fmt.Errorf("decode evaluation time: %w", err)
	}
	result := Request{OperationID: value.OperationId, AccountID: value.AccountId, EvaluatedAt: stamp, Script: value.Script, SourceName: value.SourceName, ArchiveKey: value.ArchiveKey}
	if err := result.validate(); err != nil {
		return Request{}, fmt.Errorf("validate evaluation: %w", err)
	}
	return result, nil
}

// Client is a backend-owned gRPC connection to a separately supervised worker.
type Client struct {
	// connection owns the transport lifetime; submissions are never retried.
	connection *grpc.ClientConn
}

// Dial creates a lazy private-socket connection without reading credentials.
// The caller owns a trusted private socket directory and the worker lifetime.
func Dial(socket string) (*Client, error) {
	if !filepath.IsAbs(socket) {
		return nil, fmt.Errorf("worker socket path must be absolute")
	}
	connection, err := grpc.NewClient("unix://"+socket, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxRequestBytes), grpc.MaxCallRecvMsgSize(maxPlanBytes)))
	if err != nil {
		return nil, fmt.Errorf("connect archive worker: %w", err)
	}
	return &Client{connection: connection}, nil
}

// Close releases the backend's worker transport.
func (c *Client) Close() error {
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("close worker connection: %w", err)
	}
	return nil
}

// Submit sends a small archive reference and returns before fetching or evaluation.
func (c *Client) Submit(ctx context.Context, request Request) (string, error) {
	if err := request.validate(); err != nil {
		return "", fmt.Errorf("validate worker request: %w", err)
	}
	response, err := workerpb.NewArchiveWorkerClient(c.connection).SubmitEvaluation(ctx, &workerpb.Evaluation{OperationId: request.OperationID, AccountId: request.AccountID, EvaluatedAt: timestamppb.New(request.EvaluatedAt), Script: request.Script, SourceName: request.SourceName, ArchiveKey: request.ArchiveKey})
	if err != nil {
		return "", fmt.Errorf("submit worker evaluation: %w", err)
	}
	if response.JobId != request.OperationID {
		return "", fmt.Errorf("worker submission binding mismatch")
	}
	return response.JobId, nil
}

// Job is the validated transient status of one backend-owned operation.
type Job struct {
	// ID is the backend-provided operation identifier returned by submission.
	ID string
	// AccountID binds every lookup and result to the snapshot owner.
	AccountID string
	// State is the current execution state from the worker contract.
	State workerpb.EvaluationState
	// Plan is a validated operation plan; applications still authorize its execution.
	Plan *workerpb.Plan
	// Failure is present only on a terminal failed computation.
	Failure *workerpb.EvaluationFailure
	// ExpiresAt bounds terminal result availability; running jobs have no expiry.
	ExpiresAt time.Time
}

// Get checks trusted bindings and typed operation structure, without interpreting operations.
// Applications must validate their own operation schemas before publication or execution.
func (c *Client) Get(ctx context.Context, request Request) (Job, error) {
	if err := request.validate(); err != nil {
		return Job{}, fmt.Errorf("validate result context: %w", err)
	}
	response, err := workerpb.NewArchiveWorkerClient(c.connection).GetEvaluation(ctx, &workerpb.EvaluationRef{JobId: request.OperationID, AccountId: request.AccountID})
	if err != nil {
		return Job{}, fmt.Errorf("get worker evaluation: %w", err)
	}
	if response.JobId != request.OperationID || response.AccountId != request.AccountID || response.SourceName != request.SourceName || response.ArchiveKey != request.ArchiveKey || proto.Size(response) > maxPlanBytes {
		return Job{}, fmt.Errorf("invalid worker status binding or size")
	}
	stamp, err := wireTime(response.EvaluatedAt)
	if err != nil {
		return Job{}, fmt.Errorf("decode worker evaluation time: %w", err)
	}
	if !stamp.Equal(request.EvaluatedAt) {
		return Job{}, fmt.Errorf("worker evaluation time binding mismatch")
	}
	job := Job{ID: response.JobId, AccountID: response.AccountId, State: response.State, Failure: response.Failure}
	if response.State == workerpb.EvaluationState_RUNNING {
		if response.Plan != nil || response.Failure != nil || response.ExpiresAt != nil {
			return Job{}, fmt.Errorf("running status contains terminal result")
		}
		return job, nil
	}
	expires, err := wireTime(response.ExpiresAt)
	if err != nil {
		return Job{}, fmt.Errorf("decode result expiry: %w", err)
	}
	if expires.IsZero() {
		return Job{}, fmt.Errorf("invalid terminal result expiry")
	}
	job.ExpiresAt = expires
	switch response.State {
	case workerpb.EvaluationState_SUCCEEDED:
		if response.Failure != nil {
			return Job{}, fmt.Errorf("successful status contains failure")
		}
		if err := plan.Validate(response.Plan); err != nil {
			return Job{}, fmt.Errorf("reject worker result: %w", err)
		}
		job.Plan = proto.Clone(response.Plan).(*workerpb.Plan)
	case workerpb.EvaluationState_FAILED:
		if response.Plan != nil || response.Failure == nil || response.Failure.Code < workerpb.EvaluationFailure_INVALID_SCRIPT || response.Failure.Code > workerpb.EvaluationFailure_INVALID_ARCHIVE || len(response.Failure.Message) > 256 {
			return Job{}, fmt.Errorf("invalid failed job status")
		}
	case workerpb.EvaluationState_CANCELLED:
		if response.Plan != nil || response.Failure != nil {
			return Job{}, fmt.Errorf("cancelled status contains a result")
		}
	default:
		return Job{}, fmt.Errorf("invalid worker job state")
	}
	return job, nil
}

// Cancel requests cancellation without changing an already terminal result.
func (c *Client) Cancel(ctx context.Context, account, id string) error {
	if _, err := workerpb.NewArchiveWorkerClient(c.connection).CancelEvaluation(ctx, &workerpb.EvaluationRef{JobId: id, AccountId: account}); err != nil {
		return fmt.Errorf("cancel worker evaluation: %w", err)
	}
	return nil
}

// Delete releases a terminal result after the backend has copied it.
func (c *Client) Delete(ctx context.Context, account, id string) error {
	if _, err := workerpb.NewArchiveWorkerClient(c.connection).DeleteEvaluation(ctx, &workerpb.EvaluationRef{JobId: id, AccountId: account}); err != nil {
		return fmt.Errorf("delete worker evaluation: %w", err)
	}
	return nil
}

// validOperationID accepts bounded printable backend-owned job identifiers.
func validOperationID(id string) bool {
	return len(id) > 0 && len(id) <= 128 && strings.IndexFunc(id, func(c rune) bool { return c < '!' || c > '~' }) < 0
}
