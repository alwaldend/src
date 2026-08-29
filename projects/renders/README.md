---
title: Renders
description: Reusable Blender models and rendered scenes
---

This project stores Blender source files for reusable models and rendered
scenes.

## Layout

- `blender/fumo/` contains Fumo models and Fumo-based scenes.
- `blender/misc/` contains Blender sources not specific to Fumo renders.
- `goals/<dns-id>/README.md` is the generated projection for each current
  versioned project goal.
- `docs/<legacy-name>/README.md` preserves unversioned goal histories imported
  into the current goal store. Treat those archives as historical evidence,
  not canonical writable goal state. Some archived prose links to ignored
  `out/` artifacts that were never durable; those links are provenance only
  and do not satisfy current acceptance criteria. Current goals begin with all
  imported criteria unverified and must retain or reproduce new acceptance
  evidence through the versioned goal record.
- `cmd/fumo_review/` contains the Bazel render-packet audit used to reject
  measurable Fumo regressions before visual review.

Each Blender asset or scene has its own directory beneath one of the Blender
categories.
