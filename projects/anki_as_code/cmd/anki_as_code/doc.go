// Package main composes the executable with optional cgo and WASM engines.
// It owns signals, worker entrypoints and a scoped client factory. Generic
// subprocess and in-process clients prefer cgo and otherwise use WASM.
// Qualified subprocess schemes select either compiled engine for same-executable
// workers through the host's --server-engine flag. Independent workers retain
// their caller-supplied arguments.
// With neither engine, explicit external clients remain available. Application
// libraries own collection workflows and remain independent of engine imports.
package main
