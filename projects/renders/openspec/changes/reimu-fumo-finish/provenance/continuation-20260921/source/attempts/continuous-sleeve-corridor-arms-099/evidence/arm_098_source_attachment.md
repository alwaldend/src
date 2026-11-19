# Source arm attachment in natural coordinates

**An arm-only center route is geometrically supported without changing the
cloth, once the old body-contact locus may move.** The coordinator explicitly
permits that change while keeping the body fixed and requiring a continuous
buried proximal arm. Exact old-interface triangle preservation is not a user
acceptance requirement. Its conflict with sleeve clearance is documented below
to prevent reinstating an unsupported freeze, not to require garment changes.
This report proposes no arm shape or numerical fit settings.

Pure-data diagnosis of the retained094 assembly. No Blender/native execution,
model/goal/Git writes, candidate, trial-shape sweep or altered distal target.

## Corridor-center feasibility

Two bounded geometric witnesses were measured against every actual evaluated
own-sleeve wall triangle. They are not candidate arm surfaces:

| Center-path witness | Sleeve-wall intersections | Minimum continuous path-to-wall distance |
|---|---:|---:|
| Straight segment, original proximal pole to useful095 distal pole | 0, both sides | **4.347505 mm**, both sides |
| Actual retained source row0–40 cross-section centers, connected in row order | 0, both sides | **7.608069 mm**, both sides |

The endpoint-determined segment starts at **(±10.022, 1.526, 72.307) mm**,
inside the actual body, and ends at **(±56.978, −5.026, 53.993) mm**, the
actual useful095 distal pole, outside the body. It crosses the body exactly
once, at approximately **(±26.113, −0.719, 66.031) mm**, with no ray-edge
ambiguity. Thus it provides a continuous center connection from a genuinely
buried proximal region into the sleeve.

Its root-plane entry is **(±22.276, −0.184, 67.528) mm**, inside the actual
80-point root aperture. It enters the aperture before leaving the body and
never crosses the sleeve wall afterward. The distal target's garment-axis
station is **36.205 mm** from the root, before every cuff-lip point's station
**39.066–46.263 mm**. This rules out a path that simply bypasses the sleeve
outside its root or exits beyond the cuff.

The straight path's tightest wall distance occurs near
**(±21.312, −0.049, 67.904) mm**, near the root-side transition. The existing
sleeve-center path's tightest clearance is also near row0. Distances use the
complete finite line segments against triangle faces/edges, not only sampled
center vertices. Actual-source center arrays and closest-point witnesses are
in the separate corridor JSON.

Positive margin admits a sufficiently thin connected arm volume along a
route; it does **not** prove that a chosen constant-radius/ellipsoid profile,
the original 9/10 mm transverse radii, or a visually convincing stuffed girth
fits. No width, centerline bend or contraction factor was selected. The
continuous sleeve-centered loft and shortened/reoriented ellipsoid remain
different construction hypotheses to evaluate, not tested shapes here.

## Natural frame and measured body attachment

Coordinates are obtained with the **exact inverse actual arm object matrix**.
Let `(a,b,u)` be local unit-sphere coordinates: `u=−1` is the proximal pole and
`u=+1` the distal pole. Physical transverse coordinates are approximately
`(9a,10b)` mm; axial distance from arm center is `s=34.518184u` mm.
The centers are `(±39.500, −6.000, 56.000)` mm. The distal axial direction is
approximately `(±0.853981, −0.218038, −0.472415)`. Matrices, inverse rows and
all 1,106 natural vertex coordinates per arm are preserved in JSON.

| Complete source triangle classification | Left | Right |
|---|---:|---:|
| Body-surface-crossing triangles | 104 | 104 |
| Wholly inside body | 700 | 700 |
| Wholly outside body | 1,404 | 1,404 |
| Total source arm triangles | 2,208 | 2,208 |
| Actual arm/body crossing pairs | 372 | 344 |

The independent scalar arm/body contacts reproduce native source counts;
all 372/344 pairs resolve as noncoplanar crossings. After removing their arm
triangles, each arm has exactly two edge-connected triangle components.
Strict three-ray parity classifies the 700-triangle component inside and the
1,404-triangle component outside the actual closed/self0 body. All witness
votes agree, with no boundary or edge ambiguity. A connected surface disjoint
from the body cannot change body containment without crossing it.

| Attachment bound | Natural `u`, bilateral approximate | Axial `s`, mm |
|---|---|---|
| Actual body-intersection segment locus | −0.662878 to −0.240594 | −22.881 to −8.305 |
| All vertices of the 104 complete crossing triangles | −0.707107 to −0.130526 | −24.408 to −4.506 |
| Vertices of wholly buried triangles | −1.000000 to −0.258819 | −34.518 to −8.934 |
| Vertices of wholly exterior triangles | −0.608761 to +1.000000 | −21.013 to +34.518 |

These overlapping axial ranges show why body attachment is not one clean
latitude of the oblique source sphere. The body-contact contour depends on
the transverse coordinates as well as `u`. “Inside” is a solid-containment
classification, not proof of sewing or a requirement to preserve every buried
triangle unchanged.

## Attachment triangles also cross the sleeve

| Source arm/sleeve relationship | Left | Right |
|---|---:|---:|
| Arm triangles crossing own sleeve | 102 | 97 |
| Also body-surface-crossing | 39 | 36 |
| Body-crossing triangles with a body-exterior sleeve interval | **19** | **20** |
| Wholly buried sleeve-crossing triangles | 27 | 24 |
| Wholly exterior sleeve-crossing triangles | 36 | 37 |

