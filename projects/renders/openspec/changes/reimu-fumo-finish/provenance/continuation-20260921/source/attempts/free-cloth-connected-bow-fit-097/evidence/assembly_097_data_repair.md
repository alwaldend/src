# 097 sole exact-surface knot data repair

First build stopped at native knot hygiene before pile/save. Both free hair
shells already passed native surface and all receiver contacts; no new blend.
Original arrays304f882b760abee30f0ab8b836fb2bcd538271ecd2e6f2b8fe375327e5a6f00a,
preflight50260a3b5a2da92349f4a28259db345078f91486a2e96879e3214d0ff8d29719.

The retained knot's source volume is negative -1.94412765989e-6m3: its winding
is inward. Its two constant-Z96-vertex cap polygons are convex, not duplicate
pole rings. Ring radii3.596–3.855um; collapse would change the surface and is NOT
used. All96 subgate-area triangles belong to the caps; none is exactly zero.
Their minimum twice-area2.894e-15m2 is below the frozen1e-14 numerical gate.
A center fan in the SAME cap polygons raises minimum twice-area to7.747e-13m2.

Use the one allowed data repair: retain all4704 original coordinates exactly;
reverse each noncap native triangle explicitly, preserving its surface; add
one interior point per planar cap and triangulate the same convex polygon.
Preserve directed boundaries, noncap ancestry and all original attachment
triangles (verify none is a cap). Full native surface gates stay unchanged.
Correcting winding also makes normalY<0 select the actual front; recompute the
same frozen fit and dependentD, not the first run's orientation-bugged number.
Independent decision review agrees this is representational hygiene, not a
new shape parameter or geometry strategy. No pole collapse or extra retry.

Also complete the intended source-attachment gate by recording canonical bow/
knot triangle-vertex pair identities and requiring exact equality after the
shared rigid translation. Protected-band membership alone was insufficient.
This adds no geometry change. Preserve original builder/log/arrays and use
fresh repair output names. Original builder614fcc61b9f99d259ed43be1e78a3b3c3eeb18efd70ad4e5824d1de96c794453;
pure knot diagnostic4aeba5956e26ec742b213bbc91188c8c618c28deb9e16f487282fa1d09240132.
