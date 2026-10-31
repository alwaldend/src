package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

type Startable interface {
	Start(ctx context.Context) error
}

type StartableFunc func(context.Context) error

func (self StartableFunc) Start(ctx context.Context) error {
	return self(ctx)
}

func (self StartableFunc) Stop(_ context.Context) error {
	return nil
}

func StartableFunc0(f func() error) StartableFunc {
	return StartableFunc(func(_ context.Context) error {
		return f()
	})
}

type Stoppable interface {
	Stop(ctx context.Context) error
}

type StoppableFunc func(context.Context) error

func (self StoppableFunc) Start(_ context.Context) error {
	return nil
}

func (self StoppableFunc) Stop(ctx context.Context) error {
	return self(ctx)
}

func StoppableFunc0(f func() error) StoppableFunc {
	return StoppableFunc(func(_ context.Context) error {
		return f()
	})
}

type Manageable interface {
	Startable
	Stoppable
}

type State int

const (
	StateInit    State = 0
	StateStarted State = 0x1
	StateStopped State = 0x10
	StateError   State = 0x100
)

var StateName = map[State]string{
	StateInit:    "init",
	StateStarted: "started",
	StateStopped: "stopped",
	StateError:   "error",
}

type Manager struct {
	mx           sync.RWMutex
	managed      []Manageable
	managedState []State
	stopped      bool
	parallelStop bool
}

var _ Manageable = (*Manager)(nil)

// NewParallelStopManager is for independent resources whose shutdown may run
// concurrently. Dependency-ordered resources must use the zero-value Manager.
// Stop still waits for every callback and aggregates all cleanup errors.
func NewParallelStopManager() *Manager {
	return &Manager{parallelStop: true}
}

func (self *Manager) AddState(state State, vals ...Manageable) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.stopped {
		return errors.New("lifecycle manager is stopped")
	}
	for _, val := range vals {
		self.managed = append(self.managed, val)
		self.managedState = append(self.managedState, state)
	}
	return nil
}

func (self *Manager) Add(vals ...Manageable) error {
	return self.AddState(StateInit, vals...)
}

func (self *Manager) Run(ctx context.Context, stopTimeout time.Duration) error {
	if err := self.Start(ctx); err != nil {
		return fmt.Errorf("start error: %w", err)
	}
	<-ctx.Done()
	if err := self.StopTimeout(stopTimeout); err != nil {
		return fmt.Errorf("stop error: %w", err)
	}
	return nil
}

func (self *Manager) StopTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := self.Stop(ctx); err != nil {
		return fmt.Errorf("could not stop with timeout: %w", err)
	}
	return nil
}

// Start starts registered resources concurrently. Any failed start rolls back
// every attempted resource, including the resource whose Start failed.
func (self *Manager) Start(ctx context.Context) error {
	startCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	self.mx.Lock()
	if self.stopped {
		self.mx.Unlock()
		return errors.New("lifecycle manager is stopped")
	}
	var wg sync.WaitGroup
	res := make([]error, len(self.managed))
	for i, managed := range self.managed {
		if self.managedState[i] != StateInit {
			continue
		}
		self.managedState[i] = StateStarted
		wg.Go(func() {
			if err := rescueLifecycleFunc(managed.Start)(startCtx); err != nil {
				self.managedState[i] = StateError
				cancel()
				res[i] = fmt.Errorf("could not start %s: %w", objName(managed), err)
			}
		})
	}
	wg.Wait()
	self.mx.Unlock()
	if err := errors.Join(res...); err != nil {
		return errors.Join(err, self.StopTimeout(10*time.Second))
	}
	return nil
}

// Stop releases resources in reverse registration order. Register dependencies
// before their consumers, so consumers drain before their dependencies stop.
// NewParallelStopManager explicitly opts independent resources out of ordering.
// Stop also visits unstarted resources: constructors may already own resources.
// Every registered Stop must therefore be safe when Start has not been called.
// Callbacks must honor ctx; they are never abandoned in background goroutines.
func (self *Manager) Stop(ctx context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.stopped = true
	res := make([]error, len(self.managed))
	var wg sync.WaitGroup
	for i := len(self.managed) - 1; i >= 0; i-- {
		if self.managedState[i] == StateStopped {
			continue
		}
		// A cleanup attempt is final, even when it fails. Report failure to the owner.
		self.managedState[i] = StateStopped
		stop := func() {
			if err := rescueLifecycleFunc(self.managed[i].Stop)(ctx); err != nil {
				res[i] = fmt.Errorf("could not stop %s: %w", objName(self.managed[i]), err)
			}
		}
		if self.parallelStop {
			wg.Go(stop)
		} else {
			stop()
		}
	}
	wg.Wait()
	return errors.Join(res...)
}

func rescueLifecycleFunc(f func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("lifecycle panic (%T)\n%s", r, debug.Stack())
			}
		}()
		return f(ctx)
	}
}

func objName(v any) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Pointer {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

// StopProcess signals a started child, waits for graceful exit, and kills and
// reaps it when its grace period expires. The caller must send cmd.Wait()'s
// result to wait exactly once after Start succeeds, and must not consume it.
func StopProcess(ctx context.Context, cmd *exec.Cmd, wait <-chan error) error {
	return stopProcess(ctx, cmd, wait, false)
}

// StopResourceProcess accepts termination by the SIGTERM this cleanup sends.
// Use it for resource helpers that have no cleanup protocol; plugins must use
// StopProcess so an interrupted or failed plugin cleanup remains an error.
func StopResourceProcess(ctx context.Context, cmd *exec.Cmd, wait <-chan error) error {
	return stopProcess(ctx, cmd, wait, true)
}

func stopProcess(ctx context.Context, cmd *exec.Cmd, wait <-chan error, expectedTerm bool) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	var errs []error
	select {
	case err := <-wait:
		return err
	default:
	}
	signalErr := cmd.Process.Signal(syscall.SIGTERM)
	if signalErr != nil && !errors.Is(signalErr, os.ErrProcessDone) {
		errs = append(errs, fmt.Errorf("signal process: %w", signalErr))
	}
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	select {
	case err := <-wait:
		var exitErr *exec.ExitError
		if expectedTerm && signalErr == nil && errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() && status.Signal() == syscall.SIGTERM {
				err = nil
			}
		}
		errs = append(errs, err)
	case <-waitCtx.Done():
		errs = append(errs, fmt.Errorf("process shutdown deadline: %w", waitCtx.Err()))
		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			errs = append(errs, err)
		}
		errs = append(errs, <-wait)
	}
	return errors.Join(errs...)
}
