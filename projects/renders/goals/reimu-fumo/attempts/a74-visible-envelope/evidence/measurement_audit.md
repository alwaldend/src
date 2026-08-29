# A74 measurement audit: visible head/hair envelope

## Verdict

**Conditionally usable, but only for visible-pixel constraints.** The A73
packet is accurate enough to drive an A74 visible head/hair coupon when it is
used as an outer silhouette, landmark, and layer-order dossier. It does not
measure a naked receiver, hidden seam, hidden support depth, or factory panel
pattern. The A74 builder should therefore begin with the visible brown field
and independent visible leaves/locks, then fit a subordinate support inside
them. It should not construct another analytic receiver from the A73 profile
table.

The strongest disconfirming evidence is the source itself: red bow cloth and
brown hair obscure every receiver boundary in the profile and rear frames.
The A73 prose mostly acknowledges this correctly. Any later claim that the
packet proves receiver-only `F(t)`, `R(t)`, a vertical receiver registration,
or a literal gusset layout would be unsupported.

## Replay and source-integrity audit

- Every one of the 19 source files still matches the SHA-256 recorded in
  `reference_measurements/source_hashes.csv`, including the canonical GIF,
  canonical front, eight controlling turn frames, three user stills/GIFs,
  and four extracted qualitative witnesses.
- I copied the exact A73 generator to an isolated directory and reran it
  without editing the program. The front CSV, rear CSV, same-crown
  counterexample, three overlays, and all-reference packet reproduced
  byte-for-byte.
- The profile overlay reproduced byte-for-byte and all 36 numerical/ownership
  rows reproduced semantically. The regenerated profile CSV is not
  byte-identical because the checked-in generator now emits the explicit
  `min_x_owner` and `max_x_owner` columns while the preserved A73 CSV predates
  those two columns. This is provenance drift, not a numeric disagreement.
  A74 should consume the regenerated schema or pin the artifact hash rather
  than assume the current script reproduces the old CSV bytes.
- Exact replay hashes for the principal deterministic outputs were:

  | Artifact | SHA-256 |
  | --- | --- |
  | front scanlines | `621d2d33e04a8be9b72ca40dcc9dcc11df4ba801301a49df7f5d2a1237fdff51` |
  | rear scanlines | `970ca4c08888599147a5ec1e3e42e66fa6bb6c2f22b67f441aefae26a987fa15` |
  | same-crown counterexample | `b88d1d674d911c671469cae347b60cac8b0c9f36d6bec44ca2e56f611a680172` |
  | canonical-front overlay | `989105117a2b57b659bd63826895b78ddd3cced6cc4318139f11d9a771995b04` |
  | profile overlay | `1d64553e3d3642aa9c5a96652516d3d37353c0dfdacd367e373a3314e9c71cd6` |
  | rear overlay | `89e6333a19a9b1122f6a62fd17a2c5742532ceff1ee07a7dd6c8b6c8ec8febcc` |
  | all-reference packet | `f4fea8cc35c5e99ed17e4a83262de384ec79bdb438692c8494fd7aab796860db` |

The isolated replay is under `out/a74_measurement_replay/`.

## Accuracy and semantic limits

### Canonical front

The arithmetic is correct and the six upper scanlines visibly land on the
brown outer contour. They are manually frozen semantic edges, not results of
the color-mask code. That is acceptable for a review contract, but their
uncertainty includes the human boundary decision as well as antialiasing.
Use them as outer projection targets, not as cross-sections of a solid.

The lower thresholded rows also align with visible brown pixels, but they are
unions of front field, cheek locks, and rear/nape locks. Their own CSV labels
correctly say they impose no receiver constraint. They are useful only as a
whole-module silhouette veto.

The normalization `Wh=368 +/- 4 px`, center `x=485`, crown `y=231` agrees with
the corrected canonical landmark dossier. `Wh=132 mm` is an inherited project
scale calibration: the script does not derive millimetres from the pixels.
It is consistent with the 25 cm plush and the separately established roughly
145 mm crown-to-lowest-lock height (`1.098 Wh`), but it must not be presented
as a new A73 pixel measurement.

### Turn profiles

The normalized outer-span arithmetic is correct. The cyan compact field and
magenta free rear leaf are also visually plausible component assignments in
frames 10/11 and 25/26. However, threshold components are separated by red bow
occlusion and the script chooses component indices manually. A component gap
is not empty 3D space, and the per-row `core_*` and `rear_leaf_*` endpoints do
not expose their hidden continuation.

