# Concrete A02 semantic and mechanics receipt

Pinned Blender 5.2.1 clean-opened candidate
`a9e7d470ae3de3236c1b808b087ee36191dcc357fcdda821340eeeb25fa84ccf`
without resaving it. The frozen parent is
`04dda70fc072c44e6eb0820eb0d159c6b9e6f08b552f049b0efa0916e9cf7b1c`.

The semantic snapshot diff contains exactly two allowed components:

- `Placeholder Warning Label.location[0]`: `0` to `0.2170832306`
- `Placeholder Warning Label.location[1]`: `-0.1150000021` to `0.1609133482`

The mutation contains one precomputed location assignment. No render-driven
second adjustment was made. Within the declared snapshot scope, camera,
timeline, gravity, rigid-body world, object ownership, geometry, materials,
animation, lights, world, proxy, impact marker, and replacement interface are
unchanged.

Mechanics reproduce A01 exactly:

- 24 fps, frames 1--72;
- held and kinematic through frame 12, released at frame 13;
- 0.845343 m sampled descent;
- contact at frame 22;
- minimum bottom Z of -0.000687 m; and
- 0 m center-motion span over frames 60--72.

Tracked Reimu and Sisyphus source hashes remain exact.
