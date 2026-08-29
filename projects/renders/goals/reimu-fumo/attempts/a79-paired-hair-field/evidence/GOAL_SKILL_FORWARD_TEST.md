# Goal skill forward test

## Scope

This was a read-only, isolated reasoning exercise. No protected baseline,
candidate geometry, render, goal record, tracked source, or intended task
change was inspected or modified. This ignored report is the only persisted
artifact produced by the exercise.

## Verdict

The next work unit should test the dominant uncertainty: whether the isolated
geometry is visually plausible. Exhaustive hardening and promotion machinery
must not become prerequisites for obtaining that evidence.

The updated goal skill produced the controlling distinction: preserve source,
authorization, and safety invariants and run only the cheap checks needed to
create a safe disposable artifact *before* the artifact; defer exhaustive
correctness, reproducibility, and promotion gates until *after* the artifact
has survived direct inspection.

## Plan

- Bind the attempt to the exact protected-baseline identity, isolated candidate
  identity, criteria revision, and fixed render settings.
- State the hypothesis: "The latest geometry is plausible enough in the two
  most diagnostic views to justify further engineering."
- Use one likely-failure view and one complementary whole-form view, not merely
  flattering angles.
- Produce a clearly labeled, non-promotable two-view render in disposable
  task-local storage.
- Inspect it immediately. A categorical visual miss closes the attempt early
  with `refine` or `reset`; a plausible result earns the later validation work.
- Record the artifact identity, verdicts for affected visual criteria, dominant
  remaining defect, and next decision.

## Gates before the disposable artifact

- Preserve the protected baseline and keep the candidate isolated.
- Confirm rendering is authorized, non-destructive, and cannot overwrite or
  promote canonical state.
- Bind the exact candidate, inputs, cameras, and criteria.
- Run only cheap checks necessary for safe rendering: a parse/load smoke check,
  required dependency availability, and gross geometry sanity.
- Ensure the selected views can expose the dominant visual failure.
- Mark the output diagnostic and non-promotable.

## Gates after the disposable artifact

These gates occur only for a plausible survivor:

- Freeze the surviving representation and interfaces.
- Perform exhaustive topology validation and repair.
- Finalize and validate the schema.
- Run the full clean-reopen and reproducibility checks.
- Build and validate promotion automation.
- Produce an exact promotable result and rerun the complete acceptance and
  regression plan against that same identity.

A minimal isolated load check may occur before rendering because renderability
is a safety prerequisite; exhaustive clean-reopen validation belongs after the
probe. Likewise, only topology or schema checks strictly required to render
safely occur before the probe.

## Workstream decision

Keep this work unit sequential. The artifact takes about a minute and creates
an immediate review bottleneck, while parallel exhaustive work risks hardening
a geometry hypothesis that the render may promptly reject. Only preparation
that can affect the imminent review should run concurrently.
