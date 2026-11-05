# Native manual head pre-save failure receipt

## File identity

The sole writer started from protected A157 at SHA-256
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
That protected file had the same digest after the session.

The writer made a clean pre-edit Save As at
`out/reimu_fumo_finish/attempt_015_native_manual_head/`
`a157_native_manual_head_working.blend`, SHA-256
`5f49968d81d965580020e7c3e5d3a8870213471f019e7e887ad300efe5770c36`.
No candidate was saved at the planned
`a157_native_manual_head_015.blend` path. The rejected live state was never
written and was discarded by reopening the pre-edit working copy.

The detailed sole-writer scratch receipt is `writer_diagnosis.md`, SHA-256
`e001518c56d84763a52d6d436cfe76a598bc9bf03aef4523bdc57d1d0db262c0`.

## Causal route exercised

Blender 5.1.1 was the foreground MCP authoring host. The named sole writer,
`native_head_artist`, used Blender operators and did not assign vertex
coordinates directly. The unsaved work:

1. applied the head-cushion lattice through Blender's modifier operator;
2. used native proportional Edit transforms for crown rise, front flattening,
   and nape taper;
3. created four grid stocks through Blender's mesh operator, comprising one
   crown/rear form, two side/forelock forms, and one central fringe;
4. shaped each stock with native Edit rotate, resize, and proportional
   translate operations, and applied Catmull-Clark natively; and
5. hid fourteen legacy hair objects without deleting them while leaving the
   body, costume, bow, face witnesses, lights, and cameras unedited.

One unsaved transform-order failure made the first grid about two metres tall.
The writer diagnosed the combined primitive transform as the cause, deleted
that unsaved grid, and rebuilt it with separate rotation, global resize, and
placement operators. This was a causal repair, not an equivalent replay.

## Automatic pre-save rejection

The complete front and right viewport state failed the plan's automatic
representation gate. It showed:

- a separate horizontal crown rail and side frame;
- a bulbous hourglass or armor-like central fringe;
- rounded-box side masses with stem-like roots;
- numerical overlap that did not read as sewn attachment; and
- parallel plate edges in profile.

Catmull-Clark inflated the stocks into pillows. Simple subdivision created
long crease spikes at bent Solidify and Bevel boundaries. A bounded
inverse-profile pass widened roots, narrowed ends, and reduced depth, but did
not remove the pillow, card, rail, or floating-root defects. The writer then
stopped as required by the plan.

Rejected-state topology, observed before discard, was 11,184 vertices on the
head; 2,385 on the crown/rear form; 1,617 on each side form; and 1,517 on the
fringe. These counts describe an unsaved rejected state and are not candidate
evidence.

## Review and acceptance state

The resetting reviewer was `native_head_artist`, role `sole model writer`.
`PROCESS.md` permits the writer to reset a candidate before retention. No
implementation-blind retain review was requested because no candidate was
frozen.

- Criterion 001 fails for this attempt: the complete state retained major
  silhouette and identity-form defects and no exact reviewable candidate was
  produced.
- Criterion 002 is unverified: there are no candidate bytes, calibrated
  renders, or candidate-bound measurements.
- Criterion 003 fails for this attempt: the complete state retained explicit
  pillow, card, armor, rail, and floating-root construction failures.
- Criteria 004 through 008 were outside this bounded attempt and remain
  uncovered.

There was no fixed-camera review render, no render-then-correct cycle, and no
technical claim about unsaved geometry. The working copy and detailed receipt
remain ignored task scratch in this worktree. No rejected binary or visual
artifact is promoted.

## Next causal action

Do not make another grid, shell, lobe, scalar-cut, or broad-transform variant.
Before another autonomous model attempt, prove on one disposable padded panel
that the live route can switch reliably among Grab, Smooth, Scrape or Flatten,
and a volume brush; isolate root and free-edge regions; deliver controlled
multi-angle strokes; form a smooth seated root, shallow plane, and tapered soft
tip; save and clean-reopen exact bytes; and undo natively. A skilled Blender
artist is the alternative resume condition.
