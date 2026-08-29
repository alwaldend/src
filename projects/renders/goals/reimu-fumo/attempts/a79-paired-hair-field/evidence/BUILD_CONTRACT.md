# A79 Variant A Blender build and clean-reopen contract

## Status

This contract defines the only safe path from a **revised Variant A GO** to a
disposable Blender candidate. It does not authorize the current Variant A.
The selector currently says `REVISE`, the current pure preflight says
`build_authorized: false`, and therefore no Blender process or candidate save
is permitted yet.

Before construction, the coordinator must freeze a revised representation
specification, literal coordinate/preflight report, and selector verdict. All
three must say `PASS`/`GO`, and their exact hashes must replace the provisional
contract-input hashes below in a single `authorized_inputs.json`. A changed
file without a matching manifest update is a hard failure, not an implicit
approval.

## Frozen inputs at this publication

The two protected Blender inputs and the toolchain definition are immutable
for this attempt.

| Role | Workspace-relative path | SHA-256 |
| --- | --- | --- |
| protected parent | `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/reimu_fumo_working_rung_003.blend` | `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b` |
| tracked reusable asset | `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend` | `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76` |
| pinned Blender toolchain definition | `tools/blender/binary_toolchain.json` | `9512b68b50512cdcaabb1eeb12b665656258188b07def91f27c4ced46e193bb1` |
| A79 goal plan | `projects/renders/goals/reimu-fumo/attempts/a79-paired-hair-field/plan.md` | `258d1cd26a78d83d7cc232a5bba88544ad680445663f3ae31a6b3887aed949f2` |
| A78 interface inventory | `out/reimu_fumo_attempt_078_head_rest/interface_inventory/interface_inventory.json` | `43cf0cc3ab40737f5a3d7b7ec45c1da18dee489f7d0352faf983c5e8232f8d4c` |
| A79 retained-band baseline | `out/reimu_fumo_attempt_079_paired_hair_field/interface/retained_band_baseline.json` | `629f67fa36f6413efcea9f39c900cca30dc93f89af70b954d11530e2cee78374` |
| retained-band measurement source | `out/reimu_fumo_attempt_079_paired_hair_field/interface/measure_retained_bands.py` | `cbda4b4e747f2f9decb7a7f5bfdc011e9af259d2189e3e186b186b7b766df788` |

These pre-GO artifacts are exact audit snapshots but **not authorized**.
Variant work is active, so a later on-disk revision may already have another
hash. The builder must use the post-GO `authorized_inputs.json`, never infer
authority from one of these rejected snapshot hashes:

| Role | Workspace-relative path | SHA-256 | Current verdict |
| --- | --- | --- | --- |
| adversarial selector | `out/reimu_fumo_attempt_079_paired_hair_field/ADVERSARIAL_SELECTOR.md` | `b15701f7d5241d42ee4c3244ce0b261e1507a6d903860d437058416eef97b8b7` | `REVISE` |
| Variant A specification | `out/reimu_fumo_attempt_079_paired_hair_field/variant_a/TOPOLOGY_REPRESENTATION_SPEC.md` | `e51f1761580a9d7597cb66fd5e12afd60af913f1a7e17f1461040bd19279a52b` | `REVISE` |
| Variant A pure builder/preflight | `out/reimu_fumo_attempt_079_paired_hair_field/variant_a/builder.py` | `3ed65640107e709b04c16e239463111b18a5543c669ab0a0c890cfe036a3c638` | hard build stop |
| Variant A literal control nets | `out/reimu_fumo_attempt_079_paired_hair_field/variant_a/control_nets.json` | `c961eddd61833b8b47b2f1169e4adc25ffc3b5fed12b10f9a1a3f9510cf1e55f` | receiver-concentric rear shell; rejected |
| Variant A preflight report | `out/reimu_fumo_attempt_079_paired_hair_field/variant_a/preflight_scaffold.json` | `ac6d35918a28173449e3f9ea00c40a14a801693b8db9a4575c149e2de1fba3e8` | `build_authorized: false` |

The build must also pin every controlling image used by the revised pure
preflight. The current reference hashes are:

