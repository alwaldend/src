# A79 brown short-nap hair material report

## Verdict

**Use the proposed shader only after the A79 hair geometry passes its neutral
silhouette and construction gate.** It is a materially better starting point
than `Hair brown clay`, but this swatch is preparatory evidence, not a likeness
pass and not authorization to integrate materials.

The main improvement is physical scale. The current shader's two bump
distances are `0.018 m` and `0.028 m` on a `0.25 m` plush. They are therefore
18--28 mm, far larger than the photographed sub-millimetre short nap. The
proposed shader uses `0.00010 m` and `0.00032 m` instead. It also replaces
per-object normalized `Generated` coordinates with metric `Object`
coordinates, so a small lock and a large crown no longer receive different
apparent fibre sizes.

## Controlling evidence

- Canonical color and front-surface response:
  `projects/renders/blender/fumo/reimu_fumo/references/`
  `canonical_front_25cm.png`, SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`.
- Physical edge, fold, and grazing-light response:
  `projects/renders/blender/fumo/reimu_fumo/references/physical_side.png`,
  SHA-256
  `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d`.
- Current material graph inventory:
  `out/reimu_fumo_attempt_044_render_pipeline/`
  `five_view_calibrated_host_01/review_state_actual.json`, SHA-256
  `13126b7d5d021a6bc63713afd5d79a98495a01bce93da094b5480ea9c44ec6c7`.
- Prior materials-only cycle:
  `out/reimu_fumo_attempt_042_material_module/apply_materials.py`, SHA-256
  `e27f9f0fadf28e28d12d63383911e491d30308d948bc647f35330e097bf4c355`.
  Its implementation-blind review rejected the complete plush-material gate
  at `4.5/10` intended-medium read and explicitly found a major wrong-material
  failure. The same review says geometry remained the dominant cause.

The canonical crown sample at crop `180x90+410+260` has mean display sRGB
`(0.266857, 0.195392, 0.171157)` and linear display value
`(0.058482, 0.031973, 0.025105)`. These values are image-space observations,
not raw material albedo; the proposed albedo was calibrated under the fixed
neutral review rig.

## Audit of `Hair brown clay`

| Property | Current value | Finding |
|---|---:|---|
| Mapping | `Generated`, normalized 0--1 per object | Fibre size changes with each object's bounds. |
| Macro noise | scale `4.2` | Object-relative rather than physical. |
| Nap noise | scale `48` | Object-relative rather than physical. |
| Micro noise | scale `175` | Object-relative rather than physical. |
| Nap bump | strength `0.10`, distance `0.028 m` | 28 mm is not short nap. |
| Micro bump | strength `0.13`, distance `0.018 m` | 18 mm contributes to clay/orange-peel shading. |
| Base linear RGB | `(0.045, 0.012, 0.005)` | Too red under the neutral review rig. |
| Color breakup | micro noise multiplied into base color | Risks pixel speckle rather than fibre response. |
| Roughness | `0.74`--`0.96` | Plausible range, but cannot rescue the scale errors. |
| Principled sheen | weight `0.32`, roughness `0.78` | Broad response is useful and should remain restrained. |

On the controlled normal swatch, the current material's mean sRGB is
`(0.441068, 0.247233, 0.171296)`, versus the reference crop's
`(0.266857, 0.195392, 0.171157)`. This is not a camera-aligned color proof,
but it exposes the strong red bias under the same neutral light used for both
swatches.

## Proposed exact shader

Material name in the swatch: `A79_PROPOSED_Brown_short_nap`.

### Mapping

- Use `Texture Coordinate.Object`, with object scale applied before rigging.
- Interpret coordinates in metres because the project model is built at
  physical scale.
- Use isotropic 3D noise; panel rotation does not change the response.
- For a heavily deforming final rig, replace the coordinate source with a
  metric-density UV or rest-position attribute while retaining the same
  physical frequencies. Object-space procedural texture can swim under large
  deformations. This is an integration check, not a geometry request.

### Base response

| Principled input | Value |
|---|---:|
| Base center, linear RGB | `(0.0220, 0.0100, 0.0065)` |
| Metallic | `0.0` |
| IOR | `1.46` |
| Specular IOR Level | `0.14` |
| Coat Weight | `0.0` |
| Diffuse Roughness | `0.76` |
| Sheen Weight | `0.30` |
| Sheen Roughness | `0.82` |
| Sheen Tint, linear RGB | base multiplied by `1.30` |

The base color is driven only by a macro ramp:

- Noise scale `68 m^-1`, detail `2.0`, roughness `0.52`, lacunarity `2.0`,
  distortion `0.08`.
- Ramp positions `0.20` and `0.80`.
- Dark ramp color `(0.0198, 0.0090, 0.00585, 1.0)`.
- Light ramp color `(0.02376, 0.0108, 0.00702, 1.0)`.

No micro noise drives base color. This avoids high-frequency colored speckle.

### Normal and roughness response

| Layer | Noise settings | Bump settings |
|---|---|---|
| Nap clumps | scale `430 m^-1`; detail `3.0`; roughness `0.58`; lacunarity `2.0`; distortion `0.06` | strength `0.30`; distance `0.00032 m` |
| Short nap | scale `1850 m^-1`; detail `2.0`; roughness `0.64`; lacunarity `2.0`; distortion `0.03` | strength `0.52`; distance `0.00010 m` |

Feed the short-nap Bump normal into the nap-clump Bump normal, then feed the
combined normal to Principled. Map the nap-clump noise from input
`0.18`--`0.82` to roughness `0.78`--`0.90`, clamped.

## Swatch result

The standalone `.blend` contains identical support meshes with current on the
left and proposed on the right. No fibres, displacement, curves, particles,
or geometry nodes were added.

- Swatch:
  `out/reimu_fumo_attempt_079_paired_hair_field/material/`
  `a79_hair_material_swatches.blend`, SHA-256
  `bc919fc9b2966bcc9c1a4bd4855e0cb09d64668604724f96cc7c76febeb32536`.
- Normal render:
  `renders/normal_current_left_proposed_right.png`, SHA-256
  `e25c6e127d798ca7563b8ae09e9b38c24b838e5ed73b0b601afb32aead164157`.
- Grazing render:
  `renders/grazing_current_left_proposed_right.png`, SHA-256
  `641aaced1bda877d2405fab98e3436c9b857c14debb6793b2741cd64c562e521`.
- Both renders are `900x600` PNGs made with Blender `5.2.1 LTS` and
  `BLENDER_EEVEE` under the fixed Attempt-42 light scale.
- Proposed normal-view mean sRGB is
  `(0.328331, 0.235488, 0.196482)`. Its simple RGB distance to the canonical
  crop is about 57% lower than the current material's under this swatch rig.

The proposed surface is darker, less red, and quieter. Its close and grazing
views show bounded short-nap breakup rather than the current centimetre-scale
bump. It still cannot create a fuzzy silhouette or correct a helmet-shaped
mesh. That limitation is expected because geometry-as-material was explicitly
excluded.

## Performance and integration trade-offs

- Current graph: 12 nodes, including three Noise textures and two Bump nodes.
- Proposed graph: 10 nodes, also with three Noise textures and two Bump nodes.
- The two combined-material `900x600` renders took `8.70 s` and `7.42 s`.
  These are packet timings, not isolated per-shader benchmarks.
- The proposal adds no polygons, curves, particles, displacement, external
  image textures, downloads, or runtime dependencies.
- Higher numeric Noise scale does not itself add geometry or samples; shader
  evaluation cost is comparable to the current graph.
- Normal-only nap will not affect the silhouette. If a final close-up demands
  visible flyaway fibres, that is a separate, explicitly authorized geometry
  decision and should not be smuggled into this material gate.

## Integration gate

Do not apply the shader while the neutral A79 hair field still fails macro
silhouette, rear coverage, paired-panel construction, thickness, contact, or
visual-veto checks. After the neutral geometry passes:

1. duplicate, rather than overwrite, `Hair brown clay` in the disposable A79
   candidate;
2. apply the exact graph above with scale-one hair objects;
3. render fixed front, both three-quarter, worse-profile, rear, normal close,
   and grazing close views;
4. reject visible colored speckle, plastic highlights, clay dimples, texture
   scale changes between panels, or animation swimming; and
5. accept only if the full model improves intended-medium read without hiding
   a geometry regression.

## Verification

Pinned Blender clean-reopen validation is `PASS` with no external image files:
`validation.json`, SHA-256
`f663f1b052faae7bdb4a47a04c548c3816ed9e14373140e7f9b8f10c83de99cc`.

Protected assets were unchanged after all swatch work:

- Rung 003 SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Tracked reusable asset SHA-256
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
