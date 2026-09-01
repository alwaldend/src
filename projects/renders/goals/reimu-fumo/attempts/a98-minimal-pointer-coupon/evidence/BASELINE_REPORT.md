# A98 exact-snapshot baseline audit

## Verdict

The required input is intact and structurally suitable for one isolated
target-coordinate coupon. The saved file unambiguously identifies
`A94_Y_LiveGrab_Coupon` as a single-user, modifier-free, shape-key-free mesh in
Sculpt mode with fixed topology and an identity world matrix. A post-candidate
checker and exact target-coordinate sidecar are ready.

The saved interaction and render state must **not** be trusted without
foreground reconfiguration and receipts. The fixed right and three-quarter
images are nonblank and useful for silhouette comparison, but their clay
surface is severely clipped to white, so they are weak evidence for subtle
surface-form quality.

No foreground Blender, Xvfb, XTest, or save operation was used. All inspection
and rendering used repository-pinned Blender 5.2.1 LTS in background mode. The
source Blend rehashed unchanged after the audit.

## Exact subject

- Input:
  `out/reimu_fumo_attempt_094_real_pointer_sculpt_coupon/author/plus_y/snapshots/plus_y_before.blend`
- Input SHA-256:
  `02dd81b24a23a135462044c8b15a7498f743442f71d4de05ae21dae8ba9a1331`
- Scene: `Attempt41_Manual_Head_Maquette`
- Target: `A94_Y_LiveGrab_Coupon`
- Source: `A71_MacroClay`
- Inventory: 185 objects, 86 mesh objects, 14 camera objects

## Geometry and identity baseline

Target facts:

- collection: `A94_LiveProbe`
- vertices / edges / loops / polygons: `3863 / 7722 / 15444 / 3861`
- topology digest under the A98 schema:
  `b89e3f053ae95f37d10f39e79d2cae7caf79a5ce438f1f58c2328adc297a9813`
- local/world coordinate digest:
  `8893cfd6d8b5d96742fcef2f6ea8fc13750eb1418e40f9e9c6ff4197881a334b`
- world matrix: exact 4-by-4 identity
- mesh users: 1
- modifiers: none
- shape keys: none
- mirror axes: all disabled
- dynamic topology: disabled
- material: `A94_NeutralClay`

Source facts:

- vertices / edges / loops / polygons: `3863 / 7722 / 15444 / 3861`
- local/world coordinate digest:
  `020b8ef1526a768be9b5291b1c9d6e6c44f753e25f6cda4d8a4269d9a2753327`
- topology digest: same as the target
- world matrix: exact identity

Aggregate invariants:

- all non-target object fingerprint digest:
  `7dc74076317423dff8dd5128996a21c9723164f5028614434ada8e77de741812`
- all camera-object fingerprint digest:
  `4bf5153a90d74072038cf91547298dcb706ed499ceb99952b981121f210c019e`
- complete baseline fingerprint digest:
  `c8aa2b7e580ea7f05a717d535ebb9c334600bb42beba4f9e14fffae75f16d73d`

The complete per-object records, scene settings, saved screens, camera data,
matrices, mesh counts, coordinate hashes, topology hashes, material slots,
modifiers, and bounds are in `baseline_manifest.json`.

## Saved sculpt and view readiness

Observable state encoded in the exact file after a pinned clean reopen:

- active scene/screen/workspace: `Attempt41_Manual_Head_Maquette` / `Sculpting`
  / `Sculpting`
- active object and mode: `A94_Y_LiveGrab_Coupon`, `SCULPT`
- saved Sculpting viewport: orthographic right view, quaternion
  `(0.5, 0.5, 0.5, 0.5)`, view distance `0.2800000012 m`, solid shading,
  overlays enabled
- brush: `Grab`, strength `0.400000006`, smooth distance falloff, sphere
  falloff, front-face restriction disabled
- X/Y/Z symmetry disabled; dynamic topology disabled

