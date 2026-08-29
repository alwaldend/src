# A85 S12 final head-foundation pixel gate

## Scope and binding

This is a render-only review of the pinned 512 x 512 S12 front/profile packet
against every `HEAD_MEASUREMENTS.md` gate those silhouettes can measure. S08
is the direct survivor baseline. I used the established 55% beauty-render
luminance threshold and a 45--80% stability sweep, without inspecting Blender
topology or authoring intent.

- Blend SHA-256:
  `982da6404ea6edcbb4432903e67dad4ee5c130a203a5a5727a374b773fc9ad8a`
- Front SHA-256:
  `2d4a9f4c81780a71a4575a857ae144b250abb0e5043697f07c37af36ecf30d93`
- Profile SHA-256:
  `84943201026e9ce9e14752a06d30d9140a490104fca6f662b7c0190b6f8d5d7a`
- Front component-ID SHA-256:
  `b8de9edc0a861fe025b2fc15963060c2bef9e5fc4eec7b6f8c61d354971e794c`

At 55%, the S12 front bound is `382 x 373 px` and the profile bound is
`240 x 374 px`, so `Wh = 382 px`. The threshold sweep leaves all decisive
width, depth, and crown-shoulder rows unchanged; it changes anti-aliased
vertical extrema by at most two pixels.

## Full measurable gate

| Landmark | Target / gate | S08 | S12 | Verdict |
| --- | ---: | ---: | ---: | --- |
| Maximum front width | `1.000 Wh` | `1.000` | `382 px = 1.000` | Pass by normalization |
| Cushion height | target `.94`; `.89-1.00 Wh` | `.992` | `373 px = .976` | Pass; closer to target |
| Maximum profile depth | target `.62`; `.55-.68 Wh` | `.628` | `240 px = .628` | Pass; unchanged |
| Lower width at `Y=.06 Wh` | target `.82`; `.78-.87 Wh` | `.825` | `315 px = .825` | Pass; unchanged |
| Full cheek band `Y=.18-.43 Wh` | `.96-1.00 Wh` throughout | `.969-1.000` | `370-382 px = .969-1.000` | Pass; unchanged |
| Crown shoulder at `Y=.88 Wh` | `.65-.76 Wh` | `.738` | `282 px = .738` | Pass; unchanged |
| Crown rise above shoulder | about `.06`; at most `.10 Wh` | about `.110 +/- .003` | `36 px = .094 +/- .003` | **Pass; corrected** |
| Frontmost/rearmost projected planes | front at most `+.34 Wh`; rear provisional | about `+/-.314` | about `+/-.314` | Provisionally inside stop; datum unregistered |
| Central facial-zone depth variation | at most `.04 Wh` | unverified | unverified | Not measurable from these silhouettes |

The complete S12 cheek band was scanned row by row. Its minimum remains
`370 px = .969 Wh` at the lower edge near `Y=.18 Wh`; its maximum remains
`382 px = 1.000 Wh` near `Y=.43 Wh`. The crown correction therefore preserves
the S08 cheek pass without overshoot.

The crown shoulder remains at `y=97`, width `282 px = .738 Wh`, while the
front top moves from S08's `y=55` to S12's `y=61`. The shoulder-to-top rise is
therefore `36 px = .094 Wh`, reduced by about `.016 Wh` from S08 and now below
the `.10 Wh` stop. Total front height also remains safely inside its gate.

## Protected front and profile metrics

S12 preserves S08's sampled component-ID underside exactly:

- center `x=256` reaches `y=436`;
- `x=200/312` reaches `y=434`, two pixels higher;
- `x=180/332` reaches `y=432`, four pixels higher.

The shallow front arc remains `.005-.010 Wh`, and the `.825 Wh` lower width
is unchanged. Maximum profile depth remains `240 px = .628 Wh`. Profile depth
at `Y=.06 Wh` also remains `234 px = .613 Wh`, so S12 introduces no new
profile regression relative to S08, although it does not restore S05's
shallower `.560 Wh` lower-profile row.

The projected depth bounds remain centered at about `+/-.314 Wh` if the fixed
profile image center is treated as `Z=0`, provisionally inside the front-plane
stop. As in every prior gate, that datum is not independently registered, and
the binary silhouettes cannot verify central face-plane depth variation or
the hidden rear construction.

## Verdict

**S12 passes every measurable head-foundation silhouette gate and corrects
S08's crown-rise failure without regressing the protected front or maximum
profile metrics.** It is the strongest metric survivor and can replace S08 as
the measured foundation checkpoint. This is not an absolute likeness or
constructed-plush approval: facial-plane variation and hidden rear form remain
unverified, and the retained `.613 Wh` lower-profile row should remain visible
in subsequent construction review.
