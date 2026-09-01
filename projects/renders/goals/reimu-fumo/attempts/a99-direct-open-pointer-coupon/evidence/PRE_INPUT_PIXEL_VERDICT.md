# A99 pre-input pixel verdict

## Verdict

`STOP`. The externally captured Blender window is valid, nonblack, mapped to
the exact expected X11 window, and visibly shows the complete target in the
right-orthographic Sculpting viewport. It also exposes a decisive live-state
contradiction that invalidates input.

The READY-bound plan reports native Grab as:

```text
locked_size = SCENE
unprojected_radius_m = 0.05000000074505806
strength = 0.4000000059604645
```

The actual mapped Blender pixels show the active Grab UI as:

```text
Size = 100 px
Size Unit = View
Strength = 0.400
```

The visible brush cursor is also consistent with the smaller view-locked
radius rather than the plan's computed 260.42 px scene radius. READY therefore
attests a data-block state that is not the active live tool state receiving a
pointer gesture.

The mismatch is acceptance-visible in
`harness/run/evidence/external_before.png`, SHA-256
`fb13cc11e7308f85607ac26b021f03b32eb0f1f9c99add1fe5ae649daf88d9b2`.
The external receipt binds mapped window `2097154`, dimensions
`2560 x 1440`, and nonblack RGB range `0–255`.

A99 permits no patch, relaunch, or input after any pre-input mismatch. The
inject action was never called.

