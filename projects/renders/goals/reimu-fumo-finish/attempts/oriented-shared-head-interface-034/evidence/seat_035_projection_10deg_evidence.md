# Seat 035: far-wall rejection is projection-dependent

Yes: one plausible downward projection removes the previously identified
far-red-wall successor mechanism. The earlier pad-only rejection must be
qualified as applying to the recorded level and shallow-pitch review views,
not treated as a projection-independent conclusion about the photographed
reference. This diagnostic does not approve a pad edit or a new camera.

With unchanged source 033, all 1,824 level-view seat-first rays still lead
to the rear red skirt when the seat is excluded. At a single 10° downward
orbit, none of the still-seat-visible projected rays has a red-skirt
successor: 128 hit white hem and 341 hit the review floor. This can explain
why a photograph taken from above does not expose the same far red wall.
It can also expose a floor-visible opening; whether that reads correctly
was not rendered or visually evaluated here.

## Frozen input and exact transformation

Input: `foot_033_candidate.blend`, SHA256
`98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48`.
Blender 5.2.1 LTS, build `9e2066aef7ef`, pinned background invocation through
`bazel_agent`; normal exit and unchanged input hash.

No Blender camera object, geometry, material or visibility was modified.
Only mathematical ray origins and directions changed in the diagnostic:

```text
Target T = (0, 0, 0.13) m; orbit radius r = 0.8 m
Elevation theta = 10 degrees, front azimuth unchanged
C = T + (0, -r*cos(theta), r*sin(theta))
  = (0, -0.787846208, 0.268918544) m in evaluated float arithmetic
Camera Euler XYZ = (pi/2 - theta, 0, 0)
Ray direction = (0, cos(theta), -sin(theta))
              = (0, 0.984807730, -0.173648223)
Orthographic scale = 0.292 m; image = 512 x 512
dZ/dY = -tan(10 degrees) = -0.17632698
```

Over 100 mm of depth, the downward ray drops about 17.633 mm. That is
large compared with the 7–17 mm heights where the previous level rays
reached the rear red wall. The floor is at approximately Z=−0.700 mm.

The source-033 level control uses theta=0 with the same target, radius and
scale. Its seat-first count exactly reproduces the prior 032 level count.
This is one changed angle and a same-input control, not an angle sweep.

Root's extracted reference frames 11/16 were visually inspected; they
support considering an above-level camera as a causal hypothesis. They
were not camera-calibrated, and 10° is not measured reference truth.

## Projected between-feet sampling

Start with the existing inclusive level ROI x=214–297, y=438–474.
From its 3,108 pixel-center rays, retain the 1,824 first-seat-hit world
points. Project those same points into the 10° camera, round to nearest
pixel centers, and merge duplicates. This gives an irregular 1,145-pixel
mask within x=214–297, y=447–466. No favorable replacement crop was chosen.

The 1,824-to-1,145 count reduction is projection/rounding compression,
not a visibility improvement score. The diagnostic counts first surfaces
with the seat present, then the nearest non-seat surface from the same
ray origin after excluding the complete seat BVH. It does not use the
pad's back face as a successor or change any object state.

| Condition | Seat still first | Seat-excluded successor for those rays |
| --- | ---: | --- |
| Level, original ROI | 1,824 | Rear red skirt 1,824; everything else 0 |
| 10°, projected mask | 469 | White hem 128; floor 341; red skirt 0 |

For the entire projected mask, first-hit counts with the seat present are:
front red skirt 134, white hem 494, seat 469, floor 48. With the seat
excluded they are: front red skirt 134, white hem 622, floor 389. The red
hits are on the front lap (Y=−60.327 to −54.001 mm, Z=26.735 to
28.664 mm), not the far rear wall. Some original seat points are therefore
covered by the existing front cloth in the elevated projection before
any pad correction.

## Representative successor depths

Depth is additional distance along the camera ray, not vertical clearance.

| Elevated successor | Count | Depth min / median / max, mm |
| --- | ---: | ---: |
| `Hem026 curled cotton strip` | 128 | 96.964 / 105.686 / 111.583 |
| `Review floor` | 341 | 78.759 / 100.442 / 126.684 |

At pixel (268,456), the seat hit is Y=−43.154, Z=21.497 mm; its successor
is the distant white hem at Y=+60.803, Z=3.167 mm, 105.560 mm farther
along the ray. At (268,462), the seat hit is Y=−38.846, Z=17.263 mm;
its successor is floor at Y=+63.029, Z=−0.700 mm, 103.446 mm farther.
Thus “no far red wall” does not mean that a new near garment surface fills
the opening. The diagnostic changes which distant surface can be seen.

## Interpretation and limits

This qualifies the broad inference in `seat_035_successor_evidence.md`:
the recorded level/shallow-pitch successor evidence remains correct, but
cannot alone rule out a pad-only correction based on a reference photographed
from above. The fixed review views and `review_contract.json` remain
authoritative and unchanged. This result is not permission to choose a
flattering camera, relax a gate or disregard a real defect in those views.

No model change, candidate render, support-patch enumeration or shape trial
occurred. All rays treat evaluated geometry as opaque, use a 2 m distance
limit, and omit shading, transparency, antialiasing and subpixel integration.
Projection to nearest pixel centers introduces boundary quantization;
the mask is not the whole skirt or whole image. Excluding the entire seat
is not a prediction of any particular local pad deformation. Neither the
floor-visible opening nor the white-hem successor has passed visual review.

Root owns whether this causal qualification changes the next construction
choice. No additional angle, review-contract change or repair is proposed
by this bounded evidence packet.

## Artifacts

Compact diagnostic JSON: `seat_035_projection_10deg_evidence.json`, SHA256
`093d89a3cb572bfd255a7bce22657040ed3b8558fdc5085c2b279e04c609fa95`.

Full level/projected witnesses: `seat_035_projection_10deg_witnesses.json`,
SHA256
`630203e6240a967345735a8672bb68d6f6dcbc724f60a2d3cf3f491f4718f986`.

Read-only script: `probe_seat_035_projection.py`. All remain task-local
scratch; root owns canonical evidence links and final decision records.
