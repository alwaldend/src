# Bow 030 macro-proportion preflight

## Scope and evidence binding

Read-only preflight. No Blender datablock was saved and no model, helper, or
goal file was changed.

- Candidate: `hand_029_candidate.blend`
- Candidate SHA-256:
  `93b2c163bff3e0a6e0d478b7b2312dd8578610d618b3570140cdf76ca400f386`
- Available five-view directory: `hand_029_eevee_review/` (the supplied name
  `all5hand_029_eevee_review/` is not present in this worktree)
- Front render SHA-256:
  `5e36e8209c01f78de8e39d131560c2186244a4c86a4478fa256c2ecc4e585ac0`
- Render contract SHA-256:
  `4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`
- Canonical front SHA-256:
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`
- Physical front SHA-256:
  `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51`
- Physical side SHA-256:
  `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`

The canonical front controls exact front proportions. The physical images
only support construction and pose-variation interpretation. The candidate
render is orthographic; the canonical photograph has residual perspective.

## Verdict

The evidence supports a loop-only correction and does not support resizing
the head. The upper bow loops are materially too wide and modestly too
shallow. The independent tail span and highest bow point already match their
controls.

## Measurements

`Wh` is the outer head-and-hair width. The reference datum is `368 +/- 4 px`.
The evaluated candidate hood width is `0.117439255 m`, about `206 px` in the
fixed front camera.

| Measure | Canonical control | Candidate 029 | Delta |
| --- | ---: | ---: | ---: |
| Upper-loop outer span | `1.50 +/- 0.03 Wh` | `1.7640 Wh` | `+0.264 Wh` |
| Upper-loop height | `0.61 +/- 0.03 Wh` | `0.5353 Wh` | `-0.075 Wh` |
| Highest bow point above crown | `0.22 +/- 0.01 Wh` | `0.2142 Wh` | about `-0.006 Wh` |
| Complete bow/tail span | `2.038 +/- 0.05 Wh` | `2.0065 Wh` | `-0.0315 Wh` |

Canonical front endpoint reads for the upper silhouette were approximately
`x=209..760 px` and `y=150..375 px`. The lower loop edge is partly occluded,
so its vertical value carries the larger uncertainty. The candidate's exact
evaluated upper loop plus loop-ruffle bounds are:

- X: `-0.103582658..+0.103582181 m`
- Z: `0.156068280..0.218927681 m`

The fixed render gives the same conclusion: the red-only upper mask spans
about `354 px` (`1.72 Wh`), while the evaluated white outer ruffle brings the
full upper silhouette to `1.764 Wh`. This difference is much larger than
raster, endpoint, or reference-perspective uncertainty.

The physical front has a different soft pose and mild perspective. Its upper
loop silhouette is roughly `1.4-1.55 Wh`, supporting the canonical result but
not replacing it. The oblique physical side cannot control frontal span; it
only argues against changing panel depth in this unit.

## Exact upper assembly and interfaces

The bounded upper-loop dependency sets are driven by two `4 x 2 x 4` local
lattices:

- `A154 Left loop macro cage`
- `A155 Right loop macro cage`

Each local cage drives its side's:

- `A42 Left/Right constructed bow loop`
- `A42 Left/Right narrow gathered loop ruffle`
- `A42 Left/Right white zigzag applique`
- `A42 Left/Right root fold 1`
- `A42 Left/Right root fold 2`

The loop and loop-ruffle meshes are parented to
`Head_Cushion_Manual_Target`, weighted to the `Bow` bone, and use the
`ReimuFumoRig` armature. The zigzag and root-fold curves have no armature
modifier; their common rest-pose attachment is the matching local cage plus
`022 bow proportion cage`. This is why a loop object transform would be the
wrong shared interface and a local-cage deformation is the coherent one.

The following are independent of the local loop cages and must remain exact
controls:

- `A42 Left/Right independent draped bow tail`
- `A42 Left/Right narrow gathered tail ruffle`
- `A42 flattened gathered center tie`
- `022 bow proportion cage`
- all head geometry

Current knot bounds are `x=-0.00775..+0.00775 m` and
`z=0.191650..0.202538 m`. The loop ruffles reach inward to about
`x=+/-0.0039 m`, so they overlap the knot in projection rather than leaving a
root gap.

## One bounded correction

Deform only `A154 Left loop macro cage` and
`A155 Right loop macro cage`, with zero displacement at the inner/root column
and a monotone increase toward the outboard column:

1. Bring the evaluated outer ruffle extrema from about `+/-0.10358 m` to
   about `+/-0.0881 m`. This is a root-relative horizontal factor of roughly
   `0.84-0.85` and yields an upper span near `1.50 Wh`.
2. Keep the current highest evaluated point, `z=0.218928 m`, fixed. Lower only
   the outboard lower region by about `7.5-10 mm` so evaluated loop height
   reaches `0.60-0.62 Wh`.
3. Keep cage Y coordinates unchanged. Keep the knot, tails, global bow cage,
   head, and inner root witnesses unchanged.

This is not a uniform bow scale. A uniform or global-cage scale would destroy
the already-passing complete-tail span and top/crown relation.

## Acceptance and falsification

The first corrected front render should satisfy all of:

- upper-loop outer span `1.47-1.53 Wh`;
- upper-loop height `0.58-0.64 Wh`;
- highest bow point above crown `0.20-0.23 Wh`;
- complete bow span remains within contractual `2.038 +/- 0.05 Wh`;
- head envelope, knot, tails, and tail ruffles have zero unintended evaluated
  displacement;
- no new root gap, clipping, or accidental tangency above `0.02 Wh` in side
  or either three-quarter view.

Falsifier: if a local-cage-only test reaches the upper span target but the
head still reads narrow in an implementation-blind front review, then the
perceptual complaint is not explained by upper width alone. At that point,
measure head-to-body controls directly; do not infer a head resize from the
same complaint.
