package backend

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Manager lends prepared backends within scopes and retains their lifecycle ownership.
// Callbacks finish all operations before returning and must not call Start or Stop.
// Implementations preserve callback errors and release resources before returning,
// including failed startup and panic unwinding. A collection scope requires an
// exclusive native backend; implementations can serialize access to shared servers.
type Manager interface {
	// WithBackend lends a running backend until the callback completes, then cleans up.
	WithBackend(ctx context.Context, operation func(Backend) error) error
}

// withBackend starts an owned backend, lends it to operation and stops it after
// every startup attempt. Cleanup retains context values and uses timeout
// independently of caller cancellation. Panics propagate after deferred cleanup.
func withBackend(ctx context.Context, owned Backend, timeout time.Duration, operation func(Backend) error) (err error) {
	if operation == nil {
		return fmt.Errorf("backend scope requires an operation")
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start backend scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
		defer cancel()
		if stopError := owned.Stop(cleanup); stopError != nil {
			err = errors.Join(err, fmt.Errorf("stop managed backend: %w", stopError))
		}
	}()
	if err := owned.Start(ctx); err != nil {
		return fmt.Errorf("start scoped backend: %w", err)
	}
	if err := operation(owned); err != nil {
		return fmt.Errorf("run backend scope: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish backend scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	return nil
}
