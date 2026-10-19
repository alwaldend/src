## Why

Review native adapter extraction and WASI build preparation independently from the WASM server implementation.

## What Changes

- Move the existing Go/C ABI adapter to pkg/ankiservercgo.
- Prepare command target names and a cgo-only backend setting without exposing an unavailable server.
- Use built-in Rust platform selection for the existing SQLite native/WASI acceptance fixture.

## Capabilities

No specification changes: this extracts native composition and build preparation while preserving runtime behavior. `skip_specs: true`.

## Impact

Native gRPC CLI, cgo adapter, Bazel targets and storage build fixture.
