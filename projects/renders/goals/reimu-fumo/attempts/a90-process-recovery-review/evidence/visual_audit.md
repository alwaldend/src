# Implementation-blind visual audit of recent Reimu Fumo work

Observed: 2026-09-01

## Scope and evidence

This review judges the rendered pixels against the supplied references before
using topology, object names, or author intent as evidence. Candidate names
and their declared neutral-sculpt/module stages were visible, so this is
context-light rather than perfectly anonymous. Local measurements are used
only after the absolute visual decision.

Controlling evidence:

- exact-variant front and 25 cm scale:
  `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png`
  (`864b597117c79e5556fcf360333a798584ed6964e0fdcfe97e002a34013ed63c`);
- exact-variant depth, layer order, and hidden silhouettes:
  `projects/renders/blender/fumo/reimu_fumo/references/canonical_turn_180.gif`
  (`0d774eaa7f75828e388df4fb886cda7c563ce3bcd4ccb38d9885997a0846af30`);
- selected canonical turn frames:
  `out/absolute_x_pixel_review/canonical/frame_003.png`,
  `frame_010.png`, `frame_011.png`, `frame_017.png`, `frame_018.png`,
  `frame_019.png`, `frame_025.png`, and `frame_026.png`;
- construction and material support:
  `projects/renders/blender/fumo/reimu_fumo/references/physical_front.png`
  and `physical_side.png`;
- consolidated reference sheet:
  `out/a74_measurement_replay/reference_measurements/all_relevant_reference_packet.png`.

Candidate evidence:

