# A73 macro head/hair reference measurements

## Verdict for the builder

The references support an explicit **outer hair-envelope profile**, but they do
not expose a naked receiver at any height.  Use the measurements below in two
different ways:

1. the canonical-front scanlines are hard outer projection vetoes; and
2. the turn/rear scanlines reserve depth and width for separately visible hair
   layers.

Do not relabel either set as measured receiver rings.  The A73 receiver's
`W(t)`, `F(t)`, and `R(t)` remain a construction hypothesis inside those
envelopes.

The most important actionable result is that the visible crown widens much
more slowly than an isolated full-width cushion.  If the A73 receiver uses the
current profile table, its crown cannot share the visible hair-crown datum.
Moreover, no pure vertical offset makes the current `1.000 Wh`-wide table fit
all sampled front rows within their recorded uncertainty.  The best sampled
fit occurs near a `.141 Wh` hidden-crown inset but still exceeds a controlling
row by `.0285 Wh` beyond tolerance; a global width scale of about `.971` or an
equivalent local narrowing is required for the sampled outer veto.  This is a
candidate/interface fit, not a photographed seam or a source-derived receiver
dimension.  Leaving the registration undefined makes an integrated
protrusion or mattress read very likely.

## Normalization and source ownership

- Canonical front: `Wh = 368 +/- 4 px`, center `x = 485 px`, brown crown
  `y = 231 px`.  It alone controls exact front variant and 25 cm scale.
- Canonical turn: fixed `498 x 498 px` canvas; use the unchanged turn-local
  frontal width `Wg = 244 +/- 5 px`.  Frames `10/11` own profile A,
  `25/26` own profile B, and frame `18` owns the rear with `17/19` as yaw
  brackets.  Nominal phase uncertainty is one frame, about `12 degrees`.
- Clean and physical fronts support fringe shape, stuffing, and panel
  construction only; they are different cloth states and do not override the
  canonical front.
- Physical side supports a shallow receiver followed by separately draping
  hair, plus thin padded-felt edge behavior.  Its oblique perspective carries
  roughly `.10 Wh` depth uncertainty and is not an orthographic ruler.
- The older turn and sofa GIFs qualitatively veto rigid cards, thick plates,
  bald gaps, and spherical stuffing.  Their lower resolution and moving crop
  do not control A73 dimensions.

All relevant sources are visibly attached in
[`reference_measurements/all_relevant_reference_packet.png`](reference_measurements/all_relevant_reference_packet.png).
Their exact paths and hashes are in
[`reference_measurements/source_hashes.csv`](reference_measurements/source_hashes.csv).

## Canonical-front per-height half widths

`v` is distance down from the observed brown crown in `Wh`.  Left and right
values are independent reaches from `x = 485`; the asymmetry is retained.

| `v` | left reach | right reach | full width | visible owner | uncertainty |
| ---: | ---: | ---: | ---: | --- | ---: |
| `.052` | `.182` | `.226` | `.408` | continuous crown/front-hair field | `.011 Wh` |
| `.133` | `.321` | `.353` | `.674` | continuous crown/front-hair field | `.011 Wh` |
| `.242` | `.402` | `.435` | `.837` | continuous crown/front-hair field | `.011 Wh` |
| `.351` | `.438` | `.470` | `.908` | continuous crown/front-hair field | `.011 Wh` |
| `.459` | `.465` | `.484` | `.948` | continuous crown/front-hair field | `.011 Wh` |
| `.568` | `.484` | `.492` | `.976` | continuous crown/front-hair field | `.011 Wh` |

The exact source pixels and lower lock-inclusive rows are in
[`reference_measurements/front_outer_scanlines.csv`](reference_measurements/front_outer_scanlines.csv),
and the lines are visible in
[`reference_measurements/canonical_front_profile_overlay.png`](reference_measurements/canonical_front_profile_overlay.png).

Rows below about `v=.60` are no longer one cushion witness.  The front hair
field, cheek locks, and rear/nape locks overlap in projection.  Their union is
about `.93-.99 Wh` wide through `v=.68-.95`, but this is a **visible hair-layer
constraint only**.  It does not authorize a `.99 Wh` receiver through that
whole band.

### Current A73 table against the front envelope

With no vertical inset, linear evaluation of the current A73 half-width table
would project approximately:

| outer-crown `v` | A73 full width | observed outer width | result |
| ---: | ---: | ---: | --- |
| `.052` | `.578` | `.408` | protrudes by `.170 Wh` |
| `.133` | `.852` | `.674` | protrudes by `.178 Wh` |
| `.242` | `.961` | `.837` | protrudes by `.124 Wh` |
| `.351` | `.994` | `.908` | protrudes by `.086 Wh` |
| `.459` | `1.000` | `.948` | protrudes by `.052 Wh` |
| `.568` | `.986` | `.976` | within about `.010 Wh` |

This does not by itself reject the receiver representation, because its hidden
crown datum is not observed.  It does reject **implicit same-crown
registration** and the claim that an offset alone has already solved the fit.
The reproducible offset sweep is
[`reference_measurements/a73_front_fit_offset_sweep.csv`](reference_measurements/a73_front_fit_offset_sweep.csv),
with per-row values in
[`reference_measurements/a73_front_fit_detail.csv`](reference_measurements/a73_front_fit_detail.csv).
Add the receiver-to-hair vertical offset as an explicit interface parameter,
narrow the violating rows rather than widening the source-owned hair, and
render the later integrated front overlay.

