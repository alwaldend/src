## Candidate and ownership

The published PR2 candidate was based on master
`98559dd11cbf65dc0658feb3bd80be0b468d6862`, following
the merge of PR195. It originally began at
`a2c655b20432b78e4feb75e16da1a95b73b3d30c` and extracts browser-only source from
`b3c926655884ab070d628314b907950707306cc9`.
Prior full-branch results do not validate this split candidate or the new test.
PR1 owns `internal/directclient`, `internal/timezone`, the shared backend ABI,
and both shared BUILD files, including browser visibility grants. PR2 includes
none of those files and requires no supplemental visibility patch.
PR2 preserves the original standalone WASI guest visibility grant to
`//projects/anki_grpc/pkg/ankiclient:consumers`, copied exactly from the immutable
source. This retains the existing build contract during the split.

The prerequisite [PR195](https://github.com/alwaldend/src/pull/195) merged at
2026-10-24 16:52:49 UTC. Remote metadata and fetched master were checked; master
contains published core `cd110f71f7a88187979a45037f94f3ac25cae0ff`'s direct-client
blob `00ec3bdb02f23120351833e8b73603cc8f16a69b`. The local rebase preserved the
browser source and both native and browser startup documentation without
conflicts. Core files are inherited from master, not duplicated in the PR2 diff.

## Corrected Emscripten provider composition

The adapter owner's incremental patch was applied without changing shared BUILD
files or this composition's documentation. Its source-ready-2 manifest SHA-256 is
`299c58f9a62984be535738d56c58866c65fac477df53eccbcbebfb1a6635ffba`.
Adapter blob `c5f135735a9d4f0ccb31a5c1d111988540616025` and provider fixture
blob `d561c335ceff620303e1f9c1310a2e3c72f102a9` match the owner exactly.
All four new groups failed on the old adapter, and all twelve provider groups
pass on the correction. Those red/green artifacts are retained separately from
the original composition below. No provider implementation was edited here.

The affected PR2 provider/Chromium batch passed both targets. Chromium executed
with the corrected adapter SHA-256
`88db5dc4fe456421d5bf7e3c7344038e7680659789998ba5eacc8c85c5212554`;
the provider test reused its matching cached result. The report verifies seven
workflow groups, two Rust instances sharing a 139264-byte collection, exact
Unicode note/binary/concurrent media persistence and zero remaining descriptors.
Generated MEMFS remains the exact 50202-byte artifact identified below. Existing
WASI host, Go helper, Rust and shared direct-client source blobs are unchanged.

Retained reports under `out/browser-pr-split/emscripten-composition/corrected-acceptance/`:

- Chromium `browser-storage.json`: SHA-256
  `e51ef9ac27b0e851382e4c91b58034caa8d4bebcdfb350fbb3d68e5800e9ca8a`.
- Provider `emscripten-filesystem.json`: SHA-256
  `a7cc38f7bd8a1475ea3ca581af2d6140d4482b7fa7e58c5c9c765c74e79a68b6`.

Fresh independent Astra ultra review of the full adapter correction, fixture,
changed documents and actual generated MEMFS/Go/WASI consumers returned no
findings. The parent accepted this read-only verdict for c5f135; no tests were
run by the reviewer. Review task suffix:
`browser-emscripten-adapter-corrections-review-20261024-1`.
Corrected composed SDK tree `a751c95acafad540ab77129ff5259c72dff677bb`
passed forced Chromium execution: 19 workflow groups, 20 workers, Go/Rust WASI
collection operations, shared MEMFS, IndexedDB restoration and native artifact
comparisons. Ten backend/provider runtime source blobs were checked against this
composition and match. SDK code was not imported into PR2. The browser report
SHA-256 is `e13f48222532d7b3ab5dc6f0f0c4f84c5fd24f37aad5eb16522c5747c947efc5`;
the new SDK archive is `dea93048ce7b6822d598779109d0d209824f37b65fa8750995b517e3a1e14625`
and differs from the initial archive only in emscripten.js. Generated MEMFS is
unchanged and has no filesystem Wasm sidecar.

The six provider-independent SDK Node groups were explicitly reused, not rerun:
their executed archive members, harness, toolchain and Go sources are byte-identical.
The integration owner's immutable handoff has SHA-256
`59f0647c4ff238a157b89d8830f6d9a52439203d0b0a175d684b995444fb91ec`.
Its report and reuse explanation are retained under
`out/browser-pr-split/emscripten-composition/sdk-corrected/`, with the source bridge
in `sdk-corrected-bridge.json`. This remains composed behavior evidence, separate
from PR2's independent Chromium run and exact-candidate gates. No site stress,
browser-process restart or Rust Emscripten-target support is claimed.

The parent accepted the review and corrected behavior evidence and released the
publication hold. Root quality, affected lint and consumer checks remain mandatory
for the exact prepared candidate; the local delivery receipt and handoff own
their results and the final Git/publication state.

## Historical initial Emscripten composition at adapter 417aec8

Build candidate `15927fded68721a4c6b3110d9d9a64f5f52b8468` was composed from
its 18-path delta against b3, with only additive shared BUILD/MODULE hunks and
reconciled current planning. The adapter's initial 10-path manifest has SHA-256
`56f2d1b010d41551acf4b1c0a9a9badef635bbb32c27217d1e147e71bc618392`.
Its exact implementation blob is `417aec826d72e599b2b4faea6e660d4b9a41fba3`.
No inherited SDK, shared-core or stale Go filesystem sources were copied.

Before importing the provider, the Chromium fixture was wired to create its
own filesystem and `/browser-backend` root. The runner retains caller injection
and serves the adapter/generated module explicitly. Initial analysis found that
rootpath arguments require direct data labels; those labels were corrected,
with the failed analysis retained. It was not a runtime behavior failure.

Independent PR2 `browser_storage_test` then passed on
HeadlessChrome/150.0.7871.24 with the initial adapter. Seven reported cases cover
invalid-path recovery, real collection/note/media writes, concurrent media,
backend replacement and timezone callbacks. Both Rust instances use the same
provider; the collection is 139264 bytes and no provider descriptors remain.
The JSON SHA-256 is
`2e94888c7e035ee8a84b8c337f182ab0e5932307543b8b55b647ea16fdac7760`.
Generated MEMFS is exactly 50202 bytes, SHA-256
`97a0bd83047f70c80ed30ca04d61206e6fcff9c8825201131a22464e0239c41d`.
The raw-artifact and initial provider E2Es also passed. Reports and original
inputs remain under `out/browser-pr-split/emscripten-composition/`.

Initial composed SDK tree `f3bfa93c699cb1d38ade5b28f0980c7b5062d8e2` passed
both SDK tests: 19 Chromium workflow groups with native collection verification
and six separate Node groups. Its exact report, source manifest and archive
hashes belong to the integration owner's `emscripten-browser-validation`
artifacts. This is composed behavior evidence, not independent PR2 delivery.
The same initial provider has four accepted review findings; the adapter owner
is correcting them. These initial passing reports do not validate its future
correction. The final source bridge and affected rerun remain pending.

Rust remains WASI; the resolved architecture uses standalone Emscripten MEMFS
for shared storage. SDK/site ownership and application snapshots remain separate.
Fresh provider review and exact-candidate delivery gates must finish before
publication; the parent holds publication during this coordinated work.

## Host close-error correction at bf220235

Fresh independent review of the preceding local correction
`b2c11f772ed1169b49370778243426dcf836ff23` found that aggregating a retained
sync error with a recognized close error made both terminal. Four failure-first
cases now exercise recognized and unknown close errors, each before and after
actual provider release. The old corrected host produced two retained
`Sync and descriptor close failed` failures; both unknown-error controls passed.

The narrow correction returns the provider close error in preference to the
retained sync error. Recognized errors remain recoverable; unknown exceptions
remain terminal. All four cases assert one close attempt, an invalidated guest
descriptor and no shutdown retry. For an injected failure before release, only
the fixture reclaims the deliberately open provider fd; zero remaining fixture
descriptors does not claim that the host can resolve an uncertain close failure.

Final wasi.js blob: `82f8220661346b2a77400e2484bbbd805867cc6f`.
Final filesystem fixture blob: `8e04837e3839c06d9a41c0ea5c06ae970a27da70`.
The failure-first patch is `out/browser-pr-split/close-tests-first.patch`.
Old-source JSON `close-red/test.outputs/filesystem.json` has SHA-256
`53dafc3cf1dd693b5539f986e56f2b3f7b48e7b6faa5256e456e775fd379207a`.
The final formatted fixture passes all 22 behavior cases; its JSON at
`close-green/filesystem_test/test.outputs/filesystem.json` has SHA-256
`9a93182ce46e95b06c331a0f49d0604359bb00988f1bf7d42c1187bcfca90731`.
Both node_storage_test and browser_backend_test also pass, the latter using
unchanged cached lifecycle inputs. All three reports and storage artifacts were
copied immediately. Affected semantic lint, configured formatting and strict
OpenSpec pass. The local handoff and preparation receipt identify the resulting
candidate and full patch; no runtime tests are represented as publication gates.

This changes only wasi.js, the existing behavior harness and related documents.
API, shared core, SDK, site and provider source are unchanged. Parent-coordinated
fresh review, provider-backed Chromium/SDK composition and exact-candidate
delivery gates remain pending. The selected Emscripten work and runtime
alternative investigation remain separately owned; no provider implementation
or dependency is included in this correction. No push or PR change was made.

## Earlier production-host I/O correction at b2c11f7

The parent accepted three source-confirmed P2 findings against
`a8e42622467d4d77efcf4516883d06f2b85a5979`: discarded partial I/O progress,
ignored final-symlink lookup flags and incomplete filesystem errno mapping.
The existing production-host harness was extended before production edits.
Its initial old-source run retained 17 failing assertions. After correcting a
fixture DataView that became detached when real Rust grew memory, the final
fixture was run against the same old host and retained the same 17 failures.
This fixture-only correction did not remove any assertion. Old host blob:
`e0f41cd40bad6feb4f794dbb30a1582447c6d726`.

The corrected host passes all 18 added behavior cases: partial reads and writes
across 64-KiB and iovec boundaries, exact resumed bytes and positions,
zero-progress errno, observable post-write sync errors, regular/dangling final
symlink protection and explicit follow, a native no-follow replacement race,
six recognized provider errnos and terminal unknown exceptions after progress.
Both complete host runs leave zero owned provider descriptors. The
node_storage_test also passes with collection/media persistence and concurrent
writes; browser_backend_test passes using unchanged cached lifecycle inputs.
Affected-package semantic lint, configured source formatting and strict OpenSpec
validation pass. No new provider, shared core, SDK or site source was added.

Corrected wasi.js blob: `8dee291c6f19140a20f871f848d0fa855026fee4`.
Final filesystem fixture blob: `592a57801ccb21a5ba3a47e25f878ed279c3ebf7`.
The matching old-source JSON has SHA-256
`927a8021c7c5369411b623141c85963a18d424af82f6eb3b73a4a1c06db94e8b`;
the green filesystem JSON has SHA-256
`c9c258cccab207de013bbcae7712e437522656450d1d116e9b4799ee7a73185d`.
`out/browser-pr-split/p2-source-ready.json` identifies all three Node reports,
the exact two-file correction patch, and source identities. Raw red/green
artifacts remain under `p2-red-final/` and `p2-green-final/`.

The errno cases inject symbolic provider errors over actual node:fs; they do
not induce host disk quotas or stale NFS handles. WASI numbers were checked
against the pinned SDK 34.0 errno header. Pinned Node/libuv may emit unknown
system-error strings for quota/stale failures; this correction does not parse
unknown strings or convert programming exceptions into ordinary errors.

The parent reports six SDK shared_filesystem_test groups passed against a8 in
composition `a0e3311e76c1a60601afd9c72011ebd213adabe7`. That is historical
composition evidence, not validation of this corrected source. The user has
selected Emscripten; its isolated build and adapter owners will supply the
provider and explicit Chromium fixture inputs. No provider source is imported
here yet. Fresh independent review, corrected SDK/browser composition and final
delivery gates remain pending parent coordination. Eventual delivery remains
authorized; this correction is prepared locally without a push or PR change.

## Historical caller-owned JavaScript filesystem preparation at a8e4262

The user authorized replacing the published Afero filesystem callback host at
`2151fa59b933aab01729c902906ae0454b420f6c`. NewClient now accepts only the
factory; the factory captures fs/root and receives a separate Go timezone host.
Rust calls the provider's synchronous JavaScript filesystem API directly.
`server.js`, the shared core and SDK source remain unchanged. The existing
per-backend command mutex and joined-factory startup behavior are retained.

The production-host filesystem failure cases were written before replacing
wasi.js. With the old source, the real reactor test failed with
`A Go Afero filesystem endpoint is required`; its JSON is retained under
`out/browser-pr-split/js-fs-red/`. The new host passed filesystem_test,
node_storage_test and browser_backend_test. The final production-host rerun
checks descriptor cleanup after initialization failure, close and a real Rust
trap; provider metadata, flags, positions, rights and memory bounds; and the
existing collection/media workflow over NewOsFs/BasePathFs. Two Rust lifetimes
share a 139264-byte collection, unchanged binary and concurrent media, and four
actual Rust timezone calls. No owned provider descriptors remain open. The
existing eleven command/lifecycle cases also pass, including delayed startup
cancellation and concurrent Stop/Done cleanup.

The pinned Go 1.26.5 wasm_exec.js and syscall/fs_js.go contracts were inspected.
The reusable :go helper installs borrowed fs/process/path globals after importing
wasm_exec.js and before Go.run, preserving provider receivers and runtime output.
Source identities: wasi.js `e0f41cd40bad6feb4f794dbb30a1582447c6d726`,
go.js `b6a42e9ea62e550cd0ca15fc26ac8fbc82af2253`,
and storage_js.go `27a2498f625f27a81276270963e23307a31d3f8c`.
`out/browser-pr-split/js-filesystem-validation.json` records all runtime/test
blobs and exact JSON hashes. The closed database and binary media are retained
with the Node storage report.

ZenFS was explicitly rejected. No substitute dependency or backing store was
added. Chromium remains required acceptance, and its runner now requires an
explicit caller-supplied provider module exporting fs/root. The default Bazel
target has no such module and is not claimed passing. Runnable browser
composition, composed SDK tests and fresh review remain pending with the parent.
Node success does not replace Chromium acceptance. Publication is held by parent coordination pending coherent browser composition,
SDK validation, fresh review and delivery gates; eventual delivery remains
authorized. PR196 remains at the previously published candidate. Its historical review and
delivery results below do not validate this replacement. A later master advance
touched only host_bot; the local preparation receipt records synchronization.

## Historical independent master-based verification of the Afero host

The rebased candidate `9c08789daf4c4c1865fef96f1cdba75956217159` passed full
configured formatting without changes and both browser E2Es independently.
Their artifacts were copied immediately to `out/browser-pr-split/independent-passing/`.
A fresh Astra ultra review then found a real concurrency gap in the sequential
storage coverage: Rust filesystem callbacks can schedule another Go command and
reenter the same reactor while its collection mutex is held.

Before changing production code, the existing Chromium fixture gained two
simultaneous `AddMediaFile` calls. The test failed with `JavaScript error:
unreachable`, retained in `concurrent-media-red/test.outputs/browser-storage.json`.
Stable defect: `browser-reactor-concurrent-entry`. The initial fixture compile
also exposed a configured loop-closure analyzer requirement; passing the index
explicitly into each goroutine fixed the fixture before the runtime experiment.
Its earlier copied output is labeled stale and is not failure evidence.

The browser backend now holds a per-instance command mutex across the complete
JavaScript call. Shared command cancellation and Stop/Done ownership remain
unchanged. Both E2Es then passed: all eleven Node workflow cases and seven
Chromium cases, including distinct concurrent-media bytes after both backend
lifetimes. The retained reports are under `concurrent-media-green/` with
SHA-256 values `3fb33147a599a655afc6e44a67b54274c75342af77f8e533e105f5d7e09ce0e7`
for the Node JSON and
`edc32b60785ef31d451b57f8e87f5e2947bf5a25563ba6a1ce83073f91df40f5` for Chromium.
These independent runs use merged master; the SDK result below remains historical
composition evidence and does not validate this correction.

Independent candidate `a4a2d01a2107bc96c72bfbb35c70c097fc9b8945` passed all
delivery checks at 2026-10-24 17:10:31 UTC: 28 tests including both E2Es and root
repository quality, semantic lint of all seven affected packages, native CLI
and standalone WASI guest builds, and native/shared Rust Clippy. The browser
reactor separately passed the Rust Clippy aspect. Strict OpenSpec validation
and configured formatting passed. The delivery test command reused matching
browser outputs from the independent uncached run above; 25 of its 28 tests
executed afresh. JSON contents and hashes were inspected, not just exit codes.

A fresh Astra ultra correction review returned no findings and resolved the
prior P1. It verified browser client blob
`26421bfcf26216be2e33d52fb5e258cb452f4e5b` and storage fixture blob
`3fea987f3919e739257cf31c35a61e3db3154aab`. Queued-command cancellation, Stop
and trap interactions were checked against implementation; the E2E does not
directly combine those cases with concurrent writes. No further source edit
was required. The existing startup fixture remains
`c0e56eb31f6473be082c5c7742b719a4d2cb2b71`.

`out/browser-pr-split/independent-validation.json` records source identities,
reports and validation receipt paths. `review-initial-result.json` and
`review-correction-result.json` retain both independent review rounds. Delivery
revalidates evidence amendments against the exact prepared head before publishing
to master; its receipts own the resulting publication identity.

[PR196](https://github.com/alwaldend/src/pull/196) opened against master on
2026-10-24 at initial published head
`86bff15e55a34bab64f910401603b380a013a832`, after a fresh passing delivery run.
It is registered in the child thread and the parent was notified. This
publication-record amendment preserves the reviewed runtime/test blobs; the
final receipt records its revalidated head and remote verification. No stack,
PR191 mutation, deployment, new dependency or merge was performed.

## Historical prerequisite checks on 2026-10-24

- OpenSpec strict validation passed for `add-browser-wasm-backend`.
- The configured formatter applied JavaScript, Go and BUILD formatting only
  within the new test sources. Rustfmt failed because `guest.rs` references
  PR1's absent `internal/backend/ffi.rs`; the full format gate remains blocked.
- Analysis-only selection of `browser_backend_test` and `browser_storage_test`
  failed on missing `internal/directclient`, `internal/timezone`, and
  `internal/backend:ffi.rs`. No tests ran in this standalone PR2 checkout.
- The guest-visibility follow-up restores only the exact original BUILD line
  and updates the scoped records. Known prerequisite failures were not rerun.
- No new dependency, deployment, PR191 mutation, push or PR creation occurred.

Logs and delivery preparation receipts are under ignored
`out/browser-pr-split/` in the dedicated child worktree.

## Historical composed behavior validation

Both `browser_backend_test` and `browser_storage_test` passed unchanged on
2026-10-24 at 15:28:53 UTC in worktree `t3-1b486911`, with core
`5e839d9611d293bbde51469234c7e1e539161cc7`. The exact tested composed tree is
`c1e893c248d5cb3af30dea51f80b4f3cfc08b5df`: original full source plus PR2's
five-file test projection and the core's corrected `directclient/client.go`.
The report identifies PR2 source `7253ef82b07dfd63551abaa03f60e9e09af9bafd`.

Blob comparisons verified all five test files against the tested tree and local
head `53c36e910104f50dc6ad0e09ca091f10af959a2a`. Across both PR2 heads and the
tested composition, the complete test package tree is
`0d67957a5b6d44ffe0c2c0e91c260b850210b3a8` and browser runtime tree is
`0d7b9ffc917941d9a8f0e77049725955f28c50d3`. The visibility/documentation amendment
therefore preserved the tested browser source through
`dea2027c67ca79ae39323973b20b97a6c3da3413`. Those patch-byte differences did not
represent source changes. The subsequent startup-contract change below alters
the delayed-start test, so this historical result does not validate that change.

`initial-passing/browser-storage.json` records HeadlessChrome/150.0.7871.24,
Go exit 0, two Rust instances, six completed checks, and no workflow error or
console messages. A Unicode note and binary media survive collection close,
backend shutdown and reopen in one Afero instance; Rust also observes Go-created
media. The final database is 139264 bytes. Winter/summer offsets are
-18000/-14400 seconds and four actual Rust timezone callbacks are recorded.
The Node test separately passed command, buffer, startup-cleanup and trap cases.
The composition report also records passing package formatting and configured
lint aspects; no source corrections were required.

Source evidence is `out/browser-pr-split/validation-report.txt` and
`initial-passing/browser-storage.json` in `t3-1b486911`. Inspected copies are
retained under this worktree's `out/browser-pr-split/composed-passing/`;
`composed-validation.json` records full paths, evidence SHA-256 values and the
five matching blob IDs. This establishes composed behavior only. Master-based
PR2 validation and publication still require PR195 to merge and all delivery
gates to pass against the resulting exact candidate.

## Startup completion follow-up

The user now requires Start to wait for factory completion, checking its context
before and after initialization. Cancellation during initialization returns only
after completion; a stuck factory delays startup failure and cleanup. Commands
retain their cancellation behavior. Concurrent Stop retains its own waiting
context and Done observes completion of active startup and cleanup.

The existing `main_js.go` startup case was adapted before any core production edit
in this worktree. It cancels Start while the JavaScript Promise is unresolved,
observes that Start stays pending during a bounded scheduling interval, rejects
connection access, and checks canceled concurrent Stop with Done still open.
After releasing the Promise, it requires Start's context cancellation error,
exactly one backend close, no published connection and completed Done. The
standalone storage test and non-startup behavior remain unchanged.

Phase1 paired old core client blob
`be927e53496b08f33cf0fc374913e23674774f36` with startup fixture blob
`4d2a484bb9ad1dbe1e7989ce74b1d5fd3ec7b4d7`. The run failed, but teardown raised
`Go program has already exited` before retaining the intended contract failure.
The raw log is `out/browser-pr-split/phase1/old-core-after-refetch.log` in
`t3-1b486911`, beginning at line 329. Stable defect:
`browser-startup-failure-teardown`.

The harness now defers factory release, joins a separate Start-completion
channel, then calls Stop and waits for Done under a fresh cleanup context before
returning. Cleanup errors join the original failure rather than replacing it;
the notification callback remains registered until cleanup finishes. The 50 ms
assertion and concurrent Stop path are unchanged. New fixture blob:
`c0e56eb31f6473be082c5c7742b719a4d2cb2b71`.

Only browser comments/documentation and the existing startup E2E change here.
The standalone storage test and other runtime behavior are unchanged. Targeted
formatting and OpenSpec validation are recorded locally; known
absent-prerequisite failures are not rerun.

## Composed startup validation

With the cleanup fixture and old core `5e839d9611d293bbde51469234c7e1e539161cc7`,
tree `35f2aa72c76561a4505b8eff470ff9f1bb760fee` failed with the intended
`startup returned before its canceled factory settled` assertion. The retained
JSON contains that original context-cancellation failure without a joined
cleanup error; the delayed backend closed exactly once. This confirms the
harness correction exposes the behavioral failure rather than losing it during
Go shutdown.

Core `cd110f71f7a88187979a45037f94f3ac25cae0ff`, client blob
`00ec3bdb02f23120351833e8b73603cc8f16a69b`, then passed both backend E2Es and
`//projects/anki_as_code/test/browser:browser_test` in composed tree
`de7883261b2bcf9478295217ff510b48eda0d107` at 2026-10-24 16:34:18 UTC.
The Node report records all eleven workflow cases, including canceled startup
waiting for the factory with concurrent Stop, zero delayed-backend commands and
exactly one close. Chromium records all six storage checks, Go exit 0, two Rust
instances sharing Afero, persisted note/media and actual timezone callbacks.
The SDK test passed fifteen workflow groups and native verification of the
browser-exported collection package: note identity, card scheduling, review
history and binary media were preserved, with the expected edited field.

The report uses browser runtime from PR2
`044e51b502e6f2faf02d9d3e1a0191c4a31be22a` plus the exact cleanup fixture blob
`c0e56eb31f6473be082c5c7742b719a4d2cb2b71`. PR2
`5a8548a7d89de65f8cf84e95830994c6a32693a6` matches the composed browser runtime
tree `162226f57796f413488265967144469337297bf9` and complete test package tree
`26a93ce540deb0b651f6e1c82a6cf5ea2be7b895`; all fifteen file blobs were compared.
This evidence amendment changes only OpenSpec tasks and verification, preserving
those tested inputs. No production or test source correction was required.

The source report is `out/browser-pr-split/phase2/validation-report.txt` in
`t3-1b486911`. Inspected copies are retained under this worktree's
`out/browser-pr-split/phase2-passing/`; `phase2-validation.json` records source
identities, paths and all nine verified artifact-manifest hashes, including the
native package. Key retained SHA-256 values:

| Evidence                 | SHA-256                                                            |
| ------------------------ | ------------------------------------------------------------------ |
| Validation report        | `2f0754725af24c740021f636923ca0ae72826fb700e8dad3633c7d9b57ce3036` |
| Old-core backend JSON    | `56f4e18ebca7666b619dfa43565601a1db63edd8382e4e060b393f917cbe8823` |
| Passing backend JSON     | `518d116a27ba8eb0e17f0720134065dded329afd99cee50c88397931a5121168` |
| Passing storage JSON     | `4b651f4508b25676572a1ac05559cf8b592f51d27c98dbb7a4c1467472857b58` |
| Passing SDK browser JSON | `4564b791a8336bc2913107cdd79233653eaa1c23bfe6966eac5c531b732bb77d` |

This is composed behavior validation using the full original source for missing
prerequisites. It does not satisfy independent master-based PR2 delivery gates.
The evidence amendment reruns only proportional documentation/OpenSpec checks;
already-passed runtime tests are not repeated.

## Historical Afero-host browser acceptance procedure

After PR1 is available, run both targets in `test/browser_backend`:
`browser_backend_test` retains Node command/lifecycle coverage;
`browser_storage_test` runs the shared Go harness in pinned Chromium using
production `wasi.js` and Afero, without application SDK imports.
The latter retains `browser-storage.json` under its Bazel test outputs. It records
artifact hashes, browser identity, host operations, timezone offsets, persisted
note fields, closed database/media hashes, completed cases and failures.

Composed startup validation passed before the independent review correction.
PR195 is now merged and PR2 synchronized with master. Independent gates and
review passed as recorded above; historical composition results do not substitute
for those gates. Publication uses the final validated delivery receipt.

## Session ergonomics

One read-only research task checked the existing runner and protobuf patterns.
One analysis-only invocation established the prerequisite failure before any
heavy browser build. Initial worktree/catalog reads were overbroad, and a raw
formatter tail included a large path list; subsequent diagnostics used bounded
error extraction. Resume from this record, the scoped commit and local receipts;
do not repeat unchanged missing-prerequisite checks or copy prerequisite source.
