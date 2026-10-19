## Backend selection diagram

- [x] Add the CLI, App, Factory, Client and Server diagram with connection selection, lifecycle, gRPC and executable-owned registration edges.
- [x] Render and inspect a dark WebP preview before publication, keeping publication pending user review.
- [x] Add ankigrpc's TCP/Unix registration edge, route gRPC from App to Server and refresh the dark WebP preview.
- [x] Position ankigrpc to the right of Factory and inspect the refreshed dark WebP.
- [x] Split TCP and Unix registration into separate edges and inspect the refreshed dark WebP.
- [x] Format every registration label with the action on the first line and address on the second; inspect the refreshed dark WebP.
- [x] Qualify Factory and Client with their Go package names, verify TCP/Unix registration ownership and inspect the refreshed dark WebP.

The reviewed dark WebP preview is retained under `out/backend-selection/`.
The diagram records registration ownership, client selection and lifecycle
access; its maintained source and SVG render are packaged with the project docs.

## Collection scope simplification

- [x] Extend existing network API E2E before implementation for explicit managed paths, persistence, collection close and reuse of the still-running server.
- [x] Remove backend ownership introspection and its managed API gate; retain exclusive borrowed connections and managed collection lifecycle restrictions.
- [x] Remove discarded changed results from collection/application scope callbacks without changing publication or meaningful domain change results.
- [x] Reconcile documentation and specifications, review the simplification, and validate affected E2E and pure builds before delivery.

Failure cases recorded before implementation: a managed API path with TCP/Unix is rejected; a completed scope leaves the collection open or stops the shared server; borrowed connections accept extra sources; removing a discarded change result suppresses archive writes after callback errors or cancellation; a meaningful fixture export predicate or domain change result is removed accidentally. Explicit Path selects managed collection lifetime; Address selects its client, regardless of transport.

## Abstraction review corrections

- [x] Extend the API and endpoint E2E before implementation to cover an injected unregistered scheme and relative Unix socket paths containing spaces.
- [x] Validate managed API ownership using the constructed client through Store, preserve network-source rejection before startup and use the canonical asynchronous cleanup contract.
- [x] Replace repeated default-selection wrappers with configured factories, update documentation/specifications and verify native, network and pure execution.

Failure cases recorded before implementation: global classification rejects a valid injected factory; an injected factory maps a network-looking address to a backend with different ownership; a network-backed managed API request starts the server connection before rejection; collection scopes lose factory input preservation, cancellation causes, publication or source reservations.

## 1. Established application behavior

- [x] 1.1 Extend application and CLI E2E before implementation for in-process selection, managed API calls and archive publication.

## 2. Independent integration

- [x] 2.1 Delegate through the opt-in address-only constructor, preserve existing lifecycle names and retain asynchronous source reservations; verify the scoped diff against master.
- [x] 2.2 Adapt existing native workflows to explicit optional imports before implementation and verify direct collection, API and publication artifacts.
- [x] 2.3 Verify project tests, pure application builds, generated CLI reference and formatting; record evidence for independent delivery.
- [x] 2.4 Use ankigrpc.ClientFactory for application and collection-store constructor fields instead of repeating the function signature.
- [x] 2.5 Preserve the configured subprocess shutdown timeout when converting existing worker settings into an address, and verify rejection of an invalid policy through the native scope workflow.

Final candidate validation and publication are owned by the paired backend change's repo-delivery receipts and Git.

## Package responsibility correction

- [x] Consolidate native Go integration in ankigrpc/native, keep build composition in the command, and correct address/default ownership descriptions.
- [x] Verify native engine selection, registration, pure-library builds and existing persistence/cleanup artifacts against the final candidate.

Failure cases recorded before this revision: a native import leaks into pure callers; moved native entry points stop linking; command composition introduces address parsing; startup comments promise interruption of synchronous native initialization.

## Fresh review corrections

- [x] Preserve scoped IPv6 TCP addresses through URL normalization and verify a real scoped-loopback connection.
- [x] Keep local/remote classification beside registered constructors; verify a custom owned scheme with a managed API collection.

Failure cases recorded before correction: percent signs in IPv6 zones are reparsed as invalid URL escapes; a registered custom local scheme constructs correctly but managed API access rejects it because another package hard-codes strategy names.

## Simplification pass

- [x] Rely on the shared constructor default in workflow fixtures and remove redundant imports/dependencies; keep routing simplification in the paired backend change.
- [x] Verify existing native, external and malformed-address workflows and inspect persisted artifacts; keep final pure-build and publication evidence in the delivery receipts.

Failure cases considered before this pass: scoped IPv6 zones lose escaping; a hostname equal to a registered scheme changes routing; canonical query-bearing subprocess URLs lose arguments; qualified unavailable engines fall back; normalization depends on which optional adapters are imported; using the default factory changes collection cleanup.

## Final review boundary cases

- [x] Preserve plus signs in scoped IPv6 zones while continuing to reject unavailable qualified engines; verify the endpoint workflow before and after correction.

Failure cases recorded before correction: a plus sign in an IPv6 interface zone is confused with an engine qualifier; Stop cannot observe its deadline while startup holds the lifecycle mutex; interrupted startup publishes a connection after terminal shutdown or signals Done before native cleanup.

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

- [x] Register subprocess schemes explicitly in the executable that owns its worker entrypoint; keep only unregistered subprocess construction in ankigrpc and direct in-process registration in native.
- [x] Reject empty library addresses and subprocess URLs without an executable, preserve exact argv and configured factory defaults, and verify native and engine-free executable composition.

Failure cases recorded before implementation: importing native advertises a worker entrypoint it does not own; an empty address discovers an unconfigured executable through PATH; owner registration is missing in the installed command; library fixtures rely on incidental registration by another test; zero worker args gain an implicit serve command.

## Generic constructor adapters

- [x] Parameterize ClientFactoryFunc by its concrete client result, remove interface-conversion wrappers, and retain interface returns only at runtime dispatch and consumer boundaries.
- [x] Verify concrete subprocess/network adapters, object factory injection, nil-function validation and native persistence with the existing E2E suites.

Failure cases recorded before implementation: concrete-returning functions cannot satisfy the adapter; a generic adapter no longer implements ClientFactory; nil instantiations evade registration checks; runtime scheme dispatch is incorrectly constrained to a caller-selected concrete type.

## Rebase after PR #186

- [x] Preserve merged executable composition and client-factory injection while integrating local API selection and Done-aware collection reservations.
- [x] Inspect the rebased application diff and include native, network, archive and pure-command coverage in its delivery plan.
