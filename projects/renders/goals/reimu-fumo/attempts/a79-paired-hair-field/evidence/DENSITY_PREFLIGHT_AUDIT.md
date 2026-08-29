# A79 variant A — independent density/preflight audit

## Verdict

**REVISE / RESET. Do not construct the frozen handoff.**

The exact frozen component payload passes the narrow mechanical tests for
evaluated sampling, paired thickness, winding, manifold topology, and local
quad quality. That is useful evidence, but it is not a representation pass.
The producer's earlier `PASS` generalized outer-perimeter measurements to the
whole padded pocket, rewrote the authored inner skins before measuring them,
and did not enforce the root, receiver-independence, or leaf-width contracts.

The sparse `9 x 7`, `9 x 9`, and `7 x 7` authoring controls are **not** an
independent rejection reason. The literal frozen meshes have sufficient
evaluated density. Add local controls only where a registered reference or
whole-pocket convergence measurement demonstrates a form failure.

## Audited snapshot

This verdict applies to the final rejected Variant A snapshot:

| Artifact | SHA-256 |
| --- | --- |
| `geometry_preflight.py` | `b8b471b535e87fa8e725c32fb63747bc02dbd55fdf8c49b8566b41df772ab95e` |
| `geometry_preflight_report.json` | `1a76e455b41ba4f9a0c18412762fa417c213bef34a50d3de9a2aa7309c38aa74` |
| `frozen_geometry.json` | `ee9708c076bd443d59b2828c727b4ae85784c8043a87aebfdb8df896c0e44812` |
| `control_nets.json` | `a660ad4dfd022eb36fe6f83e8a3c7e0e36a5ce656d6a449c250e0cb8e33d39f2` |
| `TOPOLOGY_REPRESENTATION_SPEC.md` | `0a5d9ec39a9e6c2ec51c91828956cc6200c1539856e08fd834f3d653dc462d4e` |
| `RESET.md` | `c9e6706c12ac733bb3ab04c90f1edb75c3595870388383c7e4b5181096590e3d` |
| pinned interface contract | `590c363c8d8675623e77f36fdf5468982719b1d66a15eae599fb81a6e47b25c7` |
| pinned camera inventory | `43cf0cc3ab40737f5a3d7b7ec45c1da18dee489f7d0352faf983c5e8232f8d4c` |

I reran the exact `geometry_preflight.py` CLI into independent output paths.
It exited `3`, reported ten interface failures, and reproduced the rejected
frozen JSON byte-for-byte at SHA-256 `ee9708c...`. The independent report is
`density/independent_reset_report.json`.

The published report itself is not byte-reproducible by the documented CLI.
It contains manually added `source_hashes.parent` and
`source_hashes.tracked_asset` entries which the audited program never emits.
After deleting those two entries and the output-path field, the published and
independently generated reports are identical. This provenance defect does
not change the already-`REVISE` result, but a future handoff must make every
reported input an explicit CLI input or omit the non-generated fields.

## Independent literal-mesh results

`density/independent_frozen_audit.py` does not import the producer. It parses
the exact frozen JSON and independently checks counts and attributes, directed
edge winding, connectedness, Euler characteristic, volume, all three face
roles, outer and inner correspondence at both cells and vertices, and
projection through the five frozen cameras. Its exact result is
`density/independent_current_rejected_audit.json`.

### Evaluated density and paired thickness

| Component | Actual spans | Ordinary floor | Uniform 2 px fallback | Worst outer-perimeter span | Cell-centre distance | Minimum outer / inner normal ratio |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| front/crown | `204 x 99` | `136 x 66` | `204 x 99` | `1.939 px` | `2.083–3.464 mm` | `.858 / .836` |
| rear base | `155 x 125` | `94 x 76` | `141 x 114` | `1.974 px` | `2.486–3.636 mm` | `.876 / .852` |
| dominant leaf | `133 x 64` | `89 x 43` | `133 x 64` | `1.903 px` | `2.174–3.396 mm` | `.888 / .879` |

All corresponding-vertex distances also pass: front/crown
`2.074–3.465 mm`, rear `2.476–3.637 mm`, and leaf `2.144–3.397 mm`.
Both outer- and inner-skin vertex-normal ratios stay above `.8`.

This confirms that the separately evaluated dense skins avoid the rejected
closed-pocket Catmull-Clark failure. The old whole-pocket Catmull refinement
shrunk the rim toward roughly `1.1–1.4 mm` and produced widespread
misorientation; the frozen separate-skin payload does not.

### Exact topology and local quad quality

Each object is one consistently wound, closed genus-zero component. There are
no non-two-face edges, same-direction shared edges, duplicate faces, exact
duplicate vertices, or non-positive volumes.

