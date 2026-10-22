## Why

The merged configuration API can be documented independently of collection
reconciliation. Extract the configuration model and operation diagrams from
PR #128 into a focused documentation change.

## What Changes

- Describe the existing Load, Dump, and Export Go APIs with their process chart.
- Show configuration fields, path references, and filesystem locations in a class diagram.
- Package light/dark site SVGs and dark WebP previews through maintained Bazel targets.
- Keep collection architecture and reconciliation documentation in #128.
- Preserve the three existing unpublished preview edits in the original worktree.

## Capabilities

No behavior or requirements change; this is a documentation extraction.
`skip_specs: true` applies.

## Impact

Only projects/anki_as_code documentation, its build wiring, and this change
record. Both PRs target master. Config behavior, collection code, dependencies,
and data remain unchanged.
