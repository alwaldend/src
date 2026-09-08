## Context

This is the native OpenSpec continuation of the durable legacy goal
`reimu-fumo-finish`, titled "Finish a reference-faithful reusable animated
Reimu Fumo". The source snapshot has goal resource version `58`, generation
`1`, lifecycle generation `15`, and criteria revision `4`. Its last state
observation is `2026-09-04T12:58:04.377051268Z`.

Outcome remains **open** and execution remains **blocked**. No attempt is
active; `acceptedAttemptID` and `acceptedResultDigest` are empty. There are no
parent, dependency, or supersedes relationships. All eight model acceptance
criteria remain unresolved. Accepted capability experiments and an accepted
measurement correction are not accepted Reimu candidates.

The byte-preserved [goal](provenance/source/goal.yaml),
[criteria](provenance/source/criteria.yaml), historical criteria revisions,
and all eighteen attempt directories under `provenance/source/attempts/`
retain the original record and evidence. The repository migration receipt
owns source-to-destination mapping and byte verification. This document is a
readable continuation, not a replacement for those historical receipts.

## Goals / Non-Goals

The intended result is one reusable animated asset that meets every criterion
on the same exact bytes. Native requirements in
`specs/project-renders/spec.md` preserve the current criteria and associate
them with their legacy IDs and revisions.

The present migration does not resume modeling or capability trials, alter
references or thresholds, accept a donor or rejected artifact, install or
configure Blender, or expose an accepted asset target. Historical binaries
and images described as ignored scratch have not been recovered or verified
by this migration.

## Decisions

### Keep the work open and preserve the execution blocker

The latest closed attempt,
[flatten-dose-response-018](provenance/source/attempts/flatten-dose-response-018/attempt.yaml),
is an investigation with review decision `reset` and no criterion verdicts.
Its [result](provenance/source/attempts/flatten-dose-response-018/result.md)
rejects plan `flatten-dose-response-v16` and blocks autonomous execution.
Do not convert its closed-attempt state into a completed change or mark any
model acceptance task complete.

Resume only with a sponsor-approved skilled Blender artist or a genuinely
different organic authoring capability independently proven outside the
failed broad-transform and fixed-dose families. Generic desktop input
transport, additional Flatten dose, and further parameter tuning do not
meet the condition. A migration instruction is not approval to resume those
operations.

### Preserve the controlling acceptance contract

All candidate evidence must bind the exact candidate, packet, and review
contract bytes. The current review contract is
`projects/renders/assets/reimu_fumo/review_contract.json`, SHA-256
`4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`.
Its landmark file SHA-256 is
`133d741d9252f34cee4e27990b0ea8d0710dc7408f43eb2a38b47d3c52a98140`.

The five fixed views are front, rear, side, three-quarter, and mirrored
three-quarter, using the contract's orthographic cameras at 512 by 512 pixels
and `0.292 m` orthographic scale. `Wh` is the outer front head-and-hair
envelope width excluding bow and isolated ribbon tips. Camera calibration
must be within `0.02 Wh` for a valid comparison. Critical canonical-view
landmarks must be within `0.03 Wh`, silhouette extrema and major gaps within
`0.05 Wh`, and uncertain depth landmarks inside the contract's frozen bands.
Visible floating, clipping, disconnected construction, or accidental tangency
above `0.02 Wh` is a major failure regardless of averages.

The measurement correction in
[attempt 008](provenance/source/attempts/measurement-band-correction-008/result.md)
sets stuffed-head depth to `0.66-0.82 Wh` and outer-hair depth to
`0.71-0.87 Wh`. It accepted evidence correction only, left the rejected
integrated-hair candidate unchanged, and did not pass criterion 002.

### Retain strategy outcomes without reopening rejected families

The [goal source](provenance/source/goal.yaml) preserves the complete strategy
text and rejection reasons. This table maps every plan to its continuing
status and the relevant attempts; no row supplies new execution authority.

