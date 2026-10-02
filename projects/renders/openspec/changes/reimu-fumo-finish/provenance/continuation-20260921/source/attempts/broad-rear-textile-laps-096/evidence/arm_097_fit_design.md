# Arm 097 fit design review

Observed 2026-09-06 02:51 UTC. Recommended verdict: **revise the arm's
transition representation while preserving the useful distal placement and
measured buried shoulder contact**. A continuous centerline-based arm with
fitted transverse sections is the smallest credible next hypothesis. The
coordinator owns the verdict and any implementation. This report supplies
no new coordinates, contraction factors, candidate, or parameter grid.

## What the failed 095 evidence establishes

I read `arm_095_collision_diagnosis.md`, `arm_095_plan.md`, and
`arm_095_result.md` in full and inspected the fixed diagnostic front, side,
and mirrored three-quarter pixels. The previously inspected controlling
canonical front and side frames 10/12 remain the reference evidence. This
review has implementation context and is not an independent visual pass.

In the 095 front image, the conspicuous outboard hand protrusions disappear.
In the side image, the hand occupies the upper opening with cloth below it.
That is useful evidence for distal placement. The known crossings still
reject the surface, and the images alone do not establish calibrated likeness.

| Measured fact | Consequence for the design |
| --- | --- |
| All 28 source cuff-return/rim crossing pairs per side disappear. There are no failed cuff-region pairs or pairs on arm triangles wholly in the fully translated source abs(X)≥58 mm region. | Preserve the successful distal placement as a working target; changing the cuff to solve the new failure has no present support. Zero crossing is not a full containment/minimum-gap certificate. |
| The actual distal pole is (±56.978, −5.026, 53.993) mm; its unsigned distance to the entire sleeve is 11.054 mm. | The distal pole itself has space. Its point distance does not describe every distal vertex or the midarm. |
| The dominant 56 new pairs per side lie at actual sleeve rows 11.630–14.366, upper/rear, abs(X)=40.682–43.398 mm, Z=69.180–72.460 mm. | Fit the lifted transition inside this measured part of the retained sleeve. This is not a root-ring or cuff-endpoint failure. |
| Those contacts originate on source-arm surface points at abs(X)=49.771–54.845 mm and have applied blend weights 0.757–0.954. | The translation field moves transverse arm material inward/upward/rearward into the cloth. Clearing the distal tip alone cannot resolve it. |
| Seven additional lower/front pairs have weights only 0.000416–0.005807 at rows 12.986–15.000. | Nearly frozen is not unchanged or clear. A smaller translation weight is not a sufficient corrective principle. |
| Source sleeves already cross the arms in 310/296 triangle pairs. Failed 095 retains 244/228 pairs on wholly frozen arm triangles. | Clearing the 63 newly added pairs alone cannot establish a usable whole arm/sleeve fit. The inherited crossing geometry also needs an explicit disposition. |

All failed crossings end by row 15.036; the actual root boundary is row 0
and the outgoing cuff is row 40. The diagnosis's construction quadrants and
row correspondence are measured, but its half B is not established as the
inner wall. Pair counts depend on triangulation; their segment locations,
not just the totals, should guide the fit.

## Why a continuous fitted arm is credible

The failed 095 field translates individual vertices by a smoothstep of
world |X|. That preserves topology and passes self-intersection checks, but
does not preserve or control cross-sections perpendicular to the arm's
direction. Different points on one oblique stuffed-arm section can receive
different blend weights. The observed upper/rear migration and almost-fixed
lower/front contacts are consistent with this field producing an unsuitable
transition envelope. The evidence does not require a different distal target.

A plush arm can be a short, continuous stuffed pod with a gently bent or
flattened middle. A narrowed section beneath a sleeve is physically credible
when it remains smoothly filled and joins the shoulder and hand without a
pinched neck, hard elbow, or rigid peg. The reference does not reveal an
exact hidden arm profile, so that profile should follow measured available
space and continuity, not invented anatomy.

The next representation should control a centerline and its transverse
envelope together:

1. Preserve the actual body-contact geometry required for the buried
   shoulder, and preserve the successful distal pole/visible bulb placement
   as working targets. First verify which distal surface region has adequate
   finite clearance. Do not infer complete clearance from its pole distance.
2. Parameterize intermediate arm sections by axial position along the
   source arm or a continuous replacement centerline, rather than each
   vertex's world-X coordinate. Connect the retained shoulder to the distal
   region with a smooth centerline inside the available sleeve corridor.
   Section frames should vary smoothly so the arm does not twist abruptly.
3. Fit the intermediate transverse extent to the evaluated sleeve interior.
   Use the measured upper/rear and lower/front contact zones to identify
   where that extent must change, then extend the check over the entire
   arm/sleeve overlap, including the inherited contacts. Moderate flattening
   can retain a soft stuffed profile while removing upper/rear excess. A
   symmetric global shrink is not justified merely because most new pairs
   occur in one quadrant.
4. Join the fitted sections continuously into the preserved shoulder and
   distal regions, with smooth radii and tangent changes. This can be a
   reparameterization of the existing connected mesh; a fitted loft does not
   inherently require disconnected pieces or more topology. Evaluate the
   final finite triangles for crossing, foldover, and self-intersection.

