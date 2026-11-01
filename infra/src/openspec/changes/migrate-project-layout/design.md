## Context

See proposal.md for scope. The clean linked worktree uses
`t3/migrate-project-layout`, initially at the supplied repository revision.
The migration is repository structure work owned by `infra/src`.

## Goals / Non-Goals

Source roles become visible in paths while component behavior and dependency
pins remain unchanged. Upstream source, historical change records, and
ecosystem-required layouts are outside the relocation inventory.

## Decisions

- Move command packages to `cmd/<command>` and implementation to `internal/`.
  Move schemas to `api/`, Starlark to `rules/`, resources to `assets/`, and
  patches to `patches/`. Keep frontend source in `internal/frontend`.
- Keep the distributable `autoscroll` Python package under `pkg/autoscroll`
  with its implementation named `internal`; retain package and wheel semantics.
- Migrate product and tool components together so repository consumers can use
  canonical paths directly. The user explicitly excludes compatibility aliases.
- Rewrite fully qualified labels, Go and Python imports, component-relative
  paths, nested-module labels, and relative documentation links together.
  Preserve archived OpenSpec evidence as historical references.

## Risks / Trade-offs

- Implicit target names and Go import boundaries can break after moves:
  inspect generated BUILD changes and compile affected consumers.
- Relative resources, Gradle paths, and wheel contents can drift: inspect
  representative outputs and source-set declarations.
- Abandoned components may already have missing build inputs: preserve failure
  evidence and distinguish pre-existing defects from relocation errors.

## Migration Plan

Inventory tracked inputs and record the old-to-new mapping in ignored scratch;
move files with collision checks; update consumers; run formatting, graph,
build, test and lint checks; inspect representative outputs; deliver through
the guarded repository workflow. No deployment is involved. Reverting the
task commit restores the previous layout.

## Candidate and Continuation

Initial revision: `628ce0ab8209582ffc17347819843d8bc81aea41`.
Logs, the 575-file migration mapping, artifact manifests, and delivery receipts
belong under `out/project-layout/`. See validation.md for acceptance evidence.
Next action: run the guarded exact-candidate delivery plan and publish the
feature pull request when its gates pass.
