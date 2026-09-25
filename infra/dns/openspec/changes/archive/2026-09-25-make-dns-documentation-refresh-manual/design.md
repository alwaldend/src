## Context

The dump CLI already separates preview, `--write`, and `--check`. The coupling
comes from `TestGeneratedDNSPagesAreCurrent` in the ordinary repository test.

## Decision

Remove that freshness assertion and its unused dependency. Keep the ownership
test and generator fixture tests. Operators choose when to refresh or compare
the snapshots through the existing CLI; no replacement automatic job is added.

## Verification

Before implementation, appending a harmless marker to a generated page made
the repository test fail specifically at its freshness assertion. Repeat this
probe after the change and verify that ordinary checks pass without rewriting
the page. In a public fixture, explicit `--check` must reject stale pages,
`--write` must regenerate both views, and subsequent `--check` must pass.
Restore task-owned probes before delivery and retain repeatable inputs/results
under `out/split/dns-manual/`.
