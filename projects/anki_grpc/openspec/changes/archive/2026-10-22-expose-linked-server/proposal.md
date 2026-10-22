## Why

Go applications currently need a separate executable or extraction to run the Rust server. A linked server entry point makes a single installed application executable possible without duplicating protobuf APIs.

## What Changes

- Extract the current Rust server and transport parser into one reusable library, retaining the standalone anki_grpc CLI.
- Expose one argc/argv C ABI entry point and a Go wrapper linked through existing Bazel Rust and cgo rules.
- Remove the client library's implicit standalone-worker runtime input; native tests declare their fixtures explicitly.
- Add the Rust serve subcommand with root/serve help and preserve stdio, TCP, Unix, native error trailers and process failure behavior.

## Capabilities

### New Capabilities

- `linked-server`: Run the existing server through a linked Go entry point.

### Modified Capabilities

None.

## Impact

Rust server source ownership, build declarations, native build consumer visibility, and a small Go library. Reuse existing dependencies and toolchains; no dependency additions or schema changes.
