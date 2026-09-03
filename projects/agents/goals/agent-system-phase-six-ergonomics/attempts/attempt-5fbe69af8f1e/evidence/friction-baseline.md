# Phase 6A friction baseline

Five representative task archetypes produced bounded FrictionRecord evidence.
The aggregate (`agent_system friction`) groups records by stable defect
signature and sorts by measured avoidable cost.

## Source records

The five records below are task-owned scratch evidence under
`out/phase-six-ergonomics/friction/`. Each aggregate cites the exact record
paths it consumed:

- `code-change.json` — code-change archaeology archetype.
- `goal-workflow.json` — goal workflow archetype.
- `catalog-regeneration.json` — catalog regeneration archetype.
- `skill-authoring.json` — skill-authoring archetype.
- `bazel-triage.json` — Bazel triage archetype.

## Aggregate

| Defect signature | Avoidable reads | Avoidable commands | Latency (ms) |
|---|---:|---:|---:|
| `goal-check-validation-error-suppressed` | 2 | 4 | 180000 |
| `catalog-updater-target-undiscoverable` | 2 | 2 | 60000 |
| `delivery-labels-missing-validation-mapping` | 1 | 2 | 45000 |
| `bazel-test-log-path-undiscoverable` | 1 | 1 | 35000 |
| `skill-validation-error-underspecified` | 1 | 1 | 30000 |
| **Total** | **7** | **10** | **350000** |

## Interpretation

The highest-cost defect (`goal-check-validation-error-suppressed`, 180 s) is
already traced by the Phase 5 learning proposal to delivered revision
`499bd74d`. The next-ranked defects are catalog-updater target discoverability
(60 s), delivery label-to-command mapping (45 s), Bazel test log path
discoverability (35 s), and skill-validation error specificity (30 s).

All five archetypes produced at least one record, satisfying the plan's
baseline acceptance condition.
