# Infrastructure 3x-ui

## Purpose

Describe 3x-ui provisioning, host deployment, node routing configuration, and
the subscription URL transformer. This baseline describes checked-in source,
not verified live node access. The advertised host-only subscription command
has a source limitation: `sub/main.go` also attempts to read `sub_file` when it
is empty, so successful host-only operation is not a baseline guarantee.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [component documentation](../../../README.md),
[root wrappers](../../../BUILD.bazel),
[provisioning package](../../../tf_setup/BUILD.bazel),
[configuration package](../../../tf/BUILD.bazel),
[host playbook](../../../ansible/playbook_deploy.yaml),
[Traefik routes](../../../ansible/files/traefik_dynamic.toml),
[node definitions](../../../tf/nodes.tf),
[node inbounds](../../../tf/node/inbounds.tf),
[node routing](../../../tf/node/routing.tf), and
[subscription transformer](../../../sub/main.go).

## Requirements

### Requirement: Separate provisioning, deployment, and service configuration

The package SHALL keep infrastructure creation in `tf_setup`, host deployment
in Ansible, and 3x-ui API resources in `tf`. The deployment playbook SHALL apply
host, Traefik, and 3x-ui roles to the control-plane and data-plane inventory
groups.

#### Scenario: Select a 3x-ui workflow

- **WHEN** an authorized operator selects the setup or main Terraform wrapper
- **THEN** its AL plugin selects the corresponding Terraform stage
- **AND** the Ansible playbook remains the owner of host and service installation

### Requirement: Terminate proxy traffic at the declared local services

The Traefik route template SHALL forward the panel to loopback port 2053 and
subscription requests to loopback port 2096. It SHALL define HTTP inbound routes
for ports 40000 through 40149 and TCP inbound routes for ports 40150 through
40299, forwarding to corresponding loopback ports.

#### Scenario: Render a subscription route

- **WHEN** the checked-in dynamic Traefik template is rendered
- **THEN** requests matching the configured subscription path prefix route to
  `http://127.0.0.1:2096`
- **AND** the router selects the configured TLS certificate resolver

### Requirement: Declare explicit node and outbound routing relationships

Terraform SHALL configure the `njalla1` and `yc1` nodes through the shared node
module. The node routing declaration SHALL block the advertising category and
private IP destinations before routing the freedom inbound directly and each
`mullvad_min` or `http_proxy` inbound through its named WireGuard outbound.
The enabled `mullvad_min_lb` inbound has no active dedicated balancing rule;
commented balancing rules SHALL NOT be treated as enabled routing behavior.

#### Scenario: Inspect a configured Mullvad inbound route

- **WHEN** a node's selected relay creates a Mullvad inbound
- **THEN** the inbound listens on loopback
- **AND** its routing rule selects the matching `out-mullvad-min-` outbound tag
- **AND** the node declaration supplies its own relay selection to the module

### Requirement: Normalize subscription URLs from readable inputs

For readable local subscription files, the transformer SHALL ignore blank and
comment lines, parse each remaining URL, set `security=tls` and `sni` to the URL
hostname, sort URLs by fragment, and write one URL per output line. It SHALL
return a nonzero exit status when input reading or URL parsing fails.

#### Scenario: Transform a local subscription file

- **WHEN** `fix_subs` receives a readable `--sub_file` containing subscription URLs
- **THEN** output URL queries contain `security=tls` and the matching hostname SNI
- **AND** output is ordered by URL fragment with blank and comment lines omitted
