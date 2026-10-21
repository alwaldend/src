# Validation observations

## WASM ownership and readiness

The command handle table now owns resource closure. Stream pumps report results
through bounded channels, with one outstanding write and command-owned progress;
the writer mutex/counter and per-resource close guards are removed. Guest tasks
and I/O share the retained host notification, replacing the guest atomic wake
flag and separate scheduling checks. Source review found no issues in this delta.
Five large-protobuf exchanges took 1.23 seconds in total versus a 1.30-second
baseline, keeping the measured runs within the baseline range.
All 12 optimized E2E suites and all four public engine configuration checks
passed, including native reopening of the large WASM package. Evidence and
delivery receipts live under `out/anki-wasm/ownership-*`.

## Package ownership cleanup

Optional engine adapters now live under `pkg/ankigrpc/native` and
`pkg/ankigrpc/wasm`. Shared Rust dispatch and command-ID generation belong to
`internal/backend`; their source bytes and crate identity are preserved.
Consumers, visibility grants, Clippy coverage and the rendered package diagram
follow the new owners. All 12 optimized E2E suites and all four configured
engine graph/linkage/RPC checks passed. The large WASM-produced package reopened
through native Anki. Source review found no issues in this delta.
Evidence and delivery receipts live under `out/anki-wasm/layout-*`.
Earlier observations retain the target names used when they were recorded.

## Composition and runtime simplification

Go build constraints now own optional source selection; explicit variant
settings replace nested Bazel conditions. The WASI event loop polls its root
server future directly, and read pumps hand off bytes with any terminal error.
The review found no actionable issues in this delta. All 12 optimized E2E
suites passed, including malformed-input shutdown, reconnection and native
reopening of the large collection package. All four configured engine graphs,
executable RPCs and pure-Go linkage checks passed. Retained evidence lives in
`out/anki-wasm/simplify-*`; delivery receipts own exact publication status.

## Trunk engine composition integration

The user approved retaining trunk's optional-engine composition and scoped
client factories. Independent cgo/WASM booleans now preserve native-only,
WASM-only, combined and external-only application builds. The shared Rust
server delegates to trunk's backend dispatch library; the WASM guest still
executes the full Rust CLI without a Go gRPC server or Anki dispatch ABI.

Public-target builds, selected dependency graphs and executable RPC checks
passed for all four flag combinations. Pure-Go configurations have no native
ELF interpreter; combined builds retain native in-process calls while the
worker returns the explicit WASM offline error. Corrected native/WASM/combined
CLI E2E suites pass. The large WASM collection workflow passed native reopening
with note identity, scheduling, review history and media intact. Its retained
`large-field-roundtrip.colpkg` contains the modern collection and media members.
The updated dependency diagram rendered and was visually inspected.

Retained failures: `ANKI_ENGINE_WASM_STUB_FILENAME` identified Go's architecture
suffix filtering of the disabled-engine stub, corrected by renaming it.
`ANKI_ENGINE_CGO_FIXTURE_SCOPE` identified an unconditional in-process CLI case
in the WASM-only suite; run that case only when cgo is compiled and continue
checking rejection without cgo. `ANKI_BAZEL_DANGLING_MANIFEST` identified one
shared repository-cache symlink into another output base. A force refetch of
that pinned crate with repo-contents caching disabled restored the task's
managed output; dependency pins and host configuration stayed unchanged.

Review: no remaining findings. Exact-candidate build, quality, lint, Clippy,
optimized E2E and publication status remain owned by delivery evidence under
`out/anki-wasm/`, rather than this narrative. The ergonomics review retained
bounded failure evidence and obtains configured binary paths from cquery:
Go configuration transitions can place them under suffixed output directories.

## Earlier validation observations

Recorded on 2026-10-19 in the task worktree originating at `03a151e54`.
This record owns validation observations and limits; Git and delivery receipts
own publication status.

- Native `ankiclient_test` and application `cli_workflow_test` passed after the
  cgo package move. The initial invocation also contained a misspelled transport
  target; `transport_test` passed in the corrected separate invocation.
