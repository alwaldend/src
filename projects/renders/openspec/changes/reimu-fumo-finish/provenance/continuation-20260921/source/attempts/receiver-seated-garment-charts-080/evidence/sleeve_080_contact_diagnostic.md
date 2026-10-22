# Sleeve 080: root-chord discretization diagnosis

One pinned Blender 5.2.1 LTS (`9e2066aef7ef`) process exited 0. It executed the hash-bound wrapper AST and unchanged construction through its contact loop, stopping before the final failing assertion, exception handlers, file writes and save/receipt code. No parameter variants, retries, saves or renders. All listed inputs remained unchanged.

## Result

This is a correctly seated vertex boundary whose straight chords cut through the arm, not evidence of a wrong free-shell offset.

- Each sleeve has exactly 18 arm intersection pairs: nine inner-surface pairs and nine shoulder/root-bridge pairs. All are resolved noncoplanar crossings with strict signed plane straddling beyond 0.1 micron, not BVH coplanarity noise.
- Only four sleeve triangles are implicated per side: 7757, 7759, 15437 and 15439. Their conservative chart extent is `u=0..1/48`; no outer-surface, cuff-bridge or farther-free-sleeve pair occurs.
- All 80 inner-root vertices are outside the arm at approximately 100 microns. However, two of their 80 edge midpoints are inside: edges 3958–3959 and 3959–3960. These are the back-chart root arcs `v=0.95–0.975` and `0.975–1` at `u=0`.
- Four implicated triangle centroids per side are outside: minimum signed clearances 216.93 microns left / 222.82 microns right. Vertex/centroid-only checking misses the edge-localized crossings.
- Among ten implicated edge midpoints per side, two are inside. Maximum midpoint penetration: 60.521 microns left / 60.518 microns right.
- A denominator-eight barycentric lattice on the four contact triangles increases the sampled maximum to 65.617 microns left / 65.615 microns right. Every inside lattice/midpoint sample lies on `u=0`. Total contact-neighborhood coverage is 201 samples per side, including duplicates on shared edges; 18 inside classifications are not 18 independent penetrations.

Three skew-ray parity votes agree on every sample classification. The independently retained nearest-normal signs agree at the reported penetration witnesses.

## Witnesses and bounds

Deepest sampled point, world millimetres: `(±17.201191, -3.030835, 59.851140)`, on root edge 3959–3960. Left nearest arm point: `(-17.154757, -3.042163, 59.806183)` mm; signed depth -65.617 microns. Right is mirrored within float precision.

Worst midpoint: `(±17.226342, -3.392835, 59.980743)` mm; approximately 60.52 microns inside. The adjacent chord midpoint near `(±17.303528, -0.767318, 59.464991)` mm is approximately 2.31 microns inside.

Conservative bounds of the implicated sleeve triangles, world millimetres:

| Side | X | Y | Z |
|---|---|---|---|
| Left | -17.4813 to -16.5599 | -5.3556 to 0.4102 | 58.2297 to 60.4992 |
| Right | 16.5599 to 17.5972 | -4.8408 to 0.6498 | 58.1909 to 60.4992 |

Maximum actual intersection-segment spans are 0.8374 mm left / 0.8540 mm right. These are intersection lengths, not penetration depths.

## Causal repair direction, not implemented

Target the root chords, not the free chart/thickness system. A receiver-conforming refinement of these two root arcs, with consistent adjacent inner/root-bridge topology, addresses the demonstrated cause while retaining intended root seating. If the 80-point topology must stay frozen, inner-root allowance is the causally relevant parameter: the existing construction has at least a 65.62-micron sampled clearance deficit between endpoints. No specific replacement allowance is proven safe by this run; validate complete edges/faces, not only the reseated vertices, after any authorized repair.

This is a bounded actual-arm audit. Existing garment contacts remain unresolved and were not independently analyzed. Sampling does not bound the continuous maximum penetration. No visual acceptance is implied.

## SHA256 receipts

- Source 076: `c7aeaf157f7d451d658c302c6a9300125ab145a7551b5fa8713288052c747050`.
- Wrapper 080: `dab097ff5d3e923a0716fc28da7288811ae5acf6c350dc115d9e69a4ef88870a`.
- Geometry helper 080: `6eb29487b7862b88bd4972b1679956315dc8da0841c72b37fea0ac0ead3e2c3b`.
- Diagnostic code: `12e93b8700764820181a4ef269088a3e35ef426b478cfc29c61f7386dbc3fc2d`.
- Diagnostic JSON: `7eaa440a4248de93cc67e4bc79a600f1714cca2908602763f9a82cfd4b43a83c`.

The JSON preserves every contact pair and sampled witness. Its diagnostic code hash binds the executed script; the report hash is supplied separately to avoid a self-reference.
