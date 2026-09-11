## ADDED Requirements

### Requirement: Reimu reference fidelity

One exact candidate SHALL match the canonical Reimu Fumo variant in every
fixed view defined by review contract
`sha256:4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`,
with no major identity or silhouette defect. Evidence SHALL be a
self-contained candidate packet inspected against the exact contract bytes,
with the candidate and packet digests recorded. This preserves legacy
`criterion-001`, revision 3, required by criteria revision 4.

#### Scenario: Review the exact variant in every fixed view

- **WHEN** a candidate is proposed for acceptance
- **THEN** its exact packet is compared with the controlling references in the contract's front, rear, side, three-quarter, and mirrored three-quarter views
- **AND** the recorded candidate and packet digests identify the inspected bytes
- **AND** any major identity or silhouette defect prevents acceptance

### Requirement: Reimu measured landmarks

Calibrated comparisons SHALL put every critical canonical-view landmark
within `0.03 Wh`, silhouette extrema and major gaps within `0.05 Wh`, and
uncertain depth landmarks inside the frozen bands bound by review contract
`sha256:4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`.
Aligned measurements SHALL be published and tied to the exact candidate,
contract, landmark file, references, cameras, and render digests. This
preserves legacy `criterion-002`, revision 3, required by criteria revision 4.

#### Scenario: Validate calibrated measurements and uncertainty bands

- **WHEN** the exact candidate's landmark evidence is evaluated
- **THEN** calibration and annotated endpoints follow the frozen contract and landmark file, with camera calibration within `0.02 Wh` for a valid comparison
- **AND** all critical landmarks meet `0.03 Wh`, silhouette extrema and major gaps meet `0.05 Wh`, and uncertain depths stay inside the frozen bands, including stuffed-head depth `0.66-0.82 Wh` and outer-hair depth `0.71-0.87 Wh`
- **AND** the aligned measurements identify the candidate, contract, landmark file, references, cameras, and renders by digest
- **AND** a failed calibration invalidates the comparison instead of being scored as a model failure

### Requirement: Reimu plush construction

The neutral model SHALL read as sewn and softly stuffed, with plausible panel
thickness, seams, gathers, compression, attachment, and contact and no
helmet, card, armor, hollow-tube, cone, floating, clipping, or
disconnected-primitive failure. Review SHALL inspect the exact fixed-view
and presentation renders in the immutable candidate packet. This preserves
legacy `criterion-003`, revision 2, required by criteria revision 4.

#### Scenario: Reject visible construction failures

- **WHEN** the neutral candidate's fixed-view and presentation renders are inspected
- **THEN** its parts show plausible sewn and softly stuffed construction, attachment, and contact
- **AND** none of the listed construction failures is present
- **AND** visible floating, clipping, disconnected construction, or accidental tangency above `0.02 Wh` is a major failure regardless of averages, as specified by the frozen landmark contract

### Requirement: Reimu independent visual review

Two implementation-blind reviewers identified by stable ID and role SHALL
recognize the exact variant, score every applicable visual category at least
`8/10`, and report no major visible failure on the same candidate and
review-contract digests. Both independent category matrices and reviewer
identities SHALL be published in the immutable candidate packet. This
preserves legacy `criterion-004`, revision 2, required by criteria revision 4.

#### Scenario: Obtain two complete reviews of the same candidate

- **WHEN** final visual acceptance is proposed
- **THEN** two implementation-blind reviewers with recorded stable IDs and roles each recognize the exact canonical variant
- **AND** each review scores every applicable visual category at least `8/10` and reports no major visible failure
- **AND** both category matrices and identities are in the immutable packet and refer to the same candidate and review-contract digests

### Requirement: Reimu reusable structure

The deliverable SHALL contain one clearly named reusable collection with
intentional meshes, materials, armature, actions, controls, transforms, and
no review-only content or missing dependency. Verification SHALL append only
the reusable collection into a blank file, clean-reopen it, and publish the
structure and dependency audit. This preserves legacy `criterion-005`,
revision 2, required by criteria revision 4.

#### Scenario: Use the collection in a clean file

- **WHEN** only the reusable collection is appended into a blank Blender file and that file is clean-reopened
- **THEN** its intentional meshes, materials, armature, actions, controls, and transforms are present and usable
- **AND** no review-only content or missing dependency is found
- **AND** the structure and dependency audit is published for the exact candidate

### Requirement: Reimu animation readiness

The same candidate SHALL deform without visible tearing or clipping in
seated neutral, actual head-yaw turn, arm-wave, and combined validation
actions. Verification SHALL exercise the stored actions and publish
fixed-view pose renders tied to the exact candidate digest. This preserves
legacy `criterion-006`, revision 2, required by criteria revision 4.