| Plan                         | Legacy state | Attempt evidence and continuing constraint                                                                                                                                                                                                                                                                                                                 |
| ---------------------------- | ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `recovery-clay-v1`           | superseded   | [001](provenance/source/attempts/recovery-capability-001/result.md), [002](provenance/source/attempts/camera-capability-002/result.md), [003](provenance/source/attempts/head-hair-rebuild-003/result.md): recover immutable A157 as visual scaffold and A202 as technical donor; recovery and transport evidence do not accept either as the final asset. |
| `curved-cap-v2`              | rejected     | [004](provenance/source/attempts/curved-cap-004/result.md): second shell or panel-box failure retired the head-covering representation family.                                                                                                                                                                                                             |
| `integrated-textile-v3`      | rejected     | [005](provenance/source/attempts/integrated-hair-005/result.md): integrated stuffed substrate remained a monolithic helmet; blind clay readiness was 3/10.                                                                                                                                                                                                 |
| `shingled-locks-v4`          | rejected     | [006](provenance/source/attempts/shingled-hair-006/result.md): evaluated head lattice embedded analytic panels, exposing crown and detached oval side pads.                                                                                                                                                                                                |
| `evaluated-lobes-v5`         | rejected     | [007](provenance/source/attempts/evaluated-lobes-007/result.md): segmented cap, giant side pad, and overlapping pillows retired the padded-lobe family after two failures.                                                                                                                                                                                 |
| `event-simulated-sculpt-v6`  | superseded   | Capability-coupon planning lineage in the goal source; a proposed event-path failure would require isolated capability review, never generic logged-in desktop control.                                                                                                                                                                                    |
| `event-simulated-sculpt-v7`  | superseded   | Capability-coupon planning lineage in the goal source; the unmapped `0x0` pinned window was a launch-host defect addressed by v8, not proof that sculpt events cannot work.                                                                                                                                                                                |
| `mapped-gui-sculpt-v8`       | accepted     | [009](provenance/source/attempts/mapped-gui-sculpt-009/result.md), [010](provenance/source/attempts/macro-sculpt-010/result.md): mapped loopback MCP and timer-spaced native Grab established bounded transport and macro-authoring capability, exact save, undo, pinned reopen, and render. No Reimu criterion passed.                                    |
| `a157-bow-span-v9`           | rejected     | [011](provenance/source/attempts/a157-bow-span-011/result.md): assembled-view extremum Grab failed shape distribution and right-side picking; repeating it cannot produce a valid symmetric bow.                                                                                                                                                           |
| `a157-tail-assembly-v10`     | rejected     | [012](provenance/source/attempts/a157-tail-assembly-012/result.md): native modal numeric input completed with zero displacement; key-stream spelling changes are equivalent retries.                                                                                                                                                                       |
| `a157-tail-vector-v11`       | rejected     | [013](provenance/source/attempts/a157-tail-vector-013/result.md): explicit vectors reached span but made fin-like tails, abrupt roots, and weak ruffles.                                                                                                                                                                                                   |
| `native-crown-patch-v12`     | rejected     | [014](provenance/source/attempts/native-crown-patch-014/result.md): scalar cut left a horizontal slab, broad pale gaps, and rigid disconnected-looking hair.                                                                                                                                                                                               |
| `native-manual-head-v13`     | rejected     | [015](provenance/source/attempts/native-manual-head-015/result.md): native shaping retained pillow, card, armor, crown-rail, and floating-root failures; the unsaved state was discarded and protected A157 unchanged.                                                                                                                                     |
| `localized-sculpt-stack-v14` | superseded   | [016](provenance/source/attempts/localized-sculpt-coupon-016/result.md): disposable multi-brush capability investigation; the follow-up changed diagnosed empty-scene dependencies only.                                                                                                                                                                   |
| `localized-sculpt-stack-v15` | rejected     | [017](provenance/source/attempts/localized-sculpt-coupon-017/result.md): three isolated Flatten strokes achieved only 4.8236 percent plane-variance reduction versus the frozen 35 percent gate; root and tip stages withheld.                                                                                                                             |
| `flatten-dose-response-v16`  | rejected     | [018](provenance/source/attempts/flatten-dose-response-018/result.md): fixed-dose response reached the 27-stroke ceiling at 32.4431 percent with terrace-like ridges; autonomous modeling blocked.                                                                                                                                                         |

### Preserve the terminal failure and exact evidence identities

Attempt 018 changed only cumulative identical Flatten dose in three-stroke
blocks on the disposable attempt-017 fixture. It opened no Reimu model. The
fixture, mask, front view, PLANE brush, strength `0.65`, scene size `0.58`,
trajectories, timing, metric, threshold, and renderer were frozen. The
[plan](provenance/source/attempts/flatten-dose-response-018/plan.md) required
at least 35 percent plane-variance reduction within at most 27 strokes.

Its hard stops were non-plane or `CONTROL_plane` displacement above `1e-6`,
cumulative maximum displacement above `0.20` scene units, any block-to-block
variance increase, two consecutive block reductions below
`0.0003428438941380965`, or failure to reach the gate by 27 strokes. The
terminal state hit the last condition. Variance fell monotonically from
`0.014215173048885853` to `0.009603326433540808`, a
`32.44312678772848%` reduction. Maximum cumulative displacement was
`0.11849009312857924`; non-plane and control displacement were exactly zero.
The final block was the first below the low-response floor, not two
consecutive low-response blocks.

