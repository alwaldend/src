// Package ankiclient supplies managed gRPC connections for generated Anki clients.
// ClientSTD owns explicitly configured process startup, shutdown, crash reporting
// and reaping. Its Start context controls startup; Stop bounds graceful shutdown
// and forces termination through a backend-owned process context when needed.
// ClientGRPC validates its endpoint during construction, owns only its remote
// client connection and never stops the server. AsyncClient supplies Done when
// Stop can return before final cleanup;
// callers keep source reservations until that channel closes.
// ClientCollection borrows an already-started backend, opens databases or restores
// packages, gates collection clients and closes native storage without stopping
// the backend. Callers compose these lifetimes and stop collections before workers.
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
