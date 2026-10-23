## Why

Anki as Code needs one client for native access and collection lifetime. Separate
transport, collection and reservation clients obscure ownership and duplicate the
same lifecycle contract. Factory-created clients can own optional storage state
directly while keeping the transport reusable between collections.

## What Changes

- Give the factory-created Client OpenCollection, CloseCollection and Done beside
  Start, Connection and Stop. Keep storage state private to that same client.
- Adapt process, network and in-process transports through the existing factory.
- Close and publish collection storage without stopping its transport; terminal
  Stop closes storage then stops the transport under independent cleanup budgets.
- **BREAKING:** replace the separate collection client constructors with collection
  methods and CollectionOptions on the existing client. Migrate repository callers.
- Keep source reservations and session cancellation in the application's store,
  observing Client.Done without returning a reservation-client wrapper.

## Capabilities

### New Capabilities

- `owned-collection-client`: One client owning transport and optional collection
  storage, with bounded terminal cleanup and observable native completion.

### Modified Capabilities

- `collection-backend`: Preserve collection storage, publication and native-call
  guarantees through OpenCollection/CloseCollection on a reusable client.

## Impact

Changes belong to pkg/ankiclient, factory adaptation in pkg/ankigrpc, and their
callers, tests and documentation. No dependencies, native protocol changes or
storage migrations are introduced. Application integration is tracked by
projects/anki_as_code/openspec/changes/add-mcp-server.
