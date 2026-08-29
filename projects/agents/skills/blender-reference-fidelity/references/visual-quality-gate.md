# Absolute visual-quality gate

Use this gate before presenting, exporting, or describing a reference-matched
candidate as successful. It catches technically valid models that remain
visibly unlike the controlling references.

## Prepare review evidence

Provide the reviewer with only:

- the controlling reference images, labeled by view;
- fixed-camera front, side, rear, and three-quarter renders;
- one uncropped presentation render at the intended approval stage;
- the current stage, such as neutral sculpt, materials, or pose.

Do not initially provide source code, object names, measurements, topology,
the previous candidate, the intended fix, or the author's own verdict. Those
details bias the reviewer toward implementation success and relative progress.

When an independent agent or image reviewer is available, use a context-light
review task and prohibit edits. If one is unavailable, perform the same review
from the rendered pixels before consulting implementation diagnostics. Never
ask the user to be the first reviewer.

## Review the candidate absolutely

Answer these questions from the images alone:

1. Would the candidate be recognized as the same subject and reference variant
   without a label?
2. Does it read as the intended medium at the current stage? A neutral plush
   sculpt must already read as constructed soft fabric rather than plastic,
   armor, a helmet, or disconnected primitives.
3. What are the five largest visible discrepancies, ordered by impact?
4. Are any identity-defining forms, attachments, contacts, or occlusions
   visibly wrong?

Score each applicable category from 0 to 10:

| Category                   | Evidence                                      |
| -------------------------- | --------------------------------------------- |
| Overall reference likeness | Same subject and variant without explanation  |
| Silhouette and proportions | Bounds, masses, spacing, and gaps             |
| Construction               | Parts join and deform as referenced           |
| Identity features          | Face, hair, accessories, or equivalents       |
| Contact and occlusion      | Clean depth order; no floating or clipping    |
| Intended-medium read       | Correct fabric, hard-surface, or anatomy read |
| Presentation               | Entire subject under diagnostic-neutral light |

Relative improvement is not a category. Do not compare with the previous
candidate until this absolute review is complete.

After the absolute review, make a second fixed-camera A/B comparison between
the candidate and the last accepted baseline. Hide implementation details and
randomize or neutralize labels when practical. This second stage controls the
internal keep-or-undo decision; it does not weaken the absolute approval gate.

## Pass or reject

An approval candidate passes only when:

- the unlabeled-recognition answer is yes;
- every applicable category is at least 8/10;
- no identity-defining discrepancy is classified as major;
- the reviewer cannot point to visible clipping, floating, accidental tangency,
  disconnected construction, or a wrong material read for the current stage;
- fixed-view checks show no regression hidden by the presentation camera.

One major visible failure overrides the average score. Technical validity,
measured landmark tolerance, majority votes, and improvement from a worse
baseline cannot override rejection.

For an internal intermediate cycle, record lower scores and continue working;
do not present it as an approval candidate. If the same category remains below
8 after two reviewed cycles, change the diagnosis or bounded local hypothesis
instead of mechanically repeating the edit. Replace the subsystem's
representation only when reviewed isolated edits establish a structural limit,
unless direct evidence makes the accepted baseline categorically nonviable.

## Record the verdict

Record:

- reviewer independence or context limitations;
- category scores;
- the ordered discrepancy list;
- pass or reject;
- the next smallest local correction if rejected, or the recorded structural-
  limit evidence and discriminating test when a representation change is
  justified.

Only after a visual pass should technical validation authorize export.
