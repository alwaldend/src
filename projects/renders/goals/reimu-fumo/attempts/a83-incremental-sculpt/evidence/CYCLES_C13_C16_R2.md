# A83 evidence: C13--C16 and Stage R2

Prepared from the current exact artifacts without publishing canonical goal
state.

## Preparation binding

- Goal `reimu-fumo`: resourceVersion `37`, generation `1`, lifecycle
  generation `1`, criteria revision `1`, outcome/execution `open / active`.
- Active attempt `a83-incremental-sculpt`: resourceVersion `5`, state `open`.
- Criteria digest:
  `sha256:c5522700389e76975e7978515c586433ca2058a6d5012ef45fbbadcb78a5740c`.
- The canonical A83 result currently ends at C12 and Stage R1. This packet is
  scratch input for a coordinator-owned checkpoint, not a canonical mutation.

## Protected and retained state

The retained appearance checkpoint remains exact C1b:

- `out/reimu_fumo_attempt_083_incremental_sculpt/live_author/a83_C1b_coupled_cap_receiver_narrow.blend`
- `sha256:d2357588b42b18285f31fcf780f2be5e76111a002a25b9ac25cd569be6cbf8d1`

C1b remains an internal comparison and recovery point, not an approved model.
No C13--C16 appearance artifact survived review.

The tracked reusable asset remains unchanged:

- `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend`
- `sha256:489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

Neither Stage R2 nor any appearance probe was promoted into it.

## C13: coupled panel/support

C13 reproduced C12's visibility-proven broad panel exactly and added only a
closed 28-vertex brown underlap ribbon beneath its free edge. The cap and head
fingerprints remained unchanged and all fixed views retained coverage, but the
new support produced only a faint crease. It did not create independent filled
depth, low rear ownership, or a material reduction of the continuous helmet.

Verdict: **categorical UNDO / NON_CANDIDATE**. The narrow support owner was
pixel-inert for `D-HEAD-HELMET`; no placement, thickness, material, or receiver
rescue was justified.

- Review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C13_coupled_panel_support/REVIEW.md`,
  `sha256:ac3f877d1b6dee74c7deb22300c3fbc284b9fd8d800c6802fd32795d45f0f025`
- Diagnostic blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C13_coupled_panel_support/C13_coupled_panel_support_NON_CANDIDATE.blend`,
  `sha256:f897bfdd8892a7ea2968332914041be9c4b4b20a2587d9be220aff73d6607af3`
- Contact sheet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C13_coupled_panel_support/contact_sheet.png`,
  `sha256:901edeec8d5fab5455c5fe0cbda19a39d6ca44fc2ceefff15463c09f97710b10`

## C14: replacement sleeve appearance module

C14 replaced only the camera-right sleeve with an axial superellipse loft.
It removed the old hollow tunnel and passed bounded BVH contact checks, but its
fixed pixels read as a narrow downward banana/croissant with a nub-like insert,
not the canonical short broad bell sleeve made from softly stuffed cloth
panels. Collision cleanliness did not compensate for the wrong silhouette and
construction language.

Verdict: **categorical UNDO / NON_CANDIDATE**. More loft-parameter tuning would
preserve the tube representation. A future appearance replacement would need
separately constructed broad front/rear panels, a shallow stuffed gap, and a
wide folded cuff; C14 itself must not be refined or promoted.

- Review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C14_sleeve_module/REVIEW.md`,
  `sha256:cb58e6a230882e4432a730d3d18c46eef6ce7773ed980958ea1c1ed5b4d84ee0`
- Diagnostic blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C14_sleeve_module/C14_camera_right_stuffed_sleeve.blend`,
  `sha256:f383156ae0d40c2bcab5c534f0e068071dd30b966c995fc661cbb6505852bda3`
- Contact sheet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C14_sleeve_module/contact_sheet.png`,
  `sha256:be0735d2bca692600c43263a8c10021eb9e875c62cef649afb5bd823d8f4f80e`

This appearance rejection is independent of Stage R2's later deformation
interface result: R2 does not validate the retained sleeve's visual fidelity.

## C15: live cap-only whole-field deformation

C15 used one reversible shape key on a task-owned exact C1b copy. It moved
8,434 cap vertices by up to `4.93 mm`, narrowing the upper/middle field,
widening the lower field, and compressing the rear highlight. The front stayed
coherent, but the first fixed three-quarter exposed a large beige receiver
region at the crown/bow-root/side and the lower field still read as a helmet.

Verdict: **categorical UNDO / NON_CANDIDATE**. This closes cap-only whole-field
deformation: a change large enough to affect the global field loses coverage
before it creates independent fabric layers. Do not tune the shape key, shrink
the receiver, or patch the gap.

The initial artifact came from the live Flatpak 5.1 MCP host, but the exact
candidate was clean-reopened and the categorical failure was independently
reproduced with repository-pinned Blender `5.2.1`. The deliverable judgment is
therefore not based only on the live host.

- Review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/REVIEW.md`,
  `sha256:bf6a98ef60d44e2cf2dd5225064366fa48c59faf69f9be8d9d37bf1a214423ee`
