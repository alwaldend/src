// Package al_plugin coordinates AL plugin processes, RPC requests, and owned resources.
// Callers select plugin labels and retain the manager until its consumers stop.
// Absolute executables do not require Bazel runfiles. Selected environment exports
// must agree on shared keys before a consumer starts; diagnostics omit their values.
package al_plugin
