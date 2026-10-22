## Why

Collection commands always launch the embedded server even when an Anki gRPC
server is already available. Allow explicit server selection while preserving
the embedded server as the default.

## What Changes

- Add optional `col --address ADDRESS`, inherited by collection and state commands.
- Use the same `--address` flag for the merged native `api <service> <method>` command.
- Inject one App into CLI and pass one CollectionOptions value containing path,
  address, borrowed connection and namespace per operation.
- Reuse the existing TCP/Unix gRPC client for path-based collection operations,
  including builds, without starting or stopping the external server process.
- Preserve collection preparation, cleanup, archive publication and source flags.
- Document the shared filesystem and exclusive server-access requirements.

## Capabilities

### New Capabilities

- `external-server-selection`: Optional CLI and application backend address with
  embedded defaults and unchanged managed collection behavior.

### Modified Capabilities

None.

## Impact

Changes are confined to the Anki as code application, CLI composition,
collection store, documentation and native CLI E2E coverage. Collection workflows and APIRequest use one concrete `CollectionOptions` value;
the source interface and constructors are removed. Native RPC schemas remain unchanged.
No new dependencies, server protocol changes or host configuration are required.
