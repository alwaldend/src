# Plan

Use the exact untouched coupon from `mapped-gui-sculpt-009` to test whether
the now-working native event route supports useful macro authoring.

1. In the mapped XWayland Blender 5.1 authoring host, clean-open the untouched
   coupon and retain its coordinate baseline.
2. Activate Blender's actual essential Grab-like Sculpt brush through the
   asset system; do not substitute direct mesh or `brush_stroke` edits.
3. Project the target surface into the live region and schedule one
   press/move/release sequence across timer ticks so the modal operator
   receives intermediate motion.
4. Require a visible, localized silhouette displacement affecting more than a
   trivial three-vertex neighborhood. Save the changed file before undo.
5. Clean-open and render the exact changed bytes with repository-pinned
   Blender 5.2.1, then verify native undo and unchanged source bytes.

Stop if the brush asset cannot be activated, timed native events do not make a
macro silhouette change, or pinned Blender cannot reopen the saved candidate.
Correct only a newly evidenced mechanism; do not repeat an equivalent stroke.
