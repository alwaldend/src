# A79 canonical-turn hair-envelope probe

## Verdict

The 30 canonical frames provide useful **per-view visible brown-hair
constraints**, especially for rear coverage, profile layer order, and the
large trailing leaf.  They do **not** support a defensible calibrated 3D visual
hull.

Use the masks as fixed-camera regression evidence for the **assembled hair**
(new field plus retained fringe and locks).  Do not convert their pixel bounds
directly into a watertight volume or use them as evidence for a concentric
helmet.

## Evidence and method

- Source: all 30 local canonical turn frames at 498 x 498 pixels.
- Phase: frame 03 is front; frame 10/11 bracket +90 degrees; frame 18 is rear;
  frame 25/26 bracket -90 degrees.  Adjacent-frame uncertainty is 12 degrees.
- The mask classifier selects high-confidence dark-brown palette pixels while
  rejecting the red bow, beige face, grey background, and black feet/eyes.
- A head-region bound and connected-component filter remove most dress and
  floor noise.  Bow-separated brown islands are retained when substantial.
- Cyan in `hair_silhouette_contact_sheet.png` is the selected visible-brown
  evidence.  Yellow is its per-view bounding box.
- `measurements.json` contains every mask bound, scanline, occupancy value,
  phase, and selected-component bound.  `projected_envelope.csv` is the compact
  per-frame table.

This is deliberately a conservative colour segmentation, not semantic
matting.  It cannot distinguish same-brown crown, bangs, cheek locks, rear
panels, or deep brown skin shadows at contacts.  Lower front scanlines include
some retained lock/contact pixels and must not control the replacement cap by
themselves.

## Critical measured views

| View | Frame(s) | Visible brown bbox (px) | Union width (px) | Key constraint |
|---|---:|---|---:|---|
| Front | 03 | `[91, 40, 354, 350]` | 264 | Broad rounded crown; face opening and locks interrupt the lower mask. |
| + profile | 10 | `[150, 50, 428, 350]` | 279 | Main field plus a bow-separated trailing island extending 80 px farther right than the largest connected field. |
| + profile | 11 | `[160, 53, 439, 326]` | 280 | Bow-separated trailing island extends 91 px farther right. |
| Rear | 18 | `[161, 50, 389, 350]` | 229 | Nearly solid brown rear field through its middle and lower crown; no beige rear opening. |
| - profile | 25 | `[56, 52, 354, 337]` | 299 | Two major islands bracket the red bow tail: `[56,52,212,337]` and `[271,91,354,312]`. |
| - profile | 26 | `[57, 60, 345, 350]` | 289 | Main field `[198,80,345,350]`; trailing island `[57,60,164,299]` plus a small lower lock. |

The three-frame front median width is 264 px.  The bracketed full-hair profile
widths are 279.5 px on one side and 294 px on the other, or 1.059 and 1.114
times the front width.  These profile numbers include the large trailing leaf;
they are **not** core-head depth.

The rear three-frame median width is 229 px, 0.867 of the front median.  This
is a useful camera-space target but not evidence that the physical object
changes width: bow occlusion, phase error, perspective, and the front locks all
affect the visible-brown bounds.

## Shape evidence by scanline

The following are outer visible-brown spans; occupancy describes how much of
the span is actually brown rather than bow/face/background.

### Front frame 03

| Image y | Span (px) | Occupancy | Reading |
|---:|---:|---:|---|
| 102 | 214 | 0.958 | Upper crown is almost continuous. |
| 148 | 230 | 1.000 | Crown reaches a broad, non-pointed maximum. |
| 195 | 230 | 0.796 | Face opening begins to interrupt the brown field. |
| 242 | 247 | 0.462 | Outer temples/locks remain wide while the beige face occupies the center. |
| 288 | 208 | 0.563 | Lower fringe/locks taper but do not collapse to one central point. |

### Rear frame 18

