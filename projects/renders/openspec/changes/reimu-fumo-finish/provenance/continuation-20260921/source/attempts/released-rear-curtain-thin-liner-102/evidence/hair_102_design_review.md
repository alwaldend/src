# 102 static geometry review

Observed 2026-09-06, before the coordinator's first 102 native execution.
This is a bounded independent critique under the decision-review skill. The
coordinator owns the numbers, final verdict and native/model/goal writes.
Only this report was created. No native process, candidate, helper change,
parameter grid or new scalar shape trial was performed.

## Advisory conclusion

PROCEED with the existing gated trial after the confirmed core-set correction.
The mirrored-Y core fit has the right inequality and finite-triangle bound.
The coordinator's corrected builder distinguishes the complete rear core
from the covered rear chart, addressing the concrete omission found here.
No static contradiction now requires a different shape or fitting method.

This recommendation is not assembly approval. The 101 evidence concerns two
rear sheets, while mixed faces, the core, laps and retained objects remain
dependent on the already frozen full native gates and visual review. The
lap displacement is a construction rule, not a proved clearance from the
actual newly triangulated cover.

## Reviewed inputs and scope

- Frozen `hair_102_plan.md`:
  `58daf7b84b70354ac2c07a40e3fe7253a1ce6fc36b9c293389e1bf4d3b6a33a5`.
  The coordinator subsequently confirmed canonical pre-execution
  clarification RV338 for the core/cover set distinction below.
- Final `hair_101_liner_feasibility.md`:
  `aee1c11b349c109134e9439c6a69c48a6627a443c2b1f1dcdaaf3e89cc22697f`.
- Immutable `bow_097_graph_fit.py`:
  `754a3a20f350b7c1b5fadf5d6e293ee7e85ffb05763f944965b774fb13584d89`.
- Exact `hair_094_repair_arrays.json`:
  `492c89b457110d01273a1d8de0616dc266d18c5c441aa634c871395c54946b7f`.
- Exact `hair_096_repair_arrays.json`:
  `bee174d66dcf5f765dcdfd99b1351439c387676e822b9d4ded67eb04df2bed1b`.
- Corrected `hair_102_build.py`:
  `bdee174c09316be70635a52879a7801fd523843daa342e2f72a4693f8bbb53a1`.
  Inspection was limited to source-set selection, cover/liner input dataflow,
  core fitting and preservation assertions. This is not a full builder audit.

The formulas, 0.5 mm liner, 0.3 mm core gap, 10 nm arithmetic guard, protected
sets and allowed object changes are the coordinator's frozen choices. This
review supplies no authority or additional acceptance criteria.

## Core inward fit and winding

Let C be the old core Y and I the actual new inner-cover Y at an overlapping
XZ point. Mirroring both gives the helper's unguarded deficit

```text
(-I) - (-C) + gap = C - I + gap.
```

For each core triangle it bounds that affine deficit over every vertex of
every triangle-intersection polygon, including point and segment contacts.
Its per-vertex displacement is the maximum bound of every incident subject
triangle. Thus the interpolated displacement on any subject triangle is at
least that triangle's uniform bound. Applying `coreY -= delta` produces the
intended `newCoreY <= innerY - gap`, with the separate arithmetic guard.
Core XZ remains unchanged, so this application does not change the projected
overlap domains used by the proof. The corrected builder uses this sign,
applies the array once and asserts exact X/Z preservation and inward-only Y.

Y reflection does not change signed XZ area or the normal's Y component.
Actual inner-shell triangles normally have the opposite winding from outer
triangles; mirroring alone would not satisfy the helper's +Y receiver rule.
The inspected code correctly uses the outer rear graph's triangle order on
the actual inner coordinates for the helper. Native closed-shell triangle
order remains unchanged. Source bow/core winding may have either sign, but
receiver graph winding must remain strictly +Y.

The majorant is deliberately conservative and can move a shared vertex more
than its own local deficit. Clamping a forbidden displacement or smoothing
the returned array would invalidate the bound; the plan correctly disallows
that. A front/protected-set conflict must remain an assertion failure.
The helper's double arithmetic is guarded, not certified interval arithmetic;
finite native coordinate conversion and closed surfaces retain their planned
checks.

## Corrected source sets and coverage limits

Static inspection of the exact 094 packet establishes:

