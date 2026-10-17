package backend

import (
	"context"
	"fmt"
)

// ManagerSTD creates and owns a fresh stdin/stdout worker for each backend scope.
// Construct it with NewManagerSTD; its copied configuration supports independent
// overlapping scopes without retaining idle workers or native collection state.
type ManagerSTD struct {
	// options supplies immutable process configuration for each new worker.
	options BackendSTDOptions
}

// NewManagerSTD copies process options without starting a worker.
func NewManagerSTD(options BackendSTDOptions) *ManagerSTD {
	options.Arguments = append([]string(nil), options.Arguments...)
	return &ManagerSTD{options: options}
}

// WithBackend starts a worker, lends it to operation and reaps it before returning.
// Cleanup retains the configured shutdown bound independently of ctx cancellation.
// Callback and cleanup errors remain discoverable; panics propagate after cleanup.
func (manager *ManagerSTD) WithBackend(ctx context.Context, operation func(Backend) error) error {
	worker := New(manager.options)
	if err := withBackend(ctx, worker, worker.options.ShutdownTimeout, operation); err != nil {
		return fmt.Errorf("manage scoped worker: %w", err)
	}
	return nil
}
