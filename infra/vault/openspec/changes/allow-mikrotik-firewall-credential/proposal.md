## Why

The printer firewall workflow needs RouterOS credentials. An existing non-DNS
credential is already maintained by ingress, so reuse it without copying secrets
or giving the MikroTik workflow the broader ingress identity.

## What Changes

- Grant `src_infra_mikrotik` read access to the existing ingress RouterOS secret.
- Add the global view to its existing DNS credential policy for public printer DNS.
- Preserve existing DNS policies, AppRole membership, and state access.
- Keep policy deployment separate from firewall deployment and source validation.

## Capabilities

### New Capabilities

- `mikrotik-firewall-access`: Named credential reads for the printer firewall and global DNS.

### Modified Capabilities

None.

## Impact

Only the existing MikroTik AppRole module and this owner's change record are
modified. The paired printer change is maintained under
`infra/mikrotik/openspec/changes/allow-wifi-printer-access`. No secret value,
new account, dependency, wildcard grant, or live mutation is introduced by source changes.
