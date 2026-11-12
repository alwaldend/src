## Context

See proposal.md. The inventory command already reads project README YAML and
generates the central index. Hugo consumes landing front matter separately.
All 30 direct projects have READMEs; 27 have registered landings and the apex
site is separate. Existing status lists contain one value when present.

## Goals / Non-Goals

Keep one authored status per project, retain Hugo's list-shaped statuses field,
and enforce consistency through the existing checkout freshness gate. No
release certification, deployment, activity-age heuristic, or feature audit.

## Decisions

- A shared projects/statuses.yaml catalog owns identifiers and definitions.
  Generate its explanatory README table instead of duplicating the enum in code.
- README statuses own each project's declaration. Inventory update rewrites only
  landing statuses, marks the field as generated, and checks it on --check.
  Avoid a second Hugo-specific data channel or a new generator command.
- Retire in_progress and review its meaning per project. Active means substantive
  feature development is currently underway; unfinished work, historical feature
  additions, routine maintenance, and unchecked plans alone do not qualify.
  The user confirmed the active set on 2026-11-12: AL, Alwaldend.com, Terraform
  modules, Vast.ai Terraform provider, and X profile
  cleaner. Other previously active projects move to maintenance. Anki gRPC's
  historical WASM work does not establish current active development.
  Rename finished to completed and abandoned to archived. Completed preserves
  the existing completion declaration without implying current maintenance.
- Keep legacy taxonomy URLs as Hugo aliases on replacement term pages.
- Use the current isolated inventory E2E fixture for invalid declarations,
  stale projections, repeatability, and output artifacts, before implementation.

## Risks / Trade-offs

- Lifecycle is an owner declaration, not proof of health. Retain limitations
  and record review evidence; do not infer abandonment from old timestamps.
- Front matter surgery could damage content. Restrict rewriting to the top-level
  YAML statuses field, verify content in the E2E, and compare all remaining
  landing bytes against the baseline during candidate review.
- New projects must declare status before regeneration; errors identify the path.

## Migration Plan

Update declarations and generator together, regenerate twice, run quality and
site checks, inspect representative rendered taxonomy/index pages, and deliver
through the normal feature PR. Rollback is reverting the source change.

## Candidate and Next Action

Task baseline: dc1ca7025. The user-confirmed activity correction follows candidate
ac3e164e6 and is implemented and validated; see review.md and validation.md.
Local AL receipts own candidate identity and delivery continuation.
