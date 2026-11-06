## Current source

The account root owns team identities and account settings. The separate
alwaldend root owns team settings and full-access keys using its operator-provided
seed. Both roots read the environment map from config.json.

The flat infra/vastai/rentals_test root owns one contract in the existing
rental_test backend. It uses VAST_API_KEY directly, has separate interactive
rent and release commands, and contains no child module or lifecycle runner.
The unused ordinary rental backend selector and TF_VAR credential injection
are removed. No temporary import or moved blocks remain.

## Evidence

[Live acceptance](validation.md) records the preceding implementation's completed
rent-and-release tests, CLI verification and remaining imported-team deletions.
The flat revision does not require another paid rental. Offline validation loads
its actual root with the packaged provider; the provider acceptance suite owns
local API lifecycle coverage. Candidate check results and publication identity
belong to delivery receipts.

## Scope and limits

Two imported teams still need management credentials for deletion. They remain
tracked in account state; no unchanged deletion request was retried during this
source simplification. The operator must run release after rent or a failed
creation apply and verify empty state and contract absence through the CLI.

## Session ergonomics

`vastai-rental-overengineering` (configured): the generic rental module, duplicate
backends and command runner exceeded this
task's single-contract acceptance need. The user's flat-root correction removes
those layers and keeps the existing state and credential references.
