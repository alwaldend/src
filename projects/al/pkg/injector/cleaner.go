package injector

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Cleaner owns paths until Stop returns. Registration and deletion are serialized.
type Cleaner struct {
	// mx serializes path registration and deletion.
	mx sync.Mutex
	// paths records owned paths in registration order.
	paths []string
	// closed prevents registration after cleanup begins.
	closed bool
	// stopErr retains the cleanup result for idempotent Stop calls.
	stopErr error
}

// NewCleaner creates an empty invocation-owned path registry.
func NewCleaner() *Cleaner { return &Cleaner{} }

// Add registers paths and rejects registration after cleanup.
func (self *Cleaner) Add(paths ...string) error {
	self.mx.Lock()
	defer self.mx.Unlock()
	if self.closed {
		return fmt.Errorf("could not add cleanup paths: closed")
	}
	self.paths = append(self.paths, paths...)
	return nil
}

// Start initializes the invocation lifecycle without fetching resources.
func (self *Cleaner) Start(context.Context) error { return nil }

// Stop releases owned resources and reports cleanup failures.
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
