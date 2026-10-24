## Why

The browser Go and Rust consumers need the same caller-supplied JavaScript
filesystem. The user selected Emscripten classic MEMFS as the backing store;
its FS API needs a small Node API translation layer.

## What Changes

- Expose an asynchronous factory returning one independent filesystem with
  callback and synchronous operations over the generated MEMFS module.
- Preserve upstream storage, streams and inode metadata; report unsupported
  operations with conventional errors.
- Add a standalone provider E2E with retained JSON and binary artifacts.

## Capabilities

### New Capabilities

- `emscripten-filesystem-adapter`: Node API translation over classic MEMFS.

### Modified Capabilities

None.

## Impact

Changes are scoped to the browser adapter, its files target and provider E2E.
The generated Emscripten dependency is owned by a concurrent build change.
Go/Rust hosting, browser integration and durable persistence remain with
their existing owners. No npm provider or polyfill dependency is added.
