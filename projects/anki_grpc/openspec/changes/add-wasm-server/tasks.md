## 1. Package and build boundaries

- [x] 1.1 Record the decision, dependency approval, and build/protocol/lifecycle/collection failure scenarios before implementation; verify the artifacts describe the preserved process boundary and cgo-disabled build requirement.
- [x] 1.2 Move the native Go/C ABI host to `pkg/ankiservercgo`, update its application consumer and documentation under its owner-local change, and verify existing native worker/CLI E2E checks.

## 2. WASM execution

- [x] 2.1 Add approved wazero v1.12.0 through the existing Go dependency workflow and verify the pinned module metadata and Bazel dependency target.
- [x] 2.2 Specify host-memory and guest-failure behavioral cases before writing the adapter; implement the separate WASM adapter and verify those cases through an executed guest fixture.
- [x] 2.3 Build the pinned Anki backend for WASI with descriptor-owned dispatch and explicit unsupported network-service handlers under the approved offline-first scope; verify the guest artifact and preserve native behavior.
- [x] 2.4 Embed the guest and serve its protobuf dispatch through the existing worker protocol; verify CLI/stdin/stdout/TCP/Unix behavior, native error trailers and cancellation against the exact worker.

## 3. Application and acceptance

- [x] 3.1 Wire a distinct cgo-disabled application target to the WASM adapter and verify it starts the same executable as a supervised child without native linkage.
- [x] 3.2 Execute real collection/package workflows including media, identity/review preservation and more than 10 MiB of collection data; retain a repeatable artifact reopened by native Anki.
- [x] 3.3 Update current documentation, validate OpenSpec and repository-required checks, and record verified handoff evidence distinguishing native and WASM results. Git and delivery receipts own publication status.

## 4. Configured command selection

- [x] 4.1 Add the validated backend setting and consolidate standalone command sources behind one configured public target; preserve private offline fixtures and separate adapter packages.
- [x] 4.2 Verify default/native and explicit WASM public targets, reject invalid values, inspect selected dependency graphs, and run existing transport/offline checks with current documentation and required delivery gates.

Before implementation, acceptance must detect accidental native linkage in the
WASM graph, missing embedded guest input, an executable name change, invalid
backend values, a changed host platform, lost offline-service errors, and a
WASM executable mistakenly used for the native interoperability fixture.

## 5. Guest packaging and review

- [x] 5.1 Move production guest embedding into the public WASM adapter with `Run(arguments)` and retain caller-supplied guest execution in an internal host; verify actual worker, ABI and offline E2E cases.
- [x] 5.2 Review the complete PR with review-agent, inspect each actionable finding and concrete simplification, implement justified in-scope changes, and recheck affected behavior.
- [x] 5.3 Validate both configured public builds, standalone execution, dependency boundaries, OpenSpec, semantic lint and repository quality before updating the PR.

Preserve coverage for missing guest packaging, premature HTTP/2 disconnects,
memory/protocol failures, traps, cancellation, optional compiler cache errors,
and accidental production-guest dependencies in the executed CLI fixture worker.

## 6. Review corrections

- [x] 6.1 Add E2E coverage before the fix for non-UTC collection timing and Chrono timezone behavior at winter/summer, ambiguous, nonexistent, half-hour and skipped-date boundaries; retain a verification artifact.
- [x] 6.2 Bridge UTC and local-time offset lookups from the guest to the existing Go host and patch pinned Chrono's WASI fallback, preserving native code and dependency versions; verify the E2E cases and native/WASM workflows.
- [x] 6.3 Remove redundant guest copying, the single-consumer command composition package and forwarding wrappers; re-review the corrections and complete required delivery gates.

Failure scenarios include overwriting a non-UTC collection offset with UTC,
using a startup-only offset for historical dates, reversing ambiguous offsets,
accepting a nonexistent local time, losing a half-hour transition and failing a
24-hour skipped date. Cover Chrono itself so backup local-time consumers share
the same corrected implementation as scheduling and timestamp formatting.

## 7. Shared Rust CLI and prerequisite extraction

