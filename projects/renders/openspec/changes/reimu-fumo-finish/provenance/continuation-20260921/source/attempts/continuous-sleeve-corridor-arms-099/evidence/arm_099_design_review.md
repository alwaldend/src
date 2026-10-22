# Arm 099 bounded design review

Observed 2026-09-06 at approximately 03:55 UTC. Recommended verdict:
**proceed with the coordinator's updated section-based fit contract**.
It directly addresses both the 095 deformation failure and the exact-interface
preservation conflict proven by 098. This is approval of a bounded hypothesis
for testing, not acceptance of geometry. The coordinator owns the final
verdict, dimensions, implementation, and acceptance.

This report reviews the original proposal and the coordinator's subsequent
clarifications together. In particular, the final proposal uses a consistent
right-handed basis, conservative distance-bound marching instead of
bisection, actual full-sleeve queries on both sides of the root, and a
connected body-attachment classification. The superseded versions are not
reported below as unresolved defects.

## Evidence and causal fit

`arm_098_source_attachment.md` was read in full. It proves that 19/20
complete body-crossing arm triangles also have sleeve-crossing intervals
outside the body, at actual sleeve rows 6–14. Preserving all original
interface triangles and the cloth would necessarily preserve these
intersections. The coordinator now allows the body-contact locus to change
while retaining the body, cloth, and a continuous buried arm connection.
That removes the demonstrated incompatibility; no hidden-root inference or
garment relocation is required by this evidence.

The prior 095 diagnosis located the dominant new crossings in the translated
upper/rear transition, with the distal cuff contacts removed. Reparameterizing
the source unit sphere by longitudinal `s=(u+1)/2` and moving complete
sections coherently replaces the per-vertex world-X displacement field that
caused that transition failure. Fitting the transverse section envelope then
addresses the space available inside the retained sleeve. Keeping the useful
095 distal pole preserves a meaningful position target, although rebuilding
its surrounding sections does not preserve the entire failed 095 bulb shape
or grant it a visual pass.

The strongest case against an arm-only fit was the fixed-interface conflict.
The new authority resolves that conflict, and the full-contact gates can
reject an unsuitable fitted arm. Changing the cloth now would add externally
visible scope without evidence that a credible continuous arm cannot fit.

## Scrutiny of the updated construction

**Section coordinates and orientation.** The supplied normalized sleeve
directions A and C are perpendicular. The clarified transverse direction
`Y = -sign * worldY` gives `C × Y = A`, so the `(C,Y,A)` basis preserves
orientation on both arms. This avoids reflecting one sphere's winding when
the same source topology is reused. Globally linear axial position and
transverse smoothsteps meeting with zero transverse derivative give the
stated C1 spine at the nominal sleeve root. The root station must follow the
same axial coordinate, rather than separately resetting each segment to an
unrelated source-sphere fraction. No new numeric station is supplied here.

These are sections perpendicular to the fixed sleeve axis A. They are not
necessarily perpendicular to the changing spine tangent away from the root.
That is a valid coherent loft representation and still addresses the world-X
warp; it should not be described as a transported anatomical frame.

**Positive-station lumen membership.** Retain the original explicit test that
each positive-axial section center lies inside the sleeve lumen. Unsigned
distance to the material shell alone cannot distinguish the hollow interior
from external space; both are outside the cloth material. Crossing-free
centerline chords alone also do not establish that distinction through an
open garment. Use the actual sleeve's sectional boundary/layer evidence to
establish which enclosed section is the lumen, with no assumption that the
previously named mesh half B is the inner wall. The proposed chord checks
then supply additional continuity evidence. This is clarification of the
proposal's membership requirement, not a new numerical constraint.

**Radial fitting.** The updated march is mathematically appropriate. Distance
to the union of the actual sleeve triangles is 1-Lipschitz. At a queried
point with nearest distance `d`, a unit-ray step no greater than
`d - 0.8 mm`, reduced by the stated numerical guard, cannot skip across the
0.8 mm offset boundary. Unlike an arbitrary distance bisection, it follows
the connected safe interval from the section center even when the distance
function is nonmonotone. Querying the full evaluated sleeve also includes
root and cuff rims rather than treating one presumed wall as authoritative.

