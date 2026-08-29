# Reimu Fumo goal-history audit through A89

Observed: 2026-09-01. This is a read-only audit of the durable goal history
through closed attempt A89. It does not include the recovery attempt opened to
record this review.

Unless written in full, `attempts/...` paths below are relative to
`projects/renders/goals/reimu-fumo/`.

## Executive verdict

The goal has produced extensive negative knowledge, reliable reference and
render infrastructure, and three cumulative provisional head/hair modules.
It has not produced an accepted plush, an accepted untextured sculpt, or even
a demonstrated improvement to the best complete-character comparator.

The hard acceptance state is unambiguous:

- `acceptedAttemptID` and `acceptedResultDigest` are empty and the goal is
  still open (`projects/renders/goals/reimu-fumo/goal.yaml`).
- The eight required criteria remain at revision 1
  (`projects/renders/goals/reimu-fumo/criteria.yaml`).
- Across the 22 structured attempts A68--A89 there are **18 resets, 4
  refines, 0 accepts**, and **176 criterion verdicts: 0 pass, 50 fail, 126
  unverified**.
- Including the imported history snapshot as the single structured record it
  actually is gives **18 resets, 5 refines, 0 accepts** and **184 verdicts:
  0 pass, 50 fail, 134 unverified**. The snapshot's `refine` is migration
  bookkeeping, not a legacy model acceptance.
