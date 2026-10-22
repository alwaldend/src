# Seat 035 occlusion correction: plan only

The causal target is `Garment42 compact internal seat pad`, not global skirt
size. Its exposed lower/front surface makes the dominant near-black band.
A local, closed-envelope retraction is worth considering, but pure inward
translation cannot hide that surface in the fixed frontal projection, and
removing it does not by itself close the opening behind it. No candidate,
shape parameters, construction helper or model write was produced here.

## Evidence and exact construction boundary

Source: retained `head_032_candidate.blend`, SHA256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
This note reuses `skirt_035_preflight.md` and its first-hit evidence, SHA256
`8d9437f11c5160b5256df53d9eaa7db3ecdcc38700890664f97e47427c376442`.
Any future edit must rebind to the exact retained post-033/034 input.

The pad is an inherited A42 mesh. The nearby 022/023 recipes preserve its
intrinsic mesh; neither constructs a new pad primitive. Current receipt:

- Parent: `ReimuFumoRig`.
- Material: `Dress red cloth.004`.
- Modifier order: `ReimuFumoRig` ARMATURE → `022 non-destructive rest
  proportion` LATTICE → `023 narrow waist, fixed height` LATTICE.
- No Solidify, Subsurf or live skirt/hem dependency is listed on the pad.
- Evaluated bounds: X=−52.500 to +52.500 mm, Y=−44.000 to +40.000 mm,
  Z=10.500 to 34.579 mm.
- Evaluated geometry digest in the 032 writer receipt:
  `315398f1f84fe71da607a479e439b3c08f8ba78e4891b4ba9776a51d079c5292`.

The original A42 primitive recipe and per-vertex attachment groups are not
exported in the bounded evidence read here. Do not infer them from the box
or assign fresh unit weights by assumption. Preserve the existing mesh's
attachment data, parent and modifier order when authoring any correction.

`body_skirt_022b.py` explicitly includes this pad in `SUPPORTS` alongside
the original cream leg roots and black pods. It casts downward rays to
those surfaces, then bakes skirt height with a nominal 1.8 mm support offset
and rounded maximum. Subsequent 022/023 fields alter the evaluated geometry;
the 1.8 mm construction constant is not a measured current clearance.
There is no live cloth simulation or shrinkwrap linking skirt to pad in
the inspected stacks. Moving the pad would not automatically re-drape the
retained red skirt or Hem026.

## Proposed narrow target and retained interfaces

Consider only a local retraction of the exposed lower/front pad envelope
into the existing concealed seat volume. Keep it a closed, outward surface;
do not delete its front faces and leave an open mesh. This is not a whole-pad
translation, a global scale, a material darkening/lightening fix or a new
skirt shape trial.

The following interfaces must remain protected for such a test:

- The pad's upper surface that supplied skirt support. Preserve its actual
  retained support witnesses, not merely its maximum Z or a guessed plane.
- Any pad regions currently embedding the retained cream leg roots or
  contacting the foot assembly. The existing evidence does not enumerate
  that contact patch; exact preserved vertices must be selected from the
  retained input before a sole writer moves adjacent geometry.
- `Skirt022 joined gathered panels`, including both evaluated boundary
  loops and the existing front opening; `Hem026 curled cotton strip`,
  including its sewn root; the current leg/foot surfaces; bodice, sleeves,
  rig, and shared 022/023 cages. Do not regenerate them from old recipes.

The disjoint exposed patch is the only intended moving region. If it cannot
be separated from the required support/attachment surfaces, a pad-only
correction is not feasible under those guards. Report that conflict instead
of silently moving the hem, feet or upper seat.

## What inward/upward motion can and cannot do

At front pixel (250,454), the pad is first hit at Y=−38.128 mm,
Z=16.793 mm. The front red hem near the center is approximately
Z=26.736 mm. The opening below that raised front lap is real: lower rays
already reach the rear skirt at positive Y. Moving a pad point farther
back in Y while retaining its X/Z does not change its orthographic image
position and leaves it first hit until another surface overtakes it.
An arbitrary small inward nudge therefore has no demonstrated mechanism.

Retracting the lower/front envelope upward could put those points behind
the existing front skirt, provided the resulting closed surface stays
below the protected upper support and does not lose leg embedding. The
roughly 10 mm difference between the sampled exposed seat height and the
front-center hem is a scale warning, not an authored displacement or a
verified available clearance. The white hem's scallops and oblique views
require actual first-surface checks, not that single central height.

Even successful removal of the seat from front rays may expose a larger
far-wall red strip or an empty cavity. Keeping the hem and feet fixed means
that outcome cannot be ruled out by pad geometry alone. Thus the honest
prediction is narrower than “the band disappears”: the visible pad band
should cease to dominate, and its replacement view must also read as a
plausible seated garment. No evidence here establishes that both are
simultaneously achievable with a pad-only edit.

## Strongest objection and first-hit falsifier

The strongest objection is that the internal pad currently fills a real
front opening. Hiding it may exchange one wrong dark surface for a deeper,
more conspicuous cavity while destroying its support or leg attachment.
A coverage correction involving outer cloth would then be a different
construction scope for root to consider, not an automatic expansion.

For a future authorized test, reuse the exact camera and the recorded
x=214–297, y=438–474 pixel-center grid. Baseline: 1,824 of 3,108 rays first
hit the seat, 656 the rear red skirt, 470 white hem, 92 foot pods and 66 no
surface. These counts are diagnostics, not acceptance thresholds.
Reject the mechanism if the dark-band witnesses remain exposed seat, or
if their successor hits become a visibly deep rear wall/empty cavity.
Reject any changed protected support/attachment, detached cloth, new visible
intersection or improvement requiring hem/foot motion. Inspect front and
both fixed three-quarter pixels as well as first-hit labels; do not certify
occlusion from bounds or a lower seat-hit count alone.

Root owns the final proceed/revise decision. This note proposes a bounded
test mechanism and its strongest counterexample, not acceptance or authority
to construct it. Red interior drape and the rigid white collar remain
separate deferred questions.
