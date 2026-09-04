# Macro038 tessellation API repair

First execution stopped before saving any candidate: pinned Blender5.2.1
tessellate_polygon returns triangle vertex indices, not iterable vectors.
Correct the index dereference in the patch builder; retain vector support
for explicit compatibility. No shape parameter changes. Expected effect:
polygon geometry construction completes instead of TypeError at line133.
Original failure preserved in macro_038_build.log; repaired execution uses
macro_038_build_indices.log. No protected input was saved or overwritten.