- [x] 7.1 Extend CLI E2E cases before implementation to check shared help and Rust-only argument validation, including duplicate options, without initializing the Anki guest; retain verification output.
- [x] 7.2 Invoke the shared Rust CLI in a small WASI command and pass only its validated transport selection to the Go host; remove the Go Cobra server tree and verify worker/application CLI and transports.
- [x] 7.3 Publish the setup-only prerequisite against master, preserve the WASM implementation in #181, and record that overlap removal follows prerequisite merge.

## 8. Built-in WASI platform selection

- [x] 8.1 Remove wasi_module and its package, use rules_rust platform attributes and direct guest outputs, and verify default/optimized guest builds, host-native binaries and existing E2E artifacts.
- [x] 8.2 Document normal compilation-mode behavior and configure required quality/lint validation for both master-based PRs.

## 9. Merged setup and dependency extraction

- [x] 9.1 Rebase onto merged #182, remove setup-only overlap and review the complete remaining diff with review-agent; no findings.
- [x] 9.2 Publish the guest dependency patches, annotations and generated lock metadata in master-based prerequisite #183; verify its native/WASI consumers independently and keep the server in #181.
- [x] 9.3 Remove dependency overlap from #181 after the prerequisite merges; rebase and validate the resulting candidate.

## 10. Full Rust CLI execution

- [x] 10.1 Record command/transport/lifecycle failure cases and add a direct WASI CLI E2E before implementation; verify it fails against the backend-only guest and retains CLI output after the fix.
- [x] 10.2 Execute the shared Rust CLI and gRPC server as one WASI command using a single-thread executor and host byte streams; verify stdio/TCP/Unix behavior, reconnect state, native errors and native compatibility.
- [x] 10.3 Remove the allocation/dispatch ABI, separate CLI guest and Go gRPC implementation; verify actual command cancellation/traps, help without backend filesystem/cache requirements and large native-readable collection artifacts.
- [x] 10.4 Review the completed diff for correctness and simplification with review-agent, apply justified changes, then validate and publish the exact candidate through required gates.

Acceptance must detect a guest that only exports backend functions, a separate
Go command parser or gRPC dispatcher, HTTP/2 deadlock under fragmented/large
I/O, connection changes recreating backend state, unbounded stream queues,
idle busy polling, leaked host sockets/goroutines, lost native trailers,
ignored command cancellation, automatic replay after a trap, and help depending
on Anki initialization or a writable cache. Network I/O remains a host
capability; upstream network services remain explicitly unsupported.

The review added a genuine process-I/O gap: malformed HTTP/2 input must make the
installed worker exit even when its peer retains the stdin writer. Record the
pre-fix timeout before making inherited protocol pipes interruptible.

## Trunk engine composition integration

- [x] Preserve trunk client factories, native in-process clients and shared Rust dispatch while integrating the full WASM command.
- [x] Verify all four cgo/WASM flag combinations, qualified worker schemes, external-only failures, pure-Go linkage and native reopening of the large collection artifact.
- [x] Review the integrated candidate and configure required quality/lint gates for guarded publication.

Failure scenarios recorded before integration: unavailable engine schemes must fail clearly; a WASM-only executable must not link native code; enabling both engines must retain native in-process clients while its default worker uses WASM; disabling both must still support external clients and reject local hosting; explicit native interoperability fixtures must remain native independently of public flags. Client API renames and dispatch extraction must preserve collection lifetime, error trailers and process shutdown.

## Runtime simplification

- [x] Poll the root WASI server future directly and combine read data/error handoff without changing readiness or bounded stream ownership.
- [x] Verify existing command, transport, cancellation, shutdown and offline collection E2E checks, retaining their native-readable artifacts.

Failure scenarios: tasks spawned while polling must run before waiting; host
notifications must not be lost; a final read containing both bytes and EOF
must deliver its bytes before EOF; cancellation and malformed protocol input
must still release blocked pumps.

## Package ownership cleanup

- [x] Move the WASM adapter to `pkg/ankigrpc/wasm` and shared Rust dispatch plus command-ID generation to `internal/backend`; update consumers, visibility, Clippy and current documentation.
- [x] Verify native/WASM builds, engine dependency isolation, transport and collection workflows, and source review.

