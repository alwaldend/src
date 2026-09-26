---
title: Download
description: Public static files and websites on XCP-ng and Yandex Cloud
---

![Static hosting architecture](assets/architecture.svg)

One VM per environment runs the shared Traefik service in front of native
Nginx. Clients connect directly on ports 80/443; Nginx listens only on
`127.0.0.1:8008`. Public DNS selects Yandex and dc1 DNS selects XCP-ng.
These are independent stores with no replication, failover, or backups.

## Content and publication

Each VM has a 100 GiB Btrfs content disk mounted at `/srv/download`:

```text
projects/<project>/releases/<version>/<uploaded files>
sites/alwaldend.com/releases/<version>/<extracted website>
sites/alwaldend.com/current -> releases/<version>
staging/
state/
```

`download.alwaldend.com/projects/` returns JSON directory listings and serves
file bytes, including website archives. `alwaldend.com` serves the selected
extracted website. `www.alwaldend.com` redirects permanently to the apex,
preserving the path and query. Staging and service state have no public route.
Directory listings never use a release's `index.html` as an index page.

The `download` account owns content and staging; Ansible does not configure
its `authorized_keys`. Administrators use existing SSH access and sudo to run
publication as this account. The host role installs rsync. Nginx reads content
and cannot write it. Ansible never publishes or selects a release.

