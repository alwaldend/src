# Frozen 027 material-ID failure

Read-only clean reopen under pinned Blender 5.2.1 LTS confirmed a saved
material-slot fault. No geometry audit, render, visibility change or model
save was performed.

Sources remained unchanged:

- 026: `56efb16739c746153c5a562195b221645865e0ae4a6c78a5f491783b2c700882`
- 027: `8c38a46141eae10d44a00d61de1d745e134717a5ce2fe42a9ad3433b487f9ae4`

## Confirmed evidence

All 13 snapshot-derived new objects have null persistent material slots
and null evaluated material slots after reopening 027. Every evaluated
face on these objects refers to a null slot: the fringe, mouth, six
eye/embroidery pieces, both cheek locks, and all three rear cloth panels.

The two exceptions are `Head027 sewn cushion` and
`Hair027 crown and back hood`. Their persistent material slots remain valid
members of `bpy.data.materials`, and all their evaluated faces have materials.
Their construction explicitly selected database materials rather than
retaining the snapshot's evaluated material references.

On 026, every inspected original object's persistent material slots are
valid, non-evaluated database IDs. Its `evaluated_get(...).to_mesh()`
material slots instead expose evaluated IDs: `is_evaluated == True`, not
the same pointer as the named `bpy.data.materials` entry. Each evaluated
material's `.original` resolves to that valid database material.

This supports the identified construction fault: the snapshot retained
evaluated material IDs for later attachment to persistent new meshes.
Resolving those references to their persistent database materials addresses
the observed distinction. The existing writer appearance record did not
detect this because it recorded the global material datablocks, not the
new meshes' slot references; the datablocks themselves were still intact.

## Limits

The missing saved slots are proven. This probe did not reproduce or prove
the native-crash cause. The corrected writer's dependency-disable execution
is a separate causal test. No final material or geometry acceptance follows.

[Durable identity report](head_027c_machine_evidence.md) records every
inspected persistent/evaluated slot, database/original identity relation,
node-tree presence, face material index, and file-preservation check.
