# 071 reveal the existing chin by correcting upper-body occlusion

User: "its chin looks receded, kind of weird. Continue". Root is sole model
and canonical-goal writer in the verified linked task worktree. Use exact
provisional070 a3f025b21b74104a29eab2c2a5afb0da662c9ca702732ecd9ee0b0e84366d27c,
new chin_071_candidate.blend only. No acceptance inherited; all earlier
unresolved construction/material failures remain. Preserve source bytes.

## Diagnosis and decision review

Independent image-only chin_071_reference_diagnosis finds the lower face tucked
back but no defensible excess chin-center depth: frame10 has similar cheek
rollback and side hair occludes the chin. Its stronger evidence is a short
visible mouth-to-underside span and red crescent above the collar.
Primary control physical_front; canonical front corroborates spacing; frame10
controls visible cheek profile with yaw/occlusion uncertainty. Approximate
image spans0.07–0.09Wh candidate versus0.10–0.12Wh reference are not calibrated.

Root read-only chin_071_before_profile.json and chin_071_occlusion.json bind
exact070 and identify the causal interface. At X0,Z82.5mm cream chin is
Y-18.248mm while red upper body isY-25.746mm,7.498mm in front. At X10mm the
red core is6.410mm ahead. AtZ85mm head is already frontmost. The stored pattern
has40px mouth-to-cream-underside (~0.109Wh); the visible shortening is partly
occlusion, not proof of missing chin volume. This disconfirms immediate large
forward extrusion and weakens the image-only downward-extension proposal.

Root verdict PROCEED with a one-object interface correction first. Alternatives:
global head pitch/inflation changes correct upper landmarks; local head
extension may remain hidden by the same red core and may overlengthen the
face; leaving the core preserves a proven foreground occluder. Tucking its
upper front behind the existing head is the smallest causal correction.
Reopen this verdict if fixed pixels retain a recessed underside or the tuck
creates a shoulder/collar gap. No claim that this solves every head issue.

## Frozen edit and controls

Change only mesh coordinates of Macro038 leaning chest and seated hips.
For each world point(x,y,z), add18mm*wx*wz*wy toY, where
wx=1-smoothstep((abs(x)-13mm)/15mm),
wz=smoothstep((z-75.5mm)/5.5mm),
wy=1-smoothstep((y+25mm)/50mm), with each input clamped0..1.
This is a tucked upper-neck/core region, not a translation of the torso.
AllZ<=75.5mm and |X|>=28mm remain exact. X/Z, topology, materials, transforms,
modifiers, head/eyes/mouth, hair/pile, collar, bow, sleeves and lower body remain
exact. The Y map has derivative>=0.46, avoiding a folded-over mapping.
18mm is a bounded displacement needed at the deeper lower seam, not a
reference-measured anatomical chin projection. Maximum frontmost head depth
and every head vertex stay unchanged.

Pre-save: finite closed body mesh, positive volume/no nonadjacent surface
overlaps; exact immutable head/graphics and all non-target fingerprints;
front ray order at X0/10mm,Z82.5mm must become cream head ahead of red body.
Record remaining red/white/head order at nearby rows rather than concealing
an exposed seam. If geometry preflight fails, stop before save. One causal
pre-save data/API repair allowed, no parameter sweep or post-render variant.

Pinned5.2.1/build9e2066aef7ef background writer, factory startup, automatic
scripts disabled,4threads, Python exit2; task-local config/cache/TMPDIR.
Save once, clean-reopen exact candidate, render fixed512px front/side/3q.
Compare visible chin span and neck crescent while preserving fixed cameras,
shader/light settings and upper-face landmarks. Nominal photo alignment is
diagnostic, not a calibrated criterion002 pass. A plausible local fast pass
requires independent blind pixel review and the full fixed packet before
module retention. Whole asset remains unaccepted. A bounded technical verifier
may check only changed core, protected head/collar and target occlusion/contact.

## Checkpoint routing

The optional portable-plan registry reached its64-entry bound before mutation;
the plan-only CLI call failed and RV240 remained unchanged. Read-only source
inspection confirms the supported compact checkpoint accepts an attempt-local
frozen plan without another portable PlanID. Use that path with this complete
plan and criteria bindings. Do not delete history, relabel a rejected plan,
raise tool limits, or hand-edit canonical manifests. This narrow workflow
choice preserves the user's goal-priority instruction and all evidence/gates.
