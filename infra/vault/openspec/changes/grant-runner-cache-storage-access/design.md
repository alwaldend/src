## Decisions

Permit the NAS AppRole to look up the canonical role groups/entities and publish
only NAS account credentials under shared/read/approle/+/nas. Existing AppRole
shared-read policies already restrict each role to its own credential path.
XCP-ng reads only its own NAS account for the storage connection.
Other AppRoles cannot read the NAS API administrator credential. The provisioning controller
uses the NAS identity in the shared and independently managed AppRole NAS roots.
VM provisioning does not request that identity or API key.

Vault owns truenas_admins and truenas_user membership as direct entity sets.
The NAS AppRole belongs to truenas_admins. truenas_user projects the canonical
AppRole groups excluding NAS; NAS provisioning reads only these two NAS access
groups. Account creation and NAS group membership follow their entity members,
including administrator precedence for an identity explicitly in both groups.

The existing self-scoped policy covers XCP-ng's own published NAS account.
Remove the added development-credential policy and its AppRole binding.