The method needs a feasible center: if its actual distance is below the
claimed radial margin, collapsing the section to radius zero does not make
that center feasible. Stop and report that condition. At an iteration limit,
stop at the last verified safe radius or report an unresolved fit; never
jump to the raw endpoint or label the offset boundary converged without
evidence. These are conservative termination semantics, not alternate doses.

**Negative-axial root transition.** The coordinator removed the hypothetical
extrapolated-corridor rule. Negative-axial sections now query the real full
sleeve and may fit against its actual rim. This closes the earlier evidence
gap: space behind a nominal root is not assumed free. Preserve the original
proximal pole as specified and verify its actual body containment, while
allowing the interface surface to refit. The historical source body-contact
triangles are no longer a frozen geometric boundary.

**Smoothness and finite surfaces.** A C1 spine does not make independently
capped angular/axial radial samples C1. Nearest-feature switches can produce
ridges or a narrow waist, and connecting separately safe samples can create
unsafe faces. The planned finite-mesh checks must therefore operate on the
actual final triangles; a later smoothing step would invalidate their result
and require renewed checking. Smoothness and stuffed volume must also be
judged in the pixels. No new radius floor or smoothing parameter is invented
here, and no additional sweep is justified.

The clarified claim boundary is sound: 0.8 mm is the conservative radial
point/path fitting margin, while full native zero arm/cloth intersections is
the decisive whole-mesh collision gate. The proposal does not claim a global
0.8 mm triangle-to-triangle gap. Report that distinction so safe sampled rays
are not promoted into a stronger whole-surface clearance certificate.

**Buried connection.** The updated body check is substantially stronger than
testing only two poles or a nonzero overlap. Requiring one connected band of
body-crossing arm triangles, then exactly two edge-connected components after
removing that band, with strict three-ray classification identifying one
inside and one outside the unchanged closed/self0 body, directly tests the
intended continuous attachment. Pole classifications add useful endpoint
witnesses. Its reliability is bounded by complete receiver contacts,
unambiguous classifications, and the declared finite BVH coverage. An
unresolved classification or truncated contact query must not be treated as
a pass. The body's fixed geometry and the remaining actual receiver checks
retain their stated scope.

## Conditions for the recommendation

No further endpoint, root, or radius choice is needed from this review before
the updated contract is frozen. The implementation must retain the explicit
lumen test, conservative march termination, actual negative-root queries,
and finite-surface/body-component checks described above. If the required
centers cannot occupy the lumen, radial fitting collapses the arm into an
implausible shape, or the connected buried attachment fails, reject the
hypothesis using that evidence; do not lower the dose, ignore the interface,
or infer permission to move the cloth.

A successful numerical fit still requires the fixed front, side, and
three-quarter pixels to show the intended high hand, soft continuous arm,
and seated shoulder without exposed clipping or pinching. This report does
not establish calibrated reference acceptance or resolve other asset defects.

## Input identity and scope

- `arm_098_source_attachment.md` SHA-256:
  `8bcbc0c9c76ea3de01aa12aefceea18ad54354c172940bc78403d0868365e26f`.
- Its bound attachment JSON SHA-256, reported by that diagnosis:
  `6d88701d6bd20faa4a2ac622fb58d98871f718dfb9e809c67ef60d912335ce1b`.
- Prior `arm_097_fit_design.md` SHA-256:
  `d0f9ac500def3fd7116f2c35f1bd3edc565aa19f961199cb8b83107c8056cab6`.

The detailed 099 proposal and updated fitting/body rules were supplied by
the coordinator in the current task messages. The coordinator identifies
source 098 as bow-only with the arms and sleeves unchanged from 094; this
review did not reopen or independently hash that model. The earlier reports
bind the inspected 095 images and geometry diagnoses. No geometry arrays
were recalculated, and no Blender/native process or candidate was run.

The linked worktree remains on `t3code/continue-fumo-desktop-use`, HEAD
`c7601f0fc80e0a94b585910459639ffccbdbdbd4`. Only this ignored scratch report
was written. Decision-review supplied the causal challenge and revision of
the proposal; reference-fidelity keeps these geometric checks distinct from
visual acceptance.
