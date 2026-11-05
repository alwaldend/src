# Intermediate body proportion retention; full model still fails

Retain corrected 022b only as an intermediate proportion/cloth checkpoint.
No stage passes and no goal acceptance pointer is set. Preserve 021 and
the rejected first022 state for comparison.

Final candidate: out/reimu_fumo_finish/desktop_astra/body_022b_candidate.blend
SHA-256 `96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014`.
Source021 SHA-256
`9cfbd356b2cd2e377f566304886d380eff2e91b0e8ecee3c3e9004fe1b3e2f22`.
Protected A202 SHA-256
`a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5`.
All remain unchanged after pinned clean-reopen verification.

## What changed

Joined skirt/hem cloth replaces eight donor panels and eighteen stitches.
Terminal native lattices shorten the body, translate the head22.1mm down,
and reduce upper-bow height about its knot. Existing meshes, modifier
stacks and rig data remain; animation is not verified.

One disposable setup test caught non-unit native lattice grid coordinates
before opening Reimu. Normalizing observed extents repaired it; mesh and
curve fixtures then followed the intended field within0.021mm, kept X/Y
fixed, and recovered exactly when the modifier was removed.

First022, SHA
`0ff8cab54dfa5be530a12524f86e7238c13f90920ea3fd1864a27a6e6cdeef92`,
failed image review: zero-slope interpolation at every drape row produced
horizontal bands. Technical review also found the mouth assigned to the
body, displacing it5.5mm relative to the head. Corrected022b uses monotone
cubic row interpolation and an explicit mouth-to-head assignment. No drape
control values, width targets, or macro deformation values were changed.

## Image and measurement evidence

Pinned5.2.1 build9e2066aef7ef rendered all five fixed contract views after
opening the exact saved candidate. Approximately7seconds covered save,
reopen and all views. Diagnostic material colors are not final materials.
Final front PNG SHA
`419bd3e8aad478477936371a642f9d5d2313a8734bad7b80ebbc4add2a876308`;
three-quarter PNG SHA
`995e7a95a711f6cbc5e9d32fbc89522de34cedf24cf0341dd1a34441e29f66fc`.
Both are under out/reimu_fumo_finish/desktop_astra/body_022b_review/.

Independent image measurement by fumo_context used canonical front
Wh368px/crown231 and first022 front Wh203px/crown145. These global
landmarks are unchanged by the two022b source repairs. This is transferred
global evidence, not a claim that the first022 pixels equal022b pixels.

| Quantity, in head widths | Reference | 022 global shape |
| --- | ---: | ---: |
| Crown to ground | 1.666 | 1.670 |
| Crown to chin | 0.973 | 0.946 |
| Crown to skirt/hem boundary | 1.427 | 1.438 |
| Bow top relative to crown | -0.220 | -0.212 |
| Mid-bodice width | 0.486 | 0.709 |

Global-height uncertainty is approximately0.03Wh. The former+0.196Wh
height excess is resolved within that uncertainty. Torso width remains
too large by0.223Wh, approximately0.04Wh uncertainty. A front projection
that looks seated does not establish actual head/collar contact.

Reviewer body_022_blind inspected reference images, then all022b views,
before seeing021 baseline pixels or any implementation. Absolute scores
out of10: likeness5, silhouette5, construction4, identity6, contact3,
plush read4, presentation7. Completed-model and stage gates fail.

After baseline comparison, the reviewer retains022b as intermediate:
shorter body/sleeves and less rigid skirt contour improve the seated form,
without major new clipping or detachment. Remaining/regressed fidelity:
torso too wide, weaker white hem/spread, shorter chunkier tie, preexisting
floating collar/tie, helmet hair, rigid bow/sleeves, and shiny hard feet.
Root agrees with this limited retention, not with whole-asset acceptance.

## Technical verification and limits

Independent pinned clean-reopen audit confirms all14head intrinsic meshes
unchanged; rigid head motion maximum residual46.28nanometers, with zero
X/Y drift. Mouth residual16.47nanometers after its corrected role. Feet,
floor and inspected rig records exactly match021. New skirt and hem are
finite, consistently wound closed evaluated cloth, Euler0, without zero
edges or degenerate faces; all base vertices have positive skin weights.

First022 contact samples found collar-to-head nearest gaps18.202mm left
and18.160mm right despite5.204mm overlapping height bounds. These unchanged
collar/head shapes remain disconnected in022b pixels. Skirt/foot minimum
gap1.577mm; hem/foot1.328/1.434mm; hem/floor3.181mm; feet/floor0.100mm.
Those contact values were not fully resampled for022b; its separate audit
was scoped to repairs, topology, intrinsic geometry and protected inputs.
No full rig, UV, animation, materials or final delivery criterion passes.

## Next bounded work

Use022b as the experimental checkpoint and retain021 for comparison.
Narrow the torso/waist toward the canonical0.49Wh target, fit collar and
tie to the actual front body/head surfaces, and restore their cloth
proportions. Preserve the now-matching global height. Then reassess the
skirt spread and gathered white hem, sleeves, and head/hair against the
references. Do not reopen screenshot or GUI investigations for these
background-scriptable operations. Root remains sole model/goal writer.
