// Package main composes the anki_as_code executable with its native cgo server
// adapter. It owns signals, constructs App with a same-executable worker, and
// runs the CLI. internal/cli owns Cobra commands; App owns their operations.
// The server adapter owns native CLI execution; this package implements no
// collection workflows.
package main
