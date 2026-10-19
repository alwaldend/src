## Why

The Anki WASM server needs reproducible WASI adaptations of its already pinned dependencies. Build inputs and native composition landed separately in #182; this change owns only guest-specific dependency patches and annotations.

## What Changes

- Register the Anki and FSRS WASI patches and the Chrono host-timezone patch.
- Scope Tokio's existing unstable WASI opt-in to the WASI target.
- Regenerate module metadata without changing dependency versions or repositories.

## Capabilities

No separate behavioral delta: guest behavior belongs to projects/anki_grpc/add-wasm-server. `skip_specs: true`.

## Impact

The Anki third-party package owns patches and crate annotations; MODULE.bazel.lock is their generated projection. Existing approved SDK/runtime pins remain owned by the merged build-input change.
