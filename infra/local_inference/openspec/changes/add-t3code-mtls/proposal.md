## Why

The local inference VM needs the repository's browser-based T3 Code interface
behind Traefik and client-certificate authentication. Remote access must not
enable cloud inference or hosted Codex tools.

## What Changes

- Deploy checksum-pinned standalone T3 Code through its shared Ansible role.
- Run T3 Code normally as its own service account on loopback.
- Add direct Traefik mTLS ingress and scoped Vault ACME enrollment.
- Declare the component's DNS delegation and hostname with existing providers.
- Permit HTTPS and ACME ingress while keeping application ports private.
- Preserve the local-only Codex provider and disabled remote tools. Configure
  Auto with local review and allow Full access, as requested.
- Remove the extra account-level egress firewall; Codex owns sandbox policy.

## Capabilities

### New Capabilities

- `t3code-ingress`: Browser access to local Codex through mTLS-protected T3 Code.

### Modified Capabilities

None. The original prepared host change remains a separate implementation record.

## Impact

Changes this component, its Vault PKI declarations, and the shared T3 Code
role. The role keeps its existing npm mode and gains an explicit standalone
mode. No infrastructure is deployed and no secret values are created.
