# Reimu Fumo attempt 36

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 36 — sparse coupled outer-cage reset

**Candidate:** rejected disposable W0 preflight only. No persistent Blender
candidate was created. Review the
[wire/topology sheet](../../../../../out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/wire_preflight_sheet.png),
[machine metrics](../../../../../out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/wire_metrics.json),
and
[cage manifest](../../../../../out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/cage_manifest.json).

### Failure targeted and hypothesis

Attempt 35 used a `256 × 128` homothetic polar field with a high-valence crown
pole, analytic ring widths and depths, horizontal bands, and automatic rear
ownership. Attempt 36 asked whether one manually authored, low-cardinality,
non-radial outer cage could instead express the coupled brown crown and beige
face before any thickness, construction, material, or persistent candidate.

The test was deliberately narrow: one `13 × 9` nonuniform H-grid with exact
`117 V / 212 E / 96 Q`, Euler characteristic `1`, one `40`-edge exterior
boundary, and a shared `13`-vertex asymmetric hairline. Rows `0..3` owned
brown crown, row `4` the brown/beige junction, and rows `5..8` beige face. Rear
V4, locks, bow, body, and the complete rear profile remained unowned.

### Frozen authority and post-result corrections

- Canonical exact-variant registration is `Wh = 368 ± 4 px`, center `x=485`,
  crown top `y=231`, with central fringe tip `(.588,.677)` and `.101 Wh`
  viewer-right offset. The old `Wh=189` and tip `(.49,.57)` measurements are
  supporting physical-front-variant evidence only; they are not averaged into
  this canonical datum.
- The six canonical crown widths are
  `.408/.674/.837/.908/.948/.976 Wh` at
  `v=.052/.133/.242/.351/.459/.568`, each with `.03 Wh` tolerance.
- `.603 × .603 Wh` describes the visible beige pixel mask, not the inferred
  complete or hidden face opening.
- `1.098 Wh` is a lock-inclusive crown-to-lowest-lock composite regression.
  This lock-excluded cage ended at `v=.990` and was never entitled to use
  `1.098 Wh` as its own height.
- The execution contract reserved only `.25 Wh` for the independent rear
  leaf. The later all-reference audit corrects that band to `.35 ± .05 Wh`.
  Attempt 36 reaches only the obsolete `.25 Wh` witness and fails the corrected
  reservation; later attempts must not inherit the old value.

### Work and immutable evidence

The deterministic W0 builder emitted only four disposable files under
`out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/`:

| Evidence | SHA-256 |
| --- | --- |
| `wire_preflight_sheet.png` | `9ecc6042fe15daa1dc439aa06972ebf38b94211e83d08aeeed70459ffc68ec92` |
| `cage_manifest.json` | `34824e4e0440be7beeb07105457d56de06416f54899b9f759f8675ba53ca1aa7` |
| `wire_metrics.json` | `5ad2ea85c3ab85403001c5263feee0009ea3d93bae08259ee07a4f7f2715d967` |
| `build_w0_preflight.py` | `687d18bb335ed33d51f987fcd61953e5c3efd7052d764336e22103a009258b23` |

The protected tracked Reimu blend remained byte-identical at
`489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
The packet completed in `6.763 s`; it saved no persistent blend and did not
run G0 or G1.

### Machine result

Bookkeeping passes the declared exact topology: one component and boundary,
one brown and one beige component, consistent winding, exact hairline control
positions, and zero duplicates, degenerates, non-manifold edges, loose edges,
or reported 3D self-intersection.

The shape gates fail decisively:

| Gate | Result | Required |
| --- | ---: | ---: |
| Minimum scaled Jacobian | `-.097250` | `>.15` |
| Edge-ratio p95 | `9.503721` | `<3.5` |
| Edge-ratio maximum | `15.750680` | `<6` |
| Projected crossings, front / `-48°` / `+48°` | `4 / 52 / 77` | `0 / 0 / 0` |
| Front-contour RMS | `.116900 Wh` | `≤.03 Wh` |
| Front-contour maximum | `.142886 Wh` | `≤.05 Wh` |
| Level-2 hairline maximum error | `.043405 Wh` | `≤.03 Wh` |

All six measured crown courses are too narrow, with errors from
`-.086418 Wh` to `-.142886 Wh`. Visible beige exposure is
`.611 × .618 Wh`, which is locally compatible with the visible-pixel target,
and the explicit hairline controls are exact. Those isolated passes do not
offset an invalid outer form.

### Pixel review and acceptance result

The raw front is visibly too narrow and pointed. Its rectangular patch
collapses into a pole-like apex even without a literal polar pole. Both
three-quarter wires expose folded and warped flow instead of a fair plush
crown, while the side views show an abrupt peaked depth profile and an
artificial wall-like return. This is an absolute human veto independent of the
numeric failures. More rows, smoothing, or subdivision would hide the same
macro error instead of fixing it.

Criteria 1–8 remain failed or unverified. No reusable or animatable model was
produced, no visual acceptance gate passed, and no tracked asset changed.
Attempt 36 is an absolute W0 **NO-GO**.

### Process and approach audit

- **Improved:** the cheap wire-first gate reached a trustworthy rejection in
  `6.763 s`, before a candidate copy, downstream render packet, or model edit.
  Exact hashes make the verdict auditable.
- **Failed approach:** changing from a radial disk to a rectangular H-grid did
  not solve ownership or parameterization. The coupled single-surface premise
  still forced incompatible crown width, face depth, side return, and open-rear
  behavior through one chart.
- **Contract defect discovered:** `.25 Wh` rear reservation was itself stale;
  the controlling multi-view band is `.35 ± .05 Wh`. The active authority also
  needed to distinguish visible beige exposure, lock-inclusive composite
  height, and the supporting physical-front variant.
- **Whole-result check:** this cage is nowhere near an acceptable complete
  Fumo. Fast rejection is useful process evidence, not artistic progress and
  not permission to call a topology-only pass a model improvement.
- **Resource verdict:** Blender execution was not the bottleneck. The next
  iteration should spend its effort selecting a source-defensible ownership
  and construction representation, then use the same cheap gate. Parallel
  workers should review independent evidence or tooling, not produce stale
  variants of the rejected cage.

### Decision and exact next action

Retire the `13 × 9` H-grid, its row ownership, and its coupled rectangular
single-surface premise. Do not move controls, add rows, densify, subdivide, or
derive construction from it.

Before another Blender candidate, compare two bounded source-audited
construction witnesses:

1. a thin brown crown/fringe physically seated on a shallow stuffed support,
   with explicit contact and clearance so it cannot repeat the floating-card
   failures; and
2. a co-sewn sparse multi-patch outer assembly whose joins occur only at
   observed or explicitly inferred hidden boundaries, without inventing a
   visible central-bang seam.

Both witnesses must use the canonical front and controlling turn frames,
reserve `.35 ± .05 Wh` for the independent rear leaf, and expose front plus
bilateral three-quarter silhouette, depth, contact, and fold evidence. Freeze
one ownership/hidden-join contract only if it survives those gates. This is a
representation-selection reset, not permission to try another H-grid.
