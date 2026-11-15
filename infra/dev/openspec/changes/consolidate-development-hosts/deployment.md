# Deployment record

The user confirmed PR #305 merged and authorized rebase and deployment, then
requested OpenHands removal and secure1 recreation without disk migration.
All live changes used the owning Terraform or packaged Ansible workflow.

## Ownership and DNS

The shared src_infra_dev identity owns both hosts, DNS and its XO resource set.
SSH signing is restricted to the two hostnames, TLS issuance to dev1 services,
and DNS credential reads to dc1 and global views. It uses its own state
namespaces and has no access to either retired backend.

Seven RouterOS and five Cloudflare records were released from the old states,
imported by existing provider ID and renamed in place. No DNS objects were
recreated. Both old backends were confirmed empty after the old secure VM and
its four lookup bindings were removed. Both legacy Vault identities, their
remaining grants and the old runner XO resource set have been retired.
Temporary migration inputs and deployment targets were removed.

At the user's request, secure1 was recreated again with 192.168.10.102 in its
canonical DNS declaration and cloud-init. Two duplicate-address ARP probes
received no replies before assignment. The interactive apply replaced only
xenorchestra_vm.secure1; a separate apply changed only its existing RouterOS
DNS row. The VM retains 8 CPUs, 16 GiB RAM and one 500 GiB boot disk.
Final VM and DNS plans reported no changes on 2026-11-15 at 18:01 MSK.

## Hosts and ingress

The new secure1 completed cloud-init with the declared hostname and .102
address. Initial IPv4 SSH timed out; strictly verified IPv6 SSH allowed the
shared host role to reconcile its firewall. IPv4 SSH then became reachable.
Ansible now generates missing standard host keys before hardening, because
cloud-init's seeded Ed25519 identity did not create the ECDSA key expected by
the upstream role. Bootstrap trust came from the created VM's public Terraform
output; subsequent SSH uses the normal signed host certificates.

Dev1's deployment configured the shared hostname, workstation services,
known_hosts and Traefik. OpenHands was subsequently stopped, disabled and
masked; its executable environment and secret file are absent. Its newly
staged Python 3.12 layering request and Vault service policy were removed.
Existing uv tooling and application work/state remain. Dev1 was not rebooted.

Ingress routes target dev1 with matching TLS SNI and retain public mTLS.
Dev1 Traefik's ACME account store is keyed by directory URL. The old account
belonged to a different Vault directory and caused HTTP 400 despite correct
EAB enrollment. Ansible restarted Traefik; the user also cleared the previous
account data and confirmed the public T3 URL works. A final direct request returned HTTP 200 with the checked-in root CA and TLS verification.
ComfyUI retains its separately launched, on-demand application process.

The final secure1 deployment completed with 212 successful tasks, 53 changes,
zero unreachable hosts and zero failures. Its reboot task was skipped.

## XO prerequisite

The user explicitly approved repair of XO's HTTPS certificate. Its managed
file contained a hostname-less self-signed certificate, and the listener left
XO's autoCert default enabled. Ansible now validates certificate chains and
names, reissues only missing or invalid certificates, and declares
"autoCert = false" in the inspected TLS listener. The approved deployment
restarted XO management only; guest VMs continued running. A strict TLS
request returned HTTP 302 and the development OIDC bootstrap succeeded.

## Validation

Existing repository quality, DNS, Terraform formatting, owner-local OpenSpec,
packaged Ansible syntax, documentation E2E and semantic lint checks cover the
source. Live outcomes above are separate evidence. Native Ansible SSH askpass
was verified against XO and the hypervisor; the pinned Ansible runtime supports
the hosts' Python 3.15. Detailed local command receipts remain in ignored task
scratch; no raw Terraform plans or state are retained as task evidence.

## Subsequent maintenance

The user explicitly requested a secure1 reboot after package updates. The
existing Ansible reboot task completed with two successful tasks and zero
failures; strict IPv4 SSH verification then confirmed the new boot and kernel
7.3.3-200.fc44.x86_64. Dev1 was not rebooted and the default reboot policy
remains disabled.

The user requested removal of the completed OpenHands cleanup block. Normal
deployment now contains no OpenHands task or secret input. Packaged syntax and
the dev1 task listing verify that the one-time removal logic is gone. The
service retirement already completed before this source cleanup.
