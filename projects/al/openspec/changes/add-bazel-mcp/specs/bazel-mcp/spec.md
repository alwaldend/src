## Purpose

Provide agents with discoverable Bazel build and inspection operations through
MCP while preserving ordinary Bazel workspace configuration and execution.

## ADDED Requirements

### Requirement: Workspace-bound stdio tools

AL SHALL offer `al bazel mcp serve` over stdio with query, cquery, aquery, build, test,
and info tools bound to the selected enclosing Bazel workspace. Protocol output
SHALL remain separate from child-process diagnostics. Existing affected-test
selection SHALL retain its behavior.

#### Scenario: Discover and invoke tools

- **WHEN** a client initializes and lists tools
- **THEN** it receives input schemas for all six operations and can invoke them
  against the selected workspace using its normal Bazel configuration

#### Scenario: Invalid invocation

- **WHEN** a tool receives invalid arguments or the server cannot locate a workspace
- **THEN** AL reports an error without starting a Bazel operation

### Requirement: Bounded results and execution

Tool results SHALL include workspace, command, exit status, bounded stdout and
stderr, truncation indicators, and timeout/cancellation indicators. Nonzero or
interrupted execution SHALL be reported as a tool error. Requests SHALL have a
finite operator-selected deadline and honor MCP cancellation and server shutdown.

#### Scenario: Failed or excessive output

- **WHEN** Bazel fails or produces output exceeding the configured byte limit
- **THEN** the response preserves its exit status and marks failure or truncation
  without presenting partial output as a complete successful query

#### Scenario: Deadline and cancellation

- **WHEN** the deadline expires, a request is cancelled, or the server shuts down
- **THEN** active Bazel clients are cancelled and the stdio session can terminate
