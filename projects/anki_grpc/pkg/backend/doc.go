// Package backend supplies a managed gRPC connection for clients generated from
// Anki's upstream services, carried by a worker's stdin and stdout. Go owns process
// startup, shutdown, and crash reporting; Anki owns collection semantics.
// New accepts BackendSTDOptions without launching a process; Backend exposes Start,
// Connection, and Stop. The Start context bounds startup only; Stop uses its own
// context to bound graceful shutdown and always completes process reaping.
// Its expiry cancels a backend-owned command context for forced termination.
// Connection reports an error when the server is unavailable.
// Callers close an open collection before stopping its backend. RPC cancellation
// stops waiting without guaranteeing rollback. Failed requests are not replayed.
package backend
