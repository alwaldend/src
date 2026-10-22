// Package main exposes Anki as code application workflows to a browser worker.
// Each runtime owns an Afero memory workspace shared with scoped Rust WASI servers
// through the browser backend. Requests call the existing App and protobuf codec;
// the worker owns ordering and IndexedDB persistence. Completed workspace snapshots
// exclude reserved /tmp runtime storage and require all backend cleanup to finish.
package main
