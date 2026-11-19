## Context

OpenHands ships three services with different responsibilities and trust
levels. This repository deploys hosts natively through Ansible with Vault ACME
certificates and shared Traefik ingress, and scopes every component to its own
Vault AppRole.

## Goals

- Deploy canvas, agent server, and automation server as separate hosts.
- Keep the agent server's execution boundary explicit and documented.
- Reuse the repository's AppRole, DNS, Terraform, Ansible, and Traefik flow.
- Reach canvas from outside without weakening the ingress mTLS policy.

## Decisions

### Run the components natively, not in containers

The repository's other components install packages, users, units, and files
through Ansible roles. Native deployment uses that service lifecycle and lets
the agent server run `--import-modules` against a plain virtual environment.
Its operating-system account permissions define its access to the host; the
workspace path is a working directory.

### Serve the packaged canvas assets directly

The canvas unit invokes the installed package's `scripts/static-server.mjs`
with Node, `--host 127.0.0.1`, the package's `build/` directory, and the declared
canvas port. The operator supplies the backend URL and session key in the
browser; the static server receives neither credential nor backend metadata.
TLS and API proxying stay with Traefik.

Inspection of the published `@openhands/agent-canvas` 1.18.0 package on
2026-09-12 found that the `--frontend-only` launcher starts a static server on
`::` at port 3001 and adds its own ingress process. The launcher does not
forward a host option. Direct use of the shipped static server supplies the
explicit loopback binding required by this deployment.

### Let the Traefik role own ingress

Component roles install and run their service on loopback; the deployment
declares the dynamic configuration that exposes it. This matches the existing
`host` plus `traefik` playbook pattern and keeps certificate and entry-point
configuration in one role.

### Proxy backend paths through the canvas origin

The operator configures the browser backend with the canvas origin. Traefik
routes `/api/automation`, `/api`, `/sockets`, and `/server_info` on that origin
to the automation and agent server services. Explicit priorities place
automation at 30, other API routes at 20, and the frontend route at 1.
This prevents the frontend's longer Host expression or the broader `/api`
route from capturing backend requests. The browser then issues no cross-origin
request, so the agent server needs no CORS allowance.

Traefik preserves the original canvas Host header. Agent-server and automation
routers therefore also match the canvas hostname; their TLS connections use
the hostname in the site-local backend URL.

### Use site-local backend targets

Backend targets use the site-local `*.openhands.alwaldend.com` names rather
than the public name. The public name resolves through `ingress.alwaldend.com`,
which requires a client certificate, and an internal hop must not depend on
that. Each backend URL already matches the hostname its certificate is issued
for, so no SNI override is needed.

### Reuse the existing public ingress contract

The [ingress specification](../../../../ingress/openspec/specs/infra-ingress/spec.md)
already requires every service entry to select a site-local backend, use the
service hostname for backend TLS, and retain verified client-certificate
access. The new `openhands-canvas` entry uses that behavior unchanged, so it
needs no ingress specification delta. OpenHands owns its DNS declarations;
ingress retains ownership of routing and client authentication.

### Populate the OpenHands resource set before provisioning

The new OpenHands entry in
[`infra/xcp_ng/tf/inventory.tf`](../../../../xcp_ng/tf/inventory.tf) supplies
pool-scoped template, storage, and network names. It uses the existing
[XO resource-set contract](../../../../xcp_ng/tf/README.md); adding this entry
requires no change to that contract.

Deployment depends first on applying the component AppRole in Vault. The XO
workflow must then synchronize its OIDC user and apply its resource set with
that exact identity and the assigned infrastructure objects. Only after those
prerequisites exist can `infra/openhands/tf_setup` provision the VMs with the
component's own credentials. A populated resource set without a synchronized
subject remains inaccessible.

Vault secret writes require an identity with write access to the component
paths before Ansible deployment. DNS records and the shared ingress service
entry must also be deployed for the declared canvas URL to work. These
prerequisites describe the dependency order; each live operation still needs
its exact operation and scope authorized under repository policy.

### Share one session API key

The agent server, the automation server, and the browser all authenticate the
same `X-Session-API-Key` value, so one secret is stored rather than three.

### Deploy the unsecured agent server with its host's owner

An agent server without a session key is appropriate only on a host that
already constrains access. The host project
(`users/simeonwarren/host_bot`) owns that host, its Traefik, and its firewall,
so it also owns the unsecured agent server. That server binds to loopback and
has no Traefik route. The OpenHands inventory contains only the three VM
components.

