## Context

The implementation already uses protobuf-generated types directly for serialized data. Runtime metadata lives in model wrappers and snapshots.

## Goals / Non-Goals

Extract these packages unchanged and keep them independently buildable. This change does not implement CLI operations, TOML parsing, archive handling, or reconciliation.

## Decisions

Copy the exact API and model tree entries from the validated PR #128 candidate. The API owns serialized field definitions; model depends on generated collection messages and the existing upstream Anki note-type messages. Model contains no business logic or parallel serialized structs.

Use existing Bazel generation rules and dependencies. Verify generation, Go compilation, semantic lint, repository quality, and OpenSpec validation. The implementation PR continues to own end-to-end collection verification.

## Migration Plan

Publish this prerequisite against master and identify it in PR #128. After it merges, rebase PR #128 to drop identical API/model additions without removing progress. No stacked base or temporary consumer rewrite is needed.

## Risks / Trade-offs

The shared files appear in both PRs until the prerequisite merges. Exact source equality makes that overlap explicit and allows a subsequent rebase to remove it.
