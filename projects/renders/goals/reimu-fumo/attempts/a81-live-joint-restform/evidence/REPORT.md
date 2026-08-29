# A81 two-view pinned batch render

- Source: `../live_author/a81_live_joint_restform_pair.blend`
- Source SHA-256 before and after:
  `bc1cf8fb4fb669076ac3199210fedcf43b534cea70d0048b12f4a1f4d6e197f2`
- Renderer: repository-pinned Blender 5.2.1 LTS through
  `//projects/renders/cmd/fumo_review:render_packet`
- Scene: `Attempt41_Manual_Head_Maquette`
- Resolution: 512 x 512
- First pixel artifact completed:
  `2026-08-31T16:53:21.301378600+03:00`
- Two-view packet completed:
  `2026-08-31T16:53:26.918902624+03:00`
- Blender elapsed from source open to front save: 12.289 seconds
- Blender elapsed from source open to pair save: 17.906 seconds

## Artifacts

- `packet/front.png`
  - camera: `Review_front_Camera`
  - SHA-256:
    `145f0c4ba76581d4947b83759824f353404fa66d2271a59bc4aab494f82eebe8`
- `packet/three_quarter.png`
  - camera: `Review_three_quarter_Camera`
  - SHA-256:
    `a1d26222a1ded976deb8e341dfcea9ed8f7650adf83b001855dcbf5f55100252`
- `packet/manifest.json`
  - SHA-256:
    `3d23e1c235529dc3b0c3ea588d4fa833a238bb065ffcf725b09f4c392c927f86`

Both images are nonblank, framed 512 x 512 beauty renders of the complete
plush. The renderer did not save the source.

## Protected-input verification

- Rung 003 remained
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Tracked reusable asset remained
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.

## Blocker observed

The first sandboxed invocation stopped before Blender because the configured
Bazel output base under `/var/cache/bazel` was not writable. The authorized
rerun with access to that existing Bazel cache succeeded. No rendering or
source mutation occurred during the failed invocation.
