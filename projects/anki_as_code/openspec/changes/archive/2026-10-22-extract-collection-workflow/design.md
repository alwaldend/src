## Context

See proposal.md. Master already provides the native process/collection backends, typed state store and codecs. PR #128 has a lifecycle workflow and state adapters that can operate without its reconciliation engine.

## Goals / Non-Goals

Extract the existing per-operation workflow and state integration, preserving its tested behavior. Keep collection export, plan, apply and build as master placeholders. Do not alter backend internals, state semantics, serialization, resource models or dependencies.

## Decisions

Retain app.runCollection and expandPath in internal/app. Compose BackendSTD and BackendCollection directly; the backend owns restore files and connection readiness. Stop the collection then the process with a bounded independent cleanup context. Invoke native ExportCollectionPackage only for changed archive output.

Make App.WithCollection own the scoped lifecycle instead of requiring separate Start/Stop calls. Construct namespaced stores directly with state.New inside its callback. Stores borrow the connection and require no cleanup, so no separate state-scoping or combined callback helper is needed. Keep application-owned collection/namespace options. Each state callback returns whether it changed stored bindings. Rename --input to --collection across the existing command scaffold without retaining an unreleased alias. CLI parsing, formats and streams remain in cmd.

Extract the existing native lifecycle and state E2E cases before production changes. Keep only their needed native fixture helpers; exclude export/reconciliation-specific helpers and tests. Include state API docs and its existing diagram. Leave #128 unchanged until this independent prerequisite merges.

## Risks / Trade-offs

Direct database errors can leave earlier committed writes; package errors before export discard restored changes. These are the existing native contracts and are explicitly covered by E2E fixtures. Some fixture code overlaps #128 temporarily; the later rebase removes that overlap.
