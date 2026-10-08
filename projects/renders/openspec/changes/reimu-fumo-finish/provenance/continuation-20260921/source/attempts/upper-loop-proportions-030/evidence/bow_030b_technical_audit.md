# Frozen bow030b winding delta

The saved030b correction makes `A42 Left constructed bow loop` outward
without changing its geometry. All 85 other visible geometry controls,
recorded rig/rest pose and appearance match independently reopened030.
Existing030 geometric contact evidence carries; no contact rescan was needed.
This is not a visual, animation or whole-asset acceptance result.

Source030 SHA256:
`4bf89ee268361802c4f0d778c470769e0a7201e9ee90282a96bd24815877072b`.
Candidate030b SHA256:
`d69f0325355fc767bccb98f75affee4b70106dbd3ac5e488ae0a70ad0f9de2a6`.
Both clean-reopened with pinned Blender 5.2.1 LTS, build `9e2066aef7ef`,
background and automatic file scripts disabled. Both hashes remained intact;
the audit performed no save or render.

## Exact delta

- All base and evaluated world coordinates are bit-identical; maximum
  displacement is zero. All 1,900 unoriented triangles, including
  multiplicities, are identical: zero added or removed triangles.
- All 950 base/evaluated polygons and all 1,900 oriented triangles reverse
  winding up to cyclic starting-corner rotation. No quad diagonal changes.
- Material regions, weights and modifier stack are unchanged. There are no
  UV layers in either file. The material `Bow proxy red.002` remains a valid
  persistent database ID with no missing face bindings.
- The unchanged modifier stack is A154 lattice, rig armature, global022
  lattice; there is no Solidify or normal-derived position offset.
- The 85 controls match both source030 and the writer receipt. Recorded
  appearance and rig state also match source; no rig refit is claimed.

The loop remains closed: 952 vertices, 1,900 edges, 950 faces, Euler 2;
every edge has two incident faces. Both base and evaluated meshes have zero
nonfinite vertices, degenerate faces below 1e-14 m², zero-length edges below
1e-10 m and inconsistent shared-edge winding. Evaluated signed volume changes
from -2.94324742165e-5 to +2.94324742165e-5 m³, confirming the outward reversal.

## Contact carry and limits

The exact triangle multiset, exact world coordinates and exact surrounding
85 controls preserve030's unsigned separation, edge-crossing and 116
knot-root witness results. Consequently its 28 scoped lowered-patch tests
still have zero crossings; dimensions and knot-root positions are unchanged.
This carries the original bounded evidence, not an exhaustive collision
guarantee. Surface-normal signs and shading deliberately change and are not
claimed equivalent. No additional global geometry was audited.

## Machine evidence

The complete compact delta in [machine evidence](bow_030_machine_evidence.md), SHA256
`18e1c7f7afe2c5c0605ceafc1847420608e4e34df50d9805bd8378255c1c1e3f`,
contains the complete compact delta report and binds the carried030 audit
(`7a0cbf817b55558b090b8b04e9dff34c0b3e4f4d565fe0f0177b1137d052b51b`)
and separate root witnesses
(`5c2215d133f68a73cd3cf3fa577930ad540a77bbbc9f80ac611d9ee27b8162fe`).
The separate raw root witnesses remain ignored and hash-bound; the compact
delta and carried030 summary are included in the canonical machine packet.
