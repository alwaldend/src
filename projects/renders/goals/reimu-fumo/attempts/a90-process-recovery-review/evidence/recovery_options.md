# Reimu Fumo recovery options after A68-A89

## Verdict

**Revise the strategy.** Do not resume A89 parameter tuning and do not make
the A85/A87/A88 module chain the sole path forward. Freeze those three exact
survivors as a comparison branch, then run a short bakeoff against a fresh,
complete, low-detail whole-plush blockout. Use direct manual sculpting only
after a live input-support probe proves that it can move broad forms; defer
retopology and rigging until a complete neutral form wins a visual gate.

This is not a recommendation to lower the goal. The evidence shows failure of
the current sequencing and authoring representations, not that a faithful
Reimu Fumo is infeasible in Blender. Stopping or weakening acceptance is a
fallback only after the bounded recovery test below fails.

## Scope and evidence

This review covers the durable attempt record from A68 through closed A89,
the exact surviving and protected assets, the current renders, and every
controlling reference. It does not mutate or promote a Blender asset.

### Outcome and constraints

The required result remains a reusable, animatable, 25 cm Reimu Fumo that:

- is recognizably the canonical supplied variant in front, side, rear, and
  three-quarter views;
- reads as sewn, lightly stuffed fabric rather than anatomy, plastic, armor,
  a helmet, a cone, or assembled primitives;
- preserves the canonical broad head, compact seated body, layered hair,
  asymmetric bow, face, sleeves, gathered skirt, and small feet;
- eventually passes every fixed landmark and absolute image-quality gate;
- supports a clean armature and deformation tests; and
- remains separate from review-only cameras, lights, and references.

No full-goal criterion currently passes. The recent `refine` outcomes are
provisional module selections, not sculpt approval or whole-asset progress.

### Reference authority used

| View or property | Controlling source | Supporting veto sources |
| --- | --- | --- |
| Exact variant, scale, front silhouette, graphics, seated proportions | `projects/renders/blender/fumo/reimu_fumo/references/canonical_front_25cm.png` | `clean_front.png`, `physical_front.png` |
| Side, rear, three-quarter silhouettes, depth, layer order, seated volume | `projects/renders/blender/fumo/reimu_fumo/references/canonical_turn_180.gif` | `physical_side.png`, `turn.gif` |
| Fabric pile, panel thickness, soft contact, applique seating | physical front/side and `sofa.gif` | canonical pair remains dimension authority |

The references are not averaged into a fictional hybrid. The canonical front
wins direct conflicts about front identity and 25 cm scale; the canonical
turn wins hidden-side construction and depth.

### Protected evidence and exact parents

| Artifact | Status | Identity |
| --- | --- | --- |
| Tracked standalone `reimu_fumo.blend` | Rejected migrated baseline; unchanged | `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76` |
| Rung 003 complete model | Frozen, unaccepted comparator; recognizable but visibly wrong | `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b` |
| A85 S12 head receiver | Provisional receiver only | `982da6404ea6edcbb4432903e67dad4ee5c130a203a5a5727a374b773fc9ad8a` |
| A87 S04 front-hair frame | Provisional three-piece module only | `da300595b73185d7c4deece12f7e98e7d6889e8f738c94709c828bd601e58d2a` |
| A88 S07 crown | Provisional two-panel crown only; active immutable parent | `0a30e2af3142081648bb3137ad75d6d1cc73de55e9f830f85bad1f85e92c8788` |
| A89 S06 rear-hair diagnostic | Rejected; do not use as a parent | `960fb7b1bc0805ed7f7d8eb5817c40478e25fa57e53abb7028730fa923996845` |

The A85/A87/A88 files should remain byte-exact. Any recovery test using them
must begin from a disposable copy.

## What A68-A89 actually established

### Attempt-phase evidence

