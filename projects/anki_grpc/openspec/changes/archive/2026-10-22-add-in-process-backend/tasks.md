## 1. Native behavior established before the split

- [x] 1.1 Extend native E2E before implementation for persistence, native errors, concurrency, cancellation, shutdown and independent instances.
- [x] 1.2 Share Rust dispatch and expose owned handles and result buffers through the C ABI.

## 2. Independent feature candidate

- [x] 2.1 Preserve the combined source checkpoint, move broad refactoring to PR #186 and retain only native execution plus minimal address integration; verify the diff against master.
- [x] 2.2 Move the native Client to an explicit subpackage and adapt the existing E2E before changing its implementation; verify independent instances, cancellation and persistence artifacts.
- [x] 2.3 Verify all affected project tests, pure generic/application builds, formatting and native lint; record candidate evidence for independent delivery.
- [x] 2.4 Expose the address-only constructor signature as the pure ankigrpc.ClientFactory interface and ClientFactoryFunc adapter.
- [x] 2.5 Correct the native handle ownership comment to name Server.

Failure cases recorded before implementation: native imports leak into core/application packages; moving the adapter changes errors or call options; cancellation frees active Rust state; stopping releases a source reservation before native completion; identical labels share state; splitting imports symbols available only in PR #186; a global rename changes upstream protocol names.

Final candidate validation and publication are owned by repo-delivery receipts and Git.

## Package responsibility correction

- [x] Consolidate native Go integration in ankigrpc/native, keep build composition in the command, and correct address/default ownership descriptions.
- [x] Verify native engine selection, registration, pure-library builds and existing persistence/cleanup artifacts against the final candidate.

Failure cases recorded before this revision: a native import leaks into pure callers; moved native entry points stop linking; command composition introduces address parsing; startup comments promise interruption of synchronous native initialization.

Additional address failure case identified before correction: a bare host:port whose hostname equals a registered scheme is misrouted to that scheme instead of TCP. Query-bearing normalized URLs must continue to reach their registered constructors.

## Fresh review corrections

- [x] Preserve scoped IPv6 TCP addresses through URL normalization and verify a real scoped-loopback connection.
- [x] Keep local/remote classification beside registered constructors; verify a custom owned scheme with a managed API collection.

Failure cases recorded before correction: percent signs in IPv6 zones are reparsed as invalid URL escapes; a registered custom local scheme constructs correctly but managed API access rejects it because another package hard-codes strategy names.

## Simplification pass

- [x] Normalize addresses before a single registry lookup, retain parsed URL data instead of parallel address/scheme results, and remove redundant default factory injection from fixtures.
- [x] Verify existing native, external and malformed-address workflows and inspect persisted artifacts; keep final pure-build and publication evidence in the delivery receipts.

Failure cases considered before this pass: scoped IPv6 zones lose escaping; a hostname equal to a registered scheme changes routing; canonical query-bearing subprocess URLs lose arguments; qualified unavailable engines fall back; normalization depends on which optional adapters are imported; using the default factory changes collection cleanup.

## Final review boundary cases

- [x] Preserve plus signs in scoped IPv6 zones while continuing to reject unavailable qualified engines; verify the endpoint workflow before and after correction.
- [x] Register synchronous startup as active native work outside the lifecycle mutex; verify concurrent Stop deadlines and final cleanup before and after correction.

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

- [x] Return concrete subprocess and network client types and adapt them only at registry boundaries.
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

- [x] Retain merged Client names, pure library compatibility, optional executable composition and existing fixture coverage while replaying native execution onto the fetched master.
- [x] Inspect the complete rebased diff and record its exact validation and delivery plan.

Integration failure cases: obsolete Backend aliases conflict with Client; merged endpoint and reserved-worker-path cases disappear; native error decoding has two owners; a pure executable imports the native engine; source reservations are released before Done. Exact-candidate checks and publication remain owned by delivery receipts.

## Collection scope simplification

- [x] Restore the minimal Start/Connection/Stop client contract by removing ownership methods.
- [x] Remove the one-field network options type and forwarding constructor; register the concrete address constructor directly while preserving endpoint validation.
- [x] Reconcile current contracts and verify native/network E2E, private factories, cancellation, source reservations and pure imports.

Failure cases recorded before implementation: backend ownership becomes a transport capability despite identical RPC support; simplifying constructor adaptation changes parsing or optional imports; removing wrappers erases error causes or registration guards. Keep AsyncClient because its final completion signal protects actual resources.

## Abstraction review corrections

- [x] Extend existing E2E before implementation for unregistered injected factories, independent configured factories, defaults and literal relative Unix socket paths.
- [x] Report backend ownership on clients, remove registration ownership metadata, publish the asynchronous cleanup interface and clarify lifecycle guarantees.
- [x] Add scoped factory configuration and centralize network parsing in construction; migrate executable and library consumers without introducing dependencies.
- [x] Validate affected workflows, pure builds, lint and representative persistence artifacts; review the final source for delivery.

Exact-candidate publication remains owned by repo-delivery receipts and Git.

Failure cases recorded before implementation: managed API calls bypass an injected factory for classification; ownership declarations disagree with constructed clients; two applications cannot register distinct implementations of the same scheme; defaults require repeated wrappers or leak into later explicit selections; Unix paths containing spaces fail URL authority parsing; IPv6 zones are decoded twice; constructor changes leak native dependencies; cleanup completion loses source reservations. Synchronous native startup remains non-interruptible; Stop may finish waiting before asynchronous cleanup completes.
