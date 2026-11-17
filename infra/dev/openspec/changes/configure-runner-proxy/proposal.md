## Why

Secure1 cannot fetch several external build dependencies without VPN egress.
The user requested reuse of the previous host HTTP proxy and direct access to
`*.alwaldend.com`.

## What Changes

- Reuse the former T3 Code proxy account through its existing Vault reference.
- Supply proxy environment variables to runner jobs, with local-domain bypass.
- Forward a loopback HTTP endpoint over verified TLS to the existing proxy.
- Install the missing runner `zip` prerequisite through Ansible.

## Capabilities

### New Capabilities

- `runner-proxy`: declared external egress and local bypass for secure1 jobs.

### Modified Capabilities

None.

## Impact

Only secure1's Ansible configuration and the development deployment's exact
Vault read grant. No new external software, proxy account or remote proxy
deployment. Live policy application and host deployment remain separate steps.