| Reference | SHA-256 |
| --- | --- |
| `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png` | `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c` |
| `projects/renders/blender/fumo/reimu_fumo/references/canonical_turn_180.gif` | `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30` |
| `projects/renders/blender/fumo/reimu_fumo/references/physical_front.png` | `f8c7d0f9911dbff1ef7f5d75601f9b10825015aecb367381971c076a5a3e7b51` |
| `projects/renders/blender/fumo/reimu_fumo/references/physical_side.png` | `cbb39e70f95fa464f6dc94862e0300d15771f3ff4c046d005849891aca55a19d` |
| `projects/renders/blender/fumo/reimu_fumo/references/clean_front.png` | `37813e03e04e4966f1dbe914e03a25a5f5ae561dcbf58b72677195c513ea48ca` |
| `projects/renders/blender/fumo/reimu_fumo/references/turn.gif` | `b42368e921bd055d73fbbb7bf65c2509a9aaf190cab02f89824b92b4cb75ece4` |
| `projects/renders/blender/fumo/reimu_fumo/references/sofa.gif` | `7c9173f91e6b6c801a1c77e50f9135e86fc89319f3c0262c10312320b1af8589` |

`authorized_inputs.json` must include the absolute workspace root, every path
and hash above, Blender version exactly `5.2.1`, a selector verdict exactly
`GO`, a pure-preflight verdict exactly `PASS`, and the revised candidate
object/data names. It must reject symlinks, path escape from this task root,
duplicate paths, and unexpected extra contract inputs.

## Exact mutation boundary

Exactly these six existing source objects may change, and only their
`hide_viewport` and `hide_render` booleans may change from `false` to `true`:

1. `A44 continuous hair cap with smooth opening`
2. `A42 Left asymmetric rear lock`
3. `A42 Off-center main rear lock`
4. `A42 Short right rear lock`
5. `A42 Main lock left seated seam`
6. `A42 Main lock right seated seam`

The receiver `Head_Cushion_Manual_Target` must remain present and must retain
its original visibility. The seven face witnesses, five front fringe/temple
panels, two cheek locks, `Subtle crown center seam`, bow, body, garments,
feet, materials, cameras, lights, worlds, animation, transforms, collection
membership, custom properties, and every other source ID are frozen.

Skipping the six objects in a generic object comparison is insufficient. For
each one, compare a full before/after payload and a second payload with only
the two allowed visibility fields removed. The normalized payload must match,
both allowed fields must make the exact `false -> true` transition, and the
set of full-payload differences must equal the six-name allowlist exactly.

## Process and append/copy workflow

1. Verify every authorized-input hash before importing `bpy`. Refuse a
   current selector other than `GO` or a pure report other than `PASS`.
2. Run only through the repository-pinned target in background mode, with
   factory startup, automatic file scripts disabled, and Python exceptions
   propagated:

   ```sh
   bazel_agent run //tools/blender:blender -- \
     --background \
     --factory-startup \
     --disable-autoexec \
     --python-exit-code 2 \
     --python /absolute/path/to/build_candidate.py \
     -- --authorized-inputs /absolute/path/to/authorized_inputs.json
   ```

3. Require `bpy.app.background` and version tuple `(5, 2, 1)`.
4. Start from an empty staging scene. Do **not** open the protected parent as
   the current main file and do not make a filesystem copy beside it. Append
   only scene `Attempt41_Manual_Head_Maquette` from the parent with
   `bpy.data.libraries.load(parent, link=False)`. Set the appended scene
   explicitly, remove only the task-created staging scene, update the view
   layer, and require `bpy.data.filepath` not to resolve to either protected
   input.
5. Record the complete source inventory and source fingerprints before adding
   any candidate data. Verify the six hide objects, receiver, retained roots,
   cameras, scene, and expected source materials all exist.
6. Create one fresh task collection beneath the appended scene. This single
   additive child collection and its declared new IDs are the only permitted
   collection-graph additions. No source object may be reparented or relinked.
7. Create fresh, single-user candidate mesh datablocks and objects. Copy a
   source hair material to a new candidate material if a neutral diagnostic
   material is needed; never edit or directly repurpose the source material.
8. Apply the six visibility changes only after all candidate geometry exists
   and all pre-save gates pass. Re-capture and validate the exact visibility
   delta before the save call.

Source-scene append is preferred over opening and overwriting a copy because
it makes the main-file identity unambiguous and removes any path by which a
failed save could target the parent. Any source material or reusable setting
must be copied into a new datablock before modification.

## Explicit mesh construction

The revised GO manifest owns the literal outer/inner coordinates,
connectivity, semantic face regions, root courses, and outer/inner
correspondences. The Blender builder may consume those values; it may not
invent, fit, smooth, resample, or repair them except for the exact bounded
hidden-root fitting operation authorized by the revised report.

