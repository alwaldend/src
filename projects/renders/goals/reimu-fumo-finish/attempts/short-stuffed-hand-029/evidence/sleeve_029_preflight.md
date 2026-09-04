# Local sleeve-mouth and stuffed-hand preflight

Read-only source: `head_028_candidate.blend`, SHA256
`c4ab72a53eb12e64f7f5d2bb216ea1a1734f0bb43cf8e19393f532624aa671b6`.
Pinned Blender 5.2.1 LTS, build `9e2066aef7ef`; source hash checked before and
after. `sleeve_029_probe.py` wrote `sleeve_029_probe.json`, never saved a blend.
The JSON contains all 44 exact sleeve object names, evaluated bounds,
materials, modifier stacks, bone/rest-pose data and front-ray witnesses.
No helper, model, goal or canonical record was edited.

## Finding

The hands are not absent or hidden. The two visible, render-enabled meshes
named `Sleeve44P L asymmetrically seated stuffed arm insert` and
`Sleeve44P R asymmetrically seated stuffed arm insert` are closed skin-fabric
inserts recessed inside the cuff. No separate hand/palm-named object was
found. The fixed side/three-quarter images show their small shaded ends
inside disproportionately deep openings; the front view largely hides them.

The canonical front controls the current down/out arm angle, which is already
approximately right. Physical front has a different pose and does not justify
rotating the arms. Physical side supports a pale rounded stuffed end nested
in a soft flared sleeve, not a large fist or an empty rigid trumpet. Keep the
existing shoulder-to-cuff axis and the front X/Z sleeve silhouette.

All numeric scene witnesses below are millimeters; Wh = 116.5 mm. The cuff plane
is a reproducible diagnostic plane through the combined cuff-edge AABB center,
normal to shoulder-root-center -> cuff-center. It is not a fitted sewing plane.

| Witness | Left | Right |
| --- | --- | --- |
| Shoulder-root center XYZ | -27.700,-2.997,76.056 | 27.700,-2.997,76.073 |
| Cuff center XYZ | -60.620,-11.931,46.820 | 60.524,-11.931,46.715 |
| Root-to-cuff length | 44.925 | 44.935 |
| Insert signed axial extent from cuff | -18.111 to -6.077 | -18.098 to -6.138 |
| Cuff-edge signed axial extent | -4.423 to +3.021 | -4.421 to +2.925 |
| Sleeve-only front-ray occlusion | 60 of 65 hits | 61 of 65 hits |

Thus the most distal insert surface is still about 6.1 mm behind the mean cuff
plane and about 9.1 mm behind its most outward lip. The cuff's combined Y span
is 36.60 mm, while each insert's Y span is only 16.65–16.67 mm. This geometric
recession and excess opening depth explain the missing-hand read without
changing visibility or making assumptions from a dark render alone.

| Evaluated insert bounds | Left | Right |
| --- | --- | --- |
| X | -61.623 to -38.978 | 39.057 to 61.585 |
| Y | -17.923 to -1.257 | -18.112 to -1.459 |
| Z | 41.022 to 60.180 | 41.047 to 60.119 |

## One small construction correction

Retain the two sewn white panels and existing arm angle; reshape only the
short stuffed insert and distal mouth. Do not rebuild or repose the whole arm.

1. Reshape each insert into a short closed, gently flattened stuffed mitten
   capsule with no fingers or separate fist. Preserve its concealed proximal
   end near axial -18 mm. Extend the distal rounded cap from about -6.1 mm to
   axial +4.5..+5.5 mm, only 1.5–2.5 mm beyond the current outermost lip. This
   is a roughly 10.6–11.6 mm extension of the hidden insert, not an equivalent
   exposed-hand length. Keep the cross-depth around 15–17 mm and preserve the
   present cross-sleeve width; the target is a small visible curved cap.
2. Soften the opening around that support by reducing its excess depth,
   without altering front X/Z. In the local axial coordinate a measured from
   the cuff center, a bounded first field is:
   w=smoothstep(clamp((a+.015)/.012,0,1));
   Ynew=Y-.34*w*(Y-Ycuff). This leaves the proximal 30 mm or so unchanged and
   reaches full depth compression 3 mm before the mean cuff plane. It predicts
   an outer mouth depth near 24–25 mm instead of 36.6 mm. Fit the final inner rim
   to the actual new capsule: allow a small crescent of opening and two
   unequal existing pinches, not an airtight circular gasket. Keep the
   existing 0.78 mm cloth thickness and softly folded rim.
