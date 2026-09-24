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
permit anonymous reads, redirect HTTP to HTTPS, and present certificates
trusted by normal public browsers. Serving SHALL NOT depend on shared ingress.

#### Scenario: Read without client credentials

- **WHEN** a visitor requests either hostname on the selected VM
- **THEN** the visitor can read public content without a login or client certificate
- **AND** HTTPS terminates on that VM's Traefik instance

#### Scenario: Validate the local hostname certificate

- **WHEN** local DNS resolves the hostname to XCP-ng while public DNS points to Yandex
- **THEN** the local VM can obtain and renew a trusted certificate for that hostname

### Requirement: Preserve the existing website alias

Each VM SHALL accept `www.alwaldend.com` over HTTP and publicly trusted HTTPS
and permanently redirect it to `https://alwaldend.com`, preserving the path
and query. The existing `www` CNAME SHALL continue to follow the apex in both
DNS views. The alias SHALL NOT select a separate website or release.

#### Scenario: Follow an existing website link

- **WHEN** a visitor opens `https://www.alwaldend.com/projects/?page=2` through either DNS view
- **THEN** the TLS connection succeeds without a browser warning
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
Ansible configuration reruns. Replacement of a VM SHALL retain and reattach
its content storage unless deletion is explicitly requested. The publisher
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
