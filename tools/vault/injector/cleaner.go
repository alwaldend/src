package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

// Cleaner owns paths until Stop returns. Registration and deletion are serialized.
type Cleaner struct {
	mx      sync.Mutex
	paths   []string
	closed  bool
	stopErr error
}

func NewCleaner(_ *al.CmdCtx) *Cleaner { return &Cleaner{} }

func (self *Cleaner) Add(paths ...string) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.closed {
		return fmt.Errorf("could not add cleanup paths: closed")
	}
	self.paths = append(self.paths, paths...)
	return nil
}

func (self *Cleaner) Start(context.Context) error { return nil }

func (self *Cleaner) Stop(context.Context) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.closed {
		return self.stopErr
	}
	self.closed = true
	for i := len(self.paths) - 1; i >= 0; i-- {
		if err := os.RemoveAll(self.paths[i]); err != nil {
			self.stopErr = errors.Join(self.stopErr, fmt.Errorf("could not remove temporary resource: %w", err))
		}
	}
	self.paths = nil
	return self.stopErr
}
