## Context

The bootstrap root owns a generated team, project key, and Vault token entry.
The provider test root separately owns its reader key, variable, and role.
The personal account also exposes three legacy teams, manual keys, and the
API-masked ENABLE_HTTPS variable. See proposal.md for the reconciliation scope.

## Decisions

Add import to the existing command map, preserving its AL wrapper and declared
runtime inputs. Use the existing state and resource addresses; do not duplicate
test-root ownership or manufacture Terraform state.

Keep all team configuration in the single teams.json catalog. The provider-owned
importer adopts numeric account IDs without
credentials; declare those teams by name with prevent_destroy. State owns
their imported account IDs; do not duplicate them in configuration or add
account-ID postconditions. Use a name-keyed map of objects, initially empty,
so per-team configuration can be added without changing resource addresses.
Use `teams` for the map and resource key; a moved block migrates the previous
resource address without recreating teams. Entries with an api_keys field
retain their team/key module; entries without it use the direct team resource.
Secret targets derive from the optional key maps in that same catalog. Preserve
the field's presence for existing entries to retain their owning resource path.
Remove the import blocks after successful adoption as the
user requested. Do not copy a masked environment value into source or infer
it from the variable's name.

## Risks / Trade-offs

- Unrecoverable team credentials limit imported teams to identity tracking.
  Refuse destruction and document unavailable team-scoped inventory.
- ENABLE_HTTPS cannot be reproduced from a masked API response. Obtain its
  intended value before declaring or importing it. The user selected `true`;
  import the existing identifier and apply only the reviewed value update.
  Masked reads cannot detect later external changes to that value.
- Account-level inventory does not enumerate resources inside inaccessible
  legacy teams. Report that coverage limit rather than claim full convergence.

## Migration Plan

Verify both existing roots with read-only plans. Expose and inspect the import
launcher. After the environment value is established, declare the variable,
import its canonical name, and review a fresh plan before any resource updates.
No compute, team, key, or secret mutations are part of the preparatory work.

## Team-key cleanup exception

The user explicitly requested revoking team keys outside IaC and preserving
account keys. The only unmanaged team API key in the observed inventory is
36032627, associated with team 846045. Revoke that ID through the component
CLI, then verify its absence and preservation of all other IDs. This is the
user-authorized imperative cleanup exception; no permanent cleanup command or
replacement credential is introduced into IaC. Preserve account keys 36011374
and 36036521, all Terraform-owned keys, and automatic owner context records.
The earlier primary-key credential request is withdrawn; authentication is
unchanged. Existing teams are adopted through the provider importer instead.
