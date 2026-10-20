// Package ankigrpc constructs managed Anki clients from connection strings.
// The generic package is pure Go regardless of enabled native build settings.
// NewClient connects to TCP/Unix endpoints or selects explicitly registered schemes. Client owns readiness and cleanup; RPC access uses the existing
// grpc.ClientConnInterface. Optional engines register scheme factories through
// explicit subpackage imports. Executables register worker constructors for their
// owned entrypoints. The registry retains constructors, never clients.
// WithClient owns a raw client callback scope with bounded, independent cleanup.
// The caller owns collection workflows, process signals and executable packaging.
package ankigrpc
