## 1. Delegation

- [x] 1.1 Extend the existing endpoint E2E with an explicit subprocess address before implementation; verify its database and state artifacts.
- [x] 1.2 Delegate App and Store construction through one address-only constructor and keep optional native imports and hosting composition in the executable; verify alternating endpoint workflows and pure application builds.
- [x] 1.3 Use ankigrpc.ClientFactory for application and collection-store constructor fields instead of repeating the function signature.
- [x] 1.4 Migrate existing workflow fixtures to the injected factory before removing default executable/address options, update examples and verify native workflows.
- [x] 1.5 Exercise cgo registration through the shared constructor in the endpoint E2E, remove executable constructor-wrapper use and update documentation.

## 2. Delivery

- [x] 2.1 Run affected integration checks and format gates, record the validation evidence and prepare the independent refactoring PR.

Final candidate validation, push and pull-request state are owned by repo-delivery receipts and Git.

## Package responsibility correction

- [x] Consolidate native Go integration in ankigrpc/native, keep build composition in the command, and correct address/default ownership descriptions.
- [x] Verify native engine selection, registration, pure-library builds and existing persistence/cleanup artifacts against the final candidate.

Failure cases recorded before this revision: a native import leaks into pure callers; moved native entry points stop linking; command composition introduces address parsing; startup comments promise interruption of synchronous native initialization.

## Fresh review corrections

- [x] Preserve scoped IPv6 TCP addresses through URL normalization and verify a real scoped-loopback connection.

Failure cases recorded before correction: percent signs in IPv6 zones are reparsed as invalid URL escapes.

## Simplification pass

- [x] Rely on the shared constructor default in workflow fixtures and remove redundant imports/dependencies; keep routing simplification in the paired backend change.
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

- [x] Remove application/store default worker settings; supply the default through the injected client factory and centralize addressed API cleanup in the client library.
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
