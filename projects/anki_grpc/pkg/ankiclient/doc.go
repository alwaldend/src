// Package ankiclient supplies managed gRPC connections for generated Anki clients.
// ClientSTD owns explicitly configured process startup, shutdown, crash reporting
// and reaping. Its Start context controls startup; Stop bounds graceful shutdown
// and forces termination through a backend-owned process context when needed.
// ClientGRPC connects to a validated network endpoint or caller-supplied stream
// dialer. It owns the connection; its caller owns the server lifetime.
// NewClient gives a transport one Client lifecycle with optional private collection state.
// Start prepares its transport; OpenCollection opens a database or restores a package.
// Connection lends generated-client access, gated by the current collection's readiness.
// CloseCollection publishes and closes storage while retaining the transport for reuse.
// Stop closes storage before stopping the transport, with independent bounded budgets.
// Its context bounds waiting; Done observes final native completion. AsyncClient is
// the optional completion capability of a transport whose Stop can return early.
// Source reservations and other application policy remain with callers.
// Package sources remain unchanged unless WriteBack is explicitly enabled;
// request outcomes do not suppress enabled publication. Failed publication or
// uncertain native completion retains restored data for recovery.
// This package does not discover a default executable, link or host the native
// server, parse application configs or own callback scopes, source reservations
// or reconciliation policy. Anki owns storage semantics. RPC cancellation does
// not establish rollback; failed requests are not replayed.
// The optional ankigrpc/native subpackage supplies direct native execution; importing
// this package or the generic ankigrpc constructor never links that engine.
package ankiclient
