# Reimu Fumo artifact log

[Canonical goal projection](../../goals/reimu-fumo/README.md)

This is a noncanonical, human-facing local-session review log. Entries are
newest-first and link to actual local artifacts. The canonical goal and
attempt records own lifecycle, acceptance, resource-version, and retained-
evidence state; this log may include newer uncheckpointed observations and
must not be used as acceptance evidence. At this 2026-09-01 log revision, the
canonical cutoff is goal resource version 51 and A88 attempt resource version
1. A85 is closed as refine, A86 as reset, and A87 as refine; A88 is the active
two-panel crown-interface attempt. Files under `out/` are
disposable working
evidence and are intentionally not committed; final approved deliverables
will live under `projects/renders/blender/fumo/` and be tracked with Git LFS.

## 2026-09-01 — current turn

- [A88 S07 six-view front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/front.png),
  [profile](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/profile_left.png),
  [rear](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/rear.png),
  [component IDs](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/front_ids.png),
  [manifest](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/manifest.json),
  [source Blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/snapshots/A88_crown_s07.blend),
  [blind review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/A88_S07_BLIND_REVIEW.md),
  and [six-view measurement gate](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/A88_S07_MEASURE_GATE.md)
  — the bounded two-panel crown coupon now has continuous front contact, no
  exposed receiver at the apex, and a defined rear overlap after S00/S01
  intersection failures and S02–S06 contact/clearance diagnostics.
  **Verdict:** preserve S07 as the provisional crown module; lower rear hair is
  the next separate owner, and this is not standalone sculpt approval.

- [A87 exact S04 survivor front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/packets/s04/front.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/packets/s04/three_quarter.png),
  [mirror](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/packets/s04/three_quarter_mirror.png),
  [source Blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/snapshots/A87_hair_frame_s04.blend),
  [reference measurements](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/A87_REFERENCE_MEASUREMENTS.md),
  and [close review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C21_front_hair_frame/A87_CLOSE_REVIEW.md)
  — the front hair is now three thin constructed panels with a short asymmetric
  fringe and two tapered cheek locks. **Verdict:** A87 closed as refine with
  S04 preserved provisionally; the later S07 overfit was rejected because its
  locks became long anime strands.

- [A86 collar-only blind review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C20_collar_body_interface/A86_BLIND_REVIEW.md)
  and [S02 diagnostic Blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C20_collar_body_interface/snapshots/A86_coupon_s02.blend)
  — a compact collar cannot hide the wide receiver underside without becoming
  a hard head-width yoke. **Verdict:** A86 closed as reset; lateral and rear
  coverage belongs to hair and garment panels, not one collar band.

- [A85 S12 pinned front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/packets/a85_head_cage_s12_pair/front.png),
  [profile](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/packets/a85_head_cage_s12_pair/profile_left.png),
  [manifest](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/packets/a85_head_cage_s12_pair/manifest.json),
  [source Blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/snapshots/A85_head_cage_s12.blend),
  [measurement gate](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a85_review/S12_MEASURE_GATE.md),
  [blind review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a85_review/S12_BLIND_REVIEW.md),
  and [full cycle record](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a85_review/A85_CYCLE.md)
  — the new directly authored quad cage now passes every measurable bare-head
  foundation silhouette gate, including the complete cheek band and corrected
  `.094 Wh` crown rise. **Verdict:** freeze S12 as the provisional receiver for
  collar/hair integration, not as standalone sculpt approval; the shaded
  underside remains an explicit interface risk that the collar must cover
  cleanly without clipping or gaps.

- [A85 constructed-cage front/profile packet](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/packets/a85_head_cage_s01_pair/manifest.json)
  and [front comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a85_review/A85_front_comparison.png)
  — the effect-verified six-row cage replaces the rounded-cube scaffold with
  broad cheeks, an explicit chin taper, restrained crown ownership, and the
  protected shallow depth. It is the first A85 controlling pair and remains
  under measurement and implementation-blind review before sculpt polish.

- [A84 front progression](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a84_review/A84_front_progression.png),
  [profile progression](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a84_review/A84_profile_progression.png),
  and [cycle verdict](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a84_review/A84_CYCLE.md)
  — Elastic Grab produced one useful lower-taper delta, but the rounded-cube
  topology could not form a controllable crown: one packet was exactly
  zero-effect and the corrected stroke was negligible. **Verdict:** close A84
  as reset and change the coarse topology rather than tune another brush.

## 2026-08-31 — current turn

- [A83 C18 hair four-view owner-ID comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C18_hair_layered/owner_id_renders/C18_owner_id_four_view_contact_sheet.png),
  [pinned pixel inspection](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C18_hair_layered/owner_id_renders/PIXEL_INSPECTION.md),
  and [representation review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C18_hair_verdict/NEXT.md)
  — the hidden-liner, close-mantle, and three-padded-leaf interface fails its
  first packet with a large beige side gap, helmet mantle, high inflated rear
  pods, and disconnected roots. **Verdict:** undo C18 hair and close this owner
  map without a brown-material or rescue pass.

- [A83 C18 sleeve comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C18_pouch_sleeve/contact_sheet.png)
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C18_pouch_sleeve/REVIEW.md)
  — a non-planar pouch fixes C17's flatness but becomes a small inflated
  mitten/capsule rather than a broad sewn bell sleeve. **Verdict:** undo C18
  sleeve; change the next authoring mode to direct contour sculpting instead
  of another whole generated sleeve representation.

- [A83 C17 final comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/final_contact_sheet.png),
  [owner-ID comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/owner_id_contact_sheet.png),
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C17_panel_sleeve/REVIEW.md)
  — two independently owned closed cloth pockets pass the ownership and
  collision checks, but become a stiff tapered paddle with weak shoulder
  attachment, insufficient stuffed depth, and a flat cuff in the controlling
  material views. **Verdict:** undo C17; preserve only the evidence that a
  panel split alone does not create a credible stuffed bell sleeve.

- [A83 C16 owner-ID front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v0_owner_id/front.png),
  [owner-ID three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v0_owner_id/three_quarter.png),
  [rejected direct-pass three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v1_direct/three_quarter.png),
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/REVIEW.md)
  — three separately owned hair surfaces proved complete initial coverage, but
  their first direct sculpt pass exposed a large beige side strip and turned
  the return leaf into a hard card. **Verdict:** undo C16 and retain only the
  v0 owner map as interface evidence.

- [A83 C15 pinned front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/pinned_5_2_1/front.png),
  [pinned three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/pinned_5_2_1/three_quarter.png),
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/REVIEW.md)
  — a live whole-field cap edit changed 8,434 vertices but exposed a large
  beige receiver region before it created layered construction; pinned Blender
  5.2.1 reproduced the failure. **Verdict:** undo C15 and close cap-only
  whole-field deformation.

- [A83 C14 sleeve comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C14_sleeve_module/contact_sheet.png)
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C14_sleeve_module/REVIEW.md)
  — the collision-free replacement removed the hollow tube but became a
  narrow curved limb with a pale nub rather than a broad sewn bell sleeve.
  **Verdict:** undo C14 and retire the axial-loft sleeve representation.

- [A83 Stage R2 pose sheet](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/rig_stage2_pose_contact_sheet.png),
  [sleeve close-up](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/rig_stage2_sleeve_closeup.png),
  [isolated blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/reimu_fumo_rig_stage2.blend),
  and [report](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/REPORT.md)
  — one owner-local lattice adds verified non-rigid motion to all 22 right-
  sleeve components while preserving the four Stage R1 actions, a fixed root,
  and trim registration. **Verdict:** retain as bounded soft-deformation
  scaffolding; full criterion 006 remains incomplete.

- [A83 C13 comparison](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C13_coupled_panel_support/contact_sheet.png)
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C13_coupled_panel_support/REVIEW.md)
  — one localized underlap support retained C12's visibility-proven panel and
  prevented receiver leakage, but added only a faint crease while the same
  spherical cap still owned the silhouette. **Verdict:** undo C13 after the
  first fixed-view packet; another local support refinement cannot resolve
  `D-HEAD-HELMET`.

- [A83 Stage R1 pose sheet](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/rig_stage1_pose_contact_sheet.png),
  [isolated rig blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/reimu_fumo_rig_stage1.blend),
  and [verification report](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/REPORT.md)
  — the isolated C1b copy now has separated asset/review collections, nine
  controls, nine deform bones, stable ownership for 116 visible objects, and
  the four exact required actions. **Verdict:** retain as coarse rig
  scaffolding and candidate-bound criterion-005 evidence; criterion 006 still
  fails because deformation is rigid and append/every-frame gates are absent.

- [A83 C12 owner-ID view](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/owner_id_pass2/three_quarter.png),
  [brown three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/brown_pass2/three_quarter.png),
  [brown rear](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/brown_pass2/rear.png),
  and [report](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/REPORT.md)
  — a cap-preserving constructed panel owns a connected `64 x 82 px` region
  with no beige leak and keeps one readable free edge in brown. **Verdict:**
  pass only the panel representation/ownership question; the coarse geometry
  is non-candidate and the whole hair mass still fails the helmet gate.

- [A83 C10 topology-stable front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C10_shapekey_packet/front.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C10_shapekey_packet/three_quarter.png),
  and [blind verdict](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C10_KEEP_UNDO.md)
  — reversible paired shape keys halved two upper-crown width errors while
  preserving topology, vertex order, the face opening, bow seat, lower rows,
  and positive cap/head separation. **Verdict:** undo C10; the fixed pixels
  show no clear fidelity gain and retain the spherical helmet construction.

- [A83 animation acceptance packet](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/animation_acceptance/ANIMATION_ACCEPTANCE.md),
  [current audit](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/animation_acceptance/CURRENT_AUDIT.md),
  and [production deformation design](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/rig_deformation_design/DEFORMATION_DESIGN.md)
  — C1b is confirmed static and the disposable rig remains only a rigid-
  cluster feasibility proof. **Verdict:** criteria 005 and 006 remain failed;
  the final asset must pass exact-action, non-rigid deformation, all-frame
  contact, clean-append, and implementation-blind animation gates.

- [A83 disposable rig pose sheet](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/NON_CANDIDATE_rig_pose_contact_sheet.png),
  [rig-readiness report](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/REPORT.md),
  [production rig contract](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/RIG_CONTRACT.md),
  and [probe blend](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/NON_CANDIDATE_C1b_rig_probe.blend)
  — an isolated non-candidate proves the modular head/hair/bow and sleeve
  clusters can execute a bounded head turn and arm wave without visible
  tearing or clipping. **Verdict:** C1b itself remains static and fails the
  animation criteria; the probe is a coarse rigid-cluster feasibility test,
  not a production rig or a promotable model.

- [A83 C11 constructed-panel three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/renders/three_quarter.png),
  [rear](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/renders/rear.png),
  and [review](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/REVIEW.md)
  — a disposable panel coupon did not create the required broad lapped plane;
  it preserved the helmet read and exposed beige at the root. **Verdict:**
  undo C11 without polish or canonical follow-up.

- [A83 C10 native-sculpt result](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/native_sculpt_c10_crown/RESULT.md)
  — the paired fixed-front crown test passed targeting and freeze checks, but
  the first receiver Grab produced zero coordinate change. **Verdict:** stop
  before later strokes, renders, or save; no candidate exists and native Grab
  is retired for this crown location.

