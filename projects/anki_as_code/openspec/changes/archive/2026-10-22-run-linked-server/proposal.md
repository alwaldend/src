## Why

A linked Rust server lets the CLI run its own worker from one installed executable without writing an executable to temporary storage.

## What Changes

- Add anki_as_code server run as a generic passthrough to the linked Rust CLI, including its serve subcommand; provide server serve as a delegating shortcut.
- Configure offline CLI scopes with the exact current executable and server serve arguments; require explicit worker settings for library collection operations.
- Remove production server discovery and the implicit standalone-worker runtime dependency.
- Preserve separate process lifetime, gRPC framing, and existing collection behavior.

## Capabilities

### New Capabilities

- `linked-worker`: Run offline operations and an explicit server command from one linked application executable.

### Modified Capabilities

None.

## Impact

CLI entry point, explicit application/collection worker configuration, packaging documentation and native E2E. This PR targets master independently of the API command change. No dependency additions, temporary executables or live deployment changes.