Failure scenarios: moved packages must not lose guest embedding, generated
command IDs, direct native calls or Clippy coverage; WASM-only builds must
remain free of native linkage; both engines must retain their existing behavior.

## WASM ownership and readiness simplification

- [x] Make the command handle table the single resource-close owner, acknowledge one outstanding write through a completion channel, and unify guest wakes with host readiness notifications.
- [x] Review ownership and wake ordering; verify large-message throughput, flush/error propagation, command cancellation, stdio/TCP/Unix behavior and native-readable collection workflows.

Failure scenarios before implementation: a spawned or self-woken task must
not sleep waiting for unrelated I/O; flush must wait for actual write completion;
a failed write must remain failed; closing a command must interrupt pending
reads/writes and join pumps; dropped streams and command cleanup must not close
the same owned resource twice. Preserve bounded queues and idle blocking.

Decision: proceed with one close owner and one notification path. The guest
command already serializes handle operations. A single outstanding write removes
shared mutex/counter state at the cost of a smaller write pipeline; compare the
existing large-message E2E and collection round trip before retaining it. Keeping
the current implementation preserves pipeline depth but keeps redundant state;
a new runtime or executor dependency would expand scope without being necessary.

## Runtime threading documentation

- [x] Clarify the native Go embedding purpose and current guest threading limitation against the selected Rust target and pinned wazero source; scope the design rationale to this integration.
- [x] Add the sourced comparison of major runtimes usable from Go, including dated star counts, guest threading capabilities, Go-binding caveats and cgo requirements.

## In-process WASM client

- [x] Add production-guest E2E cases before implementation for lifecycle, instance isolation, large messages, native errors, offline media/backups and application scheme selection.
- [x] Add the managed in-memory gRPC client and explicit WASM scheme while preserving native direct clients and the default worker.
- [x] Verify collection workflows, both engines, pure-Go linkage and final cleanup; update documentation and the strategy diagram.
- [x] Configure required quality/lint checks and guarded publication; Git and delivery receipts own the final validation and publication outcomes.

Failure scenarios: cancellation during compilation/readiness must not leak an
instance; cancellation after readiness must not stop it; Stop before/during
startup and repeated Stop must be terminal and safe; retained connections must
reject calls after shutdown; independent clients must not share backend state;
large responses, native trailers and offline errors must survive in-memory
transport; a failed guest must not reconnect/replay; compile/cache failure must
be observable; defaults, native selection and unavailable-engine errors must
remain explicit. Verify actual collection artifacts, not only readiness.

## Shared gRPC transport

- [x] Let the existing gRPC client accept a caller-supplied connection dialer and retain terminal close errors.
- [x] Keep only guest execution and stream ownership in the WASM client; reuse shared readiness, connection borrowing and native error decoding.
- [x] Configure production-guest lifecycle/large-message/offline and native, WASM and combined CLI workflows as exact-candidate delivery gates, and complete the PR source review.

Failure scenarios before refactoring: a one-use memory stream must never be
redialed or replayed; startup cancellation must still bound waiting; Stop must
reject retained clients, cancel the guest and expose final cleanup; a guest that
fails before readiness must retain its actual error; repeated Stop must preserve
close failures; TCP and Unix reconnect behavior must remain intact. Reuse the
existing real-guest and network E2E cases instead of adding helper-mirroring tests.

## CLI exit status compatibility

- [x] Preserve the Rust CLI's nonzero exit status in both the production and fixture WASM command hosts, and compare exact statuses in the installed CLI E2E.

Failure scenario before correction: invalid native CLI arguments exit with code 2,
but the Go WASM entry point replaces that with code 1. Shell callers must receive
the guest status; host setup and execution failures without a guest status still
exit with code 1. The existing real CLI parity E2E owns this behavior.

## Subprocess option ownership

- [x] Expose the existing subprocess parser's options for executable composition without duplicating URL validation or changing standalone argument vectors.
- [x] Verify standalone and same-executable application workers through the focused installed CLI E2E and configure existing client and collection workflows as guarded delivery gates.

Failure scenarios before implementation: extracting parsed options must preserve
empty arguments, executable path encodings, shutdown timeouts and existing
validation failures. Application engine selection must not inject host flags
into an independent standalone Rust executable.
