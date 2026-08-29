# A83 evidence after C9: C10--C12 and animation scaffolding

Prepared from the current files without publishing canonical goal state.

## Canonical binding at preparation time

- Goal: `reimu-fumo`, generation `1`, lifecycle generation `1`, criteria
  revision `1`, Goal resourceVersion `36`.
- Active attempt: `a83-incremental-sculpt`, resourceVersion `4`, state
  `open`.
- Outcome / execution: `open / active`.
- Criteria digest:
  `sha256:c5522700389e76975e7978515c586433ca2058a6d5012ef45fbbadcb78a5740c`.
- The canonical A83 result currently ends at C9. This packet is prepared as
  evidence for a later coordinator-owned checkpoint; it is not itself a
  publication operation.

## Retained appearance state

The retained appearance checkpoint remains exact C1b:

- C1b blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/live_author/a83_C1b_coupled_cap_receiver_narrow.blend`,
  `sha256:d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`;
- Fixed front:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C1b_packet/front.png`,
  `sha256:8c4ad6aedce38e26cda48c56b41ebeeac838e9354389276b41c5083ee0774c7b`;
- Fixed three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C1b_packet/three_quarter.png`,
  `sha256:0e4ff48da73a24bbddcc1d3546e950c61dae22b88d82fa50c4b99bf6fbc2ab5c`.

C1b is an internal checkpoint, not an approved model. It still fails the
absolute reference-likeness and constructed-plush gates: the head/hair reads
as a deep smooth helmet, the bow as rigid slabs, the sleeves as cones, and the
skirt/feet/contact stack as hard disconnected forms.

The tracked reusable asset remains unchanged:

- Tracked `reimu_fumo.blend`:
  `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`,
  `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.

No C10--C12 or rig artifact was promoted into it.

## C10: upper-crown correction

### Native Grab path: no candidate

The bounded native-sculpt run stopped at its first effect gate. Two screen-ray
targeting attempts aborted before mutation; the third reached a verified
receiver surface point, but the first Grab produced zero coordinate change and
zero moved vertices. Right-side, cap, render, and save operations were
correctly skipped. There is no candidate identity or appearance evidence.

- Native-sculpt result:
  `out/reimu_fumo_attempt_083_incremental_sculpt/native_sculpt_c10_crown/RESULT.md`,
  `sha256:17095e7e10eed0c907b43ba87373342a1080d47a8e77e50890abde1fae0507c2`.
- Verdict: **NO CANDIDATE**.

This is an interface failure, not evidence that increasing brush strength or
guessing another point would improve the model.

### Reversible paired shape key: UNDO

The paired cap/receiver shape key mechanically corrected approximately half of
the measured upper-crown width error while preserving topology, the face
opening, bow seat, lower rows, and non-owner objects. Its exact candidate was:

- C10 shape-key blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/live_author/a83_C10_upper_crown_shapekey.blend`,
  `sha256:1d904d9b3db87e8227c651106942e54dc1eff49ea2b7f082a1eb0bab187eb23c`;
- Author review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C10_shapekey_packet/SELF_REVIEW.md`,
  `sha256:5dd6dc844af67dd03b3df6566c37a9ace1b38f97909c79ec283023705cf9b171`;
- Implementation-blind verdict:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C10_shapekey_packet/BLIND_VERDICT.md`,
  `sha256:a92cca2b7bfd69ca0b13108010816f0412248afdad9af6a36ab0e23fcbf95d54`;
- Standing review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/standing_review/C10_KEEP_UNDO.md`,
  `sha256:37e1d2959b9476b2eca4e91c02122540606f5aa5a0da7c348ec13bda653fc00a`.

The independent fixed-view review found no clearly visible
reference-directed improvement at normal scale and retained the same helmet,
crown/temple dent, bow, garment, sleeve, and foot failures. **UNDO C10**. The
author's initial relative preference is superseded by the blind verdict.

## C11: constructed panel with cap recess

C11 replaced one short rear lock with a closed seven-by-nine panel cage and
recessed 334 covered cap vertices. Nominal geometry was large enough, but the
fixed pixels still showed an uninterrupted helmet, no broad owned plane or
lapped edge, and a disqualifying beige receiver leak at the root.

- C11 review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/REVIEW.md`,
  `sha256:8728374db6e8307c6b5f34eade36e7c8fe99ffeec2eba1eb121a1cd8c78f0e84`;
- Non-candidate blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/C11_panel_coupon_NON_CANDIDATE.blend`,
  `sha256:f46164eb577ae3187658c1e1c168ad2b6863c60f2cd9cbf743e7ced15bbc695e`;
