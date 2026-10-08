# Arm 100: lower-root interpolation collision

## Verdict

The scalar reconstruction reproduces **62 actual arm/sleeve triangle pairs per side**, all resolved noncoplanar crossings. Every pair is at retained sleeve **cell row 0**, not the cuff. Fifteen arm triangles join source UV rings 7 and 8; their 17 vertices and complete center-to-vertex radial paths remain at least **0.800219 mm** from the sleeve, but **14 connecting arm edges per side cross the finite sleeve wall**. This is a skin interpolation failure between clear fitted spokes at the lower root entrance. It is not a failure of the distal opening or evidence that the radial distance march missed these same spokes.

No Blender/native process, model modification, parameter trial, or geometry candidate was used. The diagnostic reads hash-bound actual arrays and writes only its script, JSON, and this report.

## Locations and complete contact evidence

All dimensions below are millimetres. A is the 100 builder's sleeve-local longitudinal coordinate, measured from retained root `(±24, -3, 70)`, along `normalize((±0.82, 0, -0.572))`. Source natural UV-sphere local Z identifies the original latitude ring; it is not world X. Ring indices are zero-based in the 25-section proximal-to-distal list, including poles.

| Measured intersection-segment bounds | Left | Right |
| --- | --- | --- |
| World X | -19.947162 to -19.499676 | 19.499043 to 19.889555 |
| World Y | -6.452540 to 1.564433 | -7.038791 to 0.804394 |
| World Z | 63.510567 to 64.015144 | 63.510533 to 63.945866 |
| A | 0.021694 to 0.333299 | 0.021195 to 0.329428 |
| Triangle pairs / distinct arm triangles | 62 / 15 | 62 / 15 |
| Sleeve half B / root-rim pairs | 31 / 31 | 31 / 31 |
| Crossing connecting edges | 14 | 14 |
| Crossing-edge lengths | 3.570757 to 3.921400 | 3.570757 to 4.095285 |
| Minimum complete affected radial-path gap | 0.800225 | 0.800220 |

Half B is the evaluated sleeve's second index half; this report does not independently relabel it inner/outer. Its contact faces are immediately adjacent to the root rim. The cuff opening begins at A = 38.825035 mm, far from all these contacts. There are no own-sleeve contacts elsewhere in the exhaustive triangle-pair reconstruction.

The JSON records every arm and sleeve triangle index, actual intersection endpoints, arm vertex IDs, endpoint barycentrics, A interval, and verified sleeve row/chart coordinates. Sleeve index interpretation is supported by complete source-067/068 quad/rim mapping and paired evaluated-half midpoint correspondence, not a visual guess. Arm/sleeve topology identity and protected sleeve fingerprints are checked against the actual 095 export and 099 evidence.

## The spokes surrounding the crossings

| Source ring | Natural local Z | A in 100 | A in passed 099 | Relevant vertices per side |
| --- | --- | --- | --- | --- |
| 7 | -0.608761430 | -1.871630 | -3.200795 | 9 |
| 8 | -0.500000000 | 1.161947 | -0.536716 | 8 |

The failed faces span **3.033577 mm axially**, joining the last pre-root section to the first enclosed-lumen section. Relevant angular directions are approximately 150–210 degrees in source natural XY: the lower transverse root/throat neighborhood, spanning depth directions. Both rings use raw capsule factor 1, so this is the full middle profile, not either rounded axial cap.

- Ring 7 vertices: `745, 768, 791, 814, 837, 860, 883, 906, 929`. All retain radial scale 1; point gaps are 1.818514–2.032633 mm.
- Ring 8 vertices: `767, 790, 813, 836, 859, 882, 905, 928`. Fitted scales are 0.781332–0.859562 left and 0.781332–0.849476 right; point gaps are 0.800220–0.800292 mm.
- All 17 actual center-to-vertex segments per side have zero sleeve hits and independently measured continuous clearance at least 0.800219 mm. Reconstructing their actual positions from source natural radial direction, recorded section center, capsule profile, and fitted scale agrees within 4.61 nm.
- Each of the 14 crossing connecting edges has two sleeve-surface hits. Example edge `(814, 836)` is 3.921400 mm left / 4.095285 mm right; it joins these clear ring-7/ring-8 vertices through the finite wall. The JSON includes every crossing edge and hit position.

## Why 099 passed does not certify 100's interpolation

The same scalar whole-surface check independently finds **zero pairs for 099 on both sides**, with identical mesh indexing. However, its corresponding rings 7 and 8 were both before A = 0 and only 2.664079 mm apart axially. Their corresponding affected edge lengths were 2.725730–2.925355 mm left and 2.780010–3.401648 mm right. In 100, the changed endpoint/profile moves the same UV sampling across the root feature and retains full middle girth behind it before the next ring is fitted smaller inside.

Thus unchanged UV topology is not unchanged physical feature coverage. Clear radial spokes and clear section-center chords do not guarantee that the surrounding straight skin edges or triangular faces avoid a varying, nonconvex receiver boundary. Here actual edge crossings directly demonstrate that missing guarantee; face-interior-only speculation is unnecessary.

The supported causal direction is receiver-aware skin interpolation/coverage across the root entrance. This report does not prescribe a radius dose, ring count, sweep, or approved replacement geometry. Merely subdividing the same frozen intersecting faces preserves their intersections; additional samples would need a corresponding refit/construction change. The full triangle gate correctly rejected 100 and must not be weakened or waived at the root.

## Bindings and limits

The script extracts only hash-guarded pure scalar definitions from prior diagnostics; it does not execute old builders. Its JSON contains all input digests, observation time, mapping checks, raw contact evidence, per-vertex path evidence, and comparison data. Source 098's sleeve geometry is bound through unchanged sleeve fingerprints to the actual retained 094 geometry in the 095 export; rejected 35/38 mm sleeve-root trials are not used.

| Artifact | SHA-256 |
| --- | --- |
| `arm_100_collision_diagnosis.py` | `826a8de325858705d4c2661b8f87b43382066231abf2d61030dfc12953d972dd` |
| `arm_100_collision_diagnosis.json` | `aee81cfcc7c3f185657c89e6940b670b7900f7df82bc27458a4d6e6384193d40` |
| `arm_100_arrays.json` | `a8fc4c8f14b764091bf38e1b43ad901960d6c728adaaa0daee8d274b6ef521fe` |
| `arm_100_preflight.json` | `41e0fcc2f52b185c12ad57e5bcca6cc34b6684b0234911362a43ddc3599cf42a` |
| `arm_099_arrays.json` | `151523770feaafef4c7ff1650180b8fc5ea9e7a6b51ece4c671f38529c336806` |
| `arm_099_preflight.json` | `d1fe957f8936e8dea59ad18610c5c916b354866872e084b11dd2b428c1f2bce4` |
| `sleeve_095_actual_assembly.json` | `4a23e0161fa80d992aedce2598ac7943b887bf83c2d5f6b1f3f10cd74783c3bf` |

This diagnosis covers actual own-sleeve crossings and the implicated point/path neighborhoods, not a new full-asset acceptance run or a claim that any proposed replacement arm passes.
