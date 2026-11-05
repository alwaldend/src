# Rear liner feasibility: own-XZ ambient shear and supplied thickness blend

Immutable bounded review, observed 2026-09-06T04:54:24.299268+00:00.
This report concerns a future hair representation, not the arm101 candidate.
No Blender/native process, model change, parameter sweep, or shape-parameter
selection was performed. Only this new report is written. The primary agent
owns the material decision and any subsequent frozen plan.

## Verdict

**Proceed conditionally with representation feasibility, not assembly approval.**
The own-XZ ambient shear fixes the known base6656 liner inversion in the actual
094 rear arrays. Both that shear and the one supplied 500 micrometre lower-gap
blend pass complete finite rear-sheet PL overlap tests, not merely vertex
sampling. The blend also reduces the measured approximately 5 mm ambient heel
gap to approximately 0.5 mm at the counterexample.

These findings do not certify the mixed side/gusset faces, connecting rims,
closed cover, core fit, retained assembly contacts, or appearance. A global
ambient-map argument cannot waive those existing full-surface checks. The
500 micrometre value remains the coordinator's proposed hypothesis, not a
parameter selected or frozen by this review.

## Inputs and exact scope

- Source: `hair_094_repair_arrays.json`, SHA-256
  `492c89b457110d01273a1d8de0616dc266d18c5c441aa634c871395c54946b7f`.
- Associated retained094 model SHA-256:
  `af5a61921ae309a69cfe98b0da092d206b2d0ff29c9123ac2a4886de5a0f9add`.
- Complete-overlay helper: `bow_097_graph_fit.py`, SHA-256
  `754a3a20f350b7c1b5fadf5d6e293ee7e85ffb05763f944965b774fb13584d89`.
- Rear selection uses the original barycentric ancestry: every nonzero support
  belongs to an original vertex ID at least 3298. The actual cover graph has
  3,318 selected vertices and 6,494 fully rear triangles per sheet.
- The inner sheet is taken from the actual source shell arrays, not recreated
  at the paired outer XZ. Source and final diagnostic coordinates are rounded
  to IEEE float32. The two source/helper hashes were rechecked after calculation.

The diagnostic ran as an inline pure-scalar Python calculation, exited zero,
and saved no mesh, candidate, helper, or numerical-shape grid.

## Formulas tested

Distances below are in metres. Let `S` be clamped smoothstep and define:

```text
P(x) = 0.034 - 0.020 * (x / 0.059)^2
w(z) = S((0.166 - z) / 0.038)
h(z) = z - 0.010 * (1 - S((z - 0.085) / 0.040))
G(x,y,z) = (x, y + w(z) * (P(x) - q_old(x,z)), h(z))
```

`q_old` is the actual original rear outer PL graph. Each inner point queries it
at that inner point's own original XZ. At an outer vertex, `q_old` equals that
vertex's source Y. This distinction is essential: paired vertices do not have
identical XZ in the inherited radial liner.

As a continuous map on a valid graph domain, `G` preserves vertical ordering:
its Y derivative is one and `h` is strictly increasing. This is not by itself
a theorem about triangles obtained by mapping only vertices, nor about a
piecewise construction which leaves front-ancestry points fixed.

For the supplied alternative, query the actual mapped outer PL graph `q_new`
at each mapped inner point's own XZ, then use:

```text
inner_blend_Y = (1 - w) * inner_ambient_Y + w * (q_new - 0.0005)
```

The diagnostic uses the original inner Z in `w`; it also checked that this is
equal to using mapped Z here. The Z map acts only below 125 mm, where `w` is
already one. At a vertex, the resulting gap is a convex combination of the
ambient gap and 0.5 mm. This identity does not establish a uniform 0.5 mm
minimum after interpolation. At points originally at or above 166 mm, both
maps are identity; exactness of any larger semantic protected set still needs
its actual vertex-set comparison. The proposed thickness is Y-directional,
not constant normal thickness.

## Measured coverage and finite-sheet results

All 3,318 original inner XZ queries were covered by `q_old`; all 3,318 mapped
inner XZ queries were covered by `q_new`. Point location used nonnegative
barycentrics without clamping an outside query into a triangle.

There were zero nonpositive projected triangles in source outer, source inner,
mapped outer, ambient inner, and blended inner. Each source/mapped rear domain
has 140 boundary edges, one directed boundary loop, and Euler characteristic
one. Exact predicates on the binary coordinates found zero nonadjacent boundary
self-intersections on either sheet and zero outer/inner boundary intersections,
both before and after mapping. Together with the unchanged disk topology,
positive projected orientations, and covered inner vertices, these checks
support containment of the entire inner projected disk, not just its vertices.
This is rear-chart evidence, not a full-shell topology audit.

