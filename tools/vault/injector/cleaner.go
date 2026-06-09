package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
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
	self.mx.Lock()
	defer self.mx.Unlock()
	close(self.consume)
	self.closed = true
	errsCh := make(chan error, len(self.toClean))
	errs := []error{}
	for _, path := range self.toClean {
		go func() {
			self.logger.Printf("Cleaning path %s", path)
			err := os.RemoveAll(path)
			if err != nil {
				err = fmt.Errorf("could not clean path %s: %w", path, err)
			}
			errsCh <- err
		}()
	}
	for range len(errsCh) {
		select {
		case err := <-errsCh:
			errs = append(errs, err)
		case <-ctx.Done():
			return fmt.Errorf("shutdown timed out")
		}
	}
	return errors.Join(errs...)
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