- `//projects/anki_grpc/pkg/ankiserverwasm:backend_test` passed with an actual
  Rust WASI module: byte ownership, native error bytes, malformed modules,
  invalid ranges/statuses, traps, and CPU-bound cancellation.
- `//projects/anki_grpc/test/wasm:worker_test` passed with `pure = "on"`: a copied
  standalone executable hosted the Rust guest over child-process stdio and
  TCP/Unix gRPC. Native error trailers and graceful reaping were preserved.
  The test produced `wasm-worker.json`; the ABI test produced `guest-response.pb`.
- The original LLVM header probe failed because no WASI `stdio.h` was available.
  The user approved the WASI SDK 34.0 sysroot after the concrete probe.

## Approved offline guest

The user approved offline-first WASM on 2026-10-19. The production guest and
application targets now build and execute with cgo disabled. Descriptor-derived
handlers return native error details explicitly identifying unsupported network
sync, media sync, remote downloads and network services. Anki's HTTP sync-server
constructor/start entry points also reject startup in WASI; the Go host's gRPC
stdio/TCP/Unix transports remain available. Native cgo behavior remains available.

The guest uses the worker's filesystem authority and absolute-path protocol.
Callers must keep exclusive offline access: the SQLite WASI VFS cannot coordinate
locks with other processes. Backups and compression are synchronous. Guest and host builds
now use Bazel's selected compilation mode; optimized builds use `-c opt`. Existing Anki and FSRS pins receive reproducible WASI-specific patches.

The real guest passed collection open/close, local file-URL media import,
backup/await completion, and explicit network rejection. The HTTP sync-server
probe passed without filesystem or network capabilities. Both native and WASM
CLI workflow suites passed after the source port. The full 39-case application suite passed with an isolated cache, including
native reopens of WASM artifacts and the 513 MiB media workflow.

The first uncached 39-case application suite exceeded its 900-second deadline
while repeatedly starting fresh workers; the four targeted suites passed.
`ANKI_WASM_REPEATED_COMPILATION` records that measured startup cost. An opt-in,
caller-owned `ANKI_WASM_CACHE_DIR` now uses wazero's existing compiled-code cache;
default execution still requires no cache files. Worker E2E covers two separate
processes using the same cache and explicit rejection of an invalid cache path.
Application suites use isolated Bazel temporary cache directories; the installed
executable scenario disables the cache to verify standalone operation.

## Resolved port failures

- `ANKI_WASI_TOKIO_FEATURES`: Native workspace feature unification enabled Tokio
  features guarded on WASM. Its documented opt-in is scoped to WASI compilation.
- `ANKI_WASI_NATIVE_NETWORKING`: Cargo target conditions exclude networking
  dependencies; shared local storage/sync metadata stays available.
- `ANKI_WASI_TEMP_DIRECTORY`: Rust's WASI `std::env::temp_dir()` panics. Backend
  initialization supplies tempfile's explicit directory override from host
  `TMPDIR`; Anki IO's error context reads the same directory without the panic.
- `ANKI_WASI_THREADS`: Backups complete synchronously, compression never selects
  threaded mode, FSRS progress does not spawn a guest thread, and FSRS split
  evaluation runs synchronously before reading its result channel.
- `ANKI_WASM_INCOMPLETE_HANDSHAKE`: The Go stdio host observes gRPC transport
  admission and rejects EOF before HTTP/2 setup, matching native CLI behavior.

Raw build/test logs remain under ignored `out/anki-wasm/`.

## SQLite toolchain acceptance cases

Before implementing the SQLite probe, its E2E cases require: a WASI guest
using the existing pinned rusqlite dependency creates and commits an on-disk
database larger than 10 MiB with Anki's exclusive-lock/WAL configuration;
native SQLite reopens it and verifies the exact payload; a second WASI
instance reopens it; a rolled-back transaction leaves the payload unchanged;
and a missing parent directory fails instead of claiming a successful write.
The test retains the database and verification output as repeatable artifacts.