Concrete concern: the cleanly reopened file stores `use_locked_size = VIEW`,
`size = 150 px`, and `unprojected_size = 0.25 m`. This is not the 0.05 m
unprojected radius reported by A94 immediately before the event. The A98
foreground author must explicitly reactivate/configure Grab, then record the
settled brush radius, lock mode, view projection, target identity, and operator
poll result. Merely opening this snapshot does not reproduce the required
brush contract.

## Fixed baseline renders and pixel inspection

- Right:
  `baseline_packet/right.png`, SHA-256
  `ba904d9dc578ea164eefee575f9325d24dbe371e07b850334d9a5bb93cede4b1`
- Three-quarter:
  `baseline_packet/three_quarter.png`, SHA-256
  `0ee8f19a624179b5e145a854b540505445efad49439340721397b94c55b6b330`
- Packet manifest:
  `baseline_packet/manifest.json`, SHA-256
  `6ab903c487721cf26be2b59de07d9f6c3ac5d78361836d24bbc4188ea0dd799c`

Both images are valid 512-by-512 8-bit RGBA PNGs rendered from the exact input
by pinned Blender. They were visually inspected. The target is fully framed
and its silhouette is clear, so the views can expose a macro Grab deformation.
However, the surface is nearly featureless white with only a dark lower rim.
Near-white pixels occupy 32.35% of the full right image and 37.51% of the full
three-quarter image; visually this covers almost the complete target surface.
Do not use these beauty renders alone to judge subtle deformation, volume, or
surface quality. A silhouette comparison remains admissible; a future form
review needs corrected neutral exposure or a solid-view capture.

## Machine-readable comparison contract

- Baseline manifest: `baseline_manifest.json`, SHA-256
  `fe7f7d6dca70165f487ef0c3fef12450c85266d564c3b0d07e060fd72114ca97`
- Target coordinate sidecar: `target_coordinates.f64le`, SHA-256
  `8893cfd6d8b5d96742fcef2f6ea8fc13750eb1418e40f9e9c6ff4197881a334b`
- Fingerprinter: `fingerprint_snapshot.py`, SHA-256
  `cb73e1971e35816b8d441a1e9f4a9027dc48d653455cf07d5761dba8b2c2943d`
- Post checker: `compare_candidate.py`, SHA-256
  `4a3ba1beec6d3ad7776b83ae56bd8b77c7b7cb20d2ab189836c282577502c71b`

The checker expects the post-candidate Blend to be loaded by pinned background
Blender and allows only target vertex-coordinate movement. It requires:

- exact Blender version and finite target coordinates;
- exact object inventory;
- exact target identity, collection, topology, world matrix, mesh users,
  modifiers, materials, shape keys, mirror state, and dynamic-topology state;
- exact source coordinates;
- exact aggregate non-target and camera fingerprints; and
- a nonzero target coordinate change.

It reports per-vertex displacement counts above 1 micrometre, 100 micrometres,
and 500 micrometres, maximum and mean displacement, mean displacement vector,
and the baseline-space support bounds above 100 micrometres. PASS is technical
isolation evidence only, never artistic acceptance.

The checker was run against the unchanged baseline as a negative control. It
correctly rejected it with exit code 3: every preservation gate passed and the
sole false gate was `target_vertex_coordinates_changed`; all 3,863 measured
displacements were exactly zero. Negative-control report: `self_check.json`,
SHA-256
`c2a0e4c71c734775938faa3a6d486ebe9e66c1980bbd08dc05d80e5be0405a26`.

## Post-candidate use

Run the checker in one pinned background Blender invocation with the immutable
candidate file loaded, `--disable-autoexec`, this `compare_candidate.py`, the
baseline manifest, coordinate sidecar, and a fresh report path under the A98
verification directory. Then render the post candidate through the same
`render_spec.json` and compare right and three-quarter silhouettes against this
packet. Rehash this report's input before accepting any result.

## Audit decision

The input baseline itself passes the requested identity/topology/readiness
audit with the brush-state and render-exposure qualifications above. It is safe
to proceed only if the foreground harness explicitly re-establishes and proves
the brush configuration. A post-candidate must pass `compare_candidate.py`
before its pixels are considered. The fixed images can reject a wrong macro
silhouette but cannot establish surface-form quality under their current
overexposed lighting.