| Attempts | Strategy tested | Decision-bearing result |
| --- | --- | --- |
| A68-A69 | Native sculpt and coupled deformation of the old cap/head interface | Broad strokes produced dents or attachment loss while the monolithic helmet survived; representation closed. |
| A70-A73 | Constructed receiver panels and deterministic analytic cages | Mechanically valid forms remained masks, mattresses, rounded boxes, or eggs; receiver-first analytic lofts closed. |
| A74-A78 | Visible hair envelopes, isolated face carrier, live-sculpt interface, sparse paired surfaces, and head-rest preflight | Cards, canopy, exposed support, and insufficient pixel ownership recurred; A78 correctly stopped before another redundant build. |
| A79-A82 | Complete diagnostic hair forms and fast direct sparse blockouts | First-pixel latency improved dramatically, but every visible representation remained bald, faceted, cuboid, or helmet-like. |
| A83 | Reversible whole-model sculpt and local coupon ladder | The complete parent was declared categorically nonviable because head/hair/bow and seated body/garment failures were coupled under two macro owners. |
| A84-A86 | Fresh head-only scene, sparse head cage, and collar integration | A85 retained a provisional shallow receiver; the rounded-cube path and collar-only coverage were rejected. |
| A87-A89 | Front hair, crown, then lower-rear hair modules | A87/A88 retained local 7/10 modules; A89 failed all states as cape/fin construction and added no survivor. |

The main visual evidence anchors are:

- all-reference modeling sheet:
  `out/fumo_reference_sheet/fumo_modeling_contact_sheet.png`;
- best complete but rejected Rung 003 comparator:
  `out/reimu_fumo_working_ladder/rung_003_eyes_locks_sleeves/five_view_gate/contact_sheet.png`;
- A85 receiver six-view sheet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C19_process_reset/a85_review/A85_S12_six_view.png`;
- A88 provisional crown six-view packet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C22_crown_interface/packets/s07_six/`;
- A89 rejected final diagnostic packet:
  `out/reimu_fumo_attempt_083_incremental_sculpt/C23_lower_rear_hair/packets/s06_six/`.

### Measurable progress

1. **The feedback path became fast enough.** A79 required 94 minutes 30
   seconds to reach the first new whole-character pixels. A80 reduced this to
   50 minutes 42 seconds. A82 froze a simple blockout at `T+1m35s`, delivered
   its controlling pair at `T+3m35s`, and recorded a decision at `T+6m12s`.
   A87 reports roughly 18 seconds per pinned three-view packet. Rendering and
   compute are no longer the main bottleneck.
2. **The review discipline improved.** Immutable copies, fixed cameras,
   implementation-blind image review, component-ID passes, exact hashes, and
   undo/reset decisions now prevent rejected geometry from contaminating the
   tracked asset.
3. **A85 found a useful hidden receiver.** Its shallow cushion stays inside
   the measured envelope without an egg or human jaw. This is real interface
   evidence, but it remains box-like and has a visible underside cavity.
4. **A87 and A88 found locally coherent panel ownership.** The short fringe,
   cheek locks, and closed crown score about 7/10 at their bounded module
   stages and are more plausible than the older helmet-plus-cards family.

### Why visible progress has nevertheless stalled

- A68-A84 contain 17 consecutive resets. A86 is another reset. A85, A87, and
  A88 preserve only provisional modules; A89 preserves no new geometry. Across
  22 attempts there is no accepted full sculpt and no accepted goal criterion.
- A83 already showed that the whole character is governed by two coupled
  macro owners: head/hair/bow and seated torso/garment/contact. The response
  was to spend A84-A89 almost entirely on the first owner, then on increasingly
  narrow sub-interfaces inside it. The current artifact is therefore a blank
  head with partial hair, not a whole Reimu Fumo whose likeness can improve.
- A85's measured envelope and A87/A88's local 7/10 reviews are not additive.
  The present combined pixels still read as a broad cap over a pale rounded
  block; rear coverage is absent. No evidence proves that these modules form a
  viable complete head, much less a complete plush.
- A89 repeats the historic failure at a smaller scale. S05/S06 close literal
  holes but become a rigid cape/shield with vertical side fins, reverse the
  intended crown overlap, and leave a median `0.148 Wh` rear receiver band.
  Seven authored states produced no survivor because the first blind review
  arrived after the representation had already converged on the wrong cue.
- Rung 003 is globally wrong but still demonstrates an important contrast: a
  complete low-quality model can be judged as Reimu, while the cleaner recent
  modules cannot yet be judged as the requested subject at all. Local quality
  has been optimized without a subject-level acceptance-visible baseline.

