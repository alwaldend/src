## Why

The server already exposes Anki's protobuf operations. A generated Cobra command
tree makes those operations usable without writing a Go client or maintaining
handwritten command wrappers.

## What Changes

- Add `anki_rpc`, with one service command and one subcommand per declared RPC.
- Generate Go registration code from the pinned upstream descriptor set.
- Accept ProtoJSON inline, from a file, or stdin and emit ProtoJSON responses.
- Reuse the managed stdio backend, with optional collection opening for one-shot
  operations, or connect to an existing TCP/Unix gRPC server.
- Deliver independently against master, without archive or #128 changes.

## Capabilities

### New Capabilities

- `protobuf-cli`: Schema-generated Cobra commands for native Anki operations.

### Modified Capabilities

None.

## Impact

Changes belong to projects/anki_grpc with a narrow third-party schema visibility
grant. Existing Cobra, gRPC, protobuf, runfiles, and backend dependencies suffice.
No new dependency, copied schema, or per-operation business wrapper is added.
