## Decisions

Keep an explicit selected-user set containing the existing bot userpass alias.
Generate one stable 72-character password and one KV entry per selected user.
The logical path is `secrets/alwaldend.com/users/<username>/userpass`.
A small policy template grants the infrastructure role data lifecycle and
metadata read/update for those entries only; metadata deletion is unnecessary.
Vault provider 5.8.0 needs those metadata operations during refresh and updates.

Grant the existing simeonwarren entity exact data/metadata reads for the bot
userpass entry. Do not assign any storage policy to users or their group.
Default-deny provides the requested privacy; no explicit deny rules are needed.
Remove the earlier shared own-namespace policy and its users-group assignment.

Selected login resources consume passwords from successfully stored KV entries.
Their selection references existing aliases so entity binding precedes login
management. Keep existing identity memberships, direct token policies, and token lifetime
settings. Secrets remain in KV and Vault-backed
Terraform state; do not export state, save plans, or enable debug logging.

## Verification

Use disposable local Vault with exact source inputs. Check that only the
selected user has a generated entry and its password has length 72; bot tokens
cannot read/write their own data or metadata; no user/group storage grants or
deny rules exist; the human can read only the bot userpass entry and cannot
write, list the bot namespace, or read siblings. Exercise metadata refresh,
rotation, unchanged reconciliation, and entity binding. The local server must
use -dev-no-store-token and preserve the shared helper; clean fixture state.

## Deployment

The user explicitly authorized the live apply on 2026-11-16 after the earlier
automatic approval rejection. Two interactive applies deployed the selected-entry
writer/reader policies and existing human/role assignments, then used a fresh
role token to store the credential before updating the existing login. Reviewed
plans contained only the authorized resources, with no destroys. Earlier designs
were never applied and required no namespace or credential migration.

The subsequent full plan reported no changes. Live checks verified the human's
exact read access, the bot's original identity and new password, and absent
bot access to its userpass entry. The temporary verification login was revoked
without changing the shared host token helper. See verification.md for scope
and results.
