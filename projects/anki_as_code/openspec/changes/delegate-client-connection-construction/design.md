## Scope and dependency boundary

App coordinates operations and validates sources. It receives one ClientFactory dependency from the executable and supplies the same factory to raw client scopes and collection Store. Factory default selection and worker configuration are not App or Store settings. The paired backend change owns generic adapters, constructor registration and connection-string semantics.

## Raw clients and managed collections

App prepares RPCs and delegates address-only access to ankigrpc.WithClient(ctx, factory, address, callback). That library helper constructs, starts and lends the client connection, then performs bounded cleanup independently of request cancellation. It reserves or opens no collection. ClientFactory remains a construction contract; injected objects need no lifecycle methods.

Store owns WithCollection, source alias reservations, collection preparation and archive publication. It constructs clients through the injected factory, using an explicit address or passing the empty operation address unchanged for factory default selection. Without a configured factory, path-only access fails before client startup; explicit registered endpoints remain usable. Collection publication/close and transport shutdown retain separate budgets, and operation and cleanup errors remain discoverable.

## Executable composition

The executable owns subprocess schemes, its entrypoint and server arguments. It registers worker constructors explicitly through the common RegisterClientFactory API and generic ClientFactoryFunc[T] adapter. Build-selected command sources choose available worker schemes and native hosting. The configured factory selects its same-executable worker URL only for empty operation input. Neither generic nor native library imports advertise a subprocess entrypoint. Fixture owners construct their declared default workers directly or explicitly register packaged worker schemes.

The executable also injects ServerRunner and owns process-global signals. App and Store import no optional native implementation. Configuration, state, resource, export, planning, reconciliation and rebuilding keep their existing domain responsibilities. Both PRs target master independently.

Direct in-process execution and its API/path source rules remain in independent PR #184. This change owns Client naming and optional executable build composition while preserving collection publication and borrowed connection lifetime.
