# A68 C2 bilateral-Grab fixed-view review

## Verdict

**Reject C2, do not promote it, and do not build another local-Grab checkpoint
on it. Reset to C0 and change the cap representation or deformation strategy.**

C2 creates paired upper-rear dimples, but it does not materially reduce HD-01:
the cap remains an oversized rounded helmet with a deep, nearly vertical rear
wall in the side and both three-quarter views. The change affects at most 507
pixels (0.124% of a 640 by 640 frame) above a 1% RGB difference. Rear coverage
is preserved and no new bald crown opens, but this is not enough to compensate
for an unchanged identity-defining silhouette.

HD-01 has now survived two reviewed native-Grab checkpoints. Under the
reference-fidelity reset rule, another stronger point Grab is the wrong next
move: restore C0 and replace the uniform spherical shell with a shallower
constructed cap/rear-panel form, or apply a genuinely broad low-frequency
deformation whose controlling side silhouette changes without local dents.

## Implementation-blind absolute review

The absolute phase was delegated to a context-light reviewer that received
only the five C2 images, the canonical front, selected canonical-turn frames,
the physical front/side photos, and the neutral-sculpt stage. It was prohibited
from inspecting C0, C1, the blend, object names, measurements, scripts, reports,
or the intended edit. Only after that verdict did this report compare prior
checkpoints and implementation diagnostics.

- Unlabeled same-subject recognition: yes, as Reimu/Fumo; exact physical
  reference-variant fidelity remains weak.
- Intended-medium read: no. It reads as smooth hard CG/foam or molded panels,
  not sewn and stuffed fabric.
- Overall reference likeness: 5/10.
- Macro silhouette and proportions: 4/10.
- Constructed-plush construction: 3/10.
- Identity-defining features: 6/10.
- Contact, attachment, and occlusion: 3/10.
- Intended-medium read: 3/10.
- Presentation readability: 7/10.
- Major visible failure: yes.
- Absolute decision: reject.

The five largest visible discrepancies, in impact order, are:

1. The bow is too wide, rigid, planar, symmetrical, and ambiguously rooted;
   the references show compact, floppy, folded loops and tails.
2. The hair/head remains a deep, smooth helmet with a vertical rear wall;
   the references show a shallower stuffed head framed by layered thin hair
   panels.
3. The cone-like dress, blocky torso, and hollow angular sleeves do not read as
   a compressed seated plush.
4. Sleeve roots and foot/hem depth order show floating, tangency, and black
   void-like gaps.
5. Raised hard eye plaques, broad cheek bands, rigid side locks, and jagged
   temple cutouts do not match the soft reference construction.

View-specific veto checks:

- Front: no baldness, but the crown is over-smooth and helmet-like; a faint
  center dent/seam and mechanically exact symmetry remain.
- Side: HD-01 is strongest. The upper-rear change reads as a local dimple, not
  a reduced rear plane. Jagged temple gaps, a floated sleeve/cuff, detached-
  looking foot, and stacked bow tails remain.
- Rear: full hair coverage is preserved, but the center panel forms a deep
  pointed groove and the segmentation is too regular. No new C2 pixel change
  above 1% is measurable in this view.
- Three-quarter: the same helmet mass survives; the ordinary side gains a
  small local dent while temple gaps and disconnected sleeve construction
  remain.
- Mirrored three-quarter: the bilateral counterpart is visible as another
  small dimple, while the large smooth shell and isolated hanging slab remain.

## C2 versus frozen C0 and rejected C1

All comparisons use identical scene, cameras, lighting, render settings,
resolution, and byte-identical render specification. ImageMagick RMSE is the
normalized 0--1 value. The threshold count is the number of pixels whose
maximum RGB-channel difference exceeds 1%.

| View | C2 vs C0 RMSE | C2 vs C0 pixels >1% | Fraction | C2 vs C1 RMSE | C2 vs C1 pixels >1% | HD-01 result |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| Front | 0.0000686 | 1 | 0.0002% | 0.0000709 | 1 | No meaningful change |
| Side | 0.0007891 | 474 | 0.1157% | 0.0008263 | 561 | Local upper-rear dimple; rear wall survives |
| Rear | 0.0000682 | 0 | 0% | 0.0000749 | 0 | Coverage preserved; no reshaping |
| Three-quarter | 0.0004116 | 289 | 0.0706% | 0.0004435 | 323 | Small one-side dent; helmet survives |
| Mirrored three-quarter | 0.0008004 | 507 | 0.1238% | 0.0008005 | 507 | Bilateral dent; helmet survives |

The measurable change is larger than C1's ineffective one-sided stroke, but
it remains confined to tiny approximately 32--51 pixel patches. It does not
move the broad silhouette responsible for HD-01. Relative improvement is
therefore insufficient, and the paired dimples are a local surface regression.

## Technical evidence revealed after absolute review

- Edited object: `A68_BackCap_Sculpt`.
- Native operator: two `bpy.ops.sculpt.brush_stroke` calls, both `FINISHED`.
- Moved vertices: 35 per side.
- Maximum displacement: 0.008356 m per stroke.
- Mean displacement: 0.002599 m per stroke.
- Frozen non-target objects: 177.
- Frozen objects added, removed, or changed: none.
- Source cap remains recoverable: yes.
- Clean pinned reopen: pass, in object mode.

These checks prove a real bilateral native-sculpt edit and safe isolation; they
do not establish likeness. The locality explains why the render shows dimples
instead of a broad cap correction.

## Frozen inputs and tool identity

- C2 blend:
  `out/reimu_fumo_attempt_068_sculpt_coupon/c2_bilateral_grab/reimu_fumo_a68_c2.blend`
- C2 SHA-256 before and after render:
  `eac856bbdf2de7942bc2e41bb4ad92a58cf192963a7d141031fe254bf5f895e5`
- C0 SHA-256:
  `26c8613fe3eb17a1ddfcf7c8b596ed2aa264162b86d2b1e81acf7033d1fa75ba`
- Rejected C1 SHA-256:
  `67c9dbf7787749038ca168215647991f6b4df422f081c35a1b410852b9931557`
- Render-spec SHA-256:
  `d377222d84dd64aacaf7edb071f50929ecd880f781392c0fdb4060256659d1d8`
- Manifest SHA-256:
  `25f1ce270eb77c4d0788533696f61a8e85a71e36bf9eda2251439a0d6ac9c83b`
- Bazel target: `//projects/renders/cmd/fumo_review:render_packet`.
- Blender: `5.2.1 LTS`, build hash `9e2066aef7ef`.
- Render exit status: 0.

Equal before/after C2 hashes confirm that the packet renderer did not modify
the candidate.

## Pixel artifacts

Every image was inspected at its native 640 by 640 resolution.

| View | Camera | SHA-256 |
| --- | --- | --- |
| Front | `Review_front_Camera` | `809804bc452ce3d4247a7cb18c40ca9454da8d59d81cfb09e331f3027b87fff2` |
| Side | `Review_side_Camera` | `f3bf80a2412511f1174c408d71ddc7090d74b1242347976c880d8e76201a7b3d` |
| Rear | `Review_rear_Camera` | `1e0d8b3801a12c00e526da44a48fe058b122704e04c46946befb2a517ae48caa` |
| Three-quarter | `Review_three_quarter_Camera` | `565d338c70ed2e853789f4a9713395243dcd508661079a98e3b12d9d95eb4f17` |
| Mirrored three-quarter | `Review_three_quarter_mirror_Camera` | `3588940884227a13f70780cc45ee2364a77a1d1da36e687b49cad88cd8280d0d` |
