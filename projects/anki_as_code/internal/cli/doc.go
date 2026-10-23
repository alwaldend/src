// Package cli adapts command-line input and output to application operations.
// It owns the complete Cobra tree, flags, syntax, format decoding/encoding,
// input/output streams and forwarding command contexts through ankiascodeapi.App.
// Its MCP command delegates hosting through an independently injected ankiascodeapi.MCP.
// Native server arguments remain opaque.
// App owns Config, Col, ColState, API and Server execution, including serve delegation.
// This package does not construct App, host native code, manage process signals
// or collection lifetime, or implement planning and reconciliation policy.
package cli
