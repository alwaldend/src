## Assets

- [x] Add the pinned assets and regenerate the root Bazel module.
- [x] Match the Pippa download to the original FBX and document its import procedure.
- [x] Build the asset targets and verify the delivered files in Blender.
- [x] Fetch the final album's 41 Mixamo clips and verify their original checksums.

Bazel resolved ten files. Blender 5.2.1 opened both scenes, the human base
mesh library, and the Snow rig. It imported three BVH clips and two furniture
models, and loaded the supplied diffuse texture onto the matching Pippa FBX.
