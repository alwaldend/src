# Phase 6A friction-baseline plan

Instrument the 3-5 most common task archetypes to produce bounded FrictionRecord
evidence with stable defect signatures before selecting optimizations.

## Archetypes

1. Code change: edit source, validate narrowly, deliver through repo_delivery.
2. Goal workflow: checkpoint, close, validate, and regenerate goal catalogs.
3. Catalog regeneration: detect stale catalogs, run the correct updater, verify.
4. Skill authoring: create or update a skill, project the skill link, validate.
5. Bazel triage: run a failing target, diagnose, implement the narrow fix.

## Method

Record one bounded FrictionRecord per archetype using the existing schema:
stable task identity, selected operation, conflicts, avoidable reads/commands,
failed assumptions, latency, evidence tier, and defect signature. Baseline is
accepted when all archetypes have at least one record and no archetype shows
zero friction. Improvements are selected by measured cost, not anecdote.
