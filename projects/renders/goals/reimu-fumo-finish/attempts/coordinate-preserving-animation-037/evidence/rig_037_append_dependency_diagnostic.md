# Rig037 append library diagnosis

The failed `assert not bpy.data.libraries` does not identify a live external
dependency. One exact collection-only append produced **zero linked,
indirect, linked-override, missing or external-file IDs**. All 86 objects,
85 meshes, 18 materials, their 18 embedded shader node trees, one armature
and four actions are local.

The sole Library record is `rig_037_candidate.blend`, pointing to the exact
input candidate. It has `users=1`, no fake user, no parent library, no
missing flag, zero IDs with `ID.library` pointing to it, and zero referencing
IDs in `bpy.data.user_map()`. Do not interpret its internal user count as a
usable linked-data dependency.

105 local IDs retain weak append-origin references to that candidate path:
85 meshes, 18 materials, one armature and the collection. These are append
provenance, not live data links. Their machine-local path strings remain a
separate final criterion 007 packaging concern; this intermediate diagnosis
neither removes them nor grants a packaging pass.

## Exact enumeration and replacement guard

The diagnostic enumerated 239 IDs across every RNA-exposed Main collection,
added embedded node trees/master collections/shape keys, and merged all IDs
seen by Blender's own user map. This includes mesh/material/action
dependencies, not merely object linkage. The relevant guard can be copied
without purging or changing any data:

```python
ids = {}
for prop in bpy.data.bl_rna.properties:
    if prop.type == "COLLECTION":
        for block in getattr(bpy.data, prop.identifier):
            if isinstance(block, bpy.types.ID):
                ids[block.as_pointer()] = block
queue = list(ids.values())
for block in queue:
    for key in ("node_tree", "collection", "shape_keys"):
        child = getattr(block, key, None)
        if isinstance(child, bpy.types.ID) and child.as_pointer() not in ids:
            ids[child.as_pointer()] = child
            queue.append(child)
for block, users in bpy.data.user_map().items():
    for item in (block, *users):
        ids[item.as_pointer()] = item
bad = []
for block in ids.values():
    override = block.override_library
    linked_override = (override and override.reference
                       and override.reference.library)
    if block.library or block.is_library_indirect or linked_override:
        bad.append((block.bl_rna.identifier, block.name_full))
assert not bad, bad
assert not [block.name_full for block in ids.values() if block.is_missing]
```

The four local actions are `ReimuFumo_ArmWave`, `ReimuFumo_Combined`,
`ReimuFumo_HeadYaw` and `ReimuFumo_SeatedNeutral`. Machine evidence additionally
records object data/material/modifier/action/custom-ID references, material
node-ID dependencies, every strong referencing ID, and every weak origin.
The external image/movie/sound/font/volume/cache filepath inventory is empty.

## Provenance and limits

Exactly one unsaved pinned Blender 5.2.1 LTS (`9e2066aef7ef`) diagnostic
reproduced root's `libraries.load(..., link=False)` collection-only append.
It preserved candidate SHA256
`e4d42acebc89f664f5c0576b28ecba1b175daee40b3984daf24286b3fcdfcef0`
and the root append script's bytes. No save, purge, candidate mutation,
root-script modification, pose test, or second execution occurred.
Save/clean-reopen/playback verification remains root-owned.

Evidence under this directory:

- `rig_037_append_dependency_diagnostic.json`: complete ID and dependency
  records, library provenance, input/script hashes and preservation checks.
- `rig_037_append_dependency_diagnostic.py`: exact one-append diagnostic.
- `rig_037_append_dependency_diagnostic.log`: successful pinned execution.
