## Why

GPU rentals need the funded alwaldend team and an independently verifiable
cheapest-offer rent-and-release workflow. Personal-account API keys cannot be
assumed to select an existing team. Credential owners and Terraform lifecycles
must remain explicit.

## What Changes

- Keep team identities and account settings in the account root.
- Add a separate alwaldend configuration root using its operator-provided
  team key, with its own backend, full-access project keys and Vault delivery.
- Read environment defaults from one configuration section in both scopes.
- Retire other teams and their test objects through their owning states;
  preserve tracked ownership when missing credentials block deletion.
- Keep one test contract in the flat rentals_test root and its existing test state.
- Select the cheapest on-demand offer with matching storage pricing and
  verify interactive creation followed immediately by release.

## Capabilities

### New Capabilities

- `team-gpu-rentals`: Independent team authentication, cheapest-offer rental and immediate cleanup.

### Modified Capabilities

None.

## Impact

Owns infra/vastai configuration, runtime targets, tests and documentation.
Retires the former provider-project live test root after destroying its objects.
Provider and SDK implementations are unchanged. No new external dependencies,
credit transfers, saved plans or automatic approvals are introduced. The user
explicitly authorized team retirement, Vault bootstrap, settings and live rental
validation. Two imported-team deletions remain blocked by missing management
credentials and are reported rather than detached.