Consequently, the profile rows support these claims only:

- the complete visible brown projection is asymmetric between the two sides;
- a compact visible brown field and a separately visible free rear leaf both
  contribute to the outer profile; and
- a candidate may not fill the full outer projection with one inflated mass.

The inherited `.77-.85 Wh` base-depth and `.36-.38 Wh` leaf-overhang bands are
not rederived by the A73 script. Treat them as provisional calibration bands
until a perspective-matched A74 overlay passes. Do not turn the raw component
endpoints into mesh coordinates.

### Rear

The rear scanline arithmetic and overlays are reproducible. They strongly
support a field that begins narrow under the bow, widens toward about `.70` of
its visible height, then narrows into unequal lower lobes. The reported
`center_px` is the center of a threshold-mask bounding box, not a physical
symmetry axis. Full width and the vertical location of maximum width are safe;
left/right half reaches are review evidence for asymmetry, not world-space
control points.

The `.45-.50 Wh` dominant-leaf width and `.10 Wh` offset quoted in the prose
come from an older semantic audit and are not independently rederived by the
A73 rear CSV. Preserve them as low-confidence review bands only, not as A74
construction dimensions.

## Minimal source-owned contract for A74

Use only the following high-leverage measurements to drive the first visible
head/hair coupon. More parameters would overfit segmentation noise and hidden
geometry.

### Front controller: canonical 25 cm image

- Scale datum: `Wh = 368 +/- 4 px = 1.000` and project scale `Wh = 132 mm`.
- Coupled brown crown-to-lowest-foreground-lock height:
  `1.098 +/- .015 Wh`.
- Outer full widths at `v = .052, .133, .242, .351, .459, .568`:
  `.408, .674, .837, .908, .948, .976 Wh`, with left/right asymmetry retained
  from the CSV. Match critical rows within `.03 Wh` and other extrema within
  `.05 Wh`.
- Visible beige exposure after hair occlusion: approximately
  `.603 x .603 Wh`, centered at `x=485`. This is not the hidden full opening.
- Identity-critical continuous fringe: clefts at `(u,v)=(.309,.372)` and
  `(.667,.372)`; low center tip `(u,v)=(.588,.677)`. The tip is `.101 Wh`
  viewer-right of the cleft midpoint. Do not create a centered V or separate
  padded bang cards.
- Foreground cheek-lock visible widths: viewer-left `.185 +/- .015 Wh` and
  viewer-right `.174 +/- .020 Wh`; keep them distinct from the crown field and
  from rear/nape lobes.

### Profile controllers: turn 10/11 and 25/26 brackets

- Register all four frames in one unchanged turn-camera scale (`Wg=244 +/- 5
  px`); do not fit one side and mirror it.
- Complete visible brown profile A: about `1.14 +/- .05 Wh`.
- Complete visible brown profile B: `1.19-1.23 Wh`, with about `.06 Wh`
  uncertainty.
- Structural gate: the compact crown/rear field and the long free rear leaf
  must remain visibly distinct contributions to the silhouette. The complete
  depth must not become one solid egg, helmet, or box.

These profile values are projected outer bounds, not support depth. A74 should
fit its visible surfaces by rendered overlay and keep hidden support
subordinate rather than solving a receiver from them.

### Rear controller: turn 18 with 17/19 brackets

- Outer brown width: `.94 +/- .04 Wh`.
- Crown-to-lowest visible rear lobe: `1.16 +/- .05 Wh`.
- Maximum width occurs around `.70 +/- .04` of rear height.
- Preserve the observed unequal lower lobes and long diagonal overlap. Do not
  infer a hidden top root or seam from the bow-occluded region.

## Recommended A74 acceptance test

Build only a visible-envelope coupon in frozen whole-plush context: one
continuous crown/fringe/temple field, two foreground cheek locks, a compact
visible rear/base field, and at least one independent free rear leaf. This is
visual ownership, not a required literal factory-piece count. Give each
visible piece thin padded-felt behavior; keep any support fully hidden.

Gate the front row widths, beige exposure, asymmetric fringe, and lock widths
first. Then require both side brackets and the rear bracket to pass without a
filled profile, bald gap, floating root, or exposed support. If the coupon can
match only by inventing hidden geometry outside these visible constraints,
return early: the representation or camera registration is wrong.