- protected complete-character rung:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/five_view_gate/contact_sheet.png`;
- A85 S12:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a85_review/A85_S12_six_view.png`;
- A87 S04:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/packets/s04/`;
- A88 S07:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/`;
- A89 S06:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C23_lower_rear_hair/packets/s06_six/`;
- side-by-side audit sheet:
  `out/fumo_goal_review/visual_audit_sheet.png`
  (`862714c71321561490a66d7b75b1765918f5ef17744c8734e4a02af2450b33a1`).

## Absolute verdict

No recent candidate passes the absolute visual-quality gate. None reaches
8/10 in every applicable category, and each contains a major visible failure.
The protected rung is still the strongest *complete-character* visual
baseline. A88 S07 is the strongest *recent head-module* checkpoint. These are
different roles and must not be conflated.

The A85 + A87 + A88 assembly is mechanically composable and worth exactly one
bounded whole-head context test. It is not yet visually strong enough to be
treated as a production foundation or frozen indefinitely. A89 S06 is a clear
regression and must not be composed forward.

## Pixel-only scorecard

Scores are 0--10. Missing identity-defining parts lower the whole-subject
scores even when a candidate was intentionally incomplete. `Reject` is an
absolute decision, not a claim that no local improvement occurred.

| Candidate | Overall likeness | Silhouette / proportion | Construction | Identity | Contact / occlusion | Plush read | Presentation | Decision |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Protected rung 003 | 5.5 | 5.0 | 3.5 | 6.5 | 4.0 | 3.5 | 6.0 | Reject; best complete visual baseline |
| A85 S12 bare head | 1.5 | 6.0 | 4.5 | 0.5 | 5.0 | 3.0 | 6.0 | Reject as subject; retain only as provisional receiver |
| A87 S04 head + front frame | 3.0 | 5.5 | 3.5 | 4.0 | 3.5 | 2.5 | 6.0 | Reject; provisional local module |
| A88 S07 head + front + crown | 3.5 | 5.5 | 3.5 | 4.5 | 4.5 | 2.5 | 6.0 | Reject; best recent module checkpoint |
| A89 S06 lower-rear addition | 3.0 | 4.0 | 2.0 | 3.5 | 2.0 | 2.0 | 6.0 | Reject and undo |

## What the pixels show

### Protected rung 003

This remains the only reviewed candidate that reads immediately as a complete
Reimu-like character without explanation. The bow, face graphics, red/white
garment, and feet provide identity. It still misses the exact plush badly:
the body is too tall and narrow, the head and hair are a smooth hard shell,
the sleeves and skirt are rigid cones/cards, the feet are too exposed, and the
rear hair is a small set of geometric lobes rather than the reference's broad
layered drape. Its usefulness is as a whole-result comparator, not as an
approved asset or necessarily as production topology.

### A85 S12

The receiver is a better shallow cushion primitive than the earlier rounded
cube. Front width, corner softness, and depth are plausible enough for a
coupon. In beauty pixels it still reads as a uniformly smoothed rectangular
block, not a sewn, stuffed head: the front/rear planes and gusset are not
legible, the side is nearly vertical, and the lower shadow reads as a cavity.
Its measured silhouette success is real but much narrower than subject
likeness.

### A87 S04

The frontal central lock begins to establish Reimu's identity and is the first
recent change visible at ordinary viewing size. The three-quarter renders
expose the limitation: the fringe is a broad visor/card, the side locks are
thin blades attached tangentially, and the bald crown dominates. The pieces
have correct conceptual ownership but do not yet show padding, soft edge
rollover, root compression, or a sewn transition. This is local progress, not
a convincing hair frame.

### A88 S07

A88 closes the conspicuous bald top in front and therefore moves the front
image closer than A87. It also preserves the useful central fringe silhouette.
The gain is outweighed in profile and three-quarter by a large monolithic cap:
the crown is boxy, its lower boundary is a dead horizontal shelf, the side
surface is a broad slab, and the front panel reads as layered armor. The rear
view is mostly exposed receiver with only a narrow arch. Zero gap pixels prove
coverage, not plush construction or likeness.

### A89 S06

A89 does not move the whole result closer. The rear hair collapses into a
single rectangular cape with three pointed notches, while the side views show
edge-on fins and depth discontinuities. Large white receiver regions remain
visible at both rear sides. The roots and crown overlap do not read as the
reference's nested soft leaves. S06 removes one earlier front protrusion but
regresses silhouette, construction, contact, and plush read overall.

## Five largest current discrepancies

1. The head/hair assembly reads as a cuboid helmet rather than a shallow sewn
   cushion wrapped by padded fabric panels.
2. The rear hair is absent in A88 and becomes a rigid symmetric cape in A89;
   the reference has a dominant rounded central leaf plus overlapping,
   staggered side leaves with visible volume.
3. The fringe and cheek locks are too planar, sharp, and tangentially attached;
   they lack stuffed thickness, soft free edges, and root compression.
4. No current recent checkpoint includes face, bow, neck/body, or seated
   context, so whole-character proportions and recognition have not been
   tested since the reset.
5. The new head path and the protected complete-character path have no shared
   fixed-view A/B integration render; local progress therefore cannot be shown
   to improve the user's actual deliverable.

## Progress assessment across A83--A89

- A83 produced visible complete-character experiments and useful component
  evidence, but no whole candidate passed; its later hair and sleeve coupons
  remained visibly rigid.
- A84 correctly rejected rounded-cube parameter tuning, but its progression
  sheets show almost no perceptual change at whole-image scale.
- A85 made a genuine representation improvement in the bare head receiver.
- A87 made a genuine but modest identity improvement in the frontal hair.
- A88 made a genuine coverage improvement and a modest frontal improvement,
  while retaining a helmet-like construction in other views.
- A89 made no promotable visual progress.

The net recent result is therefore **a better experimental head scaffold, but
no demonstrated improvement to the complete plush**. The lack of apparent
progress is not primarily render speed. It comes from optimizing local binary
conditions (coverage, gap counts, protected pixels) on deliberately incomplete
modules without an early whole-head or whole-character recognition check.
Provisional forms were frozen after local success even while their absolute
plush read remained below 4/10, forcing later modules to route around visibly
weak shapes.

## Current best viable baselines

Use two explicitly labeled frozen comparators:

1. **Whole-character visual comparator:** protected rung 003. It is rejected
   but remains the most recognizable complete result.
2. **Recent construction experiment:** exact A88 S07. It is rejected as a
   subject but is the only recent assembly worth a bounded composability test.

Do not call A88 the new whole-model baseline, and do not promote any A89
geometry. If a single context test shows that A88's helmet/cuboid read is owned
by its frozen receiver or crown rather than missing context, reopen that owner
immediately instead of adding more modules around it.

## Recommended discriminating next cycle

1. Copy exact A88 S07 into a disposable context coupon.
2. Add only low-cost, proportionally registered proxies for flush eyes/mouth,
   the canonical bow silhouette, collar/body contact, and a rear-hair *gross
   volume*. These proxies are diagnostic and cannot be promoted.
3. Render fixed front, both three-quarters, both profiles, rear, and one full
   uncropped head/body-context view at the same camera and lighting settings.
4. Run an implementation-blind head-crop review against canonical front and
   turn frames 003, 010/011, 017--019, and 025/026. Inspect the whole silhouette
   before any owner-ID or gap metrics.
5. Continue from A88 only if the contextual render reads as the same Reimu
   Fumo variant and no major discrepancy is owned by the frozen A85/A87/A88
   geometry. Otherwise reopen the smallest failed owner; do not preserve weak
   pixels merely because a local gate once passed.
6. If A88 survives, rebuild lower-rear hair as visibly separate padded leaves:
   one broad convex central leaf and overlapping shorter side leaves, with
   unequal endpoints, non-edge-on profile thickness, and roots tucked under
   the crown. Reject from rear plus both profiles after the first probe if it
   again merges into a cape.
7. Integrate the passing head into a disposable copy of the complete body and
   compare directly with rung 003 before spending cycles on seams, fibers,
   materials, or rig refinements.

This cycle is intentionally a go/no-go test of the current architecture. Its
purpose is to restore whole-goal feedback quickly, not to one-shot the final
plush.
