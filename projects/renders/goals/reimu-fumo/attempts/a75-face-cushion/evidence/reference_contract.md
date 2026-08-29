# A75 visible beige face-cushion contract

## Scope and verdict

This contract covers only the **visible stuffed beige face surface** and the
surface region that carries the eyes and mouth. It does not define the brown
hair, rear hair, bow, body, or a hidden complete head receiver.

The strongest source conclusion is negative but important: no supplied view
shows a naked beige cushion, a beige crown, a beige rear, a perimeter seam, or
the complete side wall. The canonical front supplies a strict exposed-pixel
shape; the turn supplies a strict occlusion/roundness behavior. Any full
receiver dimensions remain implementation hypotheses and must stay hidden
under the hair.

## Source authority

| ID | Source | Face-cushion role | Limitation |
| --- | --- | --- | --- |
| S1 | `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png`, SHA-256 `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` | Exact variant, 25 cm scale, exposed beige bounds, lower cheek/chin contour, graphics | No depth or hidden edge evidence |
| S2 | all 30 frames of `projects/renders/blender/fumo/reimu_fumo/references/canonical_turn_180.gif`, SHA-256 `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30` | Front-to-profile exposure collapse, shallow face read, cheek/underside continuity, occlusion order | Perspective turn; about one frame / 12 degrees phase uncertainty |
| S3 | `projects/renders/blender/fumo/reimu_fumo/references/clean_front.png`, SHA-256 `37813e03e04e4966f1dbe914e03a25a5f5ae561dcbf58b72677195c513ea48ca` | Clean eye, lid, and mouth cross-check | Suppresses pile, compression, and contact |
| S4 | `projects/renders/blender/fumo/reimu_fumo/references/physical_front.png`, SHA-256 `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51` | Stuffing, nap, lower-cheek fullness, nearly flush applique construction | Different cloth state; its wider opening does not override S1 |
| S5 | `projects/renders/blender/fumo/reimu_fumo/references/physical_side.png`, SHA-256 `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d` | Edge roll, cap/face seating, shallow non-plaque read | Oblique and heavily occluded; not a true depth controller |
| S6 | `projects/renders/blender/fumo/reimu_fumo/references/turn.gif`, SHA-256 `b42368e921bd055d73fbbb7bf65c2509a9aaf190cab02f89824b92b4cb75ece4` | Low-resolution continuity cross-check | Different pose and too small for precise dimensions |
| S7 | `projects/renders/blender/fumo/reimu_fumo/references/sofa.gif`, SHA-256 `7c9173f91e6b6c801a1c77e50f9135e86fc89319f3c0262c10312320b1af8589` | Face pile, soft compression, nearly flush graphics | Motion blur, zoom, crop, and perspective exclude dimensions |

S1 wins front dimensions and variant identity. S2 wins profile behavior and
layer order. S3--S7 may veto a plastic, hard, or implausible construction but
may not alter an S1 landmark.

## Coordinate system

Use canonical outer front hair width, excluding the bow, as
`Wh = 368 +/- 4 px = 1.000 = 132.0 mm`. The complete plush is 250 mm tall.

- canonical head center: `x = 485 px`;
- canonical brown crown datum: `y = 231 px`;
- centered horizontal coordinate: `x_wh = (pixel_x - 485) / 368`;
- downward coordinate: `v = (pixel_y - 231) / 368`.

The crown datum registers the coupled head/hair image. It is not a visible
beige-cushion top.

## Directly observed front surface

### Exposed-pixel envelope

The canonical beige exposure has the axis-aligned witness

```text
pixel bounds: x = 374..596, y = 358..580
Wh bounds:    x = -0.3016..+0.3016, v = 0.3451..0.9484
extent:       0.6033 x 0.6033 Wh = 79.6 x 79.6 mm
```

This is an **exposed-pixel envelope only**. Its top and sides are largely made
by hair occlusion; it is not a cut pattern, a hidden face-panel outline, or a
complete cushion silhouette.