The SQLite E2E passed with the pinned rusqlite/SQLite dependency and the
approved sysroot. A committed 11 MiB payload survived rollback and was reopened
by native SQLite and a fresh WASI instance. A missing parent directory failed.
The retained artifacts are `wasi-sqlite.anki2` and `sqlite-verification.txt`.

The application package diagram's maintained SVG and WebP targets rendered
successfully; the WebP was visually inspected for readable labels and edges.
Semantic lint passed for the native adapter, application executable, WASI
host, worker fixture, and SQLite E2E targets before final preparation.

## Session ergonomics

- `ANKI_WASM_TARGET_SELECTION` (live): one initial native test invocation used
  a nonexistent default transport target. The corrected explicit target passed;
  continuation should use the labels above rather than infer defaults.
- `ANKI_WASM_EVIDENCE_VOLUME` (live): two dependency/README inspections exceeded
  useful output size. Lock inspection should report counts and changed field
  names; the final crate repository set, URLs and checksums are unchanged.
- `ANKI_WASM_PREVIEW_AUTOMATION` (live): collaborative preview navigation worked,
  but snapshot/evaluate failed on both SVG and HTML. The owning hermetic WebP
  render provided a repeatable visual artifact without changing browser setup.
- Resume from these task records and the ignored build logs. Runtime/SQLite
  fixtures and broad offline acceptance passed; delivery receipts track the final candidate.

## Offline acceptance evidence

- `offline-quality-preflight.log`: all 26 repository quality tests passed.
- The WASM application dependency query returned no cgo server or Rust static
  library targets. The package diagram was rendered and visually inspected.

- `offline-app-e2e-cached.log`: five Bazel test targets passed in 190.6 seconds;
  the application suite contains 39 top-level test cases. The uncached suite's
  900-second timeout is retained separately, without relabeling it as a pass.
- Native CLI workflow and WASM HTTP sync-server probe passed in
  `offline-app-e2e-03.log`; that aggregate command failed only at the uncached
  broad suite timeout.
- `offline-lint.log`: semantic lint built 74 affected targets successfully;
  the original guest transition built WASI modules and native consumers in the
  same graph. Built-in Rust platform attributes now own that selection.
- `large-field-roundtrip.colpkg` contains modern `collection.anki21b` storage and
  media. Its 17 MiB note field, card IDs/scheduling, review logs and media were
  checked through native Anki after WASM export.
- `large-api-roundtrip.colpkg` retains the separately verified 513 MiB media file.
  Artifact checksums are observation-specific; the test targets reproduce the
  same semantic checks on fresh API-generated IDs and timestamps.
- No Cargo crate repositories or versions were added or removed. Module metadata
  changes are generated build content and the maintained FSRS patch declaration.

Final candidate OIDs, required checks and publication state belong to the ignored
`out/anki-wasm/` delivery receipts and Git; they are not duplicated in this record.

## Configured command acceptance

- One command directory and public target per project select the implementation
  with `//projects/anki_grpc:backend=cgo|wasm`; cgo remains the default.
- `config-native-build.log` and `config-wasm-build.log` record successful public
  builds. Adapter queries select only cgo by default and only WASM when requested;
  the WASM dependency graph contains no cgo adapter or Rust static library.
- Both WASM executables retain their original names and are ELF64 x86-64 with
  no dynamic-library dependencies. Only the embedded guest targets WASI.
- `config-invalid-setting.log` records rejection of an invalid backend value.
- `config-e2e.log`: all seven targeted native/WASM E2E suites passed in
  202.8 seconds, including large-data export, native reopen, media, CLI,
  transport, and offline capability errors.
- The updated package diagram was rendered and visually inspected. Test fixtures
  explicitly select private variant targets so native interoperability checks
  cannot accidentally reopen with a WASM worker when public configuration changes.
- `ANKI_WASM_EVIDENCE_VOLUME` recurred with an unnecessary full worktree inventory
  and broad README reads. Current-worktree Git metadata and bounded README sections
  are sufficient for resuming this task; no shared policy change is needed.

## Adapter packaging and review acceptance