- The tracked reusable Blend is still SHA-256
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`,
  the same protected failed baseline recorded at A68 and A89. No structured
  attempt changed it
  (`attempts/a68-sculpt-coupon/result.md`,
  `attempts/a89-lower-rear-hair/result.md`).
- The only current visual survivors are scratch checkpoints: A85 S12's bare
  head receiver, A87 S04's three-piece front-hair frame, and A88 S07's
  two-panel crown. They are cumulative parts of one incomplete head, not three
  accepted assets. A89 added nothing.

The lack of visible progress is therefore real, not merely a perception caused
by sparse status reporting. Process latency improved substantially, but the
work kept optimizing incomplete, low-scoring head modules and mechanical
conditions while whole-character likeness, body construction, materials,
rigging, animation, and delivery remained untouched.

## Audit boundary and counting caveat

The history has two different record forms:

1. `attempts/imported-unversioned/` is one digest-bound historical snapshot of
   prose referring to numbered work from A0 through A67. Its migration map
   explicitly says it did **not** reconstruct those entries as individual
   structured attempts
   (`attempts/imported-unversioned/evidence/MIGRATION_MAP.md`). Some files
   contain several candidates or decisions, some number ranges are grouped,
   and A20--A35 remain embedded in the large current-approach record. Exact
   accept/refine/reset totals for A0--A67 would therefore be invented.
2. A68--A89 are 22 complete `GoalAttempt` resources and can be counted exactly.

The nominal labels span A0--A89, but they are not 90 equivalent visual
iterations. Some were migrations, calibration, interface investigations, or
preflight vetoes; others contain many internal variants. A0 alone reports more
than 30 prior scripted variants
(`attempts/imported-unversioned/evidence/attempt_00.md`).

## Acceptance contract and actual movement

All eight criteria are required: reference likeness, measured silhouette,
plush construction, two blind 8/10 reviews, reusable structure, animation
readiness, technical integrity, and repository delivery
(`criteria.yaml`; legacy formulation in
`attempts/imported-unversioned/evidence/acceptance.md`). The numeric thresholds
were explicitly inferred to operationalize the user's quality request, not
supplied by the user
(`attempts/imported-unversioned/evidence/requirements.md`).

### Formal structured movement, A68--A89

| Slice | Decisions | Criterion passes | Fails | Unverified | Accepted result |
| --- | --- | ---: | ---: | ---: | --- |
| A68--A78 | 11 reset | 0 | 34 | 54 | none |
| A79--A89 | 7 reset, 4 refine | 0 | 16 | 72 | none |
| Total | 18 reset, 4 refine | **0** | **50** | **126** | **none** |

The apparent change from `fail` to `unverified` in recent attempts is not a
quality gain. It mostly reflects narrower module scope. For example, A85, A87,
and A88 close with all eight whole-goal criteria unverified even though their
internal module gates preserve a provisional survivor. Criteria 005--008 have
remained unverified in every structured attempt. No final rig, deformation
test, integrity audit, clean reuse test, or delivery candidate has existed.

### Legacy movement, A0--A67

- A0 rejected the inherited 3/10 complete baseline. A1 accepted only its
  repository migration; A2 accepted only reference calibration. A3--A19 all
  reject their geometry or preflight direction. A15--A19 are five consecutive
  preflight families with no new candidate pixels
  (`attempt_00.md` through `attempt_19.md`).
- The historical record contains two genuine role-specific geometry passes:
  the hidden sewn head-support cushion V3 and A28 rear felt panel V4. A28 V4
  scored 8.0--9.5 in its isolated role
  (`attempts/imported-unversioned/evidence/current_attempt.md`, section
  “Attempt 28 rear felt panel”). By A32 both are explicitly retained as
  evidence only and no longer govern visible geometry because their combined
  interfaces could not produce a viable head
  (`current_attempt.md`, section “Attempt 32”).
- A37--A44 rejected every complete visible module or assembly. A43 reached
  only 4/10 likeness and 2/10 construction; A44's durable gains were cameras,
  review tooling, measurements, and a rig-interface prototype, not approved
  geometry
  (`attempts/imported-unversioned/evidence/attempt_37_44.md`).
- A45's working ladder accumulated three relative wins--face/hairline,
  eyes/locks, and constructed sleeves--but explicitly states that no major
  visible module passed the absolute gate and none was promoted
  (`attempts/imported-unversioned/evidence/attempt_45.md`).
- A52--A67 repeatedly rejected head, skirt, bow, receiver, and inverse-pattern
  families. No candidate was promoted
  (`attempt_52_57.md` through `attempt_67.md`; summary in
  `attempts/imported-unversioned/plan.md`).

Thus the legacy archive contains useful isolated geometry knowledge and
relative comparisons, but **zero accepted complete-head or complete-plush
result**. The structured importer correctly marks all eight criteria
unverified (`attempts/imported-unversioned/attempt.yaml`).

## Complete structured attempt ledger, A68--A89

| Attempt | Close | Formal verdicts | Acceptance-visible result |
| --- | --- | --- | --- |
| A68 sculpt coupon | reset | 4 fail / 4 unverified | Two Grab checkpoints produced negligible movement or dimples; no survivor. |
| A69 head/cap interface | reset | 4 / 4 | Directionally correct rear movement lost lock contact and retained the monolithic helmet; no survivor. |
| A70 constructed receiver | reset | 4 / 4 | Four states converged on rigid cards/mask/helmet; no survivor. |
| A71 macro sculpt | reset | 4 / 4 | Synthetic Grab support failed before an artistic sculpt; no survivor. |
| A72 parametric receiver | reset | 4 / 4 | Rounded mattress/egg, 5/10 at a 6/10 intermediate gate; no survivor. |
| A73 profile cage | reset | 4 / 4 | Mechanically clean but visually worse at 4/10; no visible-layer coupon; no survivor. |
| A74 visible envelope | reset | 4 / 4 | Face card, hard canopy, blade locks, bald rear; no survivor. |
| A75 face cushion | reset | 4 / 4 | Spherical muzzle/cavity and detached beige islands; no survivor. |
| A76 live sculpt interface | reset | 0 / 8 | Interface investigation only; no mutation or visual candidate. |
| A77 sparse head/hair | reset | 2 / 6 | Mechanical face/contact veto before rendering; no visual survivor. |
| A78 flattened head rest | reset | 0 / 8 | Causal-reach veto before geometry; no candidate. |
| A79 paired hair field | reset | 2 / 6 | Bald crown, rigid rear curtain, detached board; no survivor. |
| A80 diagnostic hair blockout | reset | 2 / 6 | Both variants remained bald/helmet/ball/capsule assemblies; no survivor. |
| A81 live joint rest form | reset | 2 / 6 | Beige plaque and deep box/egg support; no survivor. |
| A82 atomic head shell | reset | 2 / 6 | One uninterrupted convex cuboid/helmet; no survivor. |
| A83 incremental sculpt | reset | 3 / 5 | All repairs undone; C1b retained only as an unaccepted comparator. |
| A84 neutral head cushion | reset | 3 / 5 | Rounded rectangular head; S02 retained only as method evidence. |
| A85 constructed head cage | refine | 0 / 8 | **S12 provisional head receiver survives**, not standalone approval. |
| A86 collar/body interface | reset | 0 / 8 | Collar cannot hide lateral underside without becoming a rigid yoke; no survivor. |
| A87 front hair frame | refine | 0 / 8 | **S04 three-piece front-hair module survives** at mostly 7/10, medium 6/10. |
| A88 crown interface | refine | 0 / 8 | **S07 two-panel crown survives** at about 7/10; current viable parent. |
| A89 lower-rear hair | refine | 2 / 6 | No A89 survivor; shield/cape, fins, reversed depth, exposed receiver. |

The exact decisions and criterion verdicts are in each attempt's
`attempt.yaml`; the result summaries are in each `result.md` under
`projects/renders/goals/reimu-fumo/attempts/<attempt-id>/`.

## What is actually preserved

### Current cumulative scratch parent

1. **A85 S12 head receiver** -- Blend SHA-256
   `982da6404ea6edcbb4432903e67dad4ee5c130a203a5a5727a374b773fc9ad8a`.
   Its head-only silhouette gates pass, but it is a smooth rectangular block
   with an unresolved underside cavity and is explicitly not standalone
   approval
   (`a85-constructed-head-cage/result.md`,
   `evidence/S12_SIX_VIEW_REVIEW.md`).
2. **A87 S04 front hair** -- cumulative Blend SHA-256
   `da300595b73185d7c4deece12f7e98e7d6889e8f738c94709c828bd601e58d2a`.
   It adds one shallow fringe and two cheek locks. Its blind module scores are
   6--7/10, below the final 8/10 gate
   (`a87-front-hair-frame/result.md`,
   `evidence/A87_S04_BLIND_REVIEW.md`).
3. **A88 S07 crown** -- cumulative Blend SHA-256
   `0a30e2af3142081648bb3137ad75d6d1cc73de55e9f830f85bad1f85e92c8788`.
   It adds two crown panels and closes literal gaps while preserving A87. Its
   local scores are about 7/10 and it remains explicitly provisional
   (`a88-crown-interface/result.md`,
   `evidence/A88_S07_BLIND_REVIEW.md`).

These are six neutral geometry owners in one partial assembly: one head, three
front-hair pieces, and two crown panels. There is no face, bow, lower-rear
hair, body, garment, feet, material set, or production rig in the survivor.

### Not preserved as current result

- A89 S06 is diagnostic only; A88 remains the parent
  (`a89-lower-rear-hair/result.md`).
- Legacy support V3 and rear V4 were role-specific passes later demoted to
  evidence-only inputs by A32.
- The A45 rung-003 assembly remains the strongest rejected complete-character
  comparator, not an approved parent or tracked promotion.
- **Nothing has been promoted into the tracked reusable Blend.**

## Recurring defects

The legacy failure ledger already recorded lower-bound occurrence counts
through A67:

| Defect family | Recorded occurrences |
| --- | ---: |
| Weak side/rear silhouette | 50 |
| Flat or disconnected construction | 46 |
| Helmet/fat/tall head | 31 |
| Incorrect hair | 23 |
| Foot/skirt clipping or tangency | 12 |
| Human/tall proportions | 11 |
| Rigid bow | 11 |
| Missing or anatomical feet | 11 |
| Bland material/contact | 3 |
| Reimu expression loss | 2 |

Source: `attempts/imported-unversioned/evidence/failures.md`. The counts
overlap and are explicitly incomplete: A65/A66 recurrences were not folded
into them, and A68--A89 add further repeats.

The same categories continue after the import. A68--A84 repeatedly alternate
between monolithic helmet/egg/box/receiver forms and thin cards, curtains, or
disconnected leaves. A85--A88 achieve local coverage and ownership, but the
actual pixels remain block/cap/visor-like and below the absolute gate. A89
then recreates the long-running rear curtain/cape/fin defect.

The stable causal pattern is: **a mechanical or scalar gate passes for a
representation that cannot express the required visual category**. Topology,
manifoldness, coverage, gap counts, exact source anchors, or render speed then
improve while the result still reads as a helmet, card, cape, cone, or rigid
assembly.

## Process/document output versus result progress

### Durable record volume

- Complete A89-cutoff corpus: **219 text files**, **27,793 lines**,
  **1,535,007 bytes**. This count explicitly excludes the concurrently opened
  process-review attempt.
- Imported historical snapshot: **41 Markdown files**, **10,144 lines**,
  approximately **639 KB**.
- Structured A68--A89 records: **173 YAML/Markdown files**, including **107
  evidence documents**, **17,309 lines**, **876,324 bytes**.
- A79--A89 alone: **99 YAML/Markdown files**, **9,381 lines**, approximately
  **481 KB**.

Against that output, the formal result is 0 accepted attempts, 0 passed
criteria, 0 tracked asset changes, and 3 cumulative provisional modules. A
simple, explicitly imperfect productivity proxy is therefore:

- **7.3 structured attempts per current provisional module** (`22 / 3`);
- **51 structured Markdown records per provisional module** (`153 / 3`);
- **no finite records-per-accepted-result ratio**, because there is no
  accepted result.

This is not an argument against evidence. The calibration, fixed cameras,
immutable snapshots, blind review, hash binding, component IDs, and render
runner are valuable. It is evidence that archival rigor and process work have
far outpaced acceptance-visible breadth.

The ignored attempt scratch reinforces the scale: **509 primary Blend files
(447 unique hashes)** and **6,445 PNGs (6,209 unique hashes)** through A89.
A68--A89 alone account for **126 primary Blends (118 unique)** and **1,188
PNGs (1,160 unique)**. These include source copies, checkpoints, masks,
diagnostics, and duplicate views; unique bytes are not independent artistic
candidates. More importantly, most legacy visual links point into ignored
`out/`; the imported record preserves prose and hashes, not durable visual
evidence. The history is stronger as a failure/provenance archive than as a
reproducible accepted-result chain.

The flattened import also did not rewrite legacy relative links. A bounded
check of its 41 Markdown files found 628 relative links; 594 do not resolve
from their imported locations and 499 target `out/`. Many referenced files
still exist when their paths are reinterpreted from the original archive, but
that is fragile local provenance rather than durable candidate-bound evidence
(`attempts/imported-unversioned/evidence/MIGRATION_MAP.md`,
`attempts/imported-unversioned/evidence/artifacts.md`).

## Why recent work appears stagnant

1. **The active evidence is too local.** A85--A88 deliberately exempt missing
   context. Local coverage and interface gates can pass while the artifact is
   still an unrecognizable partial head. Criterion 001 requires the broad
   head *and* compact seated body, full hair, bow, face, sleeves, skirt, tie,
   and feet.
2. **Provisional 6--7/10 modules became downstream constraints.** A87 and A88
   are honest about not passing the 8/10 gate, yet later modules route around
   them. This risks protecting weak geometry because it passed a narrower
   owner test.
3. **The process optimized what was easy to count.** Gap pixels, clearance,
   topology, protected-object hashes, and render latency are objective. They
   do not prove stuffed-fabric construction or likeness. A89 is the clearest
   example: closing gaps fused separate leaves into a worse shield.
4. **Representation churn replaced direct art direction.** Analytic lofts,
   cages, shells, panels, generated cloth, and synthetic sculpt operators
   repeatedly produced the same helmet/card/cape categories. The history often
   diagnosed this correctly, then tested another neighboring generator before
   a visibly authored whole-context form survived.
5. **Whole-character feedback disappeared after A83.** A83 identified two
   coupled structural owners--head/hair and seated torso/pelvis/garment--but
   A84--A89 spent six attempts almost entirely on the head. No current body,
   hem, sleeve, or foot owner exists, so the compact seated proportion and
   important occlusion/contact interfaces remain untested.
6. **`refine` is not a progress metric.** A85, A87, and A88 mean “preserve a
   provisional module.” A89 means “pause safely while rejecting every A89
   candidate.” Counting four refines would materially overstate progress.
7. **Review placement, not rendering, is now the latency bottleneck.** A82
   reached its controlling pair in 3m35s; A87 reports an approximately
   18-second three-view render. A89 still authored S00--S06 before the decisive
   independent cape/fin review. The expensive mistake is several geometry
   states before categorical visual veto, not Blender throughput.

## Genuine accomplishments

The review should not erase real progress:

- reference authority, camera calibration, fixed views, landmark normalization,
  and all-reference review are substantially better than the early session;
- candidate copies are immutable, hash-bound, and protected from accidental
  tracked-file contamination;
- rendering and basic comparison are now fast enough for an immediate visual
  loop;
- A85 is a better shallow receiver than prior rounded cubes;
- A87 establishes more defensible front-hair ownership; and
- A88 demonstrates a useful component-ID/clearance diagnosis while preserving
  its parent exactly.

These are prerequisites and local gains. They are not evidence that the
requested plush is close to completion.

## Constraints the recovery plan must honor

1. Treat exact A88 S07 as the only current provisional parent and all A89
   geometry as rejected evidence.
2. Before another downstream polish pass, render A88 in low-cost whole-head
   and seated-body context with diagnostic face, bow, rear-hair, and body
   proxies. This is a go/no-go test of the current scaffold, not promotable
   geometry.
3. If context still reads as a cuboid helmet, reopen the smallest owning A85,
   A87, or A88 module. Do not preserve it merely because a local gap or
   silhouette subgate passed.
4. Resume the seated-body structural owner in parallel with a bounded rear-hair
   coupon. Whole-character proportion and contact cannot wait until head hair
   is polished.
5. Put the implementation-blind reviewer and component-ID pass on the first
   decision-bearing render. Permit at most one geometry state before a
   keep/undo/reset decision.
6. Use direct, controllable sparse geometry and local sculpt/edit operations
   for visible camber, stuffing, root compression, overlap, and asymmetric
   free edges. Keep analytical generation for cameras, measurements, masks,
   simple initial contours, and verification--not every visible form.
7. After the next two provisional survivors, require one cumulative neutral
   full-subject packet. No further module polish proceeds unless that packet is
   recognizable and free of categorical helmet, cape, cone, tube, clipping,
   or disconnected-part failures.
8. Keep each ordinary visual cycle to one short causal plan, one immutable
   render/ID packet, one blind review, and one concise result. Add process
   machinery only when it changes the next modeling decision.

## Final assessment

The process is now good at proving that bad candidates are bad. It is not yet
good at producing a strong candidate. The decisive recovery is not another
large generator, more documentation, or a looser gate. It is restoring
continuous whole-result visual judgment to a direct, small-edit authoring loop
and refusing to freeze locally measurable geometry that remains globally weak.
