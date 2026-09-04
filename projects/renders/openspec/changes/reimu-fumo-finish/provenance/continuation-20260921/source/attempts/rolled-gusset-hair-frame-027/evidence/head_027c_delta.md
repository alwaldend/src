# 027c coverage/continuity correction

027b fixed transient material IDs and the original-disabling crash, but its
fixed images show a pale crown strip and jagged temple depth transitions.
Root rejects027b as a retained candidate. Frozen file SHA
42c4eb962e6edcf9a496b50da582fabceea89d2e674d022d9d07568c3ffc03cf.

Causal source findings:027 removed the old structural head's brown frontal
underlayer, but its new hood covered only rear and gusset faces. The old
fringe did not cover the entire crown in3D. Moreover, reprojection changes
from ray-hit depth to nearest-point displacement at unsupported overhangs,
despite that boundary not being a cloth seam.

027c adds the old front-hair face region to the independent hood, using the
same new core surface and hood thickness. It reconstructs fringe depth from
one continuous analytic panel field over the existing XZ contour, with.85mm
stand-off and existing variable padding. It does not inherit the old mixed
ray/overhang depth residual. No new head core shape, landmark XZ changes,
non-head changes, material-node edits, or acceptance relaxation.

Prediction: crown skin is covered by brown cloth, and side fringe no longer
has ray-boundary spikes. Reject if layers expose new gaps/intersections or
the overall head still regresses. Root is requesting independent pixel
comparison of027b and a source-aware counterexample before retention.