## Canonical profile reach by height

The raw measured profile rows are in
[`reference_measurements/profile_outer_scanlines.csv`](reference_measurements/profile_outer_scanlines.csv).
The color overlay in
[`reference_measurements/turn_profile_component_overlay.png`](reference_measurements/turn_profile_component_overlay.png)
keeps the compact visible brown field cyan and the independent free rear leaf
magenta.  Red bow cloth often occludes the space between them; an apparent
segmentation gap is not physical empty space.

The table below summarizes outer front-to-rear projected depth at the most
reliable middle rows.  `v` uses the turn-local crown row and `Wg`; each source
frame remains separate rather than being averaged into a fictional symmetric
profile.

| `v` | frame 10 | frame 11 | frame 25 | frame 26 |
| ---: | ---: | ---: | ---: | ---: |
| `.307` | `.930` | `.934` | `1.033` | `.984` |
| `.410` | `.975` | `.988` | `1.074` | `1.033` |
| `.512` | `1.008` | `1.016` | `1.102` | `1.070` |
| `.615` | `1.041` | `1.057` | `1.123` | `1.094` |
| `.717` | `.980` | `.980` | `1.135` | `1.061` |
| `.820` | `.836` | `.943` | `1.164` | `1.090` |

Each row has about `3 px` segmentation uncertainty before the `+/-12 degree`
phase error.  The complete-profile bounds remain the more conservative frozen
targets: approximately `1.14 +/- .05 Wh` for profile A and `1.19-1.23 Wh`
with about `.06 Wh` uncertainty for profile B.

These depths visibly belong to two layers:

- a shallow continuous crown/rear base, whose frozen whole-form depth band is
  `.77-.85 Wh`; and
- a separate broad crown-rooted rear leaf extending about `.36-.38 Wh` beyond
  the base where visible.

The profile rows therefore **veto** assigning `1.0-1.2 Wh` to the receiver.
They do not reveal a per-height split between the receiver's front and rear
reaches because the red bow tail and outer hair cover that interface.  A73's
near-flat front reach, high rear maximum, and `.065 Wh` lower-rear undercut are
falsifiable design choices, not direct pixel measurements.

## Canonical rear half widths and visible layer

Frame `18` is the controller.  `17/19` are retained only as phase brackets.
For frame `18`, with a source-local envelope center of `x=275 px`, the outer
combined rear field measures:

| fraction down rear height | left reach | right reach | full width |
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

The underlying source coordinates, all three frames, and uncertainties are in
[`reference_measurements/rear_outer_scanlines.csv`](reference_measurements/rear_outer_scanlines.csv).
The aligned visual evidence is
[`reference_measurements/rear_scanline_overlay.png`](reference_measurements/rear_scanline_overlay.png).

This directly supports a rear envelope that starts narrow under the bow,
reaches maximum width around `.70 +/- .04` of its height, and narrows through
unequal lower lobes.  It also confirms the frozen combined bounds of about
`.94 +/- .04 Wh` width and `1.16 +/- .05 Wh` crown-to-lowest-lobe height.
The dominant foreground leaf is about `.45-.50 Wh` wide and displaced roughly
`.10 +/- .05 Wh` in the rear image; the simple outer scanlines do not encode
its long diagonal overlap edge, so preserve that edge as a separate later
hair-layer gate.

## What A73 may and may not claim

### Defensible constraints

- `Wh = 132 mm` for the eventual maximum canonical front projection.
- Receiver total depth inside `.77-.85 Wh`, leaving about `.35 Wh` for
  independent rear hair rather than filling the outer profile.
- A front projection that fits the six canonical upper scanlines after an
  explicit hidden crown offset is chosen.
- No receiver projection may widen or deepen the tracked outer hair envelope.
- The later rear hair must own the `.70`-height maximum width, unequal lobes,
  off-center broad leaf, and profile asymmetry.

### Not observed

- The naked receiver crown, underside, seam, or exact height.
- Per-height receiver-only `F(t)` and `R(t)` under the red bow and brown hair.
- A literal two-piece factory pattern or a source-visible receiver seam.
- True orthographic depth from the physical side photograph.

Treating those unknowns as measurements would repeat the session's core
failure: a clean procedural surface justified by numbers that actually belong
to occluding hair.  Keep A73 reversible and reject it at P0 if its isolated
clay read remains an egg, mattress, rounded box, or human head.

## Artifact manifest

- Reproducible measurement script:
  [`reference_measurements/measure_reference.py`](reference_measurements/measure_reference.py)
- Exact front rows:
  [`reference_measurements/front_outer_scanlines.csv`](reference_measurements/front_outer_scanlines.csv)
- A73-to-front interface fit and offset sweep:
  [`reference_measurements/a73_front_fit_detail.csv`](reference_measurements/a73_front_fit_detail.csv)
  and
  [`reference_measurements/a73_front_fit_offset_sweep.csv`](reference_measurements/a73_front_fit_offset_sweep.csv)
- Exact profile rows and component ownership:
  [`reference_measurements/profile_outer_scanlines.csv`](reference_measurements/profile_outer_scanlines.csv)
- Exact rear rows and brackets:
  [`reference_measurements/rear_outer_scanlines.csv`](reference_measurements/rear_outer_scanlines.csv)
- Visual overlays and all-reference packet: the four PNGs linked above.