The dominant bottleneck is now **visual decomposition and sequencing**, not
Blender launch, renderer throughput, agent count, validation machinery, or a
missing research paper.

## Adversarial comparison of recovery strategies

### Summary

| Strategy | Expected whole-image leverage | Time to first decision | Main risk | Reuse / rig consequence | Verdict |
| --- | --- | --- | --- | --- | --- |
| Continue A85/A87/A88 modular assembly | Low to medium until many more modules exist | 15-30 min for another rear packet; likely hours before a whole-subject decision | Local gates continue to pass while the assembled subject remains generic; rear hair repeats cape/fin failures | Good object ownership if the forms survive; many panel seams can complicate deformation | One bounded control branch only |
| Restart with direct manual sculpt, then retopo | Potentially high, but highly uncertain with current input path | 10-15 min for an input-support probe; 45-90 min for a meaningful macro packet | A68/A71 showed local dents rather than broad control; A76 found no real pointer events through MCP; later retopo can become expensive before likeness is proven | Best organic ceiling after retopo; raw sculpt alone is not reusable or animation-ready | Conditional secondary route |
| Fresh complete low-detail whole-plush blockout | High and immediate: every identity and proportion relationship becomes visible together | 20-45 min for front/three-quarter; about 60-90 min for a six-view survivor | Can repeat A37-A43's disconnected primitives if “low detail” is mistaken for balls, cones, and boxes | Strong: sparse named owners make later panel construction and rig planning easier | Primary route |
| Stop or reframe acceptance | No visual leverage | Immediate | Converts process failure into a weaker deliverable and abandons an apparently feasible goal | No reusable approved asset | Reject now; retain only as bounded fallback |

The timing estimates are decision budgets, not completion promises. They are
grounded in the observed 18-30 second render packets and A82's 3m35s complete
pair. Geometry authoring and visual judgment dominate them.

### Option 1: continue the A85/A87/A88 assembly

**Strongest case for it.** This is the only recent lineage with exact survivors
rather than only rejections. S12 passes the bare receiver silhouette bands,
S04 establishes separate fringe and lock ownership, and S07 closes the crown
without moving protected front pixels. The pieces are already isolated and
could become reusable rig parts.

**Strongest objection.** Each review explicitly withholds standalone approval.
Current front and three-quarter pixels still show a generic cap on a blank
rounded block, and the direct rear shows almost the entire receiver. The next
missing owner immediately failed as a cape. Continuing one module at a time
can consume many successful-looking local cycles before exposing that the
assembled character remains wrong.

**Falsification gate.** Give this lineage one whole-context salvage test, not
another rear-hair refinement. On an immutable A88 copy, add only coarse,
reference-normalized face, bow, seated torso/skirt, sleeves, and feet proxies.
Do not alter S12/S04/S07. Render fixed front, worse three-quarter, profile, and
rear. Reject the lineage as the primary parent if:

- unlabeled same-variant recognition is not credible at blockout level;
- the head/hair still reads as a cap, helmet, block, or human hairstyle;
- any required view has exposed receiver regions that require another broad
  concealing shell; or
- it does not beat the fresh whole-blockout branch in a blind A/B on overall
  likeness, silhouette, and construction.

This test preserves the real module work while preventing sunk-cost privilege.

### Option 2: restart with direct manual sculpt and retopology

**Strongest case for it.** The plush has broad, soft, asymmetric forms that are
difficult to infer from analytic lofts. A genuinely manual sculpt can adjust
coupled volumes by eye and then be retopologized into fabric panels and rig
loops. This has the highest artistic ceiling if the operator has responsive
viewport input.

**Strongest objection.** The current agent interface has not demonstrated that
modality. A68 and A71 produced small dents from nominally broad Grab strokes;
A76 found that the live MCP exposed a viewport but no real pointer events.
Calling scripted BMesh or synthetic operators “manual sculpt” does not change
that evidence. Starting a high-resolution sculpt without a passed complete
blockout also risks an expensive smooth blob, followed by premature retopo.

