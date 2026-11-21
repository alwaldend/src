---
title: VictoriaMetrics
---

## Architecture

[![VictoriaMetrics VM inside DC1 infrastructure, with dashboard access and Traefik scrape paths](docs/assets/architecture.svg)](docs/assets/architecture.svg)

Bazel renders the [Mermaid source](docs/architecture.mmd) in light and dark themes.

One Fedora VM on XCP-ng runs VictoriaMetrics single-node under systemd.
The guest starts with 2 vCPUs, 4 GiB RAM, and a persistent 40 GiB boot and
data disk. The database lives at `/var/lib/victoria-metrics`. Retention is
30 days, with an ingestion cutoff when free disk space falls below 5 GiB.
Bazel verifies the release archive and packages the executable for Ansible.

VictoriaMetrics scrapes its loopback endpoint every 15 seconds and the
Traefik proxies every 30 seconds. The `traefik_ingress` job uses active
ingress hosts from the WireGuard inventory. The `traefik_dc1` job covers
Forgejo, dev1, three OpenHands hosts, the 3x-ui control plane, and the
local download host. Each proxy firewall allows metrics traffic from the
VictoriaMetrics guest address. Cloud-only proxies need a private route
before they can join these scrape jobs.

Add scrape jobs to the [scrape YAML template](ansible/templates/scrape.yaml.j2).
Bazel packages the ingress inventory as an extra-vars file. Ansible renders
the active WireGuard targets from it.

## Provisioning

1. Provision the [dedicated Vault AppRole](approle/README.md), including its
   DNS policy, SSH host certificate signer, and the Vault ACME role.
   Run the [web credential setup](tf_setup/web_auth/README.md) to generate
   and store the dashboard login and verifier.
   Reconcile Vault memberships and synchronize the OIDC identity and resource set
   through the [Xen Orchestra workflow](../xcp_ng/tf/README.md).
2. Verify that the address in the [DNS declaration](dns/dnsconfig.json) is
   available on the DC1 VM network. Review the component's DNS and VM plans
   separately before applying them.
3. Once SSH is available and cloud-init has finished, verify the guest's host
   key and run `bazel run //infra/victoria_metrics/ansible`. The configured
   SSH destination on DC1 is `victoria-metrics.alwaldend.com`. The deployment
   also allows an explicit `ansible_host` override before DNS is available.
4. Verify Vault resolves the service hostname through DC1 DNS and can reach
   guest Traefik on port 80 for ACME validation. Check that Traefik serves the
   Vault certificate and shared ingress can reach the backend. Guest Traefik
   renews the certificate automatically. Shared ingress uses Let's Encrypt.

## Dashboard

Open `https://victoria-metrics.alwaldend.com/vmui/` with a trusted ingress
client certificate and the shared Basic Auth login. Members of Vault's
`global_admins` group can read the credentials at
`secrets/alwaldend.com/vault1/approles/src_infra_victoria_metrics/web_auth`.
Ingress verifies the guest's Vault-issued certificate.
On the DC1 LAN, use the same dashboard URL and Basic Auth login. DC1 DNS
points the service name directly at guest Traefik, so the browser must trust
the Vault server CA. HTTP redirects to HTTPS. Shared ingress alone uses the
`dc1.victoria-metrics.alwaldend.com` upstream name.
Select **Dashboards → Monitoring overview** for scrape health, Traefik traffic
and certificate expiry, and VictoriaMetrics disk and ingestion metrics.
The panels are defined in [monitoring-overview.json](ansible/files/monitoring-overview.json).
For a table of targets that are down, run `up == 0` in **Query** and select
**Table**.

The HTTP listener on guest loopback also remains available through SSH
local port forwarding to `127.0.0.1:8428` using the Ansible account.
This deployment has one failure domain. Guest or disk loss
can make the stored metrics unavailable. Off-VM backup storage is outside
this component.
