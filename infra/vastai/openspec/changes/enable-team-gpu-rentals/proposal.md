## Why

GPU rental support exists in the provider but has no infrastructure root using
the funded `alwaldend` team. The operator requested a cheapest-offer rental
test followed immediately by release.

## What Changes

- Use a dedicated alwaldend team key from Vault with an aliased
  provider passed into a child rental module; perform no context switching.
- Declare the full-access key in the existing team catalog and deliver it to Vault through the component identity.
- Reuse the shared team module for that credential, using provider-owned issuer authentication.
- Add separate Vault-backed rental state and cheapest on-demand offer selection.
- Provide an explicit interactive rent-and-release test with a sanitized receipt.
- Replace saved-plan-only apply commands in the affected Vast.ai roots with
  interactive Terraform apply.

## Capabilities

### New Capabilities

- `team-gpu-rentals`: Team authentication, cheapest-offer rental, and immediate cleanup.

### Modified Capabilities

None.

## Impact

Owns infra/vastai configuration, runtime targets, and documentation. Consumes
the provider and imported-team module from merged PR #267 for full-access key
generation and independent binding; the SDK is unchanged. No new external
dependencies or credit transfers are required. The user authorized rental
configuration and one cheapest-offer live rent-and-release test.
