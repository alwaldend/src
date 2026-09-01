# A02 Sisyphus macro-mask plan

## Target

Test whether one coherent diagonal rock-face silhouette and a pale sky can
recover the controlling reference's macro composition without repeating the
rejected analytical slab/pile strategy.

## Frozen inputs

- `out/fumo_sisyphus_scene_scaffold/fumo_sisyphus_scene_scaffold_v01.blend`
  (`sha256:e482952bff46e3fea6b6d67b90ffc360bada6f45f280d4283f26db647c38e9d0`)
- `projects/renders/blender/fumo/fumo_sisyphus/references/sisyphus_reference.png`
  (`sha256:3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`)
- protected tracked Sisyphus Blend
  (`sha256:c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`)
- protected tracked Reimu Blend
  (`sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`)

## Bounded work

Copy V01 into task-local scratch. Preserve its boulder, neutral placeholder,
reference, camera, and lighting. Replace only the right/lower rock pile with
one low-frequency `ROCK_SLOPE_MACRO` form and set a flat pale-cream
`ENV_SKY`. Preserve or normalize the reusable collection names required by the
goal. Do not execute the unfinished V03 generator. Do not add fractures,
microtexture, extra rocks, procedural variants, or Fumo geometry.

## Outputs and gates

Produce the candidate Blend, a quote-free 512 px render, a reference/candidate
side-by-side, a silhouette overlay, a reference-at-side viewport, a clean-open
inventory, protected hashes, and a result note. Stop after the first render if
the sky luminance, 45-55 degree diagonal, 0.40-0.55W boulder, neutral
placeholder, or contact-shadow gate fails. A pass is only a macro checkpoint,
not a finished scene.