- [A83 C9 native-sculpt front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C9_packet/front.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C9_packet/three_quarter.png),
  [stroke metrics](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/native_sculpt_c9/receiver_metrics.json),
  and [blind verdict](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C9_KEEP_UNDO.md)
  — Blender's native Grab operator moved only the paired cap and receiver
  neighborhoods, but the roughly one-millimetre displacement did not create a
  legible side plane, lap, or fabric-panel break at the fixed 512-pixel views.
  **Verdict:** undo C9; native operator activity is not visual progress, and
  the next bounded test changes one side-hair owner rather than smoothing the
  same continuous shell again.

- [A83 C8 front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C8_packet/front.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C8_packet/three_quarter.png),
  and [blind verdict](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C8_KEEP_UNDO.md)
  — two localized tangent-plane shape keys moved only the intended cap and
  receiver vertices, but remained visually indistinguishable from C1b in the
  controlling views. **Verdict:** undo C8 and retire another continuous-shell
  coordinate deformation as the owner of the helmet-like side silhouette.

- [A83 retained C1b front](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C1b_packet/front.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C1b_packet/three_quarter.png),
  [C5 rear probe](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C5_packet/rear.png),
  [C5 blind verdict](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C5_KEEP_UNDO.md),
  and [reference-construction verdict](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C5_design/REFERENCE_VERDICT.md)
  — one coupled cap-and-receiver narrowing is retained as a small internal
  improvement. Four subsequent local edits were undone: three were visually
  ineffective, and the center seam was both byte-invisible in three-quarter
  and unsupported by all canonical turn frames. **Verdict:** keep only C1b;
  the next test changes one real rear-lock overlap, not the whole head.

- [A83 C6 three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C6_packet/three_quarter.png),
  [C6 blind undo](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C6_KEEP_UNDO.md),
  [C7 three-quarter](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/C7_packet/three_quarter.png),
  and [C7 blind undo](../../../../out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C7_KEEP_UNDO.md)
  — two reversible edits proved that the separate camera-side rear lock can
  change pixels but only becomes a narrow tail/tendril while the smooth cap
  remains dominant. **Verdict:** undo both and retire lock placement as the
  current owner; next test a bounded direct cap-side sculpt or panel from C1b.

- [A81 exact-front render](../../../../out/reimu_fumo_attempt_081_live_joint_restform/batch_render/packet/front.png),
  [three-quarter render](../../../../out/reimu_fumo_attempt_081_live_joint_restform/batch_render/packet/three_quarter.png),
  [implementation-independent verdict](../../../../out/reimu_fumo_attempt_081_live_joint_restform/standing_review/CATEGORICAL_VERDICT.md),
  and [timing record](../../../../out/reimu_fumo_attempt_081_live_joint_restform/timing/TIMELINE.md)
  — the direct live-edit experiment produced an enormous exposed beige support
  and only a narrow disconnected brown rear band. The pinned fallback rendered
  the frozen pair in 17.906 seconds, but the loop missed its pair deadline by
  6m31.918s. **Verdict:** reset without repair or a third view; close the
  receiver/cap/yoke family and start a fresh atomic head-and-hair shell.

- [A80 constructed-cap sheet](../../../../out/reimu_fumo_attempt_080_diagnostic_hair_blockout/constructed_cap/a80_constructed_cap_REJECTED_sheet.png),
  [sculpted-cap board](../../../../out/reimu_fumo_attempt_080_diagnostic_hair_blockout/sculpted_cap/three_view_NON_CANDIDATE_board.png),
  [blind review](../../../../out/reimu_fumo_attempt_080_diagnostic_hair_blockout/blind_review/BLIND_REVIEW.md),
  and [speed postmortem](../../../../out/reimu_fumo_attempt_080_diagnostic_hair_blockout/process_review/SPEED_POSTMORTEM.md)
  — both complete coarse hair assemblies reached pixels cheaply once Blender
  started, but one became a bald helmet band and rear balloon while the other
  exposed a bald cuboid support. **Verdict:** reset both diagnostic variants;
  neither is preferred to rung 003.

- [A79 full five-view diagnostic](../../../../out/reimu_fumo_attempt_079_paired_hair_field/diagnostic_preview/a79_non_candidate_five_view.png),
  [pixel diagnosis](../../../../out/reimu_fumo_attempt_079_paired_hair_field/diagnostic_preview/DIAGNOSIS.md),
  and [render manifest](../../../../out/reimu_fumo_attempt_079_paired_hair_field/diagnostic_preview/manifest.json)
  — the first whole-character A79 pixels expose a huge bald crown, rigid rear
  curtain, detached board-like lock, and broken joins. The sheet is explicitly
  non-candidate; protected hashes remained exact and clean reopen passed.
  **Verdict:** reset the paired-field representation before promotion work.

- [A77--A79 process timeline](../../../../out/reimu_fumo_attempt_079_paired_hair_field/process_optimization/PROCESS_TIMELINE.md),
  [diagnostic-first loop](../../../../out/reimu_fumo_attempt_079_paired_hair_field/process_optimization/LOOP_DESIGN.md),
  and [primary process verdict](../../../../out/reimu_fumo_attempt_079_paired_hair_field/process_optimization/PROCESS_OPTIMIZATION.md)
  — A79 took 94m30s to reach pixels that rendered in 29.115s. The replacement
  loop targets decision-bearing pixels within 12 minutes, stops categorical
  failures immediately, and reserves promotion-grade validation for a visual
  survivor. The generic [goal skill](../../../goal/skills/goal/SKILL.md) now
  preserves this staged-feedback rule. **Verdict:** process correction applied
  and validated; A80 uses it immediately.

## 2026-08-30 — current turn

- [Attempt 67L full five-view sheet](../../../../out/reimu_fumo_attempt_067_lower_topology/five_view_sheet.png),
  [manifest](../../../../out/reimu_fumo_attempt_067_lower_topology/manifest.json),
  and [construction contract](../../../../out/reimu_fumo_attempt_067_lower_topology/design/CONSTRUCTION_CONTRACT.md)
  — all mechanical, contact, grain-path, seam, width, and frozen-context gates
  pass, yet exact side is still a tall cape/ramp/tent, front and rear are rigid
  trapezoids, and the low pool is visually swallowed. **Verdict:** terminal
  reject; close the generated T-yoke/apron/drop/pool representation.

- [Attempt 67H full five-view sheet](../../../../out/reimu_fumo_attempt_067_head_topology/five_view_sheet.png),
  [reference/candidate comparison](../../../../out/reimu_fumo_attempt_067_head_topology/reference_candidate_sheet.png),
  [first blind review](../../../../out/reimu_fumo_attempt_067_head_topology/BLIND_REVIEW_1.md),
  [second blind review](../../../../out/reimu_fumo_attempt_067_head_topology/BLIND_REVIEW_2.md),
  and [terminal report](../../../../out/reimu_fumo_attempt_067_head_topology/TERMINAL_REPORT.md)
  — seven of eight front-width rows fail, the crown exposes `2,751` beige
  pixels, rear H/W is `1.092`, and the form repeats a bald plaque, deep faceted
  helmet, rear block/cape, and floating cards. **Verdict:** terminal reject;
  close the analytic cushion/yoke/pocket representation after its second
  reviewed topology cycle.

- [Attempt 67H construction design](../../../../out/reimu_fumo_attempt_067_head_topology/design/CONSTRUCTION_DESIGN.md)
  and [Attempt 67L construction contract](../../../../out/reimu_fumo_attempt_067_lower_topology/design/CONSTRUCTION_CONTRACT.md)
  — both active lanes rejected a cheap near-repeat before Blender and froze a
  materially different ownership graph. H uses an open high-terminating
  fringe/crown yoke with asymmetric root-mapped rear pockets; L uses a T-yoke,
  front apron, rear/side drop, and separate low horseshoe pool. **Verdict:**
  accepted only as the one-packet Attempt 67 hypotheses; rendered pixels still
  control acceptance.

- [Attempt 66H full five-view sheet](../../../../out/reimu_fumo_attempt_066_head_receiver/five_view_sheet.png),
  [reference/candidate comparison](../../../../out/reimu_fumo_attempt_066_head_receiver/reference_candidate_sheet.png),
  [implementation-blind review](../../../../out/reimu_fumo_attempt_066_head_receiver/BLIND_REVIEW_1.md),
  and [build manifest](../../../../out/reimu_fumo_attempt_066_head_receiver/candidate_v1/build_manifest.json)
  — the front allocation moves toward the source, but the complete receiver is
  still a deep rounded box with a beige face plaque, bald crown patch, global
  mantle, regular rear petals, and card-like detached strands. **Verdict:**
  terminal reject at the first packet; the second topology must split the
  front/top cap from rear-leaf ownership rather than smooth this helmet.

- [Attempt 66L full five-view sheet](../../../../out/reimu_fumo_attempt_066_lower_stack/five_view_sheet.png),
  [annotated veto overlay](../../../../out/reimu_fumo_attempt_066_lower_stack/reference_contract/overlays/a66_exact_candidate_veto_overlay.png),
  [implementation-blind review](../../../../out/reimu_fumo_attempt_066_lower_stack/BLIND_REVIEW.md),
  [self-review](../../../../out/reimu_fumo_attempt_066_lower_stack/SELF_REVIEW.md),
  and [manifest](../../../../out/reimu_fumo_attempt_066_lower_stack/manifest.json)
  — every scalar/mechanical gate passes, including zero new foot collisions,
  but the fixed views expose a disconnected table/tent in front, a straight
  side ramp, cape/fins in obliques, long floor rails, and no compressed seated
  mass. **Verdict:** terminal reject at the first packet; change topology once
  rather than tune this ring-and-lap family.

- [Attempt 65B front render](../../../../out/reimu_fumo_attempt_065_parallel/b_return_leaf/candidate_v1_packet/front.png),
  [exact side](../../../../out/reimu_fumo_attempt_065_parallel/b_return_leaf/candidate_v1_packet/side.png),
  [rear](../../../../out/reimu_fumo_attempt_065_parallel/b_return_leaf/candidate_v1_packet/rear.png),
  and [manifest](../../../../out/reimu_fumo_attempt_065_parallel/b_return_leaf/candidate_v1_manifest.json)
  — the explicit outgoing/turn/return construction is mechanically clean and
  reaches `1.498092 Wh`, but renders as tall rectangular ears, upright oblique
  slabs, one side fin with a rolled rim, and separated rear roots. **Verdict:**
  rejected at the first packet; Attempt 65 and visible bow work are closed.

- [Attempt 66H reference-contract board](../../../../out/reimu_fumo_attempt_066_head_receiver/reference_contract/contract_review_board.png),
  [independent contract](../../../../out/reimu_fumo_attempt_066_head_receiver/reference_contract/REFERENCE_CONTRACT.md),
  and [read-only parent inventory](../../../../out/reimu_fumo_attempt_066_head_receiver/inventory/INVENTORY.md)
  — all canonical-turn frames and supporting physical references now control a
  precise fifteen-object head/hair/crown receiver boundary. The board exposes
  the parent's oversized crown, box/helmet profile, weak asymmetric fringe,
  short symmetric rear petals, and unusable bow seating. **Verdict:** accepted
  as the live Attempt 66H build-and-review contract; it is evidence, not a
  candidate pass.

