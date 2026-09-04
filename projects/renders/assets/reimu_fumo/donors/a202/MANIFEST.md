# A202 privacy-sanitized technical donor

This packet preserves the furthest useful unpublished Reimu work found after
PR #24. It is a rejected parts donor, not a visual baseline, stage pass, or
accepted asset.

The published `model.blend` is a privacy-sanitized derivative, not the exact
historical bytes. Its SHA-256 is
`a5e1e96dbbabaee9d4f23c28d95930509082644124adab4607e2757b708852b5`.
The exact recovered source remains ignored in this worktree with SHA-256
`6a9f3757facba526550e78817dc85f1d23cf85bcdad360228e113bb60d5f3aa0`.

## Derivation

Pinned Blender 5.2.1 copied and remapped the 113 unique datablocks carrying
historical library weak-reference paths, cleared one 41,024-byte embedded
script body, and cleared one absolute render path plus one relative bake-output
path. It made no other intentional content change. Inventory and the geometry,
rig, and material fingerprints match the recovered source.

A separate pinned-Blender process clean-reopened the derivative and scanned
447 ID datablocks, 5,051 reachable RNA structures, 19,138 RNA strings, and 102
custom strings. It found no local path or identity marker and no unreadable
property. The package Blender test repeats a fail-closed privacy, dependency,
rig, and camera audit against the committed file.

## Review packet

The five PNGs were freshly rendered from the sanitized derivative with the
cameras and settings bound in `manifest.json`. Their textual, EXIF, and time
metadata was removed. The packet-integrity test recomputes every model and
render digest and byte count.

Fresh five-view review still scores the visual state `5/10`. The head and hair
form a deep cubic helmet, the body is a tall cone, sleeves are rigid open tubes,
feet are disconnected pods, and the bow and hair roots are hard panels. The
rig uses coarse whole-object weights, its saved action is empty, and visible
curves and lattices remain unbound.

A202 contains useful collection/material and camera definitions, an armature
and control taxonomy, per-object bindings, and later topology cleanup. A new
artist may donate those parts selectively or start clean. No old criterion
pass is inherited; the historical measured-landmark criterion never passed.

The old lineage is internally inconsistent: later records treated A172 as
accepted after its original reset, and A203 prose claimed all criteria passed
although its structured review covered only criterion 008. This packet records
those bytes as recoverable technical material without preserving the false
acceptance claim.
