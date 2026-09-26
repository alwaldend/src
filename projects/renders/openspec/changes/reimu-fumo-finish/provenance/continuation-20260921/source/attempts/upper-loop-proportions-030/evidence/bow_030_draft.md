# Bow030 two-cage draft and unsaved evaluation

Helper: `bow_030_draft.py`; entry point `build_bow_030()` returns a
JSON-serializable receipt. Import defines functions/constants only. It never
opens, saves or creates a Blender file. Root now owns review and integration.

- Input: `hand_029b_candidate.blend`, SHA256
  `9ad353c57147831cd9440ec8ef7836f95dfb8c719da7f14fe1d122802f16f37d`.
- Executed helper SHA256:
  `cd2cdfb396304e573d937676f779688a172ba9fbfe8f23f479c5591085bd7556`.
- Pinned Blender 5.2.1 LTS, build `9e2066aef7ef`.
- Read-only probe: `bow_030_probe.py` / `bow_030_probe.json`.
- One unsaved helper evaluation: `bow_030_dryrun.py` /
  `bow_030_dryrun.json`. Source hash remained unchanged. No parameter retry,
  compensation, model save, material change or canonical write was performed.

## Exact scope

`TARGETS` is exactly:

```text
A154 Left loop macro cage
A155 Right loop macro cage
```

`AFFECTED_GEOMETRY` is exactly these five names for each Left/Right side:

```text
A42 Left/Right constructed bow loop
A42 Left/Right narrow gathered loop ruffle
A42 Left/Right white zigzag applique
A42 Left/Right root fold 1
A42 Left/Right root fold 2
```

The helper verifies the actual lattice-consumer set equals those ten objects.
It modifies only existing cage-point X/Z deformation coordinates. The head,
knot, tails, tail ruffles, global bow cage and all other geometry are controls.

## Space and construction decision

Both local cages are axis-aligned 4x2x4 KEY_LINEAR lattices. Their world X
domains are approximately -105.578..-1.839 mm and +1.889..+105.577 mm. Preserve
both innermost U columns, giving an unchanged knot-side interval through
approximately |X|=36.4 mm. This also preserves actual inner cloth witnesses,
not merely one empty control column. Move the remaining two columns with
weights 0.5 and 1 toward the center. The measured-extremum prediction requests
15.963665 mm inward at each outer control column.

Lower only the bottom two W rows with weights 1 and 0.25; leave the upper two
rows' Z exactly unchanged. The outer lower control requests 11.736511 mm
downward, yielding about 8.83 mm lower evaluated outer ruffle rather than
translating the bow. Existing soft cage padding is retained.

Order matters: local macro lattice precedes the Bow armature and shared
`022 bow proportion cage`. That shared cage is 2x2x64 KEY_LINEAR; its upper
Z mapping has slope about 0.613, while the lower control rows used here have
slope 1. The helper reads its actual sampled Z map and inverses the desired
final-world Z displacement through it. It then converts world coordinates
with each local cage's inverse matrix, retaining each original local Y
component exactly. No object transform, global bow cage or modifier is edited.

The loop ruffle's Solidify modifier follows the macro lattice. Its changed
surface normals therefore cause small evaluated thickness-offset changes
even with exact cage-Y preservation. These are measured below, not hidden
with compensating geometry.

## Measured unsaved result

Wh is the unchanged evaluated hood width 117.439255 mm, matching the preflight
datum rather than silently using the older nominal 116.5 mm.

| Metric | Before | After |
| --- | ---: | ---: |
| Upper loop plus ruffle span | 207.164839mm | 176.255435mm |
| Normalized upper span | 1.764017Wh | 1.500822Wh |
| Upper assembly height | 62.859401mm | 71.682677mm |
| Normalized upper height | 0.535250Wh | 0.610381Wh |
| Highest evaluated Z | 218.927681mm | 218.917549mm |
| Lowest evaluated Z | 156.068280mm | 147.234872mm |

The target hypotheses were 1.50 Wh span and 0.61 Wh height. This one unsaved
evaluation reaches them to about 0.00083 Wh and 0.00039 Wh respectively; visual
fidelity and contacts remain unreviewed.

Preservation evidence:

- All 76 non-target visible evaluated geometry controls had identical
  coordinate digests, including head, knot and tails. Global 022 bow cage
  data remained exact. The coordinator's writer should additionally retain
  its normal rig/material/topology guards.
- Both inner U columns and every cage Y component are exact. All 116 measured
  loop/ruffle/root-fold vertices within the knot's |X|<=7.75 mm span have zero
  evaluated displacement. No claim is made that the entire root-fold curve
  is fixed: its outboard portion intentionally follows the loop.
- Loop panels and zigzag/root-fold curves have zero evaluated Y change.
  Maximum ruffle-only Y changes are 0.064254 mm left and 0.061866 mm right.
- Highest Z shifts downward 0.010133 mm, consistent with post-lattice ruffle
  normal offsets. The original highest vertex also moves inward 13.629973 mm
  in X and 5.387 micrometers in Y. Preserving its full 3D location would be
  incompatible with the required narrower span; its height is effectively
  preserved but is not byte-identical.

## Review boundary

The next action is the coordinator's guarded save and fixed front/side/
three-quarter review. Check loop-knot continuity, loop/ruffle/zigzag seating,
the lowered outer edge against head and independent tails, and whether the
bow now reads proportionate without resizing the head. The dry run proves
bounded geometry changes and measured proportions only; it does not pass
visual likeness, clipping, motion, cloth physics or the overall goal.
