# Validation

## Initial candidate before the API CLI merged

The native CLI E2E fixture was added before production edits. Against the
initial candidate (`b31e3f3bca640c433c3bea9523a0411f75a0e83a` plus the fixture),
`TestExternalServerWorkflow/tcp` failed with `unknown flag: --address`.

The implemented CLI workflow passes for TCP and Unix listeners, including all
state operations, export, plan, apply, build, embedded defaults, malformed
addresses and cancellation while waiting for gRPC readiness. Its outputs include
`external-{tcp,unix}.colpkg`, state JSON and converged plan JSON. Inspected
archives contain `collection.anki21b`; both plans are empty objects and both
state outputs contain the expected existing deck binding. The fixture also
checks a newly created deck through native export, unchanged build sources and
byte-identical archives after converged apply under the earlier publication contract.
The rebased candidate follows the merged write-back contract described below.

Preflight validation passed all 31 selected tests (native application,
collection store, CLI, archive builder and repository quality), semantic lint
of the five affected Bazel packages, and strict OpenSpec validation. The CLI
help exposes `--address` on state commands and describes its embedded default.
The final exact-candidate checks are repeated and recorded by repo-delivery.

## Native listener defect

`linked-listener-sigpipe`: Reusing the linked server initially exposed process
termination by signal 13 when a disconnected socket was written. A direct
30-client probe with stderr redirected to a regular file completed one client
under the default SIGPIPE policy before server termination. The same probe
with inherited SIGPIPE ignore completed all 30 clients with the server alive.
The executable now owns this policy alongside its existing native signal
adapter. The passing multi-command native E2E exercises actual listener reuse.

## Session ergonomics

The session used the answer/workspace/OpenSpec/Go/Bazel/layout/decision and
delivery procedures, live Cordis repository readers and bounded host commands.
No external dependency, provider or host configuration changes were needed.

`oversized-discovery-output`: Three broad discovery/read batches were truncated,
including an unnecessary full worktree inventory. Bounded owner-local searches
and Cordis line-range reads resolved the missing context. Future continuation
should start from this change and the ignored `out/external-grpc/` logs rather
than repeat repository-wide discovery.

The first native baseline build took about two minutes; subsequent CLI E2E
completed in about seven seconds. Fixture failures were corrected without
changing production semantics: the deck needed its normal configuration and
the listener's stderr reader needed to remain active. Final publication gates
and exact candidate identities belong to the delivery receipts.

## Rebase integration

Rebased onto master `4614ffae7412b6c025b96ce6c2f5371697b91938`, including the
native API command and CLI reference generator. Review of `716f10fea` found the
factory constructor incompatible with the generator and stale byte-identity
assertions incompatible with master's archive publication policy. It also found
a convergence assertion checking a nonexistent JSON field.

Before correction, the native API probe rejected `--address` as unknown. The
single-App TCP/Unix/embedded E2E was preserved from the earlier pre-implementation
fixture. Corrective changes remove the factory, pass collection options per call,
and use `--address` on native API commands. Typed protobuf convergence checks
replace map keys and byte-identity assertions. The merged write-back policy and
independent collection-close/transport-close timeouts are preserved.

Raw diagnostics and pending-work preservation are under ignored
`out/external-grpc/rebase-review/`. Publication and validation receipts identify
the final candidate; earlier successful checks above apply only to the earlier
candidate.

Corrected native CLI and reference-generator E2E checks pass. Inspected TCP and
Unix archives contain `collection.anki21b`, state artifacts each contain one deck
binding, and repeated apply/plan results are empty objects. Native API response
artifacts from both transports contain the expected formatted timespan, and the
single-App default-state artifact contains only version 1 after alternating
external and embedded calls. Generated API and collection reference pages both
expose `--address`. Strict OpenSpec validation passes.

`merged-consumer-gap` (fixture-tested): checking the merged CLI reference
consumer exposed the factory signature break. Include the generator in final
validation whenever changing the CLI constructor. The final plan adds it to the
existing native, collection, application and archive-builder gates.
`empty-address-conflict` (fixture-tested): the native E2E caught Cobra treating
an explicitly empty address as a conflicting flag. Value-based validation now
preserves the embedded default and still rejects two nonempty access modes.
The earlier oversized-read issue recurred in large diff batches; continuation
should use the recorded findings and bounded per-package reads.

## Unified collection access

The user requested one concrete collection options value. Before production
changes, native workflow callers were adapted to Path/Address/Connection/StateKey
and failed to compile against `fd9b078377e6e36facadb1c821b6b963aed2c202` because
Path and the consolidated signatures were absent. Diagnostics are under ignored
`out/external-grpc/unified-options/before-implementation.log`.

After implementation, all five selected native/application/CLI/rebuild/reference
workflow targets pass. The inspected `collection-options-validation.json`
artifact records rejection of missing access, typed nil connections, conflicting
path/address/connection fields and borrowed-connection builds. The borrowed
archive still contains native content and the default-state artifact contains
only version 1 after alternating transport/namespace selections. TCP and Unix
archives contain `collection.anki21b` and both resulting plans are empty objects.

CollectionSource and its constructors are removed. Col methods take one
CollectionOptions value, and APIRequest.Collection carries that same value
without a second backend field. Empty path-operation addresses retain embedded
access. Raw API access retains its existing exclusive path/address/connection
modes and borrowed lifetime. Build copies the supplied options and changes only
its staged path. Strict OpenSpec validation passes; delivery receipts identify
the final candidate, mandatory gates and publication state.

`collection-selection-overlap` (fixture-tested): a source plus options split path
and address selection across parameters and silently ignored borrowed-connection
addresses. One validated value removes that ambiguity. Pinned ast-grep matched
42 existing collection calls across seven test files for the migration. A first
relative-path probe was corrected to use absolute paths because Bazel launchers
execute from runfiles. Bounded per-file reads should be preferred over combined
diffs; some broad review reads again exceeded their output budgets.

`unix-fixture-address-length` (live): full delivery validation exposed a Unix
socket fixture whose absolute pathname was 152 bytes under the longer task
TMPDIR. Direct AF_UNIX connection rejected that address as too long. The native
listener already uses a relative path; its absolute-form assertions now use the
Linux `/proc/self/cwd` alias, retaining both pathname and unix:// coverage under
long test roots. This changes only fixture addressing, not production transport.
The original failed validation receipt and diagnostics remain in ignored scratch.
