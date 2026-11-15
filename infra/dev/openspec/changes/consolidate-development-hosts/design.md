## Context

Infra/dev owns dev1 at 192.168.1.210 and secure1 at 192.168.10.102.
Dev1 retains workstation storage and services; secure1 keeps the host-only
runner baseline with Actions disabled.

## Decisions

Use one packaged Ansible invocation with separate host plays and the shared
src_infra_dev AppRole. Limit SSH signing to these hosts, TLS issuance to dev1's
services, and DNS access to the existing dc1 and global provider credentials.
The development identity owns its DNS and provisioning backends; it receives
no access to either legacy namespace.

One DNS root consumes the entire canonical DNS document. The authorized
handoff adopted seven RouterOS and five Cloudflare records by provider ID and
renamed them in place. Both old backends are now empty, so their identities,
remaining grants and temporary deployment wiring can be removed.

Traefik account storage is keyed by the complete ACME directory URL. Reusing
an account registered against another Vault role caused HTTP 400 despite valid
EAB enrollment. Managed Ansible creates the account store for the new directory
and restarts Traefik when configuration changes.

OpenHands retirement is complete: its installation and secret inputs are gone,
its service is disabled, its executables are absent and its Python 3.12 layering
request is removed. Remove the completed cleanup block from normal deployment.
Preserve existing uv tooling and application work/state.

The user chose secure1 recreation without migration. Use the existing scoped
XO workflow and preserve its 8 CPUs, 16 GiB RAM and 500 GiB boot disk. The new
provisioning root uses the shared Xen cloud-init foundation and the canonical
DNS address. Seed an Ed25519 host key through the existing private-key injector
and read its public identity from the created resource for strict initial SSH
verification. The shared SSH role creates missing standard host keys before signing and hardening
and installs the permanent host certificates. Cloud-init is creation-only;
normal plans do not rotate the bootstrap key.

The first canonical VM completed cloud-init at 192.168.10.101 but IPv4 SSH
was unreachable. The user requested another recreation and a new IP. No
checked-in allocation used 192.168.10.102 and two ARP probes received no replies.
The second VM completed cloud-init at .102; deploying the declared firewall
over IPv6 then made IPv4 SSH reachable. Use native Terraform replacement with
manual plan review for explicitly authorized recreations.

## Deployment boundaries

Vault, DNS, VM provisioning and host deployment are separate live operations.
The user authorized this rollout and the XO management certificate repair.
Dev1 must not reboot; package updates leave reboot disabled by default.
The obsolete undeployed Proxmox runner was removed from IaC only. The separate
legacy XO secure VM was destroyed under the later recreation request.
