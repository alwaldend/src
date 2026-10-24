// Package ankigrpc constructs Anki clients from connection strings.
// The generic package is pure Go regardless of enabled native build settings.
// NewClient constructs TCP/Unix clients or selects explicitly registered schemes.
// NewFactory adds a private registry and an explicit default without changing other
// callers. Constructors parse original addresses and perform no startup I/O.
// ClientFactory is a constructor function; a concrete Factory's NewClient method
// carries its registry and default selection. ToClientFactory adapts typed constructors.
// Constructors return StartableClient directly without starting it. Each client
// owns its transport readiness and cleanup. RPC access uses the existing
// grpc.ClientConnInterface. Optional engines register scheme factories through
// explicit subpackage imports. Executables register worker constructors for their
// owned entrypoints. The registry retains constructors, never clients.
// WithClient owns a raw client callback scope with bounded, independent cleanup.
// The caller owns collection workflows, process signals and executable packaging.
package ankigrpc
