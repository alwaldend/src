// Package ankigrpc constructs managed Anki clients from connection strings.
// The generic package is pure Go regardless of enabled native build settings.
// NewClient constructs TCP/Unix clients or selects explicitly registered schemes.
// NewFactory adds a private registry and an explicit default without changing other
// callers. Constructors parse original addresses and perform no startup I/O.
// Factories return one Client owning transport readiness, optional collection state
// and cleanup. Transport constructors are adapted without starting them. RPC access uses the existing
// grpc.ClientConnInterface. Optional engines register scheme factories through
// explicit subpackage imports. Executables register worker constructors for their
// owned entrypoints. The registry retains constructors, never clients.
// WithClient owns a raw client callback scope with bounded, independent cleanup.
// The caller owns collection workflows, process signals and executable packaging.
package ankigrpc
