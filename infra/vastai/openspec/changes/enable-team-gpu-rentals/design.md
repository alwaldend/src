## Context

The provider supports offers, rental contracts, team identities and project keys.
The existing infrastructure root owns the team catalog and imported alwaldend
identity. The user requires that catalog's shared team module to create its
rental key, and a cheapest-offer rent-and-release test using the funded team.

## Decisions

Proceed with one team/key catalog in `infra/vastai/tf`. Its alwaldend entry calls
`projects/tf_modules/vastai_team` and declares the full-access rental key. The
root stores sensitive key outputs directly in Vault. Select the existing
project or component Vault AppRole through provider aliases and the catalog's
`vault_auth` value. The existing project-secret resource address remains stable.
AL injects provider credentials; resources receive no issuer or Vault token input.

The separate credentials root is removed. It added a backend and cross-state
ownership transfer without an independent lifecycle requirement. Keeping that
root or duplicating the shared module's key resource would preserve unnecessary
owners. The module already owns team and key composition; the root owns Vault
identity and destination selection.

Alwaldend's imported identity moves from `vastai_team.teams["alwaldend"]` to
`module.teams["alwaldend"].vastai_team.team` within the same HTTP backend before
key apply. This migration is pending. Preserve the funded team and its account
ID, reject team creation or replacement, and remove temporary migration blocks
after reconciliation. No import or moved blocks remain in final source. The
retired credential backend owns no key after the reviewed failed-key cleanup.

## Rental lifecycle

A local child module receives an aliased alwaldend provider, with its dedicated
team key injected from the component's Vault path. It never changes credential
context. Ordinary rentals and the explicit test use distinct Vault-backed states.

Search rentable, unoccupied on-demand GPU offers, order by `dph_total`, and take
one result. Use a minimal disk allocation and a configurable price ceiling,
without GPU-model or verification filters. Preserve an existing contract's offer
ID. Release skips marketplace discovery and targets only the owned contract.

The live test uses reviewed interactive creation and destruction applies without
saved plans or automatic approval. It starts release immediately after creation
returns, including after failure, and does not wait for boot or SSH. Cancellation
propagates through the real launcher, then cleanup uses an independent deadline.
A sanitized receipt records observed lifecycle and cleanup results.

## Acceptance and limitations

The catalog fixture must first fail against the old root, then verify actual
team-module composition and distinct mocked Vault provider identities. Duplicate
secret destinations and unknown identities must fail before provisioning. Native
Terraform tests use only mocks; provider and command fixtures retain their own
behavioral evidence and repeatable receipts.

The live personal issuer creates keys, but selecting alwaldend on the generated
full-access key returned HTTP 400. The key remained tracked and was revoked
through its owning workflow. The refactor does not establish binding acceptance
or justify repeating the unchanged failed request. After resolving binding,
verify the stored key's team, full permission categories and credit before the
authorized cheapest-offer rental test.

Offers can disappear, hard termination or backend outages can interrupt cleanup,
and interactive review adds a short delay before destruction. Preserve owning
state and recovery commands; never retry creation after an ambiguous result.
