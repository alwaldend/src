// Package main executes the embedded Rust Anki CLI through a pure-Go WASI runtime.
// The server package owns transport and guest lifetime; callers own process
// supervision and collection lifetime through the existing protobuf protocol.
package main
