# Reimu Fumo attempt 18

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 18 — direct-reference front head construction blockout

**Candidate:** planned
`out/reimu_fumo_attempt_018_front_head/checkpoint_front_blockout.blend`,
untextured beige support plus brown crown/fringe and two side-lock panels,
review packet `attempt-018-raw-reference-front-head-blockout`.

**Failure targeted:** Attempt 17 followed much of the outer image boundary but
assigned all brown pixels to one closed helmet ring. That erased the physical
panel separations and collar opening that make the object read as a sewn plush.

**Hypothesis:** A compact curved beige support plus three separate reference-
mapped brown fabric pieces can match the actual front layering without asking
an outline, primitive, or nested slab to imply construction. Direct raw-source
comparison will expose whether the pieces overlap and terminate correctly.

**Plan written before implementation:**

1. Start factory-empty and byte-new. Import no rejected geometry or coordinate
   arrays. Create a non-rendering `References` collection containing the raw
   physical and clean front images as image empties placed beside the model,
   never over it. Pack the images, preserve their original aspect ratio, and
   hide the collection from evidence renders.
2. Freeze source-to-world registration from physical-front pixels:
   `x=(pixel_x-237.5)/189` and `z=(201-pixel_y)/189`, with the coupled head
   centered at world `(0,0,0)`, facing `-Y`, and `Wh=1`. Use this mapping for
   panel boundaries and external candidate/reference alignment. Do not fit a
   new curve after seeing candidate pixels.
3. Create one disposable native subdivision-four icosphere support at
   dimensions `0.86 × 0.68 × 0.80 Wh`, apply transforms, and voxel-remesh once
   at `0.025 Wh`. It is beige support clay, not an approved isolated head.
   Require one finite, positive-volume, closed component with no modifier or
   other support mesh. Its final front evidence is judged only where beige is
   visible through and below the hair pieces.
4. Use the already proved ordinary Essentials `Scrape/Fill` workflow on that
   support with exact live projection and fixed pressure `1.0`. Execute only
   two Y-symmetric front-view passes:
   `(-.14,.12)→(.14,.12)` and `(-.14,-.10)→(.14,-.10)`, each radius
   `0.24 P`, strength `.12`. Persist exact samples, raycast starts, executed
   size, topology and coordinate hashes, direction, locality, and displacement.
   Require `FINISHED`, unchanged topology, maximum displacement
   `0.004–0.055 Wh`, and final depth `0.64–0.70 Wh`. Any failure is terminal.
5. Construct the brown crown/fringe as one front fabric panel, not a solid
   cap. Its source-space outline follows the observable outer crown from
   `(149,155)` through `(171,120)`, `(205,106)`, `(236,104)`, `(278,108)`,
   `(311,136)`, and `(327,172)`, then returns along the lower boundary through
   `(302,181)`, `(285,199)`, `(259,164)`, `(235,211)`, `(203,160)`,
   `(180,198)`, and `(149,155)`. Preserve those source points exactly.
6. Construct separate left and right lock panels. The left outline is
   `(149,155)`, `(143,198)`, `(143,245)`, `(151,288)`, `(165,267)`,
   `(158,225)`, `(180,198)`. The right outline is `(327,172)`, `(331,211)`,
   `(326,253)`, `(320,292)`, `(304,270)`, `(310,225)`, `(302,181)`.
   Overlap each lock beneath the crown panel by `0.02–0.04 Wh`; neither lock
   nor crown may connect beneath the chin.
7. Build each brown piece as a dense curved surface panel projected onto the
   live support from front, then offset it forward by `0.012 Wh` and give it
   applied core thickness `0.035–0.045 Wh`. Boundary points come only from the
   registered source; interior tessellation serves curvature and may not alter
   the outline. Require separate named meshes, applied transforms, one closed
   component per piece after thickness, coherent normals, and no intersections
   deeper than the declared overlap band.
8. Keep clay colors only: matte beige support and three matte brown panels.
   Add no hair fibers, texture, seams, bow, face graphics, rear panel, body, or
   rig. Shade smooth, but do not use a render-only silhouette modifier or
   procedural displacement.
9. Save append-only setup, post-support-sculpt, and post-panel checkpoints with
   parent manifest hashes. Always reopen a verified parent checkpoint before a
   stage. On any exception, save a terminal checkpoint and JSON. Never mutate a
   prior manifest, metrics, review, or evidence file.
10. Render fixed front neutral, silhouette, and grazing views plus canonical
    three-quarter neutral at `640 px`. Outside Blender, align the front render
    to the raw physical source using the frozen affine transform and produce a
    `40%` overlay and edge difference. Compare candidate pixels directly to the
    source; do not compare against Attempt 17's rejected trace.
