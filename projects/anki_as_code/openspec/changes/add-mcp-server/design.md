## Context

See proposal.md for scope. The existing App owns every workflow, and
internal/rpc owns compiled service discovery and ProtoJSON invocation. Managed
sources are reserved per App until native cleanup finishes. Writes can persist
after a later operation fails, and archive publication can occur on reads.

## Goals / Non-Goals

**Goals:** Make the existing operations discoverable and callable over MCP
without duplicating collection logic or handwritten native method catalogs.

**Non-Goals:** HTTP hosting, authentication configuration,
transactional rollback, automatic retries, and recovery
features from the separate deferred change.

## Decisions

- Add private `internal/mcp` and a Cobra `mcp serve` command. The entry point
  constructs an MCP Server with App, then passes App and that concrete server
  into CLI. CLI owns a private mcpServer interface and calls Serve; protocol
  tools call App. App has no dependency on the protocol host. The SDK owns stdio
  transport and the command supplies cancellation. One App is shared by tools
  so existing source reservations continue to apply.
- Register workflow names matching command paths with underscores. Keep each
  workflow's schema declaration beside its handler. Native names are
  `api_<service>_<method>` for declared RPCs from the existing catalog. Omit the
  226 implicit forwarded aliases while retaining declared Backend methods and
  existing canonical names. The pinned catalog has 269 native tools and 281
  total tools. Detect collisions before accepting input.
- Following the user's explicit review decision, col_open returns a
  handle and col_close publishes its archive and stops its worker.
  App.ColOpen delegates to Store.Open, which returns an already-open
  ankiclient.Client handle. MCP borrows its Connection and calls Stop; native/workflow
  methods use their existing Connection option.
  There is no worker cache or client lease abstraction. The protocol registry
  only maps handles, serializes each collection's operations, and closes all
  sessions on disconnect. App's single Store retains source/alias exclusion
  across both explicit handles and ordinary path operations.
- Use explicit SDK tool schemas and raw JSON arguments. A single converter
  traverses protobuf descriptors using the codec's canonical protobuf field
  names, `$defs`/`$ref` for nested/recursive messages, arrays, maps, enums,
  string-encoded 64-bit integers, bytes, optional fields and oneof constraints.
  Reuse ProtoJSON as the request validator and response encoder; do not infer
  ProtoJSON from generated Go struct layout. Reject unsupported message
  constructs at startup rather than advertise incomplete schemas.
  Represent each real oneof with one branch per non-null member and one branch
  for absent/null members. Exactly one branch must match, keeping the schema
  linear in the number of fields while preserving ProtoJSON's unset semantics.
- Each native tool takes collection/address/write_back controls and a nested
  `request`. Workflows take explicit paths, state_key, and structured protobuf
  config/state values where needed. Outputs contain structured JSON plus a
  text JSON fallback. Successful write-only operations return an empty object.
- Use SDK `Server.AddTool` so payload bytes reach App without lossy generic
  number conversion. Adapter validation owns envelope keys/types/required
  values; App and ProtoJSON own native validation and lifecycle restrictions.
- Tool failures set `isError` and preserve joined operation/cleanup diagnostics.
  Client cancellation and disconnect cancel work and allow owned cleanup to
  finish before server exit. Stdout carries only MCP frames. Do not classify
  native calls or collection workflows as universally read-only/idempotent.
- Pin the requested SDK at v1.8.0 through root Go module metadata and existing
  Gazelle/Bzlmod mechanisms. Runtime descriptor conversion avoids another
  generator dependency. Individual tools honor the user's preferred discovery
  and approval granularity despite the larger catalog.

## Risks / Trade-offs

- ProtoJSON differs from Go JSON -> cover every pinned descriptor and exercise
  enums, maps, recursive messages, oneofs, bytes and 64-bit identifiers.
- The API is large -> retain deterministic names and SDK tool-list pagination;
  discovery must not start an Anki backend.
- Partial native work can persist -> return native/cleanup errors, document
  persistence, and never automatically retry failed mutating calls.
- Concurrent calls or disconnects can race cleanup -> reuse one App's guards
  and verify cancellation and process termination through a real client.
- The application suite has 48 sequential cases, including 513 MiB media, and
  exceeded its 900-second budget -> use optimized builds and opt-in Bazel
  sharding. Targets retain Bazel's default shard count; callers select their
  own parallelism with `--test_sharding_strategy=forced=N`. Go's runner partitions
  top-level cases across processes, preserving environment and cache isolation.
- Large-fixture preparation previously exported and bound an archive, publishing
  the whole media package repeatedly -> prepare configs/bindings on the source
  database before packaging. Build/import still traverse the complete archive;
  convergence inspection borrows the already-open restored connection. For the
  large-field case, assert convergence from the existing CLI plan output rather
  than issuing a duplicate plan. Payload sizes and substantive assertions remain.
- Performance validation uses the normal Bazel `-c opt` mode. Direct comparison
  reduced the 513 MiB WASM case from 111.01 to 17.37 seconds and the 17 MiB field
  case from 25.30 to 9.08 seconds. Proceed with the existing optimized mode and
  document it as the E2E verification command. A forced optimization transition
  would change normal debug-build semantics and is unnecessary for this goal.
- The SDK's schema inference accepts Go types rather than protobuf descriptors.
  Go `encoding/json` and ProtoJSON differ for 64-bit integers, enum names and
  oneof layout. Keep the small descriptor converter so schemas match the
  existing codec, and use the SDK for transport, registration and schema
  resolution rather than maintaining individual RPC schemas.

## Migration Plan

This is an additive command with no state migration or live deployment. Validate
an installed executable against disposable databases/archives before publishing.
Removing the command and dependency rolls back MCP access without changing
existing config, state or collection formats.

## Simplification decision

App.ColOpen accepts ColOpenOptions containing Path, Address and WriteBack, and
returns an already-open ankiclient.Client with Connection and Stop. This shared
capability replaces the separate application Collection interface. The application
owns one private, single-use collection scope; its raw gRPC transport implements
StartableClient, embedding Client and adding Start. Factories return concrete clients
directly. Collection preparation, archive publication and temporary restores
remain with the application's collection package.

Store owns source reservations and session-context cancellation, observing
private scope completion before releasing a source. Completed reservations are
reaped during the next open so immediate reuse is independent of observer
scheduling. WithCollection uses Open, Connection and deferred Stop.

MCP does not mirror collection shutdown state: serialized handle calls obtain
their connection from the application handle, which rejects access after Stop
begins. Disconnect cleanup runs once after SDK Run drains request handlers.
Independent handles close concurrently under the scope's bounded cleanup policy.

CLI receives the concrete mcp.New(app) value beside App through its private
mcpServer interface, declared in internal/cli/mcp.go with only
Serve(context.Context) error. This consumer-owned contract groups MCP commands
behind one CLI dependency without exporting an interface or placing it in
ankiascodeapi. The concrete server receives App, matching the command adapter's
dependency direction. Hosting is
not an App operation. Retain path and handle modes, source exclusion, independent bounded
cleanup, per-handle serialization and lossless ProtoJSON validation.
