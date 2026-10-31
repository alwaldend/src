## Context

See `proposal.md`. The pinned archive contains seven lesson source files and
`data/uber.csv`. The user requires the imported source to remain unchanged.

## Goals / Non-Goals

Preserve every imported non-CSV file byte for byte and copy the first 100 CSV
lines, including the header, to the original path. Running or correcting the
exercise is outside this change.

## Decisions

- Keep the requested `code_review` layout and copy the first 100 lines of
  upstream `data/uber.csv` without changing their bytes. Add no lesson BUILD
  files or dependency wiring.
- Keep the archive pin unchanged. Retain its MIT notice with the import record.
- Mark the imported source as upstream-owned for repository formatting and
  import checks; build integration remains checked separately.
- Verify the copied-file inventory and hashes, the exact 100-line CSV prefix, and the
  absence of dependency changes. Existing quality and semantic lint gates
  verify the repository integration.
- For this pull request, obtain user approval before replying to comments or
  resolving review threads.

## Risks / Trade-offs

The exercise retains its upstream defects and old requirements intentionally.
The CSV is a partial upstream sample rather than the full training dataset.
