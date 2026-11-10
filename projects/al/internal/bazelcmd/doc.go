// Package bazelcmd selects affected tests through Git paths and Bazel's graph,
// then delegates execution and caching to Bazel. Unmapped changes select the
// requested test scope, and query failures never yield a successful selection.
// Doctor reports bounded source, executable, configuration, and scratch identity.
// The opt-in MCP stdio adapter exposes workspace-bound Bazel operations with
// bounded output and deadlines; Bazel owns configuration and action execution.
package bazelcmd
