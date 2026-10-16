# Sleeve 088 helper proposal

Status: unexecuted geometry proposal; coordinator must freeze the helper and
plan before any native call. Only Python syntax parsing was performed. No
Blender process, simulation, object operation, model save, parameter trial,
goal operation, or Git mutation was performed by this helper task.

The source is working 086 with arm/body/skirt receiver geometry unchanged
from 084. The sole model writer supplies evaluated world-space arrays and
matching BVHs from that exact candidate. API:

```python
data = sleeve_geometry(sign, exact_087_json_bytes, {
    'arm': {'points': world_points, 'faces': triangles, 'tree': arm_bvh},
    'body': {'points': world_points, 'faces': triangles, 'tree': body_bvh},
    'skirt': {'points': world_points, 'faces': triangles, 'tree': skirt_bvh},
})
```

The helper has no scene or filesystem operations. It returns sleeve vertices,
triangle faces, outer vertex count, outer chart coordinates, 342 dash vertices
and 152 dash quads, and diagnostics. Prepare both sides successfully before
any coordinator-owned object replacement. Returned coordinates are world
space; convert through each receiving object's inverse world matrix.

The 087 JSON is guarded by SHA256
`91bda2bf80d7452a5b0ecffdec3b4b138bedd3014d5b6e1e3466c924fea6cb03`.
Its accompanying interpretation MD is
`92abc03ad845a729876d83814f47ce4ed63f363460c8597de54ebed9828353d5`.
The 160-point measured root at station center |X|=35 mm is preserved exactly.
The section plane is perpendicular to the actual arm axis; it is not a
constant-world-X plane. Its measured nonplanarity is allowed only up to 10 nm.

The outer garment consists of front and back height graphs over exactly the
same arm-local 2D domain. A is normalize((sign*47,-12,-26) mm), V is world Z
projected perpendicular to A, and N is the perpendicular direction pointing
toward world -Y. Each original root arc supplies its own monotone-V boundary
samples. Separate constrained Delaunay triangulations retain those samples
while sharing the same polygonal upper/lower seams and cuff outline. Interior
sampling is fixed at 32 longitudinal by 48 transverse divisions.

Cuff endpoints are (sign*69,-4,66) and (sign*49,-3,34) mm. Its shared domain
curve adds A*19 mm*sin(pi*t) to their linear interpolation; the front/back
height lobes give 11/10 mm in world Y, hence 21 mm maximum cuff depth. This
places the nominal front cuff span near |X|=74 mm. The upper/lower domain seams
bow outward by 1 mm*sin(pi*s); panel interiors have a 0.7 mm height billow.
These are initial construction dimensions, not measured visual acceptance.

Before fitting, sample rays and full segment projection curtains reject any
shared seam crossing an actual receiver silhouette. The cuffs and root stay
fixed. Interior vertices move only along N to the appropriate extremal actual
receiver hit plus a local tangent-plane allowance of 0.6+0.5 mm. There is no
radial or independent XYZ projection and no cloth simulation. The common 2D
domain stays fixed when height fitting changes the garment depth.

The outer graphs weld along their longitudinal seams. A 0.5 mm inward
averaged-normal displacement constructs the liner; bridges connect it only
around root/cuff boundary loops, leaving the garment apertures open. Guards
reject inverted liner triangles, inconsistent welds, non-manifold walls,
nonpositive volume, self-overlap, receiver intersections, sampled clearance
below 0.6 mm, and reversed graph ordering at interior mesh vertices. Red dashes
are sampled through chart coordinates onto the final triangulated panels.
Nine back and ten front marks avoid both shared seams; their quads follow
the sampled outward normals, and dash/wall plus dash/receiver crossings fail.

The numerical gates have explicit limits: ray offsets give local tangent-plane
spacing; nearest-distance minima cover vertices, edges' midpoints, and triangle
centers, not a continuous clearance bound. The auxiliary adjacent-fold check
shrinks each triangle by 1e-3 toward its centroid, omitting narrow boundary
strips and depending on float32 BVH precision. A zero resolved displacement
fails as an ambiguous guard. The 1e-3 fraction follows a static precision
review: a 1e-5 fraction can round away at the source coordinate scale.
Averaged-normal displacement
does not prove constant normal wall thickness. Receiver arrays/BVHs must match
and their validity is the caller's responsibility. Model-wide evaluated
geometry checks, clean reopen, fixed-view renders, and the independent visual
gate remain mandatory; a successful helper call is not acceptance.

The Blender/reference skills influenced the choice of sewn panel construction,
preserved attachment evidence, immutable inputs, and separate technical and
visual gates. No model quality verdict is implied by this proposal.
