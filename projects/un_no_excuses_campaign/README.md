---
title: UN No Excuses campaign
description: Editable Blender remakes of media from the UN No Excuses campaign
statuses:
  - in_progress
---

This project contains editable Blender scenes and rendered videos for the
campaign. The first scene remakes the supplied army spot; its stylized
reconstruction has approximate character likeness and motion.

Each scene has its own directory under `assets/` for its Blender file, render,
manifest, and scene-specific documentation. Scene generators live in separate
`cmd/<scene>_scene/` directories. The shared `cmd/encode_frames/` tool accepts
scene-specific frame counts, dimensions, timing, and final-frame holds.

The [army spot asset](assets/army_spot/README.md) describes the scene, media,
and render procedure. Maintained work is recorded in the
[project OpenSpec workspace](openspec/).
