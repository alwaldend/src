## Context

Planning compares supplied configs, native snapshots and persisted bindings. Execution needs prepared native operations, not only a change report. The two capabilities communicate through the CollectionPlan protobuf.

## Goals / Non-Goals

Expose one pure internal/plan.Build function with required caller-supplied inputs and immutable output data. Keep lifecycle, capture, native mutation and persistence with callers and the executor. Do not simulate Anki validation, normalization or card generation.

## Decisions

- Build returns CollectionPlan with a report, expected stored state, desired bindings and ordered native operation groups. Runtime models, reference resolution, member matching, temporary-name allocation, moves and ownership checks remain private to internal/plan.
- Keep all shared data in the protobuf schema and reuse upstream native messages for operation payloads. New deck/type references resolve from Anki's returned IDs during execution.
- internal/reconcile executes the supplied plan in #128. It checks expected state, preserves native metadata, follows prepared phases, binds actual IDs and persists state. It does not resolve configs, capture resources or compare declarations.
- App composes capture, Build and Apply inside one collection scope. CLI/App planning output remains the report. Native operations and CLI integration stay in #128; this extraction targets master independently.
- Captured objects are indexed by native ID and desired objects by configuration path. New objects have no invented native IDs. Note-type and template-deck references in prepared operations use configuration paths. GUIDs remain native metadata. Field and template moves rename members inside one note type; cross-type transfers are rejected.
- Resource compilation compares values and emits typed operations and report entries together. Report strings do not control execution. Expected state and proposed bindings are distinct protobuf fields; bindings are projected once.
- Member IDs are authoritative when present. Otherwise original stored names identify current native members; stored ordinals only locate a rename's original state entry. Explicit renames retain their original member name for matching; ambiguous mappings fail before a plan is returned.
- Card assignments are emitted for known destination changes, missing explicit ordinals, or relevant native generation. Converged plans contain no native operations.

## Ownership and native validation

Only tracked omissions are deleted. Protect untracked deck descendants, note/type references, template deck targets and resident filtered cards, including their original home decks. Field/template names must identify unambiguous state-map addresses. Leave native name syntax, structural rules, template validation and generation to Anki.

## Risks / Trade-offs

Native effects remain deferred and reports are partial. Callers own snapshot freshness and exclusive offline access to the same collection/namespace. Stored-state checks reject plans whose bindings changed; they do not detect native content changes that leave bindings unchanged. Actual identity binding remains the executor's responsibility.

## Verification

Native-backed E2E loads exported TOML, captures values and closes the backend before planning. Verify report behavior, deterministic executable protobufs, immutable configs/snapshots/state and unchanged closed database bytes. Cover empty GUID metadata, explicit moves, invalid references, cross-type member transfers, untracked templates and filtered home references. Emit reports, executable plans and a native archive. #128 additionally executes a serialized plan without collection discovery and verifies dynamic IDs, stale rejection and convergence.

## Filtered cloze templates

All cloze card ordinals use native template zero. Filtered-card ownership checks
require that template's prepared original ordinal to remain zero, rejecting
replacement/removal while permitting explicit identity-preserving template moves.
Native fixtures cover resident cards at ordinals zero and one.