- Three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/renders/three_quarter.png`,
  `sha256:d97b3d9d63468bb35f6a240f08ccf2dcea1e41fc08bdc3d315e251857580b89f`;
- Rear:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/renders/rear.png`,
  `sha256:3f0064f78fc4ebdda7c055e9f955abde505d687a42fd9153b19443a2a7cc3213`;
- Manifest:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C11_panel_coupon/manifest.json`,
  `sha256:e634d2f1ce43a0a6ed0faab88af2356ecc1b387ecac1bd87aae4bc0695cc7fc7`.

Verdict: **UNDO / NON_CANDIDATE**. The cap-recess-first interface was rejected;
the broader panel representation remained unproven until C12.

## C12: visibility-first constructed panel

C12 began again from exact C1b, did not recess the cap or head, and replaced
only `A42 Short right rear lock` with a topology-compatible closed panel. A
temporary owner-ID pass first proved pixel ownership. The same geometry was
then shown in the existing brown hair material, with only its free edge lifted
in pass 2.

Measured owner-ID evidence:

- one connected `64 x 82 px`, `3,738 px` owner region;
- gate: `35 x 70 px`, `1,348 px`;
- no beige receiver leak;
- exact panel geometry fingerprint:
  `sha256:98ec2e5be110822c56a5dcefc3fe3abc8c346a9e89fefd2607618635583297b7`.

The brown render retains one continuous curved lap boundary and the rear view
shows unambiguous front/back ordering. This is a **PASS only for panel
representation feasibility and pixel ownership**. The cage remains coarse,
the lobe too smooth, and the whole hair mass too helmet-like; it is not a
candidate and promotion is not authorized.

- C12 report:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/REPORT.md`,
  `sha256:5d14db7a12f90e318c3ff61bc45a81963365a49bb35ab047385a6b938d559111`;
- Exact manifest:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/manifest.json`,
  `sha256:84d3f161e2de662cc894a9aa29933b5e1e832a5e5426efd685971bd8485eaa0f`;
- Owner-ID three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/owner_id_pass2/three_quarter.png`,
  `sha256:195f4b8f8fdfe3fa774770502afa2e8135d43ad0f0673d946d8126d1b0031384`;
- Brown three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/brown_pass2/three_quarter.png`,
  `sha256:98cf28af0906be123ef226ff40da7e4e8314568547fa9550b26b3bf20c9e7b5b`;
- Brown rear:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/brown_pass2/rear.png`,
  `sha256:b8b0c0c0de2dc84092cc717ba986ab637063c445a85342b35f7b6370a10e30a3`;
- Final non-candidate blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C12_panel_visibility/C12_brown_pass2_same_geometry_NON_CANDIDATE.blend`,
  `sha256:004ebc6d8a8aa113a577c5f8bce503145ace8a2dd3c49a630e6cac327d26d6f5`.

## Animation and reusable-structure evidence

### C1b audit and disposable probe

Exact C1b is static and fails criteria 005 and 006: it has no armature,
actions, controls, bindings, vertex groups, shape keys, or parenting. The
disposable probe proved only that the current object boundaries support coarse
rigid ownership for a bounded head turn and arm wave.

- Rig-readiness report:
  `out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/REPORT.md`,
  `sha256:ad4f4dd075931cd154cf31b1963c2b822d0cc1d9ffcef74bb10b4dfae67e7a00`;
- Rig contract:
  `out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/RIG_CONTRACT.md`,
  `sha256:4a489f7d5ee62c36c10c479dccfed8ec72320d606673963b3c5442075f2641d8`;
- Disposable probe:
  `out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/NON_CANDIDATE_C1b_rig_probe.blend`,
  `sha256:387313a9e33a5aebb8f46f9061d98f953b77b536e8aa43dff5c3d955b296b40a`;
- Probe contact sheet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/rig_readiness/NON_CANDIDATE_rig_pose_contact_sheet.png`,
  `sha256:95282da033d132c8e6c92e5b71f7859323fed6c200884e5054296d0dd39d257b`.

The probe is explicitly non-candidate and does not satisfy production soft
deformation.

### Stage R1 coarse production scaffold

Stage R1 is the first isolated structural candidate worth keeping. It preserves
C1b topology and rest appearance while adding:

- separate `Reimu_Fumo_Asset`, review, and source-archive collections;
- nine animator controls and nine matched deform bones;
- stable ownership for 116 visible reusable objects;
- 51 full-weight mesh Armature bindings and 65 preserved-world curve bone
  bindings; and
- exact persistent actions `Fumo_Seated` 1--24, `Fumo_HeadTurn` 1--48,
  `Fumo_Wave` 1--64, and `Fumo_Validation` 1--120.

Clean-reopen verification reports no review leakage, missing resource,
binding failure, control failure, material change, or source geometry change.
The isolated R1 candidate therefore supplies a **candidate-bound PASS for the
criterion-005 structural requirements**. It is not promoted into the tracked
asset and is not final goal acceptance.