- Production guest generation and embedding belong to `pkg/ankiserverwasm`,
  which exposes `Run(arguments)`. The reusable byte host and ABI tests moved to
  `internal/wasmhost`; fixture workers have no production-guest dependency.
- The requested review-agent review inspected the complete merge diff against
  `03a151e5497e256235bfae79478e496dbb33b14c`. It found one actionable P2 defect,
  `ANKI_WASM_LOCAL_TIME`: Chrono's WASI UTC fallback let collection timing replace
  New York's `localOffset` of 240 minutes west with zero. The pre-fix failures are
  preserved in `timezone-before.log` and `timezone-anki-before.log`.
- The WASI-only Chrono patch uses timestamp-specific Go host timezone callbacks.
  Rust E2E fixtures cover seasonal offsets, gaps/overlaps, Lord Howe's half-hour
  transitions and Apia's skipped date. The real Anki timing call now preserves
  the non-UTC collection offset; both tests retain verification artifacts.
- The follow-up review reported no findings in the timezone correction and the
  simplifications. Reports and source snapshot digest are under ignored
  `out/anki-wasm/review-agent-*.md` and `embed-review-source.sha256`.
- Simplifications remove redundant guest copying, fold the single-consumer
  command composition package into its command, and bind adapter APIs directly.
- `embed-reviewed-e2e.log`: all 11 native/WASM test targets passed in 220.4 seconds,
  including CLI, transport, guest/SQLite behavior, native reopen of large collection
  data and 513 MiB media round trips.
- Generated dependency metadata changes only Chrono's patch declaration and
  extension projections; no dependency version or repository was added/removed.

- Final preflight: all 26 repository quality tests, semantic lint, all three
  OpenSpec validations and both configured public builds passed. The updated
  package diagram was rendered and visually inspected. Dependency queries found
  neither a production guest in the ABI fixture nor native linkage in the WASM
  application. The reviewed source snapshot remained unchanged after formatting.

## Shared Rust CLI and setup prerequisite

