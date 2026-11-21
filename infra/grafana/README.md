---
title: Grafana
date: 2026-11-21
---

One Fedora VM in DC1 runs Grafana from the Fedora package under systemd.
Terraform provisions the XCP-ng guest through Xen Orchestra with two
vCPUs, 4 GiB of RAM, and a persistent 20 GiB disk. Grafana stores its
SQLite database at `/var/lib/grafana/grafana.db`. Ansible configures the
guest, installs Grafana through the `alwaldend.main.grafana` role, and sets
up Traefik. VictoriaMetrics scrapes the proxy's Prometheus endpoint. The
guest allows metrics traffic from the collector address.

Open [Grafana](https://grafana.alwaldend.com/) from the DC1 network with
the Vault server CA trusted by your browser. Traefik serves Vault-issued
HTTPS, redirects HTTP, and proxies to Grafana on loopback. The firewall
allows administrative SSH from any routed source. Select **Sign in with
Vault** to sign in as a member of `global_admins`. Those users become
Grafana server administrators. The guest's SELinux policy allows the Grafana
process to connect to Vault's OIDC API on TCP 8200.

From outside DC1, use the same hostname with a trusted ingress client
certificate. Shared ingress verifies the guest's Vault-issued certificate
and forwards traffic through the dedicated `dc1.grafana.alwaldend.com`
backend name. Grafana still requires application sign-in.

## Deployment

1. Review and apply the [Grafana Vault AppRole](approle/README.md),
   then reconcile shared memberships with `//infra/vault/tf:tf.apply`.
   Synchronize the AppRole identity and resource set using the
   [Xen Orchestra workflow](../xcp_ng/tf/README.md).
2. Store the local recovery password and a separate stable Grafana encryption
   key in Vault. Generate both independently. The example reads protected
   files containing the values.

   ```sh
   bazel run //infra/grafana:vault.kv_put -- \
     alwaldend.com/vault1/approles/src_infra_grafana/config \
     admin_password=@"$grafana_password_file" secret_key=@"$grafana_secret_key_file"
   ```

3. Confirm the address in the [DNS declaration](dns/dnsconfig.json) is
   available. Review the VM and DNS plans separately before applying them.

   ```sh
   bazel run //infra/grafana/tf_setup:tf_setup.plan
   bazel run //infra/grafana/tf_setup:tf_setup.apply
   bazel run //infra/grafana/dns:tf.plan
   bazel run //infra/grafana/dns:tf.apply
   ```

4. Wait for cloud-init, verify the guest SSH host key out of band, and check
   that Vault resolves the DC1 DNS record. Deploy with
   `bazel run //infra/grafana/ansible`.
5. With Traefik running, verify Vault can reach guest TCP 80 for ACME
   HTTP-01. Check the HTTPS certificate and sign in through Vault to verify
   the server administrator role.
6. Once Grafana DNS and the guest firewall are ready, deploy
   [shared ingress](../ingress/README.md). Verify that ingress can reach
   the guest over HTTPS and that the public route requires a trusted client
   certificate.
7. Reapply the [VictoriaMetrics playbook](../victoria_metrics/README.md) and
   check that the Grafana proxy is up in the `traefik_dc1` job.

When adding ingress to an existing Grafana VM, apply its AppRole and DNS
changes and rerun the Grafana playbook before deploying shared ingress.

Grafana uses the local `admin` password only when initializing an empty
database. Change an existing password through Grafana and update the
recovery value in Vault. Keep the encryption key with the database backups
to preserve access to encrypted data source settings.
