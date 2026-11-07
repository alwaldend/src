---
title: Blender reference-fidelity evaluations
---

# Blender reference-fidelity evaluations

The offline Bazel target checks that the Promptfoo configuration, behavioral
cases, and staged skill load without credentials or a model call. The cases
cover honest comparison, incomplete references, user-directed delivery of an
imperfect result, and continued work toward an explicit quality goal.

A live target is omitted because representative evaluation requires reference
images, a Blender scene, rendered output, and visual inspection. Offline
validation does not prove that an agent follows the guidance or that a model
matches its reference.
