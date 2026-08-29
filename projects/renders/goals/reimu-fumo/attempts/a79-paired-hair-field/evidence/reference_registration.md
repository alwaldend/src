# A79 paired hair-field reference registration

## Result and scope

This packet freezes the source-pixel and protected-interface contract for the
A79 complete crown/rear hair owner. It is ready for representation preflight.
It does not approve geometry, infer a factory pattern, mutate Blender data, or
edit the goal record.

The source conclusion is one constructed visual assembly with three roles:

1. a close front/crown field behind the retained fringe;
2. a compact rear/base field that closes the crown and rear; and
3. a broad independent free leaf that owns the additional profile depth and
   long rear overlap.

The complete `1.14--1.23 Wh` profile may not be filled by one receiver, cap,
or shell. The compact field/free-leaf allocation is a calibrated visible-layer
contract, not a measured hidden cross-section.

Machine-readable authority is split deliberately:

- [reference_registration.json](reference_registration.json) contains source
  authority, camera mapping, interface bands, vetoes, and artifact paths;
- [contours_and_landmarks.json](contours_and_landmarks.json) contains raw
  source-pixel curves, audited scanlines, normalized targets, and uncertainty;
- [canonical_turn_masks.json](canonical_turn_masks.json) contains all 30
  frames, their groups, hashes, threshold-mask row spans, and selection
  provenance;
- [turn_frame_map.csv](turn_frame_map.csv) and
  [landmarks.csv](landmarks.csv) are compact builder tables.

## Source authority

| Source | SHA-256 | A79 role |
| --- | --- | --- |
| `canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` | Exact variant, front crown, aperture/fringe, scale |
| `canonical_turn_180.gif` | `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30` | Both 3Qs/profiles, rear, depth, continuous layer order |
| `clean_front.png` | `37813e03e04e4966f1dbe914e03a25a5f5ae561dcbf58b72677195c513ea48ca` | Graphic and clean-hairline cross-check only |
| `physical_front.png` | `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51` | Pile, stuffing, thin padded front construction |
| `physical_side.png` | `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d` | Edge roll, layer thickness, overlap, root burial |
| `turn.gif` | `b42368e921bd055d73fbbb7bf65c2509a9aaf190cab02f89824b92b4cb75ece4` | Low-resolution continuity and no-bald-region veto |
| `sofa.gif` | `7c9173f91e6b6c801a1c77e50f9135e86fc89319f3c0262c10312320b1af8589` | Pile, compression, and soft panel-edge veto |

The canonical front wins direct conflicts about exact frontal identity and
scale. The canonical turn wins direct conflicts about depth, hidden-side
silhouette, and layer order. Supporting sources may reject hard, thick,
floating, or implausible fabric, but they may not move a canonical target.

## Coordinate and camera registration

### Exact front

```text
Wh = 368 +/- 4 px = 1.000 = 132.0 mm
center x = 485 px
brown crown datum y = 231 px
x_wh = (pixel_x - 485) / 368
v = (pixel_y - 231) / 368, positive down
```

Register the candidate to left/right outer hair, center, and crown before
judging the aperture. Front alignment residual must be at most `.015 Wh`
before a `.03 Wh` critical-landmark result is meaningful.

### Canonical turn phase map

The filename says `180`, but the pixels contain a full 30-step rotation. The
phase solution is a nominal `12 degrees` per frame, anchored by frame `03` as
front and frame `18` as rear:

```text
yaw(frame) = (frame - 3) * 12 degrees, wrapped to [-180,+180]
```

This yaw is an evidence-backed phase map, not embedded camera metadata.
Absolute phase uncertainty remains `+/-1` frame (`+/-12 degrees`). Use one
unchanged turn-camera scale, `Wg = 244 +/- 5 px`, and do not rescale individual
frames.

| Frames | Nominal yaw | View group | Controller / use |
| --- | ---: | --- | --- |
| `00--05` | `-36..+24` | front continuity | `03` is turn-front; canonical still remains exact |
| `06--09` | `+36..+72` | three-quarter A | `07` (`+48`) with `06/08`; `09` is late closure |
| `10--12` | `+84..+108` | profile A | `10/11` are metric brackets; `12` is the transition curve witness |
| `13--16` | `+120..+156` | rear 3Q A | overlap emerges; zero support leak |
| `17--20` | `+168,180,-168,-156` | rear | `18` controls; `17/19` are phase brackets |
| `21--24` | `-144..-108` | rear 3Q B | opposite overlap continuation; never mirrored |
| `25--26` | `-96,-84` | profile B | direct opposite profile bracket |
| `27--29` | `-72..-48` | three-quarter B | `29` (`-48`) with cyclic `28/00` bracket |

