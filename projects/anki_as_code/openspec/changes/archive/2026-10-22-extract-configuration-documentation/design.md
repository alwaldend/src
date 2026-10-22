## Context

Master already contains Load, Dump, Export, and the protobuf config model.
The published #128 source is bfe383772f0a983eabdc3f1c9e1b8a8de02ffc5d.

## Goals / Non-Goals

Keep this branch documentation-only and independently buildable. Collection
architecture, commands, state, and the pending local preview revisions stay
with their existing owners.

## Decisions

Copy only the committed configuration diagrams and their explanatory README
text. Declare site SVG and dark WebP targets under docs and add the existing
documentation consumer edge. Use master as base; retain matching content in
#128 until this PR merges, then let its rebase remove the overlap.

## Risks / Trade-offs

Temporary overlap can require README/build-file reconciliation. Preserve exact
diagram source bytes and scope the README insertion to the config API section.
Inspect rendered outputs and build the full site. Existing unpublished previews
remain untouched in the original worktree. Delivery receipts under
out/anki-config-docs-split own validation and publication state.