| Component | Exact `V / E / Q` | Bridge angle range | Bridge max aspect | Global max adjacent-area ratio |
| --- | ---: | ---: | ---: | ---: |
| front/crown | `42,818 / 85,632 / 42,816` | `65.61–114.88°` | `2.326` | `1.861` |
| rear base | `40,992 / 81,980 / 40,990` | `64.22–119.01°` | `2.064` | `1.575` |
| dominant leaf | `18,602 / 37,200 / 18,600` | `71.32–110.38°` | `2.225` | `1.827` |

Thus the bridge omission in the producer's `strict_grid_quality()` is a gate
coverage defect, but it does not hide a local angle/aspect/foldover failure in
this exact rejected mesh. The bridge fails at the projected whole-pocket
level instead.

## Why the earlier pure `PASS` was false

### 1. The projection proof is not a silhouette proof

`projected_boundary_proof()` projects only `boundary_loop(rows, columns)` on
the outer skin. It omits:

- outer-skin interior points which can own a side-view silhouette on a curved
  sheet;
- the complete inner skin; and
- all three intermediate bridge rings.

The code then compares that same outer loop to a self-generated double-density
outer loop. Passing its span, motion, and chord gates proves only numerical
sampling of that parameter boundary. It proves neither the whole-pocket
silhouette nor reference fidelity.

Independent exact-vertex projection reproduces the two reset extrema beyond
the audited outer perimeter in the side camera: `2.720 px` for the leaf and
`1.090 px` for the rear base. The provenance is important: the leaf extreme
is inner-skin vertex `8710`, not a bridge vertex; the rear extreme is bridge
vertex `40432`. Therefore `RESET.md` attributes the leaf value too narrowly
to the bridge, but its broader conclusion is correct and stronger: the
outer-perimeter proof omits silhouette-owning parts of the pocket.

The reset also records `0.774 px` whole-pocket front/crown convergence in the
three-quarter view, over the `0.5 px` gate. That measurement is not present in
the runnable report, so it cannot be treated as independently reproducible
until a future preflight emits its samples and extrema. The structural scope
error above is independently established without relying on that value.

### 2. The measured inner skin is not the authored inner skin

The frozen path calls `aligned_inner_control()` for all three components
before evaluation. Direct measurement of that exact transformation gives:

| Component | Maximum control-point move | Maximum displacement rotation | Points changed |
| --- | ---: | ---: | ---: |
| front/crown | `1.187 mm` | `24.49°` | `63 / 63` |
| rear base | `3.593 mm` | `71.09°` | `81 / 81` |
| dominant leaf | `1.951 mm` | `49.05°` | `49 / 49` |

Consequently the excellent frozen thickness/alignment values validate the
corrected field, not the literal authored inner coordinates claimed by the
representation. This is not the old Catmull pinch, but it is still a false
authorship claim and a material form change.

### 3. The face receiver root is not the contracted loop

The exact frozen `A79_ROOT_RECEIVER_FACE_LOOP` has 369 indices, is open, and
its first and last points are `121.3899 mm` apart. It cannot correspond to the
contracted closed 196-anchor face-opening loop. A later interface-format
revision cannot repair those coordinates or that topology.

### 4. Described rear/leaf gates are not executable gates

The report records `plane_and_camber(panel)` measurements but never turns the
rear receiver-independence or camber contract into a failure. It likewise has
no executable narrowed-leaf width gate. Density cannot substitute for either
semantic requirement.

### 5. Current interface failures are real for the pinned contract, but are
not the reset basis

The exact current CLI correctly reports ten incompatibilities: bridge-only
vertices have `pair_id=0` and `skin=0`, skin faces are quads where the pinned
consumer specifies paired triangles, and one rear root group uses outer-skin
points. A revised scratch interface may intentionally permit ordered
bridge-only rings and quads. Even if it does, findings 1–4 remain and keep
this representation rejected.

## Required evidence for a successor

1. Measure the complete projected pocket, including silhouette-owning
   interior, inner, and bridge geometry, against both a doubled-density target
   and the camera-registered reference contour. Do not call a parameter-boundary
   check a silhouette check.
2. Preserve the separately authored inner coordinates, or explicitly specify
   and review any correction as part of the representation rather than hiding
   it inside a thickness gate.
3. Supply a closed face-opening root with exact registered correspondence.
4. Convert rear receiver-independence, camber, and leaf-width requirements into
   executable failures.
5. Reconcile the exact interface schema before claiming a frozen handoff, then
   run the pending self-intersection, correspondence-crossing, component-
   crossing, receiver-root, reference-mask, and clean-reopen gates.
6. Retain the current evaluated density floors. Do not increase sparse control
   counts unless the resulting whole-pocket/reference measurements identify a
   local degree-of-freedom failure.

No Blender process was run and no candidate or protected asset was modified
for this audit.