The release tool and component publishing targets are delivered in
[PR #113](https://github.com/alwaldend/src/pull/113), independently against
`master`. Their commands and transfer behavior belong
to the [release tool documentation](../../tools/release/README.md). The
[website browser change](../../projects/alwaldend.com/openspec/changes/add-download-browser-and-ssh-hosting/proposal.md)
also remains separate implementation work.

## Deduplication

A daily `download-deduplicate.timer` runs `duperemove` over `projects/` and
`sites/` on each VM. It shares identical extents while retaining independent
files, paths, permissions, and copy-on-write behavior. Repeated website assets
can share storage; separately compressed archives may have fewer identical
extents. Savings are workload-dependent and stay within each environment.

A private hash database under `state/deduplication/` makes subsequent runs
incremental. Staging, ACME state, and the database are outside the scan roots.
The job runs with one I/O and one CPU thread, idle I/O priority, a half-core
CPU limit, and a 512 MiB memory limit. Systemd serializes runs; the timer makes
up a missed run after downtime. The service requires the content mount.

The job never deletes or activates releases. Its result is available through
`systemctl status download-deduplicate.service` and the service journal.
`btrfs filesystem du` reports shared/exclusive allocation; ordinary file sizes
and download responses remain their logical sizes. The declared filesystem
setup does not force formatting over an existing filesystem.

## Infrastructure ownership

Terraform provisioning, provider assignments, and DNS declarations are merged
through [PR #108](https://github.com/alwaldend/src/pull/108). Host configuration
connects by unique inventory FQDNs and needs no packaged DNS inputs. The shared
AL configuration selects provisioning or host credentials by stage. DNS
snapshots are regenerated only on manual request.

`local`, `yandex`, and `dns` are separate Terraform roots directly under this
project. Their separate state keys are
`alwaldend.com/vault1/approles/src_infra_download/tf_backend/<stage>`.
Selecting one root cannot plan deletion of another root's resources, and no
root reads another root's complete Terraform state.

- `local` owns the XCP-ng VM and two disks. `prevent_destroy` blocks
  VM replacement because the pinned provider cannot retain and reattach its
  data disk separately. Replacement needs a future reviewed storage workflow.
- `yandex` owns the cloud VM, network, firewall, reserved public
  address, retained content disk, and its Yandex DNS zone and host record.
  VM replacement reuses the content disk and reserved address.
- `dns` owns only this component's Cloudflare/RouterOS aliases and
  delegation. Each stage has its own Vault HTTP state key.
- `infra/dns` retains apex, mail, and `www` ownership. Its public apex follows
  `download.alwaldend.com`; that stable service name follows the component's
  Yandex host. Its local apex A record is the canonical VM address consumed by
  Terraform. The unique local host CNAME and download alias follow that local
  apex. Endpoint changes therefore have one maintained input per environment.
  Ansible connects by inventory FQDN, without parsing DNS source.
  Old Pages A/AAAA records are removed through staged owner cutovers;
  unrelated mail and staging records remain.
- `ansible` mounts storage, configures SSH and service access, labels static
  content for SELinux, and deploys the shared Traefik and Nginx roles.

The local address `192.168.10.66` is unclaimed in checked-in DNS declarations;
confirm its availability and the named XCP-ng template/resource-set inventory
before deployment. The cloud image defaults to the immutable Fedora 43 image
already used by `infra/threexui`; local provisioning uses Fedora 44.

## Credentials and bootstrap

The merged `infra/vault/approles/src_infra_download` root owns the component
identity. Its applied identity and core memberships are documented by the
Vault owner. `infra/xcp_ng/tf` owns the XO resource-set assignment;
`infra/yandex_cloud/org1/tf` owns the folder and its standard folder-scoped
`admin` service account.
AL selects only the chosen stage's provider credentials and state backend.

AL reads the existing Cloudflare/RouterOS credentials for DNS provisioning,
and the component Yandex account from Vault. Traefik uses
HTTP-01 in both environments: Let's Encrypt on Yandex, and the existing Vault
ACME service with role-scoped EAB on XCP-ng. The component's Vault PKI role
permits only the apex, download, and `www` names. Local clients must trust the
repository CA; neither environment requires a client certificate or login.
Traefik configuration and private ACME state stay on the system disk under
`/opt/traefik`, independent of the content mount. DNS provider credentials
are not installed on either host.

The shared host role installs the CA needed to reach Vault. Each inventory key
is the unique host FQDN used for connection and SSH host-key checks. The
shared SSH role includes `inventory_hostname` when signing keys; no host-key
alias or address override is needed.

## Deployment sequence

The following are operator entry points, not evidence of a deployed service.
Confirm the local address, named XO inventory, template disk/interface
layout, and cloud image before deployment. Each live operation requires its
own approved scope.

1. Complete the Vault prerequisite and ensure the existing XO OIDC service is
   configured through its [owner](../xcp_ng/tf/README.md). Obtain explicit
   authorization for the first-login bootstrap of `src_infra_download`:
   it creates the external XO identity/groups. Run the existing
   [bootstrap target](../xcp_ng/cmd/xo_login/README.md) and record its result:

   ```sh
   XO_BOOTSTRAP_APPROLE=src_infra_download bazel_agent bazel run //infra/xcp_ng/cmd/xo_login:bootstrap
   ```

2. After that login, review/apply the XO resource-set assignment through
   `infra/xcp_ng/tf` to reconcile the identity's bindings in IaC. Confirm the
   AppRole is absent from `approles_pending_oidc_login` and its resource set
   has the intended subject before planning the local VM. Review/apply the
   Yandex folder assignment through `infra/yandex_cloud/org1/tf` as well.
3. Review and apply `//infra/download/local:tf.plan`/`tf.apply` and
   `//infra/download/yandex:tf.plan`/`tf.apply` independently.
4. Prepare the dc1 apex A record through `infra/dns/tf` while preserving
   public Pages records, following the [local DNS procedure](dns/README.md#local-dns-preparation).
   Apply `//infra/download/dns:tf.plan`/`tf.apply` to establish cloud delegation
   and host/service aliases. Verify both inventory FQDNs before Ansible. The
   local website is in maintenance until host configuration and publication.
5. Run `//infra/download/ansible:ansible.local` or `:ansible.yandex`. There
   is deliberately no unqualified command that configures both environments.
   Bootstrap host-key trust through the existing SSH procedure before the
   first connection; subsequent host certificates use the component role.
6. Publish the site content and verify HTTP routing with address overrides.
   Follow the [two-phase apex cutover](dns/README.md#staged-public-apex-cutover)
   through `//infra/dns/tf:tf.plan`/`tf.apply`.
   HTTP-01 requires public DNS/port 80 to reach Yandex and Vault's DNS view/port
   80 to reach XCP-ng before either issuer can validate the three names.
7. Exercise Let's Encrypt staging with an explicit directory override and
   separate ACME storage, then select production; test Vault issuance locally.
   This sequence has a certificate-bootstrap interval after DNS cutover.
   Verify both trust chains, renewal, `www` redirects, downloads, and reboot
   persistence. Keep the previous deployment available for a reviewed recovery.

Use `bazel_agent bazel run <target>` for these commands. Provider credentials
and backend state are injected by AL; no host-installed Terraform or Ansible
is required. See the [acceptance matrix](openspec/changes/add-static-hosting/acceptance.md)
for isolated checks and the remaining live acceptance cases.
[Provisioning evidence](openspec/changes/add-static-hosting/terraform-evidence.md)
retains the Terraform validation and review history.
