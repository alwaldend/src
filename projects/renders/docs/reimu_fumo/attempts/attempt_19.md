# Reimu Fumo attempt 19

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 19 — native full head support for observed fabric panels

**Candidate:** planned
`out/reimu_fumo_attempt_019_full_support/checkpoint_support.blend`, neutral
stuffed support stage, review packet
`attempt-019-native-support-raw-front-and-fixed-views`.

**Failure targeted:** Attempt 18 tried to place observed cap and lock edges on
an undersized support that did not exist beneath them. Attempt 15's earlier
full-height sculpt family was impossible and mislabeled brush radii, so it did
not test whether a correctly calibrated native support can work.

**Hypothesis:** A proportioned but under-scale disposable clay blank followed
by five positive, safely started native Grab strokes can occupy the real head
envelope. Two low-strength native Scrape/Fill passes can then broaden the
front/rear panels without another box generator. Immediate raw-source-aligned
review can accept or reject the support before any hair topology depends on it.

**Plan written before implementation:**

1. Freeze the support ownership band at width `0.92–0.97 Wh`, height
   `0.98–1.04 Wh`, and depth `0.66–0.74 Wh`. The later hair cloth may add up
   to `0.03 Wh` per visible front side and `0.04 Wh` at the crown. Do not infer
   the support from the beige opening or change these bands after rendering.
2. Use the raw physical-front source and its existing registered affine
   transform as the front occupancy controller. Candidate contour distance is
   the symmetric Hausdorff distance between 512 equally arc-length-resampled
   visible contours after frozen center/top/width alignment. Critical crown,
   widest-side, lock-root-support, and underside distances must be at most
   `0.03 Wh`; all other contour distances may be at most `0.05 Wh` because
   brown fabric later owns the final boundary.
3. Start factory-empty and byte-new. Add one native subdivision-four
   icosphere, scale it once to `0.84 × 0.68 × 0.86 Wh`, apply transforms, and
   voxel-remesh at `0.025 Wh`. The untouched blank is disposable, receives
   technical mesh checks only, and is never rendered or scored.
4. Configure ordinary Essentials `Grab` with spherical smooth falloff,
   pressure `1.0`, exact `size=round(radius_fraction×P)`, locked VIEW sizing,
   and no Dyntopo, front-face restriction, automasking, autosmooth, or topology
   rake. Require executed radius error at most `0.5/P` and persist every
   sample, setting, projected start/end, start raycast, coordinate/topology
   hash, moved set, direction, locality, and displacement. Only stroke starts
   require surface hits; outward endpoints explicitly do not.
5. Run disposable fresh-blank technical probes for Grab radii `0.28`, `0.30`,
   and `0.32 P` and Scrape/Fill radius `0.28 P`. Recompute exact asset,
   operator, radius, finite movement, topology, paired X/Y symmetry residual,
   direction, and locality predicates when loading evidence. Purge all probe
   data before setup.
6. Execute exactly five front-view Grab strokes in order. Upper shoulder:
   `(.22,.24)→(.35,.38)`, radius `.30 P`, strength `.72`, X/Y symmetry.
   Lower shoulder: `(.23,-.22)→(.36,-.38)`, `.30 P`, `.72`, X/Y.
   Lower-middle cheek: `(.32,-.02)→(.47,-.04)`, `.28 P`, `.78`, X/Y.
   Crown: `(0,.34)→(0,.51)`, `.32 P`, `.70`, Y. Underside:
   `(0,-.34)→(0,-.50)`, `.32 P`, `.70`, Y. Every start has more than
   `0.08 Wh` axial or ellipse-normal margin on the untouched blank.
7. Require each Grab to return `FINISHED`, preserve topology, move a local
   neighborhood with correct signed drag projection, have maximum displacement
   `0.060–0.190 Wh`, and strong influence radius at most `0.40 Wh`. Save an
   immutable checkpoint and child manifest after every stroke. Any exception
   saves a separate terminal checkpoint and JSON before re-raising.
8. After Grab five, require width `0.89–1.00 Wh`, height `0.95–1.07 Wh`, depth
   `0.66–0.72 Wh`, crown/underside notch at most `0.012 Wh`, finite one-piece
   manifold topology, and no self-overlap. Apply exactly one `0.025 Wh` voxel
   remesh and require every dimension to move at most `0.015 Wh`.
9. Execute exactly two front-view `Scrape/Fill` passes with Y symmetry only:
   `(-.14,.12)→(.14,.12)` and `(-.14,-.10)→(.14,-.10)`, radius `.28 P`,
   strength `.12`, pressure `1.0`. Require `FINISHED`, maximum displacement
   `0.004–0.055 Wh`, unchanged post-remesh topology for each stroke, and final
   front/rear face bow `0.020–0.085 Wh` measured by frontmost/rearmost ray hits
   at center versus `x=±0.28±0.01 Wh`, using the median across a `0.02 Wh`
   vertical band.
10. Final numeric gates are width `0.92–0.97 Wh`, height `0.98–1.04 Wh`,
    depth `0.66–0.74 Wh`, symmetric Hausdorff contour error within step 2,
    face bow within step 9, axial notch at most `0.010 Wh`, positive finite
    signed volume, one closed component, coherent orientation, no self-overlap,
    zero modifier, applied transforms, and no scale normalization.