## Risks and trade-offs

- The agent server executes with its service account's permissions on the host.
  This is documented as the execution boundary rather than mitigated by a
  container.
- Component versions are pinned to a mutually compatible set; upgrading one
  requires checking compatibility with the others.
- The static canvas service requires the operator to configure the backend URL
  and session key in the browser. It does not advertise the runtime metadata
  needed to enable the pinned frontend's embedded VSCode UI.

## Alternatives considered

- **Container deployment.** Rejected: it adds a separate runtime lifecycle to
  the repository's native Ansible deployment.
- **Per-component reverse proxy inside the component role.** Rejected: ingress
  belongs to the Traefik role and the deployment's dynamic configuration.
- **Routing backends through the public ingress.** Rejected: internal hops would
  then require a client certificate for no added benefit.

## Continuation review and verification

Decision verdict: revise the initial routing and launch configuration. Offline
requests through three instances of pinned Traefik 3.7.1 reproduced the routing
boundaries and confirmed 13 successful requests plus an unknown-host rejection
after explicit priorities and preserved-Host matching. This covers HTTP routing
with fixture backends, not TLS, live DNS, ACME issuance, WebSocket upgrades, or
OpenHands authentication. Raw fixture results remain in ignored task scratch.

Traefik's loaded DEBUG configuration omitted the original static
`tls.stores.default.acme` block. That ineffective block was removed. Backend
routers now set explicit certificate domains for their own hostnames, so their
forwarded canvas Host matcher cannot trigger HTTP-01 issuance for a name whose
challenge resolves to another VM. Canvas certificate names remain derived from
its local router hosts. The parent `openhands_domain` inventory variable is
retained as requested. Live issuance remains unverified.

The roles now restart through handlers when installed packages, environment,
units, or compatibility modules change. Routine convergence starts a stopped
service without unconditionally restarting a running one. Canvas uses the
collection's npm module; uv package changes use its installation/removal
summary. These changes received source and syntax checks, not a live second
deployment.

The infrastructure skill consolidation incorporates the demonstrated bootstrap,
metadata-disclosure, inventory-precedence, host-package-backend, launcher,
routing, and certificate-domain lessons. Offline eval configuration loading
and an independent fixture review cover bounded cases; they do not establish
production deployment behavior.

## Live Vault operation evidence

On 2026-09-12 the user authorized the OpenHands Vault identity, SSH/PKI roles,
policies, and shared memberships, followed by the previously authorized secret
writes. A later instruction allowed unrelated changes that neither destroy
resources nor break behavior. The reviewed dependency updates only add
`secret-id-accessor/destroy` permission to the affected roles' existing
bootstrap policies; their other paths and capabilities are unchanged.

The repository Terraform targets applied three saved, reviewed plans: 16
OpenHands resource creations; one membership update adding only the OpenHands
entity to `approles`; and 12 updates comprising the three remaining shared-group
memberships, eight cleanup policies, and removal of the unused Yandex policy
from OpenHands token grants. All plans contained zero destroys. A fresh plan
for the same module and four groups reported no changes after apply. This is
bounded Vault reconciliation evidence, not a whole-Vault drift audit.

The component and user `vault.kv_put` targets accepted all three required paths
with `-cas=0` and returned version 1. Four independent values with 256 bits of
entropy were generated in process memory, supplied through JSON stdin, and
never printed or read back. Only path, field names, returned version, and
success metadata are retained as evidence. Raw plans and operation scratch
were restricted and cleaned up; they are not publication artifacts.

The first provider-schema attempt failed with a 96-byte temporary directory;
the provider socket suffix exceeded the Unix-domain pathname limit. Shortening
`TMPDIR` within the same task directory resolved startup. Saved plans replayed
successfully through fresh owning AL invocations because the Vault backend
supplies its endpoint and credentials through `TF_HTTP_*`. These verified
workflow details are maintained in `repo-infra`.

## Pull request review corrections

The first remote review identified missing automation-secret setup instructions,
an unnecessary inherited Yandex grant, a guest short hostname derived from the
component key, and duplicate requirements in the active change and main spec.
The instructions now cover all three secret paths, the AppRole opts out of the
Yandex token policy, cloud-init derives the short hostname from the selected DNS
record, and only the active delta owns the new requirements until archival.
