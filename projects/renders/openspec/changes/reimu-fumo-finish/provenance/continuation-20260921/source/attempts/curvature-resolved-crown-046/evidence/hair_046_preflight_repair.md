# Adaptive projection preflight repair

The first046 writer stopped before save: after eight adaptive passes, maximum
projected edge remained2.673mm, above unchanged2mm gate. Projection at the
tangent near the crown boundary expands new chords. Increase the adaptive
pass ceiling to14, retaining the1.8mm split trigger and2mm acceptance limit.
Prediction: repeated geometric bisection resolves the remaining tangent
edges; if14passes cannot meet the same bound, stop this projection method.
No candidate was saved and no failed pixels are accepted.

The14-pass test did not converge: after pass8,136 long edges are repeatedly
recreated while vertex count grows by136 each pass. Worst endpoints are
near the crown tangent. The14-pass Blender subdivision/retriangulation method
is stopped. Replace that primitive with explicit shared midpoint red/green
triangle refinement, which retains split edges rather than reintroducing
long diagonals. Preserve the same projection,1.8mm trigger,2mm gate and
14-pass stop. This repairs the preflight algorithm, not the physical model
or its acceptance standard. Record per-pass edge maxima to prove convergence.
