// Package browser connects Go js/wasm applications to Rust Anki commands through
// a JavaScript API. The native and browser clients share protobuf adaptation,
// cancellation, call options and backend lifetime in internal/directclient.
// NewClient owns one Rust instance and its WASI descriptors over the caller's
// borrowed Afero filesystem. The application supplies a factory that resolves
// with synchronous command and close methods; wasi.js supplies the WASI imports.
//
// Filesystem operations are synchronous and must not wait for JavaScript. Sync
// flushes the working filesystem only; the application owns durable snapshots.
// Close collections before Stop. Done observes final cleanup after a Stop deadline
// or failed command host. This package does not own IndexedDB persistence.
package browser
