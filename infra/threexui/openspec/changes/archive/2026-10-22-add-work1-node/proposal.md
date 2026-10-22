## Why

Add the requested work1 data-plane node hosted by njal.la.

## What Changes

- Configure work1 through the existing Terraform node module and Ansible inventory.
- Declare its supplied IPv4 and IPv6 addresses in the owned DNS configuration.
- Document the existing Vault maps and deployment prerequisites.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-threexui`: Add work1 to the configured nodes.

## Impact

Desired state uses the existing providers, modules, relay selection, and secret
references. Architecture documentation is delivered independently in the
move-threexui-architecture change.
