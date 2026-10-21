// Package wasmcache configures an isolated compilation cache for application E2E tests.
// It owns only test setup: workers still execute their embedded WASM in separate
// processes, and the cache resides in Bazel's task-owned temporary directory.
package wasmcache
