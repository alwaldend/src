# Sisyphus A03 camera-space macro coupon plan

## Target

Test the intended single-owner Sisyphus composition cheaply in camera space
before spending another Blender iteration on terrain construction.

## Frozen inputs

- controlling reference
  `projects/renders/blender/fumo/fumo_sisyphus/references/sisyphus_reference.png`
  at `sha256:3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`
- rejected A02 render
  `out/fumo_sisyphus_attempt_002_macro_mask/render/quote_free_512.png`
  at `sha256:3452ebf95199fc76226113a2c93efcc6fe691bc4e25a22eb8ae328a0c7204f09`
- A02 boulder and neutral-placeholder camera-space silhouettes and framing

## Bounded work

Create exactly one non-promotable 512 by 529 raster/vector composition coupon.
Preserve the A02 boulder and neutral-placeholder silhouettes and their framing.
Replace every terrain pixel with one coherent right-half mask whose exposed
edge is 50 degrees, and replace the background with pale amber/cream sampled
from the controlling reference. Show the untouched reference beside the
coupon and make one silhouette overlay. Do not open or modify Blender, create
rock detail, tune lighting, introduce extra masses, or generate variants.

## Outputs and gates

Produce the coupon, reference/coupon side-by-side, silhouette overlay,
camera-space measurements, exact hashes, and a candid result. Pass only if the
single diagonal owns roughly the right half, the boulder remains 0.40-0.55W,
the left sky is visibly pale/open, the boulder visibly contacts the incline,
and the macro hierarchy is closer to the reference than A02. A miss resets the
strategy. A survivor authorizes one later Blender construction attempt but is
not scene or criterion acceptance.
