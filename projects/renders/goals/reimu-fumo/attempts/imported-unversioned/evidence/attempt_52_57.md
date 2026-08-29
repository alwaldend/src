# Attempts 52–57 — envelope, receiver, and hybrid-maquette reset

[Back to attempt index](README.md) · [Back to goal](../README.md)

## Attempt 52 — canonical-turn visual hull

**Target.** Test whether the 30-frame canonical turn can provide a cheap
cross-view envelope without being mistaken for manufactured geometry.

**Plan and work.** Preserve the fixed `498 × 498` canvas, common axis
`x = 247.5 px`, nominal `12°` frame step, and source scale. Segment all
frames, carve loose and strict `96 × 96 × 112` voxel hulls, and reproject
them into the source views. The hull was explicitly limited to an ignored
outer-envelope guide.

**Evidence.** The
[source-versus-hull sheet](../../../../../out/reimu_fumo_attempt_052_visual_hull/evidence/source_vs_blender_hull.png)
and [report](../../../../../out/reimu_fumo_attempt_052_visual_hull/REPORT.md)
show a loose-hull mean recall of `.976`, worst recall `.909`, and worst
precision `.812`. The strict core is `79.3%` of the loose volume. The loose
hull fails the required per-view `.97` recall and `15%` excess-area gates.

**Decision.** Reject the hull as final geometry, metric reconstruction,
component separator, or seam authority. Retain only the reviewed masks and
the coarse uncertainty visualization. No working-ladder or tracked asset was
changed.

## Attempt 53 — direct-source construction contracts

**Target.** Replace failed envelope-driven geometry with explicit,
source-owned module contracts for the dominant head/hair receiver and seated
lower stack.

**Plan and work.** Audit rung 3 from pixels and from its saved datablocks;
calibrate four fixed turn views with `±12°` yaw uncertainty; trace the head,
hair, and lower-stack curves; and adversarially review whether constructed
panels are justified. The frozen head contract shares seam/root coordinates
between independent patches, not one conventional loft. The lower contract
requires broad connected support, a nonmonotonic rear return, explicit
front/side ownership, and a real manufactured hem relation.

**Evidence.** The
[head curve sheet](../../../../../out/reimu_fumo_attempt_053_head_seam_network_spec/curve_contact_sheet.png),
[lower interface diagram](../../../../../out/reimu_fumo_attempt_053_lower_stack_panel_spec/overlays/coupon_interface_diagram.png),
and
[construction board](../../../../../out/reimu_fumo_attempt_053_plush_panel_adversary/evidence/reference_construction_board.png)
freeze the direct observations and uncertainty. Rung 3 remains an absolute
visual reject at `5.5/10` likeness, `3.4/10` construction, and `3.1/10`
medium read. Its exact candidate hash is retained only as the source snapshot
for selective-copy experiments.

**Decision.** Freeze the contracts for bounded coupons. Treat the Attempt 52
hull and turn masks as outer-envelope vetoes only. No geometry candidate was
accepted in Attempt 53.

## Attempt 54 — receiver-normalized module probes

**Target.** Test three independent receiver hypotheses without integrating a
complete model.

**Work and results.** Face-applique configuration B is a strong relative
module direction with a maximum nominal projection of `1.35168 mm`, but its
blind scores remain `6.0–7.5/10`; it is evidence for later receiver testing,
not an approved face. Both hull-initialized head receivers fail: configuration
A reproduces fused slab/helmet artifacts, while B becomes a generic
capsule/egg and misses the direct-mask gates. Lower-receiver B passes its
bounded hidden-support checks at `.98545` minimum allowed precision, with a
`.758 × .720 Wh` floor span and three rear-profile direction changes, but it
is only a provisional hidden fitting control.

**Evidence.** Review the
[face comparison](../../../../../out/reimu_fumo_attempt_054_face_applique/blind_review_sheet.png),
[head-receiver comparison](../../../../../out/reimu_fumo_attempt_054_hull_head_receiver/contact_sheet.png),
and
[hidden lower-receiver sheet](../../../../../out/reimu_fumo_attempt_054_hull_lower_receiver/fixed_view_sheet.png).

**Decision.** Reject both head receivers and any visible use of the lower
receiver. Preserve the face hierarchy and lower support only as bounded input
to later tests. Nothing was promoted.

## Attempt 55 — explicit constructed receivers

**Target.** Determine whether source-keyed manufactured patches can replace
the helmet head and ramp lower stack.

**Head result.** Two coherent cycles fail the front/side cheap gate. Cycle 2
still shows an angular helmet, generic jagged `W`, large beige bald crown,
hard cap/face boundary, egg/mattress cushion, floating roots, and parallel
rear sheets. Profile `1.250 Wh`, added rear reach `.415 Wh`, and rear width
`.939 Wh` also exceed their bands. The five-view approval packet was
correctly skipped.

**Lower result.** The hidden `C0/W0` support passes broad-contact checks:
`.520 × .460 Wh`, area `.1873 Wh²`, bounding-box fill `.7832`, and all four
quadrants. The visible panels fail with a doubled front tier, planar
ramp/tent side, and excessive pattern strain. Blind scores are `2/5` fabric,
`2/5` mass, and `1/5` side pooling.

