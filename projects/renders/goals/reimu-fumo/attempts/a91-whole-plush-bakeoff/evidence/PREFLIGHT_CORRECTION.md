# Preflight correction

The published A91 plan incorrectly described branch B as a rung003 salvage.
The controlling A90 recovery decision instead requires an exact A88 S07
whole-context control while rung003 remains a read-only visual comparator.

The rung003 worker was interrupted before producing decision-bearing evidence.
No rung003 or tracked file was modified. A91 continues only as the fresh whole
blockout branch. The correct A88 S07 control runs in the disjoint
`out/reimu_fumo_attempt_092_a88_context_control/` workspace and will be
recorded as A92 rather than silently rewriting the immutable A91 plan.

This early return is evidence that the correction gate prevented a known-wrong
branch from consuming the bakeoff budget.
