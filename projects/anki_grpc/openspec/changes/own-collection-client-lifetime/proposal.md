## Why

Transport clients already own startup, RPC access and shutdown. Adding reusable
collection phases duplicated that lifecycle and coupled the gRPC library to
application archive and storage policy. The user requested restoration of the
three-method transport contract, with shared connection access and shutdown
available separately from startup ownership.

## What Changes

- Define Client as Connection and Stop. StartableClient embeds Client and adds
  Start; factories return concrete StartableClient implementations directly.
- Remove the generic client wrapper, Transport duplication, collection methods,
  CollectionOptions and managed storage implementation from the gRPC library.
- Retain concrete Done methods for native work that can outlive a cancelled
  Stop wait, with private completion observation in the application.
- Express ClientFactory as a constructor function and adapt concrete constructors
  through ToClientFactory; retain the concrete Factory's registry and defaults.
- Move single-use collection ownership to Anki as Code's existing collection
  package and expose its prepared handle through the shared Client interface.
- Preserve transport error handling and native backend destruction outside the
  asynchronous serving runtime.
- Limit handwritten exports to real consumers and required ABI/framework
  boundaries; hide package-only C bridge and module-only Rust CLI declarations.

## Capabilities

### New Capabilities

- `owned-collection-client`: A small transport client contract with optional
  final native completion; applications own collection storage.

### Modified Capabilities

- `collection-backend`: Remove library-owned storage phases; callers use native
  RPCs or an application-owned collection scope.
- `client-connection-platforms`: Use constructor-function injection and typed
  adaptation while retaining registry isolation, defaults and raw-scope cleanup.

## Impact

Changes affect pkg/ankiclient, factory adaptation and repository consumers.
Anki as Code owns archive publication, temporary restores and source reservations.
No dependencies, protocol changes or storage migrations are introduced.
