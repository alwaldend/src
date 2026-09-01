# A02 first-render pixel verdict

## Exact subject

- Candidate: `candidate/fumo_sisyphus_a02_macro_mask.blend`
- Candidate SHA-256:
  `c4b2b1118e8c215e0787703ad700dac1a665742cc9862315c6dadd48cccf5bd2`
- First render: `render/quote_free_512.png`
- Render SHA-256:
  `3452ebf95199fc76226113a2c93efcc6fe691bc4e25a22eb8ae328a0c7204f09`
- Controlling reference SHA-256:
  `3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`
- Review context: the exact 512 by 529 PNG was inspected at original
  resolution before consulting implementation diagnostics.

## Gate verdict

**Reject and stop at the first-render gate.**

The boulder width passes at `0.47354 W`, the neutral placeholder remains
present, and the lower-right terrain is readable. The candidate nevertheless
fails three mandatory macro conditions:

1. The nominal diagonal is `41.3478` degrees, below the required 45--55
   degrees.
2. The sampled open-sky mean luminance is `0.21187`, below the provisional
   `0.50` first-frame gate; it reads dark brown rather than pale cream.
3. Three large V01 round masses remain visible and obscure the intended single
   coherent face. The new macro surface therefore reads as an inserted flat
   triangular wedge, not the sole right/lower rock-face owner.

There is no clean, separately readable boulder-to-incline contact-shadow
path. The strongest visible cast shadow belongs to the neutral placeholder,
while the boulder contact area is obscured by the surviving central mass.

## Absolute image review

- Unlabeled same-composition recognition: no
- Intended-medium read: simplified Blender rocks, not a convincing
  photographic or stylized weathered rock face
- Five largest discrepancies, in order:
  1. dark-brown negative space instead of pale amber/cream sky;
  2. three giant round masses dominate the right/lower half;
  3. the new diagonal is too shallow;
  4. the diagonal surface reads as a flat graphic wedge;
  5. boulder contact and cast-shadow hierarchy are not legible.
- Overall reference likeness: `3/10`
- Macro silhouette and proportions: `4/10`
- Rock-face construction: `2/10`
- Contact and occlusion: `3/10`
- Intended-medium read: `4/10`
- Presentation readability: `5/10`
- Major visible failure present: yes
- Absolute decision: reject

No side-by-side, overlay, or reference-at-side diagnostic was generated after
this verdict because the immutable plan requires an immediate stop when any
first-render macro gate fails.