The frozen metric uses the current-coordinate window
`front layer; abs(current vertex x) <= 0.65; abs(grid u) <= 0.65`.
Substituting baseline X selected different vertices and was rejected before
any stroke. Twenty-four native undo operations restored the exact partial
coordinate digest. Technical success did not satisfy the shape-control gate.

Root and tip tests were conditional on a plane pass and did not run. Their
unchanged gates were root contact width increasing at least 25 percent,
root curvature discontinuity decreasing at least 30 percent, tip width at
most 0.55 of mid-panel width, tip roughness decreasing at least 25 percent,
non-target controls within `1e-6`, and exact undo. Their historical one-block
limit does not authorize a new trial after the terminal failure.

The full durable [pinned receipt](provenance/source/attempts/flatten-dose-response-018/evidence/pinned_failure_evidence.md)
and [writer receipt](provenance/source/attempts/flatten-dose-response-018/evidence/writer_result.md)
preserve every event, block metric, image identity, and verification context
they recorded. Key SHA-256 identities are:

| Artifact or record                         | SHA-256                                                            |
| ------------------------------------------ | ------------------------------------------------------------------ |
| Baseline fixture file                      | `6f49bd4e0a8af6b45870d9d4224a520c1398e52e6e9c42f8fb5bee7b8c17118e` |
| Partial three-stroke input file            | `2428de0a0b65e572de9437a8d3ef35f1ee21c18bd9dbf27ff01de1816418c0bd` |
| Terminal 27-stroke failure file            | `b9ad59c15901b0fe22cd96208e60792e2dc35e8bc0f3e73d0e1a8181697dd6fc` |
| Baseline coordinates                       | `41ee23670f67335ac070d95bd782436f53405034f1e24efdebdad709f7d47df2` |
| Partial coordinates, also restored by undo | `f2fdcbdaea90335b9de47861e8ce59b64896ec81853d19ee2ea675725d9fb16e` |
| Terminal coordinates                       | `fc3d57571800ecade593e64d2750687f883c7da9ef6177fa9eec622a8346c2e0` |
| Attempt-018 plan bytes                     | `26921b2597c1836a953ec5b848bcb9bff5881ef260b2afd75a35c57f1d1b5480` |
| Attempt-018 result bytes                   | `686fac0462003a203494dfa7e16719dcbcef23a354ffeb2c981ac1bd1acd4412` |
| Pinned failure receipt bytes               | `a98b2600f1feb534f4a1ab23ef11e7263bbb766975ce9085fd6124571d07b2b5` |
| Writer receipt bytes                       | `a299a04ddf641882c4a1502c5a32ca375eb6effd7e2fa114d402a96de15804cb` |
| Recorded pinned audit JSON                 | `ed9e21597d44d601de6532eb472e0adc737c1ce735890113ad59e380330900e7` |
| Recorded pinned render manifest            | `73065da607856f0c6ffc92d1ab2191aa249b80db99370f26f85791b9b4f12267` |
| Recorded render `READY`                    | `f6b2aad41c0c6cc88b0d7185d64d8db2c6f89091efde3e07c199b06805a70bf0` |

The legacy attempt binds criteria revision 4 by its goal-tool canonical
digest `sha256:f5fe4b5c4efc12d5b82ad42cf0d41fd7274f48b421abf42e9db1ce5ae316d8a3`.
That is the attempt's recorded canonical digest, not the YAML file-byte
checksum. The migrated `criteria.yaml` byte checksum is
`8e9f666e12fe455b09395850611026edfa6b35455b5c10b0c26c7d596fbecd21`.
The migrated `goal.yaml` byte checksum is
`26fc3104902a160e859fd4d660678f7b5debbf29f591c79b3f0d8be62f8aff71`.

## Risks / Trade-offs

History includes a technically successful transport path and an accepted
capability plan. Treating either as asset acceptance would erase the visual
failure and violate the source criteria. Native OpenSpec tasks remain
unchecked and this change stays outside the archive to preserve the
distinction.

Historical receipts refer to ignored scratch binaries, renders, and runtime
observations. Preserving receipt bytes does not reconstruct those artifacts
or establish present host capabilities. Any future authorized continuation
must revalidate the exact inputs it needs without substituting unverified
artifacts or weakening the frozen contract.

## Migration Plan

Preserve original files byte-for-byte with the repository migration receipt,
validate this native OpenSpec change, and retain it as open and blocked.
Do not execute the unchecked modeling or acceptance tasks as part of the
record migration. On an authorized future resume, first establish the
recorded unblock condition, then record the new execution state and evidence
in this change before developing a new candidate.
