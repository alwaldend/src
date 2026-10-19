## Why

Provide approved, reproducible WASI build inputs independently from the Anki WASM server.

## What Changes

- Pin the approved wazero v1.12.0 dependency and WASI SDK 34.0 sysroot.
- Reuse the pinned LLVM compiler and enable the WASI crate-universe target.
- Regenerate the module lock through its owning workflow.

## Capabilities

No specification changes: this extracts native composition and build preparation while preserving runtime behavior. `skip_specs: true`.

## Impact

Root module metadata, Go dependency metadata, LLVM sysroot wiring and Anki crate-universe platform metadata. No upstream Anki port patches or WASM server are included.
