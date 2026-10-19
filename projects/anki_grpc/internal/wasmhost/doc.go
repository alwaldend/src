// Package wasmhost executes complete Rust WASI CLI commands in a pure-Go runtime.
// It supplies checked host byte streams, readiness, filesystem access and timezone
// callbacks. Rust owns argument parsing, gRPC framing, protobuf dispatch and Anki.
// Production adapters and integration fixtures supply command bytes. Callers own
// process supervision; cancelling execution closes the guest without replay.
// ANKI_WASM_CACHE_DIR optionally selects caller-protected compiled-code storage.
package wasmhost