- Rejected blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/C15_v1_whole_field.blend`,
  `sha256:329066759d035a317e9c66e35c6406946d455733980323979bffe37e2bc376ac`
- Pinned reproduction report:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/pinned_5_2_1/verification.json`,
  `sha256:b350c6d426229e8a8c2311bd1ecef498cc6524ad950ad21d38ac7f6d87224203`
- Pinned front:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/pinned_5_2_1/front.png`,
  `sha256:5644a082f134d11ce69df538411fb8e138b2e79a669f360fd292bad05efbf337`
- Pinned three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C15_live_multipanel/pinned_5_2_1/three_quarter.png`,
  `sha256:deb41b18c488fb8b95abb57eb50286e8ab16cbce56d777273321c760387e81f9`

## Representation review after C10--C15

The final process review closed both rescue directions for the old visible
cap:

1. leave it substantially intact and a local owner remains subordinate to the
   global helmet (C12/C13); or
2. deform it enough to change the global field and coverage fails before the
   lower round mass becomes layered (C15).

It rejected another local panel and another whole-cap deformation. It allowed
one disposable, directly authored three-role boundary test—open crown/temple
field, rear/nape pocket, and foreground return leaf—with owner-ID coverage as
the first hard gate. This was a representation experiment, not approval of a
literal three-panel topology or a canonical change.

- Process/representation verdict:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C13_process_review/NEXT_HAIR_VERDICT.md`,
  `sha256:72693a7a225ba8df31a9fb897df65f115f3872adc3470e7acdd00627e0bf42d7`

C16 executed that one allowed representation test.

## C16: three-owner boundary

### v0 owner-interface fact

C16 v0 replaced the six-object failed hair boundary with an open crown/temple
field, a rear/nape pocket, and C12's foreground return leaf. Its owner-ID front
and three-quarter pair established one useful fact: **those three owner roles
can jointly cover the head in the two first-gate views with no beige receiver
leak**. Identical v0 brown geometry retained coverage, but still read as the
same rounded helmet/balloon.

This is an interface-coverage fact only. It does not pass likeness, plush
construction, or appearance acceptance.

- v0 owner-ID blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/C16_v0_owner_id.blend`,
  `sha256:96c43f56be47ccf4c26b5333569dc51fb99b551754bb664034a1c0757d9604ac`
- Owner-ID front:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v0_owner_id/front.png`,
  `sha256:76a96becb4da93e526324c8569dd2b252599d3af2b4e41096d4da0eb259e9ca7`
- Owner-ID three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v0_owner_id/three_quarter.png`,
  `sha256:7867ee2fb11a32e11f582467883deed987ce1c9b1c572bd32e299ec160091269`
- v0 brown blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/C16_v0_brown.blend`,
  `sha256:2a8a1ddbb7b78e244c377136ff91436becd9cf66b7f611c08343fb181c402dea`

### v1 direct pass

One direct owner-level pass compacted the crown wrap, lowered the rear pocket,
and advanced the return leaf. The front remained coherent, but the fixed
three-quarter exposed a large beige receiver strip from the bow root down the
side. The return leaf became a hard vertical horn/card and the owners no longer
formed plausible stuffed coverage.

Verdict: **categorical UNDO / NON_CANDIDATE**. Per the frozen stop rule, close
this three-owner split; do not add a second placement pass, support strip, or
material rescue. Preserve only the v0 coverage fact.

- Review:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/REVIEW.md`,
  `sha256:39aba111eaa04b217bd00328692ab7e108d2c797d5357dbdea0a21644d61b91b`
- Rejected v1 blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/C16_v1_direct_layered.blend`,
  `sha256:0cd524df660983cf82d2f1c0ca9b53f42d0b2beea0c572204f9d3cbc6c5eb9bc`
- Rejected front:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v1_direct/front.png`,
  `sha256:67ac5677137b821c1be8ae7f55718b253fb7447d291f467958f3479ec921f569`
