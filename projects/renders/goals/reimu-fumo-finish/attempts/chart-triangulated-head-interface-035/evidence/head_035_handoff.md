# 035 perimeter prerequisite failed

The single authorized four-seed plus one-edge-neighbor expansion is not a
simple, single-valued XZ domain. It still crosses at both side roots.
Stopped before triangulation, numerical solving, object construction,
hiding, saving, or rendering. No `head_035_draft.py` geometry helper was
created because its prerequisite failed.

## Exact scope tested

Input: `head_032_candidate.blend`, SHA-256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
Pinned Blender 5.2.1 LTS, build `9e2066aef7ef`.

Each seed was unioned with its immediate edge-sharing polygon neighbors,
then with the prior removed set. There was no recursive expansion or
numeric cutoff adjustment. Exact seed-neighbor lists, added polygon
vertices, ordered perimeter points, and original exposed-edge witnesses
are in `head_035_domain_probe.json`.

| Object | Newly included polygons | Total removed-set count |
| --- | --- | --- |
| Hair028 traced padded fringe | 67, 68, 10938, 10941 | 1,017 |
| Hair028 crown and back hood | 7742, 7743, 7744, 7835, 7836, 7837, 7899, 7992 | 968 |

The added fringe support spans Z 134.2564–140.1991 mm, crossing the former
135 mm cutoff as authorized. Its new inner arc has 333 vertices and
endpoints 2914/3612. The added hood support spans Z 135.9000–141.5364 mm;
all eight added polygons are actual gusset faces, verified through original
core correspondence. Its receiver arc has 105 vertices and endpoints
793/841. Both arcs are unbranched. The 2,140 retained front-cover faces and
its 117-vertex boundary arc remain outside the removed set.

## Definite rejection witnesses

The combined perimeter has 438 vertices and exactly two reported
nonadjacent projected crossings:

| Root | Fringe edge | Hood edge | X, Z (mm) | Conflicting Y (mm) | Difference (mm) |
| --- | --- | --- | --- | --- | --- |
| Left | 2914–8565 | 8029–841 | −58.0870, 135.9000 | −30.5990 / −29.2649 | 1.3341 |
| Right | 11283–3612 | 793–8088 | 58.0856, 135.9000 | −30.6035 / −29.2713 | 1.3322 |

The expansion moved the conflict from the old upper root crossings to the
lower root perimeter; it did not remove incompatible chart constraints.
A solver cannot satisfy both Y values at either crossing. This rejects
the authorized patch, not all possible chart-based surface constructions.
No further ring growth or convenient-chart search was attempted.

## Limits and handoff

No inner/outer orientation, seam-angle, thickness, intersection, material,
or visual candidate checks ran because no surface state was constructed.
The probe verified the source hash after inspection and that no objects
were created. Actual tip/opening preservation was not newly certified by
this failed chart test; the JSON enumerates the small additional original
fringe boundary segment affected by the proposed support, for root review.

Any later merged-shell approach must also qualify off-support appearance:
the existing hair uses Generated texture coordinates, so merging objects
changes its implicit mapping domain despite unchanged material nodes.
Shading normals can change too. No shader or material changes were made.

Root must decide whether to authorize a genuinely different retained
boundary construction. The present authority does not permit another
expansion. Foot033 remains retained; canonical state is root-owned.

Evidence hashes:

- `head_035_domain_probe.py`:
  `b70ee52ada70375f547b70e51aee1a4fc565d8c523164adf3c6dc0e66d48828b`
- `head_035_domain_probe.json`:
  `ee3fb144c2eaf9cbe99c6b9160763063e31f53d618aa220da6826cad8666fa9f`

Work remained sequential as requested: the perimeter result directly
controlled whether any construction was authorized. No delegation occurred.
