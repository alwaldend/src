# Shared Head034 diagnostic: reject folded joins, not the mirror cut

Root's rejection is supported by the local endpoint geometry. The mirror
row175 reversal is a projected wavy welded boundary, not the old competing
sheet mechanism. A separate set of near-antiparallel fringe/bridge joins is
the decisive failure.

This was one specifically authorized diagnostic reconstruction of unchanged
helper SHA256
`c9976b58dc310001b4478ddb399b6691bd801eb3a725b12358c6422000cb1c69`
on frozen032 SHA256
`6d2d6c52a499d056f9d5a4e0fdbca53fe7588ac125d91c449d07c7fa72d3cab8`.
No helper/shape parameter changed, and no blend was saved or rendered.
Source and helper hashes remain unchanged. Pinned Blender 5.2.1 LTS,
build `9e2066aef7ef`. Full new evidence:
`head_034_shared_diagnostic.json`; probe source has the same basename `.py`.

## Mirror175 is a real cut crossing, not coincident outer sheets

The actual shared edges intersect the horizontal pixel-center row at:

| X crossing | Candidate edge vertex IDs | Adjacent triangle IDs | Geometric normal jump |
| ---: | --- | --- | ---: |
| 236.740281 | 37333–37334 | 55159 / 80861, hood/bridge | 10.573 degrees |
| 248.402866 | 496–4023 | 3720 / 74866, fringe/bridge | 17.497 degrees |
| 249.311879 | 2746–8258 | 3994 / 74876, fringe/bridge | 7.486 degrees |
| 250.022307 | 4032–8275 | 4128 / 74882, fringe/bridge | 24.477 degrees |

These are shared topological edges of the same outer skin, not two
unconnected surfaces. The exact hits explaining the region reversal are:

- Pixel249: retained-fringe triangle/polygon3855, vertices8241/7241/514.
  XYZ=(-36.393702, -31.151026, 176.235750) mm; camera-ray depth687.333381 mm.
  Geometric normal=(-0.390347, -0.757161, 0.523772).
- Pixel250: bridge triangle74882 / polygon74779, vertices4032/8275/38422.
  XYZ=(-35.837114, -31.400919, 176.234275) mm; depth687.550322 mm.
  Geometric normal=(-0.480437, -0.651060, 0.587624).

The next distinct surface is inner skin, respectively 1.438863 mm and
1.555515 mm farther along the camera ray. There is no intervening competing
outer sheet at these two rays. All precise hits from X237 through252,
including shading normals, are in the JSON. At pixel251 the stepping ray
reports the same triangle twice within approximately0.1 micrometer; this
is a numerical repeat, not another surface layer.

This does not make the local shape perfect: the crossing edges still have
finite geometric normal jumps. It does invalidate treating a repeated
region name on this welded surface as proof of overlap.

## Independent decisive failure: folded fringe endpoint joins

Worst measured upper join:

```text
Candidate edge: vertices 2567–7103
Retained fringe: triangle19746 / polygon19746
New bridge:     triangle75486 / polygon75081
Normal jump:   179.978442694 degrees

Endpoint A, mm: (55.602055043, -33.960852772, 147.183015943)
Endpoint B, mm: (55.602058768, -32.958142459, 148.144647479)

Retained normal: ( 0.943737090, -0.228899479,  0.238673896)
Bridge normal:   (-0.943793654,  0.228787825, -0.238557503)
```

Other upper examples include edge7807–37100 at Z145.571–145.584 mm,
179.516 degrees, and edge2519–37174 at Z152.959–153.984 mm,
179.466 degrees. Exact adjacent IDs, points and normals for every measured
join are in `all_seam_edges`. Some resampled edges are very short; the worst
edge above spans about1.389 mm, so the finding is not based solely on a tiny
longitudinal sliver.

Above145 mm, fringe-to-bridge joins: 388 samples, median6.547 degrees,
95th percentile40.779 degrees, maximum179.978 degrees. Bridge-to-hood joins:
405 samples, median10.496 degrees, 95th percentile13.246 degrees,
maximum17.540 degrees. Tangent continuity was not achieved.

The candidate retained-fringe triangle is identified exactly. Its original
source triangle index was not carried into the diagnostic mesh metadata;
the earlier result establishes coordinate preservation, not that missing
index mapping. The source face on the removed side of this cut was not
measured, so do not claim that the original fringe had no crease here.
What is directly demonstrated is the unacceptable interface between this
retained triangle and the newly constructed bridge.

## Bounded intersection evidence

The test compared all 6,896 outer-bridge triangles with the reconstructed
shell using a BVH broad phase, followed by segment/triangle and coplanar-area
tests. Across 952 unique broad-phase pairs it found zero nonshared contacts
above the stated tolerances: shared vertex/edge contacts within0.2 micrometer
are permitted; coplanar projected overlap area threshold is1e-12 square
meters. This is not a whole-scene collision audit. No envelope-overshoot
audit ran.

No intersection does not imply a usable join: a connected surface can fold
back across its boundary tangent without intersecting another triangle.
Thus this negative intersection result does not override the measured
near180-degree seam reversal.

## Causal assessment and next decision, not implementation

The current helper projects the cross-strip chord into each endpoint's
smooth vertex-normal plane and normalizes the result to the chord length.
That constrains a plane, but does not constrain the derivative's sign
relative to the retained adjacent face. It also does not remove the
longitudinal component along the cut boundary. Normal-plane membership
alone is therefore insufficient on the stepped graph-distance cut.

The first-state topology check establishes opposite shared-edge winding.
With that winding, the nearly opposed normals at the worst join imply that
the bridge's first face goes back toward the retained side of the cut in
the local cross-edge direction, rather than departing into the removed
region. This is stronger evidence than merely observing a shading groove.
The actual chord and Hermite derivative vectors were not logged, so the
size of their longitudinal component is unavailable; excessive longitudinal
drift remains a plausible contributor, not a separately measured cause.

A boundary co-normal oriented away from the retained adjacent face is the
credible causal correction. Its sign should come from the actual oriented
cut edge and retained triangle, with a positive local cross-strip Jacobian
checked before construction. Review correspondence between the unequal,
stepped fringe and hood paths as well: changing the tangent direction alone
does not prove that the longitudinal mapping avoids twists or collapse.
These are conditions for a new root decision, not permission to run a
variant. No correction or replay followed this diagnostic.
