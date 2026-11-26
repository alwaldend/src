---
title: CMU Breakdance Motion Capture
description: Three Blender-importable BVH clips from CMU subject 85
---

These clips come from subject 85 of the
[CMU Graphics Lab Motion Capture Database](https://mocap.cs.cmu.edu/),
converted to BVH by Bruce Hahne and published as individual downloads in
[Shriinivas/cmubvh](https://github.com/Shriinivas/cmubvh/tree/3bc9c75c481ac6e622f3c378369b293d8de9be45/Sequence-081-085/85).

| Bazel target       | Clip                    |
| ------------------ | ----------------------- |
| `:handstand_kicks` | 85_05, HandStandKicks   |
| `:helicopter`      | 85_08, Helicopter       |
| `:long_sequence`   | 85_12, LongSequenceGood |

Import the resulting `.bvh` file with Blender's Motion Capture import. Each
clip includes its own animated skeleton. Retarget the animation to the
destination character's rig.

See the [CMU usage terms](https://mocap.cs.cmu.edu/) and the
[BVH conversion usage terms](https://sites.google.com/a/cgspeed.com/cgspeed/motion-capture/the-motionbuilder-friendly-bvh-conversion-release-of-cmus-motion-capture-database/readme-file-for-the-bvh-conversion-release).
CMU requests this acknowledgment in published work:

> The data used in this project was obtained from mocap.cs.cmu.edu. The
> database was created with funding from NSF EIA-0196217.
