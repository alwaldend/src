## Purpose

Provide public static file and website hosting on equivalent local and cloud
VM deployments, with explicit content publication and split-horizon routing.

## ADDED Requirements

### Requirement: Independent single-VM environments

The service SHALL provide one VM on XCP-ng and one on Yandex Cloud, each with
100 GB of persistent content storage and independently selectable deployment
and state. Applying configuration to one environment SHALL NOT delete or
reconfigure the other environment's VM or content.

#### Scenario: Deploy one environment

- **WHEN** an operator selects the local deployment
- **THEN** only the local VM and its selected environment resources change
- **AND** the cloud VM and published content remain unaffected

### Requirement: Direct public HTTP and HTTPS

Each VM SHALL accept traffic directly through its Traefik service and serve
static content through Nginx on loopback. Both configured hostnames SHALL
permit anonymous reads and redirect HTTP to HTTPS. Yandex SHALL use publicly
trusted Let's Encrypt certificates; XCP-ng SHALL use Vault certificates trusted
by clients with the repository CA. Both SHALL use HTTP-01 validation. Serving
SHALL NOT depend on shared ingress.

#### Scenario: Read without client credentials

- **WHEN** a visitor requests either hostname on the selected VM
- **THEN** the visitor can read public content without a login or client certificate
- **AND** HTTPS terminates on that VM's Traefik instance

#### Scenario: Validate the local hostname certificate

- **WHEN** local DNS resolves the hostname to XCP-ng while public DNS points to Yandex
- **THEN** the local VM can obtain and renew a Vault certificate for that hostname
- **AND** Vault validates HTTP-01 against the local DNS destination

### Requirement: Preserve the existing website alias

Each VM SHALL accept `www.alwaldend.com` over HTTP and HTTPS trusted by its
environment's configured CA
and permanently redirect it to `https://alwaldend.com`, preserving the path
and query. The existing `www` CNAME SHALL continue to follow the apex in both
DNS views. The alias SHALL NOT select a separate website or release.

#### Scenario: Follow an existing website link

- **WHEN** a visitor opens `https://www.alwaldend.com/projects/?page=2` through either DNS view
- **THEN** the TLS connection succeeds without a browser warning on a client trusting
  the environment's CA
- **AND** the response permanently redirects to `https://alwaldend.com/projects/?page=2`

#### Scenario: Follow the alias over HTTP

- **WHEN** a visitor opens an HTTP URL on `www.alwaldend.com`
- **THEN** redirects reach the same path and query on `https://alwaldend.com`

### Requirement: Split-horizon hostname routing

`download.alwaldend.com` and `alwaldend.com` SHALL resolve to Yandex Cloud in
public DNS and to XCP-ng in local DNS. Each hostname SHALL have one canonical
record owner across both views. Public DNS SHALL NOT distribute traffic
between the two independent stores.

#### Scenario: Resolve from each view

- **WHEN** the two hostnames are resolved through public and local DNS
- **THEN** public answers identify Yandex and local answers identify XCP-ng
- **AND** neither view points to the shared ingress component

### Requirement: Public release files and directory data

The public storage layout SHALL be
`projects/<project>/releases/<version>/<uploaded files>`, with no additional
`files/` directory. Completed uploaded files, including website archives,
SHALL be anonymously downloadable. Directory requests SHALL return JSON
listings usable by the website, including names, entry types, dates, and
available file sizes. Index files SHALL NOT suppress release listings.

#### Scenario: Browse a release containing HTML and an archive

- **WHEN** a client requests a release directory containing `index.html` and an archive
- **THEN** the response lists both files as JSON
- **AND** requesting either file returns its uploaded bytes

#### Scenario: Fetch from the main site

- **WHEN** JavaScript on `alwaldend.com` requests a public listing without credentials
- **THEN** CORS permits the request and refresh observes completed uploads

#### Scenario: Resume a file download

- **WHEN** a client requests a valid byte range of an uploaded file
- **THEN** the response contains the requested bytes and appropriate range metadata

### Requirement: Selected website serving

