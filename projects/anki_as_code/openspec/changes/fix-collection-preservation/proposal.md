## Why

Lexical path cleanup can select the wrong collection through symlink-parent
paths, and removing the required default deck can lose its surviving native
binding. Collection ownership also belongs in the application rather than an
expanded generic gRPC client. Export destinations are caller-owned.

## What Changes

- Preserve filesystem traversal order when resolving public collection paths.
- Reject removal of Anki's required default deck while preparing a plan.
- Remove protected-source export parameters and source/output comparisons;
  overwrite declared destinations under caller-owned output policy.
- Return a prepared collection through the shared ankiclient.Client interface,
  which exposes Connection and Stop; use StartableClient for transport startup,
  keeping storage, archive publication and source reservations in the collection
  package instead of generic transport clients.
- Inject constructor functions and keep optional completion observation private.
  Let CLI prefix native serve arguments before calling ServerRun.
- Remove redundant options/format types and test-only helpers; keep resource
  service fields private when only their package uses them, preserving public
  snapshot data and meaningful cleanup configuration.
- Exercise these behaviors through real native, WASM and MCP workflows.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `configuration-filesystem`: Preserve symlink-parent semantics and treat export
  locations and existing contents as caller-owned.
- `collection-planning`: Refuse removal of required native deck identities.
- `application-server`: Keep one ServerRun operation and adapt serve arguments
  in CLI before application-owned native hosting.

## Impact

Changes belong to application path handling, configuration output, collection
scope ownership and planning. The linked add-mcp-server collection-source API
record describes the prepared handle. No dependency, storage-format or live
infrastructure changes are required.
