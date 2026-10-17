## Context

The main feature branch has an App-private registry, reference selectors and prepared collection entries. Master has the native lifecycle and merged resource, state and export packages, with CLI export and reconciliation still implemented on separate branches.

## Goals / Non-Goals

**Goals:** Keep source reservations and archive publication in internal/collection.Store. Delegate native lifetime to backend.ManagerCollection and lend its prepared connection directly. Preserve application workflows and publish a master-based prerequisite.

**Non-Goals:** Implement reconciliation or placeholder CLI commands in this PR. Coordinate external writers. Add dependencies, SQL, manual archive handling or a second native lifecycle implementation.

## Decisions

Store owns source reservations and delegates native preparation and cleanup to backend.ManagerCollection. WithCollection takes a source path and injects grpc.ClientConnInterface into a callback. A private map records canonical source paths and file identities until native cleanup finishes. Duplicate aliases fail, and native work runs outside the reservation mutex. Callers complete operations sequentially before returning; retained connections reject calls after cleanup.

Decision: remove Collection and its service caches. Resource/state constructors create lightweight client wrappers around the borrowed connection. Workflows construct the services they need, and multi-operation callers retain those instances as local variables. External writers and background work beyond callback completion remain caller-owned.

Store publishes changed packages back to their source archive after successful scope completion. Callers export other destinations directly through Anki as their final operation and return false to skip source publication. The store has no destination selection or output reservation API. Database mutations persist immediately, even before later errors.

Store.BackendManager optionally selects a transport manager before use. A nil value resolves the bundled server or PATH and creates ManagerSTD configuration. Each scope composes ManagerCollection over that transport manager and obtains its prepared connection. Source reservations last until the manager has completed cleanup. Package publication remains inside the prepared callback, after the operation and cancellation checks and before native shutdown.

App exposes typed ColStateDump, ColStateImport, ColStateUpdate and ColStateReplace methods with explicit source paths and namespace keys. It uses collection.Store internally and does not expose WithCollection or raw gRPC callbacks. Import loads configs in App; each state workflow constructs its store on the injected connection and returns a publication decision internally. CLI flag bindings are private command data; the CLI parses arguments, decodes/encodes state objects and routes output. Keep existing command spellings, status text and formats. Callers needing a custom multi-operation scope use collection.Store directly. Keep state and export implementations independent. Reconciliation remains in #128.

## Risks / Trade-offs

Connections and services must be used inside a WithCollection callback; retained clients do not extend collection lifetime. Callers assemble local services explicitly instead of receiving a shared container. Callbacks report whether the source archive needs publication, so no-ops and explicit exports retain its bytes. Native cancellation and cleanup retain the backend library's existing guarantees. The unreleased App callback API is removed in favor of explicit workflows; CLI behavior remains compatible. Pending CLI and reconciliation PRs can adopt the extracted package without copying their implementation into this prerequisite.
