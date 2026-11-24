package worker

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"path/filepath"
	"regexp"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/sqlite"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/jobrepository"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/jobrunner"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/worker/internal/workload"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// validID permits an MLflow run ID to be used as one filesystem component.
var validID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

// service adapts the protobuf endpoint to job runner operations.
type service struct {
	// UnimplementedWorkerServiceServer supplies the generated RPC defaults.
	jobpb.UnimplementedWorkerServiceServer
	// jobs owns workload execution and results.
	jobs *jobrunner.Runner
}

// listener holds the resources of one rental-local listener invocation.
type listener struct {
	// lock excludes another listener from the same state directory.
	lock afero.File
	// socket accepts local worker RPCs.
	socket net.Listener
	// server owns gRPC handler lifetimes.
	server *grpc.Server
	// db owns the connection borrowed by the job repository.
	db *sqlite.DB
	// jobs borrows the database until its workloads have finished.
	jobs *jobrunner.Runner
}

// openListener prepares the socket, database, and gRPC handlers for a worker invocation.
func openListener(ctx context.Context, filesystem afero.Fs, directory string, workloads *workload.Runner) (_ *listener, err error) {
	if err := createStateDirectory(filesystem, directory); err != nil {
		return nil, fmt.Errorf("create worker directory: %w", err)
	}
	lock, err := lockFile(filesystem, filepath.Join(directory, "worker.lock"))
	if err != nil {
		return nil, fmt.Errorf("acquire worker listener lock: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, closeResource(lock, "worker listener lock"))
		}
	}()
	db, err := sqlite.Open(ctx, filepath.Join(directory, "jobs.sqlite3"))
	if err != nil {
		return nil, fmt.Errorf("open worker database: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, closeResource(db, "worker database"))
		}
	}()
	repository := jobrepository.New(db.Connection())
	if err := repository.EnsureSchema(ctx); err != nil {
		return nil, fmt.Errorf("prepare worker job repository: %w", err)
	}
	jobs, err := jobrunner.Open(ctx, repository, filepath.Join(directory, "jobs"), workloads)
	if err != nil {
		return nil, fmt.Errorf("open job runner: %w", err)
	}
	socket, err := listen(filesystem, socketPath(directory))
	if err != nil {
		jobs.Close()
		return nil, fmt.Errorf("listen for worker RPCs: %w", err)
	}
	server := grpc.NewServer(grpc.WaitForHandlers(true))
	jobpb.RegisterWorkerServiceServer(server, &service{jobs: jobs})
	return &listener{lock: lock, socket: socket, server: server, db: db, jobs: jobs}, nil
}

// serve handles RPCs and joins their handlers before returning.
func (listener *listener) serve(ctx context.Context) error {
	serveCtx, stop := context.WithCancel(ctx)
	stopped := make(chan struct{})
	context.AfterFunc(serveCtx, func() {
		defer close(stopped)
		listener.server.Stop()
	})
	defer func() {
		stop()
		<-stopped
	}()
	if err := listener.server.Serve(listener.socket); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve worker gRPC: %w", err)
	}
	return nil
}

// close stops workloads before releasing their database and lock.
func (listener *listener) close() error {
	listener.jobs.Close()
	err := closeResource(listener.socket, "worker socket")
	if errors.Is(err, net.ErrClosed) {
		err = nil
	}
	return errors.Join(err, closeResource(listener.db, "worker database"), closeResource(listener.lock, "worker listener lock"))
}

// listen replaces a stale socket only while the caller owns the listener lock.
func listen(filesystem afero.Fs, path string) (net.Listener, error) {
	info, err := filesystem.Stat(path)
	if err == nil {
		if info.Mode()&fs.ModeSocket == 0 {
			return nil, fmt.Errorf("worker endpoint %q is not a socket", path)
		}
		if err := filesystem.Remove(path); err != nil {
			return nil, fmt.Errorf("remove stale worker socket: %w", err)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("inspect worker socket: %w", err)
	}
	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("bind worker socket: %w", err)
	}
	if err := filesystem.Chmod(path, 0o600); err != nil {
		return nil, errors.Join(fmt.Errorf("restrict worker socket: %w", err), closeResource(listener, "worker socket"))
	}
	return listener, nil
}

// SubmitJob starts a new workload or returns the existing observation for a replay.
func (s *service) SubmitJob(ctx context.Context, request *jobpb.SubmitJobRequest) (*jobpb.Job, error) {
	id, workload := request.GetJobId(), request.GetWorkload()
	artifact := workload.GetProjectArtifactPath()
	if !validID.MatchString(id) || strings.TrimSpace(workload.GetCommand()) == "" ||
		strings.TrimSpace(workload.GetTrackingUri()) == "" || strings.TrimSpace(workload.GetExperimentId()) == "" ||
		artifact == "." || !fs.ValidPath(artifact) || strings.ContainsRune(artifact, 0) {
		return nil, status.Error(codes.InvalidArgument, "job ID and complete workload with a run-relative artifact path are required")
	}
	record, err := s.jobs.Submit(ctx, id, workload)
	return jobResponse(record, err)
}

// GetJob reads the latest committed observation.
func (s *service) GetJob(ctx context.Context, reference *jobpb.JobReference) (*jobpb.Job, error) {
	id := reference.GetJobId()
	if !validID.MatchString(id) {
		return nil, status.Error(codes.InvalidArgument, "invalid job ID")
	}
	record, err := s.jobs.Get(ctx, id)
	return jobResponse(record, err)
}

// CancelJob requests an active workload's cancellation and returns its current state.
func (s *service) CancelJob(ctx context.Context, reference *jobpb.JobReference) (*jobpb.Job, error) {
	id := reference.GetJobId()
	if !validID.MatchString(id) {
		return nil, status.Error(codes.InvalidArgument, "invalid job ID")
	}
	record, err := s.jobs.Cancel(ctx, id)
	return jobResponse(record, err)
}

// jobResponse translates persisted state and domain errors into the wire contract.
func jobResponse(response *jobpb.Job, err error) (*jobpb.Job, error) {
	switch {
	case errors.Is(err, jobrunner.ErrNotFound):
		return nil, status.Error(codes.NotFound, "job not found")
	case errors.Is(err, jobrunner.ErrConflict):
		return nil, status.Error(codes.AlreadyExists, "job ID already belongs to a different workload")
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, status.FromContextError(err).Err()
	case err != nil:
		log.Printf("worker RPC: %v", err)
		return nil, status.Error(codes.Internal, "worker state unavailable")
	}
	return response, nil
}