#### Scenario: Exercise every validation action

- **WHEN** the exact candidate's stored seated-neutral, actual head-yaw-turn, arm-wave, and combined validation actions are exercised
- **THEN** no action shows visible tearing or clipping
- **AND** the published fixed-view pose renders identify the same candidate digest used for the other acceptance checks

### Requirement: Reimu technical integrity

The same candidate SHALL have no unintended non-manifold boundary,
degenerate geometry, duplicate object name, broken modifier, dependency
cycle, local path disclosure, or unpacked required resource. Verification
SHALL run the pinned-Blender clean-open audit against the exact candidate
and publish its machine-readable result. This preserves legacy
`criterion-007`, revision 2, required by criteria revision 4.

#### Scenario: Audit the clean-opened candidate

- **WHEN** pinned Blender clean-opens and audits the exact candidate
- **THEN** no unintended non-manifold boundary, degenerate geometry, duplicate object name, broken modifier, dependency cycle, local path disclosure, or unpacked required resource is found
- **AND** the machine-readable audit result is published for those exact bytes

### Requirement: Reimu exact-byte delivery

A self-contained packet SHALL prove pinned-Blender rendering from the exact
committed Git LFS object, preserve the exact controlling reference identities
with the repository fan-work notice, and ensure
`//projects/renders:reimu_fumo` exposes only the accepted reusable asset.
Verification SHALL hydrate and rerender the committed LFS object, compare
packet and object digests, verify the fan-work notice and reference
identities, and inspect the final Bazel target. This preserves legacy
`criterion-008`, revision 3, required by criteria revision 4.

#### Scenario: Verify the committed asset before exposing it

- **WHEN** the reusable asset is prepared for final delivery
- **THEN** the exact committed Git LFS object is hydrated and rerendered with pinned Blender and its object and packet digests are compared
- **AND** the self-contained packet preserves the exact controlling reference identities and repository fan-work notice
- **AND** `//projects/renders:reimu_fumo` exposes only the asset that passed every required criterion on those exact bytes
- **AND** a donor, rejected candidate, or successful capability coupon is not exposed as the accepted reusable asset

### Requirement: Preserve the blocked Reimu continuation

The migrated change MUST preserve the open outcome, blocked execution, and
absence of an accepted candidate or result recorded at migration. It MUST
remain open with no accepted candidate or result until all required model
criteria pass. Autonomous modeling MUST NOT resume unless a sponsor-approved
skilled Blender artist
or a genuinely different organic authoring capability, independently proven
outside the failed broad-transform and fixed-dose families, is available.
The migration itself MUST NOT be treated as authorization to resume modeling
or capability experiments. The eight acceptance requirements above preserve
[criteria revision 4](../../provenance/source/criteria.yaml); this continuation
constraint preserves the [terminal result](../../provenance/source/attempts/flatten-dose-response-018/result.md)
and does not replace an acceptance criterion.

#### Scenario: Migrate the terminal failed capability investigation

- **WHEN** the legacy goal and attempt `flatten-dose-response-018` are migrated to OpenSpec
- **THEN** the outcome remains open, execution remains blocked, and every required model acceptance task stays unchecked
- **AND** the failed 32.44312678772848 percent reduction at 27 strokes is preserved against the unchanged 35 percent gate, with root and tip tests recorded as withheld
- **AND** accepted transport or capability evidence does not imply that any Reimu model criterion passed

#### Scenario: Evaluate a proposed resumption

- **WHEN** a future executor considers continuing the unfinished asset
- **THEN** the executor first establishes the recorded artist or different-capability unblock condition and the authority to execute
- **AND** generic desktop input transport, more Flatten dose, parameter tuning, or repetition of the rejected families does not satisfy that condition

### Requirement: Continue the process contract without the removed goal tool

The authoritative [process contract](../../../../../assets/reimu_fumo/PROCESS.md) SHALL
remain executable after the legacy goal tool was removed. Its attempt, plan,
checkpoint, blocker, and closure state SHALL be recorded in this OpenSpec change
and its durable evidence, not in a goal store. The contract SHALL NOT require a
removed goal command, goal store, or goal record format for any required stage
gate.

#### Scenario: Execute a stage gate after the goal tool removal

- **WHEN** an authorized coordinator executes a required Reimu stage gate
- **THEN** the attempt, plan, checkpoint, blocker, and closure state are recorded in this change and its evidence
- **AND** no removed goal tool, goal store, or goal record format is required to satisfy the gate
