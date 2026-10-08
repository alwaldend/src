# 036 closure decision: boundary and visibility evidence first

Plan only. Keep both successful cylindrical patches and their source
boundary edges fixed. Do not select a cap shape until the boundary itself
has been checked against the surrounding sheets. Root requested the next
diagnostic and owns its execution; it has not been run here.

## Causal conclusion from existing arrays

The failing primitive is the straight outer-to-inner termination cap,
not the cylindrical solve. The 116 arbitrary brown-cover cut edges were
closed with 232 triangles spanning the captured hood thickness. That
constructs a new wall at a cut made inside an overlapping cloth assembly;
it does not follow an established sewn edge or check which side of the
fringe inner sheet each boundary lies on.

The 307 raw records comprise:

- 86 positive-length hidden-cover/fringe-inner contact pairs and 210
  one-point reports;
- four positive-length inner-chart/cover-closure pairs and four one-point
  reports;
- three one-point left root-cap/retained-hood reports.

Positive length here means the recorded endpoints are more than 0.2
micrometers apart. A one-point report is neither automatically numerical
noise nor evidence of a finite penetration segment. There are no recorded
contacts involving the new outer chart.

The left root disk was also reduced to two chord triangles. Contacted
triangle 75938 uses global array vertices `[14199, 2914, 23411]`, joining
the fringe inner/outer endpoint to the hood outer endpoint. Its contacts
with retained hood source triangles 15300/15301 lie about 0.049 mm from
the nearest triangle edge; the contact with triangle 15303 lies about
0.093 mm from the nearest edge. Thus the recorded points lie in the new
cap interior, not at preserved shared vertices or edges. A planar chord
across this nonplanar root wedge can enter the retained hood even when the
two fixed chart surfaces are well-oriented.

All contacted obstacles are retained or fixed chart geometry, but the
intersecting caps are new. Their indices alone do not prove the original
source sheets intersected at those positions. No unavoidable inherited
overlap is claimed. The two zero-clearance reports are separately known
to involve original rim polygons, not an inner-skin collapse.

## Exact closure-local boundaries

The cover closure is the 116-edge path stored at
`head_035_domain_probe.json.objects["Hair028 crown and back hood"].front_cover_cut_edges`,
plus its captured inner counterpart with source vertex offset 9,114.
Its outer source endpoints are 793 and 841. No chart edge may be split or
moved as a workaround; their CDT preservation guard remains active.

The two root disks use these fixed outer/inner endpoint loops, in global
036 array indices (actual orientation follows the existing face winding):

- Left: `[2914, 23411, 32525, 14199]`.
- Right: `[3612, 23363, 32477, 14897]`.

These 236 closure triangles are the initial change scope. No head core,
fringe tip, eye opening, neck, bow, body, or foot edit is proposed. The
2,140 retained brown-cover faces are an implementation safeguard, not an
independent acceptance requirement; a hidden redundant subset may be
proposed for removal only with explicit new ownership and coverage proof.

## Next authorized diagnostic, not another solved state

`head_036_closure_diagnostic.py` is prepared for root to read and execute.
It reconstructs immutable helper SHA
`8b2dae5451e643a59ae49e81f3b4d76249e7e0c32be161834c3e989b1cd46fb9`
on the same frozen032 input. A diagnostic callback captures full arrays
then invokes the original checks. It requires exact equality with the
recorded first-state new arrays and contact records. No helper bytes or
geometry parameters change; no scene object is needed.

It will produce:

1. Full candidate arrays and original source-triangle correspondence,
   with explicit outer-skin, inner-skin, and original-rim classification.
2. Every fixed closure boundary edge, derived from incidence against the
   retained/chart surfaces, tested against all other fixed sheets. Shared
   topological edges/vertices are distinguished from crossings. An
   original retained edge against an original retained face is sufficient
   to identify inherited boundary contact without calling a new cap
   interior contact inherited.
3. Contact dimensionality and distances to fixed closure boundaries.
4. Geometric first-hit visibility of contact endpoints and finite-segment
   midpoints through all five contract views. Other visible render-enabled
   scene geometry participates; replaced fringe/hood originals are
   excluded. The replacement arrays enter the ray scene directly, without
   creating or hiding an object. This is visibility evidence, not a render
   or proof that an entire cover face is hidden.

Outputs are `head_036_closure_diagnostic.json` and
`head_036_closure_diagnostic_arrays.json`, created without overwrite.
Run with the same pinned Blender and single-threaded numerical settings as
the original036 execution. No second chart state or cap trial is authorized.

## Conditional construction choice for root

If the fixed boundaries are clear, the smallest distinct construction is
a closure-only fair membrane constrained to remain on the free side of
the measured receiver sheets, replacing the two-triangle thickness chords.
It retains all boundary coordinates and both solved chart meshes. Its
interior is determined by boundary tangency and obstacle inequalities,
not an arbitrary depth dose. The diagnostic must first demonstrate a
consistent free side; no collision-free membrane is asserted yet.

If the fixed boundary itself crosses a sheet, bending a cap interior
cannot remove that crossing. Do not try successive folds or offsets.
Root must instead choose a local surface-ownership change: identify the
redundant brown-cover piece adjoining the trapped hem and establish that
its removal cannot expose pale core or alter a protected silhouette or
face opening. Contact-point occlusion is insufficient for that claim;
candidate face-area first-hit coverage and coverage of any newly exposed
gap are required. The successful chart surfaces remain fixed. No whole
cover deletion or support expansion is pre-authorized by this plan.

## Successor verification requirements

Preserve the raw036 helper/evidence. Correct nearest-inner selection by
source polygon class, explicitly excluding original rims and closures.
Test every new closure triangle, not only chart faces: finite area,
consistent shell winding, orientation relative to its intended material
interior/free side, and retained/closure as well as closure/closure join
angles. A cap does not generally have a radial normal, so chart-radial
orientation alone is not its validity test. Report the full angle
distribution and any reversal or fold; do not hide it behind manifoldness.

Recheck actual local intersections, distinguish inherited fixed-boundary
contacts from new penetration, verify both retained chart arrays exactly,
and require first-hit coverage plus the five unchanged views before any
visual retention. No helper edit, geometry, save, render, or delegation
occurred during this plan task.
