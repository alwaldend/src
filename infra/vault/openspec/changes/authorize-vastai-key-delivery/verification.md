## Authorized deployment

On 2026-11-09, the user authorized both the team writer's exact-entry permission and the human login's own-namespace access. The interactive owning Vault apply displayed exactly one policy addition and two existing policy attachments, with no destruction. That reviewed plan was applied: 1 added, 2 changed, 0 destroyed. No mount settings or unrelated resources were changed. This was scoped deployment, not a full Vault reconciliation.

## Live acceptance

- The existing human login can create/read/update/delete/list/patch within its own namespace. The authenticated UI mount catalog now includes secrets/. Listing its own namespace succeeds and includes imagegen/.
- The human login remains denied on the infrastructure team's issuer-secret path.
- A fresh team-role login has create/read/update/delete only on the generated-token data entry and read on its exact metadata entry. It remains denied on the sibling imagegen/web secret and the user's imagegen Terraform backend.
- The existing generated-token entry remains version 1. No secret value was printed or rotated.
- The alwaldend root, using only its own Vast and Vault identities, refreshed the existing key 36816419 and the same user KV entry and reported no changes. A reviewed refresh-only apply saved the default-provider association with 0 remote resources added, changed, or destroyed.

Ignored out/imagegen/permission-access-evidence.json contains the sanitized repeatable capability, mount-discovery, folder-listing, and metadata results. Source checks cover default-false member access for other roles. No GPU was rented.

## Navigation correction

The initial acceptance checked direct namespace listings but missed the UI's path from the mount root. The user reported an empty secrets view. Fresh checks confirmed deny on secrets/metadata and on the namespace root without a trailing slash, while the imagegen folder and token entry already existed.

The already-authorized personal UI access deployment was completed with exact list-only metadata navigation rules. A read-only plan included the default-false Vast component namespace policy as a negative control: only the personal policy changed. The fresh interactive owning apply was reviewed and applied with 0 added, 1 changed, 0 destroyed; its only target was module.users_simeonwarren_approle.vault_policy.approle_secrets.

Live acceptance passed 12 root-to-leaf navigation checks: capabilities and actual LIST calls on six parents, each with and without a trailing slash. Each listing contained the next expected folder or token entry. Six value-isolation checks confirmed personal token access, denial on other namespace values and metadata, team-writer token access, and denial on the web password and user backend. The generated-token entry remains version 1. Listings and values were not printed or retained; ignored out/imagegen/ui-navigation-evidence.json stores only sanitized results. Browser rendering and a paid GPU boot were not exercised.

Stable defect: VAULT-KV-ANCESTOR-NAVIGATION. A successful leaf listing is insufficient acceptance for UI navigation; verify the complete mount-root path walk. The correction reuses the policy owner and existing change records, with no duplicate policy, new mount, or shared authentication change.
