# A91 whole-plush bakeoff

## Bindings

- Goal: `reimu-fumo`, generation 1, lifecycle generation 4, criteria r1.
- Start checkpoint resource version: `58`.
- Work type: `candidate`.
- Canonical tracked asset is protected and remains unchanged:
  `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`, SHA-256
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Exact-variant front controller:
  `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png`,
  SHA-256
  `864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`.
- Exact-variant depth and layer controller:
  `projects/renders/blender/fumo/reimu_fumo/references/canonical_turn_180.gif`,
  SHA-256
  `0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`.
- Frozen high-water whole baseline: rung003, SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- Frozen partial construction evidence: A88 S07, SHA-256
  `0a30e2af3142081648bb3137ad75d6d1cc73de55e9f830f85bad1f85e92c8788`.

## Stable defect and causal reach

The latest artifacts lost the complete subject while optimizing head and hair
interfaces. Rung003 preserves a recognizable whole Reimu Fumo but reads as a
rigid, simplistic figurine with incorrect head/body proportions and weak
constructed-fabric cues. A whole-form blockout can change the required subject
silhouette, mass hierarchy, pose, and construction category; another isolated
hair edit cannot.

The first discriminating artifact is a fixed front plus the worse of the two
three-quarter views at 512--640 px. It must show the complete seated plush and
all identity-defining masses. Stop either branch immediately if it loses the
complete subject, reads as human anatomy or a plastic mascot, clips the skirt
and feet, or cannot plausibly beat rung003 without a representation rewrite.

## Parallel candidate branches

1. **Fresh construction-aware whole blockout.** Build a complete neutral plush
   from directly controllable soft cushion and fabric-panel forms. Use the
   canonical front for the large silhouette and the canonical turn plus
   physical references for depth and manufacturing logic. Do not derive this
   candidate from rejected generated head fragments.
2. **Rung003 salvage control.** Copy the exact rung003 source and make only the
   smallest coupled macro changes needed to correct its plush mass hierarchy:
   shallow broad cushion head, compact seated torso, pooled skirt, tucked foot
   pods, and coherent hair/bow panels. Preserve recognizable graphics and the
   complete subject. Reset rather than expanding if the inherited structure
   prevents those changes.

Each worker owns only its isolated `out/.../candidate_*` directory and may not
write canonical goal state or the tracked asset. Both use all relevant frozen
references and the same fixed front and three-quarter camera contract. Target
first pixels within 30 minutes of active authoring; the render packet itself
should remain under 60 seconds. Do not build detailed seams, materials,
retopology, rigging, or six-view promotion evidence before this gate.

## Decision gate

When both first artifacts exist, open one correction episode bound to this
attempt, the stable whole-form defect, rung003 as the calibrated historical
comparison, and the two exact candidates. Route exactly one fresh independent
reviewer without implementation details. Require whole-subject recognition,
plush-medium read, macro silhouette, construction, and Reimu identity. A
survivor must score at least 6/10 internally and clearly beat rung003; this is
only a provisional parent threshold and does not weaken the final 8/10 gate.
Tie favors the simpler fresh construction-aware candidate. Allow at most one
causally named correction to the selected branch.

## Criteria and regressions

- Directly affected: criteria 001--004.
- Protected now: criteria 005--008 remain unverified; no claim is made until a
  whole-form candidate survives.
- Fixed regressions: complete subject in both views, no major clipping or
  floating, correct small-foot/no-human-leg read, coherent rear hair coverage,
  and unchanged protected input hashes.

## Process control

The work is intentionally limited to two independently reviewable hypotheses
because more branches would exceed immediate review capacity. Record worker
start, first-pixel, and decision times. If neither branch survives, close this
attempt as a reset and use the visual evidence to define one smaller sculpt
control coupon; do not return to a larger analytical generator.
