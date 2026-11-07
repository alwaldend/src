## Context

The MikroTik AppRole already owns state and receives scoped DNS credential
access. Its own secret subtree contains no RouterOS credential. A read-only
listing confirmed that ingress already has a `mikrotik` secret; its maintained
injector expects `mikrotik_username` and `mikrotik_password`.

## Goals / Non-Goals

Reuse the existing non-DNS RouterOS account while retaining the MikroTik identity.
Do not copy secret values, replace host authentication, or join a broad AppRole group.

## Decisions

**Verdict: proceed with an explicit exact-secret read policy.** This follows the
existing named shared-secret access pattern used for DNS. Sharing the whole
ingress AppRole would expose unrelated privileges. Provisioning a separate
RouterOS account is unnecessary for the user's requested reuse.

Add one policy to the existing MikroTik module and append it to its current
policy list. The paired firewall injector reads the original ingress path.
The account remains owned by ingress; rotation continues at that one location.

For the user's global DNS requirement, enable the `global` view alongside
`dc1` in the existing DNS access module. This adds read on the named Cloudflare
DNS secret through the existing shared policy owner without duplicating its
path or broadening access to other secrets.

## Risks / Trade-offs

- The shared RouterOS account couples the two workflows' rotation and router
  permissions. The policy grants access to that credential, not server-side
  per-rule privilege separation. The printer rule itself remains host-specific.
- Before code changes, reject these failure modes in review: wildcard path,
  write capability, replacement of existing policies, switching AppRole identity,
  mismatched secret field names, or storing a copied credential.

## Migration Plan

Validate and review the narrow policy change. With explicit live approval, use
the owning Vault Terraform apply scoped to `module.src_infra_mikrotik` and review
its interactive plan. Then run global DNS and firewall plans through the MikroTik identity,
review and apply the single printer rule under explicit authorization, and verify
from Wi-Fi. No Vault policy or firewall apply is authorized by this record.
Rollback the read grant through Vault Terraform after removing the firewall rule;
existing issued tokens may retain policy names until their normal expiry.
