# Foot033 independent saved-file technical audit

Verdict: **PASS with a smooth-normal qualification**.

The pinned Blender 5.2.1 LTS audit reopened
`head_032_candidate.blend` and `foot_033_candidate.blend` independently. It
did not import the Foot033 helper or writer, save either file, render, or
evaluate animation.

## Artifact binding

- Source SHA-256:
  `6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`
- Candidate SHA-256:
  `98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48`
- Writer receipt SHA-256:
  `9ae2e7ec2440502f2bc627865a501fcf3f3d1892ffdad5b4d7ad4b19611b9eb9`
- Independent audit script SHA-256:
  `48b1a08f520c549f8476ffeb90004a707fb1ef8388a2482b0d2ca4a2872dddf7`
- Successful audit invocation:
  `0b9bcbbc-0ede-4498-a5c8-bfbc548ebcd5`

Both source and candidate hashes were unchanged after the audit.

## Scope and controls

The only added objects are the two reviewed Foot033 replacements. The only
added meshes are their exact named mesh datablocks, and the only added
collection is `Foot033 conformal two-material pods`. No object, mesh,
collection, or material was removed; no material, armature, camera, or light
was added.

All 84 receipt-bounded non-target mesh/curve controls independently match
Head032 in evaluated geometry, material-slot names, transforms, parenting,
visibility, and recorded modifier settings. The rig pose also matches
exactly. The visible set changes only from the two donor pods to the two
replacement pods: both source and candidate contain 86 visible mesh/curve
objects under that definition.

Both original donor pods remain exact in their audited base-mesh, UV,
transform, weight, parent, and modifier records. Each is hidden from
render and the active view layer; neither has the global `hide_viewport` flag
set.

## Surface partition

| Check | Left | Right |
| --- | ---: | ---: |
| Original indexed vertices | 1,962 | 1,962 |
| Maximum original evaluated vertex delta | 0 m | 0 m |
| Added seam vertices | 144 | 144 |
| Maximum seam vertex distance from donor edge | 7.68e-9 m | 7.45e-9 m |
| Candidate triangles | 4,208 | 4,208 |
| Unmatched candidate triangles | 0 | 0 |
| Maximum donor-plane distance | 3.61e-9 m | 3.01e-9 m |
| Maximum donor area-coverage relative error | 7.92e-6 | 8.73e-6 |
| Maximum geometric face-normal angle | 0.0402 deg | 0.0396 deg |
| Maximum UV interpolation error | 1.33e-7 | 1.49e-7 |
| Maximum normalized material-side residual | 2.50e-7 | 2.95e-7 |

Every output triangle is contained in and area-accounted to an original
evaluated donor triangle within the listed float32 tolerances. The first
1,962 evaluated vertices are exactly indexed to the donor. Each mesh is one
closed genus-zero component with 2,106 vertices, 6,312 edges, 4,208 faces,
Euler characteristic 2, and no edge with other than two incident faces.

The material seam is one closed loop with 144 shared vertices and 144 edges;
every seam vertex has degree two. Each foot has 2,894 black faces and 1,314
cream faces. No bottom-20-percent vertex is cream. The normalized field-side
residual is below `3e-7`; it is numerical interpolation error, not a visible
surface offset.

## UV, rig, and rest representation

Each replacement retains one `UVMap`, active at index 0 and active for render.
The per-corner values match barycentric interpolation from the containing
donor triangle to at most `1.49e-7`.

Each object is parented to `ReimuFumoRig`, has identity world matrix, exactly
one full-weight `Leg_L` or `Leg_R` group, and zero wrong or non-unit vertex
weights. Its sole Armature modifier matches the corresponding donor's name,
rig target, visibility flags, vertex-group/envelope settings, volume setting,
and remaining audited scalar settings.

The replacement stores inverse-Leg-pose rest coordinates. Raw local values
therefore differ by about `0.073 m` from the donor object's separately framed
local coordinates; that number is diagnostic, not displacement. Against the
inverse-pose source surface, indexed rest coordinates match exactly, new rest
vertices lie on donor edges within `7.68e-9 m`, and the saved modifier's
frame-1 forward-pose round trip has zero measured error.

## Smooth-normal qualification

The surface and geometric face normals are preserved within float32 tolerance,
but smooth vertex normals are recomputed when the mesh is triangulated and
receives new seam vertices.

- Original indexed vertex normals differ from the donor by at most `0.0499°`
  left and `0.0517°` right; p95 is about `0.034°`.
- New seam normals differ from simple endpoint-normal interpolation by median
  `1.38°/1.37°`, p95 `4.13°/3.06°`, and maximum `6.19°/4.10°` left/right.
- Every face is smooth, neither replacement has custom normals, and there are
  zero sharp edges, including along the material seam. The measured seam
  normal difference is therefore not a hard normal split, though it can cause
  a small local shading change and remains subject to rendered review.

## Limits

This audit intentionally did not perform a broad hidden-object geometry scan,
animation evaluation, material-node/light/world comparison, save, render, or
visual acceptance. The 84-object comparison is the bounded visible-control
set from the receipt; writer evidence owns appearance preservation. Technical
cleanliness does not establish reference fidelity.
