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
release archives, ordinary artifacts, extracted sites, and service state.
Use a retained data disk mounted at `/srv/download`; size the boot disk using
the existing VM pattern. Yandex uses a VM-attached block disk, not object
storage. Retaining the data disk across VM replacement protects redeployment,
but is not a backup.

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
Nginx binds loopback, conventionally `127.0.0.1:8080`. Both public names are
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
before cutover as well as certificates for the two content hostnames.

The SSH publisher can write release and site storage. Nginx has read access
only. Ansible creates storage, accounts, configuration, and services; the
release tool owns content upload and activation. Neither a configuration
rerun nor a service restart removes published content.

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

Within `tf_setup`, give download DNS one state-owning stage and keep VM state
separate per environment. That DNS stage consumes the declared addresses for
both VMs, avoiding two states that manage the same record. Reuse the existing
combined Cloudflare/RouterOS module rather than adding view-filter behavior.

Use public ACME certificates on both hosts. DNS-01 validation through the
existing DNS provider is the selected approach because the local host is not
the public A/AAAA destination. The repository Traefik role supplies the binary
and service; component configuration supplies the public resolver, provider
secret references, and persistent ACME state. Confirm provider support in the
pinned binary, exercise issuance with the ACME staging service, and verify
renewal for both hosts before cutover. Do not copy a Vault-private certificate
pattern into an anonymously browsable public service.

[Traefik's ACME reference](https://doc.traefik.io/traefik/v3.3/https/acme/)
describes DNS challenges and provider credentials. Independent renewals for
the same names require testing concurrent TXT-record handling; do not share a
writable ACME state file between hosts.

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
- Public DNS points away from the local VM -> test DNS-01 issuance and
  renewal on both hosts, including challenge overlap, before DNS cutover.
- A symlink change is atomic, but browser asset requests span time -> site
  builds must provide a coherent asset strategy; acceptance checks both HTML
  and its assets after activation rather than claiming a transactional session.

## Migration Plan

1. Implement and validate the linked role, identity, release tool, and website
   changes; keep the current public site reachable while preparing candidates.
2. With separately authorized live scope, bootstrap Vault and provider
   assignments, provision both VMs, and apply the shared Ansible deployment.
3. Obtain trusted certificates and deploy a website archive to each explicit
   SSH target. Test both content hostnames and the `www` HTTPS redirect using
   address overrides before DNS changes.
4. Apply reviewed split-horizon records through their owners and verify both
   public and local resolution, TLS, downloads, the main website, and `www`
   path/query-preserving redirects.
5. Redeploy any retained site version through the same release command when a
   different selection is needed. Do not add a rollback operation.

## Acceptance and continuation

All implementation tasks are initially unchecked. Run an isolated, repeatable
SSH-to-HTTP fixture before live rollout, with its scenario matrix written
before implementation. Retain commands, candidate identity, request results,
artifact checksums, and browser evidence. Live DNS, certificate renewal, and
reboot checks require authorized environments and separate evidence.

The first implementation action is to complete the linked Nginx and Vault
plans, then assemble component provisioning and the end-to-end fixture.
New OpenSpec owner packaging and aggregate validation registration for this
component belong to implementation; this delivery writes change artifacts only.