- [Attempt 66L read-only lower-stack inventory](../../../../out/reimu_fumo_attempt_066_lower_stack/inventory/INVENTORY.md)
  and [machine-readable audit](../../../../out/reimu_fumo_attempt_066_lower_stack/inventory/INVENTORY.json)
  — the current feet cross the front ruffle by `239/259` triangle pairs, the
  internal seat floats `11.2 mm` above the floor, and contact is only two foot
  patches plus three tiny hem points. **Verdict:** accepted as the isolated
  replacement boundary and regression contract; first candidate pixels are
  pending.

- [Attempt 65 strategy review](../../../../out/reimu_fumo_attempt_065_parallel/STRATEGY_REVIEW.md)
  — the adversarial whole-process comparison finds the same side/root bow
  failure across materially different representations and restores the crown-
  receiver dependency. **Verdict:** close visible bow work after the final
  already-running first packet; make head/hair/crown primary and the compact
  seated lower stack parallel.

- [Attempt 65C full five-view sheet](../../../../out/reimu_fumo_attempt_065_parallel/c_native_surface_relief/five_view/contact_sheet.png),
  [self-review](../../../../out/reimu_fumo_attempt_065_parallel/c_native_surface_relief/SELF_REVIEW.md),
  and [manifest](../../../../out/reimu_fumo_attempt_065_parallel/c_native_surface_relief/manifest.json)
  — native relief reaches the exact `1.497000004 Wh` upper span and improves
  frontal folds, but exact side remains an upright fin, both obliques expose a
  paddle, and the rear root becomes a rounded knot. **Verdict:** rejected at
  its mandated first packet; no tuning or promotion.

- [Attempt 65A front render](../../../../out/reimu_fumo_attempt_065_parallel/a_lowres_sculpt_cage/candidate_v1_packet/front.png),
  [side render](../../../../out/reimu_fumo_attempt_065_parallel/a_lowres_sculpt_cage/candidate_v1_packet/side.png),
  [rear render](../../../../out/reimu_fumo_attempt_065_parallel/a_lowres_sculpt_cage/candidate_v1_packet/rear.png),
  and [attempt report](../../../../out/reimu_fumo_attempt_065_parallel/a_lowres_sculpt_cage/ATTEMPT.md)
  — the parallel low-resolution sculptable-cage branch passes topology and
  isolation but becomes mouse-ear pads, oblique paddles, a vertical side
  petal, and hollow rear shells. **Verdict:** rejected at its mandated first
  packet; no tuning or promotion.

- [Attempt 64 full five-view sheet](../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/five_view_sheet.png),
  [implementation-blind review](../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/BLIND_REVIEW.md),
  [second blind review](../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/SECOND_REVIEW.md),
  and [candidate manifest](../../../../out/reimu_fumo_attempt_064_sewn_fan_lobes/candidate_v2_manifest.json)
  — two independent sewn fan patterns have shared-vertex turned trim and the
  corrected `1.497 ± .014 Wh` upper span, but exact side exposes an open V/cup
  and both obliques remain paddles. **Verdict:** rejected at the neutral early
  gate before saddle, applique, or material polish.

- [Attempt 63 full five-view sheet](../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/five_view_sheet.png),
  [first blind review](../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/BLIND_REVIEW.md),
  [second blind review](../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/SECOND_REVIEW.md),
  and [complete parent-bow inventory](../../../../out/reimu_fumo_attempt_063_constructed_upper_bow/PARENT_BOW_INVENTORY.md)
  — a continuous upper pocket fixes A62's detached surface graphics but becomes
  one rectangular pillow in front/rear and one thick upright block in side.
  **Verdict:** reject after one packet and retire the hourglass representation.

- [Attempt 62 parent/candidate five-view matrix](../../../../out/reimu_fumo_attempt_062_bow_silhouette/parent_candidate_matrix.png),
  [first blind review](../../../../out/reimu_fumo_attempt_062_bow_silhouette/BLIND_REVIEW.md),
  and [second blind review](../../../../out/reimu_fumo_attempt_062_bow_silhouette/SECOND_REVIEW.md)
  — narrower span is a useful signal, but the incomplete nine-MESH transform
  leaves six world-space curves behind and exposes card, trim, and cuboid-root
  failures. **Verdict:** reject and close the affine family.

- [Attempt 61 parent/candidate five-view matrix](../../../../out/reimu_fumo_attempt_061_head_allocation/parent_candidate_matrix.png),
  [calibrated blind review](../../../../out/reimu_fumo_attempt_061_head_allocation/BLIND_REVIEW.md),
  and [independent relative review](../../../../out/reimu_fumo_attempt_061_head_allocation/SECOND_REVIEW.md)
  — reviewers disagree on a modest relative head-allocation gain, but both
  reject absolutely and the calibrated side contact grows from about `6 px`
  to `42 px`. **Verdict:** reject; no promotion from a measured contact veto.

- [Attempt 60 parent/candidate five-view matrix](../../../../out/reimu_fumo_attempt_060_direct_module/parent_candidate_matrix.png),
  [first blind review](../../../../out/reimu_fumo_attempt_060_direct_module/BLIND_REVIEW.md),
  and [second blind review](../../../../out/reimu_fumo_attempt_060_direct_module/SECOND_REVIEW.md)
  — the head-share direction improves, but complete vertical compression
  detaches/crosses sleeves, buries feet, and hardens the skirt/floor contact.
  **Verdict:** reject and retire the cross-junction affine interface.

- [Attempt 59 runtime decision](../../../../out/reimu_fumo_attempt_059_method_reset/SOLVER_RUNTIME_DECISION.md),
  [solver divergence audit](../../../../out/reimu_fumo_attempt_059_method_reset/SOLVER_DIVERGENCE_AUDIT.md),
  [scale decision](../../../../out/reimu_fumo_attempt_059_method_reset/SOLVER_SCALE_DECISION.md),
  [mantle/lock report](../../../../out/reimu_fumo_attempt_059_method_reset/mantle_lock/REPORT.md),
  and [all-reference landmark board](../../../../out/reimu_fumo_attempt_059_method_reset/reference_inputs/reference_board_owned_landmarks.png)
  — the exact structural solve diverges from its material gate and the fabric
  overlay independently fails release/contact/collision gates. **Verdict:**
  reject before candidate pixels, preserve the reference packet and negative
  method evidence, and stop the numerical branch.

- [Attempt 57 review packet](../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/review_packet.png),
  [implementation-blind review](../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/BLIND_REVIEW.md),
  and [report](../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/REPORT.md)
  — the candidate passes all frozen scalar bands and reopen validation but
  reads as a circular helmet, raised face plaque, hard M/W cut, and rigid rear
  slab. **Verdict:** reject the only State 1, make no correction or State 2,
  and retire the brown-core/beige-plaque representation.

- [Attempt 58 anonymous beauty matrix](../../../../out/reimu_fumo_attempt_058_layered_head/blind_comparison/beauty_matrix.png),
  [form matrix](../../../../out/reimu_fumo_attempt_058_layered_head/blind_comparison/form_matrix.png),
  and [blind review](../../../../out/reimu_fumo_attempt_058_layered_head/blind_comparison/BLIND_REVIEW.md)
  — the three isolated copies become an egg/helmet with remote leaves, a
  rectangular mattress with rigid cards, and an exposed-scalp armor/cape
  assembly. Independent reviews disagree on their relative order. **Verdict:**
  reject and retire all three; no least-bad mesh is selected or promoted.

- [Attempt 59 controlling contract](../../../../out/reimu_fumo_attempt_059_method_reset/ATTEMPT.md),
  [method audit](../../../../out/reimu_fumo_attempt_059_method_reset/METHOD_AUDIT.md),
  [source/pattern specification](../../../../out/reimu_fumo_attempt_059_method_reset/SOURCE_PATTERN_SPEC.md),
  and [immutable context specification](../../../../out/reimu_fumo_attempt_059_method_reset/CONTEXT_SPEC.md)
  — the proof changed to a deterministic local inverse-pattern structural
  solve plus separately constructed unpressurized mantle and one lock. The
  corrected `.38 Wh` uniformly opened seed passed its preflight, but full
  replay and release gates later failed as recorded above. **Verdict:**
  superseded by the terminal Attempt 59 rejection; no candidate or promotion.

- [Attempt 56A State 2 head sheet](../../../../out/reimu_fumo_attempt_056a_sculpted_head_cushion/state_2/fixed_views/contact_sheet.png),
  [56B sewn-cushion comparison](../../../../out/reimu_fumo_attempt_056b_sewn_head_cushion/contact_sheet.png),
  and
  [56C skirt-state comparison](../../../../out/reimu_fumo_attempt_056c_sculpted_skirt_shell/packet/state_a_b_review.png)
  — dense sculpting remained a rounded mattress, the pressure solves became a
  tufted pillow or foam block, and the skirt shell became a padded bench/prow.
  **Verdict:** reject all visible A56 geometry; retain only measurements,
  replay methods, and the hidden lower-support control.

- [Attempt 55 head Cycle 2 gate](../../../../out/reimu_fumo_attempt_055_constructed_head_receiver/cycle_2/cheap_gate/contact_sheet.png)
  and
  [lower front/side/x-ray sheet](../../../../out/reimu_fumo_attempt_055_constructed_lower_receiver/packet/front_side_xray_review.png)
  — the head still has an angular helmet, bald crown, and floating parallel
  sheets; the lower support is broad, but its visible garment is a doubled
  tier and planar ramp. **Verdict:** reject both visible receivers; no
  integration or promotion.

- [Attempt 54 face comparison](../../../../out/reimu_fumo_attempt_054_face_applique/blind_review_sheet.png),
  [head-receiver comparison](../../../../out/reimu_fumo_attempt_054_hull_head_receiver/contact_sheet.png),
  and
  [hidden lower-receiver sheet](../../../../out/reimu_fumo_attempt_054_hull_lower_receiver/fixed_view_sheet.png)
  — the face hierarchy is a useful receiver-normalized direction below the
  absolute gate, both head receivers fail, and lower receiver B is only a
  hidden fitting control. **Verdict:** retain bounded evidence only; no
  candidate is approved.

- [Attempt 53 head curves](../../../../out/reimu_fumo_attempt_053_head_seam_network_spec/curve_contact_sheet.png),
  [lower interface diagram](../../../../out/reimu_fumo_attempt_053_lower_stack_panel_spec/overlays/coupon_interface_diagram.png),
  and
  [constructed-panel reference board](../../../../out/reimu_fumo_attempt_053_plush_panel_adversary/evidence/reference_construction_board.png)
  — direct pixels freeze independent patch roots, broad lower support, and
  receiver-owned contact rather than a shared loft or generic padded API.
  **Verdict:** frozen as a bounded review contract, not geometry.

- [Attempt 52 source-versus-hull sheet](../../../../out/reimu_fumo_attempt_052_visual_hull/evidence/source_vs_blender_hull.png)
  and [report](../../../../out/reimu_fumo_attempt_052_visual_hull/REPORT.md)
  — the loose guide reaches mean recall `.976` but fails the per-view recall
  and excess-area gates; its confidence band occupies `20.7%` of loose-hull
  volume. **Verdict:** reject as reconstruction, component separator, seam
  guide, or active model geometry; retain only coarse uncertainty evidence.

