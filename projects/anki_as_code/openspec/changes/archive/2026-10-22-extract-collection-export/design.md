## Context

Export borrows a prepared connection to read native data and project configs. The native backend, config files, codec and state store are already available on master.

## Goals / Non-Goals

Export must accept a caller-owned connection, preserve its naming behavior and write files directly, and introduce no reconciliation or lifecycle dependency. CLI orchestration and planning/apply remain in #128. No SQL, new dependencies, schema changes or scheduling-preset edits.

## Decisions

Exporter.New(resources, states) borrows resource.Store and state.Store. Export reads a native snapshot from resource.Store and existing bindings from state.Store, assigns config paths, projects editable messages, and calls config.Export for file output. It only reads existing state for names and preset aliases; state generation, preview, serialization and output are separate operations. Runtime data remains private to export.

Export borrows resource.Store for generated clients and a fresh snapshot of upstream protobufs, sharing enumeration with reconciliation. Keep state and config projection outside that store. Do not wrap each native method or introduce generic CRUD models. The resource store is available on master after #161; this extraction adds export and its required generated rendering client.

New note labels use Anki's HtmlToTextLine RPC with media filenames omitted. Keep only filesystem sanitization, Unicode-safe length limits, fallback names and collision allocation in Go. Bound filenames bypass text conversion. Generate the rendering client from the existing pinned upstream schema; no upstream source or dependency versions change.

## Risks / Trade-offs

Sequential per-entity API reads do not promise an atomic concurrent snapshot: callers must provide exclusive offline access, as before. Native read and projection failures precede file writes. Filesystem errors may leave partial output. There is no output-tree transaction, cleanup or rollback; callers can export to a temporary directory and synchronize externally.

## Migration Plan

Publish this library extraction against master, independently of #128. Keep the identical prerequisite files in #128 until this PR merges, then rebase to remove overlap. No compatibility support is needed for the unreleased constructor.

## Verification

Write native E2E failure cases before implementation for exact fields/presentation, IDs only in state, namespaces, stable bound/new names, filtered-card exclusion, retained entries, ordinary overwrite behavior, partial output after filesystem errors and borrowed connection usability. Verify no generated state files and unchanged caller state files; emit TOML configs and native package artifacts. Also run the existing reconciliation E2E in #128. Delivery receipts own exact-candidate gates and publication evidence.

## Caller-owned state store

Exporter.New borrows the caller's resource.Store and state.Store for the same prepared connection.
The caller chooses the namespace when constructing the store for that collection.
Export constructs neither native clients nor stores and owns no lifecycle.

## Shared resource access

resource.New(connection) constructs generated native clients once. Each Snapshot
call reads fresh native protobufs, including unmanaged and filtered-card values.
Consumers retain config projection, planning, mutation order and state policy.
No cached snapshot, generic CRUD wrappers or backend lifecycle is added.
