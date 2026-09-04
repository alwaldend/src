# Reimu Fumo reference landmarks

`Wh = 1.000` is the outer front head-and-hair envelope width, excluding the
bow and isolated ribbon tips. Front comparisons align the envelope center,
crown, and width before measuring. Candidate cameras must be calibrated to
within `0.02 Wh`; otherwise the comparison is invalid rather than a model
failure.

## Canonical exact-variant controls

The canonical front has an active datum of `Wh = 368 +/- 4 px`, horizontal
center `x = 485 px`, and crown `y = 231 px`.

| Landmark                   |   Target or band | Tolerance |
| -------------------------- | ---------------: | --------: |
| Lock-excluded head height  |       `0.986 Wh` | `0.03 Wh` |
| Crown to lowest cheek lock |       `1.098 Wh` | `0.05 Wh` |
| Visible beige face width   |       `0.603 Wh` | `0.03 Wh` |
| Visible beige face height  |       `0.603 Wh` | `0.03 Wh` |
| Central fringe tip         | `(0.588, 0.677)` | `0.03 Wh` |
| Complete bow span          |       `2.038 Wh` | `0.05 Wh` |
| Combined sleeve span       |       `1.277 Wh` | `0.05 Wh` |
| Red skirt width            |       `0.916 Wh` | `0.05 Wh` |
| Ruffled hem width          |       `1.065 Wh` | `0.05 Wh` |
| Visible foot width         | `0.212-0.217 Wh` | `0.03 Wh` |
| Foot-center separation     |       `0.560 Wh` | `0.04 Wh` |

The central-fringe coordinates are normalized from the head envelope's upper
left. The exact endpoints used by a comparison must be shown on its annotated
reference; prose values alone are not sufficient evidence.

## Supporting depth and construction bands

The physical side is oblique. It controls bands and depth order, not an exact
pixel overlay.

| Landmark                      |   Allowed band |
| ----------------------------- | -------------: |
| Stuffed base-head depth       | `0.66-0.82 Wh` |
| Outer hair-and-lock depth     | `0.71-0.87 Wh` |
| Independent rear-hair reserve | `0.30-0.40 Wh` |
| Bow panel core thickness      | `0.03-0.05 Wh` |
| Bow with ruffle thickness     | `0.07-0.10 Wh` |
| Seated skirt/hem depth        | `0.99-1.15 Wh` |
| Foot-pod depth                | `0.32-0.42 Wh` |

The base head must read as a shallow gusseted cushion with a visible front
plane. Hair, fringe, bow, sleeves, and garment pieces must read as separately
rooted padded fabric panels, not a monolithic helmet, cards, armor, hollow
tubes, or a cone.

## Fixed review views

The asset faces `-Y`, uses `Z` up, and contacts the ground at `Z = 0`.
The machine-readable authority for measurement cameras is
`review_contract.json`. All five are square orthographic views at 512 by 512
pixels, 100 percent scale, and `0.292 m` orthographic scale:

| View                 | Camera location (m)     | XYZ Euler rotation (rad)   |
| -------------------- | ----------------------- | -------------------------- |
| Front                | `(0, -0.8, 0.13)`       | `(1.570796, 0, 0)`         |
| Rear                 | `(0, 0.8, 0.13)`        | `(-1.570796, 3.141593, 0)` |
| Side                 | `(0.8, 0, 0.13)`        | `(1.570796, 0, 1.570796)`  |
| Three-quarter        | `(0.52, -0.52, 0.135)`  | `(1.563997, 0, 0.785398)`  |
| Three-quarter mirror | `(-0.52, -0.52, 0.135)` | `(1.563997, 0, -0.785398)` |

A separate perspective camera may be used for presentation but not
measurement.

## Uncertainty and exclusions

- Do not infer depth from front images.
- Do not claim front-view precision for side or rear measurements.
- The canonical turn controls exact-variant side/rear silhouette and layer
  order; the physical views control fabric behavior and construction.
- The clean front controls facial graphics only where it does not conflict
  with the canonical front.
- Every critical canonical-view landmark must be within `0.03 Wh`; silhouette
  extrema and major gaps must be within `0.05 Wh`.
- Any visible floating, clipping, disconnected construction, or accidental
  tangency above `0.02 Wh` is a major failure regardless of averages.