- [Attempt 49 bow-root report](../../../../out/reimu_fumo_attempt_049_bow_hub_pixel_gate/README.md),
  [Attempt 50 blind comparison](../../../../out/reimu_fumo_attempt_050_bow_root_skin/evidence/blind_ab_sheet.png),
  and [Attempt 51 blind comparison](../../../../out/reimu_fumo_attempt_051_bow_authored_panel/evidence/blind_ab_sheet.png)
  — A49 proves a connected, crown-seated four-neck scaffold, but it remains
  unsuitable as visible fabric. A50's offset skin became a perched slab and
  A51's authored panel became circular ears plus a profile fin. **Verdict:**
  retain A49 only as hidden support; reject both visible skins and stop this
  branch until an accepted crown receiver exists.

- [Cycle-23 foot-pod comparison](../../../../out/reimu_fumo_attempt_046_foot_pod_coupon/cycle_23/reference_comparison.png)
  and [review](../../../../out/reimu_fumo_attempt_046_foot_pod_coupon/REVIEW.md)
  — the pod reaches the measured width, depth, and height bands but reads as
  an end disc inside a band or a padded wheel. **Verdict:** reject the visible
  pod and resume feet only inside an accepted hem/connector receiver.

- [Traced lower-stack report](../../../../out/reimu_fumo_attempt_046_lower_stack_traced/REPORT.md)
  and [bounded cloth gate](../../../../out/reimu_fumo_attempt_046_lower_stack_cloth/GATE.md)
  — explicit traced panels produced a ramp and bench; two coarse cloth
  configurations produced tents with only four near-floor vertices. **Verdict:**
  retire both representations rather than tune their parameters.

- [Upper-body panel-junction comparison](../../../../out/reimu_fumo_attempt_048_upper_body_panel_junction/cycle_1/contact_sheet.png)
  and [review](../../../../out/reimu_fumo_attempt_048_upper_body_panel_junction/REVIEW.md)
  — near-perfect tangent and normal continuity still renders as a rigid padded
  slab. **Verdict:** reject the geometry and retain continuity only as a
  regression check.

- [Alternative unified-hair gate](../../../../out/reimu_fumo_attempt_046_unified_hair_alt/cycle_2/review_front_side_b/contact_sheet.png)
  and [report](../../../../out/reimu_fumo_attempt_046_unified_hair_alt/REPORT.md)
  — all requested scalar guards pass, but independent front/side panels form
  an open cuboid with wedge bangs and a detached rear strip. **Verdict:** reject
  both candidates; the next head representation needs shared three-dimensional
  seam curves rather than separately extruded cards.

- [Stuffed-panel proof](../../../../out/reimu_fumo_attempt_046_stuffed_panel_proof/contact_sheet.png)
  and [API rejection](../../../../out/reimu_fumo_attempt_046_stuffed_panel_proof/API_REJECTION.md)
  — broad-contact deformation reached `8/10`, but the general panel still has
  ambiguous construction and penetrates its floor. **Verdict:** retain only
  the bounded contact-weighting hypothesis, not the mesh or a reusable API.

- [Whole-model macro audit](../../../../out/reimu_fumo_attempt_045_working_ladder_review/MACRO_AUDIT_RUNG2.md)
  — implementation-blind measurement finds that eye spacing, bow span, skirt
  width, and foot-pair span are already broadly close; simple global rescaling
  would waste time. The dominant errors are representation, layering, and
  contact: helmet hair, upright lower stack, flat bow depth, rigid sleeves,
  and detached feet. **Verdict:** stabilize unified head/hair/crown first,
  then the seated lower stack; defer receiver-dependent polish.

- [Current working-ladder five-view render](../../../../out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/five_view_gate/contact_sheet.png)
  and [implementation-blind rung review](../../../../out/reimu_fumo_attempt_045_working_ladder_review/RUNG3_AB.md)
  — the complete eyes/locks plus constructed-sleeves candidate is preferred
  over rung 2 with no critical regression. Overall likeness rises from
  `5.6/10` to `5.9/10` and constructed-plush read from `3.7/10` to `4.5/10`.
  **Verdict:** provisional working-rung promotion only; retain the 30-object
  fabric-panel core, reshape cuff scaffolds, replace both arm plugs, and keep
  the tracked blend unchanged until the final absolute gate passes.

- [Attempt 45 integration gate](../../../../out/reimu_fumo_attempt_045_reference_eval/INTEGRATION_GATE.md)
  — implementation-blind review orders the dependency graph and vetoes
  wholesale merging of Face Cycle 7, the panel sleeve, bow, skirt, and head
  FFD. **Verdict:** visible hair and the seated lower stack are the two largest
  likeness levers; face appliques are receiver-dependent.

- [All-reference construction specification](../../../../out/reimu_fumo_attempt_044_reference_spec/CONSTRUCTION_SPEC.md),
  [ownership sheet](../../../../out/reimu_fumo_attempt_044_reference_spec/reference_ownership_sheet.png),
  and [30-frame turn sheet](../../../../out/reimu_fumo_attempt_044_reference_spec/canonical_turn_all_frames.png)
  — canonical front owns exact identity, the turn owns depth/layer order, and
  physical stills own sewn-medium behavior. **Verdict:** accepted as review
  authority, not as hidden-construction fact.

- [Adversarial process decision](../../../../out/reimu_fumo_attempt_044_process_review/DECISION.md)
  and [bounded arXiv review](../../../../out/arxiv_plush_modeling_research/REVIEW.md)
  — generic reconstruction cannot repair a wrong manufactured form. The one
  locally relevant FFD experiment worsened landmark RMS by about twelve times.
  **Verdict:** continue explicit constructed modules and multi-view gates; no
  external model-generation service.

- [Attempt 44 face Cycle 7 comparison](../../../../out/reimu_fumo_attempt_044_face_hairline/cycle_7/multiview_reference_contact.png)
  and [review](../../../../out/reimu_fumo_attempt_044_face_hairline/REVIEW.md)
  — improved front identity but retained helmet crown, paddle locks, tall
  eyes, and rigid contacts. **Verdict:** absolute reject; retain only
  receiver-dependent eye/mouth/fringe directions.

- [Attempt 45 support-swap front/side gate](../../../../out/reimu_fumo_attempt_045_accepted_support_integration/front_side_gate/contact_sheet.png)
  and [verdict](../../../../out/reimu_fumo_attempt_045_accepted_support_integration/VERDICT.md)
  — a previously role-accepted hidden cushion is not a drop-in surface under
  the old cap; large beige side/crown gaps appear. **Verdict:** reject and
  rebuild the support/hair receiver coherently.

- [Attempt 45 restricted face-subset gate](../../../../out/reimu_fumo_attempt_045_face_subset_integration/front_3q_gate_02/contact_sheet.png)
  and [verdict](../../../../out/reimu_fumo_attempt_045_face_subset_integration/VERDICT.md)
  — exact cherry-picking onto the old cap clips the eyes and breaks the
  fringe/lock transition. **Verdict:** reject; reproject only after the unified
  hair field passes.

- [Attempt 45 panel-sleeve comparison](../../../../out/reimu_fumo_attempt_045_sleeve_shell/panel_cycle_8/reference_contact_sheet.png)
  and [scorecard](../../../../out/reimu_fumo_attempt_045_sleeve_shell/SCORECARD.md)
  — separate sewn panels improve the nozzle representation, but side/rear
  remain small, high, planar, and too mechanical. **Verdict:** directional
  reject; preserve only measured contact/interface ideas.

- [Attempt 47 connected bow-root evidence](../../../../out/reimu_fumo_attempt_047_bow_connected_hub/evidence/reference_comparison.png)
  and [report](../../../../out/reimu_fumo_attempt_047_bow_connected_hub/README.md)
  — one connected hub passed topology and branch-path gates, but excessive
  crown burial hid the lower branches. **Verdict:** reject coordinates; retain
  the topology validator for a shallow-depth reset.

- [Attempt 47 ruffle excess-tuck comparison](../../../../out/reimu_fumo_attempt_047_ruffle_excess_tucks/reference_comparison.png)
  and [Attempt 50 bicubic reset](../../../../out/reimu_fumo_attempt_050_ruffle_curved_pockets/reference_comparison.png)
  — true excess length/depth became cardstock, batch cloth shredded, and the
  final pocket strip flattened. **Verdict:** stop the detailed ruffle branch
  until the seated lower receiver exists.

- [Attempt 45 material evidence](../../../../out/reimu_fumo_attempt_045_fabric_materials/review/material_evidence_sheet_final.png)
  and [report](../../../../out/reimu_fumo_attempt_045_fabric_materials/REPORT.md)
  — directional nap/weave improvements preserved all geometry, but blind
  review still scored only `4–5.5/10`. **Verdict:** defer materials until the
  constructed geometry passes.

- [Attempt 36 W0 wire/topology rejection sheet](../../../../out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/wire_preflight_sheet.png),
  [machine metrics](../../../../out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/wire_metrics.json),
  and
  [cage manifest](../../../../out/reimu_fumo_attempt_036_sparse_quad_cage/preflight/cage_manifest.json)
  — the disposable `13 × 9` cage has its promised exact topology but fails at
  negative minimum scaled Jacobian, severe edge ratios, `4/52/77` projected
  crossings, and `.11690/.14289 Wh` front-contour RMS/maximum. The sheet makes
  the absolute failure visible: a narrow pointed crown, pole-like apex pinch,
  folded bilateral three-quarter flow, and abrupt wall-like side returns.
  **Verdict:** NO-GO at W0; no persistent blend, G0, or G1. Retire the coupled
  rectangular single-surface premise rather than densifying this cage.

- [Attempt 35 first-gate rejection sheet](../../../../out/reimu_fumo_attempt_035_coupled_front_skin/attempt_035_rejection_contact_sheet.png),
  [G1 beauty](../../../../out/reimu_fumo_attempt_035_coupled_front_skin/g1_front_beauty.png),
  and [persistent rejected blend](../../../../out/reimu_fumo_attempt_035_coupled_front_skin/reimu_fumo_attempt_035_coupled_front_skin.blend)
  — the polar single-surface proof renders as a horizontally banded generic
  helmet with visible beige `.802 Wh` wide against the `.603 Wh` target. Its
  semantic raster reports three brown components, and implementation review
  finds a high-valence center pole, homothetic radial courses, and incorrect
  rear ownership. **Verdict:** hard reject at G0/G1; preserve the 15-second
  early-veto workflow, not this representation or geometry.

- [Goal workflow with remote checkpoint and resource gates](../../../../.agents/skills/goal/SKILL.md)
  — requires a verified commit and push after every active-goal turn, allocates
  workers from the dependency graph instead of the slot count, caps work in
  progress to review capacity, and cancels work after a stale input or failed
  veto. It also preserves a protected main asset while hash-addressed candidate
  copies are tested and approved, and rejects premature library-file splitting.
  **Verdict:** packaged skill build and root Buildifier test pass.

## 2026-08-29

- [Attempt 34 first-gate rejection sheet](../../../../out/reimu_fumo_attempt_034_front_crown_field/attempt_034_rejection_contact_sheet.png),
  [build and review notes](../../../../out/reimu_fumo_attempt_034_front_crown_field/BUILD_NOTES.md),
  and [persistent rejected blend](../../../../out/reimu_fumo_attempt_034_front_crown_field/reimu_fumo_attempt_034_front_crown_field.blend)
  — the broad front field is present, but its pixels remain a faceted rigid
  cap with a center puncture, hanging side strips, a dark contact channel, and
  exposed bilateral closures. Raw gates also find `210` forbidden fringe
  pixels, `686` guide-halo pixels, `925/843` visible closure pixels at
  `+48/-48`, `70` boundary edges, and `124` invalid-incidence edges.
  **Verdict:** hard reject; retain only measured controls and fixed cameras,
  not this mesh, mapping, guide, or closure construction.

