# 104 side/rear paired-ownership clipper

Final immutable `hair_104_side_rear_boundary.py` SHA-256:
`c6f4281c9005acedff5f4c14c4f0be2664bf894ecf07a5c5cb565e5ac45eda39`.

Read and implemented the complete frozen `hair_104_plan.md` construction.
This new file derives from the immutable 103 helper; the 103 file remains
unchanged at SHA-256
`db656380842ede1d19f2581fc0a45a6c732ff7d415165499a8a9b0cf0a2b1c5d`.
No native process, model edit, deformation, shape variant, goal write or old
helper overwrite occurred. Only this helper and note were created.

```python
packet = build_side_rear_cover(
    outer_vertices, inner_vertices, triangles, face_kinds,
    root_z=0.160,
    rear_hem_values=rear_field,
    side_hem_values=side_field,
    zero_tolerance=3e-8,
)
```

Coordinates, root Z, both scalar fields and tolerance use metres. Source
triangles use outer winding with matching paired inner indices; face kinds
are `front`, `rear` and `mixed`. All original root/high polygons and low front
polygons remain. Low mixed polygons are clipped only by the supplied side
field. Negative-root low rear nodes get canonical independent ownership
before clipping by the supplied rear field. Root/high nodes remain shared.

Original paired source coordinates and one source-weight map are preserved
for every copy. Root and hem cuts interpolate both actual source sheets with
the same weights. Scalar values are interpolated at new nodes and classified
near zero; no analytic field is reevaluated from inserted world coordinates.
No source point is snapped or moved. Shared edge cuts propagate to retained
neighbors before source-chart ear triangulation, including partial retained
edge intervals. Positive slivers remain; actual per-face collapse or reversal
fails. Only exact lower-dimensional retained fragments are omitted and counted.

`node_owners` has exactly three values:

- `shared_upper`: nonnegative classified outer root scalar; paired inner is
  locked too, regardless of its own Z.
- `front_side`: negative-root original/shared front or mixed ownership.
- `rear`: negative-root independent rear ownership.

The function returns the existing 103 wrapper keys: paired
`source_outer_vertices`/`source_inner_vertices`, `outer_triangles`,
`triangle_source_faces`, `face_kinds`, `face_zones` (`high`/`low`),
`source_weights`, `node_keys`, `locked_upper`, `shell_faces`, `shell_labels`,
and `quality`. It also returns root/rear/side scalar arrays and oriented
`boundary`/`boundary_loops`. `hem_values` aliases the rear scalar array for
compatibility. Source indices address the supplied arrays; the caller composes
older whole-head ancestry.

Shell vertex ordering remains outer followed by inner. Outer faces, reversed
paired inner faces and each boundary's own two-triangle paired rim close the
shell. There is no underside bridge. Checks enforce consistent shared-edge
winding, connected vertex fans, manifold directed boundary loops, one outer
component, one closed shell component, and no duplicate indexed faces.
Coincident undeformed source copies are allowed and are not welded. In
particular, a topological pass is not a global geometric self-intersection or
signed-volume certificate. Actual mapped crown termination, wall intersections,
sheet ordering and all native/receiver checks remain the coordinator's gates.

Bounded synthetic checks on 2026-09-06 passed a side/rear slit ending at a
shared root, exact coincident paired copies with different ownership, inherited
source weights, root locking with inner Z below root, continuous front/mixed
ownership with propagated side-hem cuts, retained front area, owner/face-kind
consistency and no duplicate indexed shell faces.

The clipping-only smoke used the unchanged actual 094 arrays, SHA-256
`492c89b457110d01273a1d8de0616dc266d18c5c441aa634c871395c54946b7f`,
with the two frozen fields supplied once. No source field deformation was run.

| Actual undeformed result | Value |
| --- | ---: |
| Compact paired node instances | 5,387 |
| Outer / shell triangles | 10,398 / 21,544 |
| Boundary edges / loops | 374 / 1 |
| Outer / shell Euler characteristic | 1 / 2 |
| Distinct retained source vertices / instances | 4,986 / 5,019 |
| Used duplicate source-weight groups | 33 |
| Root / rear-hem / side-hem cut nodes | 220 / 142 / 6 |
| Shared-upper / front-side / rear node instances | 1,704 / 1,219 / 2,464 |
| Low mixed polygons processed | 181 |
| Omitted zero-dimensional fragments | 4 |
| Positive slivers discarded | 0 |
| Minimum outer / inner altitude | 3.524548 / 3.489303 micrometres |
| Maximum per-front-face area difference | 7.624e-21 square metres |
| Maximum source weight-sum error | 0 |
| Runtime | 0.482 seconds |

Original coordinate values, input bytes and the old 103 helper were checked
unchanged. `git diff --check` passed. These bounded cases support the frozen
104 clipping input; arbitrary alternate scalar/topology arrangements were
not exhaustively tested and may correctly fail the manifold gates.