A transverse squeeze about the failed centerline is a sufficient special
case only if that centerline already lies inside a usable corridor and the
required reduction leaves a plausible stuffed section. If the centerline
itself exits the available interior, reducing radii around it cannot solve
the problem. In that case, adjusting the centerline within the same continuous
arm representation is necessary; this is not a reason to try a grid of squeeze
factors. The current reports do not locate that centerline relative to the
interior, so they cannot freeze a numerical squeeze or profile.

## Preserved shoulder versus preserved clipping

The 095 plan froze every native arm vertex with world |X|≤32 mm; the
result reports 459 protected vertices per arm and exact unchanged body
contacts. That broad coordinate band is an implementation constraint. It is
not equivalent to the subset of geometry required to preserve the real
buried body attachment.

Before fixing a new protected region, classify the recorded source/failed
arm-sleeve crossing segments against the actual unchanged body-contact
triangles and receiver surface. This requires actual triangle and contact
evidence, not a guess from the hidden reference throat. The question is
whether an offending sleeve segment lies on a triangle that must remain
fixed for the source body connection, or on adjacent free arm surface that
was frozen only by the |X| band.

If a sleeve-crossing triangle was frozen only by that broad band, allowing
its free portion to join the continuous fit can preserve the meaningful
shoulder connection without preserving the same clipping. If an exposed
midcloth crossing lies on a triangle that really must stay exactly fixed,
then exact preservation of that triangle and removal of its crossing with
an unchanged sleeve are incompatible. An arm-only design cannot claim to
solve that case. Surface labels such as buried, source, or frozen do not by
themselves waive sleeve crossings, and valid source arm/body overlap does
not grant a blanket arm/cloth intersection allowance.

This classification is the feasibility boundary for the proposed arm-only
fit. It should precede choosing new numerical geometry, because an overly
broad frozen region could make the task impossible before fitting begins.

## Comparison with changing the cloth

| Option | Benefit | Decisive limitation |
| --- | --- | --- |
| Repeat the world-X blend with a smaller displacement or altered weights | Small code change. | Changes the useful high-hand placement or redistributes crossings without modeling the available volume. The diagnosed failure does not support another dose. |
| Squeeze intermediate sections about a verified interior centerline | Smallest deformation when only transverse excess causes contact; can preserve terminal regions and mesh connectivity. | Needs measured space and full-contact checks. A fixed outside centerline or an impossible frozen contact patch cannot be repaired by squeezing. |
| Fit a smooth centerline and continuous transverse loft | Independently controls arm route and girth while retaining shoulder and distal placement. This directly addresses the diagnosed transition envelope. | Requires local interior identification, smooth joins, and whole-surface verification; excessive thinning would cease to read as stuffed fabric. |
| Reshape the retained cloth locally around the midarm | Can resolve a corridor that cannot contain a credible arm or an incompatible required shoulder patch. | Changes an externally visible garment and risks its silhouette, drape, cuff orientation, and attachments. The present evidence does not establish that it is necessary. |

Prefer the continuous arm fit if the contact classification and measured
interior admit it. If they do not, the smallest justified cloth change would
address the measured contact region while preserving the established root
and cuff where feasible. Relocating either endpoint has no new supporting
evidence. Any cloth change must preserve a credible fabric thickness and
alter inner/outer surfaces consistently; silently erasing an inner wall or
thinning it to bypass collision would not resolve the construction problem.

The practical acceptance condition is a high contained hand, continuous
buried shoulder, no visible pinching, and a clean arm/cloth fit over both old
and new contact zones in the actual evaluated geometry. Fixed front, side,
and three-quarter pixels must confirm that the geometric repair preserves
the placement benefit. This bounded result would not clear inherited
skirt/hem defects or confer whole-asset acceptance.

## Evidence identity and scope

All following paths are relative to this report. Inputs were read in full;
there were no new native runs, scalar geometry calculations, sweeps, model
edits, goal writes, or Git mutations. Only this ignored report was written.
The previously verified linked worktree remains on
`t3code/continue-fumo-desktop-use` at
`c7601f0fc80e0a94b585910459639ffccbdbdbd4`.

| Inspected input | SHA-256 |
| --- | --- |
| `arm_095_collision_diagnosis.md` | `9d2b7009955ce9505ba84ef99ed764f0ed8fa060f544fe2f3598b9a4c546aaf7` |
| `arm_095_plan.md` | `92a645e158756be40049086ae2904ea82137340d5b2a7bc0f8d6e0854859ea58` |
| `arm_095_result.md` | `3a848e77f77fee5f3a1c1d28c6b01a1e210aaf3d6974e84f6aab303ed0994d21` |
| `arm_095_failed_review/front.png` | `21b98371160785b4ec4e05f260d129da5c734e3afa4149272e35df358a0115c5` |
| `arm_095_failed_review/side.png` | `0da131bcf11912a72c4355f53c2ee39314b9b044c465ca5a4ffb16010bd7f205` |
| `arm_095_failed_review/three_quarter_mirror.png` | `90bfb53f9b6a792e53610ceec4799925fedcf6075969f7bcdd588b14ba999c6b` |

The collision report binds the failed arrays, preflight, and retained sleeve
export; this reviewer did not independently recompute those measurements.
Earlier image identities and the corrected retained-root provenance are in
the 094 review and `sleeve_095_assembly_correction.md`. The decision-review
and reference-fidelity skills informed the rejection of repeated scalar
tuning and the distinction between placement evidence and surface acceptance.
