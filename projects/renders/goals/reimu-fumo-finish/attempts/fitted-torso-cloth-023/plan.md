# Narrow seated torso with attached cloth collar

Continue the existing project goal, generation 1, lifecycle 16, criteria 4,
from frozen body022b SHA-256
`96e6deea298308573174a35699ea4cf7b99e827260b2c108de43f8f0c1266014`.
022b is an intermediate checkpoint, not a stage pass. Preserve its corrected
global height, head shape and motion, bow, feet, floor, and source bytes.

Canonical front mid-bodice width is 0.486 Wh versus 0.709 Wh in 022,
uncertainty about 0.04 Wh; Wh = 0.1165 m. Collar-to-head samples show about
18.2 mm separation, consistent with the visible floating collar and tie.
The previous height correction cannot fix these X/depth mismatches.

Decision: proceed with a macro torso-width correction plus replacement of
the three detached collar/cravat solids with thin, surface-fitted cloth.
Leaving the old parts preserves a major reference defect. Translating only
the whole collar/tie would retain its thick polygonal shape and weak roots.
The compact stuffed core remains suitable after narrowing; its garment
outline must have a narrow waist and unchanged seated hem, not a barrel.

Append a non-destructive native lattice to the exact prior body role set,
excluding the three collar/cravat targets. X scale is 0.69 above Z 0.051 m,
1 below 0.028 m, with smooth interpolation between. Y/Z stay unchanged.
This narrows bodice/waist and draws sleeve roots inward while leaving the
low hem and feet intact. Replace Left soft collar, Right soft collar, and
Folded yellow cravat with two tapered soft cloth flaps and a gathered tie,
fitted to the narrowed body's evaluated front and the actual neck/head
surface. Initial collar span 0.055–0.059 m, tie width about 0.0135 m,
visible tie height about 0.025–0.027 m. These are construction hypotheses,
not acceptance facts. Root will review the draft before the sole save.

Risks: sleeve-root distortion, skirt/seat crossings, chin crowding, flat
collar cards, or a tie embedded into the torso. Check front, side, rear and
both three-quarter views; verify unchanged X/Y/Z controls where claimed.
Target torso width within 0.03 Wh of the canonical target; attachment roots
must be visually seated and have small measured surface gaps, not merely
overlapping bounding boxes. Any remaining major shape/contact defect fails
the stage. Retention may only be an intermediate module decision.

Root is sole writer using pinned background Blender 5.2.1. Reuse the settled
native lattice capability, normalizing actual base-coordinate extents.
No GUI, listener, add-on, host configuration, or old output is changed.
One causal source-bug correction may be recorded within this hypothesis;
a failed geometric representation requires a new plan or reconstruction.

Parallel work: root owns width field and integration; visual_direction
prepares an unexecuted collar/tie draft. After freezing, independent pixel
review and technical checks may run concurrently. No parallel model saves.
Fixed review contract remains
`4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`.
