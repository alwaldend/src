package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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
	managedErrs  []error
}

var _ Manageable = (*Manager)(nil)

func (self *Manager) AddState(state State, vals ...Manageable) {
	self.mx.Lock()
	defer self.mx.Unlock()
	for _, val := range vals {
		self.managed = append(self.managed, val)
		self.managedState = append(self.managedState, state)
		self.managedErrs = append(self.managedErrs, nil)
	}
}

func (self *Manager) Add(vals ...Manageable) {
	self.AddState(StateInit, vals...)
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

func (self *Manager) Start(ctx context.Context) error {
	if err := self.transition(ctx, StateInit, StateStarted, func(m Manageable) lifecycleFunc { return m.Start }); err != nil {
		return fmt.Errorf("could not start: %w", err)
	}
	return nil
}

func (self *Manager) Stop(ctx context.Context) error {
	if err := self.transition(ctx, StateStarted, StateStopped, func(m Manageable) lifecycleFunc { return m.Stop }); err != nil {
		return fmt.Errorf("could not stop: %w", err)
	}
	return nil
}

func (self *Manager) transition(ctx context.Context, fromState State, targetState State, selector selectorFunc) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	var wg sync.WaitGroup
	res := make([]error, len(self.managed))
	for i, managed := range self.managed {
		if self.managedState[i] != fromState {
			continue
		}
		name := objName(managed)
		if err := self.managedErrs[i]; err != nil {
			res[i] = fmt.Errorf("could not transition %s: already failed: %w", name, self.managedErrs[i])
			continue
		}
		wg.Go(func() {
			if err := rescueLifecycleFunc(selector(managed))(ctx); err != nil {
				self.managedErrs[i] = err
				self.managedState[i] = StateError
				res[i] = fmt.Errorf("could not transition %s: %w", name, err)
			} else {
				self.managedState[i] = targetState
			}
		})
	}
	wg.Wait()
	if err := errors.Join(res...); err != nil {
		return fmt.Errorf("could not transition from %s to %s: %w", StateName[fromState], StateName[targetState], err)
	}
	return nil
}

type lifecycleFunc = func(context.Context) error
type selectorFunc = func(Manageable) lifecycleFunc

func rescueLifecycleFunc(f func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		local := make(chan error, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					switch rTyped := r.(type) {
					case error:
						local <- fmt.Errorf("lifecycle panic: %s\n%s: %w", rTyped, string(debug.Stack()), rTyped)
					default:
						local <- fmt.Errorf("lifecycle panic: %s\n%s", r, string(debug.Stack()))
					}
				}
			}()
			local <- f(ctx)
		}()
		select {
		case err := <-local:
			return err
		case <-ctx.Done():
			return fmt.Errorf("timed out")
		}
	}
}

func objName(v any) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Pointer {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
