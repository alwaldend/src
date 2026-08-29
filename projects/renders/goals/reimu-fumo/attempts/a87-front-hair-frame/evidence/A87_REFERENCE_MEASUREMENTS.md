# A87 S02 front-hair reference measurements

## Scope and method

- Candidate: `packets/s02/front.png`, `front_ids.png`,
  `three_quarter.png`, and `three_quarter_mirror.png`.
- Primary frontal controller: tracked `canonical_front_25cm.png`.
- Primary depth/layer controller: sampled frames from tracked
  `canonical_turn_180.gif`.
- Cross-checks: tracked `clean_front.png`, `physical_front.png`,
  `physical_side.png`, `turn.gif`, and `sofa.gif`.
- Coordinates below are image-space ratios normalized by visible head/hair
  width `W`. `x=0` is the left head edge and `y=0` the hair crown. Positive
  `y` points downward. Canonical uncertainty is about `0.01--0.02 W`;
  physical/animated supporting views are `0.03--0.06 W` because of pile,
  perspective, blur, and occlusion.
- The canonical brown-hair connected component measures about 369 px wide
  (`x=301..669`) and 404 px tall (`y=232..635`). The clean-front component
  independently measures about 174 px wide (`x=143..316`) and 191 px tall
  (`y=110..300`). S02's face/head ID is 371 px wide (`x=72..442`); its center
  hair panel is 328 px wide (`x=92..419`) and its complete three-panel hair
  frame is about 402 px wide (`x=55..456`).

## Normalized landmark comparison

| Landmark | Canonical target | Clean-front cross-check | A87 S02 | Difference / verdict |
| --- | ---: | ---: | ---: | --- |
| Center-panel crown coverage | no exposed head-base arc; about `1.00 W` at the crown/temple silhouette | no exposed arc; about `1.00 W` | center panel `0.884 W`; exposed pale arcs total about `0.11 W` | **Fail:** cap is visibly too narrow and reads as a front plaque. |
| Left forehead notch `(x,y)` | about `(0.31, 0.37)` | about `(0.25, 0.39)` | `(0.25, 0.47)` | S02 notch is about `0.10 W` too low versus the canonical and too far outward. |
| Main bang tip `(x,y)` | about `(0.57, 0.67)` | about `(0.55, 0.67)` | `(0.48, 0.68)` | Depth is good; **tip is about `0.08--0.09 W` too far left/too centered**. |
| Right forehead notch `(x,y)` | about `(0.67, 0.38)` | about `(0.66, 0.37)` | `(0.65, 0.48)` | `x` is good; notch is about `0.10 W` too low. |
| Main bang asymmetry | broad left sweep; tip lies much closer to right notch than left (`~0.26 W` vs `~0.10 W`) | same directional skew (`~0.30 W` vs `~0.11 W`) | weak skew (`~0.23 W` vs `~0.17 W`) | **Fail:** reads generated/symmetric rather than the canonical swept fringe. |
| Left side-lock bounds | width `~0.20 W`, root `y~0.56 W`, tip `y~1.08 W` | width `~0.20 W`, tip `y~1.07 W` | width `~0.19 W`, root `y~0.36 W`, tip `y~0.89 W` | Size is close, but whole lock is about `0.19--0.21 W` too high. |
| Right side-lock bounds | width `~0.20 W`, root `y~0.56 W`, tip `y~1.09 W` | width `~0.22 W`, tip `y~1.09 W` | width `~0.19 W`, root `y~0.37 W`, tip `y~0.85 W` | Same high placement; right tip is about `0.24 W` too high. |
| Side-lock lateral placement | left center `x~0.15`, right center `x~0.83` | left `~0.17`, right `~0.79` | left `~0.05`, right `~0.94` | **Fail:** both are about `0.09--0.11 W` too far outboard. |
| Side-lock inward lean | left about `13 deg`, right about `17 deg` toward face | visibly inward on both sides | left about `6 deg`, right about `8 deg` | Too vertical; reinforces earmuff/capsule read. |
| Lock end shape | tapered, angular fabric-leaf points with unequal ends | same | broad rounded capsules | **Categorical fail:** wrong manufactured silhouette, not merely a small proportion error. |
| Complete hair-frame width | about `1.00 W` | about `1.00 W` | about `1.08 W` | Overall frame is too wide even though center cap is too narrow: width is misplaced into the side capsules. |
| Complete hair-frame height | about `1.09 W` | about `1.10 W` | about `0.89 W` | About `0.20 W` too short, entirely explained by the locks starting too high. |

## Cross-view construction evidence from all references

- `canonical_turn_180.gif`: front and three-quarter frames preserve a broad,
  irregular fringe with several staggered fabric tips. Side hair continues
  downward toward the collar and presents as thin constructed leaves. At
  three-quarter, the cap is a continuous wrapped hair mass rather than a
  front sheet ending in a long vertical cutoff.
- `physical_front.png`: proportions vary, so it does not override the table,
  but it confirms a pile-covered, irregular crown, asymmetric central sweep,
  and tapered lower side pieces. It provides no support for rounded vertical
  capsules.
- `physical_side.png`: the hair and bow are visibly thin front/back fabric
  panels with free edges, layering, and modest stuffing. The lower hair edge
  forms points and planes, not inflated rods.
- `turn.gif`: despite low resolution and a different pose, its front and rear
  frames confirm multiple staggered lower hair points and continuous side/rear
  coverage. The front identity does not reduce to three equal lobes.
- `sofa.gif`: blur and perspective preclude exact dimensions, but the visible
  side locks hang below the face framing and compress as fabric. Their roots
  do not start near the crown shoulder.
- S02 three-quarter renders: the center panel seats against the face, but its
  lateral/rear edge is a tall, nearly vertical cutoff and each side lock reads
  as a convex tube. Missing rear/crown hair is outside this coupon's scope;
  nevertheless the front-frame interface must leave a plausible overlap seam
  for that future owner rather than a detached plaque edge.

## Honest gate and smallest reachable correction

**Reject S02 as the A87 survivor.** The main tip depth is already correct and
must be protected, but four identity-level failures remain: narrow crown
coverage, weak/symmetric fringe sweep, locks placed about `0.20 W` too high
and `0.10 W` too far outboard, and capsule-shaped lock ends.

The smallest causally adequate next edit stays within the same three panel
owners:

1. widen/wrap the center panel by about `0.05--0.06 W` per side at crown and
   upper temple while preserving its current `~0.68 W` main-tip depth;
2. move the main tip right by about `0.08 W`, raise both forehead notches by
   about `0.10 W`, and keep the right notch near its current `x`;
3. translate both side locks down `0.20 W` and inward `0.10 W`, increase their
   inward lean to roughly `13--17 deg`, then directly shape each lower two
   rows into unequal tapered fabric-leaf ends;
4. render front plus both three-quarter views. Keep only if the pale crown arcs
   disappear, the main bang has the canonical right-skew, and neither lock
   remains capsule-like or creates a detached/intersecting seam.

Do not add fibers, wrinkles, rear hair, bow, face graphics, or other modules
to compensate for these macro front-frame failures.
