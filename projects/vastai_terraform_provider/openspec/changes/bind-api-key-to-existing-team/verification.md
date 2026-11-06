## Offline evidence

The Terraform CLI fixture failed before implementation because `team_id` was
unsupported. A second pre-fix run demonstrated that empty permission strings
were rejected. After implementation, success, rejected activation, missing
ownership, and scoped-input rejection scenarios pass. The fixture requires omitted permissions for default
full access, rejects issuer-context changes, and verifies generated-key cleanup.

## Live evidence, 2026-11-06

The configured infra/vastai Vault key created full-access key 36348922 through
the reviewed owning credential apply. Activation of that new key for alwaldend
returned HTTP 400 from PUT /users/save-context/. The key remained in state and
the dependent Vault entry was not created. A reviewed credential-root destroy
revoked that exact key with one deletion and no other changes. The shared
issuer, existing team, credit, and compute were not changed.

Binding remains unresolved. Fixture success establishes provider failure
handling, not compatibility with live existing-team activation.

## Provider authentication and module reuse

The provider-only E2E cases omit the resource credential and supply an explicit
provider key while the environment fallback is unusable. Lifecycle and cleanup
pass. Migration removes a legacy issuer input, verifies access through the
provider, and moves the existing key into the actual shared module with no key
creation or revocation during those steps. A denied issuer update retains the
previous state and key until a valid provider is restored. The same-apply
bootstrap tests remain covered. A separate case verifies that a new explicit
bootstrap issuer still plans replacement; only removal of redundant issuer
metadata is a state-only update.

Actual-module cases cover existing-team success, rejected activation, invalid
team IDs, explicit provider aliases, and key-only cleanup. The managed-module
upgrade fixture retains the team and both key identities, records zero API
mutations during upgrade, and produces `team-module-upgrade.json`.
