# Read-only 066-to-067 sleeve audit

The six-object change scope and evaluated manifold/finite checks are supported.
The audit flags 324 nonadjacent triangle-overlap candidates at each turned hem.
Arm-to-torso embedding remains intact; the measurements do not establish the
intended cuff-to-hand contact. No visual, module or whole-asset acceptance is
granted.

Protected input `macro_066_candidate.blend` SHA256:
`39cf9b6384b375f78b15ab65ac3a266494a20f2683a4722e68c3ddb42746b98e`.
Candidate `macro_067_candidate.blend` SHA256:
`fe1bf37cb77d044ad6023e0e019eff3e6cc2fae776d0cb25e66c94f22c10e3a4`.
The sole audit process used pinned Blender 5.2.1 LTS, build `9e2066aef7ef`,
background/factory startup/automatic scripts disabled/four threads, with a
120-second bound. It exited 0. Observation: 2026-09-05 18:24:11 UTC.
Both model hashes remained exact during and after the process. No save, render,
additional process or nested worker was used.

## Scope and validity

Exactly the intended six objects changed: each `Macro066` sleeve changed
geometry/topology, each sleeve dash object changed coordinates only, and each
stuffed arm changed transforms only. All other covered geometry, transforms,
materials and visibility match 066. Inventory remains 405 objects; all 43
protected visible mesh/curve controls and all 48 material property/node graph
fingerprints match. There are no linked libraries; both scenes remain at
frame 1. All base mesh coordinates are finite.

Each sleeve has 3,600 base vertices, 3,520 faces, exactly 160 intentional
boundary edges and no other nonmanifold/wire edges. Its 0.8 mm SOLIDIFY result
has 7,200 vertices/faces, no boundary/nonmanifold/wire edges and positive signed
volume (left 3.576036 cm³; right 3.576001 cm³). Topological closure does not
exclude self-intersection.

Saved arm transforms reconstruct the intended nominal torso start and distal
endpoint, with maximum start error 0.0000081 mm and end error 0.0000042 mm.
The distal Y is -12 mm and cross-section radii remain 9/10 mm. These nominal
endpoints remove the specified 7 mm cap allowance from the arm half-length;
they are not the extreme stuffed-surface positions.

## Turned-hem flag and witnesses

Each sleeve has 324 BVH overlap pairs between triangles that share no vertex.
All pairs involve cuff rows 34–39 and turned-hem rows 40 onward. The following
are stored pair-centroid witnesses, not exact intersection endpoints:

| Side | Triangle indices | Sleeve row span | World centroid (mm) |
| --- | --- | --- | --- |
| Left | 13294 / 13614 | 39–42 | (-68.668798, 8.729649, 60.396224) |
| Left | 13295 / 13616 | 39–42 | (-68.034828, 9.408500, 60.057640) |
| Right | 13158 / 13638 | 38–42 | (58.143713, 15.697304, 47.622707) |
| Right | 13323 / 13643 | 39–42 | (55.656713, 15.627027, 42.424604) |

These witnesses localize a cuff/return overlap issue for the next inspection.
BVH counts can include tangent/coplanar contacts; no second process was used
to classify each pair more finely. The raw report preserves eight witnesses
per sleeve and flags `left_sleeve_self_overlap_candidates` and its right-side
counterpart. No other scope/manifold/finite/endpoint check failed.

## Attachment and cloth/arm contacts

Three-ray parity finds 377 arm vertices inside the torso on each side, all in
the proximal nominal arm region (t <= 0.35), the same count as 066. Maximum
depth changes from about 14.148 mm to 14.082 mm. This is the intended torso
attachment, distinct from sleeve cloth occupying the stuffed arm.

The sleeve/arm surfaces still intersect proximally. On each side, 263 sampled
sleeve vertices lie inside the arm in shoulder rows 0–6 and 193 in shaft rows
7–33; maximum depths are 4.758 mm and 3.084 mm respectively. Corresponding
066 values were 314/342 vertices and 5.195/4.459 mm. There are 282 left/268
right sleeve-arm crossing triangle pairs, all in shoulder/shaft regions. These
counts do not by themselves distinguish desired compression from accidental
cloth breakthrough; the stored witnesses locate the actual overlap.

| Contact witness | Left world position (mm) | Depth (mm) |
| --- | --- | ---: |
| Deepest arm vertex in torso | (-11.430959, -0.610757, 66.050775) | 14.081486 |
| Deepest shoulder cloth vertex in arm | (-19.726066, -1.565743, 63.723326) | 4.758305 |
| Deepest shaft cloth vertex in arm | (-24.589773, -2.981244, 57.431959) | 3.084450 |
| Deepest proximal arm vertex in sleeve material | (-22.321301, 5.504186, 64.117253) | 0.384203 |
| Deepest distal-region arm vertex in sleeve material | (-31.138465, -11.171331, 52.808903) | 0.169196 |

Right-side witnesses are mirrored to within floating-point tolerance and are
stored separately in the JSON. There are 17 arm vertices inside the sleeve
material per side (13 proximal, 4 beyond t=0.35). These are separate from the
377 arm vertices intentionally embedded in the torso.

No evaluated sleeve-arm crossing pair or unanimously inside sleeve vertex was
found in cuff/turned-hem regions. The minimum sampled sleeve-vertex distance
to the arm is 1.788428 mm in cuff rows 34–39 and 0.478863–0.478865 mm at the
turned hem, reduced from 9.335/11.981 mm in 066. The cuff has moved closer, but
contact is not demonstrated. The closest-point location was not retained by
this single-process audit; only its distance is available. No all-triangle
minimum-clearance proof is claimed.

## Dashes and reported red patch

For each side's 342 dash base vertices, nearest distance to the sleeve center
surface is 0.472042–0.579998 mm; nearest distance to the evaluated cloth is
0.072282–0.180529 mm. Across 3,914 evaluated dash-shell samples per side
(vertices, triangle centroids and unique edge midpoints), sampled cloth
clearance is 0.039786–0.210521 mm. No dash/cloth triangle overlaps or
unanimously inside dash samples were found; four left/five right ray
classifications remain indeterminate. These results support external seating
at the sampled locations, not continuous sewn contact.

The coordinator reported a red patch visible inside the cuff in the side
image. Its pixel identity is unverified by this audit: no image was reviewed,
and no pixel-to-object correspondence was captured. The source retains the
white sleeve and red dash materials without adding a red object/material,
which is insufficient to identify that particular patch.

## Limits and evidence

No unchanged bow/head/hair contact calculations were repeated; only their
fingerprints were checked. Fingerprints omit arbitrary attributes/UVs,
custom ID properties, external image bytes and exhaustive datablock fields.
Solidify vertex-block correspondence was verified with maximum paired
midpoint error 0.00000373 mm before assigning ring regions. Ring/arm regions
are construction-coordinate diagnostics rather than independently labeled
physical attachment zones. Three-ray parity excludes samples within 1
micrometer of a surface and leaves disagreement indeterminate. Triangle
overlaps can include tangency/coplanarity; between-sample depths, complete
contact distances and visual severity remain unmeasured. The evaluated
sleeve self-overlaps further limit interpreting global inside/outside votes
as a clean solid-material classification.

`audit_067_regression.json` retains all measurements and witnesses;
`audit_067_regression.py` and `.log` preserve the sole read-only run. The
audit is complete with the turned-hem flags unresolved; no refinement or
Blender process is pending.
