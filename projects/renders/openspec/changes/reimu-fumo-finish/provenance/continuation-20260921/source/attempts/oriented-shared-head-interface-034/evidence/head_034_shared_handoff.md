# Shared Head034 first-state handoff

`head_034_shared_draft.py` exposes `TARGETS` containing exactly the old fringe
and hood, plus `build_head_034_shared()`. It creates one
`Hair034 shared crown fringe shell`, with one ordinary live full-weight Head
Armature. It captures the two existing thickness skins; no Solidify, GN,
second rig, or post-armature correction is added. Import only defines
functions. The helper does not open, save or render a model.

## Frozen execution

- Source SHA256:
  `6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
- Evaluated helper SHA256:
  `c9976b58dc310001b4478ddb399b6691bd801eb3a725b12358c6422000cb1c69`.
- One unsaved state, pinned Blender 5.2.1 LTS, build `9e2066aef7ef`.
- Full evidence: `head_034_shared_dryrun.json`; short bundle:
  `head_034_shared_metrics.json`.
- Source bytes remained exact. No save, render, second state or shape sweep.

The helper binds exact032 head-input records, the cut-set JSON, and utility
code hashes. A later source additionally requires root's explicit path/hash
and must still match all 15 head inputs. Root owns any saved candidate.

## What the state demonstrates

Both exact cut arcs were resampled into 432 common longitudinal knots.
Retained boundary triangles were split in their existing planes; 366 source
triangles form that support. The new outer and inner Hermite strips share
their endpoint vertices with the retained skins. The old 306 fringe rim
faces are removed, the retained hidden brown-cover boundary is capped, and
the two side connectors are closed. The resulting shell has 44,366 vertices,
81,704 faces, zero open/non-manifold edges and no inconsistent edge winding.

All 84 non-target evaluated controls and the rig pose remained exact.
The 72,718 indexed off-support triangles have zero measured coordinate
error. Materials retain their source slots and datablocks. Both source
targets have no UV layers, so the UV-corner check is vacuous; this does not
claim that the asset has a completed UV workflow. Source off-support smooth
shading normals were recomputed after triangulation and are **not** claimed
equal. That is a separate visual/normal-preservation risk.

The 3,888 paired strip samples have skin separation 0.711–1.158 mm, median
0.912 mm. Boundary chord length is 0.820–10.326 mm, median 8.155 mm.
These are paired-point measurements, not minimum shell thickness, absence
of self-intersection, or proof against envelope overshoot.

## Ownership result and its qualification

All sampled side and regular three-quarter rows now progress through the
retained fringe, shared bridge and retained hood without a repeated region.
No sampled row exposes the hidden front cover or inner allowance.
Face-region IDs are stored in the new mesh's `head034_face_region` attribute;
checking the common object name alone is explicitly insufficient.

The strict first-state checker reports one mirror row:

| Mirror row175 pixel X | First-hit face region |
| --- | --- |
| 214–236 | Retained hood outer |
| 237–248 | Shared outer bridge |
| 249 | Retained fringe outer |
| 250 | Shared outer bridge |
| 251–279 | Retained fringe outer |

That is a two-pixel return across a **welded** boundary. It may be the
projected shape of the discrete inner cut, rather than two competing outer
sheets. The root correctly distinguished that possibility from the earlier
tuck failure. The checker therefore supplies a review flag, not proof of
the old overlap mechanism and not an automatic visual rejection.

The completed process did not record the per-hit candidate polygon/triangle
IDs, local depth or adjacent normals at pixels249/250. Those details are
unavailable; they must not be inferred from the region names. Exact source
cut edges are available in `head_034_shared_topology_probe.json`, but they
do not identify which candidate triangle won those rays. No strip
intersection, envelope overshoot or junction-tangent audit ran. No replay
was made to fill these gaps after root requested a stop.

## Suitability and next authority

The helper's strict region checker returns `STOP before writer`. Treat the
candidate as held for root review of that checker/geometry distinction,
not as a passed visual state. It demonstrates a closed shared interface and
exact retained coordinates, but does not yet establish sufficient tangent,
intersection, shading-normal or visual quality evidence for acceptance.

Root may choose a specifically authorized diagnostic of this unchanged
state before deciding whether the flagged boundary is benign. Do not alter
the cut, strip depth, thickness, tangents or topology to silence the checker
without a new causal decision. No saved candidate or render is supplied.