11. Extract visible candidate contours by camera-space silhouette. Require
    every controlling panel boundary point to project within `0.03 Wh` of its
    registered source point, outer coupled front width `0.97–1.03 Wh`, and
    height `0.99–1.06 Wh`. Require an open collar gap at least `0.24 Wh` wide,
    separate lock bottoms within `0.04 Wh` of their targets, and the broad
    central fringe within `0.03 Wh` of `(235,211)` after alignment.
12. Implementation-blind review receives only raw reference, aligned overlay,
    silhouette, grazing, and three-quarter pixels, described as an untextured
    front head-construction blockout. It records unlabeled recognition,
    intended medium, five ordered discrepancies, silhouette, panel ownership,
    fringe/lock likeness, softness potential, three-quarter layering, and a
    major-failure boolean.
13. Hard-reject closed helmet ring, under-chin bridge, generic bob, rounded-box
    support silhouette controlling the result, thin central spike, cat ears,
    flat graphic slab, identical locks, tangent-only overlap, deep collision,
    panel floating, visible triangulation, hard bevel, or source deviation over
    `3% Wh`. Require silhouette, panel ownership, and fringe/lock likeness at
    least `7/10`, softness and three-quarter layering at least `6/10`, no major
    failure, and no numeric veto.
14. A PASS accepts only the front support/crown/lock construction for a later
    rear-hair and side-depth stage. It does not approve a finished head, model,
    texture, rig, or tracked `.blend`. A failure ends Attempt 18 without tuning
    a point, stroke, thickness, overlap, camera, alignment, or threshold.

**Structural difference from Attempt 17:** the brown regions are three actual
curved fabric pieces with an open collar, not one colored ring; every precise
edge maps directly to the raw source; the support is judged only where visible;
and the first Blender gate is small enough to review before rear or body work.

**Result:** Terminal preflight rejection. No driver or Blender candidate was
created. The registered panel-overlay SHA-256 was
`76e7ec23cc3e35ed5508ff5dc292e6ad9f745be451343298e72338deff453997`;
the audit SHA-256 was
`eb6f8213fceb5eb9a5befcd7d97d4f72f297650f40b1110b5467270f998b7a60`.

All three polygons were simple, but the planned support could not receive
them. Against the `0.86 × 0.80 Wh` support projection, all seven outer crown
points, six of seven left-lock points, and five of seven right-lock points lay
outside the ray-hittable ellipse. Examples included crown top normalized
ellipse value `1.647` and both lock bottoms above `2.45`; Scrape/Fill would not
expand the silhouette.

The source overlay also showed that the left outline was a `0.196 × 0.704 Wh`
crown-to-chin wedge and the right a `0.153 × 0.635 Wh` wedge rather than short
isolated locks. Each shared an exact seam edge with the crown, encoding tangent
contact rather than the required hidden underlap. Solidify direction and the
support-contact surface were unspecified, so thickness could float or
penetrate. The collar-gap gate also lacked a fixed measurement height.

**Criterion results:** Polygon simplicity passes. Every artistic, projection,
contact, model, reuse, and animation criterion fails or remains unverified.

**Decision:** Reject before Blender. Do not make panels until a support volume
actually occupies the head beneath the cap and a later lock stage supplies
explicit seam allowance and free-hanging construction.

### Progress and approach audit after attempt 18

- **Improved:** the direct source overlay exposed piece ownership and support
  occupancy before a dense Blender panel was built.
- **Regressed or unchanged:** no rendered geometry improved, and the proposed
  locks remained more schematic than the physical construction.
- **Absolute result:** the polygons are mathematically valid but physically
  unbuildable on the planned support. That is a terminal failure.
- **Invalid ordering:** thin cap and lock panels cannot be projected before the
  stuffed head surface beneath them exists at approximately the coupled outer
  envelope. Shrinking that support to an inferred visible beige opening made
  most observed hair boundaries unsupported.
- **Evidence against continuing:** a projection modifier cannot create missing
  support volume, and a seam thickness cannot turn a shared tangent edge into
  an underlap.
- **Highest-leverage unresolved problem:** create one soft, full head support
  whose front projection occupies the measured physical head envelope without
  becoming an egg, rounded box, or two-lobe clay artifact.
- **Approach decision:** return to native Sculpt Mode only for that support,
  now using the exact raw outer boundary as an occupancy target, safe interior
  stroke starts, true live-`P` radii, immediate aligned pixels, and no hair
  panels in the same family. Attempt 19 must pass before panel construction.
