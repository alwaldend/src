---
title: Blender reference-fidelity evaluations
---

# Blender reference-fidelity evaluations

This suite records the behavioral contract for revising a Blender asset when
fidelity to supplied references is the acceptance criterion. It covers
rejecting unmeasured apparent progress, using a reversible pixel-visible delta
before escalating to a representation rebuild, and detecting repeated visual
stagnation before another whole-model generation cycle. It also requires a
proportional causal-reach check so a local sculpt is stopped when its safe
effect cannot create or materially approach the measured reference cue. The
required offline Bazel target validates the Promptfoo configuration, referenced
cases, and staged skill without making a model call.

A live target is omitted because a representative evaluation needs reference
images, a Blender scene, fixed-camera renders, overlays, and inspection of the
resulting pixels. A tool-free response cannot prove that geometry changed or
that the rendered candidate meets its visual gates. Promptfoo validation
therefore proves only that these evaluation assets load; behavior needs an
isolated Blender fixture with image-based review before a live target would be
meaningful.
