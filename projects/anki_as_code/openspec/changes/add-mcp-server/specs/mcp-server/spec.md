## Purpose

Expose Anki as Code workflows and the supported native Anki API to MCP clients
with schema-backed discovery and the existing collection persistence guarantees.

## ADDED Requirements

### Requirement: Stdio MCP command

The executable SHALL provide `mcp serve` over stdin/stdout. Protocol discovery
SHALL work without an Anki backend or collection. Stdout SHALL contain only MCP
messages; command errors SHALL use stderr and a nonzero exit status. CLI SHALL
receive independent App and MCP interfaces. The MCP server SHALL receive App for
its tools, and App SHALL NOT depend on MCP hosting.

#### Scenario: Discover an installed server

- **WHEN** a client starts an installed executable with `mcp serve`
- **THEN** it can connect and enumerate tools without opening collection storage
- **AND** the executable does not need PATH or Bazel runfiles for bundled access

### Requirement: Individual native tools

The server SHALL advertise one deterministic `api_<service>_<method>` tool for
every declared compiled RPC, including methods declared on Backend services.
Implicit forwarded aliases SHALL NOT create duplicate tools. Existing canonical
declared-method names SHALL remain unchanged. Request and response schemas SHALL
derive from the pinned protobuf contracts and describe
nested messages, enums, repeated fields, maps, oneofs, bytes and 64-bit values
using the application's ProtoJSON conventions. Invalid requests SHALL fail
before native execution. Unsupported schema constructs SHALL fail startup.

#### Scenario: Discover declared methods without forwarded duplicates

- **WHEN** a client enumerates the native tools for the pinned Anki descriptors
- **THEN** all 269 declared RPCs are available with their canonical names
- **AND** 226 implicit forwarded aliases are omitted, yielding 281 total tools
  with the ten workflows and two session controls

#### Scenario: Call a native method

- **WHEN** a client calls an advertised native tool with a collection or address
  and a valid structured request
- **THEN** the response contains structured native JSON and a text fallback
- **AND** backend capability and collection-lifetime restrictions remain enforced

#### Scenario: Reject malformed native input

- **WHEN** a request has unknown fields, an invalid enum, conflicting oneof
  members, or invalid scalar values
- **THEN** the tool returns an error before invoking the native method
- **AND** the MCP session remains usable

### Requirement: Workflow tools

The server SHALL expose config dump/export, collection export/plan/apply/build,
and state dump/import/update/replace as individual tools. Collection options
SHALL remain explicit per operation, and state namespaces SHALL retain their
current default. Tools SHALL return structured protobuf JSON or an empty success
object for operations without a returned value.

#### Scenario: Reconcile through MCP

- **WHEN** a client exports configs, plans an edit, applies it and plans again
- **THEN** the resulting collection and stored bindings agree with the configs
- **AND** native identities, media and review history retain existing guarantees

#### Scenario: Build an independent archive

- **WHEN** a client builds a new archive from an existing source and configs
- **THEN** the output can be read through the native API
- **AND** the base source and configuration files remain unchanged

### Requirement: Existing persistence and cleanup semantics

Path-based calls SHALL retain current source reservation, native lifetime, archive
publication, and partial-write behavior. Tool errors SHALL preserve operation
and cleanup diagnostics without implying rollback or performing automatic
retries. Cancellation and client disconnect SHALL cancel request work and allow
bounded native cleanup before process exit. Tool annotations SHALL not promise
universal read-only or idempotent behavior for collection/native operations.

#### Scenario: Cancel an active request

- **WHEN** a client cancels or disconnects during a collection operation
- **THEN** owned native resources are cleaned up under the existing bounded policy
- **AND** subsequent collection access is possible after cleanup completes

#### Scenario: Native operation fails

- **WHEN** a valid native request fails during execution or cleanup
- **THEN** the tool reports failure with the available diagnostics
- **AND** later requests can continue without an automatic mutation retry

### Requirement: Explicit collection sessions

The server SHALL expose col_open and col_close. Open SHALL return
an opaque handle for a prepared App collection session. Native and collection
workflow tools SHALL accept that handle instead of path/address controls and
pass its connection to existing App methods. Calls on one handle SHALL not
interleave complete collection workflows. Ordinary path calls SHALL retain
per-call lifetime; col_build SHALL continue to require a base archive path.
Open sources SHALL remain reserved against other managed accesses. Enabled
archives SHALL publish on explicit close or disconnect, and their workers SHALL
remain reusable until that session ends. Database writes remain immediate.

#### Scenario: Reuse and publish an open archive

- **WHEN** a client opens an archive, calls native/workflow tools with its handle,
  and closes it
- **THEN** the tools reuse its prepared connection and worker
- **AND** the archive reflects the changes after close, while it remains
  unpublished before close

#### Scenario: Disconnect with open handles

- **WHEN** an MCP connection ends with handles still open
- **THEN** each session closes under App's existing bounded cleanup policy
- **AND** enabled archives publish and their workers stop
