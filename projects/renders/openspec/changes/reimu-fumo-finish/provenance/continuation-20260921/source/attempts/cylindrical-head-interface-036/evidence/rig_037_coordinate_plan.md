# Rig037 one-sleeve coordinate and normal probe

Prepared, not executed. The coordinator will review and run an unsaved
two-view sleeve coupon before any reusable-copy conversion. No action,
reparenting, additional rig, saved Blender file, or whole-asset bake is part
of this probe.

## Fixed input and representative

- File: `out/reimu_fumo_finish/desktop_astra/foot_033_candidate.blend`.
- SHA-256:
  `98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48`.
- Object: `Sleeve44P L front padded fabric panel`.
- Mesh: `Sleeve44P L front sewn panel Mesh.003`.
- Existing rig/bone: `ReimuFumoRig` / `Arm_L`; all weights exactly 1.
- Saved frame 1, identity object world matrix, no UV layer or shape keys.
- Exact order: Subsurf 1; Solidify at 0.00078 m and offset 0; Bevel;
  Armature; `022 body proportion cage`; `023 sleeve root L`.
- Materials, in order: `Sleeve44P outer warm cotton.002`, then
  `Sleeve44P inner warm fleece.002`. Both use Noise with unlinked Vector.
  Preserve complete modifier settings and material slots in the companion,
  including Solidify's material-dependent faces.

The nonlinear final lattices make this a useful test of freezing rest shape
without equating it to the original texture coordinates. The inspection
source is [rig_rest_bake_preflight.json](rig_rest_bake_preflight.json).

## Source-backed coordinate semantics

Blender builds ORCO along a parallel constructive-modifier path, skipping
deformation-only modifiers. Its final data-layer creation uses the original
object's texture space. This differs from carrying an arbitrary vertex
attribute through the mesh or normalizing the final deformed bounds.
[Pinned modifier evaluation](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/mesh_data_update.cc)
and [mesh texture-space implementation](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/mesh.cc)
establish these operations.

The public `MeshVertex.undeformed_co` getter reads internal `CD_ORCO`, then
returns `location + orco * size`. It therefore returns an unnormalized rest
position, not raw CD_ORCO. With no ORCO layer it silently returns current
position. The helper's normalization reverses that getter conversion before
mapping to the shader Generated convention; it does not normalize twice.
The texture-space RNA getters call `BKE_mesh_texspace_ensure`, so reading the
properties also resolves lazy automatic calculation on unchanged source
data. See the exact getter implementations in
[pinned RNA source](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/makesrna/intern/rna_mesh.cc).

The helper samples live evaluated ORCO before mesh conversion.
`preserve_all_data_layers` is not an ORCO guarantee: the conversion requests
`CD_MASK_MESH`, whose definition excludes CD_ORCO. See
[conversion](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/mesh_convert.cc)
and [mask definition](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/customdata.cc).
An evaluated-ORCO path is accepted only when topology/corner correspondence
matches and `undeformed_co` differs from current coordinates on this known
deformed sleeve. Otherwise exact capture is considered unavailable.

The candidate field is `0.5 * ((rest - location) / size + 1)`. It uses
texture-space properties, not recalculated baked bounds. This maps to the
Generated conversion in
[Cycles mesh export](https://raw.githubusercontent.com/blender/blender/v5.2.1/intern/cycles/blender/mesh.cpp)
and [texture-space helper](https://raw.githubusercontent.com/blender/blender/v5.2.1/intern/cycles/blender/util.h).
For the actual Eevee coupon, the engine's emitted Generated field remains
the comparison authority; cross-engine equivalence is not assumed.

For the sleeve's implicit vectors,
[Noise](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/nodes/shader/nodes/node_shader_tex_noise.cc)
uses the
[default texture-coordinate helper](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/nodes/shader/node_shader_util.cc),
which requests CD_ORCO for an unlinked coordinate input. Its replacement by
an explicit attribute still requires the coordinate and full-shader tests.

## Conditional companion fallback

If live evaluated ORCO is absent, the helper stops unless the caller sets
`allow_constructive_approximation=True`. That opt-in creates a temporary
source-object/data copy, retains Subsurf/Solidify/Bevel and both material
slots, and removes only Armature/Lattice from that copy. It samples the
constructive result using the source data's ensured texture space.

This approximates the source's parallel ORCO path without accessing internal
C structures. It is explicitly conditional because ordinary modifier
evaluation is not the internal `MOD_APPLY_ORCO` context. Vertex count,
polygon loop ranges, every corner vertex index, and face material indices
must match the captured real surface exactly. A mismatch stops the probe;
there is no nearest-neighbor transfer or alternative whole-asset fallback.
Even a matching companion remains an approximation until emission testing.

## Helper interface and bounds

[rig_037_coordinate_helper.py](rig_037_coordinate_helper.py), SHA-256
`e757680dcc64f464ef6b309679b115fb95a6cc35da6abfbb03354ae45be4d927`.

- Import performs no Blender or filesystem operations.
- `prepare_sleeve(prefix="Rig037_Sleeve",
  allow_constructive_approximation=False)` requires the exact source already
  open and returns `source`, `probe`, `collection`, `report`, and `baseline`.
- `make_coordinate_material(name, use_attribute=False)` creates an
  unassigned material emitting Blender Generated.
- The same call with `use_attribute=True` emits `rig037_generated`, a POINT
  vector attribute stored on the probe mesh.
- The caller owns render setup and temporary visibility/material overrides.
  The helper never opens, saves, renders, poses, changes frames, or changes
  source materials. It removes only its own intermediate objects/data.

The helper captures the complete evaluated surface, face material indices,
smooth flags, and corner normals. It inverse-neutral-binds the captured
positions to the existing Arm_L transform, transforms the normal preimage
with the corresponding transpose, and restores explicit custom normals.
The probe has one final armature modifier and one unit-weight bone group.
It copies only the two source materials, rewriting Generated links or
implicit Noise vectors to the explicit coordinate attribute on those copies.
The source geometry, rig, and material datablocks remain intact.

## Comparison order before any larger bake

1. **Numeric neutral check.** Require identical topology, material indices,
   and smooth flags. Maximum world-position error is at most 1e-7 m;
   maximum corner-normal vector error at most 1e-4; attribute copying error
   at most 2e-7. These test preservation of captured data, not equivalence
   between the candidate coordinate field and Blender's shader field.
2. **Coordinate emission.** Root prepares two isolated 512 px sleeve views
   at identical source/probe placement and camera settings. Render Blender
   Generated on the source and explicit-attribute emission on the probe,
   using the same engine/settings and linear floating-point output. Compare
   matching surface pixels and alpha coverage; report maximum/RMS RGB
   differences separately for interiors and antialiased boundaries. Freeze
   the harness's numerical acceptance bounds before its first run. A
   material-noise render cannot replace this field comparison.
3. **Actual shader and normals.** Only after the field test succeeds, compare
   the original source materials against the export-only copies in the same
   views, lights, engine, sampling, and color management. Inspect inner and
   outer faces, rims, highlights, and bump detail. Unchanged shader settings
   and passing numeric normals alone do not establish identical shading.

The report leaves `coordinate_shader_equivalence` and
`full_shading_equivalence` explicitly unverified. A failure remains local to
this sleeve. Preserve its metrics and images, identify coordinate,
correspondence, normal, or material-index cause, and revise only that cause.
Neither an accepted approximation nor a clean technical execution is a
whole-model visual, rig, or animation pass.

Static validation only: Python AST parsing succeeded and found no top-level
calls; whitespace checking produced no diagnostics. The helper was not
imported, executed, or tested in Blender by this worker.
