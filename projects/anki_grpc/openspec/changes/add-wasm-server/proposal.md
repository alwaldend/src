## Why

Applications currently need cgo and native Rust linkage to embed the Anki
worker. Add a WASM execution path while preserving the supervised child process
and the existing protobuf/gRPC interface. Separate Go packages must keep native
linkage out of the WASM executable's dependency graph.

## What Changes

- Build on the native cgo adapter, shared Rust CLI and dependency inputs merged
  in #182 and the guest dependency adaptations merged in #183. Add WASM
  execution and integration while preserving native behavior.
- Add `pkg/ankigrpc/wasm` for a pure-Go runtime hosting an embedded Rust WASM
  command. Execute the shared Rust CLI and gRPC server; Go supplies only WASI
  execution, host byte streams and timezone callbacks. Keep application and
  process lifecycle in their existing owners.
- Preserve trunk's optional-engine composition with independent
  `//projects/anki_grpc:cgo` and `//projects/anki_grpc:wasm` booleans. WASM-only
  builds disable native linkage; combined builds retain native in-process
  clients while WASM hosts the worker. With neither engine, retain external
  clients. Keep scoped client factories and exact subprocess argument vectors.
- Keep the same-executable worker model and gRPC over stdio. Retain explicit
  TCP/Unix serving and native error details at the transport boundary.
- Add an explicit managed in-process WASM client over in-memory gRPC streams,
  retaining the default worker and existing direct native client.
- Port and validate the pinned Anki backend for the guest environment before
  declaring the WASM worker usable. Unsupported behavior cannot silently fall
  back to cgo or be accepted as equivalent behavior.
- Implement the user's approved offline-first guest scope: local collections,
  packages and media remain supported; network sync, media sync, HTTP serving
  and network-dependent services return explicit unsupported errors. The cgo
  backend retains full upstream behavior.

## Capabilities

### New Capabilities

- `wasm-worker`: An embedded Rust WASM backend hosted in an isolated Go worker
  process, with dependency-isolated backend packages and reproducible builds.

### Modified Capabilities

None. Existing unsynchronized native change records remain owned by their
original changes.

## Impact

The server packages, worker build targets, tests, and documentation in
`projects/anki_grpc` own the execution changes. `projects/anki_as_code` owns its
entry-point and worker selection changes under a corresponding change record.
`third_party/com_github_ankitects_anki` owns upstream acquisition, guest patches,
and Cargo dependency configuration.

The proposed new Go dependency is `github.com/tetratelabs/wazero` v1.12.0. No
suitable WASM runtime was found in the existing root Go manifests. This runtime
avoids native linking but adds executable size, runtime initialization, memory
transfer overhead, and ongoing security/version maintenance. Specific dependency
approval was granted for wazero v1.12.0 in this session. Additional external
dependencies still require their own specific approval.