| Subject or receiver set | Triangles | Referenced base vertices |
| --- | ---: | ---: |
| Covered rear chart / inner receiver | 6,494 | 3,318 |
| All rear-ancestry core subject | 6,528 | 3,326 |
| Additional core support beyond covered set | 34 | 8 |

Rear ancestry means every nonzero original support ID is at least 3298.
The first row additionally requires `face_covered`. The extra core vertices
are below the old cover cutoff; the complete rear core's minimum original Z
is 79.448133707 mm, versus 81.12 mm for its covered subset. Their being outside
the old covered vertex set is not evidence that they are outside the new
extended receiving chart.

The initial builder reused covered faces as the core subject. The corrected
builder now separately constructs all 6,528 core faces, asserts 3,326 core
vertices, queries every one at its own fixed XZ and passes all those faces to
the helper. Cover construction and receiver topology correctly retain the
6,494/3,318 selection. Front ancestry is independently selected by any nonzero
original support below 3298, asserted disjoint from the core move set, and
compared exactly after fitting. Mixed nongraphic faces remain shared and may
deform through their rear vertices, as expressly allowed by the plan.

An empty overlap returns zero displacement in this helper. The explicit
receiving-graph queries therefore matter. Covered corner queries alone do
not prove that every interior point of a projected triangle lies in a
possibly nonconvex receiving domain. The helper establishes its complete
bound on actual overlaps; it does not prove graph uniqueness or full-domain
containment. Report support with that limit and retain the frozen full
core/cover contact, closed/self and exterior-component checks. No additional
coverage diagnostic was run here.

## Liner preservation and lap-to-cover mapping

The strongest evidence for this representation is the final 101 complete PL
minimum: 0.483617440 mm for the own-XZ blended liner, or 0.483607440 mm after
the 10 nm guard. That exceeds the plan's 0.45 mm rear-sheet threshold. The
same-paired-outer displacement has a documented inversion counterexample;
the inspected 102 dataflow instead queries each actual original inner point
at its own XZ, constructs its ambient point and queries actual new outer
triangles at the mapped inner XZ. Upper and non-rear coordinates are copied
and compared exactly. This preserves the supported 101 construction.

Neither that rear-sheet result nor the continuous-map argument certifies
mixed gussets, rims or a closed cover. The plan retains those checks and
acknowledges the substantial allowed lower transition motion.

For a source lap point `(x,y,z)`, write F for the proposed curtain Y map,
d for the original sign-dependent 096 shift, and A for the lower drop weight.
The lap's final inner point has

```text
Y = F(x,y,z) + d(z) - 0.0005,
Z = z - 0.004*A(z).
```

Its gap from the actual new cover is therefore

```text
F(x,y,z) + d(z) - 0.0005 - q_new(x, z - 0.004*A(z)).
```

It is not simply `d - 0.0005`: the central cover uses the different 10 mm Z
drop, and nonlinear evaluation at clipped lap nodes is not generally the
same as barycentric interpolation of the mapped original cover triangles.
The old lap ancestry guarantees correspondence, not clearance after these
different maps. Actual lap/cover/lap and retained-bow triangle gates already
required by the plan are decisive. Replacing the lap Y with a fresh q_new
offset would be a different construction, not a silent implementation fix.

Static source counts also show 54 left and 52 right lap triangles straddling
Z=166 mm. Exact preservation of their upper vertices does not alone preserve
every upper point of those triangles when a lower corner moves. Thus
"upper exact096" is supported as a vertex guarantee; retained bow safety
must use actual final triangles. The plan already requires that comparison.

For pile transfer, preserved source ancestry should resolve the anchor's
original patch membership, then barycentrically interpolate the final lap
triangle. Direct analytic deformation of an interior anchor can disagree
with the finite mapped triangle. The plan's actual anchor-residual checks
and unchanged front/other strands address this distinction.

## Decision trade-off and handoff

Reusing the mirrored helper gives a complete overlap majorant with bounded
execution cost and no new fitting machinery. A vertex-only correction would
lose edge-intersection extrema; the previously disproved paired liner move
would restore a known inversion. Neither is a supported substitute.

Proceeding with the corrected sets and the existing fail-before-save gates
is supported. Failure of actual receiving coverage, preservation, shell
ordering, closed/assembly geometry or required views changes that conclusion;
it does not authorize a new shape search or a post-fit correction. No 102
native result or asset acceptance is claimed here.
