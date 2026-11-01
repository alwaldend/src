---
title: Download
description: Public static files and websites on XCP-ng and Yandex Cloud
---

![Static hosting architecture](assets/architecture.svg)

One VM per environment runs the shared Traefik service in front of native
Nginx. Clients connect directly on ports 80/443; Nginx listens only on
`127.0.0.1:8008`. Public DNS selects Yandex and dc1 DNS selects XCP-ng.
These are independent stores with explicit local-to-public content sync;
there is no automatic replication, failover, or backup.
The diagrams describe checked-in configuration and operator procedures;
they do not establish current deployment or service health.

## Content and publication

Each VM has a 100 GiB Btrfs content disk mounted at `/srv/download`:

```text
projects/<project>/releases/<version>/<uploaded files>
third_party/<name>/releases/<version>/<uploaded files>
sites/alwaldend.com/releases/<version>/<extracted website>
sites/alwaldend.com/current -> releases/<version>
staging/
state/
```

`download.alwaldend.com/` returns a JSON index of `projects/` and `third_party/`
without redirecting. Both namespaces provide JSON directory listings and serve
file bytes, including website archives. `alwaldend.com` serves the selected
extracted website. `www.alwaldend.com` redirects permanently to the apex,
preserving the path and query. Staging and service state have no public route.
Directory listings never use a release's `index.html` as an index page.

The public host redirects ordinary HTTPS requests from `git.alwaldend.com`
to `github.com` in Traefik, preserving the path and query string. Requests
with `go-get=1` instead reach Nginx and receive static Go discovery metadata
for `git.alwaldend.com/alwaldend/src` and its descendants, identifying the Git
repository `https://github.com/alwaldend/src.git`. Unknown discovery paths
return 404. The metadata keeps Go vanity module paths usable through this
public hostname.
Yandex obtains its certificate through the existing Let's Encrypt resolver.
The redirect is disabled locally, where the Forgejo-owned DNS record continues
to select Forgejo. Its public CNAME also remains owned by `infra/forgejo`.

The `download` account owns content and staging; Ansible does not configure
its `authorized_keys`. Administrators use existing SSH access and sudo to run
publication as this account. The host role installs rsync. Nginx reads content
and cannot write it. The host-configuration playbook never publishes or selects a release.

The release tool and component publishing targets are available in this tree.
Their commands and transfer behavior belong
to the [release tool documentation](../../tools/release/README.md). The
[website browser change](../../projects/alwaldend.com/openspec/changes/add-download-browser-and-ssh-hosting/proposal.md)
also remains separate implementation work.

## Deduplication

A daily `download-deduplicate.timer` runs `duperemove` over `projects/`, `third_party/` and
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

## Operations

The content disk is optional for OS boot (`nofail`, ten-second device timeout),
while Nginx and deduplication require it. SSH and Traefik remain on the system
disk. Deployment installs the declared TCP 22 management rule before deleting
the inherited `SSH`/`mDNS` rules listed in `firewall_delete_rules`, including
multicast destinations. Other rules remain. Retiring an allow requires an
explicit deletion entry; removing its declaration alone does not revoke it.

