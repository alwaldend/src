## Why

The Go CLI needs a WebAssembly build independent of the bundled Rust servers.
Its unconditional executable-path lookup currently prevents startup under WASI.

## What Changes

- Add `//projects/anki_as_code:wasm` to select a server-free `wasip1/wasm` variant through the public target's `select()`.
- Keep platform and server-engine selection local to the selected CLI; mixed builds retain their other targets' normal configuration.
- Resolve the application executable only when a bundled server needs it.
- Build schema-generation inputs for the execution platform during cross-compilation.
- Verify the actual WASI command with file-only configuration workflows and document runtime limitations.

## Capabilities

### New Capabilities

- `wasi-cli`: Build and execute the server-free Go CLI under WASI Preview 1.

### Modified Capabilities

None.

## Impact

Project build targets, CLI startup and schema generation, project documentation, and WASI
workflow validation using the existing pinned wazero runtime. No new dependencies.
