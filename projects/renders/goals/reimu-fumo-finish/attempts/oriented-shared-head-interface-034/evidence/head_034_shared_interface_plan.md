# Head034 shared-interface proposal — no geometry trial authorized

Recommend one welded outer cloth surface joining the retained fringe to the
retained hood/gusset. Remove the redundant sheets in the transition before
adding the shared strip. Do not attempt another overlapping-sheet tuck.
Root owns the thickness choice, final scope and authorization to implement.

## Exact frozen topology

Source: `head_032_candidate.blend`, SHA256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
The read-only `head_034_shared_topology_probe.py/.json` contains complete
source face/edge index sets. Pinned Blender 5.2.1 LTS, build `9e2066aef7ef`;
source bytes unchanged. The probe classified one cut proposal only: no
deletion, join, modifier change, helper variant, render or blend save.

The only proposed replacement targets are:

```text
Hair028 traced padded fringe
Hair028 crown and back hood
```

All other objects, including the head core, face, locks, rear cloth, neck
contacts and bow, remain protected. A later foot candidate may supply the
scene only if all 15 head inputs match the frozen032 guards.

| Frozen face set | Count | Remaining outer faces | Cut topology |
| --- | ---: | ---: | --- |
| Fringe upper-root strip | 1,013 | 9,931 | One 337-vertex inner cut arc |
| Hood redundant upper-front strip | 960 | 7,970 | One 212-edge closed hole |
| Hood side of that hole adjacent to gusset | — | — | One 97-vertex outer receiving arc |
| Other side of that hole | — | — | One 117-vertex retained-front-cover arc |

Both removed sets are single connected patches. There are no cut branches
or islands. The fringe cut endpoints are source vertex IDs 3620 and 4318;
the hood/gusset arc endpoints are 793 and 841. Exact ordered paths can be
traversed from these endpoints through the degree-two edge graphs in the
JSON; do not invent a different nearest-vertex correspondence at execution.

The fringe set is selected from the existing 307-vertex upper outer boundary
arc, guarded by Z > 135 mm and elliptical radius > 0.90. Faces wholly above
135 mm whose mean graph cloth distance is below 7 mm are removed. Its
source-index list digest is
`cff1ce4c829289ccbc9113c6eada529874d45cc400a6e4e7f28aa9fabe4ac15d`.
The removed patch begins at Z=135.4333 mm. The original upper open arc has
306 edges; its replacement inner cut has 336 edges.

The hood correspondence is verified exactly against the retained core
front/rear/gusset recipe. Only hood faces corresponding to original brown
front faces, wholly above 135 mm and with mean core elliptical radius at
least 0.86, enter the removal set. Its list digest is
`75c0c406227065e61422cc277009e49f0348104897094dec63951cfd0bb6363f`.
The patch starts at Z=135.8979 mm. The receiving arc has 96 edges adjacent
to the retained gusset; the other 116 cut edges face retained brown cover.
No gusset or rear face belongs to the removal set.

These are frozen coarse cut paths, not a claim that an arbitrary connecting
surface will have the desired silhouette. Endpoint caps and transition
coverage still require construction and first-state verification.

## The brown cover dependency must remain

Exactly 2,140 original front brown hood faces remain outside the local cut.
Preserve their coordinates and material, including the inner forehead cover
behind the fringe openings. Do not delete all front hood faces, recolor the
core, or remove a whole upper cap because the fringe usually occludes it.
Those shortcuts can reveal pale core through the eye/forehead cutouts.

The remaining front cover has an internal 116-edge boundary after the local
cut. It is not a second visible crown interface: keep it underneath the
retained fringe and close its thickness there. Its exposure would fail the
first-hit/pale-hole gate. Existing overlap away from the local transition
is preserved; no broad Boolean union or interior cleanup is required.

## Welded transition construction

1. Preserve the two source objects as recoverable hidden donors only after
   a replacement validates. Snapshot their evaluated outer and inner skins,
   materials and Head attachment before any edit. Preserve all retained
   geometry outside the cut and at most one adjacent seam-support ring.
2. Remove the two exact local face sets. The fringe inner arc and the hood
   gusset arc are the visible strip's two ends. They are connected by new
   cloth faces, not by co-located independent surfaces.
3. Orient both arcs consistently from one side endpoint to the other.
   Parameterize each by arc length. Use the union of their longitudinal
   breakpoints and split only boundary edges as required: 337 versus 97
   vertices does not require a fan pole or moving an existing endpoint.
   Reconcile the adjacent retained faces to avoid T-junctions. Retain the
   original evaluated triangulation outside the local support ring.
