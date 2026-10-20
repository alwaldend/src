## 1. Behavior and implementation

- [x] 1.1 Record failure cases and adapt the existing native E2E to address-only construction before implementation; verify persistence artifacts remain produced.
- [x] 1.2 Implement Client naming, pure address construction and explicit optional subpackages; verify the native E2E and pure builds with cgo enabled.
- [x] 1.3 Expose the address-only constructor signature as the pure ankigrpc.ClientFactory interface and ClientFactoryFunc adapter.
- [x] 1.4 Remove deprecated Backend compatibility aliases and their build declaration; verify all current consumers use Client types.
- [x] 1.5 Address review comments: name the stdio constructor NewSTD, dispatch parsed URL schemes through registered ClientFactory values, let optional adapters register on import, and keep the C ABI bridge and native hosting together in ankigrpc/native.

## 2. Delivery

- [x] 2.1 Verify enabled and disabled executable compositions, project checks and formatting; prepare the refactor independently against master and retain validation evidence.

## Failure cases recorded before implementation

- Generic or application imports accidentally acquire native dependencies when cgo is enabled.
- Address encoding changes an executable path containing spaces or reserved characters, or loses argument boundaries.
- An unavailable explicit scheme silently falls back to subprocess or another strategy.
- External client shutdown terminates the server or prevents a subsequent connection.
- Empty per-operation addresses retain a previous endpoint instead of the configured default.
- A cgo-disabled executable still links the native server or cannot use an external server.
- Client renaming accidentally changes upstream protobuf Backend service identifiers.
- Optional adapter registration is missing, overwrites an existing scheme, or executes a factory while holding the registry lock.
- Malformed URL escapes are ignored, or a bare TCP address / absolute Unix socket path loses its existing interpretation.
- Removing constructor wrappers leaves the executable or a workflow fixture calling a deleted API.

Final candidate validation, push and pull-request state are owned by repo-delivery receipts and Git.

## Package responsibility correction

- [x] Consolidate native Go integration in ankigrpc/native, keep build composition in the command, and correct address/default ownership descriptions.
- [x] Verify native engine selection, registration, pure-library builds and existing persistence/cleanup artifacts against the final candidate.

Failure cases recorded before this revision: a native import leaks into pure callers; moved native entry points stop linking; command composition introduces address parsing; startup comments promise interruption of synchronous native initialization.

Additional address failure case identified before correction: a bare host:port whose hostname equals a registered scheme is misrouted to that scheme instead of TCP. Query-bearing normalized URLs must continue to reach their registered constructors.

## Fresh review corrections

- [x] Preserve scoped IPv6 TCP addresses through URL normalization and verify a real scoped-loopback connection.

Failure cases recorded before correction: percent signs in IPv6 zones are reparsed as invalid URL escapes.

## Simplification pass

- [x] Normalize addresses before a single registry lookup, retain parsed URL data instead of parallel address/scheme results, and remove redundant default factory injection from fixtures.
- [x] Verify existing native, external and malformed-address workflows and inspect persisted artifacts; keep final pure-build and publication evidence in the delivery receipts.

Failure cases considered before this pass: scoped IPv6 zones lose escaping; a hostname equal to a registered scheme changes routing; canonical query-bearing subprocess URLs lose arguments; qualified unavailable engines fall back; normalization depends on which optional adapters are imported; using the default factory changes collection cleanup.

## Final review boundary cases

- [x] Preserve plus signs in scoped IPv6 zones while continuing to reject unavailable qualified engines; verify the endpoint workflow before and after correction.

Failure case recorded before correction: a plus sign in an IPv6 interface zone is confused with an engine qualifier.

## Native naming and exact subprocess arguments

- [x] Extend existing worker E2E before implementation to require zero URL arguments to remain zero, keep the default command explicit, and reject Go engine flags in Rust.
- [x] Rename the optional integration package to ankigrpc/native; remove Rust engine selection, injected worker flags and the no_args URL parameter.
- [x] Reconcile current documentation and specifications and retain the delivery plan for exact-candidate native/pure validation and independent publication.

