# A01 concrete camera-fit and clean-reopen plan

## Target

Test whether the current mechanically successful neutral scaffold needs only
a one-property camera correction to expose the complete initial drop.

## Frozen input

`out/fumo_concrete_drop_scaffold/attempt_02/fumo_concrete_drop_scaffold.blend`
(`sha256:a9488c220c5076a3202e61c9897cf3710f24b1abe74fb9edfc4750bfaebfdc26`)
is immutable. The tracked Reimu and Sisyphus Blends are protected at
`489213b7...` and `c5bd58ed...` respectively.

## Bounded work

Copy the exact input into this attempt. With repository-pinned Blender 5.2.1,
change only `Drop Camera Data.lens` from 45.0 mm to 35.0 mm, reset to frame 1,
and save a candidate. Do not rebuild the scene or change camera transform,
target, lighting, animation, physics, interface, materials, or labels.

## Outputs and gates

Clean-open the exact saved candidate without resaving it. Emit a semantic
delta, machine audit, protected hashes, exact renders at frames 1, 12, 20, 22,
28, 40, 56, and 72, and a labeled contact sheet. Reset immediately if pinned
5.2.1 cannot open it, anything besides the lens/file-format metadata changes,
mechanics regress, the interface differs, protected hashes change, or any
required pixel touches the image boundary; frames 1 and 12 require at least
12 px top and bottom subject margins.
