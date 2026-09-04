# Explicit-vector bow checkpoint receipt

The writer opened a task-owned copy of protected A157. Its SHA-256 remained
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
The saved candidate is
`out/reimu_fumo_finish/attempt_013_tail_vector/a157_tail_vector_013.blend`,
SHA-256
`56d86d0caa316db21d9c87d8a19cab60d7f57cda1db35543c34f5fed6fe22073`.
It remains rejected task scratch and was not promoted as an accepted asset.

## Native edit receipt

In multi-object Edit Mode, deterministic source-space selection produced 132
of 650 red-tail vertices and 67 of 332 white-ruffle vertices on each side.
Blender's native proportional translation operator returned `FINISHED` for
explicit global-X values `-0.0267 m` and `+0.0267 m`, using connected topology,
smooth falloff, and radius `0.045 m`.

The selected extrema moved by the requested X vector within floating-point
tolerance. Protected roots moved `0.0 m`; Y and Z deltas were `0.0 m`. Only
the four planned red-tail and white-ruffle meshes changed. The white ruffles'
Solidify modifier states, all five fixed cameras, all non-target objects,
object names, linked libraries, and image dependencies were unchanged.

The pinned Blender 5.2.1 clean-open audit is
`out/reimu_fumo_finish/attempt_013_tail_vector/technical_audit.json`, SHA-256
`e8faf6b0ba5bc4c9376f608d313f592b0350c76e5bc04e99004fca697dc1bd17`.

## Frozen fast-pair evidence

Pinned Blender 5.2.1 rendered the candidate at 512 by 512 with the fixed front
and side cameras into a fresh directory:

- candidate front: SHA-256
  `dd3f50dfc3b494d3a73581f331e7f4fe8fd109e254be47d214b226851fb95d1f`;
- candidate side: SHA-256
  `6cc153fec885c845d64a1762f6167171eb281329a421e18221c6bcb12bd3e995`;
- candidate manifest: SHA-256
  `66722302085d667b57630c2430bda9da09708633254f72ec3735a334737d2269`;
- A157 baseline front: SHA-256
  `ed01afd0c9336f1fec3cbe6be44110e30d446476918e8c4c005b5a8fa8f3ae88`;
- A157 baseline side: SHA-256
  `f619709c29f234780d4a99a959c5498bc03ac1b30b90e65b2ebc9e8b3b2d2ba3`;
- baseline manifest: SHA-256
  `55d22768ab20df09dd4fca8c48fa784173d7324cc0f80d00e802aec7b0c8dc06`.

Aligned pixel measurement found head-and-hair envelope width `Wh = 247 +/- 1
px` and candidate bow span `505 +/- 3 px`, or `2.045 Wh` with conservative
range `2.024..2.065 Wh`. This is inside the specified `2.038 Wh +/- 0.05 Wh`
band. Bow-center offset was `-0.5 +/- 0.5 px`; upper-loop and knot silhouette
shift was at most one pixel. No pixel touched the frame, although the robust
left and right margins were only three and four pixels. The contract defines
no separate numeric crop-margin or lower-tail asymmetry threshold, so these
measurements are evidence rather than invented acceptance gates.

An implementation-blind paired reviewer returned `REFINE` and withheld the
remaining views. The reviewer found the enlarged tails nearly frame-filling,
taut, and horizontal; they read as wing-like fins instead of soft down-and-out
panels. Root transitions remained abrupt and the white edging was too thin to
read as a gathered ruffle. Symmetry, the upper loops, the knot, and side-view
attachment showed no meaningful regression.

The numeric span subtarget passed, but the full likeness and sewn-construction
criteria did not. In accordance with the fast-pair gate, rear and both
three-quarter views were not rendered.
