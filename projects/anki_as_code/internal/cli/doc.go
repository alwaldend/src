// Package cli adapts command-line input and output to application operations.
// It owns the complete Cobra tree, flags, syntax, format decoding/encoding,
// input/output streams and forwarding command contexts through ankiascodeapi.App.
// Its MCP command group consumes the injected server through its private mcpServer interface.
// The server serve command prefixes serve before calling App.ServerRun;
// remaining native server arguments stay opaque.
// App owns Config, Col, ColState, API and native Server execution.
// This package does not construct App, host native code, manage process signals
// or collection lifetime, or implement planning and reconciliation policy.
package cli