All 30 frames participate through the disjoint groups above. Frame `00` also
acts as the cyclic phase bridge for the second front 3Q comparison; that does
not change its single ownership in the frame map.

The turn source is perspective and its intrinsics are unavailable. A valid
candidate overlay must fit one perspective camera/scale family to invariant
crown, head, and ground witnesses, keep it unchanged through the sequence,
and select only the nearest phase within the one-frame band. Turn alignment
residual must be at most `.03 Wh`. A raw orthographic side-by-side or a
view-specific image rescale is not calibrated overlay evidence.

The older A53 profile/rear curves were stored with a `233 px` local scale.
A79 retains those curves as raw pixels but normalizes every turn value with
the audited A73/A74 `Wg = 244 +/- 5 px`. Values normalized by `233` may not be
mixed with the A73/A74 scanline ratios.

## Front crown, aperture, and retained fringe

The asymmetric crown scanlines below are outer projection targets. The source
measurement uncertainty is `.011 Wh`; the integrated-render acceptance band
is `+/- .03 Wh` after registration.

| `v` below crown | Left reach | Right reach | Full width |
| ---: | ---: | ---: | ---: |
| `.052` | `.182` | `.226` | `.408` |
| `.133` | `.321` | `.353` | `.674` |
| `.242` | `.402` | `.435` | `.837` |
| `.351` | `.438` | `.470` | `.908` |
| `.459` | `.465` | `.484` | `.948` |
| `.568` | `.484` | `.492` | `.976` |

The refined central-fringe clefts are approximately `(u,v)=(.309,.372)` and
`(.667,.372)`, where `u` is measured from the canonical left hair bound. Their
separation is `.358 Wh`. The low blunt tip is `(u,v)=(.588,.677)`, or signed
`(x,v)=(+.087,.677)`. It is `.101 Wh` viewer-right of the cleft midpoint and
`.305 Wh` below the cleft row. Its lower third must remain at least `.12 Wh`
broad; a triangular spike fails.

At the refined cleft row, shape-preserving interpolation of the asymmetric
crown rails gives:

```text
left reach  = .443 Wh
right reach = .473 Wh
full width  = .916 Wh
```

The visible beige aperture has pixel bounds `x=374..596`, `y=358..580`, or
signed bounds `x=[-.3016,+.3016]`, `v=[.3451,.9484]`. Its exposed extent is
`.6033 x .6033 Wh = 79.6 x 79.6 mm`, with `+/- .03 Wh` tolerance. This is a
visible-pixel envelope, not a hidden face-panel pattern.

The seven-point full hairline, strict 13-point cheek/chin arc, 13-point crown
rails, and both independent cheek-lock contours are stored as raw pixels in
[contours_and_landmarks.json](contours_and_landmarks.json). The refined
clefts are apex witnesses; the older seven-point hairline contains shoulders
and sweeps and must not substitute its shoulder samples for both apexes.

A79 preserves the five fringe/temple panels and both cheek locks. Therefore:

- the new crown field owns the upper outer silhouette;
- the retained fringe owns the complete free aperture/hairline;
- the new visible outer surface must terminate behind the retained fringe;
- apparent stand-off over beige remains `.02--.05 Wh`;
- unintended root separation is at most `.01 Wh`; and
- no second brown line, beige island, parallel cavity rim, or gap wider than
  one review pixel may appear inside the aperture.

The retained cheek locks remain separate foreground pieces. Their visible
widths are `.185 +/- .015 Wh` and `.174 +/- .020 Wh`; crown to the lowest lock
is `1.098 +/- .03 Wh`. They may not become the side wall of the new field.

## Profile field versus free leaf

| Target | Profile A | Profile B | Status |
| --- | ---: | ---: | --- |
| Complete brown projected depth | `1.14 +/- .05 Wh` | `1.19--1.23 Wh`, uncertainty `.06 Wh` | Direct outer silhouette |
| Compact front/crown/rear field | `.77--.85 Wh` (`101.6--112.2 mm`) | Same construction band, independently overlaid | Calibrated visible allocation |
| Free-leaf overhang beyond field | `.36--.38 Wh` (`47.5--50.2 mm`) | Same role, unequal curve | Calibrated visible allocation |
| Free-leaf apparent edge | `.015--.03 Wh` | `.015--.03 Wh` | Supporting fabric band |

