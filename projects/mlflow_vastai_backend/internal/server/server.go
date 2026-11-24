package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	vastai "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/sqlite"
	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
)

const (
	// stateReady waits for an invocation of the job server.
	stateReady machine.State = "coordinator-ready"
	// stateServing owns the listener and reconciliation loop.
	stateServing machine.State = "coordinator-serving"
	// start begins a coordinator invocation.
	start machine.Event = "serve"
)

// Config describes one coordinator invocation.
type Config struct {
	// Listen is a literal loopback address with a TCP port.
	Listen string
	// Directory holds the private coordinator database and lock.
	Directory string
	// Workers lists allowed rental IDs or exact labels in selection order.
	Workers []string
	// VastEndpoint is the versioned Vast.ai API root.
	VastEndpoint string
	// TrackingURI is the HTTP or HTTPS endpoint shared with submitted jobs.
	TrackingURI string
	// SSHUser is the worker account on each rental.
	SSHUser string
}

// invocation carries the settings and filesystem used by the server.
type invocation struct {
	// config holds the server configuration.
	config Config
	// filesystem owns coordinator file operations.
	filesystem afero.Fs
}

// listenerMachine owns a coordinator listener invocation.
type listenerMachine struct {
	// definition drives service startup and shutdown.
	definition *machine.Machine[invocation]
}

// newListenerMachine configures serving without opening resources.
func newListenerMachine() *listenerMachine {
	listener := &listenerMachine{definition: machine.New[invocation](stateReady)}
	listener.definition.Configure(stateReady).Permit(start, stateServing)
	listener.definition.Configure(stateServing).OnEntry(listener.serve)
	return listener
}

// Serve runs the job server until the context is cancelled or serving fails.
func Serve(ctx context.Context, filesystem afero.Fs, config Config) error {
	listener := newListenerMachine()
	state := stateReady
	ctx = listener.definition.WithInvocation(ctx, &state, invocation{config: config, filesystem: filesystem})
	if err := listener.definition.FireCtx(ctx, start); err != nil {
		return fmt.Errorf("serve coordinator: %w", err)
	}
	return nil
}

// Graph describes the coordinator and its job operations.
func Graph() machine.GraphSection {
	return machine.GraphSection{
		ID: "coordinator", Title: "Job server", DOT: newListenerMachine().definition.ToGraph(),
		Calls: []machine.GraphCall{
			{From: stateServing, Label: "coordinates jobs", Target: jobGraph()},
		},
	}
}

// serve opens resources, starts the gRPC service, and drains work on shutdown.
func (listener *listenerMachine) serve(ctx context.Context, _ ...any) (err error) {
	current, err := listener.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read coordinator settings: %w", err)
	}
	config := current.config
	endpoint, err := config.validate()
	if err != nil {
		return fmt.Errorf("validate coordinator settings: %w", err)
	}
	api, err := vastai.New(vastai.Options{Endpoint: config.VastEndpoint, APIKey: os.Getenv("VAST_API_KEY")})
	if err != nil {
		return fmt.Errorf("configure rental discovery: %w", err)
	}
	directory, err := filepath.Abs(config.Directory)
	if err != nil {
		return fmt.Errorf("resolve coordinator directory: %w", err)
	}
	lock, err := lockDirectory(current.filesystem, directory)
	if err != nil {
		return fmt.Errorf("lock coordinator directory: %w", err)
	}
	defer func() { err = errors.Join(err, lock.Close()) }()
	db, err := sqlite.Open(ctx, filepath.Join(directory, "jobs.sqlite3"))
	if err != nil {
		return fmt.Errorf("open coordinator database: %w", err)
	}
	defer func() { err = errors.Join(err, db.Close()) }()
	repository, err := openRepository(ctx, db.Connection())
	if err != nil {
		return fmt.Errorf("open coordinator repository: %w", err)
	}
	jobs := newCoordinator(repository, rentals{api: api, allowed: config.Workers, user: config.SSHUser},
		newTracking(endpoint, os.Getenv("MLFLOW_TRACKING_TOKEN")))
	socket, err := net.Listen("tcp", config.Listen)
	if err != nil {
		return fmt.Errorf("bind coordinator listener: %w", err)
	}
	grpcServer := grpc.NewServer(grpc.WaitForHandlers(true))
	jobpb.RegisterJobServiceServer(grpcServer, jobs)
	runCtx, stop := context.WithCancel(ctx)
	stopped, reconciled := make(chan struct{}), make(chan struct{})
	context.AfterFunc(runCtx, func() {
		defer close(stopped)
		grpcServer.Stop()
	})
	go func() {
		defer close(reconciled)
		jobs.reconcile(runCtx)
	}()
	defer func() {
		stop()
		<-stopped
		<-reconciled
	}()
	if err := grpcServer.Serve(socket); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve coordinator gRPC: %w", err)
	}
	return nil
}

// sshUser accepts the unquoted account name passed via OpenSSH's -l option.
var sshUser = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.-]*$`)

// validate restricts job submission to a trusted local service and fixed tracking endpoint.
func (config Config) validate() (*url.URL, error) {
	address, err := netip.ParseAddrPort(config.Listen)
	if err != nil || !address.Addr().IsLoopback() {
		return nil, errors.New("listen address must be a literal loopback IP and TCP port")
	}
	if config.Directory == "" || len(config.Workers) == 0 || !sshUser.MatchString(config.SSHUser) {
		return nil, errors.New("state directory, allowed workers, and a valid SSH user are required")
	}
	for _, value := range config.Workers {
		if strings.TrimSpace(value) != value || value == "" {
			return nil, errors.New("worker selectors must be nonempty IDs or exact labels")
		}
	}
	endpoint, err := url.Parse(config.TrackingURI)
	if err != nil || endpoint.Hostname() == "" || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return nil, errors.New("tracking endpoint requires a host and must omit credentials, query, and fragment")
	}
	if endpoint.Scheme != "https" && !(endpoint.Scheme == "http" && net.ParseIP(endpoint.Hostname()).IsLoopback()) {
		return nil, errors.New("tracking endpoint requires HTTPS or literal loopback HTTP")
	}
	return endpoint, nil
}

// lockDirectory confines durable state to its owner and excludes other coordinators.
func lockDirectory(filesystem afero.Fs, directory string) (afero.File, error) {
	if err := filesystem.MkdirAll(directory, 0o700); err != nil {
		return nil, fmt.Errorf("create state directory: %w", err)
	}
	info, err := filesystem.Stat(directory)
	if err != nil {
		return nil, fmt.Errorf("inspect state directory: %w", err)
	}
	if !info.IsDir() || info.Mode().Perm()&0o077 != 0 {
		return nil, errors.New("state directory must be private to its owner")
	}
	file, err := filesystem.OpenFile(filepath.Join(directory, "server.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open coordinator lock: %w", err)
	}
	descriptor, ok := file.(interface{ Fd() uintptr })
	if !ok {
		return nil, errors.Join(errors.New("coordinator lock requires a file descriptor"), file.Close())
	}
	if err := unix.Flock(int(descriptor.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		return nil, errors.Join(fmt.Errorf("acquire coordinator lock: %w", err), file.Close())
	}
	return file, nil
}
