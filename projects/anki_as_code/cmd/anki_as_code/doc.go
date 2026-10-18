// Package main composes the anki_as_code executable and its process adapters.
// It constructs App, supplies the current executable and linked native runner,
// installs process signal handling, executes the CLI and reports exit failures.
// internal/cli owns all Cobra commands; App owns the operations they invoke.
// This package does not implement server dispatch, collection workflows,
// configuration policy, planning or reconciliation.
package main
