# A68 C0 fixed-view baseline render packet

## Verdict

The packet is a technically valid pre-sculpt baseline. All five calibrated
views rendered at 640 by 640 pixels, are nonblank, contain the complete subject,
and were visually inspected. It is not a visual acceptance pass.

The baseline clearly exposes the intended A68 defect: in side and both
three-quarter views, the brown rear/crown cap reads as an oversized spherical
helmet with a near-vertical rear wall. The front remains broadly recognizable,
and the rear view confirms complete hair coverage, so subsequent sculpt work
must reduce the spherical side depth without reopening a bald region or
damaging the fringe, bow, face, locks, body, cameras, or lights.

## Frozen input

- Blend:
  `out/reimu_fumo_attempt_068_sculpt_coupon/candidate_c0_pre_sculpt/reimu_fumo_a68_c0.blend`
- SHA-256 before render:
  `26c8613fe3eb17a1ddfcf7c8b596ed2aa264162b86d2b1e81acf7033d1fa75ba`
- SHA-256 after render:
  `26c8613fe3eb17a1ddfcf7c8b596ed2aa264162b86d2b1e81acf7033d1fa75ba`
- Scene: `Attempt41_Manual_Head_Maquette`
- Render specification SHA-256:
  `d377222d84dd64aacaf7edb071f50929ecd880f781392c0fdb4060256659d1d8`

The equal before/after blend hashes confirm the read-only packet renderer did
not modify the candidate.

## Tool identity

- Bazel entry point: `//projects/renders/cmd/fumo_review:render_packet`
- Repository toolchain: `//tools/blender:blender`
- Blender: `5.2.1 LTS`
- Build hash: `9e2066aef7ef`
- Build date: `2026-08-25 02:12:34`
- Render process exit status: `0`

## Commands

Tool identity:

```sh
bazel_agent run //tools/blender:blender -- --version
```

Five-view render, performed in one pinned background-Blender process:

```sh
bazel_agent run //projects/renders/cmd/fumo_review:render_packet -- \
  --blend-file \
  out/reimu_fumo_attempt_068_sculpt_coupon/candidate_c0_pre_sculpt/reimu_fumo_a68_c0.blend \
  --spec \
  out/reimu_fumo_attempt_068_sculpt_coupon/candidate_c0_pre_sculpt/render_packet/spec.json \
  --output-dir \
  out/reimu_fumo_attempt_068_sculpt_coupon/candidate_c0_pre_sculpt/render_packet/packet
```

The target supplies background mode, factory startup, disabled auto-execution,
and a nonzero Python exception exit code.

## Pixel artifacts

Every output is a non-interlaced 8-bit RGBA PNG at exactly 640 by 640 pixels.

| View | Camera | Bytes | SHA-256 |
| --- | --- | ---: | --- |
| Front | `Review_front_Camera` | 522244 | `f91065e25ead7e6536772e592832e74209d4c15a769ddde24f0f94676e4326e5` |
| Side | `Review_side_Camera` | 474626 | `7b6cfbdf60481606ac8c54ba159b6fa46ca182e728d6273d40aba323791d83c4` |
| Rear | `Review_rear_Camera` | 506852 | `3f899f64f22c330b697011961312b7ed42b04ff9f9ebcf96383a09a249ba4ad9` |
| Three-quarter | `Review_three_quarter_Camera` | 503593 | `1fcf249e81c675197791e8b4360e6e700147a54e28d19222c400fce9fc122e80` |
| Mirrored three-quarter | `Review_three_quarter_mirror_Camera` | 508822 | `1588cc0b8bc4ce2f6d6e500d6c388acd9950b6bcd0075b1ec1266c9b385f1d5e` |

The machine-readable packet manifest is
`packet/manifest.json`, SHA-256
`e7544d2e65900cf55963cd20df45bd6d758887701a61fcefc3332abe0151c875`.

## Visual inspection notes

- Front: full subject and bow are framed; face and costume are readable; no
  blank render or camera clipping. The cap looks broad and smooth.
- Side: full subject is framed; the cap projects too far rearward and terminates
  in a near-vertical wall, producing the strongest spherical-helmet evidence.
- Rear: full subject is framed; hair coverage is complete, with no bald patch;
  the central and flanking rear pieces remain visible.
- Three-quarter: full subject is framed; excessive rounded cap depth dominates
  the head and separates visually from the flatter fringe and lock language.
- Mirrored three-quarter: the same defect survives bilaterally, ruling out a
  one-camera or one-side artifact.

No blocker occurred. This packet is suitable as C0 evidence for direct visual
comparison with later sculpt checkpoints.
