# 099 continuous sleeve-corridor arm fit

Decision-review PROCEED. Use retained098 bow_098_candidate.blend
2bee72c54ba571662e36fead903c9516ecb0f13b04d9138fbe897e0c1ffb6193.
Only the two continuous stuffed arm meshes may change. Body, sleeves, dashes,
all hair/head/bow, pile, materials, cameras and other objects remain exact.
Root sole model/goal/Git writer. Body-interface locus may change: preserving
complete old interface triangles was disproved as compatible with clearance,
and is not a user requirement. Preserve a meaningful continuous buried root.

Final arm_098_source_attachment.md cad825ec17c11891ade2a5b7bf3c2586d792b72920447b3aec79a06032fa5f5e
and corridor500dba0b129c1324e8348dbce1e45952a40ce19bed87465c83567b5b95648f21
support actual center routes: direct endpoints4.347505mm and retained sleeve
center route7.608069mm minimum full segment-to-wall distances, zero crossings.
These are not complete arm fits. arm_099_design_review.md recommends this
updated hypothesis; its source-report hash names an interim wording revision,
so root separately read the final report and preserves the final numeric
evidence and relaxed interface authority here. No additional arm dimensions
are inferred from photographs of a hidden shoulder.

Reparameterize original unit UVsphere, preserving topology/material ancestry.
Use sleeve basis A=normalize(sign*.82,0,-.572), C=normalize(sign*.572,0,.82),
Y=(0,-sign,0), so C cross Y=A. For s=(localZ+1)/2, axial coordinate interpolates
from original proximal pole to the useful095 distal pole (old distal plus
(-sign*12,+8.5,+14.3)mm). Nominal sleeve root is (sign*24,-3,70)mm, axial0.
For negative axial coordinates interpolate proximal transverse offsets to
zero with clamped smoothstep; for positive coordinates interpolate zero to
distal transverse offsets likewise. This is one C1 spine through the root,
with sections perpendicular to fixed sleeve A, not a transported frame.
Raw transverse coordinates are source localX*9mm*C + localY*10mm*Y.
The proximal and distal poles remain exact endpoint targets.

Verify every positive sampled section center lies inside BOTH actual nested
wall contours cut by its plane; use full evaluated geometry, not an assumed
inner mesh-half label. Require two closed degree2 section loops, nonzero
areas and unit winding about the center, with no ambiguous plane vertices.
Verify all centerline chords against the full sleeve. Negative stations also
query the actual full sleeve/rim; no extrapolated empty-corridor claim.

Fit each raw radial segment from its center by conservative distance stepping:
target sampled gap0.8mm with0.2um extra numerical guard; advance at most0.9
times(current full-sleeve nearest distance minus guarded gap). End at raw
radius or when the remaining margin<=0.1um; at256 iterations fail unresolved.
Initial center must have the margin. The1-Lipschitz distance bound prevents
skipping an unsafe interval even when unsigned distance is nonmonotone.
Keep the last safe sample; no post-fit smoothing or shape sweep. Report radii,
contractions, iterations, and sampled gaps. Reject any radial scale below0.6
as an explicit implausible-thinning guard, not a reference measurement.

Complete native arm meshes must be closed consistent positive single
components with zero self/duplicate/degenerate triangles. Evaluate ALL visible
mesh receivers: only actual body intersection is allowed, no inherited sleeve
crossing waiver. Body-crossing arm triangles must form one edge-connected
band; removing it leaves exactly two components. Classify robustly chosen
component witnesses with strict three-ray parity against closed/self0 body:
exactly one buried and one exterior, zero ambiguity. Proximal pole inside,
distal outside. Preserve source bytes and all nontarget fingerprints.
The0.8mm claim covers fitted radial samples/paths only, not a certified global
triangle gap; actual whole-surface zero intersections is mandatory.

One pinned5.2.1/four-thread build, one causal API/data repair maximum, no
post-render numerical edit. Save only after native gates. Clean-reopen fixed
five and presentation; primary plus image-only review before retention.
Require high contained hands, attached shoulder, no exposed pinching/clipping,
and no silhouette regression. This does not pass whole clay/asset acceptance.