- [Whole-result baseline context audit](../../../../out/reimu_fumo_resume_033/baseline/WHOLE_RESULT_CONTEXT.md)
  — prevents the active hair panel from becoming the whole goal: the seated
  footprint remains 35–40 percent too short, the bow is too narrow/upright,
  and the rear hair leaf is missing, while head-base and foot size are already
  near target. It freezes the post-head module order and the attachment space
  Attempt 34 must preserve.

- [Attempt 34 source-relative crown/fringe control rows](../../../../out/reimu_fumo_resume_033/front_panel_design/CONTROL_ROWS.md)
  — freezes the corrected crown courses, the full asymmetric `L0…T2` free
  edge, depth/contact bands, no-go face volume, and three-quarter vetoes before
  geometry. **Verdict:** buildable only as one filled paired-skin disk with
  three lower subpatches; a rectangular grid or full row below the clefts is
  prohibited.

- [Attempt 33 first-pixel rejection sheet](../../../../out/reimu_fumo_attempt_033_visual_hull_cage/early_rejection_contact_sheet.png),
  [build and absolute-review notes](../../../../out/reimu_fumo_attempt_033_visual_hull_cage/BUILD_NOTES.md),
  and [persistent rejected blend](../../../../out/reimu_fumo_attempt_033_visual_hull_cage/reimu_fumo_attempt_033_visual_hull_cage.blend)
  — scalar bounds matched the corrected targets, but pixels scored only
  `1/10` likeness. The connected front-to-rear field became a thick suspended
  W fringe and broad rear cape, exposed a massive beige crown, formed an
  under-chin rail, and left pill-shaped locks floating. **Verdict:** reject at
  the cheap head-only gate; do not render, texture, or tune this representation.

- [Corrected canonical head targets](../../../../out/reimu_fumo_resume_033/landmarks/HEAD_TARGETS.md)
  and [turntable view-role dossier](../../../../out/reimu_fumo_resume_033/references/TURNTABLE_VIEW_ROLES.md)
  — the previous `Wh = 395 px` box included red/background. The controlling
  datum is `Wh = 368 ± 4 px`; `1.098 Wh` is a lock-inclusive composite rather
  than a lock-excluded cage height. Profile depth is a `.77–.85 Wh` shallow
  base plus a crown-rooted leaf reaching
  `1.14–1.23 Wh`, not one inflated volume.

- [Tracked-baseline fixed-view review sheet](../../../../out/reimu_fumo_resume_033/baseline/baseline_review_sheet.png)
  and [absolute baseline verdict](../../../../out/reimu_fumo_resume_033/baseline/BASELINE_VERDICT.md)
  — a read-only background render of the current standalone blend exposes the
  upright doll proportions, rounded shell head, flat fringe, parallel side
  walls, and bald rear. It is the before-image and a composition regression,
  not acceptance evidence.

- [Attempt 32 candidate comparison](../../../../out/reimu_fumo_attempt_032_visible_head/initial_candidate_comparison.png)
  and [absolute rejection](../../../../out/reimu_fumo_resume_033/attempt32_rejection/VERDICT.md)
  — two independent pixel audits reject V1 and V1-alt at about `2/10`
  likeness. V1 is a short inflated hood around a face ball; V1-alt is a card
  fringe on a voxel helmet. Neither geometry is retained for Attempt 33.

- [Calibrated support-versus-hair envelope diagnostic](../../../../out/reimu_fumo_attempt_029_front_top_overlap/support_hair_interface_audit/diagnostic.png),
  [geometric verdict](../../../../out/reimu_fumo_attempt_029_front_top_overlap/support_hair_interface_audit/VERDICT.md),
  and [source-only verdict](../../../../out/reimu_fumo_attempt_029_front_top_overlap/support_hair_interface_audit/REFERENCE_VERDICT.md)
  — evaluated support V3 projects outside exact O over `82.18%` of its arc,
  reaches `-.122583 BU`, and fails `14/17` keyed landmarks. **Verdict:** reset
  only the hidden crown/temple receiver; do not widen O or inflate the hair.

- [Alternative Stage-A four-view Blender render](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_front_lift_alt/contact_sheet.png)
  and [verdict](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_front_lift_alt/VERDICT.md)
  — a continuous native-gusset O lift reduces the inherited `.4591 BU`
  depth jump to `.00397 BU` and passes orientation/edge-ratio gates, but exact
  O samples alternate from `-.12238` to `+.04613 BU` relative to the support.
  Pixels show a severe comb/sawtooth crown and temple crossings. **Verdict:**
  reject constant-native-v lifting; the interface needs a jointly solved
  exterior bridge or a support-interface reset.

- [Regularized-seating Stage-A four-view Blender render](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_front_lift_witness/contact_sheet.png)
  and [reset diagnosis](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_front_lift_witness/RESET.md)
  — all seated collar samples clear the support by
  `.00444..00615 BU` with zero overlaps and positive orientation, but the
  independently solved midpoint depths tear O at `52.64°` tangent and
  `46.47°` split-normal maximum error. **Verdict:** reject the local solve;
  solve the complete narrow strip jointly with fairness and derivative terms.

- [Rejected V3 radial/geodesic preflight diagnostic](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_back_skin_preflight/geodesic_back_skin_diagnostic.png),
  [hard-gate verdict](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_back_skin_preflight/VERDICT.md),
  [native-support replacement contract](../../../../out/reimu_fumo_attempt_029_front_top_overlap/v3_back_skin_preflight/NATIVE_SUPPORT_MAPPING.md),
  and [candidate stop evidence](../../../../out/reimu_fumo_attempt_029_front_top_overlap/candidate_v3_back_skin/STOP_EVIDENCE.md)
  — the executed fixed-cardinality field inverts, crowds the rear apex, and
  intersects both accepted parents. **Verdict:** rejected before Blender mesh
  creation. Preserve exact front X/Z/topology, reset only the O-near depth
  collar, and use hidden all-quad center reduction on the native support atlas.

- [Current V2 progress render](../../../../out/reimu_fumo_attempt_029_front_top_overlap/candidate_v2/progress_contact_sheet.png),
  [technical report](../../../../out/reimu_fumo_attempt_029_front_top_overlap/candidate_v2/validation_report.md),
  and [diagnostic blend](../../../../out/reimu_fumo_attempt_029_front_top_overlap/candidate_v2/reimu_fumo_front_top_overlap_v2.blend)
  — the corrected terminal itself reaches the intended V4 clearance, but the
  full O-to-U map crosses the support and tears at the temples. **Verdict:**
  rejected at first geometry and pixel gates; reset the 3D surface mapping,
  not the accepted parents.

- [Corrected exterior-lap interface contract](../../../../out/reimu_fumo_attempt_029_front_top_overlap/interface_reset/BUILDABLE_CONTRACT.md)
  and [layer/ownership diagnostic](../../../../out/reimu_fumo_attempt_029_front_top_overlap/interface_reset/exterior_lap_diagnostic.png)
  — retire the impossible support-gap underlap and freeze support < V4 < new
  lap < future occluders. **Verdict:** BUILDABLE for a bounded shape candidate;
  not promotable as a completed interface until final accepted bow/side-band
  geometry hides the complete terminal.

- [Independent hidden-U interface reset](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/INTERFACE_PREFLIGHT_REJECTED.md)
  — two read-only reviews agree that the proposed terminal can cross the
  accepted V4 shell and that its occlusion, integrated silhouette, and
  evaluated-continuity gates are incomplete. **Verdict:** reject/reset the
  depth and gate contract; retain the source-keyed chart; Candidate V1 may be
  diagnostic only.

- [Source-exact non-folding front-field chart](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/parameterization_chart.png)
  and [machine metrics](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/parameterization_metrics.json)
  — the three shared semantic patches preserve all `39/39` O/F keys exactly
  under the supporting physical-front `Wh = 189 px` datum, form one all-quad
  disk, and pass the original Jacobian/aspect gates without weakening.
  **Verdict:** accepted for that physical-front 3D mapping only; support
  seating, inferred U overlap, and pixels still require Blender gates.

- [Accepted-parent interface depth diagnostic](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/interface_depth_diagnostic.png),
  [depth/ownership contract](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/interface_depth_contract.md),
  and [preflight scorecard](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/preflight_scorecard.md)
  — actual V3/V4 geometry proves that V4's crown back skin is already only
  `.006 BU` outside the support. **Verdict:** V4 remains fully immutable. The
  then-proposed support-clamped successor terminal is retained only as the
  measured premise later rejected by the independent interface audit. The
  lower world-`-X` hook remains later-lock ownership, not visible front/top
  geometry.

- [All-relevant-reference front/top hair board](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/all_relevant_reference_board.png)
  and [per-view role matrix](../../../../out/reimu_fumo_attempt_029_front_top_overlap/preflight/reference_role_matrix.md)
  — both front stills, the physical side, all four turn frames, and all eight
  selected sofa frames are now explicitly assigned to silhouette, depth,
  overlap, felt/pile, or qualitative-regression roles. **Verdict:** accepted as
  Attempt 29 source control. It rejects single-view authorization, a full-
  width hidden crown seam, and the stale `Wh = 178 px` normalization.

- [Accepted rear felt-panel V4 review sheet](../../../../out/reimu_fumo_attempt_028_rear_felt_panel/v4/review_sheet_accepted.png),
  [implementation-blind verdict](../../../../out/reimu_fumo_attempt_028_rear_felt_panel/v4/BLIND_VERDICT.md),
  [registered contour overlay](../../../../out/reimu_fumo_attempt_028_rear_felt_panel/v4/rear_contour_overlay.png),
  [technical report](../../../../out/reimu_fumo_attempt_028_rear_felt_panel/v4/validation_report.md),
  and [candidate blend](../../../../out/reimu_fumo_attempt_028_rear_felt_panel/v4/reimu_fumo_rear_felt_panel_v4.blend)
  — rear silhouette IoU is `.977984`, contour RMS `.00506 Wh`, and blind
  scores are `9.5/8.8/8.1/8.0/8.1/8.2`. **Verdict:** accepted only as the
  isolated rear-hair parent. Freeze its rear XY contour; the front/top field
  must bury the proud crown tab and hard front rim and resolve the grazing-left
  hook in depth.
- [Goal skill with whole-result module guard](../../../../.agents/skills/goal/SKILL.md)
  — small modules remain the default work unit, but each is now evaluated
  against overall composition, shared proportions, adjacent interfaces,
  downstream use, the dominant unresolved failure, and a lightweight
  integrated regression when available. **Verdict:** `quick_validate.py`, the
  packaged Bazel skill build, and the root Buildifier test pass.

- [V3.4 source/candidate review sheet](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/review_sheet.png),
  [front](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/renders/front.png),
  [rear](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/renders/rear.png),
  [side](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/renders/side.png),
  [technical report](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/validation_report.md),
  and [candidate blend](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/reimu_fumo_hair_base_v3.blend)
  — exact four-patch topology survives while the crown/rear field is seated on
  the evaluated cushion. The inner supported skin has at least `.02758 BU`
  radial clearance, the crown leak is gone, and rear-camera handedness is
  corrected. The [implementation-blind verdict](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_4/BLIND_VERDICT.md)
  scores overall likeness `4/10`, plush construction `3/10`, and medium read
  `4/10`. **Verdict:** rejected; do not use it as a parent. Reset to a shallow
  source-shaped rear felt panel before rebuilding the front/top overlap.
