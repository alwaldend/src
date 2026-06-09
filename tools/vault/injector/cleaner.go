package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

type Cleaner struct {
	consume chan string
	mx      *sync.RWMutex
	logger  *log.Logger
	ctx     *al.CmdCtx
	closed  bool
	toClean []string
}

func NewCleaner(ctx *al.CmdCtx) *Cleaner {
	return &Cleaner{
		consume: make(chan string),
		mx:      &sync.RWMutex{},
		ctx:     ctx,
		logger:  log.New(ctx.Stderr, "com.alwaldend.src.tools.vault.injector.cleaner ", ctx.LogFlags),
	}
}

func (self *Cleaner) Add(paths ...string) error {
	self.mx.RLock()
	defer self.mx.RUnlock()
	if self.closed {
		return fmt.Errorf("could not add: closed")
	}
	for _, p := range paths {
		self.consume <- p
	}
	return nil
}

func (self *Cleaner) Shutdown(ctx context.Context) error {
	var wg fp.WaitGroupE
	self.mx.Lock()
	defer self.mx.Unlock()
	close(self.consume)
	self.closed = true
	for _, path := range self.toClean {
		wg.Go(func() error {
			self.logger.Printf("Cleaning path %s", path)
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("could not clean path %s: %w", path, err)
			}
			return nil
		})
	}
	if err := wg.WaitCtx(ctx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}

func (self *Cleaner) Start() {
	go func() {
		for {
			path, ok := <-self.consume
			if !ok {
				return
			}
			self.toClean = append(self.toClean, path)
		}
	}()
}