Each recorded arm/sleeve contact segment was split at all actual body ray
events, then each open interval classified by strict three-ray parity. There
were **no unresolved intervals or segment-edge ambiguities**. This distinguishes
the exposed portion of a body-crossing triangle from its buried portion.
Here “exposed” means **outside the body**, not proven visible in a render.

| Classified sleeve-contact subintervals | Left | Right |
|---|---:|---:|
| On wholly buried arm triangles, inside body | 95 | 90 |
| On body-crossing triangles, inside body | 74 | 65 |
| On body-crossing triangles, outside body | **58** | **59** |
| On wholly exterior triangles, outside body | 87 | 86 |

The 19/20 conflicting complete triangles span `u≈−0.382684..−0.130526`
(`s≈−13.210..−4.506 mm`). Their body-exterior sleeve-contact portions occupy
`u≈−0.365502..−0.130526` and local transverse `b≈−0.973302..−0.098455`,
the lower portion of the natural arm section. World contact bounds are:

| Body-exterior sleeve intervals on body-crossing arm triangles | Left, mm | Right, mm |
|---|---|---|
| X | −30.657 to −26.220 | 26.220 to 30.745 |
| Y | −12.118 to −0.527 | −12.274 to −0.527 |
| Z | 50.678 to 61.092 | 50.678 to 61.092 |
| Retained sleeve cells | rows 6–14 | rows 6–14 |

The actual sleeve root is row0 and cuff lip row40. These body-exterior
intervals cannot be relabeled as a contact at the root seam. Preserving their
complete source arm triangles while leaving the sleeve fixed necessarily
preserves those intersections, regardless of changes farther down the arm.

Preserving the 104 crossing triangles fixes 104 vertices and 208 mesh edges
per arm. The fixed-vertex closure contains exactly those 104 complete
triangles, no additional complete triangles. JSON also records sleeve-contact
endpoint witnesses on the fixed edges, so the conclusion does not depend on
ignoring fixed primitives adjacent to the triangle set.

## Consequence for an axial fit

Switching from world-X weights to natural axial sections addresses the failed
095 transition representation, but **does not remove this fixed-geometry
conflict**. Even an axial prefix merely large enough to include every actual
body-contact point reaches body-exterior sleeve intervals: 31 left / 32 right.
That is interval-inclusion evidence, not a suggested freeze threshold.

The earliest source sleeve contact is at `u≈−0.666124`; the earliest actual
body-surface intersection is `u≈−0.662878`. A deeper proximal region exists
inside the body without these measured surface contacts, but retaining only
such an anchor would allow the existing body-interface triangles/contour to
change. An arm-only rebuild with a continuous buried connection is therefore
**not disproven in general**; it is not the same operation as preserving all
104 interface triangles exactly, and its resulting connection and clearances
remain untested.

The coordinator has now chosen preservation of a meaningful continuous buried
connection, **not** the old complete interface triangles. The affected
interface triangles can therefore participate in an arm-only refit with the
body and cloth unchanged. The positive corridor witnesses above support that
direction; the earlier fixed-patch conflict does not establish a need to
change the garment. The useful095 distal movement remains the endpoint
evidence, not a newly tested complete fit here.

All other **exported095** source arm receiver contacts are zero in the bound
native preflight: opposite arm/sleeve, dashes, red skirt, ivory hem, collars,
head/cover, side cloths and bow meshes. This coverage does not include every
nonexported object, and says nothing about contacts after a future refit.

## Evidence

JSON includes exact triangle/vertex sets, every source arm/body segment,
all split arm/sleeve intervals with parity evidence, source natural coordinates,
fixed-vertex closure, and the axial-prefix conflict witnesses. The body geometry
is hash-identical to the previously guarded closed/self0 receiver. Finite
ray events merge within 100 nm and intervals at most 200 nm would remain
unresolved; neither occurred ambiguously in this result. Source pair counts
and triangle indexing are validated against the bound prior evidence.

| Artifact | SHA-256 |
|---|---|
| `arm_098_source_attachment.py` | `4b3a878a126ec37149ec73e365ac4fffde941ecb71f5090c2581d5bded702ec7` |
| `arm_098_source_attachment.json` | `6d88701d6bd20faa4a2ac622fb58d98871f718dfb9e809c67ef60d912335ce1b` |
| `arm_098_source_attachment_corridor.py` | `ae75b5d6a44077a67a4ea97974f88fd1f4caf56a1000856ba1781056a080ae71` |
| `arm_098_source_attachment_corridor.json` | `500dba0b129c1324e8348dbce1e45952a40ce19bed87465c83567b5b95648f21` |
| `sleeve_095_actual_assembly.json` | `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf` |
| `arm_095_collision_diagnosis.json` | `269759de8e6a8ba1961d1748afb23ff14dc5fcc94c0cd3979ad3e5ad5a1e0aa3` |
| `arm_097_fit_design.md` | `d0f9ac500def3fd7116f2c35f1bd3edc565aa19f961199cb8b83107c8056cab6` |

Fully read scalar helpers were AST-extracted under hash guards, never their
native builders. Standard-library computation exited 0, including the added
fixed-vertex closure check; all input hashes stayed unchanged.