The hair-owned upper boundary is the asymmetric sequence below. It is listed
because it controls which beige pixels remain visible, not because the beige
cushion owns this contour:

```text
(-0.326, .563) -> (-0.261, .380) -> (-0.103, .535)
 -> (+0.087, .677) -> (+0.149, .394) -> (+0.323, .554)
 -> (+0.378, .500)
```

### Strict visible cheek-to-chin arc

The source-observed open lower contour is:

| Pixel `(x,y)` | Centered `(x_wh,v)` |
| --- | --- |
| `(374,471)` | `(-.3016,.6522)` |
| `(375,500)` | `(-.2989,.7310)` |
| `(382,530)` | `(-.2799,.8125)` |
| `(397,553)` | `(-.2391,.8750)` |
| `(421,569)` | `(-.1739,.9185)` |
| `(450,578)` | `(-.0951,.9429)` |
| `(485,580)` | `(.0000,.9484)` |
| `(521,578)` | `(+.0978,.9429)` |
| `(550,569)` | `(+.1766,.9185)` |
| `(575,553)` | `(+.2446,.8750)` |
| `(590,530)` | `(+.2853,.8125)` |
| `(596,501)` | `(+.3016,.7337)` |
| `(596,474)` | `(+.3016,.6603)` |

Interpretation:

- maximum visible fullness lies in the lower-middle cheek, not at the crown;
- the side stays almost vertical only from about `v=.65` to `.73`;
- the lower turn begins by `v=.73-.81`, then tightens continuously;
- the chin is a broad, shallow arc: at `v=.943`, about `.19 Wh` of width is
  still visible, while the center is only `.005 Wh` lower;
- there is no flat bottom, square corner, jaw, muzzle, or separate chin ball.

Use `+/- .03 Wh` for the registered visible arc. Soft pile below about
`.01-.02 Wh` is not a separate form target.

### Observed asymmetry

The lower cheek/chin support is nearly bilateral: corresponding right-hand
points are only `1-2 px` (`.003-.005 Wh`) farther out than the left. Do not add
a large anatomical asymmetry. The **exposure** is nevertheless asymmetric
because the fringe and locks are asymmetric. Preserve low-frequency stuffing
variation and the off-center hair occlusion; do not mirror the complete visible
mask or add random surface noise.

## Profile depth and roundness

### What S2 directly shows

Across canonical frames 00--11 the beige region contracts to a narrow front
crescent and then disappears behind hair. A conservative color-component audit
gives the following screen-space exposure sequence, normalized by the turn
front width `Wg ~= 244 px`:

| Frame | Visible lower-face component width | Use |
| ---: | ---: | --- |
| 08 | about `.34 Wg` | late three-quarter witness |
| 09 | about `.26 Wg` | transition witness |
| 10 | about `.18 Wg` | near-profile witness |
| 11 | about `.11 Wg` | closest profile exposure |
| 12 | no reliable beige component | fully hair-occluded |

Treat each value as `+/- .05` because color segmentation, perspective, pile,
and yaw dominate. The values are a regression trend, not orthographic depth.
The opposite half of the turn is not a mirror: frames 27--29 must be compared
directly.

The visible surface itself reads as:

- a broad, mildly convex front carrier;
- quietest/most planar through the eye-and-mouth field;
- increasingly rounded through the lower cheek and underside;
- continuous under the chin, with no plate edge or vertical side wall; and
- completely subordinate to brown hair at crown, temples, sides, and rear.

A passing profile therefore shows a small beige front crescent at frames
10--11, not a side-on oval, disk, mask, or half-sphere.

### Coupled-depth evidence that does **not** belong to the beige surface

The following values are useful only as host-head containment context:

- canonical frame 10 compact coupled core projection: about `.74 Wh`;
- older turn side occluded cushion lower bound: about `.678 Wh`;
- physical oblique-side inner head projection: about `.810 +/- .10 Wh`.

