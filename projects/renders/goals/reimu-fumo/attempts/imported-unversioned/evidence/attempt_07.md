# Reimu Fumo attempt 7

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 7 — multi-view cushion and thin folded bow panels

**Candidate:** planned `out/reimu_fumo_attempt_007_folded_panels.blend`,
SHA-256 pending, neutral construction stage, review packet
`attempt-007-multiview-folded-panels`.

**Failure targeted:** Attempt 6 matched the front trace but used constant depth
for the head and inflated each bow silhouette into one pocket. The fixed side
and three-quarter views therefore contradicted the controlling construction
references.

**Hypothesis:** A head cage whose front and rear depths vary by row according
to the measured turn-side contour, combined with thin traced bow faces and
separate return layers, will retain the improved front likeness while removing
the box, visor, bulb, fin, wing, and perched-root reads.

**Plan written before implementation:**

1. Start from factory-empty data and copy no rejected mesh. Retain only the
   frozen front and side traces, camera contract, and proven all-quad topology
   assertions.
2. For every head cage row, intersect the normalized turn-side polygon at that
   height and use the resulting front and rear depths for the structured front,
   rear, and gusset surfaces. Keep the exact physical-front outer contour and
   evaluated `1.00 × 0.74 × 1.03 Wh` dimensions, but let the crown and lower
   cushion roll inward in side view.
3. Keep the face as a region of the unified cushion. Remove the broad separate
   fringe visor; express the traced irregular hairline as the face-region
   boundary and a subtle seated seam. Rebuild asymmetric cheek locks as short,
   broad panels that wrap around the cushion sides instead of hanging as front
   slabs.
4. Rebuild every bow loop from a thin visible traced panel and a separate thin
   return layer. Use `0.012–0.020 Wh` paired-surface thickness, no inflation,
   low-amplitude authored fold valleys, two or three root pleats, and at most
   one smoothing level. Keep the measured front outline unchanged.
5. Place the return layers at distinct measured depths, fold their outer edges
   back toward the visible panels, and compress all loop and tail roots inside
   one small band knot embedded into the crown slope. Require hair to occlude
   bow cloth in their overlap band.
6. Render calibrated front, fixed three-quarter, shaded side, side silhouette,
   aligned 42% front overlay, and grazing root close-up. Reject before body work
   unless all reviewed categories reach 8/10 with no box, visor, fin, wing,
   inflated-pocket, flat-panel, or tangent-root description.

**Work performed:** Built a factory-empty component with the exact traced
physical-front perimeter and a row-varying front/rear depth cage derived from
the turn-side contour. Kept the evaluated head at
`1.000 × 0.740 × 1.030 Wh`, removed the separate broad fringe visor, and added
shorter wrapped cheek locks. Rebuilt the visible bow loops and tails as thin
traced all-quad surfaces with separate return layers, low-amplitude valleys,
compressed roots, small embedded knot, and outer fold bridges. Produced fixed
front, three-quarter, side, grazing-close-up, overlay, contact-sheet, and
silhouette evidence. No rejected geometry entered the tracked asset.

**Evidence:** Candidate SHA-256
`b784a5190a9089bc0b997759934e5045d83dfecc5858743328ca37b8cc264d23`.
The evaluated head is exactly `1.000 × 0.740 × 1.030 Wh`, and overall bow span
is `1.693 Wh`. Every mesh has zero boundary edges, zero non-manifold edges,
and zero non-quad faces. The diagnostic head has 65,536 all-quad faces after
applied subdivision. Each base loop, tail, and return panel has 126 vertices
and 124 quads; each fold bridge has 70 vertices and 68 quads.

The fixed side image confirms that mapping both controlling contours produces
a visibly softer crown and lower-head roll than Attempt 6. The same render
also shows that the candidate remains a deep rounded helmet around an inset
face. The fixed front and three-quarter images show one capsule loop, one tall
fin loop, pill-like tails, thick transverse fold bars, and no natural gathered
convergence at the knot.

The implementation-blind reviewer narrowly recognized the same Reimu variant.
Scores were head likeness 6/10, bow likeness 3/10, silhouette 4/10, sewn
construction 3.5/10, contact and occlusion 4/10, and presentation 7/10. The
reviewer reported automatic-failure helmet, fin, wing, inflated-pocket,
flat-panel, and tangent-root reads.

**Criterion results:** Exact head bounds, the controlling front and side
outlines, and mesh-manifold checks pass their local tests. Reference likeness,
complete measured silhouette, plush construction, and presentation quality
still fail the absolute image gate. Reusable structure, animation readiness,
and full technical integrity remain unverified. Repository delivery still
applies only to the rejected migrated baseline.

**Decision:** Reject at the coupled component gate. Preserve the calibrated
front and side contour evidence, but do not reuse the head, hair, bow, or fold
geometry in the body or tracked asset.

### Progress and approach audit after attempt 7

- **Improved:** head likeness rose from 5.5/10 to 6/10, and the fixed side
  contour is visibly less rectangular while preserving the exact normalized
  head bounds and front trace.
- **Regressed:** bow likeness fell from 3.5/10 to 3/10 and sewn construction
  fell from 4/10 to 3.5/10. The independent return layers did not produce a
  folded ribbon, and rigid outer bridges introduced a new bar-like failure.
- **Unchanged:** silhouette remained 4/10, contact remained 4/10, and
  presentation remained 7/10. No acceptance criterion passed as a whole.
- **Absolute result:** the component remains a visibly artificial assembly,
  not a convincing sewn Reimu Fumo. Better contour matching is useful evidence
  but cannot justify continuing its rejected representation.
- **Approach evidence:** simultaneous front and side constraints are retained.
  Thin independent layers collapse edge-on; connecting them with a rigid
  bridge produces a hard bar. A beige material region inside one brown cushion
  still reads as a face mask set into a helmet even when its boundary is
  traced accurately.
- **Representation reset:** use one exposed beige face cushion as the stuffed
  structural head. Add brown crown, rear, fringe, cheek, and nape pieces as
  thin overlapping fleece panels that follow the cushion without surrounding
  the face as one shell. Build each bow loop as a continuous folded ribbon and
  use a real gathered root; do not reuse paired traced pockets or rigid fold
  bridges.
- **Highest-leverage problem:** prove one natural gathered ribbon loop seated
  into a small knot and one layered face/hair cushion from front, side, and
  three-quarter views before restoring the rest of the bow or any body form.
- **Next approach:** locally test Blender cloth sewing springs or a
  deterministic equivalent on a continuous ribbon pattern, and rebuild the
  head/hair layer order around the measured cushion rather than treating hair
  color as the structural volume.
