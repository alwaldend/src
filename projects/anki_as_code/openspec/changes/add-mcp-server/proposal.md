## Why

Agents need structured access to both Anki as Code's declarative workflows and
the native Anki API. MCP can expose the existing application operations and
compiled protobuf catalog while retaining their current storage ownership.

## What Changes

- Add `anki_as_code mcp serve`, using the official Go MCP SDK over stdio.
  CLI receives independent App and MCP interfaces; the MCP server receives App.
- Expose config export/dump, collection export/plan/apply/build, and all four
  state operations as individual tools with structured inputs and outputs.
- Register an individual tool for every declared native RPC, including declared
  Backend methods, with schemas derived from protobuf descriptors. Omit 226
  implicit forwarded aliases: the pinned catalog has 269 native tools and 281
  total tools. Canonical declared-method names remain unchanged.
- Preserve ProtoJSON encoding, per-operation collection access, archive
  publication, cancellation, and operation/cleanup error semantics.
- Add explicit col_open/col_close tools returning opaque handles.
  Handle-backed calls borrow the existing connection; enabled archives publish
  on explicit close or disconnect. Ordinary path calls retain per-call lifetime.
- Verify the installed executable through real MCP client workflows and retain
  repeatable collection-state and tool-catalog evidence.
- Run application and CLI test cases in isolated parallel Bazel shards. Avoid
  redundant archive publication in large-fixture setup and inspection while
  retaining the 17 MiB field and 513 MiB media round trips.

## Capabilities

### New Capabilities

- `mcp-server`: Stdio MCP transport, workflow tools, generated native tools,
  schema discovery, and errors/lifecycle behavior.

### Modified Capabilities

- `collection-source-api`: Add explicit App collection sessions using
  ColOpenOptions and one factory-created client, lending prepared connections
  through existing borrowed-operation contracts.

## Impact

Changes belong to `projects/anki_as_code`, with a private MCP adapter, Cobra
entry point, E2E tests, and documentation. Root Go/Bazel dependency metadata
will pin the user-requested `github.com/modelcontextprotocol/go-sdk` at v1.8.0
and resolve its transitive modules through the owning Go workflow. No separate
protobuf schema-generation dependency, HTTP service, host configuration, or
live collection migration is introduced.
