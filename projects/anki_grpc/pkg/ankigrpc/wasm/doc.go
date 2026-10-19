// Package wasm builds, embeds and executes the offline Anki Rust CLI
// as a WASI command through its argument-vector API. Rust owns gRPC and backend
// dispatch; internal/wasmhost supplies a pure-Go runtime and native byte streams.
// Run hosts a command on process streams. New constructs a managed in-process
// client over in-memory stdio using ankiclient's shared gRPC transport;
// each client owns an independent instance.
// Callers close collections before Stop, and Done observes final cleanup after
// a Stop deadline. No companion file or host WASM installation is required.
package wasm