- [V3.3 verdict](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_3/VERDICT.md)
  — rear handedness and visible-hem smoothing were corrected, but a pale crown
  slit proved that a straight chord between two valid support points still
  crosses a convex cushion. **Verdict:** reject the five-row chord transition.
- [V3.2 verdict](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_2/VERDICT.md)
  — direct support wrapping materially reduced the inflated helmet volume and
  seated the layer, but decorative seam protrusion, rear handedness, and hem
  sharpness remained. **Verdict:** retain the representation, not the bytes.
- [V3.1 verdict](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3_1/VERDICT.md)
  and [V3 verdict](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3/VERDICT.md)
  — V3 passed exact topology but exposed large beige rear wedges; V3.1 closed
  the direct rear while retaining a helmet/plate side opening. **Verdict:**
  retire analytic ellipsoid/superellipse seating and map to the real support.
- [V3 source-honest boundary contract](../../../../out/reimu_fumo_attempt_027_front_hair_field/v3/boundary_contract.md)
  — freezes exact central, asymmetric temple, shared shoulder, and rear-hem
  boundaries plus the `4,845 V / 4,704 Q` base and `10,530 V / 10,528 Q`
  closed paired topology. Rear camera handedness is explicit. **Verdict:**
  accepted as construction input, not as a factory-pattern claim.

- [V2 chart preflight rejection](../../../../out/reimu_fumo_attempt_027_front_hair_field/v2/PREFLIGHT_REJECTED.md)
  — Blender's pre-mesh Jacobian gate caught crossed left/center trajectories;
  correcting source-X correspondence and delaying width expansion did not
  remove the fold. **Verdict:** retire the single chart and split it into
  shared central, temple-return, and crown/rear patches.

- [Continuous hair-base source contract](../../../../out/reimu_fumo_attempt_027_front_hair_field/hair_base_source_contract.md)
  — follow-up source review confirms that a coupled visible base assembly may
  own the reliable front edge plus full crown/rear brown coverage without
  claiming hidden factory seams. **Verdict:** proceed with a bounded pole-free
  front/rear macro candidate; cheek/nape locks remain later modules.

- [Stopped V1 builder preflight](../../../../out/reimu_fumo_attempt_027_front_hair_field/v1/PREFLIGHT_REJECTED.md)
  — implementation review found an unreliable right-edge continuation, two
  artificial crown-spanning side borders, unproved world-space coverage, and
  uncontrolled temple n-gons before Blender ran. **Verdict:** rejected with no
  candidate geometry or render cost; reset the module interface first.

- [Front-hair source-boundary board](../../../../out/reimu_fumo_attempt_027_front_hair_field/source_boundary_audit/source_views.png)
  and [audit decision](../../../../out/reimu_fumo_attempt_027_front_hair_field/source_boundary_audit/README.md)
  — four controlling reference views show continuous crown/fringe pile and no
  enclosed central-bang root. **Verdict:** reject an isolated bang; authorize
  one complete continuous front crown/fringe field as the smallest valid
  module. Separate overlapping cheek/nape locks remain later work.

- [Rejected coupon V2 reference comparison](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v2/reference_comparison.png),
  [decisive-view sheet](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v2/candidate_review_sheet.png),
  [implementation-blind verdict](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v2/blind_verdict.md),
  [technical report](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v2/validation_report.md),
  and [candidate blend](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v2/reimu_fumo_front_crown_crop_coupon_v2.blend)
  — the broad explicit crop is one closed manifold with the frozen cushion and
  tracked assets unchanged, but blind scores are seating/contact `1.5/10`,
  fabric read `2.5/10`, and crop closure `2/10`. **Verdict:** rejected early.
  The underside cavity, pointed artificial closures, and rectangular card read
  falsify the crop module/interface; front and neutral three-quarter renders
  were deliberately not spent on this rejected candidate.
- [Goal-skill early-return clarification](../../../../.agents/skills/goal/SKILL.md)
  — a falsified approach, module, boundary, or interface is now explicitly a
  valid turn result after its rejection artifact and reset are recorded; the
  same turn need not force a substitute candidate. **Verdict:** quick validator
  and packaged `bazel_agent` build pass.
- [Rejected coupon V1 reference comparison](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v1/reference_comparison.png),
  [fixed-view review sheet](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v1/review_sheet.png),
  [absolute verdict](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v1/verdict.md),
  [technical validation](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v1/validation.json),
  and [candidate blend](../../../../out/reimu_fumo_attempt_026_front_crown_wrap_coupon/coupon_v1/reimu_fumo_front_crown_wrap_coupon_v1.blend)
  — the explicit coupon is one finite closed manifold and crosses the apex, but
  independent scores are contact `4/10`, apex continuity `6/10`, and stuffed-
  fabric construction `3.5/10`. **Verdict:** rejected early. Its continuous
  crown air channel and rounded rail section read as a floating rigid strap;
  the next attempt changes to a broader thin panel patch with a seated back
  surface rather than tuning or expanding this strip.
- [Rejected front-hair V4 source comparison](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v4/reference_comparison.png),
  [fixed-view review sheet](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v4/review_sheet.png),
  [absolute verdict](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v4/verdict.md),
  [technical validation](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v4/validation.json),
  and [candidate blend](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v4/reimu_fumo_front_hair_field_v4.blend)
  — explicit paired front/back surfaces remove the modifier spike and form one
  clean closed manifold, but independent scores are front identity `4/10`,
  side construction `2/10`, three-quarter attachment `2/10`, and plush read
  `3/10`. **Verdict:** rejected early. The front-only offset crescent leaks the
  beige crown and floats like a rigid visor, so its module/interface is retired
  before any adjacent part is built.
- [Accepted head-support cushion V3 three-quarter render](../../../../out/reimu_fumo_attempt_025_modular_parts/head_support_cushion_v3/renders/three_quarter.png),
  [front](../../../../out/reimu_fumo_attempt_025_modular_parts/head_support_cushion_v3/renders/front.png),
  [side](../../../../out/reimu_fumo_attempt_025_modular_parts/head_support_cushion_v3/renders/side.png),
  [technical validation](../../../../out/reimu_fumo_attempt_025_modular_parts/head_support_cushion_v3/validation.json),
  and [candidate blend](../../../../out/reimu_fumo_attempt_025_modular_parts/head_support_cushion_v3/reimu_fumo_head_support_cushion_v3.blend)
  — the work was reduced to one hidden structural module. **Verdict:** accepted
  and frozen for its role. Blind visual scores are `8–8.5/10`; the mesh is one
  closed manifold with zero boundary, wire, or non-manifold edges, preserved
  interfaces, and in-band `2.645` depth.
- [Rejected front-hair field V3 front render](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v3/renders/front.png),
  [side](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v3/renders/side.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v3/renders/three_quarter.png),
  and [metrics](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v3/metrics.json)
  — the source-registered continuous fringe survives, but the curved panel's
  modifier-derived thickness emits long needles and expands evaluated depth to
  an impossible range. **Verdict:** rejected early. Retire `Solidify` for this
  panel and rebuild only it as explicit paired sewn surfaces next turn.
- [Rejected combined head V5 front](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v5/renders/front.png),
  [three-quarter](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v5/renders/three_quarter.png),
  [side](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v5/renders/side.png),
  and [rear](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v5/renders/rear.png)
  — blind head-only review scored likeness `4/10`, proportions `4/10`, face
  `3/10`, hair `3/10`, plush read `3/10`, and rear completeness `3/10`.
  **Verdict:** rejected; rounded helmet, raised eye shields, slab locks, and a
  blank rear forced a representation reset.
- [Rejected combined head V6 front](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v6/renders/front.png),
  [side](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v6/renders/side.png),
  and [rear](../../../../out/reimu_fumo_attempt_024_head_module/head_module_v6/renders/rear.png)
  — source tracing improved the frontal hairline, but bundling cushion, front
  and rear hair, locks, face, and materials produced another rigid box, a
  failed rear field, and pasted graphics. **Verdict:** rejected; this attempt
  directly triggered the small-module workflow.

- [Source-owned front projection reset](../../../../out/reimu_fumo_attempt_022_pattern_head/front_hair_trace_reset_v1/physical_front_contours.png),
  [multi-source comparison](../../../../out/reimu_fumo_attempt_022_pattern_head/front_hair_trace_reset_v1/comparison_board.png),
  [schema-2 contours](../../../../out/reimu_fumo_attempt_022_pattern_head/front_hair_trace_reset_v1/contours.json),
  and [source decision](../../../../out/reimu_fumo_attempt_022_pattern_head/front_hair_trace_reset_v1/README.md)
  — the physical `Wh = 189 px` source now freezes one continuous asymmetric
  three-span front edge, including the broad center low band around
  `(235, 211)`, while all controlling side/rear views freeze complete brown
  coverage and layer order. **Verdict:** accepted only for a front/rear macro
  blockout. Literal one-piece brown-cushion, beige-applique, and occluded
  seam/root claims are explicitly unproven and retired.
- [Conditionally accepted sewn-panel topology](../../../../out/reimu_fumo_attempt_022_pattern_head/conforming_cap_preflight_v3/conforming_cap_construction.png),
  [editable SVG](../../../../out/reimu_fumo_attempt_022_pattern_head/conforming_cap_preflight_v3/conforming_cap_construction.svg),
  and [preflight report](../../../../out/reimu_fumo_attempt_022_pattern_head/conforming_cap_preflight_v3/README.md)
  — a brown front annulus, nearly flush visible beige face insert, and brown
  all-quad multi-patch rear disk conditionally close one shallow sewn volume;
  the construction is a reversible implementation hypothesis, not a claim
  about occluded manufactured seams. **Verdict:** authorized for one bounded
  neutral front/rear macro candidate under fail-closed topology and pixel
  gates; no geometry is accepted yet.
- [Frozen coupled-head review contract](../../../../out/reimu_fumo_attempt_022_pattern_head/unified_cushion_head_v1_review_contract/README.md)
  — written and source-corrected before candidate pixels. It requires
  full-frame unlabelled front/rear first, complete brown ownership, a soft
  non-mask face, one broad asymmetric fringe edge, rounded/tapered rear, and
  absolute `8/10` implementation-blind scores before downstream views.
- [Rejected flat-pattern gravity true-side beauty](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/flat_pattern_gravity_v2/renders/beauty/true_side.png),
  [source-aligned overlay](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/flat_pattern_gravity_v2/renders/diagnostic/source_aligned_overlay.png),
  [absolute verdict](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/flat_pattern_gravity_v2/VERDICT.md),
  and [cycle report](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/flat_pattern_gravity_v2/README.md)
  — both cuts are planar, only `82/82` waist vertices are pinned, `62`
  face-free edges sew the side seams, the `80`-edge hem is open, motion
  converged, contacts passed, and the candidate clean-reopened. **Verdict:**
  rejected at the implementation-blind true-side gate. Silhouette `3/10`,
  construction `3/10`, and intended medium `2/10`; it is a short glossy
  lampshade with a straight roof, wrong early fall, weak rear pool, and noisy
  ground stack. The homogeneous gravity-only family is retired.