Shared role behavior and the isolated security-fix evidence are recorded in
[deployment hardening](../../projects/ansible_collection/openspec/changes/harden-download-deployment/evidence.md).
The canonical package-update role owns native automatic updates and reboot
detection. Immediate and automatic reboots are disabled by default; set
`update_all_packages_reboot_enabled: true` only for an explicitly approved
maintenance reboot.
The [2026-10-26 local rollout](openspec/changes/add-static-hosting/evidence.md#authorized-local-deployment-2026-10-26)
records the complete playbook, reboot persistence, serving checks, and an
unchanged firewall/Nginx/Traefik rerun. The
[Yandex rollout](openspec/changes/add-static-hosting/evidence.md#yandex-zone-change-2026-10-26)
records the completed download-only deployment after moving to `ru-central1-a`,
trusted public HTTPS, and firewall/service persistence across reboot. The
[website rollout](openspec/changes/add-static-hosting/evidence.md#versioned-site-publication-and-mirror-2026-10-26)
records publication and mirroring of version 2026.44.0. Public apex/www routing
is enabled; the [DNS owner](../dns/openspec/changes/activate-public-download-site/evidence.md)
records the completed cutover and trusted HTTPS checks.

![Publication, maintenance, and recovery operations](assets/operations.svg)

Use `ssh.local` or `ssh.yandex` for operator access through the shared
[`ssh_binary`](../../tools/ssh/README.md) rule. Each invocation authenticates
with the component AppRole, injects a temporary Vault-signed client key, and
connects as `ansible` to the selected inventory FQDN. It verifies the host
certificate against the repository-owned SSH server CA, scoped to
`*.download.alwaldend.com`, with strict checking enabled. Client configuration
and SSH-agent identities are excluded; no persistent SSH configuration or
known-hosts file is changed. AL removes temporary keys and trust files on exit.

For example, inspect the local services without changing the host:

```sh
bazel_agent bazel run //infra/download:ssh.local -- \
  systemctl is-active nginx traefik download-deduplicate.timer
bazel_agent bazel run //infra/download:ssh.local -- \
  findmnt /srv/download
```

Invocation arguments follow the fixed SSH destination and form the remote
command. Omitting them opens a shell. Select `ssh.yandex` explicitly for cloud
access. The target supplies authentication, not authority to change a host.

Publication explicitly selects `publish.local` or `publish.yandex`. These
wrappers inject the same temporary key, signed certificate and scoped host CA
as the SSH targets. `--namespace third_party` selects versioned upstream
artifacts; project publication remains the default. Upload,
extract, and link are separate ordered steps; only link selects the website.
Both publication targets inject the same temporary SSH credentials and host
trust as the SSH and sync targets. Normal publication rejects an existing
release version and stages new uploads before making them public. Linking
removes older extracted website versions while retaining downloadable archives.
Restore a previous archive with `--steps extract,link` for rollback. Public sync
overwrites from local and propagates the extracted-site cleanup.
The release tool owns the
[commands and failure behavior](../../tools/release/README.md#choose-publication-steps).

Publish and inspect locally, then mirror that published content to Yandex:

```sh
bazel_agent bazel run //infra/download:sync.yandex -- --check
bazel_agent bazel run //infra/download:sync.yandex
```

**Local is the source of truth:** differing public content is replaced and
public-only files, releases, and site selections are removed. The target runs
`ansible/playbook_sync.yaml` with the existing inventory and temporary Vault
SSH credentials. Host trust remains strict and both content mounts are required.
Check mode reports changes and deletions without changing published content.

The playbook captures only `projects/` and `sites/` on the controller and
validates paths, current links, and site indexes before destination writes.
Credentials stay on the controller. Source and destination publishers must
remain idle throughout capture and selection. Each run needs controller space
for all source content; the private capture is removed afterward.

Native `ansible.posix.synchronize` stages changed trees privately, checksums
content, normalizes directories to 0755 and files to 0644, and reuses identical
destination files through rsync hardlinks. A transfer failure leaves public
content unchanged. The native `mv --exchange --no-copy -T` command atomically
replaces each complete project under a synchronization lock. It requires
support from the destination's GNU coreutils. Projects change independently;
a failure during selection may leave some updated and others unchanged.
Rerunning converges. The `always` cleanup removes only this run's capture and
staging; interruption or an unreachable host can prevent cleanup.

Configuration, certificates, private state, and DNS are excluded. Copying the
site does not enable disabled ingress routes. The playbook does not build
releases or run host configuration. It preserves the locally selected website
and all locally retained release archives. This is a mirror, not a backup.

The manual `//infra/download/test/sync:sync_test` target runs the packaged
playbook and release CLI against isolated directories through real OpenSSH
and rsync. It uses a private user/mount namespace to verify that site cleanup refuses
mounted release directories without deleting private state. It checks trust
and layout refusals, authoritative replacements and
deletions, failed staging, repeated execution, and archive-preserving local
cleanup. Its `sync-receipt.json` contains repeatable hashes, modes, and links.

On the selected host, inspect `systemctl status nginx traefik`,
`systemctl list-timers download-deduplicate.timer`, and
`journalctl -u download-deduplicate.service`. Check the mount with
`findmnt /srv/download`, free space with `btrfs filesystem usage /srv/download`,
and sharing with `btrfs filesystem du` on the published trees. Verify DNS,
certificate trust, website assets, JSON listings, downloads, and `www` redirects
from each DNS view. Extracted-site cleanup runs after release selection; archives remain retained.
Deduplication never deletes releases. Retained disks and earlier releases provide no
backup against disk loss.

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
  `download.alwaldend.com`. The public
  `download.alwaldend.com` name follows the component's Yandex host.
  Its local apex A record is the canonical VM address consumed by
  Terraform. The unique local host CNAME and download alias follow that local
  apex. Endpoint changes therefore have one maintained input per environment.
  Ansible connects by inventory FQDN, without parsing DNS source.
  Old Pages A/AAAA records are removed through staged owner cutovers;
  unrelated mail and staging records remain.
- `ansible` mounts storage, configures SSH and service access, labels static
  content for SELinux, and deploys the shared Traefik and Nginx roles.

The local address `192.168.10.66` comes from the apex DNS declaration and is
assigned to the managed XCP-ng VM. The cloud image defaults to the immutable Fedora 43 image
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

![Public Yandex Cloud and local XCP-ng serving paths](assets/deployment.svg)

[Download the dark WebP](assets/deployment.dark.webp). Publisher denotes the
controller-side release or synchronization workflow; administrators use SSH
to run it as the content account. The site path abbreviates the selected
release under `sites/alwaldend.com/current`. The GitHub edge represents a
browser redirect, not a proxy connection.

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

## Third-party namespace deployment

`//infra/download/ansible:ansible.local` and `:ansible.yandex` accept
`-- --tags download_publication` for mounted-storage checks, publication roots,
SELinux labels, deduplication units and Nginx. `--tags nginx` restricts a routing
update to the shared Nginx role. Both paths retain strict host-CA verification
with the component injector. Deployment requires explicit host scope.
