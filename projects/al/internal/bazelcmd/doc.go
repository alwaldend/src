// Package bazelcmd selects affected tests through Git paths and Bazel's graph,
// then delegates execution and caching to Bazel. Unmapped changes select the
// requested test scope, and query failures never yield a successful selection.
// Doctor reports bounded source, executable, configuration, and scratch identity.
package bazelcmd
