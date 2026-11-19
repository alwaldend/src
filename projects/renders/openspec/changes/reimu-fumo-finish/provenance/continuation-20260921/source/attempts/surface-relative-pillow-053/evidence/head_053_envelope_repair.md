# 053 repair 2: sampled support envelope

The corrected source guard passed and pinned Blender opened unchanged 050.
The first front-ray failure was X -0.0040153223, Z 0.1949311048 m; no candidate
was saved. Read-only diagnosis disproved the initial stale-outline hypothesis:
actual front and traced outline share the same extrema, and the exact failed
X is an outline intersection at that Z. The defect is in the generated
receiver's sampled envelope, not in those input bounds.

The angle-height sampling misses the narrow asymmetric crown contour when
the inherited head radius clamps to 1 mm above Z 0.1945. A union radius at
queried scanlines does not guarantee the polygon between sampled lines covers
the controlling contour, especially after forced exact pole closure.

Causal change: construct the convex support perimeter from the actual front
XZ points and the supporting head profile, with 0.9 mm outward padding. Use
that perimeter directly for the shared radial pillow rather than resampling
the steep crown into an inferred symmetric width field. Preserve every front
vertex's XZ, the strict all-support-rays gate, width/depth limits, and all
non-target source data. The next run should cover the previously missed
crown point and all remaining front points. This consumes repair 2 of 2;
stop this candidate if another guard or representation gate fails.
