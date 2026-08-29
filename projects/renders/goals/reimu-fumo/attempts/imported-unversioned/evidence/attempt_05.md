# Reimu Fumo attempt 5

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 5 — sculpted continuous soft-volume assembly

**Candidate:** planned `out/reimu_fumo_attempt_005_soft_volume.blend`, SHA-256
pending, neutral sculpt macro stage, review packet
`attempt-005-soft-volume-fixed-five-view`.

**Failure targeted:** Outline-driven solids matched numbers but alternated
between balloon primitives and rigid armor. They never produced soft surface
flow, a continuous bow, integrated hair framing, gathered garment volume, or
convincing contact.

**Hypothesis:** Sculpted rounded base volumes plus all-quad curved fabric
surfaces, with low-frequency bulge and fold fields and explicit shared roots,
will read as one sewn stuffed object in front and side views before texture or
face graphics.

**Plan written before implementation:**

1. Start from factory-empty data and retain no prior candidate geometry.
2. Form the head from a highly rounded subdivided cube cage with a broad flat
   face and asymmetric stuffing; apply the form before adding conforming fabric
   regions so no mask or armor gap can appear.
3. Lay the face and crown regions directly against the evaluated cushion. Add
   the continuous fringe, cheek locks, and nape as curved quad surfaces whose
   roots and depth follow that cushion.
4. Build the bow as a connected knot, two curved all-quad loop pockets, and two
   twisted quad-grid tails. Use smooth bulge and fold valleys across each
   surface; do not derive lobes by extruding their front silhouette.
5. Build the skirt as joined front and back quad fabric grids with gathered top
   width, low pooled hem, subtle pleats, and side seams. Nest the bodice under
   the head; use curved bell grids for sleeves and rounded compact feet only
   slightly forward of the hem.
6. Use one neutral diagnostic clay family and strong grazing light so geometry,
   attachment, and contact must carry the plush read. Render front and side
   early; stop the representation immediately if it still reads as a box,
   armor, a cone, slabs, human anatomy, or disconnected primitives.
7. If the early gate passes, render the full fixed set, align the silhouette,
   measure all applicable rows, and request a new implementation-blind review.

**Work performed:** Built a new factory-empty character from a
sculpt-deformed rounded head, directly assigned face region, conforming quad
fringe and locks, connected bow knot with curved loop and tail surfaces, soft
bodice, joined front/back skirt envelope, curved bell sleeves, and compact foot
pods. Added fixed diagnostic cameras, neutral materials, grazing light,
automated bounds, front and side silhouettes, the aligned physical-front
overlay, and a reference/candidate contact sheet. No rejected geometry entered
the tracked asset.

**Evidence:** Candidate SHA-256
`57c4d7e4bb5a7b65299e8363d637a227944c474685ad16cf5434f117ea6ef528`.
The head measured `1.002 × 1.034 × 0.791 Wh`, fringe width `0.818 Wh`,
skirt `1.045 × 1.055 × 0.415 Wh`, each foot approximately
`0.300 × 0.280 × 0.370 Wh`, and bow top `2.189 Wh`. The sleeves were too
wide and the bow was about `0.079 Wh` too high. More importantly, the front
image showed a rabbit-ear or four-petal bow, helmet/cuboid head, symmetric
scallop hair, slab locks, paddle sleeves, and trapezoid skirt. The side image
showed the square head over a ramp-like garment.

The implementation-blind reviewer answered no to unlabeled recognition and
classified the medium as a matte molded CG toy or inflated primitives rather
than constructed fabric. Scores were macro likeness 4/10, silhouette 3/10,
construction 2/10, contact 3/10, intended-medium read 3/10, and presentation
6/10. A separate strategy audit scored likeness 2/10, silhouette 3/10,
construction 1/10, identity 2/10, contact 2/10, and plush-medium read 1/10.

**Criterion results:** Reference likeness, full measured silhouette, plush
construction, and presentation quality fail. Several outer bounds are close,
but the image gate and construction gate fail absolutely. Reusable structure,
animation readiness, and technical integrity remain unverified. Repository
delivery still applies only to the rejected migrated baseline.

**Decision:** Reject at the early front/side macro gate. Do not add face
graphics, ruffles, fabric materials, a rig, or this geometry to the tracked
asset.

### Progress and approach audit after attempt 5

- **Improved:** macro likeness rose from 2.5/10 to 4/10 in the blind review,
  intended-medium read rose from 1.5/10 to 3/10, and the softer surfaces
  removed some of Attempt 4's explicit armor-plate read.
- **Regressed or unchanged:** silhouette remained 3/10, construction remained
  an absolute failure at 2/10, and no criterion passed. The bow became rabbit
  ears, the head remained a helmet, and contact remained weak.
- **Absolute result:** close bounding boxes did not create reference likeness.
  The candidate was not recognizable as the physical variant without color
  cues, and it still read as a molded toy rather than a sewn plush.
- **Approach evidence:** the new code changed helper names and used curved
  grids, but it retained the same representational assumptions: generic
  bulge fields, independent surfaces, regular lofts, uniform subdivision, and
  tangent intersections. Those operators smooth a primitive; they do not
  encode fabric panels, seam tension, stuffing, gathering, or shared roots.
- **Repeated-defect diagnosis:** three factory-empty reviewed attempts have
  now failed for opposite-looking symptoms generated by the same cause. Whole
  character procedural generation from outlines and analytic volumes is
  discarded, including cosmetically curved variants.
- **Highest-leverage problem:** the coupled head, hood, and bow determine most
  of the unlabeled identity and expose the failed medium read in every view.
- **Next approach:** use one deterministic all-quad sewn-pattern cage. Prove
  the head's front/rear/gusset construction and the bow's paired folded panels,
  gathered roots, return leaves, and knot attachment before any body work.
