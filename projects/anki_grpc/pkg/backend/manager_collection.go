package backend

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// collectionScopeCleanupTimeout bounds native close independently of callback cancellation.
const collectionScopeCleanupTimeout = 5 * time.Second

// ManagerCollection opens scoped collections through an injected transport manager.
// It owns native collection cleanup while the supplied Manager owns its transport.
// It retains no source reservations, publication policy, resource clients or cache.
// Callers coordinate same-source access and satisfy the transport's exclusivity rules.
type ManagerCollection struct {
	// manager lends the prepared transports used for native collection scopes.
	manager Manager
}

// NewManagerCollection selects a transport manager without acquiring native resources.
func NewManagerCollection(manager Manager) *ManagerCollection {
	return &ManagerCollection{manager: manager}
}

// WithCollection lends a prepared collection, closes it, then releases its transport scope.
// Paths follow BackendCollection's direct-database and temporary-package contracts.
// Collection cleanup has a five-second bound independent of ctx cancellation.
// Callbacks borrow the prepared backend; managers retain lifecycle ownership.
// Errors preserve native and cleanup causes.
func (manager *ManagerCollection) WithCollection(ctx context.Context, path string, operation func(Backend) error) error {
	if operation == nil {
		return fmt.Errorf("collection scope requires an operation")
	}
	if manager.manager == nil {
		return fmt.Errorf("collection scope requires a backend manager")
	}
	if path == "" {
		return fmt.Errorf("collection scope requires a path")
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start collection scope: %w", errors.Join(err, context.Cause(ctx)))
	}
	if err := manager.manager.WithBackend(ctx, func(worker Backend) error {
		prepared := NewCollection(BackendCollectionOptions{Backend: worker, Path: path})
		if err := withBackend(ctx, prepared, collectionScopeCleanupTimeout, operation); err != nil {
			return fmt.Errorf("manage scoped collection: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("manage collection %q: %w", path, err)
	}
	return nil
}
