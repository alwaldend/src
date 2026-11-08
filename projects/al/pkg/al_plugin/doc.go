// Package al_plugin coordinates AL plugin processes, RPC requests, and owned resources.
// Callers select plugin labels and retain the manager until its consumers stop.
// The manager owns built-in injection and discovers eligible resource providers
// through startup capability advertisements; provider-specific dispatch stays remote.
// Clients retain optional generic startup results independently of capability discovery.
// Absolute executables do not require Bazel runfiles. Selected environment exports
// must agree on shared keys unless explicit injection overrides them; diagnostics omit values.
package al_plugin
