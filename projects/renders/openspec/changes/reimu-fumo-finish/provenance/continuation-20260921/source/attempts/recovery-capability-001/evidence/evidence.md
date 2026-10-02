# Recovery and capability evidence

## Recovery

The current worktree contains 61 copied files under
`out/reimu_fumo_finish/recovery/` (33 MB). The protected candidates still hash
to:

- A157: `433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
- A202: `6a9f3757facba526550e78817dc85f1d23cf85bcdad360228e113bb60d5f3aa0`.

All seven tracked reference copies match their original LFS object hashes.

## MCP save and pinned reopen

The Blender MCP background interface opened the exact A157 file, reported 180
objects and no missing files, created mesh `MCP_CAPABILITY_COUPON`, and saved a
new candidate:

- candidate SHA-256:
  `7c2db35088fd9c248bdb7316a61fc5de186f98e4efb8c28ce079e6a7b2b71771`;
- marker dimensions: approximately `0.24 x 0.24 x 0.24` scene units;
- source candidates remained byte-identical after the operation.

Pinned Blender 5.2.1 clean-reopened the candidate with automatic file scripts
disabled and rendered `Review_front_Camera` successfully. The rendered image
hash is
`b1998f313c60150d34f8eddfc5c64a23d8323ce1f95e209914779a745cf98797`.
The report found no missing external images.

## Failure

Pixel inspection showed the correct complete A157 subject, but the new marker
was outside the frozen camera frame. The save/reopen path is technically
functional, yet the planned visible-effect proof failed. The placement method
must use camera projection rather than scene bounds.

