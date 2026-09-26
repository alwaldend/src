## Context

See [proposal.md](proposal.md) for scope and linked owner plans. The inspected
baseline has XCP-ng provisioning in `infra/openhands/tf_setup`, Yandex VMs in
`infra/threexui/tf_setup`, and a native Traefik role in the Ansible collection.
The download component did not exist before this change. Existing source
describes patterns, not verified live infrastructure.

The main website currently publishes to GitHub Pages. Its linked change owns
the new downloads UI and production SSH deployment. This plan owns the two
hosts and their integration, not copies of those shared behaviors.

## Goals / Non-Goals

**Goals:** Use ordinary VMs, filesystems, Traefik, Nginx, and SSH; retain one
source for each configuration and requirement; make either environment an
explicit deployment target.

**Non-Goals:** Containers, object storage, databases, upload APIs, browser
uploads, automatic replication, automatic failover, backups, garbage collection,
or a rollback command. A website can be redeployed by selecting an existing
release. Routine source work does not authorize live provisioning or cutover.

## Decisions

### Topology and state

Each environment has one VM with 100 GB of persistent filesystem storage for
release archives, ordinary artifacts, extracted sites, and deduplication state.
Traefik configuration and ACME state stay on the system disk under `/opt/traefik`;
the proxy has no dependency on the content mount.
Mount a dedicated Btrfs content disk at `/srv/download` alongside a 20 GiB boot
disk. Yandex retains an independently managed block disk across VM replacement.
The pinned XCP-ng provider owns its disks with the VM and cannot reattach a
retained disk through the VM resource. As explicitly confirmed by the user,
protect the whole XCP-ng VM with `prevent_destroy`; a future change must
define its replacement workflow. These lifecycle protections are not backups.

Use separate local and Yandex Terraform state keys and clearly named command
targets, with the existing Vault HTTP backend. Share one Ansible playbook and
component variables where they are identical. Provider inventory, VM image
IDs, addresses, and CPU/memory sizing are deployment inputs, resolved from
existing infrastructure during implementation.

```text
public DNS ----> Yandex VM: Traefik ----> Nginx ----> /srv/download
local DNS ----> XCP-ng VM: Traefik ----> Nginx ----> /srv/download
publisher ----> chosen VM over SSH ----> release files and selected site
```

Traefik binds ports 80 and 443, redirects HTTP to HTTPS, and routes by host.
Nginx binds `127.0.0.1:8008`, a Fedora SELinux HTTP port. Both public names are
anonymous read endpoints. There is no connection through `infra/ingress`.

### Filesystem and routes

```text
/srv/download/
  projects/<project>/releases/<version>/<uploaded files>
  sites/<project>/releases/<version>/<extracted website content>
  sites/<project>/current -> releases/<version>
  staging/<task-owned temporary content>
```

`download.alwaldend.com/projects/` exposes only the public project tree.
Directory requests return Nginx JSON autoindex, including directories that
contain a file named `index.html`; file requests return the uploaded bytes.
The listing route must deliberately bypass index-file rendering. It must not
expose staging, service state, credentials, or the extracted-site tree.

`alwaldend.com` serves `/srv/download/sites/alwaldend.com/current`, with normal
index-file rendering and no separate versioned website URL. Its archive stays
public at `projects/alwaldend.com/releases/<version>/`. Root `/` on the download
host redirects to `/projects/`. The themed browsing UI lives on the main
website at `/downloads/`, not on the download VM as a separate application.

Preserve `www.alwaldend.com` as an alias: Traefik accepts both HTTP and trusted
HTTPS for it and permanently redirects to `https://alwaldend.com`, retaining
the request path and query. This is a redirect router, not another website
root or release selection. Obtain a certificate covering `www` on each VM
from its configured issuer as well as certificates for the two content names.

The download account owns release and site storage, with no authorized keys
configured for it. Administrators connect using existing SSH access and run
rsync and activation through sudo as that account. Nginx has read access only.
Ansible creates storage, accounts, configuration, and services; the release
tool owns content upload and activation in its separate publication PR. Neither
a configuration rerun nor a service restart removes published content.

### DNS ownership and HTTPS

Both `download.alwaldend.com` and `alwaldend.com` use the normal global/dc1
split: public addresses select Yandex and local addresses select XCP-ng.
Do not publish both hosts as a public round-robin pool. The same release can
be deployed to both explicitly; their files and active symlinks are independent.

`infra/download/dnsconfig.json` owns the download name in both views. Existing
apex records remain in `infra/dns/dnsconfig.json`, updated through that owner's
Terraform workflow. VM addresses feed these declarations without duplicate
maintained values. Retire obsolete apex A/AAAA destinations during cutover;
preserve unrelated records and the existing `www` CNAME pointing at the apex
in both views. That alias continues to belong to the apex DNS owner.

