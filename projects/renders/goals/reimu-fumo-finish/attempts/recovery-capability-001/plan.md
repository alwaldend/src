# Recovery and authoring-capability plan

This is the first work unit of successor goal `reimu-fumo-finish`. PR #24 at
`803abe43cb9ec177b4ac94e2f2490b128952e231` and the unpublished local
continuation at `401582f36e23b72b9ac68bdde4ebf7e3bf95e698` remain the historical
archives; their 137 attempts and unsupported pass claims are not re-imported.

## Exact inputs

- A157 visual scaffold: SHA-256
  `433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
- A202 technical donor: SHA-256
  `6a9f3757facba526550e78817dc85f1d23cf85bcdad360228e113bb60d5f3aa0`.
- Seven controlling/supporting references copied without byte changes into
  `projects/renders/assets/reimu_fumo/references/`.

## Work unit

1. Verify both candidates and their useful render packets are present in this
   worktree under `out/reimu_fumo_finish/recovery/`.
2. Keep the source files immutable and use a new candidate path.
3. Through the local Blender MCP background interface, make one reversible,
   visible test edit to an isolated A157 copy and save it as a new file.
4. Clean-reopen and render the edited file with pinned Blender 5.2.1.
5. Rehash both protected recovery inputs and inspect the rendered pixels.

The test edit is a capability coupon only. It cannot pass a visual goal
criterion or become a modeling baseline.

## Stop conditions

- Reset if MCP cannot open the exact A157 bytes, cannot save a distinct file,
  or changes the protected source.
- Reset if pinned Blender cannot clean-reopen the saved file or produces a
  blank/wrong-subject render.
- Refine into the first head-and-hair cycle only if all capability checks pass.

## Parallel workstreams

Recovery inventory, process/schema review, acceptance-tooling review, landmark
review, and read-only blend inventory use disjoint outputs and run in parallel.
Only the coordinator writes canonical goal state. Model mutation remains
sequential because one `.blend` has no safe merge boundary.