- Stage R1 report:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/REPORT.md`,
  `sha256:8fa8dac8105110491244b873660e259ef77b3b131f6354a8f04b7caf4942c610`;
- Stage R1 blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/reimu_fumo_rig_stage1.blend`,
  `sha256:af1f16e567dbc8a974f5e7f0cc9c05558b66dbfa4af26d20a99f9e73f6871533`;
- Pose contact sheet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/rig_stage1_pose_contact_sheet.png`,
  `sha256:66111f8320977eb0b77a6ac87998da664ad46ddc2a95f48f68b6783746e170ee`;
- Verification report:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage1/verification_report.json`,
  `sha256:be81154379ec42ea8023b96318d4793069889ecbc3c287b1b7d5e0633be8ff8f`.

Criterion 006 remains **FAIL / INCOMPLETE**. Stage R1 is still rigid
single-owner deformation. It lacks distributed plush deformation, proven
sleeve/skirt/head/foot compression, hair and bow tip deformation, every-frame
contact checks, and clean append/replay validation.

The exact acceptance contract and minimum deformation architecture are
preserved separately:

- Animation acceptance:
  `out/reimu_fumo_attempt_083_incremental_sculpt/animation_acceptance/ANIMATION_ACCEPTANCE.md`,
  `sha256:1791bc478081c13922471971c71da2cd6a7fdb7ca57ade3d331a58c0f094dc4f`;
- Current animation audit:
  `out/reimu_fumo_attempt_083_incremental_sculpt/animation_acceptance/CURRENT_AUDIT.md`,
  `sha256:4db33f3d0e6c301b88249ebdf9921706afbc9925da8e9a62d6ae10b1164e527c`;
- Deformation design:
  `out/reimu_fumo_attempt_083_incremental_sculpt/rig_deformation_design/DEFORMATION_DESIGN.md`,
  `sha256:a060792b2e90b5cb3b3538f61b99216f91597a9bfb2af52d9aeff66e34624585`.

The design is guidance, not acceptance evidence. Its next mechanics gate is a
small two-zone sleeve cage/lattice coupon that demonstrates visible non-rigid
deformation while keeping panel trims registered and avoiding clipping.

## Criterion status at this checkpoint

| Criterion | Status | Exact conclusion |
| --- | --- | --- |
| 001 reference likeness | **FAIL** | C1b remains recognizable but materially unlike the canonical plush; C10/C11 were undone and C12 is only a representation probe. |
| 002 measured silhouette | **UNVERIFIED** | C10 measured only three upper-crown rows and was undone; no frozen complete landmark packet passes. |
| 003 plush construction | **FAIL** | C12 proves a panel can own pixels, but the retained whole model still reads as helmet/cones/rigid stacked forms. |
| 004 presentation quality | **FAIL** | The blind C10 review scores major categories below 8/10 and reports major visible defects. |
| 005 reusable structure | **PASS for exact isolated R1 candidate; not final/promotion evidence** | R1 has the named hierarchy, resources, armature, controls, actions, ownership, and review split. The tracked deliverable is unchanged. |
| 006 animation readiness | **FAIL / INCOMPLETE** | R1 demonstrates conservative rigid owner motion only; soft deformation, every-frame contact, and append/replay are absent. |
| 007 technical integrity | **UNVERIFIED for a final asset** | R1 passes its bounded structural checks, but no complete final-model mesh/modifier/dependency audit exists. |
| 008 repository delivery | **NOT SATISFIED by this checkpoint** | The tracked LFS asset and target were not updated with an accepted appearance-plus-rig result. |

## Process audit and next action

C10 confirmed that a mechanically correct local displacement can still be
visually inert. C11 showed that recessing support before proving pixel
ownership creates expensive ambiguity and can expose the receiver. C12's
owner-ID-first gate settled the representation question in two bounded passes:
the panel footprint and one free-edge interface can produce the required
visible ownership without touching the cap. Stage R1 similarly separated
structural rig feasibility from the still-rejected appearance and soft-motion
questions.

The highest-leverage next action is **C13: directly model one refined
side/back hair panel on an isolated exact C1b copy, using C12's proven pixel
footprint and free-edge ownership while discarding its coarse cage finish**.
Keep the cap and receiver unchanged, bury the root under existing upper hair,
add an asymmetric lower contour and a softly stuffed surface, and render fixed
front, three-quarter, and rear immediately. Preserve only a candidate that
keeps the front regression stable, has no receiver leak, and is clearly
preferred against the canonical turn and physical side reference at normal
review scale.

Do not generalize R1's rigid bindings yet. After one appearance owner survives,
transfer or rebuild the retained R1 scaffold on that retained appearance copy,
then run the bounded two-zone sleeve deformation coupon before expanding the
soft-deformation system.
