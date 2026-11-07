---
name: blender-reference-fidelity
description: >-
  Build or revise Blender models against supplied visual references. Use when
  resemblance to a reference matters, especially when visual iterations have
  stalled; do not use for freeform designs without controlling references.
---

# Blender reference fidelity

Use the supplied references to guide modeling and judge the rendered result.
The user's requested outcome, acceptance criteria, and delivery instructions
control the work. Do not claim an exact match when material differences or
missing evidence remain.

## Compare the model with the reference

1. Preserve the source and any accepted checkpoint. Identify the available
   reference frames or views and what each shows; note important uncertainty
   where the subject is obscured or a view is missing.
2. Render representative views and compare the pixels with the reference.
   Match framing and perspective closely enough for the judgment being made.
   Use landmarks, overlays, or fixed cameras when they help resolve a specific
   discrepancy; do not invent unavailable reference views.
3. Prioritize identity, silhouette, proportions, and visible contacts before
   fine surface detail. For animation, also check representative cuts, poses,
   timing, and transitions. Inspect other views when they can expose a likely
   regression.
4. When modeling is requested, make a coherent change and check its effect.
   Preserve useful checkpoints and record important findings. A short
   [scorecard](references/scorecard.md) can help with a complex or stalled asset; the
   [visual-review guide](references/visual-review.md) offers optional
   review prompts. For stuffed toys or fabric dolls, consult
   [plush construction](references/plush-construction.md).

## When progress stalls

If repeated edits do not improve the same major mismatch, reconsider the
modeling representation or work on another meaningful part of the task.
Do not repeat a failing approach indefinitely. An iteration limit is a cue to
reassess, not a reason to abandon a user who asked for continued improvement.
If the goal cannot be met, preserve useful work and explain the specific
remaining limitation.

## Deliver honestly

Inspect the actual asset and representative final renders for obvious defects
before delivery. Use independent visual review when it would materially help
and is available. Report material differences and uncertainty plainly, and
distinguish a working approximation from a result that meets an explicitly
requested fidelity bar.

Honor a user's request to deliver a reviewed but imperfect result. Do not
impose an additional numeric score, approval stage, or export prohibition.
If the user asks to keep improving until a stated quality level is reached,
continue meaningful work toward that level rather than presenting an unfinished
result as complete. Follow `repo-blender` for the repository's toolchain,
asset safety, and technical verification.
