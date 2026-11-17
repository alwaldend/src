# Validation

Vault formatting and offline provider validation passed, along with coordinated
quality, semantic lint and strict specification checks. The truenas_admins Vault
group contains the NAS AppRole entity. The truenas_user Vault group derives the
other canonical AppRole entities, excluding NAS; NAS reads these two direct
entity sets to provision accounts and matching group membership.

The existing scoped NAS publisher and AppRole shared-read grants remain intact.
XCP-ng reads its own account through its existing policy. The added
development-account read policy and binding are removed; other AppRoles gain no
read of the separate NAS provisioning API credential. These are source-policy
checks; no live Vault group, grant or secret was changed.
