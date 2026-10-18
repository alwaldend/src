// Package backend supplies a managed gRPC connection for clients generated from
// Anki's upstream services, carried over stdio, TCP or Unix sockets. Go owns process
// startup, shutdown, and crash reporting; Anki owns collection semantics.
// New accepts BackendSTDOptions without launching a process; Backend exposes Start,
// Connection, and Stop. The Start context bounds startup only; Stop uses its own
// context to bound graceful shutdown and always completes process reaping.
// Its expiry cancels a backend-owned command context for forced termination.
// Connection reports an error when the server is unavailable.
// BackendCollection accepts BackendCollectionOptions and borrows an already-started
// backend. It opens databases in place or restores packages to a temporary directory
// by default, gates collection clients, and closes storage without stopping its backend.
// Callers stop the collection before stopping its backend and own composition of
// these lifetimes. Package sources remain unchanged unless WriteBack is explicitly
// enabled; request outcomes do not control enabled publication. Failed publication retains
// restored data for recovery. Application callback scopes and source reservations
// belong to consumers. BackendGRPC owns only a remote client connection.
// RPC cancellation
// stops waiting without guaranteeing rollback. Failed requests are not replayed.
package backend