Frames `10/11` and `25/26` contain different source silhouettes. Profile B
must be compared directly and may not be a mirrored Profile A. Raw audited
profile scanlines for every `v=.102..922` row are preserved in
[contours_and_landmarks.json](contours_and_landmarks.json).

The compact field and free leaf must remain visibly different contributions:

1. the compact field carries the close crown/rear mass and tapers toward crown
   and underside;
2. the leaf owns the deepest reach and a bowed, tapering free lower point; and
3. red bow occlusion between threshold components is not empty 3D space.

The `.77--.85` and `.36--.38` bands are not receiver-only sections or literal
factory-piece dimensions. They are representation discriminators validated by
the two profile brackets. A candidate still fails if the compact field is an
egg or helmet that happens to hit its numeric band.

## Rear width, height, lobes, and overlap

Frame `18` is the rear controller; `17/19` bracket its phase. Its audited
threshold scanlines are:

| Fraction down rear height | Left reach | Right reach | Full width |
| ---: | ---: | ---: | ---: |
| `.10` | `.074` | `.189` | `.262` |
| `.20` | `.172` | `.250` | `.422` |
| `.30` | `.270` | `.291` | `.561` |
| `.40` | `.328` | `.361` | `.689` |
| `.50` | `.406` | `.414` | `.820` |
| `.60` | `.463` | `.447` | `.910` |
| `.70` | `.467` | `.455` | `.922` |
| `.80` | `.451` | `.455` | `.906` |
| `.90` | `.225` | `.463` | `.689` |

The `.922 Wh` row value is the frame-18 threshold mask at the sampled height.
The registered outer-silhouette target, including antialiased edge and phase
brackets, remains `.94 +/- .04 Wh = 124.1 +/- 5.3 mm`.

The complete crown-to-lowest-lobe height is `1.16 +/- .05 Wh =
153.1 +/- 6.6 mm`. Maximum width occurs at `.70 +/- .04` of that height,
equivalent to approximately `.812 +/- .082 Wh` below the registered rear
crown after propagating both uncertainties.

The lower edge has four or five broad unequal lobes. High-notch to low-tip
relief is `.10--.14 Wh`. Equal scallops, periodic spacing, a straight hem, or
separate identical paddles are categorical failures.

The foreground leaf is `.45--.50 Wh` wide and about `.10 +/- .05 Wh`
off-center. Prior prose conflicts on whether to call that offset screen-left
or screen-right, and the fixed rear camera reverses world X. A79 therefore
does not invent a signed world-X scalar: the registered raw frame `17--22`
pixels and overlap rails control the sign.

Across the rear arc, the long overlap moves laterally `.15--.25 Wh` from crown
toward the lower tip, with about `.05 Wh` uncertainty. The leaf's visible
root-to-tip drop is `.98--1.05 Wh`, uncertainty `.07 Wh`. Its exact bow-hidden
root remains unobserved by `.05--.08 Wh`. Frame-19 raw outer and overlap
polylines are retained as T-junction/free-edge witnesses; they do not require
four literal seams.

## Frozen receiver, fringe, and root relationship

The protected rung-003 receiver bounds are `X=-66..66`, `Y=-54.201..51.399`,
`Z=91.5..220.2 mm`. A79 hides exactly the old cap, three old rear locks, and
two old rear seam objects. It retains the receiver, five fringe/temple panels,
two cheek locks, crown seam, face witnesses, and every other context object.

The old cap-to-receiver baseline is `0.524..1.020 mm` positive clearance with
zero triangle intersections. For new paired skins:

- visible thickness is `2.0--5.0 mm` over at least `95%` of sampled visible
  area;
- the declared receiver root band is `+0.30..+1.30 mm` for at least `80%` of
  its samples;
- negative signed root samples are forbidden; and
- baseline-relative crossings may not increase.

The retained source already contains legacy crossings: left/right temple
transitions `218/202`, cheek locks `112/132`, and crown seam `36`. They are
frozen baseline state, not A79 repair authority. Do not demand zero crossings
from those protected objects or move them to make a new field easier.

No source or A78 probe supplies a millimetric new-field-to-retained-fringe seam
position. Each representation preflight must declare its hidden fringe root
band before build, then pass the screen-space hairline, one-pixel gap, signed
clearance, and no-new-crossing gates. Rear roots are similarly localized under
the bow but position-uncertain; visible continuous brown coverage is the
source fact.

## Absolute vetoes

Reject without averaging scores if any affected view shows:

