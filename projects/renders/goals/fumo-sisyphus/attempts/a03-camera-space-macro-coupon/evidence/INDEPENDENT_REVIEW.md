# A03 independent camera-space review

## Verdict

`SURVIVOR`, recorded as `refine`. The coupon satisfies every frozen macro
gate and is materially closer to the controlling reference than rejected A02.
This authorizes exactly one task-local, non-promotable Blender macro
construction attempt. It does not accept the scene or any goal criterion.

## Frozen-gate findings

- One continuous 50-degree right-side diagonal owns approximately 44.47% of
  the frame.
- The boulder spans 243 of 512 pixels, or 0.474609 frame width.
- The left sky is visibly pale amber and open rather than A02's dark brown.
- The boulder's lower-right silhouette visibly contacts the incline.
- The coupon removes A02's three competing round masses and restores the
  reference's sky, boulder, and single-incline hierarchy.

The strongest reset case is the terrain edge's rightward offset relative to
the reference and the visibly floating frozen placeholder. Neither defeats
this deliberately non-promotable coupon: the terrain remains within the
frozen “roughly right half” gate and the plan explicitly froze placeholder
placement.

## Claim corrections

- Terrain area, sampled sky RGB/luma, and the placeholder gap are descriptive
  metrics, not frozen thresholds.
- The reported `0.812387` value is weighted nonlinear sRGB luma, not
  linear-light relative luminance.
- The contact boolean uses a selected coordinate; inspected native pixels,
  not the boolean alone, establish visible contact.
- The generated polygon and pixels support one connected mask, although the
  stored boolean itself is hard-coded.
- The overlay alpha values differ by subject; calling it a uniform
  45-percent overlay would be inaccurate.
- A02 silhouettes were manually traced and spline-smoothed, not demonstrated
  pixel-exact copies. They are visually adequate for this macro gate.

## Narrow next work

Build one fixed-camera Blender blockout with one coherent incline, pale sky,
the frozen boulder scale, and no competing terrain masses. Reseat the neutral
placeholder onto that same incline by placement, never by adding a second
support mass. Stop after one 512 by 529 render. The measured 80.523844-pixel
value is a vertical gap at the placeholder center rather than minimum
geometric separation, but retaining that placement would visibly float the
character and must fail the next render.

## Exact subjects

- controlling reference SHA-256:
  `3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`
- rejected A02 render SHA-256:
  `3452ebf95199fc76226113a2c93efcc6fe691bc4e25a22eb8ae328a0c7204f09`
- coupon SHA-256:
  `81a0ac9128d6623b05d84038ef5f6ccc1d70f97ee1c9cada935c93b2109ad0e4`
- side-by-side SHA-256:
  `45b8dd4b15fd4b8d4b786754d7ae3efd17cb9d1f1197fb2e984ff5beb8101f0e`
- overlay SHA-256:
  `777515d9bae771aea6a3f53b57b62677c44b3753cab6c876585e6902d27e73c1`

The artifact manifest verified. Review was read-only; no file, goal, Git, or
Blender state was changed.
