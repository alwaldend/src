# Rig037 one-curve coordinate preservation probe

Prepared, not run. Root reviews and executes this unsaved coupon. It tests
one legacy beveled decoration before any conversion of all 44 curves or the
85-object reusable copy. No saved model, whole-asset conversion, rig/action
change, reparenting, new software dependency, or delegation occurs here.

## Prior result and exact representative

The [sleeve receipt](rig_037_sleeve_review/receipt.json), SHA-256
`db62cc994834ff51bccbd00cea7636e7fd921b74e9faeffeb2b86ce284a3406a`,
records zero world drift, corner-normal error 1.69245e-7, maximum coordinate
RGB difference 0.00045836, and maximum actual-shader RGB difference
0.00093669. Both two-view comparisons have p99 RGB difference 0 and alpha
difference 0. These are saved 16-bit PNG comparisons decoded to scene-linear,
including view-transform/quantization effects; they are not float Render
Result equality. Root reported no visible change in the four shader images.
The worker does not claim a separate image review from that report.

Reuse the proven neutral-binding, normal capture, numeric-comparison, and
export-only material-copy methods. Do not transfer the mesh coupon's
constructive-companion coordinate hypothesis to a legacy curve.

- Exact input: `foot_033_candidate.blend`, SHA-256
  `98e92ee9a73ff49be32695dc06518ff885e5d91016278d16fb5a8771fd8fed48`.
- Object: `A42 Left white zigzag applique`, type CURVE.
- Existing material: `Bow white trim.002`, with explicit Generated links to
  two Noise and two Wave textures.
- Bevel depth: `0.001019999966956675` m; resolution 1, render resolution 0.
- Ordered stack: `A154 shared macro cage` targeting
  `A154 Left loop macro cage`, then `022 non-destructive rest proportion`
  targeting `022 bow proportion cage`.
- Intended existing bone: `ReimuFumoRig` / `Bow`.
- Keep current world placement and saved frame 1. Source data and all lattice
  dependencies remain untouched.

This wider zigzag supplies an easier-to-read beveled surface than the thin
root-fold stitches while retaining their Generated-material and lattice
construction pattern. The selection is not a claim that every other curve
has identical tessellation or shading behavior.

## Curve semantics established from pinned source

The legacy curve post-evaluation path converts beveled curves to mesh and
then applies remaining deformers directly. Unlike ordinary mesh evaluation,
this function does not maintain a parallel ORCO mesh. Pre-tessellation
modifier handling can also depend on Apply-on-Spline flags. Capturing the
actual evaluated surface includes that choice without reconstructing it.
See [pinned displist evaluation](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/displist.cc),
especially `curve_get_tessellate_point` and `curve_calc_modifiers_post`.

`mesh_copy_texture_space_from_curve_type` copies Curve texture-space values
and disables automatic mesh texture-space calculation. Curve bounds use
control points/handles, so recalculating bounds on the tessellation can
change Generated coordinates. The evaluated-curve conversion copies the
already evaluated mesh rather than entering the ordinary mesh all-data
reevaluation path. See
[pinned conversion](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/mesh_convert.cc).

Curve texture-space RNA reads ensure lazy calculation, and the curve
evaluation path synchronizes evaluated values to original data. See
[Curve implementation](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/blenkernel/intern/curve.cc)
and [Curve RNA](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/makesrna/intern/rna_curve.cc).
The helper requires the captured mesh's auto texture-space flag to be false
and its values to match the source Curve values. It never normalizes the
captured mesh's newly computed bounds.

For an evaluated curve-derived mesh, the helper reads `undeformed_co` from
the captured vertices. If preserved CD_ORCO exists, that getter returns its
unnormalized coordinate; otherwise it returns final vertex positions.
[Pinned Mesh RNA](https://raw.githubusercontent.com/blender/blender/v5.2.1/source/blender/makesrna/intern/rna_mesh.cc)
defines that distinction. The no-ORCO final-position branch is the expected
source-backed hypothesis for this legacy curve stack, not a silent fallback
to the sleeve's undeformed companion.

The candidate field is `0.5 * ((coordinate - location) / size + 1)`, using
the inherited Curve texture space. The renderer's original Generated
emission remains the actual comparison authority. Passing internal data
copy checks does not prove this field matches the engine's shader field.

## Frozen helper and interface

[rig_037_curve_coordinate_helper.py](rig_037_curve_coordinate_helper.py),
SHA-256
`6d906fe7b3dac7e2a4dfbe41bebc1b34d1f67f29374183b004e4d1d30b62e6da`.

The helper reuses the corrected, successfully tested
`rig_037_coordinate_helper.py` utility functions only after verifying its
exact SHA-256:
`0ed12641531786067eda7aff4b094ce0f774d2a1d186eccfaace973d0e932e8c`.
It does not call that module's sleeve entry or companion logic.

- `prepare_curve(prefix="Rig037_Curve")` returns `source`, `probe`,
  `collection`, `report`, and `baseline`.
- `make_coordinate_material(name, use_attribute=False)` creates an
  unassigned original-Generated emission material.
- The same function with `use_attribute=True` emits `rig037_generated`.
- The helper captures one evaluated curve surface, stores its coordinate
  field, inverse-neutral-binds it to Bow, restores captured corner normals,
  and assigns one final rigid armature modifier. The source is unchanged.
- It creates one export-only material copy, preserving slots and evaluated
  face indices while replacing only Generated/implicit coordinate inputs.
- Prefix collisions, unexpected modifier/material state, texture-space
  mismatch, topology drift, material-index drift, or failed numeric
  tolerances stop the probe and clean up only its newly created datablocks.

## Root's comparison sequence

1. Require the same numeric gates as the sleeve: world-position maximum
   error at most 1e-7 m; corner-normal vector maximum error at most 1e-4;
   coordinate attribute copying error at most 2e-7; exact topology, face
   material indices, and smooth flags.
2. Use two isolated 512 px curve views at identical source/probe placement,
   cameras, and rendering settings. Compare original Generated emission
   with the explicit attribute before adding procedural shader complexity.
3. Inspect `generated_component_bounds` and
   `generated_outside_zero_one`. Never clamp the stored coordinate field.
   Root's prepared `rig_037_curve_render_probe.py` applies the same fixed
   emission encoding `0.125 * Generated + 0.5` to source and probe, and
   requires component values inside plus/minus 3.5. This keeps the diagnostic
   values away from clipping. Record that encoding and the actual PNG
   precision in the receipt; the thresholds below are on encoded RGB.
4. After coordinate agreement, compare the original material with the
   export-only copy in the same views. Inspect bevel highlights, smooth
   normals, texture orientation, and bump detail. Preserve rendered files
   even if direct Render Result pixel access is unavailable.

Root retains the successful sleeve harness's numerical coordinate gates:
p99 RGB at most 0.001, maximum RGB at most 0.01, and maximum alpha at most
0.01, now explicitly on encoded RGB. Report actual
max/RMS/p99 and boundary metrics; do not equate this tolerance pass with
bit-identical field reconstruction. Root's visual review remains separate.

No scene action has been taken by this worker. Static AST parsing succeeded
and found no top-level calls; the new helper was neither imported nor run.
The curve-specific coordinate and full-shader equivalence fields remain
unverified until Root executes and reviews this coupon.
