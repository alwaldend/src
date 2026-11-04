## Decisions

The personal credential remains in this component's existing injector. The
bootstrap root creates a test_ team with transfer_credit fixed to zero.
The provider generates an independent management credential, selects its team
context without changing the personal key, and generates a restricted project
key. Terraform writes the latter to the project's own Vault entry, consumed
by that project's AppRole and injector. All state is Vault-backed; generated
keys are sensitive and never exposed as outputs.

## Risks / Trade-offs

Team owner credentials are not retrievable from API-key list responses.
Bootstrap must generate and retain its own management credential in protected
state. Missing or uncertain mutation results require diagnosis, never retries.
The project role must be deployed before token delivery. Vault's provider uses
that role's standard namespace permissions; the personal-key injector and state
backend retain infra/vastai's identity. Soft deletion uses existing data-path
permissions and does not require metadata deletion.

## Decision Review

Proceed with two invocation-scoped Vault identities and no cross-role policy.
The operator already has bootstrap membership in both roles. Granting an extra
policy to the infrastructure role is unnecessary when the Vault provider selects
the destination's owning role directly. A checked-in team/key catalog calls the
reusable projects/tf_modules/vastai_team module; secret resources remain in this
root with its explicitly selected provider identity.

The catalog is plain `teams.json`, loaded with jsondecode(file(...)); main.tf
owns resource composition. The user requested removal of moved.tf. No live
bootstrap apply has been performed in this task.

## Rebase validation

Observed 2026-11-04 after SDK PR #251 and provider PR #255 merged. The remaining
integration preserves the merged provider packages, typed dependency payload,
reviewed lifecycle fixes, and fixture source. Offline validation passes all 38
selected tests, including the three Terraform format checks, owner OpenSpec
validators, repository quality, and loopback Terraform CLI acceptance. Bootstrap,
Vault, test-root, and provider launchers and the shared site build successfully.
The bootstrap runtime includes teams.json and the shared team module; no live
plan, apply, token generation, Vault write, credit transfer, or rental is run.
Git publication identity and review state belong to the delivery receipts.
