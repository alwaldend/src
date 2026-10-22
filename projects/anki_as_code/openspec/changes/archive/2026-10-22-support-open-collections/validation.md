## Verification

- `bazel_agent bazel test //projects/anki_as_code/... //:repo_quality_test`: 37 tests pass, including native App, collection-store, archive-builder and CLI workflows.
- `bazel_agent bazel build --config=lint //projects/anki_as_code/...`: passes for all 48 selected targets.
- Generated Go, Python and Java protobuf targets build with the renamed package. A source comparison confirms that only package declarations changed in the schema; message definitions and fields are identical.
- Dependency queries return no path from CLI or `ankiascodeapi` to the concrete App, `ankiclient` or `ankiserver`. CLI consumes `ankiascodeapi.App`; the typed application variable in cmd checks concrete App compatibility during compilation.
- The architecture SVG and WebP render through their owning Bazel targets. The rendered labels and path-source lifecycle edge were inspected.
- The borrowed workflow produces `borrowed-collection.colpkg` and `borrowed-state.anki.toml` in the App test's undeclared outputs. It checks explicit owner publication, source archive preservation, native identity, media, convergence, multiple collections, invalid sources, errors and cancellation.

Final candidate identity, validation logs and publication state are retained by the task-local delivery receipt under `out/support-open-collections/`.

## Resolved implementation findings

- `ANKIASCODE-GAZELLE-SCHEMA-RESOLUTION` (fixture-tested): the first generation pass guessed nonexistent `//anki` labels and omitted three upstream Go imports. Added project-owned proto and Go resolution directives for the already pinned Anki schemas; the next generation pass resolved them successfully.
- `ANKIASCODE-EXTERNAL-TEST-DEPENDENCY` (fixture-tested): the expanded test run found a missing direct API dependency in the external archive-builder test. Added it; that test passes and a scoped Gazelle diff is empty.
- `ANKIASCODE-TYPED-NIL-SOURCE` (fixture-tested): PR review identified a boxed nil connection that passed interface validation and panicked in gRPC. Expanded the existing invalid-source workflow before changing implementation, reproduced that RPC panic, and verified that nil source implementations and nil connection pointers now return validation errors before access.
- `ANKIASCODE-REBASE-WORKER-CONFIGURATION` (fixture-tested): rebased onto master at `a5c5e2f81174de9923bb1de469c610b4ebb31e04`. Five files conflicted in documentation, imports and test setup. Preserved trunk's explicit worker construction and linked CLI alongside source selection and schema renames; all 37 tests pass, including borrowed access on a zero-value App.

- `ANKIASCODE-REBASE-SERVER-BOUNDARIES` (fixture-tested): rebased onto master at `142825e94b7d2f5f756f137bec968b4993611d29` and reconciled eight conflicting files. Preserved trunk's complete CLI tree, injected native runner and expanded package boundaries. Added `ServerRun` and `ServerServe` to the shared `ankiascodeapi.App` contract, retaining the source API, schema names and caller-owned collection lifetimes. All 37 tests pass, including server dispatch, native hosting and borrowed/path workflows.

## Session ergonomics

The workspace and change identity are recorded in this owner-local change. Syntax-aware migration changed 33 existing collection call sites; other consumer edits are the protobuf rename. Initial missing-API failure evidence was captured before production implementation. An unnecessarily broad worktree listing and two corrected Cordis path requests added discovery overhead; future continuation needs only this change, the scoped source diff and the delivery receipt. No shared agent policy or host configuration was changed.

The public-interface follow-up reuses this change and PR #176. The guarded rebase stopped without changing the main worktree; three-way comparison established mechanical resolutions in five files before manual continuation. Fetched refs, failed-rebase evidence, range comparison and final delivery receipts remain in task-local scratch. Batched skill reads exceeded output budgets twice; future reads should be limited to the active phase and bounded sections.
