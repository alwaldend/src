package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

type Startable interface {
	Start(ctx context.Context) error
}

type Stoppable interface {
	Stop(ctx context.Context) error
}

type State int

const (
	StateInit State = iota
	StateStarted
	StateStopped
)

var StateName = map[State]string{
	StateInit:    "init",
	StateStarted: "started",
	StateStopped: "stopped",
}

type Handler struct {
	startMx sync.RWMutex
	start   []Startable
	stopMx  sync.RWMutex
	stop    []Stoppable
	state   State
}

func (self *Handler) Add(vals ...any) error {
	for _, val := range vals {
		valStart, okStart := val.(Startable)
		if okStart {
			if err := self.AddStart(valStart); err != nil {
				return fmt.Errorf("could not add startable: %w", err)
			}
		}
		valStop, okStop := val.(Stoppable)
		if okStop {
			if err := self.AddStop(valStop); err != nil {
				return fmt.Errorf("could not add stoppable: %w", err)
			}
		}
		if !okStart && !okStop {
			return fmt.Errorf("value is neither Startable, nor Stoppable: %T", val)
		}
	}
	return nil
}

func (self *Handler) Start(ctx context.Context) error {
	self.startMx.Lock()
	defer self.startMx.Unlock()
	if len(self.start) == 0 {
		return nil
	}
	if self.state != StateInit {
		return fmt.Errorf("cannot start state %s", StateName[self.state])
	}
	res := make(chan error, len(self.start))
	for _, start := range self.start {
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

func (self *Handler) AddStart(vals ...Startable) error {
	self.startMx.Lock()
	defer self.startMx.Unlock()
	if self.state != StateInit {
		return fmt.Errorf("cannot add Startable while in state %s", StateName[self.state])
	}
	self.start = append(self.start, vals...)
	return nil
}

func (self *Handler) Stop(ctx context.Context) error {
	self.stopMx.Lock()
	defer self.stopMx.Unlock()
	if len(self.stop) == 0 {
		return nil
	}
	if self.state != StateStarted {
		return fmt.Errorf("cannot stop state %s", StateName[self.state])
	}
	res := make(chan error, len(self.stop))
	for _, stop := range self.stop {
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

func (self *Handler) AddStop(vals ...Stoppable) error {
	self.stopMx.Lock()
	defer self.stopMx.Unlock()
	if self.state == StateStopped {
		return fmt.Errorf("cannot add Stoppable while in state %s", StateName[self.state])
	}
	self.stop = append(self.stop, vals...)
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
