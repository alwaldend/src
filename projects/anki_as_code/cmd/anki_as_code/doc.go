// Package main composes the anki_as_code executable with opt-in compiled server
// adapters. It owns signals, constructs App with a same-executable worker, and
// runs the CLI. internal/cli owns Cobra commands; App owns their operations.
// The server adapter owns native CLI execution; this package implements no
// collection workflows.
package main