3. Transfer the identical distal field to cuff edges, pinched folds, sewn
   joins and red running stitches. Keep all shoulder pleats and the gathered
   shoulder seam unchanged. This avoids a moving panel beneath floating trim.

These dimensions are bounded modeling hypotheses informed by the scene and
reference construction, not exact dimensions measured from the photographs.
Actual first side/three-quarter pixels decide whether the exposed cap is
large enough; do not enlarge the whole hand merely to make the front view
show a fist that the reference does not have.

## Exact proposed target and protected-interface scope

Thirty existing objects: prefix `Sleeve44P L ` or `Sleeve44P R ` followed by
exactly one of these fifteen suffixes. Enumerate the set in a writer guard;
do not hide or deform every object whose name contains sleeve/arm.

```text
asymmetrically seated stuffed arm insert
front padded fabric panel
rear padded fabric panel
front folded cuff edge
rear folded cuff edge
upper pinched cuff fold
lower pinched cuff fold
upper sewn panel join
lower sewn panel join
front red running stitch 1
front red running stitch 2
front red running stitch 3
rear red running stitch 1
rear red running stitch 2
rear red running stitch 3
```

Protected: all 14 shoulder-pleat/root-seam objects; all target vertices/faces
wholly proximal to a=-15 mm; bodice shoulder interface; sleeve axis; front X/Z
outline of the white panels; head/cheek locks, collar/tie, bow, skirt/hem,
feet, floor, material nodes and lighting. Record the proximal retained-face
hash and non-target evaluated hashes before changing anything.

## Materials and attachment facts

- Inserts use `Sleeve44P warm skin fabric.002` (Principled base color
  .76,.56,.36; roughness .92). Keep this material for the first geometry test.
  It is warm skin fabric, not a missing material or transparent placeholder.
- Panels use `Sleeve44P outer warm cotton.002` and
  `Sleeve44P inner warm fleece.002`; seams use
  `Sleeve44P compressed seam.002`; red stitches use
  `Dress red cloth.005`. Reuse them without node edits.
- Inserts and four panels are parented to `ReimuFumoRig`, with Arm_L/Arm_R
  weights and Armature modifiers followed by 022 body and 023 sleeve-root
  lattices. The insert has 74 weighted base vertices, evaluated to 338 vertices
  and 336 faces. Each panel has 104 assigned base vertices; its existing stack
  includes subdivision 1, solidify 0.78 mm, bevel, Armature and both lattices.
- Arm_L/Arm_R are children of root. Pose matrices equal rest matrices in
  frozen 028. Bone heads are X=+/-60 mm,Y=0,Z=60 mm; tails have Z=95 mm. These
  are not the measured shoulder positions after the static lattices. Do not
  reposition bones to match this view or double-apply the proportion fields.
- The other 38 sleeve pieces, including cuff trims, stitches and shoulder
  seams, are currently unparented and have only the static 022/023 lattices.
  They are not already proven to follow arm motion. Preserve this fact in
  the receipt. A replacement can receive the corresponding Arm bone with
  inverse-current-pose compensation, but this would still need a later
  movement test; rest-frame contact does not prove a working sleeve rig.

## Causal risk and stop condition

The main risk is over-compressing the cuff into the new hand, or extending
the cap below the white sleeve into the skirt. Validate actual inner-rim
clearance and skirt/hand separation, not only AABB overlap. Preserve the
root, use the capsule as support, and retain a soft flared outline instead
of turning the cuff into a flat disk. Reusing the existing material may still
render the cap warm/dark; do not try to correct a palette issue with more
extrusion.

The useful first-state witness is a small rounded skin-fabric end visible
inside and just beyond each soft cuff in side/three-quarter views, no deep
empty tunnel, no large fists, and unchanged canonical front arm angle and
shoulder seating. Repeat the same front-ray grid as a diagnostic only; total
front exposure is not the goal. Stop this module there for frozen visual and
contact review. No visual, rig, physics or final acceptance is claimed here.
