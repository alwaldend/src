# Reimu Fumo attempt 4

[Back to attempt index](README.md) | [Back to goal](../README.md)

## Attempt 4 — identity-bearing quad-panel macro assembly

**Candidate:** planned `out/reimu_fumo_attempt_004_macro.blend`, SHA-256
pending, neutral macro sewn-assembly stage, review packet
`attempt-004-identity-macro-fixed-five-view`.

**Failure targeted:** Attempt 3 matched many scalar dimensions but remained
unrecognizable because a single inflated-cushion representation produced a
bare box head, mouse-ear bow, capsule sleeves, wedge dress, exposed neck, and
forward cylindrical pods.

**Hypothesis:** A complete head–hair–bow assembly built from distinct
quad-dominant pattern surfaces, combined with a low front/back gathered skirt,
open-cuff bell panels, and hem-tucked foot pods, will recover the reference
silhouette and sewn construction without relying on face graphics or fabric
materials.

**Plan written before implementation:**

1. Start again from factory-empty data and copy no Attempt 3 mesh.
2. Build the shallow head cushion from quad front/back grids and a perimeter
   gusset, avoiding n-gon subdivision and radial pinching.
3. Add the identity-bearing macro hair construction now: a thin padded crown
   cap around a broad face opening, one continuous three-lobe fringe sheet,
   cheek locks, and a rear/nape panel seated against the head cushion.
4. Build each bow loop and tail from thin folded front/back ribbon grids with
   asymmetric drape and a compressed crown-seam root; reject oval ears and
   vertical slabs in front and side views.
5. Eliminate the exposed neck. Build a compact bodice nested under the head, a
   low pooled skirt from separate front/back panels and a hem gusset, detached
   bell sleeves with readable cuff edges, and short pods tucked behind the hem.
6. Render the controlling front and side views first. Continue to rear,
   three-quarter, perspective, and black silhouettes only if the first two no
   longer read as a box, cone, capsule, mouse, or human body.
7. Align the front silhouette to the physical reference, measure the frozen
   rows, and run a new implementation-blind macro review. Reject before face,
   ruffles, final materials, or rigging if any identity-defining macro failure
   remains.

**Work performed:** Built a new factory-empty candidate with a bevelled quad
head cage, flush face region, continuous three-lobe fringe, cheek and nape
panels, layered angular bow loops and tails, nested bodice, low skirt and hem,
bell panels with cuff recesses, and two corrected `±0.315 Wh` foot pods. The
first build exposed both feet at the origin through its metric report; that
placement defect was corrected before the candidate was saved and reviewed.
Rendered only the planned front and side diagnostic and silhouette gates. No
rejected mesh or tracked Blender source was modified.

**Evidence:** Candidate SHA-256
`bf75e1b01d65a68066c41d0f7a070cae12f961dd814eb42183272208a9384262`.
The head measured `0.990 × 1.030 × 0.736 Wh`, the face region
`0.780 × 0.600 Wh`, sleeves approximately `0.323 × 0.427 Wh`, feet
`0.283 × 0.280 × 0.355 Wh`, hem `1.069 × 1.052 Wh`, and overall bow span
`1.689 Wh`. The aligned 42% physical-front overlay showed that many outer
bounds were close while the shapes remained visibly rigid and incorrect. The
front diagnostic showed a raised face mask, armor fringe, paddle bow, wedge
dress, hem bar, and block feet. The side diagnostic showed a tall cuboid head,
plate-thin hair and bow, rectangular torso, ramp skirt, and shoe-like foot.

The implementation-blind reviewer answered no to unlabeled recognition and
called the intended medium a rigid modular plastic or low-poly toy. Scores
were overall macro likeness 2.5/10, silhouette 3/10, construction 1.5/10,
contact and occlusion 2/10, intended-medium read 1.5/10, and presentation
6/10. Every identity-defining construction category contained a major visible
failure.

**Criterion results:** Reference likeness, full measured silhouette, plush
construction, and presentation quality fail. Scalar bounds alone pass several
rows but cannot override the wrong silhouette segments or absolute image gate.
Reusable structure, animation readiness, and technical integrity remain
unverified. Repository delivery still applies only to the rejected migrated
baseline.

**Decision:** Reject at the early front/side macro gate. Do not render the
remaining approval packet or advance this geometry.

### Progress and approach audit after attempt 4

- **Improved:** the front hairline landmarks, bow span, garment width, sleeve
  bounds, foot separation, and head/body overlap were more explicit and easier
  to inspect than in Attempt 3.
- **Regressed:** construction fell from 2/10 to 1.5/10, contact from 3/10 to
  2/10, and intended-medium read from 2.5/10 to 1.5/10. The clearer part map
  made the plastic armor and plate construction more obvious.
- **Absolute result:** the candidate remains 2–3/10 and unrecognizable without
  color cues. No acceptance criterion passed.
- **Approach evidence:** direct silhouette extrusion can match landmark bounds,
  but its planar caps, uniform bevels, independent slabs, and hard seams are
  structurally incapable of the soft fold continuity in the references.
- **Repeated-defect diagnosis:** two clean reviewed attempts using scripted
  outline-driven solids produced opposite but equally invalid results:
  inflated blobs in Attempt 3 and rigid plates in Attempt 4. The entire
  silhouette-extrusion representation is now discarded.
- **Highest-leverage problem:** establish one convincing continuous soft-volume
  construction language for the head, bow, body, sleeves, skirt, and feet,
  with broad stuffing, seated roots, and real contact before refining identity.
- **Next approach:** use sculpted rounded base volumes and curved quad fabric
  grids with smooth bulge, fold, and twist fields. Conform the face and hair to
  the head instead of stacking plates. Build the bow as one connected soft
  assembly and the garment as joined front/back cloth surfaces rather than
  frustums or extruded polygons.
