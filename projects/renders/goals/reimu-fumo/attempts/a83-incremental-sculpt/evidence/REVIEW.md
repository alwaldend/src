# C17 constructed panel sleeve review

Status: **NON_CANDIDATE — owner-ID gate passed; final-material gate failed;
categorical UNDO / stop**.

C17 started from exact C1b SHA-256
`d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`.
It replaced only the camera-right sleeve on a task-owned copy. The old 22
`Sleeve44P R` objects remain preserved and hidden. All other C1b object
fingerprints remained unchanged.

## Owner-ID gate

The owner-ID pair proves the intended ownership change:

- the magenta front pocket owns the broad front silhouette;
- the cyan rear pocket remains independently visible as a narrow edge in the
  controlling three-quarter view;
- the lime cuff is continuous, with no circular tunnel or background hole;
- the short insert remains hidden; and
- every evaluated new-mesh versus torso/skirt BVH pair reports zero
  intersections.

The representation is genuinely two separate closed cloth-panel pockets with
a 1.4 mm inner gap. It has no axial loft and no closed tubular bridge. The
front/rear geometry remained fingerprint-identical when materials changed.

## Final-material fixed-pixel verdict

The material pair is nevertheless a categorical appearance failure:

- Front: the sleeve is broad enough to avoid C14's narrow limb, but it reads
  as a stiff, flat trapezoidal board. Its shoulder contact is visually weak,
  so the module appears pinned beside the torso rather than gathered into it.
  The canonical sleeve is a short, soft bell with visible cloth collapse and
  a fuller shoulder transition.
- Three-quarter (controlling): the separate panels collapse into one thin
  paddle once they share the white material. There is too little changing
  curvature or seam shadow to communicate stuffed depth. The cuff edge still
  reads as a rigid tube-like slot rather than a rounded compressed fold, and
  the panel remains less draped and less substantial than both the canonical
  turn and physical-side fabric reference.
- Contact: exact collision cleanliness is retained, but the visible root gap
  and rigid tangent transition still make the sleeve look detached.

Internal scores for this isolated module are: silhouette/proportion `5/10`,
constructed-fabric read `3/10`, stuffed-volume read `2/10`, cuff/opening read
`3/10`, contact/occlusion `6/10`, and presentation readability `6/10`.

Decision: **UNDO C17 and do not promote**. The owner-ID gate succeeded, but
the final result is not “clearly close,” so the authorized one-step contour
correction was not used. The failure is volumetric and attachment-related,
not a single contour error. A credible next module would need a cloth-pouch
surface with non-planar course profiles, purposeful compression at the root,
and a cuff whose fold owns actual depth; it should not merely widen or rotate
this paddle.

## Exact evidence

- pristine C1b copy, SHA-256
  `d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`
- owner-ID blend, SHA-256
  `07d043ac5014a9304fb716bcbea66e581940ff8ed638acdf0937b879958b00b4`
- final-material diagnostic blend, SHA-256
  `296f3c8fa4fb07f32d4b8166ab655bc7a8a055a19613480520b5a0ce2a5ef210`
- owner-ID front / three-quarter, SHA-256
  `37be957aabc44a028b432cc9c7ec6486c4bd02080c69590b7aee1fa2808851fe`
  / `26ab7cf117b075b905b75fe523da7e92363acc51b6bb6626bfc9d35bfcfd3f53`
- owner-ID contact sheet, SHA-256
  `283a19b6722f6de9baf9acbf618bfeeb63e11f9f89f9349dcbab111ebfe675c5`
- final front / three-quarter, SHA-256
  `1b4a6f67eaac396a788a85c12220aab3f6e33c273bdb8b2d077483167733be1e`
  / `0f1c68841ed965ccead1a57c8313963f25824c054d3c1eaf472532fa2024f1f6`
- final reference/baseline/owner/final contact sheet, SHA-256
  `92a391bf211d873bbbbd63527c47b0b080cb3a171778ee215e7a32346b50f9ce`

The exact commands were:

```sh
cp --reflink=auto \
  out/reimu_fumo_attempt_083_incremental_sculpt/live_author/a83_C1b_coupled_cap_receiver_narrow.blend \
  out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/C17_C1b_pristine_copy.blend
bazel_agent run //tools/blender:blender -- \
  -b /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/C17_C1b_pristine_copy.blend \
  --python /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/build_owner_id.py
bazel_agent run //tools/blender:blender -- \
  -b /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/C17_panel_sleeve_OWNER_ID.blend \
  --python /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/render_owner_id.py
bazel_agent run //tools/blender:blender -- \
  -b /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/C17_panel_sleeve_OWNER_ID.blend \
  --python /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/apply_final_materials.py
bazel_agent run //tools/blender:blender -- \
  -b /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/C17_panel_sleeve_final_materials.blend \
  --python /var/home/simeonwarrenbot/.t3/worktrees/src/t3code-1040a9fb/out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/render_final.py
```

All four images are 512 px fixed views rendered after clean reopen with
repository-pinned Blender `5.2.1 LTS`. The source, pristine copy, owner-ID
blend, protected rung, and tracked reusable asset remained byte-exact.
Promotion was never authorized.

Warm pinned-Blender wall times were 10.2 s for the owner build, 19.3 s for
the two owner renders/contact sheet, 4.7 s for material-only publication, and
23.6 s for the two final renders/contact sheet. The earlier relative-path
launch failed before opening any blend; switching to absolute workspace paths
removed that setup bottleneck.
