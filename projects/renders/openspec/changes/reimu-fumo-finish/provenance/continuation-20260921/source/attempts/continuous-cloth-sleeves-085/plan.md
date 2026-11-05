# 085 continuous cloth sleeves with a shorter rest shape

Coordinator verdict: PROCEED after closing 084. Root remains sole model/goal
writer in the current linked feature worktree. Pinned Blender 5.2.1 LTS,
build9e2066aef7ef, factory-startup/background/disable-autoexec, four threads.
This uses the already proven native cloth route and the documented
cloth-rest-orco-leading-deformer workaround from cloth_rest_044_evidence.md;
it is not a fresh authoring-capability claim.

Source hair_084_candidate.blend
cb6bff1fa0e654f1720236a5a90dab56346f41e36f999dfc08fe746ea26294cb.
Only two sleeves and their red dashes may change. Head, hair, face, arms,
body, bow, skirt, materials and camera remain exact. Donor sleeve080
7914792f9f94094c4728731677fcaf6336b34a2505400708985f89034ec875c8
supplies its valid 4000-vertex outer sheet per side, including the160-sample
receiver-seated shoulder boundary. Preserve all protected bytes.

082's independent receiver projection folds the outer sheet. Instead use one
combined two-sheet native cloth solve, initial geometry from080, pin only the
320 outer shoulder vertices. Define rest geometry with cuff bottom raised
17mm (17 to34mm), lower-length blend u squared times v; reduce cuff depth
31 to21mm by a u-weighted scale about its Y center. No per-vertex receiver
projection or cuff pinning. The rest edge lengths and cuff bending provide
continuous aperture support, whose success must be judged, not assumed.

Both proxies at10x scale;48frames;180s internal solver deadline. Quality8,
mass0.10,tension/compression20,shear10,bending0.12,air3. Last8% cuff bending
weight ramps to1, max0.8. Enabled Decimate COLLAPSE ratio1 BEFORE Cloth,
rest_shape_key bound to short rest shape (value0), dynamic mesh disabled.
Collision and self collision on; quality4; cloth/proxy distance0.6mm each,
self distance0.5mm, friction4/self5 (multiply by10 in simulation units).
Use actual evaluated arm and nearby body/garment collision proxies. Record
initial and final pinned/free crossings: initial embedded pinned points are
not automatically repairable. No collision-disabling workaround.

Freeze the final outer sheet and a0.5mm inner return as an explicit closed
shell, with original chart correspondence for reattaching existing red dash
geometry. Outer shoulder samples unchanged. Inner shoulder gap is measured;
up to1mm may be a concealed attachment allowance, not a visible floating seam.
No simulation proxies, caches or new runtime cloth modifiers enter the saved
candidate. Restore scene frame/range/gravity after removing only owned proxies.

One build plus one causal API/data repair after a contemporaneous checkpoint;
no post-render geometry or solver-dose variant. Save sleeve_085_candidate.blend
once if finite, closed, positive and no self-intersection; record unresolved
receiver contacts explicitly as diagnostic rejection risks. Do not suppress
failure evidence or call diagnostics a construction pass. Compare front cuff
bottom/span to canonical front and side aperture to turn frames10/12. Reject
collapsed aperture, tangled cloth, new arm exposure, detached root, excessive
wing height or important non-target changes. Fixed512triple first; full packet
and independent image review only if plausible. No whole stage/goal acceptance.