- [Rejected constructed-panel hair contact sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_constructed_panels_v2/candidate_01/contact_sheet.png),
  [reference comparison](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_constructed_panels_v2/candidate_01/reference_comparison.png),
  [trace alignment](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_constructed_panels_v2/candidate_01/trace_alignment.png),
  and [absolute verdict](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_constructed_panels_v2/candidate_01/VERDICT.md)
  — the evaluated-support cap has `6,912` quads; all fields are closed
  all-quad meshes; minimum evaluated clearances are positive; the registered
  edge error is `0.0 Wh`; and `0/14,425` rear samples are exposed. **Verdict:**
  rejected at front/rear pixels. The square side walls, hard equator, centered
  V, rear curtain/teeth, and fused locks retire both the rectangular-field
  representation and the current trace as a 3D authorization. No side or
  three-quarter render was made.
- [Tracked standalone structure audit](../../../../out/reimu_fumo_standalone_structure_audit_v1/README.md)
  and [machine-readable results](../../../../out/reimu_fumo_standalone_structure_audit_v1/audit.json)
  — read-only inspection and a factory-empty append test of the exact tracked
  asset. **Verdict:** the `FUMO` append boundary, hierarchy/scale conventions,
  and self-contained packaging are reusable; the geometry is not. `Fumo_Rig`
  is only an Empty, with zero armatures, actions, vertex groups, constraints,
  drivers, or shape keys. The evaluated primary-form probe found `35`
  crossing pairs, including skirt/foot and ruffle/foot clipping.
- [Authoritative snapshot/render READY marker](../../../../out/live_checkpoint_harness_smoke_v5/READY),
  [manifest](../../../../out/live_checkpoint_harness_smoke_v5/manifest.json),
  [front beauty](../../../../out/live_checkpoint_harness_smoke_v5/beauty/front.png),
  [rear beauty](../../../../out/live_checkpoint_harness_smoke_v5/beauty/rear.png),
  [front silhouette](../../../../out/live_checkpoint_harness_smoke_v5/silhouette/front.png),
  and [rear silhouette](../../../../out/live_checkpoint_harness_smoke_v5/silhouette/rear.png)
  — exact-parent snapshot from Blender 5.1.1 followed by scripts-disabled,
  repository-pinned Blender 5.2.1 rendering. **Verdict:** tooling smoke passed;
  candidate save plus hash/reverification occupied `.024597 s`, the complete
  live helper returned in `.047686 s`, and pinned rendering through manifest
  inputs took `8.534518 s`. This renders a known rejected parent solely to
  validate the harness; it is not positive model evidence. The stale v4 packet
  correctly failed strict schema validation and produced no lock or output.
- [Rejected support-derived hair contact sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_surface_shell_v1/candidate_01/contact_sheet.png),
  [reference comparison](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_surface_shell_v1/candidate_01/reference_comparison.png),
  [trace alignment](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_surface_shell_v1/candidate_01/trace_alignment.png),
  and [absolute verdict](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_surface_shell_v1/candidate_01/VERDICT.md)
  — a continuous `96 × 32` support-derived shell passed topology, trace,
  hash, crown-coverage, and evaluated vertex/face-center clearance gates.
  **Verdict:** rejected; variable-height circumferential rows collapse into
  central sails and spikes, locks float, the front is a centered V, and the
  rear is a helmet above a narrow bundle. Retain only the crown construction.
- [Rejected sector-skirt true-side beauty](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/sector_side_proof_v1/renders/beauty/true_side.png),
  [source-aligned overlay](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/sector_side_proof_v1/renders/diagnostic/source_aligned_overlay.png),
  [metrics](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/sector_side_proof_v1/proof_metrics.json),
  and [absolute verdict](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/sector_side_proof_v1/VERDICT.md)
  — exact calibrated boundary proof in pinned Blender 5.2.1. **Verdict:**
  rejected despite IoU `.85987`, `4 px` maximum roof excess, zero violating
  columns, and zero ground gap; the surface still reads as a taut diagonal
  cape, so target-remapped monotonic meridians are retired.
- [Recovered hair checkpoint 02 contact sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_02_after_crash/contact_sheet.png),
  [absolute verdict](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_02_after_crash/VERDICT.md),
  and [manifest](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_02_after_crash/manifest.json)
  — source-traced filled panels reconstructed from the exact hashed parent in
  isolated background Blender after the T3 crash. **Verdict:** rejected;
  triangle centers penetrate by up to `.0434 Wh`, the crown remains exposed,
  and the pixels show a centered bang, cards, and a generic helmet/curtain.
- [Physical-side skirt fall calibration](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/physical_side_fall_calibration_composited_v2.png),
  [measurements](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/fall_gate_calibration.json),
  and [parent decision](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/PARENT_REVIEW.md)
  — direct trace of the controlling side photo. **Verdict:** the proposed
  `60%` early-fall gate is false; the source measures `18.91%`. The full
  registered boundary plus `4 px` now controls the approved two-cut proof.
- [Skirt pattern-topology preflight](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/README.md),
  [flat cuts and drape diagram](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/flat_cut_and_drape_topology.png),
  [source board](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/source_evidence_board.png),
  and [first-mask contract](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/pattern_topology_reset/sleeve_context_and_first_mask_gate.png)
  — two gathered sector/trapezoid halves, waist-only pins, angled side seams,
  and one open free hem. **Verdict:** approved for a base-skirt side proof
  after calibrating its provisional fall ratio directly from the physical-side
  trace; no Blender candidate was made during this preflight.
- [Invalid live-checkpoint smoke manifest](../../../../out/live_checkpoint_harness_smoke_v2/manifest.json),
  [front beauty](../../../../out/live_checkpoint_harness_smoke_v2/beauty/front.png),
  [front silhouette](../../../../out/live_checkpoint_harness_smoke_v2/silhouette/front.png),
  and [pinned-reopen manifest](../../../../out/live_checkpoint_harness_smoke_v2/pinned_reopen/manifest.json)
  — preserved negative evidence from the superseded live-render harness.
  **Verdict:** invalid, not passed; the manifest names a live cloth-sack source
  while claiming the hair parent. It cannot establish ancestry. The replacement
  is snapshot-only and must prove exact parent identity before pinned rendering.
- [Approved front hair trace](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/source_trace_preflight/front_trace_overlay.png),
  [source-and-mask comparison](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/source_trace_preflight/front_source_and_mask.png),
  [rear trace](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/source_trace_preflight/rear_trace_overlay.png),
  and [trace report](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/source_trace_preflight/README.md)
  — editable, registered paths preserve the broad center bang, unequal side
  masses and locks, tapered rear, and four-to-five irregular lobes. **Verdict:**
  accepted as the next pattern input; no 3D hair is accepted yet.
- [Gathered-cloth checkpoint B](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_b/front_side_3q_checkpoint.png)
  and [metrics](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_b/checkpoint_metrics.json)
  — actual `21 × 19` cloth, animated gathering, collision, open hem, and
  separate ruffle. **Verdict:** rejected in `11.146 s`; the center falls and
  pools, but taut off-center columns still own a triangular side cape.
- [Macro-drape checkpoint C side](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_c/side_first_checkpoint.png),
  [edge overlay](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_c/renders/diagnostic/side_edge_overlay.png),
  and [metrics](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_c/side_gate_metrics.json)
  — full-surface correction retained local folds and stopped at the first
  side view. **Verdict:** rejected in `4.20 s`; `12` rear columns exceed the
  target roof by up to `15 px`, retiring the rectangular rear-sheet family.
- [Whole-process iteration audit](../../../../out/reimu_fumo_process_audit_v1/README.md)
  — measured reference work, live MCP, pinned rendering, low/full packets,
  review, logging, Bazel, Git, remote latency, and hashing. **Verdict:** source
  ownership, representation choice, and handoff dominate; keep every final
  quality gate while adding source packets, task capsules, candidate snapshots,
  immediate blind vetoes, and manifest-driven evidence.
- [Hair Stage 3C checkpoint 01 front](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_01_front.png),
  [rear](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_01_rear.png),
  [front hair-only mask](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_01_front_hair_only.png),
  and [metrics](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v3/checkpoint_01_metrics.json)
  — one uninterrupted source-outline rear panel. **Verdict:** tapered sides
  and no equator are retained, but the surface cuts through the head, the
  front remains a V, and the chest sphere/rear lobe contaminate the front;
  rejected in `16.80 s` before a save or later views.
- [Source-owned garment checkpoint A](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_a/front_side_3q_checkpoint.png)
  and [measurements](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/attempt_a/checkpoint_metrics.json)
  — compact support, separate authored front/rear panels, open hem, ruffle,
  legs, and sewn foot pods. **Verdict:** rejected in `8.63 s`; subdivision
  erased the authored descent-to-pool bend and restored a taut cape, necked
  front, and weak foot contact. No masks or full packet were made.
- [Hair Stage 3B checkpoint 02 front](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_02_front.png),
  [rear](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_02_rear.png),
  [front silhouette](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_02_front_silhouette.png),
  and [rear silhouette](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_02_rear_silhouette.png)
  — fitted cap plus separate full-width nape test. **Verdict:** rejected in
  `14.20 s`; crown holes, a hard horizontal equator, vertical curtain sides,
  repeated teeth, and a symmetric-V bang retire this representation before a
  save or side/three-quarter render.
- [Hair Stage 3B checkpoint 01 front](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_01_front.png),
  [rear](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_01_rear.png),
  [front silhouette](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_01_front_silhouette.png),
  and [rear silhouette](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v2/checkpoint_01_rear_silhouette.png)
  — first source-measured contour reset. **Verdict:** rejected before a save;
  its `1.2993` rear height ratio passes, but the open crown, curtain sides,
  three-lobe hem, rear-lobe front contamination, and weak central bang fail.
- [Independent hair comparison board](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_reference_audit_v2/comparison_board.png)
  and [source audit](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_reference_audit_v2/report.md)
  — physical front/side plus turntable and sofa-GIF construction review.
  **Verdict:** controls the next reset with an asymmetric three-mass fringe,
  broad cheek locks, layered side depth, complete rear coverage, and an
  irregular lower rear contour.
- [Hair Stage 3A contact sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v1/contact_sheet.png)
  and [deterministic review gallery](../../../../out/reimu_fumo_attempt_022_pattern_head/hair_construction_v1/review/review.html)
  — exact crown/back ancestry, continuous fringe, rooted locks, and complete
  rear coverage. **Verdict:** retain topology evidence only; the three soft
  scallops, thin locks, and ball-like rear are visibly generic and rejected.
- [Source-owned garment construction sheet](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/construction_hypothesis.png),
  [front annotation](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/front_annotation.png),
  [side annotation](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/side_annotation.png),
  and [preflight report](../../../../out/reimu_fumo_attempt_023_pattern_body/source_owned_garment_reset/README.md)
  — reassigns the physical side silhouette to a compact support, separate
  hanging front and pooled rear panels, open hem, independent ruffle, short
  legs, and sewn pods. **Verdict:** accepted as the next body construction
  hypothesis; deep core, cone, disc, cap, and waist-to-tail roof are retired.
- [Seated-core sculpt Attempt B sheet](../../../../out/reimu_fumo_attempt_023_pattern_body/seated_support_sculpt_v1/attempt_b/front_side_3q_checkpoint.png)
  and [Attempt A sheet](../../../../out/reimu_fumo_attempt_023_pattern_body/seated_support_sculpt_v1/attempt_a/front_side_3q_checkpoint.png)
  — one continuous ring cage after a joined/remeshed three-volume test.
  **Verdict:** both rejected from three cheap views; A is a snowman and B is a
  tapered reclining pear, with tangent-looking foot roots in both.
