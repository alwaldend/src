## Current identities and scopes

The existing src_infra_vastai role generates and delivers its team token using its own credentials. Its delivery policy grants create/read/update/delete on exactly secrets/data/alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai, and read on that entry's exact metadata path. It grants no wildcard, sibling-secret, user-backend, or policy-management access.

The personal role explicitly enables allow_member_secret_access. Its existing human members receive the already-defined personal namespace policy through the existing member group. Other roles retain the default false setting. No duplicate personal policy, additional group, AppRole, or mount is created.

## Decision

The previous recipient-token lookup and aliased user provider made the team root depend on the recipient's authentication. The team role initially reported deny on the destination; the human login also reported deny on its own secret and did not discover the KV mount. The approved correction grants the team writer only its delivery entry and reuses the personal namespace policy for the human login.

## Deployment scope

The user explicitly authorized both scoped permissions. The owning targeted apply was limited to vault_policy.vastai_imagegen_delivery, module.src_infra_vastai_approle.vault_approle_auth_backend_role.role, and module.users_simeonwarren_approle.vault_identity_group.group. The completed deployment and live acceptance are recorded in verification.md. No full Vault reconciliation or GPU rental is part of this change.

The personal policy subsequently received exact list-only metadata ancestor rules to complete normal UI navigation. That correction used a separate reviewed apply targeting only the existing personal namespace policy. Ancestor folder names are visible; values and metadata reads outside the user's namespace remain denied. The full root-to-leaf path walk replaces the earlier leaf-only navigation acceptance.
