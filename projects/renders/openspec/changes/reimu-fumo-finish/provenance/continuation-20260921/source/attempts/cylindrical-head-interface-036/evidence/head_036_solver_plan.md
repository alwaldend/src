# 036 first-state discretization, recorded before construction

Prerequisite result: `head_036_prerequisite.json` passes. No shape state
was constructed during that probe. The unchanged 035 support sets and
fixed model-axis chart remain mandatory.

Use Blender's bundled `mathutils.geometry.delaunay_2d_cdt` separately for
outer and inner chart domains. Keep all 438 original boundary vertices
and edges. Insert a fixed 1.5 mm Cartesian grid in the angle-height chart,
excluding points within 0.3 mm of the boundary. This is a discretization,
not paired fringe/hood rows. Reject unexpected boundary merges/splits,
nonpositive chart areas, or retained source-coordinate changes.

Use bundled NumPy and the piecewise-linear cotangent stiffness matrix L
with lumped triangle-area mass M. Boundary radii are exact Dirichlet data.
Minimize ||M^(-1/2)(L R - b)|| squared over interior radii, where b is the
integrated boundary flux from retained-face normal-derived outward radial
slopes. The two short object-to-object connector edges use the mean of
their adjacent source-edge gradients, not an invented surface normal.
This is a discrete minimum-bending solve with boundary slope data; it
does not promise exact finite-triangle seam normals. Use double-precision
least squares with column equilibration, report rank and residuals, and
reject a deficient solve instead of introducing a tuning regularizer.

Reconstruct interior points from radius and angle; use captured original
3D boundary coordinates directly, without round-tripping their floats.
Weld both chart patches into the captured retained source triangles, add
the unchanged front-cover cut allowances and two short endpoint closures,
and retain one live full-weight Head Armature with inverse-pose mapping.
No new Solidify, shaders, dependencies, or rig architecture.

Before scene-object creation, inspect actual 3D face orientation, all
retained/new junction angles, shell topology, and bounded actual local
triangle intersections. Reject reversed/degenerate cells or a seam above
90 degrees. Also report nearest inner/outer clearance and finite radius
ranges; positive chart areas alone do not establish shell validity.
Only one solved state is authorized. No grid, axis, depth, or slope sweep.

All 371 retained fringe-boundary vertices (and corresponding inner
vertices), the 26 original lower-trace witnesses, and indexed off-support
triangles remain fixed. Existing material nodes remain untouched, but
merged Generated texture coordinates and recomputed smooth normals can
change target appearance; exact off-support shading is not claimed.
