package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
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
	StateInit State = iota
	StateStart
	StateStop
	StateError
)

var StateName = map[State]string{
	StateInit:  "init",
	StateStart: "start",
	StateStop:  "stop",
	StateError: "error",
}

type Manager struct {
	mx           sync.RWMutex
	managed      []Manageable
	managedState []State
	managedErrs  []error
	state        State
}

var _ Manageable = (*Manager)(nil)

func (self *Manager) Add(vals ...Manageable) {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.managed = append(self.managed, vals...)
}

func (self *Manager) Run(ctx context.Context, stopTimeout time.Duration) error {
	if self.state != StateInit {
		return fmt.Errorf("cannot start from %s", StateName[self.state])
	}
	if err := self.Start(ctx); err != nil {
		return fmt.Errorf("could not start: %w", err)
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
		return fmt.Errorf("stop error: %w", err)
	}
	return nil
}

func (self *Manager) Start(ctx context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.state != StateInit {
		return fmt.Errorf("cannot start from %s", StateName[self.state])
	}
	if len(self.managed) != 0 {
		if err := self.start(ctx, self.managed); err != nil {
			return fmt.Errorf("could not start: %w", err)
		}
	}
	self.state = StateStart
	return nil
}

func (self *Manager) start(ctx context.Context, vals []Manageable) error {
	res := make(chan error, len(vals))
	for _, start := range vals {
		go lifecycleFunc(start.Start)(ctx, res)
	}
	var err error
	for curErr := range res {
		err = errors.Join(err, curErr)
	}
	if err != nil {
		return fmt.Errorf("start failed: %w", err)
	}
	return nil
}

func (self *Manager) Stop(ctx context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.state != StateStart {
		return fmt.Errorf("cannot stop in state %s", StateName[self.state])
	}
	if len(self.managed) != 0 {
		if err := self.stop(ctx); err != nil {
			return err
		}
	}
	self.state = StateStop
	return nil
}

func (self *Manager) stop(ctx context.Context) error {
	res := make(chan error, len(self.managed))
	for i, stop := range self.managed {
		go lifecycleFunc(stop.Stop)(ctx, res)
	}
	var err error
	for curErr := range res {
		err = errors.Join(err, curErr)
	}
	if err != nil {
		return fmt.Errorf("stop failed: %w", err)
	}
	return nil
}

func lifecycleFunc(f func(context.Context) error) func(context.Context, chan error) {
	run := func(ctx context.Context, res chan error) {
		defer func() {
			if r := recover(); r != nil {
				switch rTyped := r.(type) {
				case error:
					res <- fmt.Errorf("lifecycle function paniced: %s\n%s: %w", rTyped, string(debug.Stack()), rTyped)
				default:
					res <- fmt.Errorf("lifecycle function paniced: %s\n%s", r, string(debug.Stack()))
				}
			}
		}()
		res <- f(ctx)
	}
	return func(ctx context.Context, res chan error) {
		local := make(chan error, 1)
		go run(ctx, local)
		select {
		case err := <-local:
			res <- err
		case <-ctx.Done():
			res <- fmt.Errorf("timed out")
		}
	}
}