| Image y | Span (px) | Occupancy | Reading |
|---:|---:|---:|---|
| 110 | 73 | 1.000 | Narrow visible root beneath the bow knot. |
| 155 | 133 | 1.000 | Rear field expands steadily below the root. |
| 200 | 187 | 1.000 | Continuous central rear coverage. |
| 245 | 214 | 1.000 | Broad lower rear field; a bald opening is incompatible. |
| 290 | 213 | 0.962 | Coverage remains broad near the leaf tips. |
| 320 | 170 | 0.712 | Separate lower tips/notches create the final taper. |

### Profile brackets

At mid/lower head rows, the outer brown span is roughly 242--280 px, but only
about 40--60% is occupied by visible brown.  The red vertical bow tail is in
front and separates the main crown/rear field from the trailing brown leaf.
Therefore:

1. the profile must retain both lobes and their full combined envelope;
2. the visible result must not become one inflated uninterrupted side wall;
3. hidden continuity under the bow is unresolved by these pixels and cannot
   be used as an acceptance criterion;
4. the trailing leaf must project substantially beyond the crown field on
   both sides, rather than being absorbed into a symmetric helmet.

## Layer-order constraints

The turn establishes only the following defensible order relations:

1. **Bow over brown hair:** the top bow and long red side/rear tails occlude the
   brown crown and trailing leaf in profile and rear views.
2. **Brown rear field over torso/collar:** frame 18 keeps a continuous brown
   mass down to the red collar/dress, ending in several lower tips.  A beige
   rear leak or exposed receiver is a hard failure.
3. **Trailing brown leaf behind the red tail:** frames 10/11 and 25/26 show a
   large brown island outside the crown projection with the red tail between
   it and the main visible field.
4. **Same-brown panel ordering unresolved:** the RGB turn cannot reliably
   decide whether crown, rear base, and individual rear leaves overlap one
   another at hidden seams.  That must come from construction references or a
   separately controlled side/rear photograph, not this mask.

The machine-readable component bounds quantify the profile overshoot but label
components by colour connectivity, not semantic part name.

## Why a calibrated visual hull is rejected

Under a fixed orthographic turntable, a silhouette's projected width should be
nearly identical at angles separated by 180 degrees.  Across the 15 antipodal
frame pairs here, visible-brown width disagreement has:

- median: 12.54%;
- maximum: 38.35%;
- front/rear canonical pair: 13.26%;
- profile pair 10/25: 6.69%;
- profile pair 11/26: 3.11%.

That exceeds the 5% reference-fidelity silhouette tolerance before accounting
for unknown camera intrinsics, frame-to-frame registration, red-bow occlusion,
and same-colour layer ambiguity.  A 3D intersection made from these masks would
encode missing/occluded pixels as false concavities and would likely reproduce
the already-rejected helmet/card families.

## How A79 should consume this evidence

- Render the complete assembled hair at the fixed front, both profiles, rear,
  and both three-quarter cameras.
- Align by stable head/face landmarks before comparing masks.
- Treat an inward miss against a high-confidence cyan outer edge as a likely
  coverage failure.  Treat outward overshoot greater than about 5% of front
  hair width (13 px) as a likely silhouette failure where the reference edge
  is not bow-occluded.
- Require rear frame 18's continuous central/lower coverage and scalloped
  bottom taper.
- Require a distinct trailing profile leaf beyond the crown field, with the
  red bow tail visibly in front.
- Do not compare the replacement module alone to these masks: retained fringe,
  temple, and cheek-lock objects contribute to the reference silhouette.
- Do not infer millimetres from these frames.  The known 25 cm product size
  does not calibrate a seated, cropped, perspective turn without a stable
  full-object landmark and camera model.

## Artifacts

- `hair_silhouette_contact_sheet.png`: selector-facing 14-view overlay sheet.
- `measurements.json`: all 30 masks, scanlines, components, phase metadata, and
  antipodal consistency measurements.
- `projected_envelope.csv`: compact frame-level envelope measurements.
- `masks/`: one binary high-confidence mask per frame.
- `overlays/`: one source-aligned cyan/yellow overlay per frame.
- `analyze_visual_hull.py`: reproducible local-only extraction script.
