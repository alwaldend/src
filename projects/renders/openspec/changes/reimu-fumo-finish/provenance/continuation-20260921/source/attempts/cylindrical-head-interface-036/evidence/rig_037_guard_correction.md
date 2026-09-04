# Parent-state guard correction

The initial helper e757680d stopped before surface capture because it
incorrectly required an unparented source. Existing preflight JSON records
the sleeve's parent as ReimuFumoRig and world matrix as identity.
The corrected helper 0ed12641 requires that exact rig parent. It changes no
geometry, coordinate or normal code. The initial helper remains preserved
as rig_037_coordinate_helper_initial.py and the failed execution log is
rig_037_render_probe.log. A second execution is authorized only to get past
that identified contradictory guard; the predicted difference is reaching
the coordinate capture and neutral checks, not a changed candidate shape.
