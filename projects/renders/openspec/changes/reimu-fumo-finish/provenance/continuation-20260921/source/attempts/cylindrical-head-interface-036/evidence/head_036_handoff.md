# 036: cylindrical patches solve; closure contacts stop creation

One pre-recorded state was solved. The cylindrical chart removes the
previous orientation/fold failure, but the local shell contact gate stops
before object creation. No candidate was saved or rendered. No parameter,
grid, axis, boundary, or solver retry occurred.

Helper: `head_036_draft.py`, entry `build_head_036()`, exact two `TARGETS`:
`Hair028 traced padded fringe` and `Hair028 crown and back hood`.
The helper and first-state bytes are frozen for root review, not silently
corrected after this result.

## Prerequisites and construction

`head_036_prerequisite.json` confirms the captured 438-point inner boundary
is noncrossing in the fixed cylindrical chart. The root-settled outer chart
was not retested; its original world coordinates were verified exactly.
Minimum expected signed radial normal is 0.4337 at fringe retained faces
and 0.4943 at hood retained faces. No near-singular samples occurred.
All 26 original traced lower-contour witnesses lie on the 371-vertex
retained fringe open boundary, with exact source indices recorded.

`head_036_solver_plan.md` was recorded and hash-bound before execution:
independent constrained chart triangulations, fixed 1.5 mm grid and 0.3 mm
boundary clearance, cotangent minimum-bending least squares, exact boundary
radii, and compatible source-normal boundary flux. Original 3D boundary
positions were used directly. The 035 support sets were unchanged.

| Measurement | Outer | Inner |
| --- | --- | --- |
| Interior scalar unknowns / solver rank | 454 / 454 | 448 / 448 |
| Triangles | 1,344 | 1,332 |
| Equilibrated condition number | 7.532 | 7.330 |
| Radius range (mm) | 22.3033–65.8132 | 21.8394–64.7169 |
| Largest retained/chart junction angle | 33.5412° | 46.3451° |

All 2,676 chart triangles have positive actual radial orientation,
minimum signed radial cosine 0.5021. There are no degenerate chart faces,
inconsistent winding edges, or nonmanifold edges in the assembled arrays.
Across all 872 retained/new chart junctions, angles are median 2.4444°,
p95 10.9423°, maximum 46.3451°; none exceeds 90°. These are technical
measurements, not a visual acceptance claim.

## Why the first state stopped

The bounded actual intersection test compared all 2,912 new chart and
closure triangles against the captured target shell. Its 1,420 broad-phase
pairs yielded 307 nonshared segment-contact records:

| Contact category | Pairs |
| --- | --- |
| Hidden-cover closure vs retained fringe inner skin | 296 |
| Hidden-cover closure vs new inner chart | 8 |
| Root-endpoint closure vs retained hood | 3 |

None involves the new outer chart. No coplanar-overlap contact was reported.
Contact positions span Z 134.2440–184.7158 mm. Full points, array face IDs,
and retained source-triangle correspondence are in `head_036_dryrun.json`.
The test allows shared edge/vertex contact within 0.2 micrometers and uses
a 1e-12 m² projected coplanar-area threshold; these are bounded numerical
tests, not a whole-scene collision audit.

Examples: inner chart face 75691 crosses closure face 75932 near
(-57.142, -29.950, 135.903) mm. Retained fringe triangle 22119 crosses
closure face 75916 near (-55.946, -32.568, 135.979) mm. Retained hood
triangle 15300 crosses root closure face 75938 near
(-57.422, -30.153, 134.244) mm. The closure labels identify construction
roles, not verified camera invisibility. No first-hit visibility audit ran.

## Clearance classification defect, not a shape retry

The nearest-inner query also reported two zero-distance samples. Its
selection used region 4, which includes retained Solidify rims as well as
inner skin. Thus those two records do not establish collapsed thickness.
Using the unchanged all-quad fringe capture ordering and exact removed
polygon lists, retained array faces 39726/39895 correspond to source
triangles 43794/44115, source polygons 21897/22057. Both are rim polygons;
the inner-skin range ends at polygon 21887. The queried outer samples are
the retained root vertices 2914/3612, where rim contact is expected.
This correspondence is an array-order derivation, not another Blender run.

The raw sampled nearest-distance distribution is median 0.9836 mm,
p95 1.1135 mm, max 1.2080 mm. Its reported zero minimum must remain
qualified by the rim-classification issue. No corrected query was run and
no zero-clearance gate was silently removed. The 307 independent closure
contacts remain regardless of this classification defect.

## Preservation and next authority

There were zero created objects and neither target was hidden. The source
and helper hashes remain unchanged, and non-target geometry comparison
passed after the unsaved array rejection. The arrays preserve captured
off-support geometry, but no evaluated new-object preservation or motion
claim is available. Material nodes were untouched; a future merged object
would change Generated-coordinate mapping and potentially smooth normals.

Root owns whether the closure/retained-layer interface needs a different
construction or a narrower acceptance statement. This result provides no
authority to weaken the intersection gate, edit the immutable helper, or
run a second state. Foot033 remains retained.

## Evidence binding

- Input `head_032_candidate.blend`:
  `6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`
- `head_036_draft.py`:
  `8b2dae5451e643a59ae49e81f3b4d76249e7e0c32be161834c3e989b1cd46fb9`
- `head_036_prerequisite.json`:
  `3e18644cb2047579be1f38fbe509b282c9d3980b2ba454ede37d2dfc718fc0b7`
- `head_036_solver_plan.md`:
  `c5d88cd6f72f0506d1c4eca48c615fbf3ffec2e2ce9c88b3806d54c18b3d4b8c`
- `head_036_dryrun.json`:
  `257b0cbab00a0f04099ef10dc8111bb1dd396522c932545ddfcffe512a435979`

Runtime: pinned Blender 5.2.1 LTS, build `9e2066aef7ef`. Work stayed
sequential and no delegation occurred, as authorized.
