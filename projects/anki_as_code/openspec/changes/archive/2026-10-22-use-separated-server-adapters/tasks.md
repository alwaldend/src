## 1. Adapter construction

- [x] 1.1 Select the explicit native adapter package and update its documentation/diagram; verify existing linked-worker and CLI E2E checks.
- [x] 1.2 Add a separate WASM executable construction path after the server adapter is usable; verify a cgo-disabled build and same-executable child startup.

## 2. Acceptance

- [x] 2.1 Run application collection/package E2E workflows with the WASM worker and retain native-readable round-trip artifacts.
- [x] 2.2 Validate documentation and OpenSpec and record verified handoff evidence. Git and delivery receipts own publication status.

## 3. Configured command selection

- [x] 3.1 Consolidate the command package and select private adapter variants through the shared backend setting; retain existing CLI and native/WASM E2E fixtures.
- [x] 3.2 Verify both configured public builds, cgo-free linkage, native package interoperability, executable names, documentation, OpenSpec, and required delivery checks.

## 4. Adapter-owned guest packaging

- [x] 4.1 Remove guest byte handling and copy/embed rules from the application command; invoke the WASM adapter's argument-vector API and retain configured backend selection.
- [x] 4.2 Verify native/WASM CLI and collection workflows, native reopen, standalone packaging and any review-driven simplifications, then complete required delivery checks.

## 5. Signal handling rationale

- [x] 5.1 Explain why worker hosting restores SIGINT termination and ignores SIGPIPE; preserve behavior and validate formatting and required delivery checks.

## 6. Standard build configuration

- [x] 6.1 Document -c opt for optimized WASM builds after removing the custom guest transition; verify configured command builds and existing CLI/collection E2E.

## Trunk engine composition integration

- [x] Preserve trunk client factories, native in-process clients and shared Rust dispatch while integrating the full WASM command.
- [x] Verify all four cgo/WASM flag combinations, qualified worker schemes, external-only failures, pure-Go linkage and native reopening of the large collection artifact.
- [x] Review the integrated candidate and configure required quality/lint gates for guarded publication.

Failure scenarios recorded before integration: unavailable engine schemes must fail clearly; a WASM-only executable must not link native code; enabling both engines must retain native in-process clients while its default worker uses WASM; disabling both must still support external clients and reject local hosting; explicit native interoperability fixtures must remain native independently of public flags. Client API renames and dispatch extraction must preserve collection lifetime, error trailers and process shutdown.

## Composition simplification

- [x] Replace nested per-attribute engine conditions with explicit variant settings and let Go build constraints select optional source files.
- [x] Verify all four configured engine combinations and existing native/WASM/combined CLI workflows without changing their capabilities.

Failure scenarios: optional imports must stay out of disabled engine graphs;
source build constraints and selected dependencies must agree; fixed native
interoperability fixtures must remain independent of public flags.

## Engine package paths

- [x] Update WASM imports, build dependencies and the package diagram for `anki_grpc/pkg/ankigrpc/wasm`.
- [x] Verify native, WASM-only and combined CLI workflows and all configured engine graphs after the path changes.

## Strategy diagram by dependency

- [x] Define the always-available, WASM-enabled and cgo-enabled strategy groups against command composition and client registration; retain the current cgo default and worker precedence.
- [x] Replace the factory diagram with one App/Server strategy diagram containing three dependency groups, maintain SVG and dark WebP targets, and document their availability conditions.
- [x] Inspect the combined light/dark output for readable labels and container boundaries, and validate the owner-local OpenSpec record. Delivery receipts own exact-candidate quality, lint and documentation checks.

Failure scenarios: diagrams must not imply WASM in-process support without its engine, cgo-free
combined builds, native subprocess selection while WASM is enabled, or a bundled
worker that can run without either engine. Keep subprocess strategies out of
the default group; TCP and Unix clients require no optional local engine.
Keep registration details out of the strategy graphs.

## Strategy preview refinements

- [x] Arrange the three dependency groups horizontally, use connection strings alone in strategy titles, and inspect the dark WebP preview. Leave the revision local for user review; publication is withheld.
- [x] Expand the WASM execution label to show compilation of the embedded module and the WASI `_start` entry into the Rust CLI, then inspect the updated local preview.
- [x] Stack the groups vertically in WASM, cgo, default order, show the subprocess app starting the wazero runner and the runner executing the server through WASI, and inspect the local preview.
- [x] Connect stdio directly to the WASM server, label the execution chain App → Run → WASM runtime (wazero) → Run → Server, restore plain incoming junction connectors, and inspect the updated local preview.
- [x] Use only `Run` for the WASM runtime-to-server label and inspect the local preview.
- [x] Remove execution `Run` labels while retaining the C ABI connection label, and inspect the local preview.
- [x] Label the App-to-WASM-runtime arrows `Run` for both strategies, name the in-process transport `In-memory stdio`, and inspect the corrected diagram.

## In-process WASM composition

- [x] Register `inprocess+wasm://` only in WASM-enabled executables, prefer cgo then WASM for generic in-process clients, retain default worker precedence, and verify the installed CLI API and collection workflows.
- [x] Document and render the new strategy and configure full publication as requested; Git and delivery receipts own its final outcome.

## Explicit worker engine selection

- [x] Register every compiled engine's qualified subprocess scheme and select that engine in the same executable; make both generic strategies prefer cgo and otherwise use WASM.
- [x] Verify installed CLI engine selection in native-only, WASM-only, combined and external-only builds, including an invalid WASM cache that must not affect explicit cgo workers.
- [x] Show qualified schemes alone in the diagram and document generic alias precedence and explicit worker selection in the README.

Failure scenarios before implementation: a qualified cgo worker in a combined
build must not silently execute WASM; a qualified WASM worker must not execute
cgo; unavailable engines and invalid selectors must fail; root engine selection
must preserve opaque Rust CLI arguments and existing signal handling; generic
subprocess and in-process selection must both prefer cgo as requested. Verify native selection
with an unusable WASM cache, and retain an installed-executable E2E artifact.

## Worker compatibility and reference corrections

- [x] Preserve standalone executable arguments in qualified subprocess addresses while selecting engines for same-executable workers, including executable aliases.
- [x] Define the root engine flag in the shared CLI tree so generated reference pages include it.
- [x] Add diagram canvas padding and inspect the rendered dark WebP and site SVG.
- [x] Verify the focused installed-command and documentation workflows and configure full CLI, client, collection, quality, lint and documentation gates for guarded PR delivery.

Failure scenarios before implementation: a standalone `anki_grpc serve` worker
must not receive the application's `--server-engine` flag; qualified workers
using the application's executable through a symlink must still select their
engine; generated server help must include the root selector; outer diagram
titles and borders must have clear space above and below them.

The focused installed-command engine-selection, shared server dispatch and CLI
reference workflows pass after reproducing both review failures. The inspected
SVG and dark WebP use 96 units of canvas padding. Delivery receipts own the
full exact-candidate gate results and publication state.
