# Rebase conflict review: goal-skill evaluations

Observed on 2026-09-01 in
`/var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb`.
This is a read-only semantic review. It did not modify tracked files or Git
state.

The coordinator completed the rebase while this review was running. The
conflict-stage blob IDs below were captured before resolution, so they remain
immutable evidence even though `REBASE_HEAD` and the unmerged index entries no
longer exist. The final comparison was made against rebased `HEAD`
`054028a1ac7325cefc60f4d896d5ce1077dda69e`.

## Conflict inputs

| File | merge base | stage 2: upstream/onto | stage 3: rebased task |
| --- | --- | --- | --- |
| `evals/README.md` | `1fb6180b3dedd79fb5382055bd316ad741b1c3ca` | `0b255fff5c0b3ebf5c7fa27449a70edfa794752b` | `2866a1f7e79e3368268df194e7b4b7f162241543` |
| `evals/cases.yaml` | `113304dca54c8e6a41b5a95a9d450ef42cc12f3d` | `f629fdbcb29139a1ac30e775cf485a3810dc40f1` | `bf3ae0e7d71ea056f49d9d96f6370a4a6ce45831` |

In a rebase conflict, stage 2 is the fetched-base side and stage 3 is the
commit being replayed. The correct resolution is therefore not an arbitrary
choice between “ours” and “theirs.”

## Findings

1. Stage 3 is a strict semantic subset of stage 2 for both conflicted files.
   The seven additions below are byte-equivalent on both sides and must appear
   only once:

   - durable active-goal progress and pushing;
   - useful long-goal parallelism;
   - a concrete reason for sequential execution;
   - keeping process work subordinate to the result;
   - isolated copies for monolithic assets;
   - early closure of a falsified module; and
   - periodic full-history/process review.

2. Stage 2 alone adds the interruption/lifecycle distinction:

   - `Preserves an active goal across questions and additional tasks`;
   - `Honors an explicit lifecycle stop despite goal persistence`; and
   - the matching README sentence explaining that turn priority does not
     silently change goal lifecycle or authority.

   These are upstream behavior and must be retained.

3. Task-specific process lessons are not alternatives to those upstream
   cases. They cover distinct controls and should be added once:

   - `Renders a safe diagnostic before building promotion machinery` from
     `aef2f1cf257a7cc725a2bc969a0c68678c871fbc`;
   - `Detects stagnation without waiting for the user`;
   - `Improves the accepted baseline by an evidenced local delta`;
   - `Replaces a representation only after an evidenced local limit`; and
   - `Cancels a locally bounded edit that cannot reach the defect` from
     `054028a1ac7325cefc60f4d896d5ce1077dda69e`.

4. The superficially similar cases are not duplicates:

   - “Closes a falsified module” is a reactive early-close rule after decisive
     evidence; “Detects stagnation” places independent correction and a
     first-artifact veto on the live critical path before another long cycle.
   - “Performs a bounded full-history review” is a periodic retrospective;
     “Detects stagnation” requires an operational reviewer, latency measures,
     and a modality reset while those controls can still change the attempt.
   - “Keeps process work subordinate” establishes result priority;
     “Improves the accepted baseline” controls the exact baseline, smallest
     defect owner, reversible delta, and fixed before/after comparison.
   - “Renders a safe diagnostic” establishes staged validation before
     promotion; “Closes a falsified module” governs the decision after a
     diagnostic disproves an approach.

5. The nonconflicting `SKILL.md` context supports all of those cases without a
   contradictory or obsolete rule. The rebased file preserves:

   - result priority and interruption-safe lifecycle handling;
   - smallest-module work and early evidence boundaries;
   - evidenced deltas plus the causal-reach check;
   - staged feedback before promotion;
   - standing adversarial correction on the critical path;
   - inspectable and remotely preserved progress;
   - periodic whole-process review; and
   - bounded, useful parallelism with one canonical writer.

## Exact semantic union

Use the complete stage-2 files as the baseline. Do **not** concatenate or
re-append stage 3, because that would duplicate the seven shared additions.
Then make only these task additions:

1. In `evals/README.md`, retain the upstream interruption paragraph, then add
   the staged-feedback, standing-correction, delta-first, and causal-reach
   summaries once, before the live-target limitation paragraph.
2. In `evals/cases.yaml`, retain every stage-2 case in its existing YAML-list
   format and order. Insert the diagnostic and standing-correction cases after
   the full-history case, retain the two upstream interruption cases, and add
   the three delta/representation/causal-reach cases afterward.
3. Preserve the current auto-merged `SKILL.md`; no conflict-side replacement or
   extra duplicate instruction block is needed.

This resolution preserves upstream formatting and ordering, retains every
unique upstream behavior, and adds each task-specific lesson exactly once.

## Verification of the rebased result

Rebased `HEAD` already matches the recommended union:

- `SKILL.md` SHA-256:
  `7451364a2e58823a840231e3dbd169640a646da1813d5205e80aae8dccffb4dc`
- `evals/README.md` SHA-256:
  `b4ee8945374d08d6791e1e91d9a14074752d5b4da651fe2372d9663661e1ab16`
- `evals/cases.yaml` SHA-256:
  `a1826ca8fbdfd52ff434e2977534afdb16cfc12831739150b45883fb3c5545c2`
- `cases.yaml` has 23 descriptions and all 23 description strings are unique.

No tracked-file correction is recommended for these three files. The next
appropriate step is validation of the rebased skill package and eval config,
not another semantic edit.