4. Construct a short cross-strip with common vertices at both stitched
   arcs. Use the retained surface tangents for its endpoint derivatives,
   then a bounded cubic/Hermite transition between them. The derivative must
   lie in each retained boundary's tangent plane; simply interpolating the
   two positions recreates a corner. Longitudinal rows must share vertices
   and normals. Do not add thickness, a raised lip, or a shell over this
   outer strip to hide a tangent mismatch.
5. Close the two short side connectors into the existing endpoint topology.
   Their support stays above 135 mm. Their exact fit is not yet measured;
   stop if a connector requires changing lower fringe, a traced tip, or a
   protected hood/lock surface. Do not collapse dissimilar endpoints by a
   broad merge-by-distance operation.
6. Check the visible outer topology before building its thickness: every
   stitched edge has exactly two incident faces, with consistent winding
   and no duplicate outer face. Then close inner skin and hidden allowance
   boundaries. The final shell must have no boundary/non-manifold edges or
   self-intersection along the new strip.

This replaces the structural cause of034's alternating ownership: one
outer face exists at the transition. The retained hidden cover is capped
below it, not allowed to compete for visible ownership.

## Thickness options and recommendation

| Option | Benefit | Strongest objection |
| --- | --- | --- |
| Keep separate 0.7/1.1 mm Solidify panels | Least change to thickness controls | Matched open outer edges still have independent normals/rim walls; easy to recreate a hairline gap or doubled edge |
| One joined outer mesh with native variable Solidify | One seam topology; editable thickness; one Head Armature | Joining normals changes derived offsets and may disturb retained surfaces beyond the desired strip |
| Join captured outer/inner skins and bridge both | Exact retained evaluated surfaces; explicit local shell closure | Thickness becomes mesh geometry rather than a live parameter |

Prefer the welded shell with captured existing outer/inner skins for the
first strict-preservation test. Only the two targets' Solidify result is
captured, not the rig. Keep the existing Head bone and full unit weights;
inverse-map all captured world coordinates into that live attachment. One
ordinary Head Armature remains sufficient. This introduces no GN system,
second rig, or static post-armature correction.

For that option, remove the matching inner faces as well as the selected
outer faces and the old 306 upper fringe rim faces. Connect the retained
inner arcs with a corresponding inner transition. Preserve the actual
0.7/1.1 mm endpoint shell separations, smoothly joining them without a thin
pinch or crossing. Cap the hidden retained-front-cover boundary and the two
side connectors. Keep existing evaluated outer/inner vertices and original
thickness rims everywhere else exactly; compare indexed correspondence,
not only bounding boxes. Source normals/UV/material data should be retained
outside the local support region where available.

The native variable-thickness option is also credible if root prioritizes
live thickness editing. Pinned RNA confirms `vertex_group` plus
`thickness_vertex_group` as the zero-influence thickness factor: one 1.1 mm
Solidify can use factor 0.7/1.1, zero weight on retained hood, unit weight on
retained fringe, and a smooth local weight transition, with offset -1.
It must pass an evaluated-surface guard away from the one-ring seam support.
Do not silently accept a repeat of the failed tuck's 1.343 mm derived X
movement. If exact retained surfaces are the stronger requirement, choose
the captured-shell option before authoring, not after a failed sweep.

Separate matched open panels are the fallback, not the preference. They
would require shared outer coordinates and tangents, and independently
hidden allowance/rim topology. Leaving each default Solidify rim at the
visible seam defeats the reset even if the base boundaries coincide.

## First-state gate and stopping rule

Use the frozen side and both three-quarter cameras. At the old034 witness
rows, the new continuous outer mesh must be the only first-hit surface
through the transition; no exposed old fringe rim, hood/front-cover sliver,
skin leak, doubled line, or transverse intersection may appear. One object
name alone is not enough: classify retained-fringe, bridge, retained-hood
and hidden-cover face regions, and check a single ordered progression.

Measure boundary tangent/normal continuity and shell thickness on both
stitched arcs; require no new raised crest, pinched valley, or envelope
overshoot. Preserve tip/eye-opening landmarks, all surfaces outside the
approved support ring, material nodes, pose and non-target evaluated
geometry. Geometric correctness is not a visual pass: render once and
reject a replacement stripe even with perfect manifold counts.

The independent rear cloth line and broad retained hood/gusset shading band
are not included. If the shared interface clears but that broader band
remains, report the distinction without a new head smoothing/resize pass.
No further geometry attempt has been run under this proposal.
