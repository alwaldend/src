# Head032: local tangent defect, not another head-size change

Read-only preflight of `bow_030_candidate.blend`, SHA256
`4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b`.
Viewed fixed front, side and three-quarter 030 pixels, canonical front and
physical side. Inspected frozen 027c/028 construction, 028 plan and technical
audit. A single pinned Blender 5.2.1 LTS/build `9e2066aef7ef` read-only probe
produced `head_032_probe.json` via `head_032_probe.py`. Source hash stayed
exact. No model, modeling helper, goal or canonical file was changed.

## Dominant cause

The remaining crown/side ridge is principally a tangent discontinuity in
the existing core, inherited by its hood, not a missing head-width increase
or excessive cloth thickness alone. The front panel is position-connected
to the rolled gusset, but their directions do not match. The current radial
smoothstep has zero endpoint slope; at the front seam it turns the panel
nearly front-facing immediately before the gusset turns outward. 028 changed
the broad upper profile but deliberately left this seam construction intact.

Actual polygon-normal jumps across the 108 core front/gusset interface edges
above 125 mm are 44.08–76.93 degrees, median 67.34 degrees. The equivalent rear
interface has median 13.88 degrees and maximum 28.45 degrees. These are real
candidate geometry measurements, not normals inferred from a photograph.

| Core interface witness | World XYZ, mm | Normal jump |
| --- | --- | ---: |
| Right temple | 57.272,-29.448,140.528 | 76.778 degrees |
| Upper-right transition | 35.723,-24.899,180.254 | 60.456 degrees |
| Near crown center | 0.940,-22.004,192.518 | 44.085 degrees |

At the temple witness the adjacent normals change from approximately
(0.009,-1.000,0.022) to (0.972,-0.218,0.085). That is the front-mask-to-side
corner visible in the fixed side view. The physical side supports a softly
continuous stuffed contour with cloth seams, not this sharp change of
direction. It does not supply a defensible exact target normal angle.

## Separating core, hood and fringe

The isolated surface profiles provide three useful discriminators:

- At Z140mm/Y-29mm, side-ray X is 57.418 mm on the core and 57.984 mm on the
  hood. Their normals are effectively identical. The hood is carrying the
  underlying crease with about 0.55 mm normal stand-off; the lateral interface
  is joined hood topology copied from the core, not an open hood-edge gap.
- The sagittal crown core itself turns abruptly around Y-22mm/Z192.53mm.
  Across the next 0.5 mm in Y its normal changes from approximately
  (0,-.966,.257) to (.012,-.515,.857), before considering any cover. Merely
  removing hood thickness cannot eliminate that core kink.
- Fringe overlap contributes a smaller visible lip. For the crown cut at
  Y-25mm, core/hood/fringe top intersections are Z186.732/187.658/189.316mm.
  The fringe sits 1.658 mm above the hood in this particular vertical
  projection; this is not a nearest-surface clearance measurement. Near
  Y-22.5mm the fringe's free upper edge reverses its Y-facing normal as it
  rolls around its thickness. Thus fringe-root seating must be reviewed
  after correcting the receiver, but thickening or adding more padded lobes
  would not address the dominant structural cause.

The current cover recipe is already explicit: hood normal stand-off 0.55 mm
with 0.7 mm inward Solidify, fringe 0.85 mm stand-off plus up to 1.65 mm
padding and 1.1 mm inward Solidify. A new shell layer is not the proposal.

## One causal geometry change

Fair the existing outer front transition into the existing gusset tangent.
Keep its seam endpoint positions, finite gusset, back panel, broad face and
whole head scale. This is local continuity repair of the retained topology,
not another annular wall, ellipsoid, shell, or padded-hair-family rebuild.

Initial bounded construction hypothesis:

1. Work only in the outer front meridional band, approximately normalized
   radius rho=.86 through the existing front seam near rho=.988, above
   Z125mm. Fade the operation smoothly from zero at 125 mm to full at 145 mm.
   This is roughly the last 7–8 mm of front-panel radius, not the facial center.
2. Replace the current zero-slope terminal blend there with a tangent-matched
   cubic continuation. At the inner band boundary retain existing position
   and first derivative. At the existing seam retain its exact position and
   match the meridional tangent of the first gusset band. Use actual source
   correspondence, not a nearest-face branch that changes at an overhang.
3. Change depth/Y only for the core and transfer that same local receiver
   displacement to existing hood/fringe geometry. Retain existing padding
   and thickness for this first test; do not tune them independently. Keep
   the frontal landmark X/Z coordinates fixed and report any derived
   Solidify normal-offset movement rather than claiming those offsets are
   byte-identical after recomputation.

For orientation, the present gusset's normalized radial derivative at its
front endpoint is .048, while its depth derivative is the local 4–11 mm seam
width. A tangent-compatible front slope is therefore about .083–.229 meters
of Y per normalized radius, instead of the nearly zero lateral terminal
slope. This supplies a falsifiable boundary condition, not a request to
change the entire radial field. The necessary displacement and monotonicity
must be measured on the eventual candidate; no unsaved modeling trial or
parameter sweep was performed for this note.

Exact bounded targets:

```text
Head028 sewn cushion
Hair028 crown and back hood
Hair028 traced padded fringe
```

Protected interfaces: original underside faces through 89 mm and all collar
root witnesses; lower-face grid through 125 mm; eye/mouth/eyebrow positions;
front fringe tips; cheek locks; front-seam positions and all rear/gusset
geometry; rear hair; retained body/limbs/feet; all bow geometry and cages;
materials, lights and rig pose. Existing broad-face X/Z and depth outside
the fairing band stay fixed. If receiver transfer would require modifying
a protected detail, stop and report that dependency instead of silently
expanding the three-object scope.

## Locks are a separate cause, not proof that the head needs widening

The locks are already closed padded geometry: at Z80–100mm their central
depth is 4.39–5.16 mm and projected width 13.59–14.55 mm; at Z120mm the width is
18.72 mm and depth 2.89 mm. Their free lower centerline remains nearly planar
in depth: the left center Y changes only about 0.21 mm from Z80 to100mm.
That explains why more padding alone did not remove the straight-card read.
The canonical lower hair framing appears fuller, but this preflight did not
derive a reliable pixel-to-world replacement outline. Keep both locks fixed
for the seam test; their taper/drape is a distinct later design question.

## Falsifier and stopping point

First verify a materially smaller front/gusset normal jump at the same
recorded seam edges while preserving exact seam positions and the protected
controls. Then inspect fixed side and three-quarter pixels under unchanged
materials. Stop after that single causal comparison.

The hypothesis is falsified if the measured core tangent jump is removed
but the same pronounced crown/side ridge remains. That would elevate the
fringe's independent upper cut edge/overlap to the dominant cause; it would
not justify another whole-head sphere or more shell thickness. Also reject
a new bulge, shelf, pinched cap, frontal-outline change, loss of root seating
or new bow/head collision. Passing the seam metric alone cannot pass plush
likeness. This note proposes a test and claims no visual or goal acceptance.
