## Context

The approved runtime, stock WASI SDK sysroot, LLVM integration and native adapter foundation landed in #182. The guest still needs source adaptations in already pinned Anki, FSRS, Chrono and Tokio dependencies. See proposal.md for motivation; the Anki gRPC change owns the runtime contract.

## Goals / Non-Goals

Keep dependency versions and repository identities unchanged. Preserve native behavior while enabling the approved offline WASI guest. SDK acquisition and executable composition remain owned by the merged prerequisite.

## Decisions

- Keep reproducible WASI patches beside the pinned Anki package. Compile networking out at target boundaries and return explicit unsupported errors rather than attempting an unapproved network bridge.
- Patch Chrono's WASI local-time implementation to query checked Go host callbacks. UTC fallback would change collection timing in non-UTC zones.
- Apply Tokio's existing unstable WASI opt-in only to WASI. Retain native crate features and versions.
- Regenerate MODULE.bazel.lock through the pinned module workflow. Generated metadata projects the owning annotations; do not edit repository pins independently.

## Risks / Trade-offs

- Native behavior could change through a broad patch. Keep adaptations target-gated and verify native CLI, transports and collection workflows alongside WASI E2E checks.
- Successful compilation cannot establish filesystem or time semantics. Verify real SQLite databases, offline collection/package round trips and historical timezone transitions; reopen exported artifacts with native Anki.
- The unoptimized large-data workflow exceeded its existing 15-minute deadline. Validate the full workflow with normal Bazel -c opt and retain the failed debug observation separately.

## Migration Plan

Build on merged #182, extract guest adaptations into a separate master-based PR, regenerate metadata and validate both native and WASI consumers. Keep executable/server integration and consumer visibility in #181; remove its dependency overlap after this prerequisite merges. Reverting this change removes guest annotations without changing the merged build inputs. Git and delivery receipts own the exact candidate and publication state.
