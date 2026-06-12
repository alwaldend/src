package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
)

type Cleaner struct {
	consume chan []string
	mx      *sync.RWMutex
	ctx     *al.CmdCtx
	closed  bool
	done    chan struct{}
	lc      lifecycle.Manager
}

func NewCleaner(ctx *al.CmdCtx) *Cleaner {
	return &Cleaner{
		consume: make(chan []string),
		done:    make(chan struct{}, 1),
		mx:      &sync.RWMutex{},
		ctx:     ctx,
	}
}

func (self *Cleaner) Add(paths ...string) error {
	self.mx.RLock()
	defer self.mx.RUnlock()
	if self.closed {
		return fmt.Errorf("could not add: closed")
	}
	self.consume <- paths
	return nil
}

func (self *Cleaner) Stop(ctx context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	self.closed = true
	close(self.consume)
	<-self.done
	return self.lc.Stop(ctx)
}

func (self *Cleaner) Start(_ context.Context) error {
	run := func() bool {
		paths, ok := <-self.consume
		if !ok {
			return true
		}
		for _, path := range paths {
			self.lc.AddState(lifecycle.StateStarted, lifecycle.StoppableFunc(func(_ context.Context) error {
				self.ctx.Logger.Printf("cleaning path %s", path)
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("could not clean path %s: %w", path, err)
				}
				return nil
			}))
		}
		return false
	}
	go func() {
		for {
			if ok := run(); ok {
				break
			}
		}
		self.done <- struct{}{}
	}()
	return nil
}