**Falsification gate.** Before authoring a candidate, make a disposable input
coupon and require three broad, independently directed deformations to move at
least roughly 20% of the intended silhouette in the expected direction without
localized dimples, zero-effect output, or a camera-only illusion. Freeze and
render the probe. If that fails within 10-15 minutes, close this route rather
than substituting more scripted generators. If it passes, sculpt only a
complete low-resolution macro character first; no retopology starts until its
whole-subject packet wins the visual bakeoff.

### Option 3: build a simpler complete low-detail blockout first

**Strongest case for it.** This directly restores the missing acceptance-visible
unit: the whole plush. It exposes head-to-body ratio, bow span, seated height,
skirt/foot occlusion, side depth, and rear mass in one packet. A sparse complete
blockout is cheap to discard, easy to compare, and compatible with later
retopo and rig controls. It uses the lesson from A82's fast blockout loop while
avoiding the A84-A89 head-only tunnel.

**Strongest objection.** Attempts 37-43 already produced complete models made
from disconnected primitives, helmet hair, cone skirts, and ball feet. A
coarser model can look even worse if it is allowed to pass on silhouette alone.

**Falsification gate.** The blockout must be low detail but construction-aware:
one shallow cushion head, one compact seated torso/hip core, thin but volumetric
hair/bow/sleeve/skirt owners, and small stuffed foot pods with explicit overlap.
Use simple materials or owner IDs only to separate forms. Reject at the first
front/three-quarter packet if it shows anatomy, spheres/eggs, a dress cone,
tube limbs, or disconnected floating parts. Promote to six views only when:

- an implementation-blind reviewer recognizes the canonical Reimu variant at
  blockout level;
- overall likeness, silhouette, and identity are each at least 6/10 as an
  internal parent-selection threshold, never as final approval;
- the main front landmarks are within `0.08 Wh` provisionally and depth lies
  within the frozen broad bands; and
- no single major wrong-category defect makes refinement implausible.

This is the recommended primary strategy.

### Option 4: stop or reframe acceptance

**Strongest case for it.** The history is unusually long, the current tool has
limited direct sculpt control, and the final 8/10 multi-view plus animation
gate is demanding. A stylized rather than reference-faithful model would be
cheaper.

**Strongest objection.** The user did not ask for a weaker target, the canonical
references are adequate, Blender can represent the object, and the process has
not yet run the most direct whole-plush bakeoff under the now-fast loop. Lowering
acceptance would hide a sequencing failure rather than solve it.

**Falsification gate.** Reopen this option only if three genuinely distinct,
complete macro representations fail the same subject-level category under the
same calibrated views, or if a required authoring modality is demonstrably
unavailable and no remaining route can meet the requested quality within an
explicit user time budget. Any acceptance change then requires a new user
decision; it must never happen silently.

## Recommended bounded four-attempt recovery plan

These are attempt-sized decision units, not four promises to finish the asset.
Attempts 1 and 2 are independent and should run concurrently on separate
copies. One coordinator owns the comparison and any later canonical parent.

### Recovery attempt 1 — complete construction-aware blockout

- Start an empty task-owned scene with the frozen cameras, reference contract,
  and 25 cm scale; do not append rejected visible geometry.
- Author the entire neutral plush at low resolution: head, visible hair mass,
  bow mass, compact seated torso, skirt/hem mass, sleeves, and foot pods.
- Use at most one or two sparse owners per major mass; add no ruffles, stitches,
  fibers, detailed face applique, materials, armature, or simulation.
- Freeze the first front/three-quarter packet within 30 minutes of geometry
  start. Permit at most one bounded correction before the six-view decision.
- **Keep only if** it clears the Option 3 recognition, score, landmark, and
  wrong-category gates. Otherwise reset the representation, not its details.

### Recovery attempt 2 — A88 whole-context salvage control

- Copy exact A88 S07 and keep A85/S12, A87/S04, and A88/S07 byte-equivalent.
- Add disposable coarse identity/body proxies sufficient to judge the whole
  character; do not continue A89 and do not repair current modules.
- Render the same cameras, framing, lighting, and owner-ID scheme as Attempt 1.
- Run an implementation-blind randomized A/B only after both absolute reviews.
- **Keep the modular lineage only if** it clearly beats Attempt 1 on whole-
  subject likeness, silhouette, and construction with no exposed receiver or
  helmet/cap major failure. A tie goes to the simpler complete blockout because
  it has fewer inherited interfaces and lower sunk-cost risk.

