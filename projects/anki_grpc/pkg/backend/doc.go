// Package backend supplies a managed gRPC connection for clients generated from
// Anki's upstream services, carried by a worker's stdin and stdout. Go owns process
// startup, shutdown, and crash reporting; Anki owns collection semantics.
// New accepts BackendSTDOptions without launching a process; Backend exposes Start,
// Connection, and Stop. The Start context bounds startup only; Stop uses its own
// context to bound graceful shutdown and always completes process reaping.
// Its expiry cancels a backend-owned command context for forced termination.
// Connection reports an error when the server is unavailable.
// BackendCollection accepts BackendCollectionOptions and borrows an already-started
// backend. It opens databases in place or restores packages to a temporary directory
// by default, gates collection clients, and closes storage without stopping its backend.
// Callers stop the collection before stopping its backend. Manager scopes own this
// coordination: ManagerSTD.WithBackend lends one worker, and ManagerCollection
// uses a supplied Manager to open and close a collection before transport cleanup.
// Callbacks borrow prepared backends without owning Start/Stop; managers preserve
// callback/cleanup errors and use bounded cleanup independent of cancellation.
// ManagerCollection adds no path registry, publication policy or idle cache.
// RPC cancellation
// stops waiting without guaranteeing rollback. Failed requests are not replayed.
package backend