The current `control_nets.json` is specifically ineligible: its rear field is
a receiver-concentric height-slice shell even though its topology is paired.
At `z = 160 mm`, for example, receiver, inner, and outer center Y are
approximately `51.399`, `51.927`, and `55.524 mm`. Density, refinement, clean
topology, or passing thickness cannot convert that construction mechanism into
an independently authored fabric panel.

For every selected Variant A component:

- construct an indexed `MeshSpec` first and validate it in pure Python;
- create a new mesh with `mesh.from_pydata(vertices, (), faces)`;
- call `mesh.validate(clean_customdata=False)` only as a detector; if Blender
  reports any repair, delete the just-created candidate data and fail;
- preserve explicit outer and inner face labels and bridge labels in a stable
  custom-data representation;
- require finite coordinates, exact declared counts, unique names, one mesh
  user, identity object transform, one closed face-connected component,
  Euler characteristic `2`, consistently positive winding/volume, and zero
  duplicate vertices, duplicate faces, degenerate faces, wire edges,
  boundary edges, overused edges, or self-intersections;
- reject Blender's automatic `.001` name suffixes, shape keys, constraints,
  animation, Geometry Nodes, Cloth, Shrinkwrap, Solidify, or unapplied
  Subdivision; and
- if the revised representation uses deterministic geometric refinement,
  freeze/apply it into the explicit paired mesh before any evaluated gate and
  before saving. The saved object must remain a closed paired pocket with no
  live modifier. Refinement must not create a new silhouette extremum or move
  any frozen-view silhouette by more than one pixel.

Smooth shading is allowed because it does not move geometry. It cannot be
used as evidence for silhouette density, paired thickness, or plush form.

## Source and candidate fingerprints

Reuse the A78 deterministic serializers only after adapting them to A79.
Capture stable payloads for every original source object, scene setting, mesh,
curve, material, image, action, node tree, camera, light, world, collection,
and other source datablock. Include transforms, parent/constraint/modifier
state, material slots, custom properties, animation references, mesh topology
and coordinates, curve controls, and node/socket values. Exclude volatile
runtime fields, user counts, previews, evaluated caches, and backreferences.

Maintain four separate records:

1. `source_before`: exact appended-source inventory, full object payloads,
   visibility-normalized object payloads, scene settings, source collection
   graph, and source datablock payloads;
2. `source_after_allowed_hide`: the same records after the exact six-object
   visibility transition;
3. `candidate_before_save`: candidate object and mesh payloads plus candidate
   collection/material payloads after the dependency graph has settled; and
4. `protected_file_hashes`: parent, tracked asset, and all frozen contract
   inputs before and after save and again after clean reopen.

Adding the declared candidate collection necessarily adds new IDs to the file.
Do not weaken source fingerprints to accommodate that. Compare every original
ID by its stable payload, compare the original collection graph separately,
and allow exactly one additive task collection containing exactly the declared
candidate objects/data. Any other added, missing, renamed, or changed ID is a
failure.

## Fatal gates before any `.blend` save

All reports must contain raw samples and summary statistics; a Boolean alone
is not evidence.

### Representation and camera gates

- Revised pure report `PASS` and independent selector `GO` at the exact
  manifest hashes.
- At 512 px, every evaluated visible boundary span is at most `3 px`, or at
  most `2 px` at crown extrema, temple turns, nape-lobe tips, leaf tips, and
  overlap corners.
- Maximum normal chord error is at most `1 px` at those identity-critical
  points and at most `2 px` elsewhere.
- Both paired skins have positive signed Jacobians, zero foldovers, quad
  angles in `25..155 degrees`, aspect ratio at most `4:1`, and adjacent-face
  area ratio at most `4:1`.
- No extraordinary/refinement vertex lies on a critical silhouette or root;
  elsewhere its one-ring normal spread and incident area satisfy the revised
  selector's `20 degrees` and `25%` limits.
- Receiver fitting is restricted to the declared hidden root/perimeter bands.
  It may not move or generate a non-root inner point and may never move or
  generate a visible outer point.
- Run an area-weighted non-concentricity test on every component and on the
  union of their non-root visible outer surfaces. For each sample, record the
  nearest receiver point, outward receiver normal, and signed normal offset;
  fit the single constant receiver-normal offset minimizing weighted squared
  residual. Reject the representation if `80%` or more of any tested
  non-root visible outer area lies within `+/-1.0 mm` of that best-fit offset.
  Record the fitted offset, area fraction, residual quantiles, and sample
  rows. This is a fatal representation gate even when density, refinement,
  silhouette, topology, thickness, and roots otherwise pass.

### Applied-refinement paired-pocket gates

