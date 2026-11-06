## Live acceptance before the flat-root revision, 2026-11-06

The branch was rebased onto merged master 33e7a8f. The operator-provided Vault
key selects funded alwaldend account 846045. The CLI verified all nine team
permission categories, positive credit and an initially empty instance inventory.

The separate alwaldend root applied at 19:19 UTC: five additions, no changes or
deletions. These are three environment settings, rental key 36371877, and its
Vault entry. The generated key independently selects alwaldend with full team
permissions. All three variable names are present. API reads mask values; the
declarations and accepted writes establish their configured true values.

The former provider-test root destroyed its reader key, environment marker and
role at 19:16 UTC. The account root then deleted test_vastai_provider, its
project and management credentials, and its project Vault entry at 19:20 UTC.
Deletion of projects_vastai_terraform_prov (846285) and
test_vastai_terraform_provider (846772) failed because imported identities have
no recovered management credentials. Both remain tracked for deletion.
Alwaldend and the personal account key were preserved.

The account's alwaldend identity and three variable addresses were reconciled
within the same state. Read-only state listing confirmed all four new addresses;
completed moved blocks were removed immediately afterward. No temporary import
or moved blocks remain in final source and no resource changed backend.

The first independent rental test succeeded: offer 56954301, RTX 3060 Ti,
contract 59459593. Its reviewed price was USD 0.026851851851851852 per hour for
a requested 1 GB allocation. Creation completed at 19:21:56 UTC; release began
at 19:21:57 UTC and deletion completed at 19:22:15 UTC. The receipt confirms empty
test state and contract absence. The server reported 3.5 GB during release,
so the default allocation and offer comparison were corrected to 3.5 GB.
A separate CLI search with the same predicates and allocation returned the
cheapest offer first at USD 0.02753214814814815 per hour at that observation.
The final candidate test also succeeded with 3.5 GB: offer 56954110 and
contract 59460402 at USD 0.02753214814814815 per hour, matching the independent
CLI's cheapest offer. Creation completed at 19:28:59 UTC; release began at
19:29:00 UTC and deletion completed at 19:29:20 UTC. The release refresh still
reported 3.5 GB, confirming allocation and price-comparison alignment. The final
receipt confirms empty state and contract absence at 19:29:22 UTC.

The final account plan without temporary blocks reports zero additions, zero
changes and only the two tracked imported-team deletions. It preserves alwaldend
and all three account settings. No further deletion retry was made without new
credential evidence.

## Validation and ergonomics

The actual account/team root fixture passed before live deployment. It uses
separate Terraform working directories, verifies full-access key configuration
and shared defaults, and rejects invalid names. Rent/release retains its own
command fixtures and provider acceptance coverage in that revision.

Observed defect: fixture-only context switching was mistaken for a supported
live API contract. Two issuer approaches failed; one generated key was revoked
and the second request created no key. Those speculative changes were removed.
The pinned CLI's supported commands should have been inspected first. Current
component documentation now identifies each credential owner, backend and CLI
entry point, including the one-time operator bootstrap.

A targeted destroy initially included an unrelated Vault secret and was
cancelled. Automatic approval review later referenced that cancelled plan and
rejected confirmation. A fresh full-catalog plan proved the failed-key-only
scope; the subsequent reviewed cleanup removed exactly that key. No unrelated
secret was deleted during that cleanup. The later test-secret deletion was
separately authorized as part of retiring its team.

Final offline checks passed: 35 tests, including repository quality, component
fixtures, the full provider Terraform acceptance suite and both owner OpenSpec
validations. Affected semantic lint passed. Local review identified one fixture
receipt field claiming rental-state separation without loading the rental root;
that claim was removed. The account/team fixture now also changes one shared
value and verifies both roots consume it, preventing hardcoded-value duplication.

The final team-root plan at 19:36 UTC reported no changes. Independent post-test
CLI reads confirmed the personal account context and operator keys are preserved,
retired keys are absent, and the alwaldend instance inventory is empty.

## Flat-root revision, 2026-11-06

The root is now infra/vastai/rentals_test, with separate tf.rent and tf.release
commands and the same rental_test backend. No child module, generic rental state,
credential variable or Go lifecycle runner remains. No new paid rental was run.

At 20:00 UTC, its owning state listing was empty. At 20:01 UTC, the read-only
creation plan selected offer 49574214, Quadro P2000, for 3.5 GB at USD
0.027962962962962964/hour: one addition, no changes or deletions. The subsequent
destroy plan reported no changes and performed no marketplace read. No plan was
saved. The offline catalog fixture validates the actual flat root with the
packaged provider; provider acceptance retains local API lifecycle coverage.