Every value includes brown hair, an occluded support, perspective, or more than
one layer. None measures a beige-only cushion depth. Prior design bands such as
hidden support width `.90-.95 Wh`, height `.84-.90 Wh`, and depth
`.66-.82 Wh` are reversible construction guesses, not A75 likeness gates.

If A75 needs a complete hidden support, initialize it inside those broad bands,
keep it invisible outside the S1 opening, and judge it only through the coupled
front/three-quarter/profile pixels. Do not expose or score a beige crown, side
wall, underside seam, or rear plane.

## Eye and mouth carrier

The face surface must remain a low-curvature textile carrier over the complete
graphic field. It may bulge gently as one stuffed plane, but it may not crease,
ridge, or change normal abruptly under an applique.

### Canonical red eye-fill observations

| Eye | Pixel bounds | Normalized bounds | Component center |
| --- | --- | --- | --- |
| viewer-left | `x=384..448, y=450..507` | `x=-.274..-.101, v=.595..750` | `(-.180,.671)` |
| viewer-right | `x=526..590, y=450..508` | `x=+.111..+.285, v=.595..753` | `(+.206,.673)` |

These are red-fill measurements, not the complete applique. S3/S4 support a
complete sleepy-eye applique around `.21 Wh` wide by `.30-.32 Wh` high, with
its carrier center near `x=-.18/+ .18, v=.61`. Preserve the straight half-lid,
beige upper field, rounded red lower field, dark textile border, and restrained
catchlight. The canonical rightward bias is small; do not replace it with
perfect radial symmetry or enlarge the eye spacing to escape the fringe.

The canonical mouth is an almost centered stitch:

```text
pixel bounds: x=481..492, y=549..551
center:       x_wh=+.004, v=.867
extent:       .033 x .008 Wh (about 4.3 x 1.1 mm)
```

S3/S4 cross-check the mouth at `v=.84-.87`. It is a short dash, not lips, a
mouth cavity, or a raised bead.

Eye and mouth layers must conform to the evaluated face. Eye projection is at
most `.015 Wh` (2.0 mm), preferably `.003-.008 Wh`; the mouth should be still
shallower. No face perimeter rim may be introduced to seat the graphics.

## Fabric and stuffing cues

S4, S5, and S7 consistently show:

- short soft nap/fuzz with directional, low-amplitude color variation;
- one broad stuffed face plane, not porcelain, rubber, or a taut balloon;
- low-frequency cheek fullness and broad contact compression;
- a softly rolled lower edge/underside rather than a bevel;
- hair sitting close over the forehead/temples without a parallel plaque rim;
- nearly flush layered textile graphics and a stitched mouth; and
- broad compression toward the collar/torso, with no exposed human neck.

Geometry must establish the shallow plane and cheek/chin roll before fuzz or
textures are added. Random displacement, high-frequency wrinkles, uniform
subdivision, or a procedural noise normal cannot substitute for stuffing form.

## A75 acceptance and vetoes

Before any hair integration is accepted:

1. Register S1 at `Wh`, center, and crown; exposure must be
   `.603 +/- .03 Wh` wide and high.
2. The strict cheek/chin arc must remain within `.03 Wh` without a flat chin,
   square corner, or cheek sphere.
3. Eye and mouth carriers must retain the measured positions and read as one
   calm stuffed plane.
4. Fixed turn-matched renders must reproduce the progressive exposure collapse
   through frames 08--12 and the direct opposite-side return through 27--29.
5. No beige pixel may appear on the crown, rear, or outside the hair-framed
   opening; no face-panel rim may appear in side or three-quarter views.
6. The untextured form must already read as a shallow stuffed face with a broad
   cheek/chin roll. A sphere, egg, mattress, face disk, raised mask, vertical
   side wall, or anatomical skull is an immediate reject.
7. Physical and sofa references must veto plastic edges, hard tangencies,
   raised pill eyes, random lumpiness, and missing collar compression.

Hidden receiver topology, seam placement, and exact rear closure remain open
implementation choices. They are not allowed to alter these visible gates.