- Validate the final explicit, applied/evaluated mesh rather than only its
  control cage.
- Area-weighted declared outer-to-inner thickness is `2.0..5.0 mm` over at
  least `95%` of visible area. Every sample has positive local orientation.
  Record `min/P05/P50/P95/max` per component.
- Correspondence segments do not cross either skin; bridge faces do not cross
  either skin; bridge roll/width stays inside the revised component contract;
  and no connected out-of-band patch projects wider than one unfiltered pixel
  in any fixed review view.
- Each component remains a single closed two-manifold, positively wound
  pocket with explicit bridge continuity and zero self-intersections after
  refinement.

### Root and crossing gates

- Sample every declared root uniformly by evaluated arc length, not by vertex
  count or unrestricted nearest-point minima.
- Every root sample is positive, both endpoints pass, at least `80%` of each
  root's arc length lies inside its component-specific signed band, and the
  longest contiguous out-of-band run is at most `0.57 mm` and never projects
  wider than one unfiltered pixel in any review view.
- Record signed `min/P05/P50/P95/max`, total arc length, passing arc length,
  every contiguous failing run, closest triangle, closest point, outward
  normal, signed clearance, and any bounded pre/post fit coordinate.
- The front/crown and rear-base roots are checked against the unchanged
  receiver **and** the relevant retained fringe/temple composite. The leaf
  root is checked against the final rear-base outer surface.
- Candidate/receiver, candidate/retained-root, and leaf/rear-base triangle
  crossings are zero. Use BVH only as a broad phase and run an exact
  triangle/triangle narrow phase; an AABB/BVH candidate-pair count is not an
  intersection result.
- The baseline retained temple, cheek-lock, and crown-seam crossing counts and
  pair hashes do not increase or drift.

### Multi-view ownership and semantic gates

Before saving, rasterize the exact evaluated geometry through the frozen
camera matrices into unfiltered 512 px object-ID/depth masks. This is an
in-memory mechanical diagnostic, not review evidence or a beauty render.

- The front scanline widths and face exposure remain inside the frozen A79
  bands; retained fringe owns the opening silhouette.
- Compact field and independent leaf pass the selector's per-side profile
  partitions; the leaf owns farthest depth for a contiguous height of at least
  `.25 Wh`, while the rear base remains separately visible.
- Rear brown union width is `.90..98 Wh`, crown-to-lowest-lobe height is
  `1.11..1.21 Wh`, and maximum width occurs at `.66..74` of rear height.
- The front/crown-to-rear-base depth step is `1.5..3.5 mm` and `2..6 px` in
  both three-quarters and both profiles.
- Front/crown, rear-base, and leaf IDs each own a nonzero contiguous region.
  Union-only or shared-material masks cannot pass this gate.
- Within the registered target hair region, the rear mask has zero receiver,
  face-applique, or background pixels; no gap, leak, overshoot, or disconnected
  island exceeds one pixel.

Any failed gate blocks `bpy.ops.wm.save_as_mainfile` completely.

## Guarded save and failure behavior

- Use a new, unique candidate filename under this build directory. Refuse to
  overwrite any existing file and refuse a resolved path outside
  `out/reimu_fumo_attempt_079_paired_hair_field/build/` or equal to either
  protected input.
- Set `bpy.context.preferences.filepaths.save_version = 0` so no `.blend1`
  backup is created.
- Write a JSON failure report for a pre-save failure, but do not call any save
  operator and do not create a candidate `.blend`.
- After all pre-save gates pass, settle the dependency graph, capture final
  candidate/source fingerprints, set the exact scene, and save once to the
  unique candidate path.
- Immediately rehash the parent, tracked asset, toolchain/contract inputs, and
  references. A mismatch quarantines the saved candidate and forbids render
  or promotion.
- A process or I/O failure during the save cannot be made transactional by
  Blender. Therefore never reuse the filename. A partial or post-save-failed
  file remains quarantined evidence with `render_allowed: false`; do not open,
  overwrite, promote, or silently delete it.
- The successful candidate is promoted by a manifest pointer after clean
  reopen, not by overwriting or moving a protected asset.

## Independent clean reopen

Clean reopen is a new pinned Blender process, not a second function call in
the builder process:

```sh
bazel_agent run //tools/blender:blender -- \
  --background \
  --factory-startup \
  --disable-autoexec \
  --python-exit-code 2 \
  /absolute/path/to/the_unique_candidate.blend \
  --python /absolute/path/to/validate_reopen.py \
  -- --manifest /absolute/path/to/build_report.json
```