| Measurement | Ambient shear | Supplied blend |
| --- | ---: | ---: |
| Sampled inner-vertex minimum Y gap | 0.459886500 mm | 0.483617440 mm |
| Sampled inner-vertex maximum Y gap | 5.072855988 mm | 3.357133009 mm |
| Complete PL overlap minimum Y gap | 0.459356111 mm | 0.483617440 mm |
| Minimum minus helper's 10 nm guard | 0.459346111 mm | 0.483607440 mm |
| Faces with positive guarded inversion deficit | 0 | 0 |
| Base6656 Y gap at inner's own mapped XZ | 5.009120225 mm | 0.500000908 mm |

The maximum rows are vertex samples, not claimed global PL maxima. The minima
labelled complete are computed over all overlapping finite triangles.

For each variant the existing helper examined 54,139 projected candidate pairs:
20,855 empty, 233 point contacts, 222 segment contacts, and 32,829 area overlaps.
It used 24,638 exact-predicate fallback pairs. Its plane-arithmetic error
estimate was `4.440892098500626e-16 m`, below its `1.25e-9 m` declared budget;
the reported conservative margin subtracts the separate `1e-8 m` guard.
This is a guarded numerical certificate with exact clipping predicates,
not an exact-rational evaluation of all interpolated Y values.

The ambient minimum occurs at XZ approximately `(1.980000176, 133.000001311)`
mm, inner face 3209. The blend minimum occurs at approximately
`(-1.980000176, 146.889999509)` mm, inner face 2756. Thus the blend's full-sheet
minimum remains below 0.5 mm in its transition region even though it is positive.

## The base6656 counterexample and heel risk

The earlier same-paired-displacement proposal put the inner skin approximately
36.6615 micrometres beyond the outer skin at the inner point's own XZ. Moving
the head core inward could not repair that shell-internal inversion.

In this float32 diagnostic, base6656 maps to local vertex 5196. Its inner point
has mapped XZ `(-9.259772487, 70.560000837)` mm, where the actual mapped outer
graph has Y `33.503301336` mm. Ambient inner Y is `28.494181111` mm; blended
inner Y is `33.003300428` mm. Both are correctly inward, with the gaps in the
table above.

The ambient result is therefore mathematically ordered but potentially looks
like a thick padded heel: approximately 5 mm here and 5.073 mm at sampled
maximum. The supplied blend addresses this measured lower-heel risk without
independently moving paired XZ. It retains larger transition/upper gaps,
including the sampled 3.357 mm maximum; no visual cloth judgment follows from
ordering alone.

## Reusing the bow helper

Yes: the existing helper is suitable for this bounded complete PL ordering
audit. Pass mapped outer triangles as the "bow", rear inner triangles with
the required +Y receiver winding, and `gap_m=0`. Its unguarded deficit is
`inner_Y - outer_Y`; the largest such value over every clipped intersection
vertex is the negative minimum sheet gap. The computation includes point and
segment intersections. Read the unclamped witness values: zero clamped fit
bounds alone report no required correction, not the achieved clearance.

The helper does not establish receiver coverage or graph uniqueness; the
separate domain checks above matter. Nor does it audit faces omitted from its
receiver and subject lists.

For an inward correction of the inner sheet, mirroring both sheets' Y converts
the inequality into the helper's outward-fit direction:
`(-inner_Y) + delta >= (-outer_Y) + gap`, equivalently
`inner_Y - delta <= outer_Y - gap`. Verify the receiver's required projected
winding explicitly; mirroring Y alone does not change XZ signed area. A
vertex-max majorant can then bound the final finite correction on complete
overlaps. Its unconstrained result does not automatically preserve protected
upper vertices: nonzero required corrections there would be a conflict, not a
reason to silently zero the correction. No such correction was built here.

## Remaining coupled-boundary and assembly questions

The rear chart shares 114 vertices with 260 mixed-ancestry cover faces. The
proposed outer map moves 77 shared vertices in Y and 52 in Z. At base6656 the
outer motion is approximately +42.757 mm Y and -10 mm Z. This is a consequential
side/gusset deformation, not a hidden thickness adjustment.

Allowing those shared nongraphic transitions to deform while keeping original
front-ancestry and graphic vertices exact removes the earlier contradictory
freeze. Shared topology must remain shared; duplicating the join would not
solve the geometry. However, ordered rear sheets can still have intersecting
rims or mixed gussets. The rear certificate does not decide chin coverage,
exposed stuffing, sidewall folding, or whether the cloth reads as free ends.

The next implementation still needs its already required full native cover
self/contact checks, actual core/cover separation and exterior relation,
protected front/graphics and retained assembly-contact comparisons, and fixed
view inspection. The proposed core support fit must use the actual new inner
chart, include every overlapping branch, and preserve core Z. Empty projected
overlap is not evidence of coverage, and an inward core fit does not replace
the cover's own closed-surface checks. No new acceptance criterion or authority
to change additional geometry is supplied by this review.
