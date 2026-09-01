# Sisyphus A04 single-owner Blender macro blockout plan

## Target

Express the A03 survivor composition once in pinned Blender, correcting A02's
unremoved terrain masses, dark sky, and unsupported neutral placeholder. This
is a macro construction stop gate, not a finished-rock or goal-acceptance
attempt.

## Frozen inputs

- rejected A02 candidate
  `out/fumo_sisyphus_attempt_002_macro_mask/candidate/fumo_sisyphus_a02_macro_mask.blend`
  at
  `sha256:c4b2b1118e8c215e0787703ad700dac1a665742cc9862315c6dadd48cccf5bd2`
- A03 coupon
  `out/fumo_sisyphus_attempt_003_camera_coupon/coupon.png`
  at
  `sha256:81a0ac9128d6623b05d84038ef5f6ccc1d70f97ee1c9cada935c93b2109ad0e4`
- controlling reference
  `projects/renders/blender/fumo/fumo_sisyphus/references/sisyphus_reference.png`
  at
  `sha256:3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`
- tracked Sisyphus and Reimu sources at their frozen SHA-256 values
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`
  and
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

## One bounded construction

Work on a task-local copy only. Preserve the exact fixed quote-free camera,
boulder, lights, packed/reference-at-side setup, render resolution, and
neutral placeholder geometry. Inventory `ROCKWORK` from the clean-open input,
then remove every mesh except `Sisyphus Boulder Scaffold`. Create exactly one
coherent, broad, solid `ROCK_SLOPE_MACRO` owner whose camera-space exposed edge
is the A03 50-degree diagonal and whose surface sits behind the boulder. Do not
create secondary support rocks, slabs, chips, or variants.

Set the rendered sky to the A03 pale amber target with an explicit,
deterministic color-management/material conversion rather than a render-driven
tuning loop. Move the existing placeholder and contact envelope together only
along world Z by one analytically computed translation so the placeholder's
evaluated bottom contacts the incline at its center X. Do not scale, reshape,
replace, or identity-style the placeholder.

Do not add fractures, texture, displacement, microdetail, new lighting,
materials beyond the one slope and sky owners, a Fumo/Reimu asset, or more
than one candidate/render. Use pinned Blender 5.2.1 through `bazel_agent`.

## Outputs and stop gate

Save one task-local candidate, clean-reopen it, verify exact protected hashes,
module inventory, absence of legacy terrain meshes, and the single computed
placeholder translation, then render one 512 by 529 quote-free fixed-camera
PNG. Produce a reference/candidate board and candid pixel verdict only after
the render.

Survive only if actual pixels show all of the following:

- one 45--55-degree lower-left-to-upper-right incline owns roughly the right
  half with no competing round masses or slab pile;
- the boulder remains 0.40--0.55 frame width and visibly contacts/occludes the
  incline;
- pale open sky remains on the left;
- the neutral placeholder visibly sits on the same incline rather than
  floating or being hidden by another support mass; and
- macro hierarchy is materially closer to the controlling reference and A03
  coupon than rejected A02.

On any miss, stop and reset without a correction render. A survivor authorizes
a later rock-construction module but is not criterion or scene acceptance.