The verifier must require that `bpy.data.filepath` is the exact candidate and
that its SHA-256 equals the build report. It must then repeat, from the reopened
evaluated scene:

- Blender identity and every protected/frozen input hash;
- original-ID inventory, exact six-object visibility delta, normalized source
  fingerprints, source collection graph, animation, cameras, lights, and
  missing-library/external-resource checks;
- candidate object/data/material/collection fingerprints, exact names, object
  counts, topology, transforms, data users, modifier absence, face-region
  labels, finite coordinates, winding, volume, and self-intersection checks;
- applied paired-pocket thickness/orientation and connected failure-patch
  checks;
- the area-weighted non-concentricity test against the unchanged receiver,
  again on exact reopened evaluated geometry;
- arc-length root statistics, longest contiguous failure runs, exact crossing
  checks, and baseline retained-crossing equivalence; and
- the same unfiltered fixed-camera object-ID/depth semantic masks.

The reopen report records all raw samples, mask hashes, failures, and
`pass`. Only a passing reopen may set `mechanical_preflight_pass: true` and
`render_allowed: true`. Renderers receive only the exact candidate hash from
that passing report. The later five-view packet is produced in a separate
process and still requires pixel inspection; clean reopen is not visual
acceptance.

## Reusable implementation from A77/A78

The following A78 functions are useful foundations after A79-specific review:

- `sha256_file`, `_normalise`, `digest_payload`, `custom_properties`, and the
  object/mesh/curve/node/action/datablock serializers;
- `source_inventory`, `capture_source_state`, and `compare_source_state`;
- `append_protected_scene`;
- `MeshSpec`, `topology_report`, `require_valid_topology`,
  `create_explicit_mesh`, and `spec_from_mesh_object`;
- `evaluated_mesh_data`, `evaluated_bvh`, and `resample_polyline`; and
- the guarded save/manifest and separate `--validate-reopen` structure.

A77's protected-file hashes, append pattern, material-copy helper, evaluated
mesh extraction, topology summary, and post-save rehash are also reusable.

The following must **not** be copied unchanged:

- A78's `HIDDEN_SOURCE_OBJECT` and `assert_allowed_hide_delta` model permit one
  hidden receiver and skip its full object digest. A79 must keep the receiver
  visible and prove exactly six visibility-only changes.
- A78's `signed_coverage_report` accepts aggregate 80% coverage without a
  contiguous-run limit. A79 needs arc-length weighting, passing endpoints,
  all-positive samples, and the `0.57 mm`/one-pixel longest-run gate.
- A78's paired-point or unrestricted nearest-distance logic is not enough for
  thickness. A79 needs declared outer/inner correspondence, area weighting,
  positive orientation, bridge checks, and connected failure-patch masks on
  the applied mesh.
- `BVHTree.overlap` alone must not be described as exact triangle
  intersection evidence; use it as broad-phase input to a deterministic exact
  narrow phase.
- A77's Subdivision/Solidify helpers, ruled mantle/leaf construction,
  one-minimum contact summaries, and pre-gate save behavior are specifically
  forbidden for this attempt.

## Known implementation risks

1. Linking a candidate collection changes the scene's aggregate collection
   graph. Treat the exact additive child as a separate allowlist; do not hide
   it by weakening every scene/collection fingerprint.
2. Blender user counts, evaluated caches, normals, previews, and dependency-
   graph flags are volatile. Exclude those runtime fields while preserving all
   authored values and relationships.
3. Blender silently suffixes colliding names. Reject any name or datablock name
   that differs from the frozen manifest.
4. `mesh.validate()` can repair malformed input. Any reported repair is a
   fatal mismatch, not a successful cleanup.
5. Smooth shading can conceal face shading but cannot repair a coarse
   silhouette. Pixel-span and chord-error gates must inspect evaluated
   geometry before shading credit.
6. A dense paired mesh can still be a receiver-offset helmet. Topology and
   thickness gates do not test authorship of form; the non-concentricity gate
   is independently fatal.
7. Nearest-surface signs are meaningful only after receiver and candidate
   winding/manifold checks pass. Record the exact triangle and normal for every
   signed sample.
8. Material reuse can mutate source user counts or source nodes. Copy a new
   neutral material and fingerprint both source and copy.
9. A save can succeed while the wrong scene or stale file is active. Set the
   scene explicitly, record `bpy.data.filepath`, hash the exact output, and
   require a separate-process reopen.
10. An in-process semantic rasterizer is a mechanical veto only. It cannot
   replace the post-reopen fixed-camera packet or implementation-blind pixel
   review.