- a helmet, concentric shell, egg, anatomical skull, mattress, or rounded box;
- a canopy, cape, curtain, card, blade, or full-width rear sheet;
- one inflated field filling the full profile instead of a compact field plus
  independent leaf;
- a bald rear, beige/support leak, detached brown island, or exposed root tab;
- a floating root, dark cavity, or gap wider than one review pixel;
- a duplicate hairline, parallel aperture rim, small round opening, face ball,
  or bulbous beige 3Q oval;
- a symmetric rear split, equal paddles, periodic lobes, narrow rod/capsule
  leaf, horizontal nape shelf, or lost diagonal overlap;
- Profile B obtained by mirroring Profile A;
- per-frame camera rescaling or an overlay judged above alignment tolerance;
  or
- hard plastic/taut-shell read or material detail used to counterfeit shape.

One-sided improvement also fails: front, both front 3Qs, both profiles, the
rear, and both rear-transition arcs must all preserve the registered owner and
layer order.

## Visual target and contact-sheet manifest

| Artifact | SHA-256 | Purpose |
| --- | --- | --- |
| [Raw all-30 contact](canonical_turn_all_30_contact.png) | `2ae6f070dde1f212b0f3bff529a4e65e23ef5330d9c668794f59c4e03bd4f49a` | Unmodified canonical-frame overview |
| [All-30 grouped overlay](overlays/canonical_turn_all_30_group_overlay.png) | `0ebd68d60adc73024dcc22813ada53c26bb8f8d5c46d9f7e58629293f65c0865` | Every frame, phase, group, and brown mask |
| [Exact front overlay](overlays/front_exact_target_overlay.png) | `71a02b0a876624bf0bafc778b6033f310739f661aded01f9ddb01f73e7199272` | Crown rails, aperture, hairline, locks |
| [Turn-front bracket](overlays/front_turn_target_overlay.png) | `686df7ee3ae236cbb4db6b84d9a3520863ff1bcdfb7c915346da71bbbdc56631` | Frames `00--05` continuity |
| [3Q A bracket](overlays/three_quarter_A_target_overlay.png) | `8ab6e7f23aadaf67505fa24b4b29db951edf0291a4faaa885166f622a01d66e4` | Frames `06--09` |
| [Profile A bracket](overlays/profile_A_target_overlay.png) | `d1cb7808dc475bf6eebf074264b576545ab1e80ee8dc444fc079e911e065d3e7` | Frames `10--12`, compact/leaf colors |
| [Rear 3Q A bracket](overlays/rear_three_quarter_A_target_overlay.png) | `1ecf454c44b98fa04c7dba28f58f2ad35618663868ae44e7522ed8864d51d6ee` | Frames `13--16` overlap emergence |
| [Rear bracket](overlays/rear_target_overlay.png) | `6b9f65708d2cc04a08734b28a2b49b79d5d4ff0f5445fc77a1cd43f3d6916547` | Frames `17--20`, scanlines and edges |
| [Rear 3Q B bracket](overlays/rear_three_quarter_B_target_overlay.png) | `145de68ba3f8b4e5fdf367f385cc91085409bf93bdbfe7e06c10fc0cad814f2c` | Frames `21--24` opposite overlap |
| [Profile B bracket](overlays/profile_B_target_overlay.png) | `b854b19e932bc7f3f863e592f89cdecf474cf7ed2e8cc794785886f676244a8a` | Frames `25--26`, direct non-mirrored target |
| [3Q B bracket](overlays/three_quarter_B_target_overlay.png) | `2e05ad853c0222e4f68f1e02ff7c9b6b8659b6b5a7da2917973672456e3672cb` | Frames `27--29` plus cyclic `00` |
| [Supporting-source board](overlays/supporting_reference_board.png) | `55094baf6ef6d8d9ba2befb8ff0c7295ed61e0e03539d3754be60a914b87a1a7` | All noncanonical construction roles |
| [Older-turn contact](support_turn_contact.png) | `d5c819a13f124513b25f9d52dbed94d57825f14fdbb5446ba50d08164676d5ea` | All four supporting turn frames |
| [Sofa contact](support_sofa_contact.png) | `a58f0b5ae2d84c9ef1a726e2391000cce216740e210b250732148635355ca30a` | Representative motion/compression states |

The yellow/cyan threshold masks are visible-pixel measurement aids. Pile,
shadow, and red-bow occlusion make them unsuitable as hidden mesh outlines.
The curated raw curves, scanlines, source authority, and uncertainty bands in
the machine files control any disagreement.
