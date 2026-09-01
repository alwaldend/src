# A93 preflight correction — scripted proportional editing is not the test

## Trigger

After A93 started but before the author produced geometry, an independent
read-only method audit found a materially stronger, executable discriminator.
The host has Xvfb, `libX11`, and `libXtst`, so a private display can deliver
genuine pointer events to pinned foreground Blender without calling
`bpy.ops.sculpt.brush_stroke`.

## Conflict

The initial A93 plan interpreted A90's "direct sculpt input" as sparse
Edit-mode proportional transforms. That would make broad vertex motion
auditable, but it would still test code-authored deformation rather than the
manual sculpt input path. A68/A71/A76 already show that transport and operator
success are not the missing capability.

The stronger discriminator is the existing A76 real-event protocol applied to
the exact A71 uniform macro-clay source. XTest supplies actual press, motion,
and release events; a pass-through modal sentinel proves provenance. Blender
code may prepare, project, and measure, but may not synthesize the sculpt
stroke.

## Disposition

Stop and reset this attempt before geometry mutation. Start a replacement
attempt using pinned Blender 5.2.1 in a private Xvfb and three isolated real
Grab drags with pairwise non-collinear intended directions. A single invalid,
local, zero-effect, or dimpled branch closes the modality without tuning.

No tracked asset, prior attempt, reference, or protected comparator changed.