Extracted website content SHALL reside in
`sites/<project>/releases/<version>/`. The main website SHALL serve the release
selected by `sites/alwaldend.com/current` at `https://alwaldend.com`.
Selection SHALL NOT change the public URL or remove its downloadable archive.
Only designated website archives require extraction; regular release files
SHALL be served unchanged.

#### Scenario: Select a website release

- **WHEN** successful SSH deployment selects an extracted website release
- **THEN** the main hostname serves that release's HTML and assets
- **AND** its original archive remains in the public release directory

### Requirement: Persistent content and limited runtime access

Content and site selection SHALL survive service restart, host reboot, and
Ansible configuration reruns. Yandex VM replacement SHALL retain and reattach
its content storage unless deletion is explicitly requested. The XCP-ng VM
SHALL reject automatic destruction and replacement; its replacement workflow
is deferred to a future explicitly authorized change. The publisher
SHALL have write access while Nginx SHALL have read-only content access.
Staging, credentials, and service state SHALL NOT be public browse roots.

#### Scenario: Reconfigure and reboot

- **WHEN** configuration is reapplied and the VM is rebooted
- **THEN** previously published files and the selected site remain available
- **AND** the public service exposes no staging or credential files

### Requirement: Explicit retention and environment publication

The service SHALL provide no automatic backups, replication, failover, or
release deletion. Publication targets SHALL be explicit, and a failed
publication to one environment SHALL NOT be reported as success there.

#### Scenario: Publish only locally

- **WHEN** a new version is published to the local VM
- **THEN** the cloud copy remains unchanged until explicitly published there

### Requirement: Daily content deduplication

Each environment SHALL use Btrfs content storage and run a daily background
extent-deduplication job over published project files and extracted website
releases. The job SHALL use persistent incremental hash state outside its scan
roots, run at low I/O priority with bounded CPU/memory, and wait for mounted
storage. It SHALL preserve file bytes, independent file identities, permissions,
and the selected website link. It SHALL NOT scan staging or service state,
delete releases, or merge storage across environments.

#### Scenario: Deduplicate repeated assets

- **WHEN** two releases contain independently uploaded files with identical extents
- **THEN** a successful maintenance run shares their physical extent storage
- **AND** both original paths still return their original bytes
- **AND** a later write to one file does not change the other file

#### Scenario: Repeat or interrupt maintenance

- **WHEN** a daily job runs again or is interrupted
- **THEN** published content and website selection remain intact
- **AND** another run can reuse the private hash database without overlapping an active service run

### Requirement: Administrator access to the content account

The service SHALL retain its dedicated content account without configuring
its authorized keys. Administrators SHALL use their existing SSH and sudo
access to run publication as the content account. Host configuration SHALL
install rsync and provision storage without uploading or selecting releases.

#### Scenario: Configure a host without a publisher key

- **WHEN** an administrator configures the selected download host
- **THEN** the content account and rsync are available for publication
- **AND** no separate publisher key or authorization entry is required
- **AND** existing release content and website selection remain unchanged

### Requirement: Injected operator SSH access

The component SHALL expose explicit local and Yandex SSH targets using the
shared SSH rule and its own Vault AppRole. Each invocation SHALL inject a
temporary signed client key, verify the destination hostname against the
repository-owned SSH server CA, and clean up temporary material on normal
exit. It SHALL NOT modify persistent SSH configuration or accept an unknown
host key interactively.

#### Scenario: Inspect the local deployment

- **WHEN** an authorized operator invokes the local SSH target with read-only service and storage commands
- **THEN** the command connects as the configured administrator to the local inventory hostname using injected credentials
- **AND** the Yandex host remains untouched
- **AND** the invocation removes its temporary credential and trust files after the command exits

### Requirement: Explicit local-to-public content promotion

The component SHALL expose a local-to-Yandex synchronization command using its
existing inventory and temporary Vault Ansible credentials. Local
SHALL be the authoritative source for both published content trees. The command
SHALL require mounted content roots and use the component Ansible playbook
with ansible.posix.synchronize for transfers. It SHALL NOT change DNS or deploy host configuration.

#### Scenario: Mirror locally published content