- [Goal skill process-audit update](../../../../.agents/skills/goal/SKILL.md)
  — makes the whole delivery loop, critical-path bottleneck, and measured
  feedback-time optimization part of every attempt review. **Verdict:** skill
  validation passed; quick checkpoints may reject but cannot replace the full
  acceptance evidence.
- [Coupled support-and-hair 320 px sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/two_panel_mochi_v1/context_contact_sheet.png)
  — corrected two-panel support under minimal crown, fringe, cheek-lock, and
  rear-hair context. **Verdict:** the support remains provisionally useful and
  rear coverage is complete; diagnostic hair is rejected as a hemispherical
  helmet with cardboard edges, so no full bare-head packet was made.
- [Complete seated-body Attempt B front/side checkpoint](../../../../out/reimu_fumo_attempt_023_pattern_body/seated_body_assembly_v1/attempt_b/front_side_checkpoint.png)
  — adds a low hidden support beneath revised panels. **Verdict:** rejected
  after two views; the support cannot change an analytic outer loft, so the
  cape, collar, and wedge sleeves persist. The complete parametric-surface
  family is retired before Attempt C.
- [Complete seated-body Attempt A front/side checkpoint](../../../../out/reimu_fumo_attempt_023_pattern_body/seated_body_assembly_v1/attempt_a/front_side_checkpoint.png)
  — first garment test with torso, sleeves, attached white legs, grounded feet,
  and the frozen seated-depth band. **Verdict:** front contact improved, but
  side remains a triangular cape with a collar neck and pill sleeve; rejected
  after two views before masks or topology work.
- [Two-panel mochi head checkpoint 01, three-quarter](../../../../out/reimu_fumo_attempt_022_pattern_head/two_panel_mochi_v1/checkpoint_01_three_quarter.png),
  [front](../../../../out/reimu_fumo_attempt_022_pattern_head/two_panel_mochi_v1/checkpoint_01_front.png),
  [side](../../../../out/reimu_fumo_attempt_022_pattern_head/two_panel_mochi_v1/checkpoint_01_side.png),
  and [top](../../../../out/reimu_fumo_attempt_022_pattern_head/two_panel_mochi_v1/checkpoint_01_top.png)
  — first hidden-support reset at `H/W=.890` and `D/W=.775`, with two panels
  and no broad gusset. **Verdict:** material macro-form improvement, but only a
  provisional parent; the seam protrudes and the observable coupled hair/head
  silhouette is still unverified.
- [Retired Stage-2C checkpoint 02, three-quarter](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v3/checkpoint_02_three_quarter.png),
  [side](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v3/checkpoint_02_side.png),
  and [top](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v3/checkpoint_02_top.png)
  — numerically D-shaped broad-gusset cushion. **Verdict:** rejected before a
  full packet; planar crown/chin and a wide side band still read as a foam
  block, and its `H/W=1.048` gate conflated outer hair height with the hidden
  beige support.
- [Body red-mass Attempt F front](../../../../out/reimu_fumo_attempt_023_pattern_body/red_mass_feet_v1/attempt_f/renders/beauty/front.png),
  [side](../../../../out/reimu_fumo_attempt_023_pattern_body/red_mass_feet_v1/attempt_f/renders/beauty/side.png),
  and [three-quarter](../../../../out/reimu_fumo_attempt_023_pattern_body/red_mass_feet_v1/attempt_f/renders/beauty/three_quarter_right.png)
  — open-hem panel checkpoint after the earlier disc failure. **Verdict:**
  rejected before masks and a full packet; it is still a long cape with a
  human-like neck and detached shoe pods, so the isolated-skirt branch is
  retired in favor of a complete seated lower assembly.
- [Pinned-Blender render smoke beauty](../../../../out/fumo_review_render_packet_smoke/packet/studio.png),
  [component-ID pass](../../../../out/fumo_review_render_packet_smoke/packet/studio_ids.png),
  and [manifest](../../../../out/fumo_review_render_packet_smoke/packet/manifest.json)
  — reproducible packet generated through the repository's Blender 5.2.1 LTS
  toolchain. **Verdict:** smoke passed, including exact requested RGB interior
  pixels; final tracked implementation validation is recorded separately.
- [Known-bad body Attempt B Bazel audit](../../../../out/fumo_review_attempt_b/report/review.html)
  — linked six-view packet and component-gate table from the new hermetic
  auditor. **Verdict:** correctly failed the oversized registered head,
  undersized sleeves, short visible bodice, and nearly erased feet.
- [Stuffed-cushion v2 top-inclusive sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v2/contact_sheet.png)
  and [physical-reference comparison](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v2/reference_comparison.png)
  — structural cross-depth and seam-tension reset with a new top gate.
  **Verdict:** rejected; the top, side, and both three-quarter views still read
  as a rounded rectangular foam block.
- [Non-bald full proxy v3 contact sheet](../../../../out/reimu_fumo_non_bald_review_proxy_v3/contact_sheet.png)
  — six views of the latest full body with explicit crown, temple, and rear
  hair coverage. **Verdict:** the bald-rear regression is fixed in this review
  proxy, but the helmet-like hair and failed body remain rejected.
- [Stuffed-cushion v1 seven-view sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v1/contact_sheet.png)
  and [physical-reference comparison](../../../../out/reimu_fumo_attempt_022_pattern_head/stuffed_cushion_v1/reference_comparison.png)
  — the first stuffed child of the panel cage, explicitly shown without hair or
  face. **Verdict:** rejected; side and three-quarter views read as a rounded
  mattress with a belt-like gusset and paired lower pinches.
- [Paper-pattern head cage v1 review sheet](../../../../out/reimu_fumo_attempt_022_pattern_head/paper_cage_v1/contact_sheet.png)
  — front, side, rear, both three-quarters, silhouette, and cross-section of
  the planar front/rear/gusset construction proof. **Verdict:** bounded Stage 1
  representation test passed; the angular cage is not a plush checkpoint.
- [Separate-pattern body v5, Attempt F contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_f/contact_sheet.png)
  — latest full-body descendant after zero-crossing corrections. **Verdict:**
  rejected; megaphone sleeves, stacked contacts, flange ruffle, and puck feet
  remain despite almost every numeric gate passing.
- [Body pattern-piece reset report](../../../../out/body_pattern_reset/body_pattern_reset_after_v5c.md)
  — replaces perimeter rings and annuli with manufactured cloth panels and
  staged component coupons. **Verdict:** accepted as the next body hypothesis;
  no resulting geometry has passed yet.
- [Separate-pattern body v5, Attempt C contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_c/contact_sheet.png)
  — six fixed views after the first volumetric-contact reset. **Verdict:**
  rejected; angular sleeves and a broad white strip remain, with `806`
  skirt/ruffle and `150/166` ruffle/foot crossings.
- [Goal-skill forward-test final audit](../../../../out/goal_skill_forward_test/palette_optimizer/artifacts/turn_005_goal_record_audit.txt)
  — realistic six-stage validation of the persistent-goal workflow and session
  links. **Verdict:** passed; four nonblocking clarity findings were applied to
  the skill.
- [Separate-pattern body v5, Attempt B contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_b/contact_sheet.png)
  — six fixed views of the surface-relative body rebuild. **Verdict:**
  rejected; fin sleeves, bead hands, hard waist, flange ruffle, clipped puck
  feet, and `32/32` ruffle/foot crossings.
- [Attempt B front component-ID mask](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_b/body_v5_component_id_front.png)
  and [three-quarter component-ID mask](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_b/body_v5_component_id_three_quarter.png)
  — prove the pieces are separate but also expose detached roots and wrong
  depth order. **Verdict:** useful diagnostic evidence, not a visual pass.
- [Separate-pattern body v5, Attempt B Blender file](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_b/reimu_fumo_body_v5_attempt_b.blend)
  — exact disposable source reviewed above. **Verdict:** preserved rejection;
  the next body proof must replace its component representations.
- [Head and hair representation-reset report](../../../../out/head_pattern_reset/head_hair_representation_reset.md)
  — proposes sewn front/rear cushion panels, a continuous gusset, and open hair
  pattern pieces after the egg-head strategy failed. **Verdict:** accepted as
  the next modeling hypothesis; no geometry has passed yet.
- [Separate-pattern body v5, Attempt A contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_a/contact_sheet.png)
  — front, side, rear, and three-quarter review of the first body-pattern
  reset. **Verdict:** rejected; rigid stacked torso, detached sleeves, flat
  ruffle, and swallowed feet.
- [Separate-pattern body v5, Attempt A Blender file](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v5_agent/attempt_a/reimu_fumo_body_v5_attempt_a.blend)
  — disposable source for the rejected body review. **Verdict:** retained only
  as failure evidence; Attempt B also failed.

## 2026-08-29 — latest head, face, hair, and body reviews

- [Face-applique v4 contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/face_applique_v4/contact_sheet.png)
  — latest face/head deformation from Attempt 21. **Verdict:** rejected; the
  head remains egg-shaped, fringe dominates the face, and the eyes look
  faceted.
- [Face-applique v3 contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/face_applique_v3/contact_sheet.png)
  — earlier sleepy-eye layout with better graphic recognition. **Verdict:**
  eye graphic retained as reference, but the head and side silhouette fail.
- [Open-hair v5 contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/open_hair_v5/contact_sheet.png)
  — open fringe and cheek-lock construction. **Verdict:** rejected; too sharp,
  symmetric, schematic, and still attached to a bulbous cushion.
- [Rear-curtain v5 contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/rear_curtain_v5/contact_sheet.png)
  — technically open rear hair panel. **Verdict:** rejected; reads as a shield,
  cape, and wraparound mane instead of hanging plush hair.
- [Shaped sewn-head neutral contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/shaped_toile_v1/neutral_contact_sheet.png)
  — exact connected-panel head toile before hair and face. **Verdict:** topology
  is reusable evidence, but the rendered cushion is generic and egg-like.
- [Continuous-shell body v4 contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/seated_body_v4/contact_sheet_latest.png)
  — final child of the previous body representation. **Verdict:** rejected;
  lampshade skirt, paper sleeves, and extensive foot/garment overlap.
- [Five-piece bow proof contact sheet](../../../../out/reimu_fumo_attempt_021_base_cage/bow_agent_v1/contact_sheet_v2.png)
  — isolated bow construction proof. **Verdict:** rejected; petal/paddle forms
  intersect the head and do not resemble gathered fabric.

## Reference and measurement artifacts

- [Physical-reference target board](../../../../out/reimu_fumo_attempt_016_owned_cushion/guides/reference_target_board.png)
  — normalized visual target used for head silhouette calibration.
- [Head reference board](../../../../out/attempt010_reference_audit/boards/head_references.png)
  — supplied front, side, and three-quarter evidence collected for comparison.
- [Bow reference board](../../../../out/attempt010_reference_audit/boards/bow_references.png)
  — supplied bow views and overlap evidence.
- [Physical-front body crop](../../../../out/body_reference_physical_front_4x.png)
  and [physical-side body crop](../../../../out/body_reference_physical_side_3x.png)
  — enlarged supplied references used for the seated-body measurements.
