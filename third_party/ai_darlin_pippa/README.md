---
title: Pipkin Pippa
description: FBX character mesh and matching diffuse texture from Darlin
---

Build `//third_party/ai_darlin_pippa:assets` to fetch the mesh and texture
from the [Darlin](https://www.darlin.ai/) CDN. The FBX matches the attachment
used in the earlier Pippa animation. The download provides the mesh and a
4096 × 4096 diffuse texture.

Import `Pippa.fbx` into Blender. In the Shader Editor, load
`Pippa_DiffuseColor.jpg` in an **Image Texture** node and connect its **Color**
output to **Base Color** on the mesh material's **Principled BSDF**. Rig the
mesh before retargeting an animation to it.

The supplied CDN files have no accompanying redistribution license. See the
[Phase Connect Fan Work Guidelines](https://phase-connect.com/fan-work-guidelines/)
for the character's use terms.
