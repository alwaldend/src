# Flatness diagnosis of macro 064

Read-only visual advice to the coordinator, observed 2026-09-05 17:50 UTC.
The six supplied images were inspected before prior diagnostics or the writer
receipt. No builder or model was inspected. No model, source, goal, or Git
state was changed. The coordinator owns the construction decision.

## Finding

The model has substantial depth, but its front view does not communicate
enough constructed-fabric volume. The strongest actual sheet-like surfaces
are the bow wings and tails. The sleeves are already deep, yet their side
opening reads as a smooth round tube. The head is visibly rounded in side
view; increasing its whole depth is unsupported by these images.

Approximate hand-read pixels below use the current front close-cap width
`Wh ≈ 210 px` (x151–361). They describe the evidence, not a solved camera
alignment or a new acceptance target.

| Component | Image evidence | Meaning |
| --- | --- | --- |
| Head | Current side front/back envelope roughly x170–327, or .75 Wh; crown around y145 and lower cushion around y337. Frame 12 also shows a large curved front/back mass below its bow. | A globally shallow head is not the dominant failure. Confidence high. |
| Upper bow | Current front wings span x103–411, about 1.47 Wh. Their root is near x255,y143; outer edges are nearly straight. In side, the whole visible assembly largely follows a sloping band from x192,y90 to x344,y213. | The front surface has broad nearly uniform regions and weak root compression. Large depth range from tilt does not establish local folded volume. Confidence high for appearance, moderate for exact geometry. |
| Bow reference | Canonical wings span roughly x212–760 against head width ≈360 px, about 1.52 Wh. Root folds near x460–527,y212–256 visibly converge into a pinch. Physical front has strongly rolled root folds; physical side and frame 12 show turning edges and layered returns. | Preserve variant-appropriate outline; recover folds and the return into the knot. Do not combine the dramatically upright physical-side wing pose with the canonical front silhouette. |
| Sleeves | Current front shows broad triangular faces; side opening occupies roughly x214–288,y372–436 and is nearly circular. Physical-side/frame-12 sleeves are collapsed, oblique cloth openings with unequal lips and a continuous arm inside. | Depth exists. The manufactured form needs compression and seam-directed bends, not added depth. This is a secondary construction issue. |
| Garment/feet | Current frontal bodice, skirt, and lace read as stacked bands; side clearly exposes a flared skirt and projecting foot. Canonical front exposes more top surface of the thighs/hem. | Projection hides depth in front. The side also has a relatively rigid garment profile, but these unmatched views do not isolate a trustworthy numeric torso-depth correction. |
| Face/hair | Current face has a broad quiet center and cheek/chin shading; side proves curvature. Hair framing remains one smooth, closely coupled surface visually. | A front-facing face panel is expected on this plush. Small cloth overlap may help construction; making the face globular would threaten the reference. |

The level frontal view removes foreshortening and much top-surface visibility.
The photographs do not share that camera. Broad soft lighting reduces normal
contrast, but it cannot explain the bow's simple side edge or the sleeve's
round opening. No three-quarter candidate was supplied, so statements about
that view below are predictions. This diagnosis does not depend on fine pile.

## One bounded construction hypothesis

Revise only the bow assembly's **local folded shape**, keeping its canonical
front pattern and existing overall placement initially. Treat each wing as
two soft fabric faces meeting at a seam and compressed into the knot. Over
the root third, use two broad folds which converge and lose amplitude at the
knot. Let them open into a shallow cupped outer panel with a softly turning
edge. The tail should leave the same attachment through a local bend, then
hang as a thin panel; it should not become another uniformly inflated lobe.

This is a replacement of the fold/return field if the present representation
cannot express that construction, not another global tilt or depth scale.
The cuff and head observations above are not permission to patch them during
this test.

If initial amplitudes are needed, use approximately .03–.05 Wh of local wing
relief relative to its existing mean surface, narrowing toward seam and root.
Those 3.5–5.8 mm values are bounded prototype settings, **not** depths measured
from the photographs. Keep panel-edge thickness much smaller than that relief.
Avoid sharp corrugations, a balloon bow, and a dark disconnected root cavity.

- Front prediction: folds converge into the knot and create broad alternating
  surface turns while wing tips, span, zigzag path, and head overlap remain
  within the existing critical-landmark budget. Large outline shifts are a
  regression, not proof of added volume.
- Side prediction: the slanted bow retains its overall angle, but its edge
  has a real low-frequency return around the root and outer panel instead of
  one near-uniform band. No thick slab or new collision with the crown.
- Three-quarter prediction: one sees the front/back return and the pinch
  continuously; roots do not float, and the far wing is not a disconnected
  card. This must be rendered before acceptance.

## Disconfirming implementation evidence

Read only after the image diagnosis: `macro_064_writer_receipt.json` reports
head depth `.783 Wh`, sleeve aperture depth `.042 m` (`.361 Wh`), and bow
pattern depth `.08725 m` (`.751 Wh`). These support rejecting generic global
thickening. The last number is an assembly bounding extent; it cannot tell
whether each local cloth section is rolled or flat. `head_063_diagnosis.md`
likewise finds no support for a large global head-depth change and records
the front-camera projection confound. That earlier diagnosis is corroborating
evidence, not a measurement of candidate 064.

## Absolute provisional review

Same subject recognized: yes. Front/side-only scores: likeness 6, silhouette
6, construction 4, identity 6, contact/occlusion 5, plush-medium read 4,
presentation readability 7, all /10. Reject as an approval candidate.

Five largest visible discrepancies: insufficient bow root/edge folding;
round tubular sleeve construction; rigid garment with stacked frontal trim;
smoothly coupled hair/face construction and rounded lower hair ends; front
projection obscuring the spatial relation of bodice, hem, and feet. These
are ordered image judgments, not proof that all five need a geometry edit.

## Exact image inputs

SHA-256 values, paths relative to workspace:

- `out/reimu_fumo_finish/desktop_astra/macro_064_fast_review/front.png`: `3b2536076fe3a7446016e00dd3286d31fc1e30232820ee6c90671453277a7d16`
- `out/reimu_fumo_finish/desktop_astra/macro_064_fast_review/side.png`: `b2fdb369c27d3c4202611752115461d0c5a991c53f15dba65243e1efee219f0c`
- `projects/renders/assets/reimu_fumo/references/canonical_front_25cm.png`: `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`
- `projects/renders/assets/reimu_fumo/references/physical_front.png`: `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51`
- `projects/renders/assets/reimu_fumo/references/physical_side.png`: `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`
- `out/reimu_fumo_finish/desktop_astra/side_062_frame_12.png`: `4ebed9233acb351f2a8093496eb3e544771f1ee01e503cec84a1ce86f22b7365`
