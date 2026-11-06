## Live reconciliation, 2026-11-06

### Current status

The user selected `true` for ENABLE_HTTPS. Its existing identifier was imported
successfully at `vastai_environment_variable.enable_https` at approximately
11:49 UTC. The declaration sets the string value `true`.

After the user explicitly approved the live update, the exact reviewed saved
plan applied at approximately 12:11 UTC: zero additions, one in-place value
update, and zero deletions or replacements. A fresh plan returned exit 0 with
`-detailed-exitcode` and reported no changes. ENABLE_HTTPS is now reconciled.

All three additional owned teams below are now imported at
`vastai_team.teams`, with null management credentials and destruction guards.
The first CLI import succeeded; the second encountered Terraform's partial
for_each postcondition evaluation error before saving that instance. A fresh
reviewed saved plan imported the remaining two identities: two imports, zero
additions, zero changes, and zero deletions. No team or credential was created.
Import blocks were then removed as requested. A fresh plan returned exit 0
with `-detailed-exitcode` and reported no changes. State inspection verified
all three imported account IDs and ENABLE_HTTPS equal to `true`.
The final configuration declares only team names, with numeric IDs retained
in state. Removing the hardcoded IDs and their postcondition preserved the
resource addresses and destruction guards; a fresh plan again reported no changes.
Team configuration now uses the name-keyed `teams` map with an initially empty
object per team. A reviewed moved block migrated all three resources from
`vastai_team.existing` to `vastai_team.teams` with zero live resource changes.
The post-migration plan reported no changes and retained bootstrap module addresses.
All four team objects now live in teams.json and load through the single
`teams` local. The separate `bootstrap_teams` local is removed. Optional
api_keys fields retain the credential-generating module path; other entries
retain their team resources. The unified catalog's live plan reported no changes,
including the existing project key and Vault secret wiring.

The user requested direct revocation of unmanaged team keys while preserving
account keys. Key 36032627 was revoked through the authenticated component CLI.
Before/after inventories verified that it was the only removed ID. Primary
key 36011374, active account key 36036521, all owner context records, and all
Terraform-owned keys were preserved. No Vault credential was changed.

Declared resources are reconciled. Team-scoped inventories for the three
imported teams remain unavailable without their management credentials, so
this observation does not establish that their contents are empty. API masking
also prevents detecting later external changes to environment-variable values.

### Initial inventory

Read-only observations used the component CLI and owning Terraform wrappers.
Both roots returned exit 0 with `-detailed-exitcode` and reported no changes.
At that initial observation, no import, apply, credential generation, credit
transfer, or rental had run.

### Existing Terraform ownership

| Root                                           | Owned resources                                                                                        | Live plan  |
| ---------------------------------------------- | ------------------------------------------------------------------------------------------------------ | ---------- |
| `infra/vastai/tf`                              | Team `test_vastai_provider` (846810), management key 36078045, project key 36078048, Vault token entry | No changes |
| `projects/vastai_terraform_provider/test/team` | Reader key 36079589, `TEST_VASTAI_TERRAFORM_PROVIDER`, role `test_vastai_terraform_provider`           | No changes |

The management key is part of the team resource, not a separate Terraform
address. The team's automatic owner-key record also does not require a second
resource. The account's primary key is an account credential, not a generated
project key.

### Imported existing teams

API-key metadata identifies these additional teams, now declared separately
from the credential-generating bootstrap catalog. Their contents have not been
inventoried because their team credentials are unavailable through the current
authentication contexts.

| Team                             | Team ID | Automatic owner-key ID |
| -------------------------------- | ------- | ---------------------- |
| `alwaldend`                      | 846045  | 36015937               |
| `projects_vastai_terraform_prov` | 846285  | 36036241               |
| `test_vastai_terraform_provider` | 846772  | 36074058               |

The manual `infra/vastai` team key 36032627 associated with team 846045 was
revoked. The account key with the same name, 36036521, remains operator-owned
as requested. API-key detail reads returned HTTP 410. The provider now imports
team identities without recovering keys; generated API-key import remains
unsupported. Automatic owner context records represent owned teams and were
preserved rather than treated as additional manual API keys.

The personal account also contains `ENABLE_HTTPS`. Both the CLI's explicit
value-read option and the underlying API response preserve its mask. Its
adoption used the user's value selection recorded above.

### Coverage and limitations

- Personal-account instance, SSH-key, and volume collections were empty.
- The initial inventory returned ten API-key metadata records. Exactly one
  unmanaged team key was subsequently removed, leaving nine records.
- Account-context role listing returned HTTP 401 for missing `team_read`;
  this is not evidence of an empty role inventory.
- The upstream subaccount command returned HTTP 400 because it supplied an
  unsupported `owner` input; it did not establish a complete team listing.
- The existing team and key values in protected Terraform state were retained.
  No state was fabricated, copied between roots, or force-unlocked.

Raw responses, state, saved plans, and credential-bearing diagnostics are
temporary private scratch, not checked-in evidence. Resource and key IDs in
this record are identifiers, not credential values.

## Source validation

The provider's full Terraform CLI acceptance suite passes. Seven team-import
scenario artifacts record GET-only API traces, including valid and invalid
imports, refused destruction, external absence, and refresh failures. Existing
creation and credential-cleanup scenarios remain covered. Terraform formatting,
both owner OpenSpec validations, repository quality, and affected semantic lint
are delivery gates. Live acceptance is the verified imported state and final
no-change plan after removing import blocks.

## Temporary migration cleanup

On 2026-11-06 at 17:25 UTC, the owning `//infra/vastai/tf:tf.plan`
with `-detailed-exitcode` returned exit 0 and no changes after removal of the
completed `vastai_team.existing` to `vastai_team.teams` moved block. Refresh
confirmed all three imported addresses and the existing singleton test-team
module address. No import or moved blocks remain in this root; no apply or
state-changing operation was run for this cleanup. The import-capable shared
module preserves the existing managed-team and project-key addresses.
