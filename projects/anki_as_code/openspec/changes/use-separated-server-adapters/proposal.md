## Why

The application must select a server adapter without pulling native cgo linkage
into its WASM build. Preserve its existing same-executable worker lifecycle.

## What Changes

- Build on the native adapter composition merged in #182.
- Select native or WASM hosting through one command package and public Bazel
  target, using independent `//projects/anki_grpc:cgo` and `//projects/anki_grpc:wasm` booleans.
- Keep collection commands and process supervision independent of adapter choice.

## Capabilities

### New Capabilities

- `server-build-selection`: Native and cgo-disabled WASM application builds with
  the same supervised child-process boundary.

### Modified Capabilities

None.

## Impact

Application entry points, Bazel targets, documentation and E2E consumers.
The server and guest implementation are owned by the linked
`projects/anki_grpc/openspec/changes/add-wasm-server` change.