11. Render immediate `640 px` neutral and silhouette front plus neutral side,
    rear, canonical three-quarter, mirrored three-quarter, and readable grazing
    front. Outside Blender, assemble raw-source front, candidate, `40%` aligned
    overlay, edge difference, and fixed-view sheet. Run no geometry operation
    after the exact checkpoint used for evidence.
12. Implementation-blind review is told only: "neutral stuffed head support
    for a character plush; cloth hair and graphics come later." It records
    intended medium, five ordered discrepancies, front/side/rear/three-quarter,
    continuous stuffing, non-primitive form, support suitability, presentation,
    and a major-failure boolean. It must not score hair identity or seams.
13. Hard-reject sphere, egg, balloon, capsule, rounded box, mattress, two-lobe
    split, axial trench, pole pinch, brush dimple, stretched spike, flat plate,
    hard shoulder, roof, parallel wall, broken three-quarter highlight, or an
    image packet that hides the contour. Require every view, stuffing,
    non-primitive form, and support suitability at least `7/10`, presentation
    at least `6/10`, no major failure, and no numeric veto.
14. A PASS accepts only the frozen support surface for a new constructed crown
    and lock plan with explicit seam allowances. Sculpt topology remains
    disposable until fabric retopology. A failure ends Attempt 19 without
    changing a target, stroke, brush, gate, camera, or reference.

**Structural difference from Attempt 15:** all start points are safely inside;
brush pixels equal their declared live fractions; endpoints need not raycast;
the raw source controls the first pixels immediately; evidence is append-only;
and no later hair panel is forced onto an absent support.

**Work actually performed:** The complete frozen plan received a read-only
geometry, reference-ownership, operation, measurement, and checkpoint audit
before driver construction or Blender execution. No driver, candidate, render,
or Blender mutation was created. The terminal audit is
`out/reimu_fumo_attempt_019_full_support/preflight_audit.json`, SHA-256
`6c5bc11057d6008dc412371ba94e6438a9f96446fbf6b47816ef44de47f6110d`.

**Raw verification evidence:** The controlling trace contains approximately
`0.25 Wh` vertical hanging-lock runs and a `0.58 Wh` nearly flat lower rail. A
smooth `0.95 × 1.02 Wh` support ellipse has approximately `0.160 Wh` symmetric
Hausdorff error; even `0.97 × 1.04 Wh` remains approximately `0.143 Wh`, far
above the frozen `0.05 Wh` gate. Therefore a numerical pass would require the
support to copy hair-owned lock corners and the lower rail. The ideal seed does
contain all five proposed starts with nearest-boundary margins `0.0997`,
`0.1064`, `0.0994`, `0.0900`, and `0.0900 Wh`, and the predicted Grab maxima
`0.1376`, `0.1484`, `0.1171`, `0.1190`, and `0.1120 Wh` are compatible with
the displacement gate. Those local feasibility results cannot repair the
ownership contradiction. Scrape/Fill also lacks a feasible effect gate: the
prior verified probe moved only `0.000468 Wh` at strength `.34`, while this
plan demanded at least `.004 Wh` at strength `.12`.

**Acceptance result:** Criteria 1 through 7 remain `unverified` and therefore
fail the goal gate. Criterion 8 remains passed only for the unchanged migrated
baseline. No fixed-regression visual criterion can pass without candidate
pixels. The preflight itself passes its purpose because it rejected a
deterministically impossible plan before mutation.

**Decision:** Terminal preflight rejection. Preserve the frozen plan and audit
as evidence; do not construct its driver or execute its geometry. Reset the
controller from a hidden-support silhouette to the visible unified outer
cushion defined by the modeling contract.

### Progress and approach audit after attempt 19

- **Improved:** the plan's Grab start math and displacement ranges were shown
  feasible, and the new preflight rule prevented another long driver and
  Blender round trip.
- **Regressed or unchanged:** no model pixel or acceptance criterion improved.
  Five successive preflight families have now produced no candidate geometry,
  so the proof-heavy sculpt-driver workflow has poor feedback latency.
- **Absolute result:** the plan is structurally invalid. Its main numeric gate
  can pass only by producing the same lock-corner and flat-rail support that
  its own representation gate rejects.
- **Repeated assumption failure:** the physical plush does not provide evidence
  for a full hidden cushion that must match all visible hair. The frozen
  modeling contract instead assigns the beige face and brown crown/rear to one
  visible stuffed cushion, with the fringe and hanging locks as later panels.
- **Evidence against continuing unchanged:** symmetric outer-hair Hausdorff
  distance rewards wrong ownership; Y-symmetric front strokes prevent
  independent front/rear shaping; Scrape/Fill has no demonstrated effect band;
  and disposable voxel topology delays the reusable animated result.
- **Highest-leverage unresolved problem:** produce a controllable, softly
  irregular visible outer head cushion with an independently flatter face,
  fuller rear, variable side transition, and room for later rooted locks.
- **Approach decision:** discard Attempt 19. Preflight a topology-preserving
  sculpt blank or coupled all-quad sewn cage whose visible cushion—not a hidden
  support—is the reviewed form. Preserve the user's sculpt-first requirement
  by using native Sculpt Mode for macro shaping before any hair, face graphic,
  material, body, or detail work, while avoiding a second remesh and unproven
  brush family.
