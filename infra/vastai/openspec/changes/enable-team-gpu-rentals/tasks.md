## 1. Authentication and independent ownership

- [x] 1.1 Rebase onto merged master and verify the personal account and operator-provided full-access alwaldend key through the pinned CLI.
- [x] 1.2 Keep team creation and account settings in the account root; create a separate alwaldend root and backend for team keys and settings, using provider authentication.
- [x] 1.3 Store the generated full-access rental key in Vault separately from its operator-owned seed and verify its team, permissions and positive credit.
- [x] 1.4 Define ENABLE_HTTPS, WEB_ENABLE_HTTPS and WEB_ENABLE_AUTH once in config.json and apply the same true values to the account and alwaldend.

## 2. Retirement and reconciliation

- [x] 2.1 Destroy the old live test root's key, marker and role before revoking its credential; retire its source while preserving offline acceptance tests.
- [x] 2.2 Delete test_vastai_provider, its generated keys and its Vault entry through their owning state.
- [ ] 2.3 Delete projects_vastai_terraform_prov and test_vastai_terraform_provider. Their imported state has no recovered management credentials; provider deletion was refused and both remain tracked.
- [x] 2.4 Reconcile alwaldend and the account variable addresses within the existing state, verify their new addresses, and remove all temporary migration blocks.

## 3. Independent rental lifecycle

- [x] 3.1 Keep the flat rentals_test root in its existing independent test state with separate interactive rent/release targets; never put contracts into team configuration state.
- [x] 3.2 Select the cheapest rentable, unoccupied on-demand GPU offer with storage pricing and a USD 1/hour ceiling, without model or verification filters.
- [x] 3.3 Execute the initial rent-and-release test, identify the API's 3.5 GB allocation floor, and align default allocation and price comparison.
- [x] 3.4 Independently compare the cheapest offer through the CLI and execute the final 3.5 GB rent-and-release test. Verify contract absence and empty test state in its receipt.
- [x] 3.5 Flatten the root, remove the GPU child module, duplicate ordinary rental backend, credential variable and Go lifecycle runner.
- [x] 3.6 Verify the flat rental root and empty test state without another paid rental.

## 4. Validation and delivery

- [x] 4.1 Before implementation, exercise separate account/team configuration, shared defaults, full-access keys and invalid key names with actual Terraform roots and mock providers.
- [x] 4.2 Pass formatting, repository quality, affected semantic lint, owner OpenSpec checks and package tests against the final candidate.
- [x] 4.3 Complete local review and prepare delivery inputs, reporting the two remaining deletions. AL receipts own publication and remote-feedback status.
