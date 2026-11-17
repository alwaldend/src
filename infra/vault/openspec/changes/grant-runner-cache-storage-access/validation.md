# Validation

Vault formatting and offline provider validation passed, along with coordinated
quality, semantic lint and strict specification checks. The truenas_admins Vault
group contains the NAS AppRole entity. The truenas_user Vault group derives the
other canonical AppRole entities, excluding NAS; NAS reads these two direct
entity sets to provision accounts and matching group membership.

The existing scoped NAS publisher and AppRole shared-read grants remain intact.
XCP-ng reads its own account through its existing policy. The added
development-account read policy and binding are removed; other AppRoles gain no
read of the separate NAS provisioning API credential. The authorized rollout
applied the three access additions and three updates. State persistence was
recovered through the [independent AppRole migration](../isolate-approle-state/validation.md)
without restarting or reloading Vault. All four migrated resource plans are unchanged.
