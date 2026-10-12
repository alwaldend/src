## Why

PR #146 combines the Rust Anki gRPC project with its Go consumer. Extract the
server into an independently buildable PR against master so it can be reviewed
and merged separately.

## What Changes

- Extract projects/anki_grpc from the validated PR #146 snapshot.
- Include third-party Anki native build definitions required to build the
  server from master; remove the superseded shared dependency-provider rule.
- Preserve upstream command messages and approved dependency pins.
- Import the pinned upstream Cargo workspace directly and replace the custom
  dependency-provider adapter with generated crate_universe dependency helpers.
  Additional workspace-only dependency versions require exact approval before
  declaration or acquisition.
- Add CLI-selected stdio, TCP, and Unix-socket transports, with stdio as default.
- Initialize Anki at server startup and remove the custom Initialize RPC.
- Keep upstream Anki compilation in third_party; the project builds its server
  and depends on the third-party native library.
- Keep the Go API, its E2E tests, and the integration-approaches diagram in #146.
- Register the project landing page and documentation in the main site.
- Update #146 to identify this prerequisite and remove shared files after merge.

## Capabilities

### New Capabilities

- `anki-server-transports`: Select and validate server transports through CLI
  parameters while preserving the default Go consumer's stdio protocol.

### Modified Capabilities

None.

## Impact

projects/anki_grpc, third_party/com_github_ankitects_anki, the generated module
lock, and removal of tools/rust. Both PRs target master. The approved backend,
server, and additional workspace dependency versions are reused. Upstream
source and Cargo manifests are not patched.
