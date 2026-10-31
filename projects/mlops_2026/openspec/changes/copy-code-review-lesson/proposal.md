## Why

Keep the upstream code-review exercise available in the MLOps course project
exactly as requested, with a small CSV sample.

## What Changes

- Copy every non-CSV file from `lesson-08-code-review` into
  `projects/mlops_2026/code_review` without changing its bytes.
- Copy the first 100 lines of upstream `data/uber.csv`, including its header,
  without changing their bytes or adding data wiring.
- Preserve upstream provenance and exempt the copied files from general formatting.

## Capabilities

### New Capabilities

None. This imports teaching material without changing application behavior.

### Modified Capabilities

None; `skip_specs: true` records the absence of a behavioral delta.

## Impact

The new lesson directory contains the unchanged source from archive revision
`74dc45062f43f9b4abfbca738ffa6b56db4d5daa` and a 100-line CSV sample. Dependency files
remain unchanged. The supporting
`infra/src/openspec/changes/respect-vendored-python-exclusions` change makes
repository checks preserve the imported exercise as upstream-owned source.