Failure cases recorded before correction: a URL with no arg values silently gains serve; a qualified URL gains an engine flag; a Go binding choice leaks into the Rust CLI; removing implicit arguments breaks the empty-address default; moving the native package leaks optional dependencies or leaves stale import/build labels.

## Explicit scoped TCP URLs

- [x] Preserve previously accepted raw IPv6 zones in explicit tcp URLs without double-escaping already encoded zones; verify both endpoint forms before and after correction.

Failure cases recorded before correction: tcp://[::1%lo]:PORT fails URL parsing although the legacy network client accepted it; blindly escaping the prefix form changes an already encoded %25 zone.

## Client factory interface

- [x] Replace the function alias with the ClientFactory interface and ClientFactoryFunc adapter; remove the specialized subprocess registration helper, migrate registry entries and application/store construction, and verify stateful factory injection through the existing endpoint E2E.

Failure cases recorded before implementation: concrete factory objects cannot be injected; plain functions no longer satisfy the interface unless adapted; application or collection construction bypasses the injected factory; registry dispatch calls a nil function; error wrapping loses the original cause.

## Reserved characters in scoped IPv6 zones

- [x] Preserve query and fragment delimiters inside bracketed zones for bare, raw TCP URL and escaped TCP URL addresses; reject those delimiters outside the brackets and verify persisted endpoint artifacts before independent delivery.

Failure cases recorded before correction: a whole-address delimiter guard rejects valid zone characters; a fallback treats a query or fragment following the port as zone data; raw normalization double-escapes a valid encoded zone; bracketed plus signs are confused with engine qualifiers.

## Constructor and application ownership review

- [x] Return concrete subprocess and network client types and adapt them only at registry boundaries.
- [x] Verify native endpoint, API and collection lifecycle workflows and retained persistence artifacts before independent publication.

Failure cases recorded before implementation: an empty per-operation address bypasses the configured factory; an explicit address is replaced by the default; address-only API access bypasses injection or leaks a client after failure; removing worker settings changes cancellation, cleanup or source reservation behavior; concrete constructor results no longer adapt to registered factories.

## Addressed client scope ownership

- [x] Move the raw addressed-client scope into ankigrpc beside ClientFactory; retain only collection preparation, publication and reservations in Store.
- [x] Verify factory injection, addressed API startup/cleanup, native collection workflows and pure builds before independent delivery.

Failure cases recorded before implementation: moving a raw scope into Store gives collection code responsibility for requests with no collection; adding lifecycle methods to ClientFactory unnecessarily burdens injected constructor objects; library scoping loses cancellation causes, cleanup on failed startup or callback errors, bounded shutdown or the configured factory.

## Executable-owned subprocess registration

- [x] Register subprocess schemes explicitly in the executable that owns its worker entrypoint; keep only unregistered subprocess construction in ankigrpc and keep native CLI hosting separate.
- [x] Reject empty library addresses and subprocess URLs without an executable, preserve exact argv and configured factory defaults, and verify native and engine-free executable composition.

Failure cases recorded before implementation: importing native advertises a worker entrypoint it does not own; an empty address discovers an unconfigured executable through PATH; owner registration is missing in the installed command; library fixtures rely on incidental registration by another test; zero worker args gain an implicit serve command.

## Generic constructor adapters

- [x] Parameterize ClientFactoryFunc by its concrete client result, remove interface-conversion wrappers, and retain interface returns only at runtime dispatch and consumer boundaries.
- [x] Verify concrete subprocess/network adapters, object factory injection, nil-function validation and native persistence with the existing E2E suites.

Failure cases recorded before implementation: concrete-returning functions cannot satisfy the adapter; a generic adapter no longer implements ClientFactory; nil instantiations evade registration checks; runtime scheme dispatch is incorrectly constrained to a caller-selected concrete type.