The root-level `dns` stage owns the component's Cloudflare/RouterOS
aliases and delegation through the existing combined DNS module. The Yandex
stage owns its delegated `yc.download.alwaldend.com` zone and host A record,
which consumes the reserved address directly. Public service CNAMEs target
that host; Cloudflare flattens the apex CNAME. The public apex follows the stable download service
name, so the component owns the Yandex endpoint once. The apex owner's local
A record is the single local address input used by Terraform; the unique
host CNAME and local download alias follow it. Prepare that dc1 apex record
through its owner before configuring the local host, while preserving the
old public Pages records until the later two-phase public cutover. The local
website has a maintenance interval during host/content setup.
Each name has one state owner, and VM states remain separate per environment.

Use HTTP-01 in both environments, as requested in PR review: Let's Encrypt
for public Yandex traffic and Vault ACME for local XCP-ng traffic. The shared
Traefik role obtains EAB credentials for the component-scoped Vault PKI role;
Yandex disables EAB. Only the three serving names are allowed by that role,
with no subdomain or client-certificate grant. Local clients trust the
repository CA. DNS credentials are used only by provisioning.

The shared host role installs the CA used to connect to Vault. Traefik keeps
private ACME state per VM and uses environment-specific renewal timing. The
issuer must resolve each name to its corresponding VM and reach port 80;
certificate bootstrap therefore follows DNS routing, with an explicit initial
TLS transition interval. HTTP-to-HTTPS redirects remain compatible with the
challenge handler. Verify issuance and renewal in each environment.

[Traefik's ACME reference](https://doc.traefik.io/traefik/reference/install-configuration/tls/certificate-resolvers/acme/)
and [Vault's ACME reference](https://developer.hashicorp.com/vault/docs/secrets/pki/acme)
describe HTTP validation and external account binding.

### Deduplication

The user selected Btrfs with daily background deduplication. Use the native
`duperemove` package and a component-owned systemd service/timer to scan only
published project files and extracted sites. Keep a persistent hash database
in private service state for incremental scans. The kernel shares identical
extents while files remain independently writable through copy-on-write.
Hardlinks are not needed, and the SSH publisher requires no deduplication logic.

Bound the maintenance job to one I/O and one CPU thread, half a CPU core,
512 MiB memory, idle I/O priority, and a six-hour deadline. Systemd prevents
overlapping service runs and catches up a missed daily timer event. Neither
maintenance failure nor cancellation changes release selection or deletes
files. Compressed archives may offer little sharing; verify actual allocation
savings without promising a capacity multiplier.

### Public listing integration

Nginx supplies its native JSON listing shape, with directory navigation and
presentation owned by the linked website plan. Allow credential-free browser
GET/HEAD requests using CORS on the listing and download endpoints. Use short
or disabled listing caching so completed SSH uploads appear on refresh.
Uploaded file bytes retain ordinary static HTTP behavior, including HEAD and
byte-range requests. The generic Nginx role takes configuration from this
component and contains no download-specific frontend or hostnames.

## Risks / Trade-offs

- Single VM or data-disk failure per environment -> downtime or permanent
  data loss is possible; no backups or failover are requested.
- Independent environments -> publish deliberately to each; a local release
  can differ from the public release until it is deployed to Yandex.
- Archives plus extracted content share 100 GB -> report insufficient space
  before activation and leave the selected site intact; cleanup is manual.
- HTTP-01 needs routed DNS before issuance -> coordinate the certificate
  bootstrap interval; Vault must use the dc1 DNS view for local validation.
- A symlink change is atomic, but browser asset requests span time -> site
  builds must provide a coherent asset strategy; acceptance checks both HTML
  and its assets after activation rather than claiming a transactional session.

## Migration Plan

1. Implement and validate the linked role, identity, release tool, and website
   changes; keep the current public site reachable while preparing candidates.
2. With separately authorized live scope, bootstrap Vault and provider
   assignments, provision both VMs, and apply the shared Ansible deployment.
3. Deploy a website archive to each explicit SSH target and verify HTTP
   content using address overrides before DNS changes.
4. Apply reviewed split-horizon records through their owners. Once the issuers
   can reach their respective hosts on port 80, exercise public staging and
   production issuance and local Vault issuance, then verify both trust chains,
   downloads, the main website, and `www` path/query-preserving redirects.
5. Redeploy any retained site version through the same release command when a
   different selection is needed. Do not add a rollback operation.

## Acceptance and continuation

The current implementation pass covers IaC, the linked Nginx role and Vault
identity source, and the architecture diagram. The linked SSH release tool
and component publishing targets are delivered in a separate PR against
`master`; the website browser remains separate implementation work. See
[acceptance.md](acceptance.md) for the matrix written before implementation
and [evidence.md](evidence.md) for checks and their limits.

Before live rollout, run the integrated SSH-to-HTTP/browser fixture. Retain
commands, candidate identity, request results, artifact checksums, and browser
evidence. Live inventory, SELinux enforcement, DNS, certificate renewal, and
reboot checks require authorized environments and separate evidence.
