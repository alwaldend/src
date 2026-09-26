# Hair 077 causal diagnostic

One pinned Blender 5.2.1 LTS (`9e2066aef7ef`) process replayed the exact existing builder prefix, stopping before `# Existing strands split`. Exit 0; no model save, parameter variant, repair or pile construction. Source and builder hashes remained unchanged.

- Source 076: `c7aeaf157f7d451d658c302c6a9300125ab145a7551b5fa8713288052c747050`.
- Replayed 077 builder: `c088ef61980b2da9a92d56363b4a8d18a0190fb4ab0eb82f2124edf0fb1c7278`.
- Raw diagnostic JSON: `8785e437ad03e61a77d5ab3d4877f797d7f2ca9ee6d0be3b5f3d086bcbcd620c`.

## Two independent defects

### Clipping identity and seam topology

The raw midsurface has 4,750 vertices and 9,240 triangles. Every edge has one or two incident faces, but boundary vertex 3987 has degree four and two disconnected face fans: one with eight faces, the other with face 7289 alone. Bridging its boundary produces the single nonmanifold shell edge 3987–8737 with four rim faces.

Witness, world millimetres: outer endpoint `(24.904089, 30.480357, 105.181783)`, inner endpoint `(24.580274, 30.084034, 105.582498)`.

Original head triangles 10894 and 10904 share source edge 3763–4315, yet their clipped results belong to separate fans. This supports lost intersection identity during independent float interpolation/coordinate-key welding, not an intended disconnected material island.

There are 70 raw triangles with minimum altitude below 1 micron, 76 below 10 microns. The narrowest has altitude 0.867 nanometres and a 0.931-nanometre edge, versus a 2.087-mm longest edge. Its two near-identical clipped points differ only in X: `-14.000000432` and `-14.000001363` mm. Both are at Y32.342646271, Z104.999996722 mm. No triangle area is at or below 1e-16 m²; that area-only threshold would therefore miss the pathological slivers.

### Real thickness/core crossings away from the clipping seam

Raw midsurface, outer surface alone, inner surface alone, and head core each have zero non-shared-vertex self-pairs. The closed cover has 18 pairs, all resolved noncoplanar inner/outer crossings. None involves a rim. All involve outer triangles 9066 or 9071 in the rear-right patch near Z127 mm; the patch spans approximately 20.6 mm in X. Maximum resolved crossing span is 1.927 mm, with witness `(25.619680, 36.216714, 127.121799)` mm.

All 26 head/cover pairs are resolved noncoplanar core/inner-shell crossings near rear Z127 mm. They involve only inner triangles 18187/18193 on the left and 18306/18311 on the right. No involved head triangle has a protected graphic vertex. Maximum crossing span is 5.556 mm, witness `(23.377766, 36.498792, 127.184032)` mm. Spans describe intersection length, not penetration depth.

These contacts are over 20 mm above the clipping defect and cannot be eliminated by welding its rim. Individually simple skins do not establish an ordered shell: the independently displaced, piecewise-linear radial layers cross each other and the core.

## Smallest supported construction remedy

1. Make clipping topology-preserving: retain original vertex IDs, cache each crossing by its canonical original edge ID, and reuse source endpoints consistently when they are on the clipping plane. Eliminate duplicate/collinear polygon vertices before triangulation. Require one connected fan per retained vertex and degree-two boundary loops before adding shell bridges. This addresses the observed seam/sliver cause; splitting the single degree-four vertex alone would conceal the underlying crack.
2. Locally rebuild the long rear triangles implicated above, using a conforming, well-shaped patch for the outer, inner and core layers. Construct inner and core placement from that shared patch with explicit triangle-interior clearance ordering, rather than independent constant radial subtractions. Preserve the unaffected envelope and graphic vertices. Clipping repair alone or an untested global thickness reduction is not an evidence-backed solution to the rear contacts.

No alternate construction or thickness was tested, so no proposed remedy is claimed collision-free. The diagnostic excludes self-pairs sharing vertex IDs and does not establish continuous maximum penetration, visible severity, other-hair contacts, pile correctness or visual acceptance. All raw pair witnesses remain in `hair_079_diagnostic.json`.
