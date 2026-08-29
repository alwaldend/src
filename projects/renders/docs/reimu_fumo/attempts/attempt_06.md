# Reimu Fumo attempt 6

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 6 — coupled sewn head-and-bow pattern proof

**Candidate:** planned `out/reimu_fumo_attempt_006_head_bow.blend`, SHA-256
pending, neutral construction stage, review packet
`attempt-006-head-bow-fixed-component-views`.

**Failure targeted:** Smooth primitives and independent surfaces can match
width, height, and depth while still producing a cuboid helmet and petal bow.
They lack explicit panel topology, shared seam boundaries, return folds,
gathering, and physically seated roots.

**Hypothesis:** A low-resolution all-quad front/rear/gusset cushion plus closed
paired-panel bow pockets authored around matching seam loops will produce the
reference-specific broad head, horizontal folded bow, and sewn-medium read
before the body or surface decoration exists.

**Plan written before implementation:**

1. Start from factory-empty data and retain no prior candidate geometry. Build
   `9 × 9` all-quad front and rear head patches with matching 32-vertex
   perimeter loops, joined by three all-quad gusset rings. Use no ngons, cube
   bevel, voxel remesh, radial pole, or full hair shell.
2. Keep head width `1.00 Wh`, height `1.03 Wh`, and depth inside
   `0.66–0.82 Wh`. Limit front convexity to `0.02–0.03 Wh`, put most depth at
   the rear, taper the lower corners about 5%, and use paired support loops at
   the front/gusset seam. Add subdivision and a `4 × 3 × 4` lattice only for
   stuffing asymmetry and lower-head settling.
3. Derive one thin hood from the evaluated upper/rear head surface, cut a
   `0.78 × 0.60 Wh` face opening, conform it at `0.018–0.025 Wh`, and close it
   at `0.025–0.035 Wh` thickness. Extend its upper opening rows into one
   continuous fringe; root closed cheek and nape panels beneath the hood seam.
4. Build each bow loop from closed paired quad panels with about seven stations
   from gathered root to outer fold and seven vertices across. Return both ends
   beneath the knot, separate faces by `0.03–0.05 Wh`, pinch the root, press a
   center row inward for the fold valley, and author one twist in each closed
   tail. Bind the roots to the head lattice and wrap them with a flattened
   band-shaped knot.
5. Render only this coupled component: calibrated front clay, fixed
   three-quarter clay, fixed side silhouette, aligned 42% physical-front
   overlay, and a grazing-light close-up of the cushion seam and bow roots.
6. Reject before body work unless the fixed measurements remain within the
   frozen bands, every root is visibly seated, and a pixels-only reviewer calls
   the forms a sewn cushion with a folded horizontal cloth bow rather than a
   cube, helmet, petals, wings, or rabbit ears.

**Work performed:** Built a factory-empty, closed all-quad head from structured
front and rear grids joined through three gusset rings. Replaced the analytic
front outline with the measured physical-front perimeter, normalized the
evaluated head to `1.00 × 0.74 × 1.03 Wh`, and used the traced irregular face
opening and hairline. Built paired closed all-quad bow panels directly from the
measured asymmetric loop and tail traces, added compressed roots, a knot,
perimeter seams, fixed cameras, reference empties, and neutral construction
lighting. Iterated internal fold depth only until the front and side evidence
were stable enough for the formal gate. No rejected geometry entered the
tracked asset.

**Evidence:** Candidate SHA-256
`3523bd9f7c01bd97274d2a9355be6f798c83f94cd2f59c877e86cd2dee6a8cfa`.
The evaluated head is exactly `1.000 × 0.740 × 1.030 Wh`; the whole bow span is
exactly `1.693 Wh`. Its upper and lower extrema are `2.086` and `1.122 Wh`,
and the four panel cages are closed `9 × 7` paired grids with 126 vertices,
124 all-quad faces, zero boundary edges, and zero non-manifold edges each. The
head and every hair, bow, and knot mesh also have zero boundary edges, zero
non-manifold edges, and no non-quad faces.

The aligned overlay confirms that the evaluated front head perimeter, face
opening, asymmetric hairline, gather, upper loop extrema, tail extrema, and
overall bow span closely follow the physical-front control image. The fixed
three-quarter and side images expose the unresolved representation: a deep
box-like head, flat forehead visor, long symmetric side slabs, one horizontal
bulb loop, one vertical fin loop, pill-like tails, and a knot perched through
tangent/intersecting roots.

The implementation-blind reviewer narrowly recognized the same Reimu variant
from the component alone, the first reviewed clean attempt to do so. Scores
were head likeness 5.5/10, bow likeness 3.5/10, silhouette 4/10, sewn
construction 4/10, contact and occlusion 4/10, and presentation 7/10. The
reviewer reported automatic-failure helmet/box, rabbit-ear/wing, and flat-panel
reads.

**Criterion results:** Several controlled front landmarks and the manifold
all-quad construction pass their local checks. Reference likeness, complete
measured silhouette, plush construction, and presentation quality still fail
the absolute image gate. Reusable structure, animation readiness, and full
technical integrity remain unverified. Repository delivery still applies only
to the rejected migrated baseline.

**Decision:** Reject at the coupled component gate. Preserve its measured
front traces and topology evidence, but do not extend this geometry to the body
or tracked asset.

### Progress and approach audit after attempt 6

- **Improved:** unlabeled recognition changed from no to narrowly yes. Head
  likeness reached 5.5/10, silhouette reached 4/10, sewn construction reached
  4/10, and presentation reached 7/10. The physical-front overlay is the first
  close candidate/reference alignment produced in the clean rebuild.
- **Regressed or unchanged:** no score reached the 8/10 gate. Bow likeness was
  only 3.5/10, roots remained implausible, and side/three-quarter views still
  showed a box, visor, fins, wings, and flat panels.
- **Absolute result:** the component is recognizably Reimu but is not yet a
  convincing physical Fumo. No acceptance criterion passed as a whole.
- **Approach evidence:** reference-derived perimeter traces and exact evaluated
  normalization materially improve likeness and must be retained. All-quad
  topology and manifold checks prevent technical shortcuts but do not by
  themselves create fabric construction.
- **Representation failure:** the head used the correct front boundary but
  constant front and rear seam depths, so the side view remained rectangular.
  The bow used correct front boundaries but treated each entire traced region
  as one inflated closed pocket, so folds became fins and pills. A separate
  broad fringe panel became a visor.
- **Highest-leverage problem:** encode the controlling side contour and real
  layer order directly. The head needs height-varying front and rear depth;
  the bow needs thin visible and return cloth layers rather than padded lobes.
- **Next approach:** preserve the traced front boundaries, map the measured
  turn-side contour into every head cage row, remove the broad fringe object,
  and construct each bow lobe as a thin folded panel plus a separate return
  layer terminating in a compressed shared gather.
