# Phase 1 shared-contract evidence

## Subject

The tested subject is base commit
`6780d53a69e32064d648e6a04f1c0cecd7d713fd` plus the following exact source
files:

| SHA-256 | Path |
| --- | --- |
| `f96eda96c7a2f42345219db8f4b286dd55f84fedf368e2930ca581d5b46526d9` | `projects/goal/internal/fsstore/checkpoint.go` |
| `11642ecdd8f638d1925b59efbf544699f57a2808ba41cc6c89710285b16dae9b` | `projects/goal/internal/fsstore/store.go` |
| `2ebca2f2a8163364f62baf39e24a6a75174a80e361f807abdec1542267971c79` | `projects/goal/internal/fsstore/store_test.go` |
| `a6c622363fcd764df7e1ab9dcb52ea8d79d4d2711c1a43fc2139c9ae8b0b8e24` | `tools/agents/BUILD.bazel` |
| `e3d192a3cb8ec302dcd0db451ce02276a61c85fe505f86b1ed11e27649d00d03` | `tools/agents/README.md` |
| `8936968d2e01eebfcf0596de85064042310d358dd12b3770f7c242814168dfc2` | `tools/agents/api/v1alpha1/BUILD.bazel` |
| `fc3ad622e8d8bbbe4f73915b9e42f2cc2e43a5f6ad46851a5518cfa576f08226` | `tools/agents/api/v1alpha1/canonical.go` |
| `2aaad7d438773082adb1ddf11fac991326532f0baebe1396184420f5ebfad3f9` | `tools/agents/api/v1alpha1/types.go` |
| `57078c7515a34975629baa18d6c051ec03aacd3f6ac43723d48377171aaccf69` | `tools/agents/api/v1alpha1/validation.go` |
| `727f0441ddb8bb1113e9b70d38c3b4c5bf85dd0277ff78092f12d9a6864a7652` | `tools/agents/api/v1alpha1/validation_test.go` |

Goal-resource projection changes are excluded from the source identity because
closing this attempt necessarily advances them again.

## Verification

- `bazel_agent test //projects/goal/... //tools/agents/api/v1alpha1:all //:buildifier_test`
  passed 8 of 8 owner-goal tests plus the contract and Buildifier tests.
- `bazel_agent build //tools/agents/... //projects/goal/internal/fsstore:all`
  passed for all six selected targets.
- `bazel_agent build --config=lint //tools/agents/api/v1alpha1:all //projects/goal/internal/fsstore:all`
  passed for all four selected targets and their configured lint aspects.
- `bazel_agent run //:gazelle -- tools/agents/api/v1alpha1` completed; its
  generated conventional test target was reviewed, and the duplicate manual
  target was removed.
- `git diff --check` passed after the final source changes.
- The checked-out `agent-system-phase-1` record validated after the goal-store
  repair, and the same CLI successfully published this attempt using a
  workspace-relative plan path.

## Criterion observations

- `shared-contracts`: pass for the shared contract slice. Strict fixtures cover
  deterministic operation and artifact JSON round trips, unknown fields,
  malformed identities, unknown effects, authority widening on every scoped
  axis, and incompatible information flow.
- `information-policy`: pass for the shared contract slice. Disclosure,
  build-consumer, publication, and public/secret/personal information axes are
  represented independently.
- `encountered-bug-policy`: pass. Two bounded goal-store defects were fixed
  with regression coverage instead of bypassed: missing empty `attempts/`
  after checkout and workspace-relative checkpoint artifact resolution.
- `exact-candidate-validation`: pass for the affected source listed above.
- The remaining Phase 1 criteria were not implemented or claimed by this
  attempt.
