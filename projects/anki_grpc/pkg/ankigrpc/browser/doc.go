// Package browser connects Go js/wasm applications to Rust Anki commands through
// a JavaScript API. The native and browser clients share protobuf adaptation,
// cancellation, call options and backend lifetime in internal/directclient.
// NewClient owns one Rust instance and its Go timezone callbacks. The caller's
// factory captures a Node-compatible JavaScript filesystem and root, receives
// the separate timezone host, and resolves with synchronous command and close
// methods. wasi.js adapts synchronous provider operations without Go filesystem
// callbacks and owns only its descriptors; the caller retains the storage.
// Start waits for the factory to finish even if its context is canceled during
// initialization. A canceled startup never publishes its backend; a stuck factory
// delays startup failure and cleanup. Command cancellation remains independent.
//
// Applications using Go os install the shared go.js filesystem facade before
// Go.run and select their root with Afero BasePathFs over NewOsFs. The provider
// and Go globals must remain stable while Go runs. Close collections before Stop.
// Done observes final cleanup after a Stop deadline or failed command host.
// This package does not select a backing filesystem or own durable persistence.
package browser
