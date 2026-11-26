---
title: Blender Demo Assets
description: Human base meshes and sample scenes for Blender
---

The [Blender demo file catalog](https://www.blender.org/download/demo-files/)
provides these modeling and scene assets.

| Bazel target         | Asset                                                                                  | License and credit                                 |
| -------------------- | -------------------------------------------------------------------------------------- | -------------------------------------------------- |
| `:human_base_meshes` | Human Base Meshes v1.4.1, with realistic and stylized meshes for sculpting and rigging | CC0 1.0, Blender Studio and community contributors |
| `:classroom`         | Classroom, a Cycles interior with reusable desks, chairs, and other linked props       | CC0 1.0, Christophe Seux                           |
| `:cozy_kitchen`      | Cozy Kitchen, the Blender 3.5 splash scene                                             | CC BY-SA, Nicole Morena                            |

Extract the human mesh bundle and add the extracted directory as a Blender
Asset Library. Extract Classroom and open `classroom/classroom.blend` with
its `assets` directory beside it. Cozy Kitchen opens directly as a `.blend`
file.
