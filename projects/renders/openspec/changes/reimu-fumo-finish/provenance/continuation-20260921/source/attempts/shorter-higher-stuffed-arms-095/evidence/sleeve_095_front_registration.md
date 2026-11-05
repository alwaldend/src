# Front registration check, 095

Observed 2026-09-06T02:07:16Z. Image-only comparison of
[canonical front](../../../projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png)
and [087 front](chin_087_all_review/front.png). No geometry, model bounds,
scripts, preflights, camera parameters, or goal state inspected. Native
[annotated pair](sleeve_095_front_registration.svg) embeds the source images
and applies a vector transform to the candidate; the PNGs remain unchanged.

**Result: the head-datum fit does not pass independent validation at 0.02 Wh.**
The fitted datums coincide by construction, but the mouth/eye witnesses remain
high and the feet higher still. Scale/crop alone does not register this entire
figure. That observation does not identify whether the remaining mismatch is
shape, pose/projection, or both.

## Head datums and one transform

Coordinates use the original images' upper-left origins, x right and y down.
The candidate is 512 × 512; the reference is 1000 × 1000.

| Datum | Candidate observed pixels | Reference target |
| --- | --- | --- |
| Outer head width, excluding bow and independent cheek locks | `212 ±2` | `368 ±4` supplied |
| Head center x | `254.5 ±1` | `485` supplied |
| Crown y, excluding red bow knot | `143 ±1` | `231` supplied |

The continuous cap reaches approximately x149–360 around y252–262, with
partial edge coverage at x148. Outer-boundary estimates x148.5–360.5 give the
reported width/center. Bow and lower tied locks are excluded. A bounded
brown-pixel check supported the visual trace; it was not used as an unattended
whole-image segmentation. The crown's first brown row is approximately y143.

The nominal similarity transform, with no rotation or nonuniform stretching,
is:

```text
s  = 368 / 212 = 1.73584906
x' = 485 + s × (x − 254.5) = 1.73584906 x + 43.22641509
y' = 231 + s × (y − 143)   = 1.73584906 y − 17.22641509
```

Nominal fitted width, center, and crown residuals are zero; this is imposed,
not an independent calibration test. Trace uncertainty remains. The nominal
validation limit is `0.02 × 368 = 7.36` reference pixels.

## Independent validation

Eye points are the visible black top-line corners, not iris centers. Manual
trace uncertainty is approximately ±1 candidate pixel and ±2 reference pixels
for eye/mouth points. Foot contact excludes the soft cast shadow; its reference
uncertainty is approximately ±3 px and candidate uncertainty ±1 px. The darkest
sole pixels reach y843 in the reference and y483 in the candidate.

| Witness | Candidate `(x,y)` | Reference `(x,y)` | Transformed candidate | Residual `(dx,dy)` | Error / Wh |
| --- | --- | --- | --- | --- | --- |
| 1: Left eye, outer corner | `(198,263)` | `(375,446)` | `(386.9,439.3)` | `(+11.9,−6.7)` | `0.037` |
| 2: Left eye, inner corner | `(239,263)` | `(455,451)` | `(458.1,439.3)` | `(+3.1,−11.7)` | `0.033` |
| 3: Right eye, inner corner | `(278,263)` | `(522,452)` | `(525.8,439.3)` | `(+3.8,−12.7)` | `0.036` |
| 4: Right eye, outer corner | `(318,262)` | `(598,446)` | `(595.2,437.6)` | `(−2.8,−8.4)` | `0.024` |
| 5: Mouth midpoint | `(255.5,320.5)` | `(487,550)` | `(486.7,539.1)` | `(−0.3,−10.9)` | `0.030` |
| Lowest physical foot contact | `y483` | `y843` | `y821.2` | `dy −21.8` | `0.059` |

Face errors are Euclidean; the ground witness uses vertical error only.
These are nominal residuals, not precision metrology. Landmark and fitted-
scale uncertainty can move a marginal row across the limit, particularly the
right outer eye. It cannot justify declaring all validation witnesses within
0.02 Wh. No second fit to the mouth, eyes, or feet was substituted.

## Sleeve and hand under that same transform

Only visible pixels are recorded. The two sleeve/body junction witnesses are
not a reconstructed hidden throat or its physical cross-section. Each source
point has roughly ±2 px uncertainty (±3 for H3), before registration uncertainty.

| Mark | Visible candidate witness | Candidate `(x,y)` | Same-transform reference position |
| --- | --- | --- | --- |
| S1 | Upper white sleeve/body join, beside collar and beneath hair | `(305,350)` | `(572.7,590.3)` |
| S2 | Inner lower sleeve/body turn where the dark crease meets the waist | `(306,397)` | `(574.4,671.9)` |
| H1 | Hand first emerges from beneath the upper cuff edge | `(366,394)` | `(678.5,666.7)` |
| H2 | Outermost exposed hand point | `(378,410)` | `(699.4,694.5)` |
| H3 | Lower hand/cuff overlap | `(354,417)` | `(657.7,706.6)` |

S1 lands near the *visible reference attachment area*: the independently
traced white/red loop marker is around `(579,588)`. Its nominal separation is
only 6.7 px, but the different visible features and occlusion prevent treating
this as a physical throat-endpoint pass. Thus an assertion that the entire
upper root is far low/outboard is not supported after this head registration.

S2 remains approximately 11.4 px right and 7.9 px down from the reference's
visible inner lower-panel turn `(563,664)` (13.9 px, or 0.038 Wh, nominally).
That local difference survives the same scale/crop correction, though the
reference contour is soft and partly occluded. The hand points map onto the
reference's visible cuff-panel area, where the canonical front exposes no
corresponding hand. Scale/crop does not remove that occlusion discrepancy.

Therefore scale/crop explains a substantial raw image offset and weakens a
blanket low/outboard-root diagnosis, but is insufficient to explain all
remaining sleeve/hand differences. The true hidden throat endpoints remain
unmeasured. Do not infer a new root displacement from these visible witnesses
or from the unpassed whole-figure registration.

The reference-fidelity skill's datum/validation distinction and uncertainty
requirements informed this check. This is evidence for the coordinator's
camera/construction judgment, not approval of a candidate or a new procedure.

Source SHA-256:

- Canonical front: `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`
- 087 front: `18610642926ca2be36a12ff0b843ae7b7d69fd2bb0a437df93634374babeeca9`