The setup prerequisite [#182](https://github.com/alwaldend/src/pull/182) is merged.
Its native adapter composition, approved build inputs and SQLite fixture are
part of master. The rebased #181 contains the WASM-specific server, guest port,
runtime and application integration. The merged setup change records remain
unchanged; guest dependency adaptations have their own shared-build record.

The CLI parity E2E failed against the previous Go parser's different help before
implementation. The updated installed worker matches native Rust help and
validation, including duplicate options and unavailable TMPDIR/cache, and retains
`rust-cli-parity.txt`. Worker transports, backend ABI and both application CLI
suites pass. Review-agent found no issues in the CLI reuse delta. The Rust CLI
WASI command shares parsing, help and diagnostics; its checked callback only
copies transport selection, then Go opens the listener after guest completion.

The full WASM collection workflow and native Clippy/generated CLI reference also passed after Rust CLI reuse. The CLI guest is 93,991 bytes and depends on the std-only server CLI rather than the Anki backend. Remote code review of setup candidate32958d1 completed with no findings.

## Built-in platform selection and rebase

Rust guest rules select `@rules_rust//rust/platform:wasip1` directly and Go embeds
their actual outputs. The ABI fixture lives in its worker package; no generic
transition/copy wrapper remains. The approved stock WASI SDK sysroot supplies
SQLite libc headers/libraries through the pinned LLVM toolchain.

Default/debug backend, worker CLI/transports, SQLite, offline-service errors,
HTTP rejection and application CLI checks passed after the wrapper removal.
The default/debug full collection workflow exceeded its existing 15-minute
deadline (`remove-wrapper-debug-e2e.log`); it is a failed observation.
Normal `-c opt` selects optimized builds for both host and guest. Optimized
full-workflow and final candidate validation are tracked by the delivery logs.

`ANKI_WASM_EVIDENCE_VOLUME` recurred while resolving generated lockfile hunks:
a small line range contained a very large encoded metadata value. Parse JSON
and report changed field names/counts instead of displaying lockfile hunks.
The regenerated lock matches the preserved pre-rebase guest metadata; no
repository identity or dependency version changed during the rebase.

The full rebased review against merged #182 reported no findings. The user then
requested a separate dependency PR. Its scope is the Anki/FSRS/Chrono patches,
Tokio WASI annotation and generated module metadata. Server/runtime code and
new consumer visibility remain here. Both PRs targeted master. [#183](https://github.com/alwaldend/src/pull/183)
is now merged; the dependency files and their change record belong to master.

Two local orchestration errors were corrected before validation: a missing
`TMPDIR` caused the first optimized command to fail before tests ran
(`ANKI_WASM_TEMP_DIR`), and delivery inspection was attempted before its new
worktree launcher had been generated (`ANKI_DELIVERY_ENTRYPOINT_READY`). Create
and verify task scratch first, then await launcher generation before using it.
The failed command logs remain separate from subsequent successful evidence.

`rebased-opt-e2e.log` records all seven optimized E2E targets passing after the
rebase and direct-platform change. The full collection workflow completed in
about 268 seconds, including native reopen and large field/media round trips.
The independent dependency branch also passed a direct WASI Anki library build,
31 native/SQLite/quality targets and semantic lint. Both review-agent reviews
reported no findings; exact final candidate gates belong to delivery receipts.

## Merged dependency prerequisite

The dependency merge contains byte-identical copies of all ten overlapping
patch, annotation, lock and change-record files. Rebase onto merged #183 removes
those files from the PR diff; only new consumer visibility remains in the Anki
package. Server/runtime and application integration stay in #181. Rebase refs,
review results and exact-candidate validation belong to the ignored
`after-deps-merge-*` delivery records.

The local rebase completed without source conflicts. All seven optimized E2E
targets passed with unchanged runtime inputs. The tree differs from the prior
published candidate only in this change's current-state documentation.

`ANKI_DELIVERY_FROZEN_SCOPE_SHRINK` and `ANKI_DELIVERY_REBASE_INSTALL` (live):
three preparation attempts stopped at scope/index bookkeeping gates. The exact
clean-worktree, index-tree and submodule checks passed before a local Git rebase.
No source or progress was dropped; exact fetched OIDs and failure evidence are
under `out/anki-wasm/after-deps-merge-*`. Final validation and lease publication
remain owned by the delivery adapter. No tool source or safety check changed.

## Complete Rust CLI correction

`full-cli-before.log` records the direct WASI help check failing against the
backend-only guest: it produced no Rust CLI output. The full command executes
shared Rust parsing, gRPC and dispatch; Go now supplies platform streams.
The first host/worker E2E passed all command, transport, reconnect, timezone,
cancellation and trap cases. The first broad run passed ten of eleven targets.

`ANKI_WASI_NESTED_EXECUTOR` (live, corrected): local file URL media import panicked when
Anki's futures::block_on ran under LocalPool's executor-entry guard. Replace
the outer pool with direct readiness-driven task polling; retain the failing
`full-cli-e2e.log` separately from correction and final candidate evidence.
The broad collection round trip and native compatibility passed in that run;
local media acceptance is not complete until its correction passes.

`full-cli-executor-correction.log`: the real local file URL media import,
command cancellation/traps and worker CLI/network/timezone checks pass after
the nested-executor correction. Native behavior passed the broad run before
the WASI-only executor correction. The final large-data candidate run remains
a required acceptance check.

`full-cli-final-e2e.log` records all eleven optimized E2E targets passing after
both executor and inherited-pipe corrections. The run completed in about 204s.
It includes native lifecycle/transport/application compatibility, full-command
help/cancellation/traps, worker reconnects/timezones, SQLite interoperability,
offline Anki/media operations, HTTP rejection, application CLI workflows and
the large collection export reopened by native Anki. The collection artifact
retains its 17 MiB field, identities, scheduling, review history and media.
The follow-up review-agent source review reported no findings after the fixes.
Exact-candidate gates and publication state belong to the ignored
`full-cli-*` delivery receipts.

The first candidate's default-mode validation passed all 36 targets, including
repository quality and the native/WASM CLI and offline checks. Clippy then
rejected the repeated backend-initializer function-pointer signature as overly
complex. Name the initializer and queued-task types instead of suppressing the
lint; retain the failed check log separately from final validation.
`full-cli-clippy-type-check.log` records the corrected native Clippy passing.

`ANKI_WASI_STDIO_SHUTDOWN` (fixture-tested, corrected): review-agent identified
that inherited blocking pipes were not registered with Go's I/O poller. The
new real-process malformed-preface E2E retained its stdin writer and timed out
after ten seconds (`full-cli-stdio-before.log`). Unix protocol files now use
private nonblocking duplicates registered with Go before pumps start. Closing
them interrupts active reads/writes, allowing cleanup to join the pumps.
The corrected process exits with the Rust HTTP/2 protocol error in about 0.24s.
One correction run failed only because its diagnostic assertion used the wrong
case; `full-cli-review-fixes-verified.log` records all three host, installed-worker
and real offline Anki targets passing with the corrected assertion.

The simplification review also replaced FuturesUnordered and its stream-waker
registry with a task list using one shared readiness waker. The direct futures
dependency was removed from the shared server target. Bounded host pumps remain
necessary for asynchronous native I/O; the guest owns the complete CLI and
gRPC protocol. No new external dependency was added.

The closing ergonomics review reused this owner-local record, failure logs and
delivery receipts. `ANKI_WASM_EVIDENCE_VOLUME` also affects formatter progress:
one log line can contain substantial terminal output. Bound extracted output
by both line count and byte length. The wrong-case protocol assertion above
was an avoidable test correction; neither it nor the real pre-fix failures is
counted as passing acceptance evidence.

## In-process WASM acceptance

`ANKI_INPROCESS_WASM_UNREGISTERED`: the installed CLI E2E failed on the absent
`inprocess+wasm` scheme before implementation (`out/anki-wasm/inprocess-cli-before.log`).
The managed client now uses one production guest and an in-memory gRPC stream.
Five host, native, WASM and combined CLI suites passed in 238 seconds
(`inprocess-workflows.log`), covering local API selection, collection export,
state, apply, build, lifecycle and independent guest state. Generic in-process
selection prefers cgo and falls back to WASM in WASM-only builds.

The production-guest E2E retains a database with a 17,600,000-byte config value,
an offline backup package, media verification and non-UTC scheduling output.
The native direct client independently reopened that database and verified the
complete value (`inprocess-native-verification.json`). The first in-process
scheduling attempt exposed `ANKI_INPROCESS_TIMEZONE_FIXTURE`: changing `TZ`
after Go cached `time.Local` could not change the running host's timezone.
Set the test process timezone before startup; the corrected suite passed.

Both owner-local OpenSpec records validate strictly, and the rendered strategy
diagram was visually inspected. No external dependency or native transport was
added. Exact-candidate quality, lint, behavioral and documentation gates are in
`out/anki-wasm/inprocess-checks.json`; delivery receipts own their final results
and publication state. Bounded continuation evidence and observed discovery
friction are retained in `inprocess-ergonomics.json`.

## Shared transport and explicit worker selection

The WASM client now supplies one `net.Pipe` dialer to the existing gRPC client;
shared code owns readiness, connection lending and native error decoding.
Production-guest lifecycle and large-message checks pass, and the native client
independently reopened the retained database and verified its full config value.
The CLI parity E2E now compares exact exit statuses; both production commands
return 2 for invalid arguments.

Installed-executable checks verify native-only, WASM-only, combined and
external-only engine selection with an unusable WASM cache. Both generic
strategies prefer cgo and otherwise use WASM; qualified schemes select their
compiled engine. The first default-address check omitted its required collection
fixture; the corrected E2E passes in all three engine-enabled build variants.
The qualified strategy diagram was rendered and visually inspected.

`out/anki-wasm/shared-transport-checks.json` configures final quality, lint,
behavioral and documentation gates. Its delivery receipts own exact candidate
and publication results; `cgo-preference-selection.log`,
`shared-transport-native-verification.json` and `shared-transport-ergonomics.json`
retain bounded continuation evidence.
