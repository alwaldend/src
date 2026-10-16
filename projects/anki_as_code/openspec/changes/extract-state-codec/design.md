## Context

The state store and TOML/JSON implementation already exist on master. PR #128 needs state projection and moves for export, planning and apply. See proposal.md for extraction scope.

## Goals / Non-Goals

Keep state construction, moves, ownership validation and persistence behind Store. Preserve storage, encoding, CLI options, collection lifecycle and native behavior.

## Decisions

Keep the unchanged codec in internal/codec and update existing master callers. Keep both PRs against master with required overlap until this prerequisite merges.

Store.Preview accepts parsed configs for explicit moves/source checks and an optional collectionpb.Resources protobuf. Nil resources retain the moved stored bindings. A supplied message replaces the entire resource set, even when empty; individual maps have no special presence semantics. Store clones inputs and returns detached snapshots without writing. Zero identities are allowed for planned creations, but persisted writes reject them. Nil map entries are invalid, including unresolved resources.

Callers construct Resources inputs using private native-to-protobuf conversions. There is no public Bindings request or raw mutation helper. Store owns state moves, ownership validation and persistence; Import privately projects native members. Update and Replace remain the external-state write API. Preserve native member identities/ordinals and input isolation.

Decision: reuse the existing protobuf instead of duplicating its structure in a public Go type. Whole-message replacement avoids protobuf map-presence ambiguity. Both PRs remain against master with required overlap until this prerequisite merges.

## Risks / Trade-offs

- Preview must remain read-only for plan/export: verify native state before/after, input/snapshot isolation, move conflicts and rejected provisional-state writes in native E2E before implementation.
- Preserve explicit move and member identity behavior: reuse the current move implementation inside state and run #128's native workflows.
- Keep master-based PRs buildable: #128 retains the required prerequisite overlap until merge, then rebase removes it.
