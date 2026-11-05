# 085 failed cloth construction; protected working model unchanged

Root verdict RESET. One combined two-sheet48-frame solve completed in177s,
then failed pre-save on actual outer self-overlap candidates. No085candidate,
no final target mesh replacement, no renders. Source084 and donor080 unchanged.
Helper695b1f727f72a197540e52135cd0fe288ac8264eca8da44511edf565695d173d;
builder5bb874c745c70d9ad5ccfdd0a147f240d17bd8caee4fc1ae6b644a96c6d65eac.

Left/right outer self-pairs8930/14734, shell49369/75721; new outer/arm160/250.
The outer inputs had zero self pairs but617/576 torso crossings, including17
pinned-incident pairs each, and368/286 skirt crossings. Six root vertices per
side participate in initial torso crossings. The local nearest-normal root
depth estimate is10.395mm in the torso; this finite sign method is not a global
inside proof. Body pairs themselves establish conflicting initial contact.
Other open-surface nearest-normal flags are not reliable penetration claims.
No collision metric or gate was relaxed to save tangled cloth.

Cuff Z ranges41.03–59.26mm left,44.72–59.57mm right, Y widths21.56/21.31mm.
Fixed outer-root error before roundoff snap10.5/13.4nm. Both160-point root
layers stayed exact, but that preservation included the problematic buried
torso attachment. Thus good arm clearance alone does not define a valid cloth
initial state. Preserving that embedded root while enabling the full body
collider supplies incompatible constraints near the attachment; immediate
short-rest contraction and pre-crossed skirt are additional unisolated causes.

No API/data error has been identified, so the unused one-repair budget does
not authorize solver-dose or shape tuning. Retire this initial-state setup,
not native cloth generally. Next sleeve design must start outside the body/
arm union and avoid initially crossed skirt, with controlled cuff support.
Root placement near the visible canonical shoulder must be reviewed, not
silently shifted to hide failures. A worker is reviewing that bounded choice.

Meanwhile inherited rear-head marks have source evidence of copied facial
graphic triangulation on the back; one independent read-only measurement is
pending. A separate086plan may address that causal defect while preserving
084's retained lower-hair improvement. Whole goal remains open and active.