### Recovery attempt 3 — one macro refinement or proven manual-sculpt test

- If either branch survives, freeze it as an explicitly unaccepted baseline
  and diagnose its single largest whole-image defect. Make one coupled macro
  correction owned by the smallest set that can actually move those pixels.
- If neither branch survives, run the direct-input support coupon described in
  Option 2. Only a passed coupon authorizes one complete manual low-resolution
  sculpt; a failed coupon closes direct sculpt without another candidate.
- Render controlling and regression-risk views first, aligned to references.
- **Keep only if** the named defect visibly improves, protected macro landmarks
  stay within their bands, and the internal whole-subject scores rise toward
  7/10 without a new major failure. Two misses in the same category trigger the
  stop/reframe review rather than a fifth representation family.

### Recovery attempt 4 — constructed form and animation-readiness coupon

- Run only if Attempt 3 leaves a viable complete neutral parent.
- Convert one identity-critical region and one motion-critical interface into
  clean manufactured topology: recommended owners are the complete hair/head
  interface and one sleeve-to-torso or hip-to-skirt interface.
- Preserve the winning outer pixels while adding real panel thickness, roots,
  overlap, and deformation-friendly edge flow. Add a minimal armature coupon
  for head turn plus one arm motion; do not texture the whole character yet.
- **Pass the recovery phase only if** the constructed/retopologized copy keeps
  the accepted macro silhouette within `0.03 Wh`, has no clipping or floating
  surfaces in required views, and the two test motions do not tear the selected
  interfaces. Otherwise undo the retopo/rig delta and revise its local owner;
  do not rebuild the visual parent.

After this four-attempt phase, a viable neutral whole-plush parent can proceed
through full constructed sculpt, then the explicit material, seated-pose, and
final rig approval gates. If no viable parent remains, pause and present the
evidence for an explicit acceptance/scope decision.

## Process changes that make the plan materially different

1. **Whole-subject pixels are the primary progress unit.** A module may be
   developed separately, but every provisional keep decision must be followed
   immediately by a same-turn whole-context packet. Local 7/10 scores cannot
   accumulate into an assumed full-model pass.
2. **Two lanes, one comparison.** The fresh blockout and modular salvage branch
   run on isolated copies in parallel. Parallel agents may author, measure, and
   review, but one coordinator chooses the parent and nobody mutates the
   tracked asset.
3. **Beauty review precedes mechanics.** Produce the first 512 px controlling
   pixels before topology inventories, reopen matrices, promotion scaffolding,
   or full packets. Component IDs diagnose ownership only after a beauty defect
   is visible.
4. **One correction per survivor.** Do not author S00-S06 before independent
   review. A first packet gets an immediate blind verdict; a second state is
   allowed only for one named, causally reachable defect.
5. **Retopo and rig follow likeness.** Preserve low-resolution object and
   interface ownership for reuse, but do not pay for full topology, weights,
   actions, or technical export on rejected pixels.
6. **Use all references with explicit authority.** Every packet includes the
   canonical front and canonical turn views plus the relevant physical
   construction witness; reviewers must not judge against one convenient
   still or an unaligned contact sheet.
7. **Keep rejected work immutable and cheap.** A85/A87/A88 remain available as
   evidence and a reusable parts bank even if the fresh whole blockout wins.
   No result depends on overwriting or discarding them.

## Evidence that would change this verdict

- A one-packet A88 whole-context salvage rendering that clearly beats the
  fresh complete blockout would justify returning to modular assembly, but
  only with whole-context gates after every module.
- A live input coupon proving broad, predictable sculpt support would raise
  direct sculpt from a conditional route to the preferred macro-refinement
  method.
- Three distinct complete blockouts failing the same calibrated identity or
  construction gate would justify stopping for a user-authorized scope or
  acceptance review.
- Faster rendering or more parallel agents alone would not change the verdict;
  current evidence already reaches pixels in seconds to minutes. The missing
  leverage is a better whole-plush hypothesis and earlier absolute judgment.