- **WHEN** an authorized operator invokes `sync.yandex`
- **THEN** the Ansible playbook mirrors the local published trees to Yandex using injected authentication and strict host trust
- **AND** dry run reports replacements and removals without changing either remote store

### Requirement: Public GitHub redirect

The Yandex Traefik configuration SHALL redirect HTTPS requests for
git.alwaldend.com permanently to https://github.com, preserving the path and
query string. Requests with go-get=1 SHALL instead reach Nginx and return
Go discovery metadata for the git.alwaldend.com/alwaldend/src repository
prefix, including requests beneath that prefix. The metadata SHALL identify
Git and https://github.com/alwaldend/src.git. Unknown discovery prefixes
SHALL return 404 without metadata. Both public routers SHALL use the existing
public ACME resolver. The local inventory SHALL omit both public Git routers.

#### Scenario: Open a public repository URL

- **WHEN** a public client requests https://git.alwaldend.com/alwaldend/src?tab=readme
- **THEN** Traefik returns a permanent redirect to https://github.com/alwaldend/src?tab=readme with a trusted certificate
- **AND** requests for download.alwaldend.com continue using the download service

#### Scenario: Discover the source repository

- **WHEN** Go requests /alwaldend/src?go-get=1 or a nested module/package path with go-get=1
- **THEN** the public endpoint returns HTTP 200 with text/html metadata naming git.alwaldend.com/alwaldend/src and https://github.com/alwaldend/src.git
- **AND** no GitHub redirect occurs

#### Scenario: Reject an unknown discovery prefix

- **WHEN** a client requests /alwaldend/src-other?go-get=1 or another repository prefix
- **THEN** the public endpoint returns 404 without go-import metadata

#### Scenario: Preserve local routing

- **WHEN** the local inventory is rendered
- **THEN** it has no GitHub redirect router or git.alwaldend.com certificate request

### Requirement: Source-authoritative synchronization

The download Ansible playbook SHALL select the local and Yandex inventory hosts.
The destination projects and sites trees SHALL match the validated source,
including changed bytes, removed releases, and current links. It SHALL exclude
service state, staging, host configuration, and credentials from synchronization.
It SHALL verify SSH host trust and keep authentication on the controller.

#### Scenario: Replace public divergence

- **WHEN** the public store has changed bytes or releases absent locally
- **THEN** sync replaces the differing content with local bytes
- **AND** removes public-only content from the two published trees

#### Scenario: Remove a local selection or all local content

- **WHEN** the source has no current link or has empty published trees
- **THEN** the destination reflects the missing link or empty trees
- **AND** private state and host configuration remain unchanged

### Requirement: Publish completed trees

The command SHALL prepare and validate changed trees privately before replacing
public trees. Each project directory switch SHALL be atomic. The command SHALL
not claim a transaction spanning projects or sites. Unchanged trees SHALL
remain untouched. Failed source capture or destination transfer SHALL leave
public trees unchanged.

#### Scenario: Update the selected website version in place

- **WHEN** the source has different bytes under an already-selected version
- **THEN** visitors see a completed site tree after the switch
- **AND** the destination version no longer retains its differing old bytes

#### Scenario: Fail during transfer

- **WHEN** a transfer fails before public tree exchange
- **THEN** the command fails and the existing published trees remain unchanged

### Requirement: Validate publication boundaries

The command SHALL reject source release links, special files, unsafe paths,
invalid current links, and missing site indexes before destination writes.
It SHALL reject symlinked destination directories and unavailable roots.
The operator SHALL be able to require mounted content roots.

#### Scenario: Reject an escaping path

- **WHEN** a source release contains a symlink or a destination directory is a symlink
- **THEN** synchronization fails without following that link during publication

### Requirement: Read-only preview

Dry run SHALL report proposed changes and deletions without remote writes.
It MAY capture source content in private controller storage and SHALL clean
that capture on completion or failure.

#### Scenario: Preview authoritative promotion

- **WHEN** an operator previews a source that lacks a public release
- **THEN** dry run reports the public release's removal
- **AND** both remote content trees remain unchanged
