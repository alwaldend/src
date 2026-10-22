# Static gathered-sheet study 134

## Decision before execution

**Proceed with one bounded static fixture.** The failed native cloth profile
in study 130 is closed for this cycle. Study 129 could not produce gathered
fullness because it extruded one curve through the ribbon width. Here a planar
annular rest sheet supplies unequal attachment/free-edge lengths, while only
13 discrete stitches contract to a shorter seam. The free edge remains free.

Use a 72/96 mm inner/outer edge, 9 mm width and 48 mm finished seam. Static
edge and shear projections preserve the rest metric; a dihedral regularizer
resists sharp bending without prescribing a final waveform. A small initial
perturbation only breaks symmetry. This is geometric relaxation, not a
calibrated textile model or a continuation of Blender's failed dynamic profile.

The [authors' XPBD paper](https://matthias-research.github.io/pages/publications/XPBD.pdf)
explains iteration-dependent stiffness in ordinary position projection and
why compliance needs an explicit formulation. It supports distinguishing our
bounded geometric fixture from a physical fabric calibration. The
[position-based dynamics survey](https://matthias-research.github.io/pages/publications/PBDTutorial2017-CourseNotes.pdf)
distinguishes bending constraints from distance surrogates that also affect
stretching. These sources inform the representation; they do not establish
that the resulting folds match the Fumo.

The strongest failure risk is a numerically valid corrugated strip that still
looks unlike soft gathered fabric. The alternatives are a directly modeled
sewn panel or another dynamic material setup. The static fixture offers a
small, inspectable test without retuning the failed material profile. Allow
one diagnosed implementation repair; do not search numerical settings after a
valid but visually rejected result.

Before transfer require full seam contraction, finite positions, edge/shear
absolute strain P95 at most 2% and maximum 8%, anchor error at most 10 microns,
no collapsed or antiparallel hinge degeneracy, final sweep residual at most
20 microns, zero audited nonadjacent crossings and no sampled self-clearance
below rendered thickness. Native checks cover the exact evaluated mesh and
record their adjacency exclusions and sampling limits. A numerical pass is
only eligibility for image review, not permission to integrate.

Root remains the sole native writer. Pure-math authoring may run outside
Blender; its saved output and module identities must be bound in the native
receipt. Preserve source 127, SHA-256
`1834951f4efa931238a06e560153f1f4d185f47fe0aa2e89a30612e3150c5079`,
and save new ignored fixtures with pinned Blender 5.2.1, build `9e2066aef7ef`.
Inspect front, side, elevated and top views from the clean-reopened file.
No full-model or stage acceptance changes.

## Garment transfer boundary

A reference review confirms that the current waist-to-hem cone and exposed
red seat are the dominant dress failures. The red/white front widths are
already near their recorded targets. A useful replacement needs cloth slack,
localized thigh contact and a low rounded rear hem, not global rescaling.
Do not lower the entire front hem to hide the seat: the canonical front apron
is shallow and its photograph has unresolved elevation. A fixture pass alone
cannot establish full garment/body correspondence.

## Result

**Reject and do not transfer.** The frozen solver
`6c30c259636ebbeabfd97fe4693f8bd2a555cf8660eca113a853efa8a97d987a`
produced solution
`ac49fab2cfdf18515f1a2031d9b2ece664857a50a646438d1a08875f3dda7cbc`
in 16.17 seconds and 640 cycles. All 13 anchors reached the 48 mm seam,
but absolute distance strain reached P95 5.576%, maximum 27.673%, and final
cycle motion was 0.11675 mm. No triangle collapsed. Of 120 distance
constraints over 8% strain, only three touch anchors; this is not solely a
pinned-boundary artifact. Independent finite-difference audit verified the
dihedral gradients to 1.98e-9 and found no material projection or hinge bug.
No coefficient or iteration sweep followed.

Pinned Blender saved the rejected diagnostic
`e45a34a6b262e60002eace41532952e89b2a7a4af543b4db3b70be9ee2d3ad6c`.
After one subdivision level the native audit found 921 nonadjacent midsurface
crossing pairs and 4,235 solidified-shell pairs. There were 10,363 sampled
vertex/triangle pairs closer than the 0.35 mm rendered thickness, with the
closest about 0.000595 mm. Shared-vertex triangle neighbors and immediate
one-ring clearance neighbors are excluded. Counts describe bounded audit
coverage, not all possible contacts. Protected source 127 is unchanged.

All four fixture views were clean-reopened. Root rejects the pixels: the
cloth forms a broad curved flange with small tight folds near the stitches
and folded ends, rather than rounded open ruffle mouths. Numerical and pixel
failures independently prevent transfer. The dormant wall-limit cycle-count
off-by-one and hardcoded perturbation-description text do not affect this
default cycle-limited result; the source remains frozen with those limitations.

## Reopened construction decision

**Revise to directly authored transverse fabric turnbacks.** Reference and
source review establish that another smoothly interpolated skirt envelope
would repeat studies 123 and 128. A distinct bounded construction test needs
extra material visible across a fold: separated return layers gathered into
a short sewn waist region, releasing toward a broad free panel. The side and
rear folds must change the silhouette and contact, not decorate a cone.
Do not transfer failed 134 geometry or declare a physical-material pass.
The full garment/body projection conflict remains open, and cameras and
acceptance criteria remain unchanged. A detailed bounded panel plan precedes
native study 135; 127 remains the continuation source until reviewed evidence
supports a specific replacement.

### Study 135 scope

**Proceed with one coupled skirt construction candidate.** Independent
reference/source critique favors a single side-gore fixture as the cheapest
causal proof. Root chooses a complete four-sector skirt for this bounded test
because the user's outcome is the whole dressed silhouette and the waist,
side/rear fullness and hem must join coherently. This costs a fuller contact
and multi-view audit and cannot attribute any improvement solely to a pleat.

Insert four explicit side/rear material turnbacks with 6–8 mm of extra path
per fold. Retain them through a short sewn region and release over roughly
20–25 mm into broad inward troughs of 4–7 mm. These are authored construction
parameters, not measured pattern pieces. Keep the shallow front center,
canonical red/white widths near 121/141 mm and the recorded uncertain
131–152 mm total hem-depth band. A matching white band follows the actual
folded boundary; failed 134 and the extruded 129 ribbon are excluded.

Preserve bodice, original feet, head, hair, bow, materials and fixed cameras.
The existing closed red pelvis may be reshaped only as necessary to contain
it under the new sewn shell while joining the unchanged foot roots. Keep its
posterior support beneath the torso; no cavity, material substitution or
camera change may hide the unresolved front-gap contradiction. Report this
coupled body change explicitly rather than attributing it to a skirt-only test.

Require preserved-object geometry identity, closed core, no evaluated
unintended crossings and clean attachment to the unchanged waist and feet.
Allow at most one diagnosed contact/setup repair after fixed front/side
inspection. Reject if smoothing erases the actual return path, broad folds
read as hard ridges, side fall remains triangular, or another view regresses.
This construction test grants no material, whole-model or stage acceptance.

### First 135 result and diagnosed completion

The first native snapshot is
`7e289243cfc9f0465cb97a5ee582c563127923c93249e804c1ac3020ba9a2bcc`,
from module `8d474c4b9c463c3506153a68d28fc97f3eca6ec1091d837dde5b61ba81c593cc`.
Root's independent evaluated-geometry audit preserves all 62 protected meshes.
Cloth/core/foot shell pairs have zero audited crossings and no sampled signed
penetrations beyond 0.15 mm. The closed core has positive volume, no boundary
or nonmanifold edges, minimum Z 0.0183 mm and 72.97 mm² of posterior surface
below 0.45 mm. Core/foot crossings are confined to the declared proximal sewing
regions. These findings do not grant a complete contact or visual pass.

Four skirt and 42 white-ribbon self-crossing pairs reject the setup. It is
saved explicitly as a rejected diagnostic, with the new geometry revealed and
all intended old replacements hidden. Fresh five-view, clay and presentation
renders bind that exact file. Implementation-blind reviewer
`/root/model135_absolute` rejects it: likeness 6.5, silhouette 7,
construction 5.5, identity 9, contact 8, medium 4 and presentation 8.
The principal failures are smooth manufactured surfaces, a flared skirt with
thin wavy edging, rigid sleeves, sheet-like rear hair and stiff bow panels.
These image scores do not override native self-contact failures.

The root folds survive subdivision: the two side and two rear sections have
14/16/16/14 measured transverse backtracking intervals and 6.72–7.98 mm extra
upper material arc. However, source inspection proves that both upper-return
and trough terms become exactly zero after 25 mm of meridian distance. The
whole lower skirt and hem therefore revert to the smooth envelope. This is
an omitted part of the intended hanging-fold construction, not evidence that
visible lower folds were tested and failed. The rear pixels corroborate the
short waist cuts and smooth lower surface.

**Revise once to complete the intended material path, then stop for review.**
Open each upper return into an explicit broad, rounded U-shaped transverse
panel that continues to the hem. Preserve global extremal landmarks, the front
center, closed core and every protected component. Do not merely extend a
scalar depression on the old envelope. The exact rear skirt crossings occur
near material coordinates `(u=0.406/0.593, v=0.264–0.278)` around
`(X=±23.11, Y=26.78, Z=41.77)` mm, consistent with the collapsing return.
The hem witnesses lie at its free rim around `(8, −81.6, 31.05)` mm, where
the tangential cup displacement nearly cancels forward progression. Apply
one geometrically justified rim-curvature correction without thinning the
fabric. This disclosed implementation completion goes beyond a contact-only
repair; it does not authorize a coefficient sweep or change acceptance.
The first module, snapshot and packet remain frozen separately.

### Completed 135 result

**Reject promotion; retain only as comparative evidence.** The completed
module `2453909742f1fe981706bb00c9098c010a848f2d8f1844f8e3c52d83ed1f2f69`
produced native snapshot
`cf6cae4aa30e989ed7436c1ecfff9ef80e4b6e94bfec528766b0514b9e2d4369`.
Its five fixed views, two clay views and uncropped presentation were rendered
after a clean reopen. Source 127 remains unchanged; all 62 protected evaluated
meshes and the closed seated core remain identical to the first 135 candidate.

The lower material section now exists: at one side fold the evaluated release
depth is 5.486 mm and the upward hem cup is 5.505 mm. The 12 section witnesses
cover all four folds at the upper return, release and hem. They demonstrate
the missing deformation, not a calibrated rest metric: that release and hem
have respectively 0.169 and 0.260 mm less arc than their undeformed envelope
comparisons. Upper backtracking alone never established lower surplus cloth.

The unchanged-thickness hem correction removes all 42 audited self-crossings.
However, the skirt now has **461 nonadjacent crossing pairs**, versus four
before completion. Its rounded transition still violates shell clearance.
There are no audited skirt/hem crossings with the core or either foot, no
sampled signed penetration beyond 0.15 mm, and no core self-crossings. Those
bounded successes do not override the skirt's failed native gate. No further
135 repair or parameter sweep followed.

Implementation-blind reviewer `/root/model135_completed_absolute` rejects
the complete image packet: likeness 6.5, silhouette 7, construction 5.5,
identity 8.5, contact 7, medium 5.5 and presentation 8. Ordered failures are
the bulbous face and helmet-like cap, thin rigid bow construction, shell-like
skirt and exposed underside, stacked leaf-like rear hair, and the collar/tie
shape. Recognition is strong; a material-texture pass cannot correct these
construction failures. Different reviewers' scores do not establish a trend.

After its earlier absolute review, `/root/model135_absolute` compared 127
and completed 135. It favors the rounder hem and less conical side drape,
while identifying the straight front shelf, abrupt waist and newly exposed
pale rear forms as remaining problems. Its image-only recommendation to
retain 135 is overruled by root because native crossings prevent promotion.
The relative improvement remains useful evidence; 127 remains the continuation
source, and all final acceptance criteria stay open.

### Next representation decision

**Revise the next cycle to a whole-figure neutral construction blockout.**
Do not begin another lower-body substitution. The independent reference
review found that study 123 already tested a raised front underside, closed
seated core, splayed foot pods and shared garment support. Study 135 repeats
that broad family while the camera/body projection ambiguity remains open.
Neither supports inventing a new hidden body or changing the fixed camera.

At least three identity-defining forms still fail together. The next cycle
must address their coherent mass and fabric construction before adding folds
or surface detail. The proposed shared support is a shallow gusseted head:
a broad gently curved face plane with cheek/temple roll at its perimeter,
closely fitted cap/fringe and rooted hanging hair panels, and a bow compressed
onto the crown with broad center pinches and bent wings/tails. Rebuild these
forms together so their contacts follow the same support volume. This differs
from 127, which protected the anterior cap while replacing posterior hair.
Keep the lower body as explicitly rejected context, with no foot-gap objective.
Record multi-view landmarks and their uncertainty before authoring. The camera
investigation remains diagnostic; the fixed acceptance contract is unchanged.
This is a continuation direction, not an executed or accepted blockout.

## Bounded process review

Independent ergonomics review inspected six distinct artifacts, selected receipt
fields and 108 source lines. It identified two task-local controls:

- `REIMU-134-SOLVER-CAPABILITY`: preserve the explicit absence of self-contact
  response and the failed metric/contact results. A static solver's successful
  execution is not a calibrated cloth capability; no transfer or parameter
  search follows this rejected fixture.
- `REIMU-135-LOWER-FOLD-OMISSION`: upper-edge arc/backtracking witnesses are
  correctly scoped but cannot demonstrate a hanging fold. Extend the existing
  witness to the release endpoint and hem, measuring the intended broad section
  there without requiring lower backtracking. This should catch the observed
  dropout before full rendering; pixel review remains necessary.

Root also observed `REIMU-DELIVERY-MISSING-RUNFILES`: the study-131 generated
launcher failed before executing because its referenced runfiles directory
was missing. A fresh study-135 launcher build completed in 172.745 seconds.
This is not evidence of stale semantic tool content, and rebuild duration
alone is not a defect. Generated runfiles paths are temporary, so continuation
must verify their availability and use the owning refresh workflow when needed.
No host configuration, shared skill, solver library or acceptance contract was
changed. One oversized skill-read batch was truncated; only the omitted skill
content was subsequently loaded in bounded reads.
