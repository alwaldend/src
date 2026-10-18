// Package cli adapts command-line input and output to application operations.
// It owns the complete Cobra tree, flags, syntax, format decoding/encoding,
// input/output streams and forwarding command contexts through its private
// application interface. Native server arguments remain opaque.
// App owns Config, Col, ColState and Server execution, including serve delegation.
// This package does not construct App, host native code, manage process signals
// or collection lifetime, or implement planning and reconciliation policy.
package cli
