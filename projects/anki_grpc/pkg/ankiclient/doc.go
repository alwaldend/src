// Package ankiclient supplies gRPC connections for generated Anki clients.
// Client exposes connection access and shutdown; StartableClient embeds it and
// adds transport startup for callers that own initialization. Collection
// open, close, archive preparation and publication belong to callers.
// ClientSTD owns explicitly configured process startup, shutdown, crash reporting
// and reaping. Its Start context controls startup; Stop bounds graceful shutdown
// and forces termination through a backend-owned process context when needed.
// ClientGRPC connects to a validated network endpoint or caller-supplied stream
// dialer. It owns the connection; its caller owns the server lifetime.
// Concrete owned transports expose Done for completion; process transports close it
// after reaping, and in-process work can finish after Stop's caller stops waiting.
// Source reservations and other application policy remain with callers.
// This package does not discover a default executable, link or host the native
// server, parse application configs or own callback scopes, source reservations
// or reconciliation policy. Anki owns storage semantics. RPC cancellation does
// not establish rollback; failed requests are not replayed.
// The optional ankigrpc/native subpackage supplies direct native execution; importing
// this package or the generic ankigrpc constructor never links that engine.
package ankiclient
