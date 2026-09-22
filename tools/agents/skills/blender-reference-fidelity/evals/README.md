---
title: Blender reference-fidelity evaluations
---

# Blender reference-fidelity evaluations

This suite records the behavioral contract for revising a Blender asset when
fidelity to supplied references is the acceptance criterion. The required
offline Bazel target validates the Promptfoo configuration, referenced case,
and staged skill without making a model call.

A live target is omitted because a representative evaluation needs reference
images, a Blender scene, fixed-camera renders, overlays, and inspection of the
resulting pixels. A tool-free response cannot prove that geometry changed or
that the rendered candidate meets its visual gates. Promptfoo validation
therefore proves only that these evaluation assets load; behavior needs an
isolated Blender fixture with image-based review before a live target would be
meaningful.