**Evidence.** See the
[head Cycle 2 gate](../../../../../out/reimu_fumo_attempt_055_constructed_head_receiver/cycle_2/cheap_gate/contact_sheet.png)
and
[lower front/side/x-ray sheet](../../../../../out/reimu_fumo_attempt_055_constructed_lower_receiver/packet/front_side_xray_review.png).

**Decision.** Reject every visible A55 mesh. Retain direct curves, fixed
gates, and a copied `C0/W0` hidden-support control only. No integration or
promotion is authorized.

## Attempt 56 — representation-changing head and skirt probes

**56A dense sculpt.** Two directly sculpted states replace section lofts but
remain a rounded foam block or upright mattress. State 2 reaches the nominal
`.7503 Wh` depth floor and improves the front plane to `6/10`, yet scores only
`2/10` for construction, gusset behavior, and bottom closure. Reject both
states and the continuously smoothed base-head representation.

**56B sewn pressure solve.** Two explicit front/rear/gusset configurations
are finite and watertight. Radial configuration A becomes a button-tufted
pillow with a rear crater. Pole-free configuration B is stable at a
`.065 mm` final-eight-frame RMS step but becomes a uniform foam mattress.
Reject both; uniform pressure and one continuous gusset do not supply the
required differential construction.

**56C sculpted skirt shell.** State B removes A55's planar gores and passes
outer-envelope containment, exact shared seams, and distributed support. It
still reads as a padded bench from front/rear and a prow from side, while
pattern strain reaches `72.9%` maximum and `65.1%` p95. Reject both states and
the closed taut shell with pulled support pockets.

**Evidence.** The decisive artifacts are the
[56A State 2 sheet](../../../../../out/reimu_fumo_attempt_056a_sculpted_head_cushion/state_2/fixed_views/contact_sheet.png),
[56B sewn-solve comparison](../../../../../out/reimu_fumo_attempt_056b_sewn_head_cushion/contact_sheet.png),
and
[56C state comparison](../../../../../out/reimu_fumo_attempt_056c_sculpted_skirt_shell/packet/state_a_b_review.png).

**Decision.** Attempt 56 accepts no candidate geometry. Preserve the replay
mechanisms, measurements, and hidden `C0/W0` control; change construction
ownership before another head or skirt solve.

## Attempt 57 — hybrid head maquette (rejected at State 1)

**Failure targeted.** Prior heads become either a beige mattress beneath a
hair shell or a brown helmet/card assembly. The unresolved question is whether
the observable brown crown/core and nearly flush face can first be shaped as
a credible shallow target while a separate leaf owns posterior reach.

**Frozen hypothesis.** Build one disposable visual maquette from factory
startup: a shallow brown core at `.74 ± .05 Wh` depth, a nearly flush beige
face, the canonical asymmetric three-span fringe, and one clearly separate
posterior leaf witness. The complete core-plus-leaf depth is approximately
`1.02 ± .06 Wh`. The maquette is never a reusable topology parent.

**Frozen gate.** Produce only fixed 512 px front and true-side beauty, semantic
ID, and registered overlays for State 1. Visible beige remains
`.603 ± .03 Wh` wide and high; beige projection is at most `.015 Wh`; the
central tip remains within `.03 Wh` of `(.588,.677)`. Any hood, helmet,
mattress, plaque, embossed fringe, floating layer, or theft of the rear-leaf
extremum rejects the state. Each visual category must reach `8/10` with no
major defect. At most one bounded correction is allowed.

**Result.** The [strategy decision](../../../../../out/reimu_fumo_attempt_057_strategy_reset/DECISION.md)
and [candidate contract](../../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/ATTEMPT.md)
produced one preserved State 1 at SHA-256
`3acadcc30987872d214ec9bf73cd57dfe28e34935674085ec4458af815a247fd`.
Reopen validation and every frozen scalar band pass, but the
[review packet](../../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/review_packet.png)
shows a circular helmet, angular raised face plaque, hard M/W fringe, and
oversized rigid rear slab. The
[implementation-blind review](../../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/BLIND_REVIEW.md)
independently rejects it and finds that one bounded correction cannot repair
the multiple macro failures. No State 2 was built and no geometry was
promoted.

## Progress, approach, and process audit

- Attempts 52–54 improved source authority, fixed-view calibration, and safe
  module boundaries, but did not advance a complete candidate.
- Attempts 55–56 changed representations and still repeated helmet/mattress,
  floating-sheet, ramp/bench, and uniform-pressure failures. Numeric
  containment and topology repeatedly failed to predict manufactured pixels.
- The highest-leverage unresolved problem remains the coupled visible
  head/hair macro-form. Attempt 57 tests it before hidden topology, seams,
  rigging, or full packets.
- Attempt 57 kept one factory-startup candidate, two decisive views, immutable
  tracked files, and a bounded-correction rule. Its first pixels exposed a
  representation error early enough to avoid another full-model rebuild.
- Attempt 58 now compares three isolated head-layer hierarchies in parallel,
  adding a three-quarter gate while preserving the same small-module and
  immutable-main-file discipline.
