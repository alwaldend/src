# Head interface: chart-based topology proposal

Status: plan only, following rejection of the fixed ruled-strip family.
No triangulation, depth solve, object construction, save, or render was run
for this proposal. Root owns the next scope decision and any candidate.

## Decision

Do not triangulate the exact existing boundary as a single-valued Y(X,Z)
surface. It has two projected crossings with contradictory depth values.
Instead, authorize a small structural expansion at the two side roots,
removing those legacy overlap endpoints from the fixed boundary. Then test
the changed boundary once before constructing a surface. This is a topology
reset, not another strip-width, depth, or endpoint-tangent adjustment.

The proposed outer surface would use one constrained triangulation in XZ
and a smooth scalar Y field. Interior connectivity would not depend on the
old 432 paired fringe/hood knots. The first implementation remains
conditional on a simple, noncrossing expanded boundary with consistent
depth constraints; that changed boundary has not yet been selected or tested.

## Bound evidence

The one authorized read-only chart probe is
`head_034_chart_probe.py`, with full points and results in
`head_034_chart_probe.json`. It inspected frozen
`head_032_candidate.blend`, SHA-256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`,
using pinned Blender 5.2.1 LTS, build `9e2066aef7ef`.
It imported functions from the unchanged shared helper, SHA-256
`c9976b58dc310001b4478ddb399b6691bd801eb3a725b12358c6422000cb1c69`,
but did not call its builder or create an object.

The closed boundary has 434 vertices: the 337-vertex fringe inner arc,
the reversed 97-vertex hood/gusset arc, and two endpoint connectors.
Its two nonadjacent projected segment crossings are:

| Root | Fringe source edge | Hood source edge | X, Z (mm) | Incompatible Y values (mm) |
| --- | --- | --- | --- | --- |
| Left | 3620–8566 | 840–2238 | −57.9220, 138.8175 | −30.4831 / −29.6114 |
| Right | 11284–4318 | 2706–794 | 57.9228, 138.7996 | −30.4846 / −29.6119 |

The depth conflicts are 0.8716 and 0.8727 mm. There were no other reported
nonadjacent segment contacts. These cannot be repaired by silently averaging
Y, snapping endpoints, dropping a constraint, or feeding the crossing
boundary into a triangulator as though it were a simple polygon.

The chart is otherwise locally usable along the inspected retained sides:
every boundary-adjacent triangle has negative normal-Y, with minimum
absolute normal-Y 0.2162 for fringe and 0.2167 for hood; none are below
0.02. Thus the measured obstacle is the selected overlap boundary at the
roots, not a vertical surface singularity there. This does not prove that
every possible three-dimensional chart fails. A more complicated chart
merely to retain arbitrary old overlap endpoints has no demonstrated
reference benefit.

The left fringe cut begins with vertex 3620 at Z 135.5399 mm, jumps to
8566 at Z 140.1991 mm, then returns through 3619 at Z 139.5065 mm and
8565 at Z 138.8452 mm. The right root has the corresponding reversal.
This locally folded projected cut explains why preserving every old
intermediate boundary vertex is an inappropriate construction constraint.

## Smallest proposed boundary reset

Keep the two source objects as the only replacement targets:

- `Hair028 traced padded fringe`
- `Hair028 crown and back hood`

Use these exact crossing-support faces as expansion seeds, followed by one
bounded edge-sharing neighbor ring within the local root-transition patch:

| Root | Fringe retained polygon / evaluated triangle | Hood retained polygon / evaluated triangle |
| --- | --- | --- |
| Left | 68 / 137 | 7836 / 15671 |
| Right | 10941 / 21883 | 7743 / 15485 |

These seed sets identify where the new support must change; they are not
an assertion that the expanded perimeter is already valid. Build the new
perimeter from the union of the prior editable band and these bounded
local faces, then check its actual ordered points. The former fringe
endpoints 3620/4318 and their neighboring arbitrary transition vertices
may become interior degrees of freedom. Do not retain them as conflicting
Dirichlet constraints merely because the prior cut used them.

Root must explicitly authorize this local support change, including any
neighboring face crossing the old 135 mm cutoff. Preserve the actual fringe
tips, eye-opening edges, and their surrounding functional cloth. Do not
expand across those landmarks, through the neck, or into unrelated rear
cloth. If one bounded ring still leaves a crossing or inconsistent depth,
stop with its exact boundary evidence rather than growing rings repeatedly.

## Proposed construction after a valid boundary exists

1. Form a single planar straight-line boundary graph in XZ with unique
   boundary Y values. Keep required edge and landmark constraints, not a
   cross-strip pairing. Confirm the graph encloses the intended brown
   outer transition and does not remove the retained forehead cover.
2. Constrained-triangulate that domain. Use positive, nondegenerate XZ
   triangle area and bounded aspect ratio as geometric prerequisites;
   do not insert the former 432 correspondence rows.
3. Solve a smooth scalar Y field with exact boundary depth and a
   minimum-bending objective. Use retained-face normal-derived boundary
   slopes where well-conditioned. Conflicting corner value/gradient
   constraints are a reason to revise the construction boundary, not to
   force a fold or hide it with thickness. No solver or library selection
   is needed for this plan-only stage.
4. Weld this outer patch into the captured off-support shell. Retain the
   selected captured outer/inner-shell representation, existing material
   assignment, and one full-weight live Head Armature. Preserve original
   off-support indexed geometry and UV data where present. The inner
   transition and closures need their own orientation and clearance checks;
   a valid outer depth graph does not prove a valid thick shell.

A consistently oriented triangulation with scalar depth avoids the prior
ruled mapping's mixed XZ Jacobians and competing values at the same chart
point. It does not guarantee pleasing curvature, matching seam normals,
adequate thickness, or absence of intersections with other retained parts.
Those remain measured first-state gates, not implied acceptance.

## Preservation and falsifiers

Preserve head core and all non-target objects, neck/contact geometry,
eyes/mouth, bow, locks, body, and feet. Keep actual fringe tips and eye
openings fixed. Preserve all target geometry outside the newly declared
local support. Maintain the brown front cover needed behind the forehead;
eliminating overlap at the visible interface is not permission to expose
the pale core through existing openings. Reuse existing materials and Head
binding without new rig or geometry-node architecture.

Before construction, reject a new boundary that still crosses, has two
depths at one boundary point, branches ambiguously, or intrudes on protected
tips/openings. After one authorized construction, report triangle signed
areas, boundary-normal angle distributions, local shell intersections,
thickness clearance, and camera first-hit face-region ownership. Any
reversed/degenerate cell or seam above 90 degrees is a definite rejection,
not a reason to sweep a depth parameter. A welded projected ownership
change alone is not evidence of overlapping sheets; inspect actual depths,
normals, and topology as the previous mirror-row diagnostic established.

The controlling front landmarks still come from
`references/canonical_front_25cm.png`; physical front/side photographs
inform cloth construction with their pose/variant limitations. Final
side/three-quarter images must demonstrate that the ridge is removed
without replacing it with a crease or excessive crown curvature. No visual
or rig acceptance is asserted by this proposal.

## Relationship to the rejected first state

`head_034c_handoff.md` and `head_034c_diagnostic.json` retain the causal
failure: all 1,728 endpoint direction checks passed, but 456 of 6,896 strip
cells were reversed or mixed, including 343 wholly above 145 mm. Six
junctions exceeded 90 degrees, with a maximum of 167.181 degrees. The
precreation gate stopped correctly; no object existed for an evaluated
intersection audit. This plan changes surface topology and boundary
constraints instead of retrying those settled derivative checks.

This task remained sequential because the authorized boundary check
directly determined the plan, while root independently handled projection
diagnostics and canonical rejection. No candidate or canonical state was
written here.
