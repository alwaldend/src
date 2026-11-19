# 103 paired-sheet clipping helper

Final immutable `hair_103_free_boundary.py` SHA-256:
`db656380842ede1d19f2581fc0a45a6c732ff7d415165499a8a9b0cf0a2b1c5d`.

The helper implements the frozen `hair_103_plan.md` clipping construction.
It has no native imports, file access, deformation, coordinate welding or
normal offsets. Root/native/model/goal state remains the coordinator's.

```python
packet = build_free_boundary_cover(
    outer_vertices, inner_vertices, triangles, face_kinds,
    root_z=0.160, hem_values=hem_values, zero_tolerance=3e-8,
)
```

Coordinates, root Z, hem scalar and tolerance use metres. `face_kinds` contains
`front`, `rear` or `mixed` per original outer triangle. Input outer and inner
arrays share indices; `triangles` is outer winding. Every source triangle is
split at the classified original outer root field. High polygons survive;
low front survives, low mixed is removed, and low rear is clipped by the
interpolated classified hem field. The caller supplies the frozen hem values.

Original nodes retain their exact coordinate values. Root cuts share canonical
original-edge keys. Hem cuts share canonical root-split-edge keys. Every new
node carries one original-source barycentric weight map applied identically
to both actual input sheets. Scalar fields are interpolated at inserted nodes,
with explicit near-zero classification; neither field is reevaluated from the
inserted world coordinates. No point is snapped onto a world plane.

When a hem crosses a root edge, that node is also inserted into the retained
neighbor's boundary. Convex ear triangulation in original barycentric charts
preserves collinear boundary nodes and prevents a T-junction. Exact predicate
fallback resolves near cancellation. Positive slivers are retained; actual
triangle collapse/reversal fails. Zero/one-dimensional kept fragments are
omitted and counted. `FreeBoundaryError.quality` retains partial diagnostics.

The wrapper-required keys are present exactly:
`source_outer_vertices`, `source_inner_vertices`, `outer_triangles`,
`triangle_source_faces`, `face_kinds`, `face_zones`, `locked_upper`,
`shell_faces`, and `quality`. Zones are exactly `high` or `low`.
Additional fields are `source_weights`, `node_keys`, `root_values`,
`hem_values`, `boundary`, `boundary_loops` and `shell_labels`.

Source indices in weights and `triangle_source_faces` refer to the supplied
arrays, not original whole-head IDs; the caller composes inherited ancestry.
Unused nodes are compacted through an explicit canonical-key index map.
`locked_upper` is true when the classified OUTER root scalar is nonnegative.
A zero/root node locks its paired inner too, irrespective of that inner's Z.
Non-rear preservation during deformation is the caller's responsibility.

Shell vertex order is compact outer followed by compact inner. Shell faces
are outer, reversed paired inner, then two oriented rim triangles per boundary
edge; labels are `outer`, `inner` and `rim`. Every rim spans only its own
paired edge. There is no closure across a removed lower wrap.

Checks enforce consistent edge winding, one connected orientable outer sheet,
connected vertex fans, degree-two directed boundary loops, and a connected
closed paired shell. Both clipped sheets preserve each original triangle's
winding and have positive finite triangle altitudes. Rim triangles must be
nondegenerate. Global intersections, signed volume, projected graph validity,
native float conversion, deformation and receiver/core coverage are caller
checks. The helper permits multiple manifold boundary loops in general;
the actual 094 smoke output has one loop and is a disk.

Bounded validation completed on 2026-09-06. Synthetic checks passed root/hem
intersection inside a source triangle with propagation to its high neighbor;
same-weight paired reconstruction; exact source-node reuse; paired-inner
root locking below its own root Z; front retention; low mixed removal;
classified-root faces kept once; reversed source winding; positive tiny
slivers and exact-cancellation orientation; and rejection of duplicate faces,
disconnected components, zero-thickness rims and degenerate source triangles.

Undeformed actual 094 smoke input SHA-256:
`492c89b457110d01273a1d8de0616dc266d18c5c441aa634c871395c54946b7f`.
The frozen hem values were supplied once; no field deformation or parameter
trial was performed. The input bytes remained unchanged.

| Actual clipping result | Value |
| --- | ---: |
| Compact paired nodes | 5,364 |
| Outer triangles | 10,434 |
| Closed shell triangles | 21,452 |
| Boundary edges / loops | 292 / 1 |
| Outer / shell Euler characteristic | 1 / 2 |
| Retained original / root-cut / hem-cut nodes | 5,050 / 220 / 94 |
| Low mixed polygons removed | 181 |
| Omitted zero-dimensional fragments | 4 |
| Positive slivers discarded | 0 |
| Minimum outer / inner altitude | 3.524548 / 3.489303 micrometres |
| Original pure-front nodes retained exactly | 1,871 / 1,871 |
| Maximum per-front-face area error | 7.624e-21 square metres |
| Maximum front plane residual | 5.731e-18 metres |
| Maximum source weight-sum error | 0 |
| Clipping/check runtime | 0.436 seconds |

`git diff --check` passed. No old helper, native file, model, goal or tracked
source was changed. The helper and this note are final for root hash binding.