- Rejected three-quarter:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C16_live_layered_boundary/v1_direct/three_quarter.png`,
  `sha256:da460ba790833dc659ec03683d78d86f3e94888c96d62efd914d5f878e92c0e5`

The C13 process review's one allowed three-role representation test has now
been consumed and rejected. Another hair mutation requires a new
representation-level diagnosis; C16 is not a base for refinement.

## Stage R2: bounded soft right-sleeve deformation scaffold

Stage R2 starts from exact retained Stage R1 and adds a `5 x 3 x 3`
owner-local lattice under `DEF_arm.R`. It binds all 22 right-sleeve owners—3
meshes and 19 curves—while leaving the shoulder root fixed and producing
`1.979474 mm` cuff p90 displacement. Zero-strength evaluated geometry is exact,
curve/panel registration drift is only `0.016922 mm`, and the four required R1
actions, original geometry, materials, resources, and collection separation
remain unchanged. Repository-pinned Blender `5.2.1` clean-reopen verification
and fixed-view visual vetoes passed.

Verdict: **bounded KEEP as isolated right-sleeve soft-deformation scaffolding**.
It validates one shared owner-local cage interface; it is not the finished rig,
does not validate the current sleeve's appearance, and is not promoted.

- Decision:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/DECISION.md`,
  `sha256:9693d13b51118ccff7835550cad10da67a71b5666b787a847d3ae9788cc30e8f`
- Report:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/REPORT.md`,
  `sha256:445ff150edb66a32a847364321e4eafb4b77cdf4b7f670ac329948afc2275e29`
- Stage R2 blend:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/reimu_fumo_rig_stage2.blend`,
  `sha256:f609794ab6e27f7b8a5e1b93325c617dd627da2e7ac05d00e2b755fc1b5b5dde`
- Pose contact sheet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/rig_stage2_pose_contact_sheet.png`,
  `sha256:96819621c5d35d326feb4399c39f6836883d31329b8430a4a98e241ebed6ecdb`
- Sleeve close-up:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/rig_stage2_sleeve_closeup.png`,
  `sha256:40954712ce5afe4ecbbd29a9b5c575e1553c31c15185e14d3cb4401a6eb55fc2`
- Verification:
  `out/reimu_fumo_attempt_083_incremental_sculpt/production_rig_stage2/verification_report.json`,
  `sha256:79b0be52888447f5957d453955e6dd06545ef10360b1ba658e687a61adf5ff65`

Criterion 006 remains **FAIL / INCOMPLETE**. R2 proves only one sleeve at a few
discriminating frames. The opposite sleeve, head, torso, legs/feet, hair, bow,
all-frame contact/collision checks, and factory-clean append/replay remain
absent.

## Criterion status

| Criterion | Status | Conclusion |
| --- | --- | --- |
| 001 reference likeness | **FAIL** | C13--C16 produced no retained appearance candidate; C1b still has the major helmet and construction failures. |
| 002 measured silhouette | **UNVERIFIED** | No surviving complete landmark packet exists. |
| 003 plush construction | **FAIL** | C14's replacement sleeve and C16's layered boundary both failed their constructed-fabric reads; R2 changes motion, not appearance quality. |
| 004 presentation quality | **FAIL** | No exact appearance candidate has two blind reviews at 8/10 with no major defect. |
| 005 reusable structure | **PASS only for isolated Stage R1/R2 scaffolding** | The structural candidate remains unpromoted and is not final goal acceptance. |
| 006 animation readiness | **FAIL / INCOMPLETE** | Only one soft sleeve coupon is proven; required owners, actions across all frames, contacts, and append/replay are incomplete. |
| 007 technical integrity | **UNVERIFIED for a final asset** | R2 passes bounded clean-reopen checks, not the complete final-model audit. |
| 008 repository delivery | **NOT SATISFIED by this checkpoint** | The tracked asset is unchanged and no accepted result was promoted. |

## Process conclusion and next smallest useful module

The appearance loop correctly stopped C13, C14, C15, and C16 at their first
categorical evidence rather than polishing failed representations. C16 also
consumed the one three-role hair-boundary experiment authorized by the C13
process review. The next appearance action is therefore a representation reset
and causal review, not another edit to C1b, C15, or C16.

The smallest ready implementation module with a proven interface is a
**left-sleeve Stage R2 cage coupon on an isolated exact R2/R1-derived copy**.
Reuse the right-sleeve root/cuff, mesh/curve registration, zero-strength,
fixed-view, and clean-reopen gates; do not bulk-generate other zones or promote
the result. This advances criterion 006 scaffolding while the appearance
critical path performs a new post-C16 representation diagnosis. It does not
make the goal achieved and must not displace the need for a reference-faithful
appearance candidate.
