# Rest-coordinate mechanism proof

Stable defect: cloth-rest-orco-leading-deformer.042b/042c requested different
flat rest widths but produced byte-identical2016vertex coordinate arrays.
Read-only review of exact Blender build9e2066aef7ef identified that the
leading-deformation path omits the cloth-specific rest-coordinate layer.
Sources: https://github.com/blender/blender/blob/9e2066aef7ef/source/blender/blenkernel/intern/mesh_data_update.cc
and https://github.com/blender/blender/blob/9e2066aef7ef/source/blender/modifiers/intern/MOD_cloth.cc

Root runtime probe inserts enabled Decimate COLLAPSE ratio1 before Cloth.
Two independent opens of exact041b test the existing53mm and36mm patterns.
Both initial meshes preserve2016vertices/1920faces and position delta0.0.
At frame5 their maximum model-space displacement difference is57.799mm,
versus previous0.0. Thus the evaluation-path workaround restores actual
rest-pattern sensitivity. This is not a calibrated fabric or visual pass.

Probe source SHAe73b5cd0e62fc32b2ea41e23b2548fc3e7824d5b77acf7d4c04680b3739f4f6f;
helper SHAb4717b0df46b7cb5c916a8ac4e6a3736cc53de0003477e067c706950dd3583fc.
No source file or host Blender installation changed. Next: one bounded
short-pattern solve, save extracted ordinary cloth, then inspect fixed views.
