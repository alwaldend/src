# A02 result: reject at first-render macro gate

## Work performed

Starting from the exact frozen V01 Blend, the bounded author removed the
known V01 incline and chip objects, created one low-frequency
`ROCK_SLOPE_MACRO`, replaced the sky material with a flat `ENV_SKY`, preserved
the declared boulder, neutral placeholder, reference, cameras, and lights,
saved a task-local candidate, clean-reopened it, and rendered one quote-free
512-pixel frame with pinned Blender 5.2.1 LTS.

## Result

**Reset; do not promote or refine this candidate.**

The first render fails the immutable stop gate. The diagonal measures
`41.3478` degrees instead of 45--55 degrees, sampled sky luminance is
`0.21187` instead of pale/bright, and three large round V01 masses remain in
the right/lower pile. Their actual names use `Fractured Mass`, while the
author's preflight assumed the unfinished V03 script's `Fractured Slab`
names. Consequently the intended single-owner experiment was not expressed,
and the new surface appears as a flat graphic wedge behind the old pile.

The boulder width passes at `0.47354 W`, lower-right terrain luminance passes
the readability floor at `0.34003`, the placeholder remains neutral, the
reference remains packed and render-disabled, clean reopen succeeds, and all
protected hashes are unchanged. Those technical passes do not override the
macro pixel failures.

## Criterion verdicts against revision 2

- `criterion-001`: **fail** -- boulder width passes, but the 41.35-degree
  diagonal misses the required interval and does not own the right half.
- `criterion-002`: **fail** -- sky is dark brown and the boulder contact-shadow
  path is not separately legible.
- `criterion-003`: **fail** -- three giant circular masses and a flat graphic
  wedge remain visible.
- `criterion-004`: **unverified** -- packed/render-disabled inventory passes,
  but the required reference-at-side viewport was intentionally not generated
  after the first-render stop.
- `criterion-005`: **pass for this exact rejected candidate** -- only the
  neutral future-Fumo placeholder is present; clean reopen, named modules, and
  all protected hashes pass.

## Artifacts

- `candidate/fumo_sisyphus_a02_macro_mask.blend`
- `render/quote_free_512.png`
- `evidence/author_receipt.json`
- `evidence/clean_open_inventory.json`
- `evidence/protected_hashes.json`
- `evidence/render_metrics.json`
- `evidence/FIRST_RENDER_PIXEL_VERDICT.md`
- `evidence/RUN_RECEIPT.md`

The planned side-by-side, silhouette overlay, and reference-at-side viewport
are absent by design: producing them after a decisive first-render failure
would violate the stop condition and add no decision value.

## Process finding

The next attempt must inspect the frozen Blend's clean-open inventory before
freezing deletion selectors; generator-source names are not evidence of saved
datablock names. It should compute the visible diagonal in camera pixel space,
not infer it from an unchecked world-space coefficient, and should calibrate
the sky using one cheap color-management probe before a full candidate render.
These are preflight corrections, not authorization for an A02 retry.
