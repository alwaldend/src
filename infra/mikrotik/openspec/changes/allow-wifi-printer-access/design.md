## Context

See [proposal.md](proposal.md) for motivation. The recorded router1 topology has
wired `bridge1` and wireless `bridge2`; its terminal `forward` drop blocks Wi-Fi
access to wired devices. The existing `tf/` root owns DNS only. Exports are
historical documentation and cannot deploy an exception.

## Goals / Non-Goals

Provide a deployable exception without taking ownership of the rest of the
firewall. Do not change ingress, authentication, printer firmware, DHCP,
multicast discovery, or print applications. The caller supplies the printer's
address; it must remain assigned to that device.

## Decisions

**Verdict: proceed with a separate firewall root and one host-specific rule.**
Adding Wi-Fi to the general LAN forwarding list would expose unrelated wired
devices. Editing the export would not deploy anything. A separate root retains
the DNS credential boundary and limits apply to the firewall exception.

Read the printer address from the canonical DNS document instead of maintaining
a second independent address. Match wireless source subnet and input interface,
wired output interface, and a single IPv4 host destination. Use its bare address
to match RouterOS's canonical representation. Permit all IPv4 protocols to that
one device so its web interface and native printing protocols remain available.
The user requested device connectivity rather than a particular print protocol.

Find the enabled forward drop by chain, action, and its recorded comment
`drop forward`; fail unless exactly one matches. Use its live identifier for
`place_before`, retaining existing invalid-connection and WAN protections.
The pinned RouterOS provider documents this insertion mechanism in
`docs/resources/ip_firewall_filter.md` and `docs/data-sources/ip_firewall.md`
at tag `v1.99.1`.

Keep the rule lookup inline with its one consumer. A module existed only to
work around Terraform mock handling of this provider's nested response; that
test-only boundary did not justify production indirection. Remove the wrapper
and mock suite, retaining schema validation and the resource precondition.
Review the actual provider plan before deployment; offline validation does not
exercise absent or duplicate live query results.

Publish the same printer A record in both views, as requested by the user.
Enable Cloudflare injection and grant the existing DNS policy read access to
its named Cloudflare secret. Public resolution of the private address does not
provide a route to it. The local record remains in place.

Use the existing owner AppRole and a `tf_backend/firewall` path. Following the
user's request to reuse an existing credential, reference ingress's existing
non-DNS RouterOS secret with `mikrotik_username` and `mikrotik_password`.
The paired Vault change grants read on that exact secret to the MikroTik
identity. Do not copy credentials, switch to the broader ingress AppRole, or
use the shared DNS account for firewall work. The shared router account retains
ingress ownership and its existing router permissions.

## Risks / Trade-offs

- All services on the printer become reachable from Wi-Fi; other wired hosts
  remain outside this rule. Narrower service selection requires printer-specific
  protocol requirements.
- Historical topology may differ from live state; inspect the live plan and
  verify from a real Wi-Fi client before claiming acceptance.
- DNS and unicast routing do not relay mDNS/AirPrint discovery across subnets.
  Clients can address the printer by its name or IP.
- A changed printer address or missing return gateway breaks connectivity.
  Confirm stable addressing and gateway `192.168.1.1` on the printer.
- The insertion option is not continuous monitoring of administrator reordering;
  verify rule order after deployment and any manual firewall changes.

## Migration Plan

Validate offline, then review the owning DNS and firewall text plans. Operational
calls require the existing DNS credentials and the named shared-secret read policy.
Only apply an additions-only plan within the user's requested scope. Do not
bootstrap credentials or edit the router imperatively without exact approval.
After deployment, resolve the name and open the printer from Wi-Fi, confirm the
exception precedes the drop, and retain a sanitized receipt. Revert the DNS
declaration and firewall resource through their owning workflows for rollback;
established connections can persist until connection tracking expires.
